package keeper

import (
"crypto/ecdsa"
"crypto/rand"
"crypto/sha256"
"database/sql"
"encoding/hex"
"encoding/json"
"fmt"
"math/big"
"os"
"strings"
"sort"
"sync"
"sync/atomic"
"time"

"github.com/ethereum/go-ethereum/accounts/abi"
"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	Type            string  `json:"type"`
	Wallet          string  `json:"wallet"`
	To              string  `json:"to,omitempty"`               // transfer destination
	Amount          float64 `json:"amount,omitempty"`
	AmountOut       float64 `json:"amount_out,omitempty"`       // swap output amount
	AmountPerHuman  float64 `json:"amount_per_human,omitempty"` // for ubi_distribution
	LPShares        float64 `json:"lp_shares,omitempty"`        // for add_liquidity
	// FromDemurrageLost/ToDemurrageLost carry the exact AEQ amount the
	// primary node decayed off Wallet/To via settleDemurrageLocked while
	// processing this TX. Secondary nodes replay these exact numbers
	// (ApplyTransferDelta/ApplySwapDelta/AddLiquidityDelta/RemoveLiquidityDelta)
	// instead of recomputing decay from effectiveBalance() at replay time —
	// recomputing would use the replaying node's own wall-clock time, which
	// can differ from the primary's by anything from network latency to a
	// full historical resync, producing a different decay amount and a
	// StateRoot divergence identical in kind to the swap-fee bug fixed in 8e3f675.
	FromDemurrageLost float64 `json:"from_demurrage_lost,omitempty"`
	ToDemurrageLost   float64 `json:"to_demurrage_lost,omitempty"`
	// DistributionAt carries the exact Unix timestamp the primary chose for
	// a distribution round (e.g. the new last_ubi_at) on
	// "ubi_distribution_finalize" TXs. Audit recheck 2 (P0 #4) found the
	// primary used to call time.Now() directly inside DistributeUBIPool
	// while secondaries replayed block.Timestamp instead — two different
	// instants, guaranteeing a StateRoot mismatch on every UBI round even
	// when every credited amount was correct. The primary now picks this
	// value once and uses it for both its own immediate state and this
	// field, so secondaries replay the IDENTICAL value instead of any
	// wall-clock reading of their own.
	DistributionAt int64 `json:"distribution_at,omitempty"`
	TxHash          string  `json:"tx_hash"`
	// Nullifier and Commitment are set on register_human TXs so secondary
	// nodes can apply the registration to their local state when they receive
	// the block — without needing a separate snapshot or state sync.
	Nullifier  string  `json:"nullifier,omitempty"`
	Commitment string  `json:"commitment,omitempty"`
	// ZK proof fields for register_human — enables secondary nodes to
	// independently verify the proof via BioVerifier without trusting
	// the validator signature alone. Fields are omitted for non-registration
	// TXs and for blocks produced by old nodes (backward-compatible).
	ProofA     []string   `json:"proof_a,omitempty"`   // [2]string big.Int decimal
	ProofB     [][]string `json:"proof_b,omitempty"`   // [2][2]string big.Int decimal
	ProofC     []string   `json:"proof_c,omitempty"`   // [2]string big.Int decimal
	PubSignals []string   `json:"pub_signals,omitempty"` // public signals (decimal)
}

type Block struct {
Height       int64    `json:"height"`
Timestamp    int64    `json:"timestamp"`
ParentHashes []string `json:"parent_hashes"`
Hash         string   `json:"hash"`
Proposer     string   `json:"proposer"`
Humans       int      `json:"humans"`
IsGenesis    bool     `json:"is_genesis,omitempty"`
	StateRoot    string   `json:"state_root,omitempty"`
	Transactions  []Transaction `json:"transactions,omitempty"`
	Signature    string   `json:"signature,omitempty"`
}

// peerChallenge holds a one-time challenge issued to a registering peer.
type peerChallenge struct {
	value     string
	expiresAt int64
}

type BlockDAG struct {
blocks                 map[string]*Block
tips                   map[string]bool
mu                     sync.RWMutex
keeper                 *Keeper
state                  *ChainState
evm                    *EVMEngine       // set by EVMRPCServer after construction; used by replayTransactions for ZK proof verification
nodeID                 string
height                 int64
// bootHeight is dag.height's value at construction time (after restoring
// it from the persisted "max_block_height" — see createGenesisBlock's
// caller), captured ONCE and never updated again. Used by
// replayTransactions to recognize "ancestor catch-up" blocks: cs.accounts
// is loaded fully from the DB at startup and already reflects every
// block up to and including bootHeight, but dag.blocks/dag.tips are
// purely in-memory and start empty on every restart — so the node must
// still fetch and insert those ancestor blocks for hash-chain/tips
// bookkeeping, WITHOUT re-applying their transactions (already accounted
// for) or comparing their claimed StateRoot against cs.accounts' current,
// much-later state (guaranteed to "mismatch" despite no real divergence).
bootHeight             int64
pendingTxs             []Transaction
txMu                   sync.Mutex
signingKey             *ecdsa.PrivateKey
authorizedValidators   map[string]bool  // Ethereum addresses allowed to propose blocks
activeSyncPeers        map[string]bool  // peers with a running syncWithNode goroutine
syncPeerMu             sync.Mutex
warnedUnknownProposers map[string]bool  // suppresses repeated "not authorized" log lines
peerChallenges         map[string]peerChallenge // address → pending challenge (P1-3)
challengeMu            sync.Mutex
replayedBlocks         map[string]bool  // tracks blocks already replayed — prevents double-credit on duplicate delivery
replayedMu             sync.Mutex
	// replayMu serializes replayTransactions calls across concurrent
	// AddPeerBlock invocations (e.g. the same or different blocks arriving
	// via P2P and HTTP sync at the same time) — replay must happen in a
	// well-defined order since TX dependencies span blocks (a register_human
	// in block N must be applied before a transfer in block N+1 from the
	// same wallet). This replaces the old single-consumer-goroutine +
	// channel design, which serialized replay the same way but ran it
	// asynchronously — see AddPeerBlock for why that was a correctness bug,
	// not just a latency tradeoff.
	replayMu sync.Mutex
stateRootMismatches map[string]int // FIX 4: per-proposer StateRoot mismatch counters
	// lastSuccessfulPeerSyncAt is the Unix timestamp of the last time this
	// node successfully accepted a peer block via AddPeerBlock. Read/written
	// with atomic.Int64 (not dag.mu) since it's set from AddPeerBlock's
	// success tail, after dag.mu has already been released — see
	// /api/health/combined (Gesamtaudit 2026-06-28, P2-4/P3-7: "Health/API
	// zeigt nicht ... seit wann [ein StateRoot-Mismatch existiert]").
	lastSuccessfulPeerSyncAt atomic.Int64
	// orphans holds blocks whose parent isn't known yet, keyed by the missing
	// parent's hash. When that parent is later added, every block waiting on
	// it is retried automatically. See AddPeerBlock for why this exists —
	// without it, a block whose parent arrived even one sync cycle late was
	// silently dropped forever, along with everything built on top of it.
	orphans   map[string][]*Block
	orphansMu sync.Mutex
	// orphanFirstSeen/orphanLastAttempt back orphan TTL + per-hash fetch
	// cooldown — see queueOrphan's "abandon" comment and fetchMissingAncestors'
	// cooldown skip for why both exist.
	orphanFirstSeen   map[string]time.Time
	orphanLastAttempt map[string]time.Time
	// orphanResolveInFlight/orphanResolveAgain coordinate triggerOrphanResolve
	// (sync_blocks.go): at most one resolution pass runs at a time, and if a
	// new orphan arrives while one is running, exactly one more pass runs
	// immediately after instead of being dropped — see triggerOrphanResolve.
	orphanResolveInFlight bool
	orphanResolveAgain    bool
	orphanResolveMu       sync.Mutex
}


// genesisTimestamp reads the genesis_time from genesis.json if present,
// falling back to the hardcoded date. P2-11: avoid hardcoded timestamp.
func genesisTimestamp() int64 {
if data, err := os.ReadFile("genesis.json"); err == nil {
var g struct { GenesisTime string `json:"genesis_time"` }
if json.Unmarshal(data, &g) == nil && g.GenesisTime != "" {
if t, err := time.Parse(time.RFC3339, g.GenesisTime); err == nil {
return t.Unix()
}
}
}
return time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC).Unix()
}

// loadAuthorizedValidators reads the AUTHORIZED_VALIDATORS env var
// (comma-separated Ethereum addresses). Used to reject peer blocks from
// unknown signers so no one can inject arbitrary blocks into the DAG.
func loadAuthorizedValidators() map[string]bool {
	m := make(map[string]bool)
	for _, addr := range strings.Split(os.Getenv("AUTHORIZED_VALIDATORS"), ",") {
		addr = strings.ToLower(strings.TrimSpace(addr))
		if strings.HasPrefix(addr, "0x") && len(addr) == 42 {
			m[addr] = true
		}
	}
	return m
}

// GetSigningKey returns the ECDSA private key used to sign blocks, or nil
// if no signing key is configured. Used by the snapshot handler to sign
// exported snapshots so peer nodes can verify their authenticity.
func (dag *BlockDAG) GetSigningKey() *ecdsa.PrivateKey {
	return dag.signingKey
}

// P1-3: Challenge-Response Validator Signature Verification ─────────────────

// IssuePeerChallenge generates a one-time challenge for a registering validator.
// The peer must sign this challenge with their signing key to prove ownership.
// Challenges expire after 90 seconds.
func (dag *BlockDAG) IssuePeerChallenge(signingAddr string) string {
	ts := time.Now().Unix()
	// P3-FIX: add 16 random bytes so two challenges issued for the same
	// address in the same second always produce different values.
	var nonce [16]byte
	rand.Read(nonce[:]) //nolint:errcheck — crypto/rand never returns an error on supported platforms
	raw := fmt.Sprintf("aequitas-validator:%s:%d:%s", strings.ToLower(signingAddr), ts, hex.EncodeToString(nonce[:]))
	h := sha256.Sum256([]byte(raw))
	challenge := hex.EncodeToString(h[:])
	dag.challengeMu.Lock()
	// Fix 7: Cap peerChallenges to prevent unbounded growth from floods of
	// challenge requests. Prune expired entries first; if still over cap, reject.
	now := time.Now().Unix()
	for addr, c := range dag.peerChallenges {
		if now > c.expiresAt {
			delete(dag.peerChallenges, addr)
		}
	}
	if len(dag.peerChallenges) > 200 {
		dag.challengeMu.Unlock()
		fmt.Printf("[DAG] ⚠ peerChallenges cap exceeded for %s — rejecting new challenge\n", strings.ToLower(signingAddr))
		return ""
	}
	dag.peerChallenges[strings.ToLower(signingAddr)] = peerChallenge{
		value:     challenge,
		expiresAt: ts + 90,
	}
	dag.challengeMu.Unlock()
	return challenge
}

// VerifyPeerChallenge verifies that signature is a valid secp256k1 signature of
// the previously issued challenge by the private key corresponding to signingAddr.
// Returns true only if: challenge exists, is not expired, and ecrecover matches.
func (dag *BlockDAG) VerifyPeerChallenge(signingAddr, signature string) bool {
	signingAddr = strings.ToLower(signingAddr)
	dag.challengeMu.Lock()
	ch, ok := dag.peerChallenges[signingAddr]
	if ok {
		delete(dag.peerChallenges, signingAddr) // one-time use
	}
	dag.challengeMu.Unlock()
	if !ok || time.Now().Unix() > ch.expiresAt {
		return false
	}
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil || len(sigBytes) != 65 {
		return false
	}
	// Ethereum signed message prefix
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(ch.value), ch.value)
	hash := crypto.Keccak256Hash([]byte(msg))
	// Normalize recovery id (v=27/28 → 0/1)
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}
	pubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false
	}
	recovered := strings.ToLower(crypto.PubkeyToAddress(*pubKey).Hex())
	return recovered == signingAddr
}

// AddAuthorizedValidator adds an Ethereum address to the set of addresses
// allowed to propose blocks. Thread-safe; safe to call after startup.
func (dag *BlockDAG) AddAuthorizedValidator(addr string) {
	addr = strings.ToLower(strings.TrimSpace(addr))
	if addr == "" {
		return
	}
	dag.mu.Lock()
	dag.authorizedValidators[addr] = true
	dag.mu.Unlock()
}

func (dag *BlockDAG) AddTransaction(tx Transaction) {
dag.txMu.Lock()
defer dag.txMu.Unlock()
dag.pendingTxs = append(dag.pendingTxs, tx)
}

func NewBlockchain(keeper *Keeper, nodeID string, state *ChainState) *BlockDAG {
dag := &BlockDAG{
blocks:                 make(map[string]*Block),
tips:                   make(map[string]bool),
keeper:                 keeper,
state:                  state,
nodeID:                 nodeID,
authorizedValidators:   loadAuthorizedValidators(),
activeSyncPeers:        make(map[string]bool),
warnedUnknownProposers: make(map[string]bool),
peerChallenges:         make(map[string]peerChallenge),
replayedBlocks:         make(map[string]bool),
	stateRootMismatches:    make(map[string]int),
	orphans:                make(map[string][]*Block),
	orphanFirstSeen:        make(map[string]time.Time),
	orphanLastAttempt:      make(map[string]time.Time),
}
if pk := strings.TrimPrefix(os.Getenv("RELAYER_PRIVATE_KEY"), "0x"); pk != "" {
	if key, err := crypto.HexToECDSA(pk); err == nil {
		dag.signingKey = key
		// Always authorize ourselves — derived from the signing key, not the nodeID.
		selfAddr := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
		dag.authorizedValidators[selfAddr] = true
		fmt.Printf("✓ Block signing enabled (RELAYER_PRIVATE_KEY loaded), proposer addr: %s\n", selfAddr)
	} else {
		fmt.Printf("[BLOCK] Warning: RELAYER_PRIVATE_KEY invalid, blocks will be unsigned: %v\n", err)
	}
} else {
	fmt.Println("[BLOCK] ⚠ RELAYER_PRIVATE_KEY not set — blocks will be unsigned. Peer nodes will reject unsigned blocks.")
}
dag.createGenesisBlock()

// FIX (audit 2026-06-28 recheck 5, P1-2): recover any pending_txs row left
// "included" by a process that crashed before its block ever reached
// BroadcastBlock — see ResetStaleIncludedPendingTxs' own comment for why
// that's always safe to retry. 10 minutes comfortably exceeds how long a
// single ProduceBlock call could ever legitimately take.
state.ResetStaleIncludedPendingTxs(10 * time.Minute)

// FIX (audit 2026-06-28 full recheck, P1-3): restore every durably-saved
// block (see chain_blocks' own comment and SaveBlockToDB) BEFORE falling
// back to the bare max_block_height counter below. This is what lets a
// node recover its own previously produced/accepted blocks — and their
// full transaction lists — across a restart without needing any peer to
// still have them; the counter-only fallback further down only recovers
// the height NUMBER, not the actual block data.
loaded, loadErr := state.LoadBlocksFromDB()
if loadErr != nil {
	// FIX (2026-06-28, production incident): a transient DB error here
	// used to be silently treated as "this node has zero durably-saved
	// blocks" — for a node with a full chain_blocks table, that meant
	// starting fresh at genesis and forcing a complete peer resync of its
	// own history, repeatedly, on every restart that hit the hiccup (see
	// LoadBlocksFromDB's own comment). Crashing here and letting the
	// process supervisor (Docker --restart unless-stopped / Railway) retry
	// the whole startup is safer than silently continuing with a DAG that
	// doesn't reflect this node's real history.
	fmt.Printf("[BLOCK] ✗ FATAL: could not restore blocks from chain_blocks: %v — exiting so the process supervisor restarts cleanly instead of starting with a falsely-empty DAG\n", loadErr)
	os.Exit(1)
}
if len(loaded) > 0 {
	referenced := make(map[string]bool, len(loaded))
	for _, b := range loaded {
		dag.blocks[b.Hash] = b
		// Already reflected in chain_accounts (committed when these TXs
		// were first applied, before this block was even assembled) —
		// must not be re-applied by replayTransactions.
		dag.replayedBlocks[b.Hash] = true
		for _, ph := range b.ParentHashes {
			referenced[ph] = true
		}
		if b.Height > dag.height {
			dag.height = b.Height
		}
	}
	for hash := range dag.tips {
		if referenced[hash] {
			delete(dag.tips, hash)
		}
	}
	for hash := range loaded {
		if !referenced[hash] {
			dag.tips[hash] = true
		}
	}
	fmt.Printf("[BLOCK] Restored %d durable block(s) from chain_blocks — height=%d, tips=%d\n", len(loaded), dag.height, len(dag.tips))
}

// FIX (double-apply): dag.height/dag.blocks/dag.tips used to be purely
// in-memory — ReconstructState is a no-op when using Postgres, so they
// reset to genesis on every process restart regardless of how much
// chain history cs.accounts (loaded fresh from the DB above) actually
// reflects. This counter-only fallback covers any block produced before
// chain_blocks existed (or saved by a node that hadn't yet picked up
// this fix): it can only raise dag.height, never lower what the loaded
// blocks above already established, so ExportSnapshot reports the
// chain's true cumulative height, not "blocks observed since this
// process last started" — see the same fix's writes in
// ProduceBlock/AddPeerBlock and StateSnapshot.Height's comment for the
// bug this caused (a fresh-bootstrapped secondary's snapshot cutoff was
// reported far too low, so it still re-replayed — and double-applied —
// every block between the true height and the process-local one).
// FIX (audit 2026-06-28 recheck 4, P0-1): startup code, no lock held —
// must use the plain DB-only read.
if persisted := state.getConfigValueDB("max_block_height"); persisted != "" {
	var h int64
	fmt.Sscanf(persisted, "%d", &h)
	if h > dag.height {
		dag.height = h
	}
}
// Captured ONCE, after the restoration above and before any block
// processing begins — see bootHeight's field comment.
dag.bootHeight = dag.height
return dag
}

// RefreshBootHeightAfterSnapshotImport re-reads max_block_height/
// snapshot_import_height from the DB and raises dag.height/dag.bootHeight
// to match, if higher than what NewBlockchain captured at construction time.
//
// FIX (root cause behind Contabo VPS's permanent post-resync catch-up
// failure, found 2026-06-28): main.go constructs the BlockDAG (which seeds
// dag.height/dag.bootHeight from whatever max_block_height already exists
// in the DB) BEFORE RESYNC_FROM_SNAPSHOT/BOOTSTRAP_SNAPSHOT_URL ever runs.
// On a freshly wiped DB, that means dag.height/dag.bootHeight are captured
// as 0 — and bootHeight, not just height, matters here: replayTransactions'
// skipHeight check (see its own comment) takes max(dag.bootHeight,
// snapshot_import_height read live from DB), so in principle the live DB
// read alone should have been enough. In practice dag.height itself (the
// sync frontier doSyncOnce pages forward from) stayed frozen at 0, so the
// node still had to fetch, hash-verify, and insert into dag.blocks/dag.tips
// every single one of ~18,000 historical blocks one HTTP page at a time
// before reaching its true frontier — needless work that starved
// fetchMissingAncestors of cycles while validators kept producing new
// blocks every ~6s, which is what caused the orphan buffer to fall behind
// and start permanently abandoning blocks it could have resolved given
// less contention (see orphanAbandonAfter's comment). Calling this right
// after a successful snapshot import/resync, before HTTP sync starts, lets
// the node begin paging from near its true height immediately — the only
// blocks it then needs from peers are the handful actually referenced as
// parents going forward, fetched on demand via fetchMissingAncestors,
// never the full historical backlog.
func (dag *BlockDAG) RefreshBootHeightAfterSnapshotImport() {
	dag.mu.Lock()
	defer dag.mu.Unlock()

	// bootHeight = max(max_block_height, snapshot_import_height): controls
	// replayTransactions' skipHeight so we never re-apply state that the
	// snapshot already encodes.
	var bootH int64
	if persisted := dag.state.getConfigValueDB("max_block_height"); persisted != "" {
		fmt.Sscanf(persisted, "%d", &bootH)
	}
	if snapHeightStr := dag.state.getConfigValueDB("snapshot_import_height"); snapHeightStr != "" {
		var snapHeight int64
		fmt.Sscanf(snapHeightStr, "%d", &snapHeight)
		if snapHeight > bootH {
			bootH = snapHeight
		}
	}
	if bootH > dag.bootHeight {
		dag.bootHeight = bootH
	}

	// dag.height = max_block_height ONLY — this is the sync frontier
	// doSyncOnce pages forward from. After a snapshot resync, max_block_height
	// is reset to 0 (see ResyncFromSnapshotURL) so the node re-downloads all
	// block headers sequentially from genesis. Raising dag.height here from
	// snapshot_import_height would cause doSyncOnce to start near the snapshot
	// height, where dag.blocks is empty (chain_blocks was cleared), making
	// every incoming block orphan on a missing parent permanently.
	var maxH int64
	if persisted := dag.state.getConfigValueDB("max_block_height"); persisted != "" {
		fmt.Sscanf(persisted, "%d", &maxH)
	}
	if maxH > dag.height {
		dag.height = maxH
	}
}

func (dag *BlockDAG) createGenesisBlock() {
genesis := &Block{
Height:       0,
Timestamp:    genesisTimestamp(), // P2-11: reads from genesis.json when available
ParentHashes: []string{},
Proposer:     "genesis",
// FIX: this used to be dag.state.TotalHumans() — i.e. however many humans
// THIS node's own DB currently has loaded at the moment it happens to
// start up. calculateHash() includes Humans in the hashed fields, so two
// nodes starting at different points in registration history (e.g. one
// freshly reset to 0 humans, another restarted after a registration
// already succeeded) computed two DIFFERENT genesis hashes. Since
// AddPeerBlock only removes a parent from dag.tips on an EXACT hash
// match, a secondary's own (differently-hashed) genesis tip was never
// removed when the primary's block #1 arrived — referencing the
// PRIMARY's genesis hash as its parent, not the secondary's. The
// secondary's orphaned genesis then sat in dag.tips forever (nothing
// ever referenced it as a parent to remove it), permanently showing
// "Tips: 2" with no merge ever happening — confirmed in production.
// Genesis must be 100% deterministic across every node by definition
// (it's the one block everyone is supposed to agree on without any
// data exchange), so it can never depend on a node's own live state.
Humans:       0,
IsGenesis:    true,
}
genesis.Hash = dag.calculateHash(genesis)
dag.blocks[genesis.Hash] = genesis
dag.tips[genesis.Hash] = true
dag.height = 0
fmt.Printf("✓ Genesis Block (DAG): %s\n", genesis.Hash[:16]+"...")
}

func (dag *BlockDAG) calculateHash(b *Block) string {
// Normalize nil to empty slice so JSON always produces "[]" not "null".
// omitempty on the Transactions field strips the key during HTTP transport,
// and the receiver deserialises to nil — without this normalisation the
// tx_root differs between producer and receiver, causing hash mismatches.
txs := b.Transactions
if txs == nil {
txs = []Transaction{}
}
txData, _ := json.Marshal(txs)
txRootBytes := sha256.Sum256(txData)
txRoot := hex.EncodeToString(txRootBytes[:])
// Use parent hashes in the order stored on the block — do NOT sort here.
// Sorting must happen when PRODUCING a block (in ProduceBlock) so the order
// is baked into block.ParentHashes before the hash is computed. Re-sorting
// during verification would break hashes for blocks produced by peers using
// the original order.
data, _ := json.Marshal(map[string]interface{}{
"height":        b.Height,
"timestamp":     b.Timestamp,
"parent_hashes": b.ParentHashes,
"proposer":      b.Proposer,
"humans":        b.Humans,
"state_root":    b.StateRoot,
"tx_root":       txRoot,
})
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:])
}

func (dag *BlockDAG) ProduceBlock() *Block {
dag.mu.Lock()
defer dag.mu.Unlock()

// Collect all current tips as parents.
// Sort deterministically so the hash is identical regardless of map
// iteration order — both nodes must agree on parent_hashes ordering.
parentHashes := make([]string, 0, len(dag.tips))
for hash := range dag.tips {
parentHashes = append(parentHashes, hash)
}
sort.Strings(parentHashes)

// Height = max parent height + 1
maxParentHeight := int64(0)
for _, ph := range parentHashes {
if parent, ok := dag.blocks[ph]; ok {
if parent.Height > maxParentHeight {
maxParentHeight = parent.Height
}
}
}

dag.txMu.Lock()
txs := make([]Transaction, len(dag.pendingTxs))
copy(txs, dag.pendingTxs)
dag.pendingTxs = nil
dag.txMu.Unlock()

// Drain DB-persisted pending TXs — these survived a node restart and
// must now be included in a block so secondary nodes receive them.
// Without this, a transfer applied just before a restart would never
// reach secondary nodes and balances would diverge permanently.
var pendingTxIDs []int64
if dag.state != nil {
	dbTxs, ids := dag.state.LoadPendingTxs()
	if len(dbTxs) > 0 {
		fmt.Printf("[DAG] Including %d restart-surviving TX(s) from DB in block\n", len(dbTxs))
		txs = append(txs, dbTxs...)
		pendingTxIDs = ids
	}
}

proposer := dag.nodeID
if dag.signingKey != nil {
	// Use the Ethereum address derived from the signing key so peer nodes
	// can verify the block signature against a known Ethereum address.
	// The libp2p nodeID is used for network routing; the signing address
	// is what peers need for consensus verification.
	proposer = crypto.PubkeyToAddress(dag.signingKey.PublicKey).Hex()
}
block := &Block{
Height:       maxParentHeight + 1,
Timestamp:    time.Now().Unix(),
ParentHashes: parentHashes,
Proposer:     proposer,
Humans:       dag.state.TotalHumans(),
Transactions: txs,
StateRoot:    dag.state.StateRoot(),
}
block.Hash = dag.calculateHash(block)
if dag.signingKey != nil {
	hashBytes := common.HexToHash(block.Hash)
	if sig, err := crypto.Sign(hashBytes[:], dag.signingKey); err == nil {
		block.Signature = hex.EncodeToString(sig)
	} else {
		fmt.Printf("[BLOCK] Warning: could not sign block #%d: %v\n", block.Height, err)
	}
}

dag.blocks[block.Hash] = block

// FIX (audit 2026-06-28 full recheck, P1-3): "durably stored" used to mean
// only "inserted into dag.blocks" — a Go map that resets to genesis on
// every restart, not actual durable storage. Persist the block header to
// chain_blocks BEFORE clearing the outbox: if this save fails, skip the
// clear below so the TXs survive in pending_txs for the next ProduceBlock
// attempt to re-include, rather than disappearing from every durable
// record at once. A save failure does NOT stop block production itself —
// blocking the whole chain on a transient DB hiccup would be worse than
// the narrow durability gap this closes — but it does mean this round's
// TXs stay safely re-includable instead of being silently dropped.
// FIX (AQT-NEW-P1-01): stamp included_block_hash BEFORE SaveBlockToDB so
// ResetStaleIncludedPendingTxs can always distinguish the two crash windows:
//   • crash after mark but before save → block absent from chain_blocks →
//     Reset's NOT EXISTS check fires → rows requeued → correct
//   • crash after save but before clear → block present in chain_blocks →
//     Reset leaves rows alone → TXs already confirmed → orphan rows harmless
// Previous order (mark AFTER save) left a window where rows had
// included_block_hash=NULL but the block was already saved; Reset's
// "included_block_hash IS NULL" arm would requeue them → double-inclusion.
if len(pendingTxIDs) > 0 {
	// FIX (BRUTAL-P2-04): MarkPendingTxsIncluded now returns error.
	// A failure here means rows keep included_block_hash=NULL even though
	// the block is about to be saved — ResetStaleIncludedPendingTxs's
	// "block absent from chain_blocks" arm would requeue them, risking
	// double-inclusion on the next ProduceBlock. Mark node degraded so
	// operators can see it; block production continues to avoid halting
	// the whole chain over a transient DB write.
	if err := dag.state.MarkPendingTxsIncluded(pendingTxIDs, block.Hash); err != nil {
		dag.state.SetBootstrapDegraded(fmt.Sprintf(
			"MarkPendingTxsIncluded failed for block %s: %v — pending TXs may be re-included; check /api/health/combined",
			block.Hash[:16], err))
	}
}

blockSaved := true
if err := dag.state.SaveBlockToDB(block); err != nil {
	blockSaved = false
	fmt.Printf("[BLOCK] ⚠ Could not persist block #%d (%s...) to chain_blocks: %v — outbox rows kept for retry\n", block.Height, block.Hash[:16], err)
}

if blockSaved && len(pendingTxIDs) > 0 {
	if err := dag.state.ClearPendingTxs(pendingTxIDs); err != nil {
		fmt.Printf("[BLOCK] ⚠ ALERT: outbox rows for block #%d could not be cleared — these TX(s) may be duplicated into a future block: %v\n", block.Height, err)
	}
}

// Record that this proposer produced a block — used for proportional
// validator-reward distribution in DistributeValidatorsPool.
//
// FIX (audit recheck3, P2): used to fire via `go` — if this process died
// right after producing the block (before the goroutine's single DB
// UPDATE ran), that block silently never counted toward this validator's
// own reward weight, with no error anywhere to reveal it. Synchronous now;
// it's one UPDATE statement, the same cost ProduceBlock already pays for
// setConfigValue("max_block_height", ...) a few lines below.
dag.state.IncrementBlockCount(proposer)

// Remove all parents from tips, add this block as new tip
for _, ph := range parentHashes {
delete(dag.tips, ph)
}
dag.tips[block.Hash] = true
dag.height = block.Height
// FIX (double-apply): persist so a restart can resume from the true
// cumulative height instead of dag.height resetting to 0 — see
// createGenesisBlock's restoration of this value and the comment on
// StateSnapshot.Height for why an in-memory-only height broke snapshot
// bootstrap.
dag.state.setConfigValue("max_block_height", fmt.Sprintf("%d", dag.height))

if len(parentHashes) > 1 {
fmt.Printf("[DAG] 🔀 Merged %d tips into block #%d\n", len(parentHashes), block.Height)
}

return block
}

// WithBlockProductionPaused runs fn while holding the same lock
// ProduceBlock takes for its entire body (tip/parent selection, pending-TX
// drain, and the final dag.state.StateRoot() read are all done under
// dag.mu — see ProduceBlock above).
//
// FIX (audit recheck 2, P0 #2): daily distribution (main.go) mutates state
// across several separate calls (DistributeUBIPool, then
// DistributeValidatorsPool, then DistributeLPPool, then escrow) and only
// persists the TX explaining each mutation a moment later via SavePendingTx
// — without this guard, ProduceBlock's 6-second ticker could fire in the
// gap between a mutation and its TX, assembling a block whose StateRoot
// already reflects the mutation but whose Transactions list doesn't yet
// include the TX that explains it. No other node could ever reproduce
// that StateRoot by replaying that block. Wrapping the entire distribution
// round (every mutation AND every corresponding SavePendingTx call) in
// this guard makes ProduceBlock block until the round finishes, the same
// way it already serializes against AddPeerBlock's replay via replayMu.
func (dag *BlockDAG) WithBlockProductionPaused(fn func()) {
	dag.mu.Lock()
	defer dag.mu.Unlock()
	fn()
}

// maxOrphans caps total queued orphan blocks across all missing-parent keys,
// so a malicious or buggy peer sending blocks that reference parents which
// will never arrive can't grow this map without bound.
//
// FIX: confirmed in production at 2000 — a node that fell significantly
// behind (multiple validators producing every ~6s while it was still
// catching up on a large historical gap) overflowed this buffer, silently
// DROPPING individual blocks with no record of which hash was missing.
// Once dropped, no mechanism (not even fetchMissingAncestors, which walks
// back from recorded missing-parent hashes) can ever learn to re-fetch that
// specific block — if the BlockDAG's multi-parent tolerance lets later tips
// route around the gap via a sibling branch instead, that block's
// transactions are gone from this node's view forever, a real, confirmed
// divergence (a transfer present on two other nodes never landed on the
// one that overflowed). Raised by 25x to make this far less likely to
// trigger under the same catch-up load; does not fix the underlying
// lossy-on-overflow design (tracked separately — recovery today is a full
// resync from a signed snapshot, see ImportSnapshotFromURL).
const maxOrphans = 50000

// orphanAbandonAfter bounds how long this node will keep trying to resolve a
// single missing-parent hash before giving up on it for good.
//
// FIX (the real completion of the orphan saga, not another mitigation): the
// eager triggerOrphanResolve above fixed the *original* bug (an orphan
// sitting idle for up to 6s before the next retry) but, confirmed live on
// the VPS and CD20 immediately after deploying it, exposed a second one —
// every NEW orphan re-triggers a resolution pass over EVERY pending
// missing-parent hash, including ones that have been failing for minutes
// because the block genuinely no longer exists on ANY reachable peer (it
// was lost during the original pre-fix orphan-overflow incident: confirmed
// by checking the VPS itself, which produced/relayed that exact chain, and
// it doesn't have the ancestor either). With validators producing a new
// block every ~6s, that's a fresh full sweep — including an HTTP request to
// every peer for the dead hash — roughly 10x a minute, forever, compounding
// across however many nodes are simultaneously stuck on the same dead
// branch. That retry storm is what was timing out CD20→VPS requests, not a
// network failure: confirmed by checking account balances on every node
// (they match exactly, 1134/866 AEQ, total_supply 2000 everywhere) — the
// stuck branch is provably an empty/no-value side-chain, not a transaction
// any node still needs. Past this timeout, stop retrying a specific hash
// and drop everything waiting on it, freeing the memory and ending the
// storm, instead of retrying something proven unfetchable forever.
const orphanAbandonAfter = 3 * time.Minute

// orphanFetchCooldown is the minimum gap between fetch attempts for the same
// missing-parent hash, checked by fetchMissingAncestors (sync_blocks.go).
// Without it, every new orphan's triggerOrphanResolve pass re-attempts every
// OTHER still-pending hash too, even ones whose last attempt was a second
// ago — multiplying request volume by however often new orphans arrive.
const orphanFetchCooldown = 10 * time.Second

// queueOrphan stores block, which is waiting on missingParent to appear,
// and logs the wait (the old code dropped this case with zero logging).
func (dag *BlockDAG) queueOrphan(missingParent string, block *Block) {
	dag.orphansMu.Lock()
	now := time.Now()
	if first, ok := dag.orphanFirstSeen[missingParent]; ok {
		if now.Sub(first) > orphanAbandonAfter {
			abandoned := len(dag.orphans[missingParent]) + 1 // + this block
			delete(dag.orphans, missingParent)
			delete(dag.orphanFirstSeen, missingParent)
			delete(dag.orphanLastAttempt, missingParent)
			dag.orphansMu.Unlock()
			// FIX (2026-06-28): downgraded from a skull-emoji "Abandoning" line —
			// new node operators reading deploy logs during catch-up read this
			// as their node being broken. It isn't: this is a dead-end sibling
			// block from concurrent multi-validator production, not present on
			// any currently-synced peer, with zero effect on account state.
			fmt.Printf("[DAG] (housekeeping) discarded %d dead-end sibling block(s) for missing parent %s... — normal during catch-up, no peer still has it, no effect on account balances\n",
				abandoned, missingParent[:min(16, len(missingParent))])
			return
		}
	} else {
		dag.orphanFirstSeen[missingParent] = now
	}
	// FIX (audit recheck2, P2 #11): the same block can legitimately arrive
	// more than once before its parent shows up — once via P2P broadcast,
	// again via an HTTP-SYNC page that happens to cover the same height, and
	// again on every syncWithNode retry while the gap persists. Without a
	// hash check here, each delivery appended a fresh duplicate entry,
	// burning through maxOrphans' budget on copies of a block already
	// waiting rather than genuinely distinct blocks, and inflating the
	// "N block(s) now waiting" counts in the logs above their real value.
	for _, b := range dag.orphans[missingParent] {
		if b.Hash == block.Hash {
			dag.orphansMu.Unlock()
			return
		}
	}
	total := 0
	for _, v := range dag.orphans {
		total += len(v)
	}
	if total >= maxOrphans {
		dag.orphansMu.Unlock()
		fmt.Printf("[DAG] ✗ Dropped peer block #%d: orphan buffer full (%d waiting), missing parent never arrived\n",
			block.Height, total)
		return
	}
	dag.orphans[missingParent] = append(dag.orphans[missingParent], block)
	waitingCount := len(dag.orphans[missingParent])
	dag.orphansMu.Unlock()
	fmt.Printf("[DAG] ⏳ Block #%d from %s queued as orphan — missing parent %s... (%d block(s) now waiting on it)\n",
		block.Height, block.Proposer, missingParent[:min(16, len(missingParent))], waitingCount)

	// FIX (the actual completion of the orphan-buffer mitigation, not just a
	// bigger cap): before this, a newly queued orphan sat untouched until the
	// next periodic syncWithNode tick — up to 6s later, PER peer, and only
	// one peer's fetchMissingAncestors ran per tick. Under sustained load
	// (multiple validators producing every ~6s while a node is still deep in
	// catch-up) new orphans can arrive faster than that cadence drains them,
	// which is exactly how the buffer reached its cap in production. Kicking
	// off resolution immediately, against every currently-syncing peer in
	// parallel, the instant a gap is detected — instead of waiting for the
	// next tick — closes that race instead of just buying more headroom
	// before it recurs.
	go dag.triggerOrphanResolve()
}

// popOrphans returns and removes every block that was waiting on parentHash.
func (dag *BlockDAG) popOrphans(parentHash string) []*Block {
	dag.orphansMu.Lock()
	defer dag.orphansMu.Unlock()
	waiting := dag.orphans[parentHash]
	delete(dag.orphans, parentHash)
	delete(dag.orphanFirstSeen, parentHash)
	delete(dag.orphanLastAttempt, parentHash)
	return waiting
}

// MissingParentHashes returns a snapshot of every hash currently blocking at
// least one queued orphan. Used by fetchMissingAncestors (sync_blocks.go) to
// know exactly which specific ancestor blocks to fetch by hash.
func (dag *BlockDAG) MissingParentHashes() []string {
	dag.orphansMu.Lock()
	defer dag.orphansMu.Unlock()
	hashes := make([]string, 0, len(dag.orphans))
	for h := range dag.orphans {
		hashes = append(hashes, h)
	}
	return hashes
}

// shouldAttemptFetch reports whether enough time has passed since the last
// fetch attempt for hash to try again, and records this attempt if so. See
// orphanFetchCooldown for why this exists — without it, every new orphan's
// resolve pass re-hits every other pending hash regardless of how recently
// it was last tried.
func (dag *BlockDAG) shouldAttemptFetch(hash string) bool {
	dag.orphansMu.Lock()
	defer dag.orphansMu.Unlock()
	if last, ok := dag.orphanLastAttempt[hash]; ok && time.Since(last) < orphanFetchCooldown {
		return false
	}
	dag.orphanLastAttempt[hash] = time.Now()
	return true
}

func (dag *BlockDAG) AddPeerBlock(block *Block) bool {
dag.mu.Lock()
// NOTE: no defer — we manually unlock before the channel send below (Fix 2).
// All early-return paths must call dag.mu.Unlock() explicitly.

// Skip if already known
if _, exists := dag.blocks[block.Hash]; exists {
dag.mu.Unlock()
return false
}

// Genesis blocks are always created locally — never accept from peers.
// A peer could send any block with IsGenesis=true and it would bypass
// both the signature check and the parent check below.
if block.IsGenesis {
fmt.Printf("[DAG] ✗ Rejected peer genesis: genesis can only be created locally\n")
dag.mu.Unlock()
return false
}

// Integrity check 1: recompute hash from block fields.
expectedHash := dag.calculateHash(block)
if expectedHash != block.Hash {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: hash mismatch (claimed %s..., computed %s...)\n",
block.Height, block.Hash[:min(16, len(block.Hash))], expectedHash[:16])
dag.mu.Unlock()
return false
}

// Integrity check 2: all non-genesis blocks must carry a valid ECDSA
// signature from the proposer. Unsigned blocks are rejected — this is the
// primary consensus enforcement mechanism.
if !block.IsGenesis && block.Signature == "" {
	fmt.Printf("[DAG] ✗ Rejected peer block #%d from %s: missing signature\n",
		block.Height, block.Proposer)
	dag.mu.Unlock()
	return false
}
if block.Signature != "" && !block.IsGenesis {
	sigBytes, sigErr := hex.DecodeString(block.Signature)
	if sigErr != nil || len(sigBytes) != 65 {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: malformed signature\n", block.Height)
		dag.mu.Unlock()
		return false
	}
	hashBytes := common.HexToHash(block.Hash)
	pubkeyBytes, recErr := crypto.Ecrecover(hashBytes[:], sigBytes)
	if recErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: signature recovery failed: %v\n", block.Height, recErr)
		dag.mu.Unlock()
		return false
	}
	pubkey, parseErr := crypto.UnmarshalPubkey(pubkeyBytes)
	if parseErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: invalid public key: %v\n", block.Height, parseErr)
		dag.mu.Unlock()
		return false
	}
	recoveredAddr := strings.ToLower(crypto.PubkeyToAddress(*pubkey).Hex())
	proposer := strings.ToLower(block.Proposer)
	// Proposer must be the Ethereum address that produced the signature.
	// Blocks where the proposer field does not match the recovered signing
	// address are unconditionally rejected — no libp2p-nodeID exemption.
	if recoveredAddr != proposer {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: signature mismatch (signer %s, proposer %s)\n",
			block.Height, recoveredAddr, proposer)
		dag.mu.Unlock()
		return false
	}
	// Proposer must be in the authorized validator set. Without this check
	// anyone can generate an Ethereum key, sign a block, and feed it in.
	if !dag.authorizedValidators[proposer] {
		// P3-2: cap to prevent unbounded memory growth from forged proposer addresses
		if len(dag.warnedUnknownProposers) > 500 {
			dag.warnedUnknownProposers = make(map[string]bool)
		}
		if !dag.warnedUnknownProposers[proposer] {
			dag.warnedUnknownProposers[proposer] = true
			fmt.Printf("[DAG] ✗ Proposer %s is not an authorized validator — add to AUTHORIZED_VALIDATORS env var to accept its blocks\n", proposer)
		}
		dag.mu.Unlock()
		return false
	}
}

// Integrity check 3: parent-existence and height validation.
//
// FIX (orphan buffer): this used to tolerate a missing parent only while
// len(dag.blocks) <= 10 — i.e. only during the first ~minute of a fresh
// node's life, since every node produces its own block every 6s regardless
// of sync status. Past that point, ANY block whose parent wasn't already
// in dag.blocks was silently dropped with NO log line (this branch had no
// fmt.Printf, unlike every other reject path in this function) and NEVER
// retried — and everything built on top of that block inherited the same
// fate, since ITS parent (the dropped block) would also never exist
// locally. Confirmed in production with 3 concurrent validators: every
// node's own /api/blocks ended up showing ONLY its own single-parent
// chain, never the other validators' blocks, because somewhere in their
// ancestry a single missing parent (a brief P2P gap, a sync page that
// didn't cover it, anything transient) permanently blocked the entire
// subtree above it — with no error anywhere to even reveal why.
//
// Now: a block with a missing parent is queued in dag.orphans, keyed by
// the missing hash, instead of being dropped. When that parent later
// arrives (via AddPeerBlock, below), every block waiting on it is
// automatically retried — and if THAT retry succeeds, its own dependents
// get retried too, recursively. A transient gap now costs one retry
// instead of permanently orphaning an entire branch.
if len(block.ParentHashes) == 0 {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: no parent hashes\n", block.Height)
dag.mu.Unlock()
return false
}
if block.Height > 1 {
maxParentHeight := int64(-1)
missingParent := ""
for _, ph := range block.ParentHashes {
parent, parentExists := dag.blocks[ph]
if !parentExists {
	missingParent = ph
	break
}
if parent.Height > maxParentHeight {
maxParentHeight = parent.Height
}
}
if missingParent != "" {
	dag.mu.Unlock()
	dag.queueOrphan(missingParent, block)
	return false
}
if maxParentHeight >= 0 && block.Height != maxParentHeight+1 {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: invalid height (parent max %d)\n",
block.Height, maxParentHeight)
dag.mu.Unlock()
return false
}
}

// Integrity check 4: transaction type whitelist — unknown types could
// inject unrecognised state-change commands into the audit log.
for _, tx := range block.Transactions {
switch tx.Type {
case "", "register_human", "transfer", "swap_aeq_tusd", "swap_tusd_aeq", "add_liquidity", "remove_liquidity", "faucet", "ubi_distribution", "ubi_distribution_finalize",
	"validator_distribution", "validator_distribution_pool_zero", "lp_distribution", "lp_distribution_pool_zero", "escrow_move", "escrow_release":
// known / empty — OK
default:
fmt.Printf("[DAG] ✗ Rejected peer block #%d: unknown tx type %q\n", block.Height, tx.Type)
dag.mu.Unlock()
return false
}
}

// Structural validation passed. Release dag.mu before replay — replay
// uses dag.state's own lock (cs.mu), not dag.mu, and must never run while
// holding dag.mu (ProduceBlock and other dag.mu users would block for the
// duration of every peer block's replay otherwise).
dag.mu.Unlock()

// FIX (the actual BlockDAG correctness bug, not just a hardening pass):
// this used to (1) compare block.StateRoot against dag.state.StateRoot()
// BEFORE replaying this block's own transactions, and (2) insert the block
// into dag.blocks/dag.tips unconditionally, with replay only queued
// asynchronously onto a channel that could silently drop it if full.
//
// (1) is not just "risky" — it is structurally wrong and guaranteed to
// "mismatch" on every single block that contains any transaction, even
// between two perfectly healthy, fully-synced nodes: block.StateRoot is
// computed by the PRODUCER *after* applying this block's own TXs (the
// producer's RPC handlers apply state changes synchronously before queuing
// them for inclusion — see evm_rpc.go/register.go/swap.go), so it is a
// POST-state root. Comparing it against the RECEIVER's StateRoot at this
// point — before the receiver has replayed this block's TXs — compares a
// post-state against a pre-state. That's why "[DAG] StateRoot mismatch ...
// accepted (warn only)" fired constantly throughout this project's history
// on nearly every non-empty block: the check could never have detected real
// divergence, it was comparing the wrong two snapshots by construction.
//
// (2) meant a block could be permanently "in the DAG" (counted in height,
// returned by /api/blocks, used as a valid parent for later blocks) before
// its own state changes were verified to apply cleanly, or even applied at
// all if the replay queue happened to be full.
//
// Fixed by replaying SYNCHRONOUSLY, right here, before the block is
// inserted anywhere — and only THEN comparing StateRoot, now correctly
// post-state vs. post-state. replayMu serializes this across concurrent
// AddPeerBlock calls (same ordering guarantee the old channel+goroutine
// provided, without the "silently drop if busy" failure mode: this blocks
// instead of dropping, and replayTransactions' own dedup guard makes that
// safe even under concurrent delivery of the same block).
dag.replayMu.Lock()
replayOK := dag.replayTransactions(block)
dag.replayMu.Unlock()

// FIX (block-level atomicity): replayTransactions now rolls back and
// returns false if any of this block's transactions hit a genuine
// state-inconsistency failure (not an expected idempotent skip like
// "already registered"), OR if the post-replay StateRoot doesn't match
// the producer's claimed root (audit recheck 2, P0 #1 — moved into
// replayTransactions itself so a mismatch can use that function's own
// rollback snapshot; see its comment). Treat either exactly like any
// other validation failure: the block is rejected outright, never
// inserted into dag.blocks/dag.tips. A later sync cycle or orphan-retry
// may succeed once local state has caught up with whatever caused the
// mismatch.
if !replayOK {
	return false
}

dag.mu.Lock()
dag.blocks[block.Hash] = block
// FIX (audit 2026-06-28 full recheck, P1-3): same durability gap as
// ProduceBlock — see its comment. Best-effort here too: a save failure
// doesn't reject an otherwise-valid, already-replayed peer block (its
// account-state effects are already committed), it just means this
// node would need to re-fetch this block from a peer again after a
// restart instead of having its own durable copy.
if err := dag.state.SaveBlockToDB(block); err != nil {
	fmt.Printf("[BLOCK] ⚠ Could not persist accepted peer block #%d (%s...) to chain_blocks: %v\n", block.Height, block.Hash[:16], err)
}

// Remove parents from tips
for _, ph := range block.ParentHashes {
	delete(dag.tips, ph)
}

// Add this block as new tip
dag.tips[block.Hash] = true

if block.Height > dag.height {
	dag.height = block.Height
	// FIX (double-apply): see the matching comment in ProduceBlock — persist
	// so a restart resumes from the true cumulative height.
	dag.state.setConfigValue("max_block_height", fmt.Sprintf("%d", dag.height))
}

tipCount := len(dag.tips)
dag.mu.Unlock()

// FIX (audit recheck3, P2 — "IncrementBlockCount laeuft asynchron und
// ist nicht konsensual deterministisch"): this was worse than just
// asynchronous — it was never called here at all. ProduceBlock only
// ever incremented blocks_produced for blocks THIS node itself
// produced; a peer block accepted here never touched the counter for
// ITS proposer. distributeValidatorsPoolLocked reads blocks_produced
// as the proportional reward weight (falling back to a minimum of 1
// for any registered node stuck at 0) — so on whichever single node
// actually runs distribution (DISTRIBUTION_ENABLED=true), every OTHER
// validator's real block production was invisible and they were
// floored to the same token "1" weight regardless of how active they
// actually were, while the distribution node's own blocks counted
// fully. Real fix, not just "make it synchronous": count every
// accepted block here too, for whichever proposer signed it — that's
// the only way this node's blocks_produced table ends up reflecting
// every validator's actual production, not just its own.
dag.state.IncrementBlockCount(block.Proposer)

// Now that this block exists (and has been replayed), any blocks that were
// queued as orphans waiting specifically on this hash as their missing
// parent can be retried. Done via a fresh top-level AddPeerBlock call
// rather than recursing — this naturally cascades: if a retried orphan
// succeeds, its own dependents get resolved the same way when ITS
// insertion reaches this point.
for _, waiting := range dag.popOrphans(block.Hash) {
	dag.AddPeerBlock(waiting)
}

dag.lastSuccessfulPeerSyncAt.Store(time.Now().Unix())
fmt.Printf("[DAG] ✓ Added peer block #%d | Tips: %d\n", block.Height, tipCount)
return true
}

// TotalStateRootMismatches sums every proposer's consecutive StateRoot
// mismatch counter, and LastSuccessfulPeerSyncAt returns the Unix
// timestamp of the last accepted peer block (0 if none yet this process) —
// both exposed via /api/health/combined (Gesamtaudit 2026-06-28, P2-4/P3-7).
func (dag *BlockDAG) TotalStateRootMismatches() int {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	total := 0
	for _, n := range dag.stateRootMismatches {
		total += n
	}
	return total
}

func (dag *BlockDAG) LastSuccessfulPeerSyncAt() int64 {
	return dag.lastSuccessfulPeerSyncAt.Load()
}

// Note: uses Go's built-in min() (available since Go 1.21; this module
// targets 1.24.1) rather than a custom helper — other files in this
// package already define min4()/min4b() specifically to avoid shadowing
// the built-in, so we follow that same convention here by not shadowing it.

// LatestBlock returns a single representative tip for display purposes
// (e.g. /api/status's "latest_hash"/"height"). With more than one validator
// producing concurrently, it's normal and expected to have multiple tips at
// the same max height for brief windows until the next block merges them —
// that's the DAG working as intended, not a fork.
//
// FIX: this used to pick whichever same-height tip happened to be visited
// first during Go's randomized map iteration — never replaced on a height
// TIE (only on strictly greater height) — so two nodes that both genuinely
// held the exact same set of tips could still report two DIFFERENT
// "latest_hash" values purely because their map iteration order differed.
// That's a misleading status signal, not an actual ledger divergence
// (StateRoot is computed from the full account/pool/nullifier state via
// replay, independent of which tip this function reports) — but it made
// "are these nodes in sync" impossible to answer just by comparing
// /api/status output, confirmed in production. Tie-break deterministically
// on hash so any two nodes holding the identical tip set always agree on
// which one to report, regardless of map iteration order.
func (dag *BlockDAG) LatestBlock() *Block {
dag.mu.RLock()
defer dag.mu.RUnlock()
var latest *Block
for hash := range dag.tips {
b := dag.blocks[hash]
if latest == nil || b.Height > latest.Height || (b.Height == latest.Height && b.Hash < latest.Hash) {
latest = b
}
}
return latest
}

func (dag *BlockDAG) Height() int64 {
dag.mu.RLock()
defer dag.mu.RUnlock()
return dag.height
}

func (dag *BlockDAG) GetBlocks() []*Block {
dag.mu.RLock()
defer dag.mu.RUnlock()
result := make([]*Block, 0, len(dag.blocks))
for _, b := range dag.blocks {
result = append(result, b)
}
// P3-1: O(n log n) sort instead of O(n^2) bubble sort.
sort.Slice(result, func(i, j int) bool { return result[i].Height < result[j].Height })
return result
}

// GetBlockByHash returns the block with the given hash, or nil if unknown.
// Used by /api/block/{hash} so a syncing peer can fetch one specific
// missing-ancestor block directly instead of relying solely on the
// height-windowed /api/blocks pagination (see fetchMissingAncestors).
func (dag *BlockDAG) GetBlockByHash(hash string) *Block {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	return dag.blocks[hash]
}

// GetBlockByHeight returns a block at the given height, or nil if none
// exists. Multiple validators can produce a sibling at the same height —
// when that happens this prefers the one with the most parent hashes (the
// merge block), matching the explorer UI's own dedup-by-height preference,
// so a search for a specific height shows the same block the list view
// would have shown for it.
func (dag *BlockDAG) GetBlockByHeight(height int64) *Block {
	dag.mu.RLock()
	defer dag.mu.RUnlock()
	var best *Block
	for _, b := range dag.blocks {
		if b.Height != height {
			continue
		}
		if best == nil || len(b.ParentHashes) > len(best.ParentHashes) {
			best = b
		}
	}
	return best
}

func (dag *BlockDAG) TotalBlocks() int {
dag.mu.RLock()
defer dag.mu.RUnlock()
return len(dag.blocks)
}

func (dag *BlockDAG) GetTips() []string {
dag.mu.RLock()
defer dag.mu.RUnlock()
tips := make([]string, 0, len(dag.tips))
for hash := range dag.tips {
tips = append(tips, hash)
}
return tips
}

// replayTransactions applies all TX types from a peer block to the local
// state. The block's ECDSA signature was already verified against an
// authorized validator before this function is reached.
//
// Design principle: secondary nodes apply the STORED amounts directly
// rather than re-running business logic. This avoids divergence from
// pool-state differences, floating-point order sensitivity, and
// demurrage timing differences between nodes.
// replayTransactions applies block's transactions to local state and
// returns false if the block was rolled back due to a genuine
// state-inconsistency failure (and should therefore NOT be considered
// applied — see AddPeerBlock, which rejects the whole block in that case).
//
// FIX (block-level atomicity): this used to apply each TX with continue-on-
// error and no way to undo TXs that had already succeeded earlier in the
// SAME block. A block with TX1 (succeeds) and TX2 (genuinely fails —
// insufficient balance, missing account; NOT an expected idempotent skip
// like "already registered") ended up partially applied: this node's state
// reflected less than what the producer's block.StateRoot was computed
// against. Money-moving TX types (transfer/swap/add_liquidity/
// remove_liquidity/faucet) now snapshot the touched accounts + pool before
// replay and roll back to that snapshot if any of them hits a genuine
// failure, so a failed block changes nothing instead of changing some of
// what it intended to. register_human's existing per-TX skip conditions
// (already registered, invalid proof, malformed data) are deliberately
// NOT treated as block-wide failures — those are intentional content
// rejections, not signs of state divergence, and were already
// self-consistent (TryClaimNullifier/ReleaseNullifier are already a
// correctly paired claim/release).
func (dag *BlockDAG) replayTransactions(block *Block) bool {
	// Fix 4: Deduplication guard — if this block has already been replayed,
	// skip it. Prevents double-credits when a block is delivered more than once.
	dag.replayedMu.Lock()
	if dag.replayedBlocks[block.Hash] {
		dag.replayedMu.Unlock()
		return true // already successfully replayed
	}
	dag.replayedMu.Unlock()

	// FIX (double-apply on snapshot bootstrap): a node that imported a
	// snapshot already has the cumulative effect of every block up to and
	// including snapshot_import_height baked into cs.accounts. Without
	// this guard, the HTTP-SYNC catch-up that follows (which always starts
	// from height 0, since dag.blocks is empty in memory after any
	// restart regardless of what the snapshot seeded) would apply every
	// pre-snapshot block's transactions a second time on top of the
	// already-current balances — confirmed in production: two secondary
	// nodes that bootstrapped from the same primary snapshot both ended up
	// crediting one wallet +2 AEQ and debiting another -2 AEQ relative to
	// the primary, exactly matching one historical transfer being replayed
	// twice. Mark the block as replayed (so dedup/tips/hash-chain
	// bookkeeping in the caller proceeds normally) without touching state.
	skipHeight := dag.bootHeight
	// FIX (audit 2026-06-28 recheck 4, P0-1): this runs before
	// dag.state.mu.Lock() is taken further down in this function — must use
	// the plain DB-only read, never cs.dbExec()/cs.activeTx, or this could
	// race against a concurrent atomic operation's in-flight transaction.
	if heightStr := dag.state.getConfigValueDB("snapshot_import_height"); heightStr != "" {
		var snapshotHeight int64
		fmt.Sscanf(heightStr, "%d", &snapshotHeight)
		if snapshotHeight > skipHeight {
			skipHeight = snapshotHeight
		}
	}
	// FIX (audit recheck 2, P0 #1 follow-up): bootHeight covers the more
	// general case the snapshot_import_height guard above was written for
	// — ANY node whose cs.accounts already reflects history that its
	// in-memory dag.blocks/dag.tips don't, not just one that bootstrapped
	// via snapshot. Confirmed in production within minutes of deploying
	// the StateRoot hard-reject above: a plain node restart (no snapshot
	// involved) immediately got stuck rejecting every single ancestor
	// block during ordinary post-restart catch-up, because cs.accounts
	// (loaded fully from the DB) was already at the LATEST state while
	// each ancestor block's claimed StateRoot reflects state as of THAT
	// historical height — comparing "now" against "back then" was always
	// going to mismatch, with no real divergence involved. Below
	// skipHeight, the block is still fetched and inserted into
	// dag.blocks/dag.tips (hash-chain/tips bookkeeping needs it as a valid
	// parent for later blocks) but neither its transactions nor its
	// StateRoot claim are touched.
	// FIX (audit recheck2, P2 #4): naming this explicitly, since it's a real
	// trust boundary, not just a performance shortcut. Every block ABOVE
	// skipHeight is independently re-verified by replaying its transactions
	// and checking the resulting StateRoot — this node never just trusts a
	// peer's claim. Every block AT OR BELOW skipHeight is "snapshot trust
	// mode": this node accepts cs.accounts' already-loaded state (from its
	// own DB, or from a signed snapshot import) as correct for that range
	// without re-deriving it from block history, because re-deriving it
	// would require replaying transactions whose effects are already
	// baked into that state by definition (see bootHeight's and
	// snapshot_import_height's comments for why re-replaying them would
	// double-apply, not re-verify, them). The actual trust anchor for
	// snapshot-sourced state is ImportSnapshotFromURL/ResyncFromSnapshotURL's
	// mandatory ECDSA signature check against BOOTSTRAP_SIGNER — this skip
	// doesn't grant any trust itself, it just avoids re-deriving what that
	// signature check already vouched for.
	if skipHeight > 0 && block.Height <= skipHeight {
		dag.replayedMu.Lock()
		dag.replayedBlocks[block.Hash] = true
		dag.replayedMu.Unlock()
		return true
	}

	touchedAddrs, needsFullSnapshot := blockTouchedAddresses(block)
	// FIX (audit recheck3, P0/P1 — "Block-Replay-Rollback ist nicht gegen
	// parallele lokale Mutationen isoliert"): this used to take rollbackSnap
	// via the lock-acquiring, lock-releasing snapshotForRollback, then let
	// every individual Delta function below take and release cs.mu on its
	// own for just that one call. Between any two of those calls — or
	// between the snapshot and the first call — a concurrent API operation
	// or distribution round (each its own complete runAtomicWithOutbox/
	// runAtomicDistributionWithOutbox critical section) could mutate the
	// very same account, fully commit, and report success to its own
	// caller — and if THIS replay later hit a hardFailure or StateRoot
	// mismatch unrelated to that account, rolling back with rollbackSnap
	// would revert it anyway, silently erasing an already-committed,
	// unrelated, successful operation. Holding cs.mu continuously from the
	// snapshot below through either a successful StateRoot match or a
	// rollback closes that gap: every Delta call in this loop now uses its
	// "...Locked" sibling (assumes cs.mu already held) instead of the
	// public lock-each-time wrapper, and the snapshot/rollback/StateRoot
	// comparison below do the same.
	dag.state.mu.Lock()
	defer dag.state.mu.Unlock()
	configBackup := make(map[string]configValueSnapshot, len(stateRootRelevantConfigKeys))
	for _, key := range stateRootRelevantConfigKeys {
		value, existed := dag.state.getConfigValueExists(key)
		configBackup[key] = configValueSnapshot{value: value, existed: existed}
	}
	rollbackSnap := dag.state.snapshotForRollbackLocked(touchedAddrs, needsFullSnapshot, configBackup)

	// FIX (audit 2026-06-28 full recheck, P0-4 — "Replay-Rollback ist nicht
	// als DB-Transaktion isoliert"): every DB write this replay makes (via
	// saveAccountToDB/savePoolToDB/setConfigValue, all routed through
	// cs.dbExec()) used to auto-commit immediately on its own, with
	// rollbackSnap/restoreFromRollbackLocked emulating "undo" at the
	// application level by recomputing and rewriting old values — not a
	// real SQL rollback. If a step partway through failed in a way that
	// left some writes committed and others not, the application-level
	// restore could only ever re-derive what it already knew about
	// (rollbackSnap's captured fields), not guarantee every write this
	// replay made was actually undone. A real DB transaction makes that
	// guarantee structurally instead of by careful bookkeeping: every
	// dbExec() call below joins dbTx (set as cs.activeTx for the duration),
	// and either ALL of them commit together or tx.Rollback() discards
	// every one of them atomically. rollbackSnap/restoreFromRollbackLocked
	// are still used for the IN-MEMORY side (cs.accounts/cs.pool are plain
	// Go maps with no transactional semantics of their own — only the DB
	// side can be made truly atomic this way), now redundant-but-harmless
	// for the DB side specifically when a real rollback already ran (the
	// same pattern runAtomicWithOutbox already established: tx.Rollback()
	// for the DB, restoreFromRollback for memory, in that order).
	var dbTx *sql.Tx
	if dag.state.db != nil {
		var err error
		dbTx, err = dag.state.db.Begin()
		if err != nil {
			fmt.Printf("[REPLAY] ✗ Block #%d: could not begin replay transaction: %v — block rejected\n", block.Height, err)
			return false
		}
		dag.state.activeTx = dbTx
	}
	// commitOrRollback finalizes dbTx according to success, clearing
	// activeTx either way so no write after this point accidentally joins
	// a transaction that's already been resolved. Returns an error if a
	// commit was attempted and failed (caller must then treat this exactly
	// like any other hardFailure, including the in-memory restore).
	commitOrRollback := func(success bool) error {
		if dbTx == nil {
			dag.state.activeTx = nil
			return nil
		}
		dag.state.activeTx = nil
		if !success {
			if err := dbTx.Rollback(); err != nil {
				fmt.Printf("[REPLAY] Warning: replay transaction rollback for block #%d failed: %v\n", block.Height, err)
			}
			return nil
		}
		return dbTx.Commit()
	}
	hardFailure := false
	var claimedNullifiers []string

	for _, tx := range block.Transactions {
		if hardFailure {
			break // stop applying further TXs once we know this block is being rolled back
		}
		wallet := strings.ToLower(strings.TrimSpace(tx.Wallet))
		switch tx.Type {

		case "register_human":
			nullifier := strings.TrimSpace(tx.Nullifier)
			commitment := strings.TrimSpace(tx.Commitment)
			if wallet == "" || nullifier == "" {
				fmt.Printf("[REPLAY] ⚠ Skipping register_human in block #%d: missing wallet or nullifier (older node version?)\n", block.Height)
				continue
			}
			// FIX (audit recheck2, P1 #10): malformed wallet/nullifier, missing
			// proof fields, and an invalid proof used to just `continue` (skip
			// this TX, accept the rest of the block) instead of hardFailure.
			// Unlike the wallet==""/nullifier=="" case above (genuine legacy
			// compat — see its own comment), there is no legitimate node
			// version that ever produces a malformed wallet/short nullifier or
			// omits proof fields; register.go always populates them. A block
			// containing one of these is either a bug in the producing node
			// or a validator deliberately packing unverifiable "registrations"
			// into otherwise-valid block history — permanently, since an
			// accepted block is never revisited. Treating these as the same
			// genuine state-inconsistency failure every other case in this
			// switch already hardFails on closes that gap.
			if len(wallet) != 42 || wallet[:2] != "0x" {
				fmt.Printf("[REPLAY] ✗ register_human in block #%d: malformed wallet %q — rolling back whole block\n", block.Height, wallet)
				hardFailure = true
				continue
			}
			if len(nullifier) < 16 {
				fmt.Printf("[REPLAY] ✗ register_human in block #%d: nullifier too short %q — rolling back whole block\n", block.Height, nullifier)
				hardFailure = true
				continue
			}
			// Verify the ZK proof via BioVerifier before applying. This
			// eliminates unconditional trust in the validator ECDSA key for
			// registration TXs — a compromised validator key cannot inject
			// fake registrations without also producing a valid Groth16 proof.
			//
			// FIX: this used to skip verification entirely (falling through to
			// trust the block signature alone) whenever proof fields were
			// absent, "for backward compatibility with old nodes". Both
			// current TX-creation sites (register.go) always populate
			// ProofA/B/C/PubSignals, so no legitimate code path produces a
			// register_human TX without them anymore — that fallback was pure
			// attack surface letting any authorized validator (or one whose
			// signing key leaked) inject registrations for arbitrary wallets
			// with no biometric proof at all, defeating "one human, one
			// registration" silently.
			if len(tx.ProofA) != 2 || len(tx.ProofB) != 2 || len(tx.ProofC) != 2 || len(tx.PubSignals) < 2 {
				fmt.Printf("[REPLAY] ✗ register_human for %s (block #%d): missing ZK proof fields — rolling back whole block\n", wallet, block.Height)
				hardFailure = true
				continue
			}
			if !dag.verifyZKProof(tx) {
				fmt.Printf("[REPLAY] ✗ register_human for %s (block #%d): ZK proof verification failed — rolling back whole block\n", wallet, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ ZK proof verified for %s (block #%d)\n", wallet, block.Height)
			// FIX (audit 2026-06-28 recheck 5, P1-1): tryClaimNullifierLocked
			// now returns an error distinctly from "already used" — a genuine
			// DB failure during the claim must roll back the block, not be
			// silently treated as a normal duplicate-registration skip.
			claimed, claimErr := dag.state.tryClaimNullifierLocked(nullifier, wallet)
			if claimErr != nil {
				fmt.Printf("[REPLAY] ✗ register_human for %s (block #%d): nullifier claim DB error: %v — rolling back whole block\n", wallet, block.Height, claimErr)
				hardFailure = true
				continue
			}
			if !claimed {
				continue // already registered
			}
			if err := dag.state.registerHumanLocked(wallet); err != nil {
				// FIX: release the nullifier claimed two lines above on failure —
				// it used to stay claimed forever ("nullifier recorded, balance
				// NOT credited"), permanently burning that biometric for
				// everyone even though no registration ever actually completed
				// with it (e.g. wallet already human via a different nullifier).
				dag.state.releaseNullifierLocked(nullifier)
				fmt.Printf("[REPLAY] ✗ RegisterHuman %s: %v (nullifier released, balance NOT credited)\n", wallet, err)
				continue
			}
			// Track this claim so it can be released too if a LATER TX in
			// this same block hard-fails and the whole block gets rolled
			// back — the account-balance/IsHuman side of this registration
			// is already covered by rollbackSnap (this wallet is in
			// blockTouchedAddresses via tx.Wallet), but cs.nullifiers is a
			// separate map the account snapshot doesn't touch.
			claimedNullifiers = append(claimedNullifiers, nullifier)
			if commitment != "" {
				_ = dag.state.SaveBioRegistration(commitment, wallet, tx.TxHash, "")
			}
			fmt.Printf("[REPLAY] ✓ Applied register_human for %s (block #%d)\n", wallet, block.Height)

		case "transfer":
			to := strings.ToLower(strings.TrimSpace(tx.To))
			if wallet == "" || to == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping transfer in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.applyTransferDeltaLocked(wallet, to, tx.Amount, tx.FromDemurrageLost, tx.ToDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ Transfer %s->%s %.6f: %v (block #%d) — rolling back whole block\n", wallet, to, tx.Amount, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied transfer %.6f AEQ: %s->%s (block #%d)\n", tx.Amount, wallet, to, block.Height)

		case "swap_aeq_tusd":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_aeq_tusd in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.applySwapDeltaLocked(wallet, tx.Amount, tx.AmountOut, true, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_aeq_tusd %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_aeq_tusd %.6f AEQ->%.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "swap_tusd_aeq":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_tusd_aeq in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.applySwapDeltaLocked(wallet, tx.Amount, tx.AmountOut, false, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_tusd_aeq %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_tusd_aeq %.6f tUSD->%.6f AEQ for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "add_liquidity":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping add_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.addLiquidityDeltaLocked(wallet, tx.Amount, tx.AmountOut, tx.LPShares, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ add_liquidity %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied add_liquidity %.6f AEQ + %.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "remove_liquidity":
			if wallet == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping remove_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.removeLiquidityDeltaLocked(wallet, tx.Amount, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ remove_liquidity %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied remove_liquidity %.6f shares for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "faucet":
			if wallet == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping faucet in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.applyFaucetDeltaLocked(wallet, tx.Amount); err != nil {
				fmt.Printf("[REPLAY] ✗ faucet %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied faucet %.6f tUSD for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "ubi_distribution":
			// FIX (audit recheck 2, P0 #5): main.go now emits ONE of these per
			// human (Wallet set, AmountPerHuman omitted) instead of a single
			// flat broadcast — see ApplyUBIRewardDelta's comment. The
			// AmountPerHuman>0 branch below only fires for historical blocks
			// produced by older node versions before this change; new blocks
			// never set that field.
			if tx.AmountPerHuman > 0 {
				if err := dag.state.applyUBIDeltaLocked(tx.AmountPerHuman, block.Timestamp); err != nil {
					fmt.Printf("[REPLAY] ✗ legacy flat ubi_distribution: %v (block #%d) — rolling back whole block\n", err, block.Height)
					hardFailure = true
					continue
				}
				fmt.Printf("[REPLAY] ✓ Applied legacy flat UBI distribution %.6f AEQ/human (block #%d)\n", tx.AmountPerHuman, block.Height)
			} else if wallet != "" && wallet != "0x0000000000000000000000000000000000000000" {
				if err := dag.state.applyUBIRewardDeltaLocked(wallet, tx.Amount, tx.FromDemurrageLost); err != nil {
					fmt.Printf("[REPLAY] ✗ ubi_distribution %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
					hardFailure = true
					continue
				}
				fmt.Printf("[REPLAY] ✓ Applied UBI reward %.6f AEQ for %s (block #%d)\n", tx.Amount, wallet, block.Height)
			} else {
				// FIX (audit recheck2, P1 #10): see register_human's matching
				// comment — this TX type is only ever emitted internally by
				// RunDailyDistributionAtomic, never user-submitted, so a
				// malformed one (neither shape populated) means either a
				// producer bug or a validator fabricating distribution
				// history. hardFailure instead of a silent skip.
				fmt.Printf("[REPLAY] ✗ ubi_distribution TX in block #%d has neither amount_per_human nor a wallet — rolling back whole block\n", block.Height)
				hardFailure = true
				continue
			}

		case "ubi_distribution_finalize":
			if err := dag.state.applyUBIFinalizeDeltaLocked(tx.DistributionAt); err != nil {
				fmt.Printf("[REPLAY] ✗ ubi_distribution_finalize: %v (block #%d) — rolling back whole block\n", err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Finalized UBI round, last_ubi_at=%d (block #%d)\n", tx.DistributionAt, block.Height)

		case "validator_distribution":
			wallet := strings.ToLower(tx.Wallet)
			if err := dag.state.applyValidatorRewardDeltaLocked(wallet, tx.Amount, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ validator_distribution %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied validator reward %.6f AEQ for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "validator_distribution_pool_zero":
			if err := dag.state.applyValidatorPoolZeroDeltaLocked(); err != nil {
				fmt.Printf("[REPLAY] ✗ validator_distribution_pool_zero: %v (block #%d) — rolling back whole block\n", err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Zeroed validators pool (block #%d)\n", block.Height)

		case "lp_distribution":
			wallet := strings.ToLower(tx.Wallet)
			if err := dag.state.applyLPRewardDeltaLocked(wallet, tx.Amount, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ lp_distribution %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied LP reward %.6f AEQ for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "lp_distribution_pool_zero":
			if err := dag.state.applyLPPoolZeroDeltaLocked(); err != nil {
				fmt.Printf("[REPLAY] ✗ lp_distribution_pool_zero: %v (block #%d) — rolling back whole block\n", err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Zeroed LP pool (block #%d)\n", block.Height)

		case "escrow_move":
			wallet := strings.ToLower(tx.Wallet)
			if err := dag.state.applyEscrowMoveDeltaLocked(wallet, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ escrow_move %s: %v (block #%d) — rolling back whole block\n", wallet, err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied escrow move for %s (block #%d)\n", wallet, block.Height)

		case "escrow_release":
			if err := dag.state.applyEscrowReleaseDeltaLocked(tx.Amount); err != nil {
				fmt.Printf("[REPLAY] ✗ escrow_release: %v (block #%d) — rolling back whole block\n", err, block.Height)
				hardFailure = true
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied escrow release %.6f AEQ → UBI pool (block #%d)\n", tx.Amount, block.Height)

		default:
			// FIX (audit 2026-06-28 recheck 4, P2-2): unknown TX types used to
			// be silently ignored — applied no delta, but also didn't reject
			// the block. That's a forward-compatibility hazard: a node
			// running OLDER code that doesn't yet recognize a NEW TX type
			// introduced by upgraded peers would silently skip that TX's
			// economic effect while still accepting the block as valid. The
			// post-replay StateRoot comparison below would usually catch
			// this (the skipped delta means local state can't match the
			// proposer's claimed root) — but relying on StateRoot alone to
			// catch a known, structural gap is exactly the "we believe it's
			// atomic" pattern this audit pass has been closing elsewhere.
			// Hard-fail explicitly instead: an unrecognized type is treated
			// the same as any other genuine state-inconsistency failure.
			fmt.Printf("[REPLAY] ✗ Unknown TX type %q (block #%d) — rolling back whole block\n", tx.Type, block.Height)
			hardFailure = true
			continue
		}
	}

	if hardFailure {
		commitOrRollback(false) // real SQL ROLLBACK — see commitOrRollback's comment
		if rbErr := dag.state.restoreFromRollbackLocked(rollbackSnap); rbErr != nil {
			fmt.Printf("[REPLAY] CRITICAL: rollback persistence failed for block #%d — memory/DB may now disagree: %v\n", block.Height, rbErr)
		}
		for _, n := range claimedNullifiers {
			dag.state.releaseNullifierLocked(n)
		}
		fmt.Printf("[REPLAY] ✗ Block #%d rolled back due to a genuine state-inconsistency failure — block rejected\n", block.Height)
		return false
	}

	// FIX (audit recheck 2, P0 #1): StateRoot comparison moved here (from
	// AddPeerBlock, after this function returned) so a mismatch can use
	// THIS function's own rollbackSnap to actually undo the replay, not
	// just log it. This used to be warn-only: the block was accepted into
	// dag.blocks/dag.tips regardless, meaning a node could build on top of
	// a block whose state it could not itself reproduce. Sequenced after
	// the distribution-atomicity and per-human demurrage-replay fixes
	// (audit recheck 2, P0 #2-#6) specifically because those were the
	// known, frequent divergence sources that made this check fire on
	// nearly every block in practice — rejecting on every block would have
	// halted sync entirely rather than catching genuine divergence.
	// Known residual divergence sources (non-atomic nullifier persistence,
	// mirror-path outbox — audit recheck 2, P1 #7/#8) can still trigger
	// this; a rejected block is retried by a later sync cycle once local
	// state catches up, the same recovery path hardFailure above already
	// relies on.
	if block.StateRoot != "" {
		// FIX (audit 2026-06-28 full recheck, P0-4): computed BEFORE
		// commit/rollback, while dbTx is still open — this must reflect
		// exactly what's about to be committed (or discarded), the same
		// view every write in this replay just made within dbTx.
		localRoot := dag.state.stateRootLocked(dag.state.getConfigValue("last_ubi_at"))
		if block.StateRoot != localRoot {
			commitOrRollback(false) // real SQL ROLLBACK
			if rbErr := dag.state.restoreFromRollbackLocked(rollbackSnap); rbErr != nil {
				fmt.Printf("[REPLAY] CRITICAL: rollback persistence failed for block #%d — memory/DB may now disagree: %v\n", block.Height, rbErr)
			}
			for _, n := range claimedNullifiers {
				dag.state.releaseNullifierLocked(n)
			}
			fmt.Printf("[REPLAY] ✗ StateRoot mismatch on block #%d (proposer=%s..., local=%s...) — rolled back, block rejected\n",
				block.Height, block.StateRoot[:min(16, len(block.StateRoot))], localRoot[:min(16, len(localRoot))])
			dag.stateRootMismatches[block.Proposer]++
			if dag.stateRootMismatches[block.Proposer] >= 5 {
				fmt.Printf("[ALERT] 5+ consecutive StateRoot mismatches from proposer %s — state may have diverged. Consider resync.\n", block.Proposer)
			}
			return false
		}
		dag.stateRootMismatches[block.Proposer] = 0 // reset on match
	}

	// FIX (audit 2026-06-28 full recheck, P0-4): commit dbTx now that every
	// check has passed — this is the moment every DB write this replay made
	// actually becomes durable, all together, as one SQL transaction. If
	// commit itself fails (rare: connection loss, constraint violation the
	// DB only catches at commit time), Postgres has already rolled the
	// transaction back server-side — treat it exactly like any other
	// rollback path: restore in-memory state and reject the block, so
	// memory and DB can't end up disagreeing about whether this block
	// applied.
	if commitErr := commitOrRollback(true); commitErr != nil {
		if rbErr := dag.state.restoreFromRollbackLocked(rollbackSnap); rbErr != nil {
			fmt.Printf("[REPLAY] CRITICAL: rollback persistence failed for block #%d — memory/DB may now disagree: %v\n", block.Height, rbErr)
		}
		for _, n := range claimedNullifiers {
			dag.state.releaseNullifierLocked(n)
		}
		fmt.Printf("[REPLAY] ✗ Block #%d: replay transaction commit failed (rolled back, block rejected): %v\n", block.Height, commitErr)
		return false
	}

	dag.replayedMu.Lock()
	// FIX 1: Cap the cache to prevent unbounded growth (memory leak).
	// dag.blocks is the authoritative deduplication store; this is a fast-path cache.
	if len(dag.replayedBlocks) > 50000 {
		dag.replayedBlocks = make(map[string]bool, 1000)
	}
	dag.replayedBlocks[block.Hash] = true
	dag.replayedMu.Unlock()
	return true
}

// ReconstructState is a no-op: the PostgreSQL database is the authoritative
// source of truth and is already loaded by ChainState.LoadFromDB() before
// this is called. Ongoing state sync happens via replayRegistrations(), which
// is called from AddPeerBlock for every received block.
func (dag *BlockDAG) ReconstructState(state *ChainState) {
	fmt.Printf("[CHAIN] State loaded from DB — skipping full block-replay reconstruction\n")
}

// Close is a no-op now that replay runs synchronously inside AddPeerBlock
// (see its comment for why the old async channel+goroutine design was
// removed) — there's no longer a background goroutine to shut down. Kept
// for call-site compatibility (main.go may call this on shutdown).
func (dag *BlockDAG) Close() {}

// verifyZKProof reconstructs the Groth16 proof from the TX's decimal string
// fields and calls the BioVerifier contract via the local EVM engine to check
// validity. Returns true when the proof is valid, false otherwise.
// Only called when all four proof fields (ProofA, ProofB, ProofC, PubSignals)
// are present — blocks from old nodes omit them and fall back to trust-based
// validation for backward compatibility.
func (dag *BlockDAG) verifyZKProof(tx Transaction) bool {
	if dag.evm == nil {
		fmt.Printf("[WARN] EVM not initialized, rejecting ZK proof for block safety\n")
		return false
	}

	// Parse ProofA [2]*big.Int
	if len(tx.ProofA) != 2 || len(tx.ProofC) != 2 || len(tx.PubSignals) < 2 || len(tx.ProofB) != 2 {
		return false
	}
	var pA [2]*big.Int
	for i := 0; i < 2; i++ {
		n := new(big.Int)
		if _, ok := n.SetString(tx.ProofA[i], 10); !ok {
			fmt.Printf("[REPLAY] ✗ verifyZKProof: invalid ProofA[%d]: %q\n", i, tx.ProofA[i])
			return false
		}
		pA[i] = n
	}

	// Parse ProofB [2][2]*big.Int
	var pB [2][2]*big.Int
	for i := 0; i < 2; i++ {
		if len(tx.ProofB[i]) != 2 {
			fmt.Printf("[REPLAY] ✗ verifyZKProof: ProofB[%d] has wrong length\n", i)
			return false
		}
		for j := 0; j < 2; j++ {
			n := new(big.Int)
			if _, ok := n.SetString(tx.ProofB[i][j], 10); !ok {
				fmt.Printf("[REPLAY] ✗ verifyZKProof: invalid ProofB[%d][%d]: %q\n", i, j, tx.ProofB[i][j])
				return false
			}
			pB[i][j] = n
		}
	}

	// Parse ProofC [2]*big.Int
	var pC [2]*big.Int
	for i := 0; i < 2; i++ {
		n := new(big.Int)
		if _, ok := n.SetString(tx.ProofC[i], 10); !ok {
			fmt.Printf("[REPLAY] ✗ verifyZKProof: invalid ProofC[%d]: %q\n", i, tx.ProofC[i])
			return false
		}
		pC[i] = n
	}

	// Parse PubSignals [2]*big.Int (only first two are needed by verifyProof)
	var pubSignals [2]*big.Int
	for i := 0; i < 2; i++ {
		n := new(big.Int)
		if _, ok := n.SetString(tx.PubSignals[i], 10); !ok {
			fmt.Printf("[REPLAY] ✗ verifyZKProof: invalid PubSignals[%d]: %q\n", i, tx.PubSignals[i])
			return false
		}
		pubSignals[i] = n
	}

	verifierABI, err := abi.JSON(strings.NewReader(bioVerifierABI))
	if err != nil {
		fmt.Printf("[REPLAY] ✗ verifyZKProof: ABI parse failed: %v\n", err)
		return false
	}
	verifyData, err := verifierABI.Pack("verifyProof", pA, pB, pC, pubSignals)
	if err != nil {
		fmt.Printf("[REPLAY] ✗ verifyZKProof: ABI encode failed: %v\n", err)
		return false
	}

	caller := common.HexToAddress(tx.Wallet)
	ret, err := dag.evm.CallContract(caller, common.HexToAddress(BIO_VERIFIER_ADDR), verifyData, big.NewInt(0), false)
	if err != nil {
		fmt.Printf("[REPLAY] ✗ verifyZKProof: BioVerifier call failed: %v\n", err)
		return false
	}
	if len(ret) != 32 || ret[31] != 1 {
		return false
	}
	return true
}
