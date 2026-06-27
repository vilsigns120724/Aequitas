package keeper

import (
"crypto/ecdsa"
"crypto/rand"
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"math/big"
"os"
"strings"
"sort"
"sync"
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
	// orphans holds blocks whose parent isn't known yet, keyed by the missing
	// parent's hash. When that parent is later added, every block waiting on
	// it is retried automatically. See AddPeerBlock for why this exists —
	// without it, a block whose parent arrived even one sync cycle late was
	// silently dropped forever, along with everything built on top of it.
	orphans   map[string][]*Block
	orphansMu sync.Mutex
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
return dag
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

// Only now that the block carrying them is durably stored in dag.blocks —
// clear the DB outbox rows. See LoadAndClearPendingTxs's doc comment for
// why this is no longer a single, earlier delete-then-build step.
if len(pendingTxIDs) > 0 {
	dag.state.ClearPendingTxs(pendingTxIDs)
}

// Record that this proposer produced a block — used for proportional
// validator-reward distribution in DistributeValidatorsPool.
go dag.state.IncrementBlockCount(proposer)

// Remove all parents from tips, add this block as new tip
for _, ph := range parentHashes {
delete(dag.tips, ph)
}
dag.tips[block.Hash] = true
dag.height = block.Height

if len(parentHashes) > 1 {
fmt.Printf("[DAG] 🔀 Merged %d tips into block #%d\n", len(parentHashes), block.Height)
}

return block
}

// maxOrphans caps total queued orphan blocks across all missing-parent keys,
// so a malicious or buggy peer sending blocks that reference parents which
// will never arrive can't grow this map without bound.
const maxOrphans = 2000

// queueOrphan stores block, which is waiting on missingParent to appear,
// and logs the wait (the old code dropped this case with zero logging).
func (dag *BlockDAG) queueOrphan(missingParent string, block *Block) {
	dag.orphansMu.Lock()
	defer dag.orphansMu.Unlock()
	total := 0
	for _, v := range dag.orphans {
		total += len(v)
	}
	if total >= maxOrphans {
		fmt.Printf("[DAG] ✗ Dropped peer block #%d: orphan buffer full (%d waiting), missing parent never arrived\n",
			block.Height, total)
		return
	}
	dag.orphans[missingParent] = append(dag.orphans[missingParent], block)
	fmt.Printf("[DAG] ⏳ Block #%d from %s queued as orphan — missing parent %s... (%d block(s) now waiting on it)\n",
		block.Height, block.Proposer, missingParent[:min(16, len(missingParent))], len(dag.orphans[missingParent]))
}

// popOrphans returns and removes every block that was waiting on parentHash.
func (dag *BlockDAG) popOrphans(parentHash string) []*Block {
	dag.orphansMu.Lock()
	defer dag.orphansMu.Unlock()
	waiting := dag.orphans[parentHash]
	delete(dag.orphans, parentHash)
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
case "", "register_human", "transfer", "swap_aeq_tusd", "swap_tusd_aeq", "add_liquidity", "remove_liquidity", "faucet", "ubi_distribution":
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
dag.replayTransactions(block)
dag.replayMu.Unlock()

// State-root integrity check — now comparing the RECEIVER's post-replay
// state against the PRODUCER's post-state, the comparison it was always
// supposed to be. Still warn-only rather than reject-on-mismatch: any
// remaining non-determinism source (not yet found) would otherwise halt
// sync entirely rather than just logging a now-meaningful warning. The
// ECDSA signature check above already prevents forged blocks; this is a
// state-consistency signal, not a forgery defense.
if block.StateRoot != "" {
	localRoot := dag.state.StateRoot()
	if block.StateRoot != localRoot {
		fmt.Printf("[DAG] ⚠ StateRoot mismatch on peer block #%d (proposer=%s..., local=%s...) — accepted (warn only)\n",
			block.Height, block.StateRoot[:min(16, len(block.StateRoot))], localRoot[:min(16, len(localRoot))])
		// FIX 4: Track mismatches per proposer.
		dag.stateRootMismatches[block.Proposer]++
		if dag.stateRootMismatches[block.Proposer] >= 5 {
			fmt.Printf("[ALERT] 5+ consecutive StateRoot mismatches from proposer %s — state may have diverged. Consider resync.\n", block.Proposer)
		}
	} else {
		dag.stateRootMismatches[block.Proposer] = 0 // reset on match
	}
}

dag.mu.Lock()
dag.blocks[block.Hash] = block

// Remove parents from tips
for _, ph := range block.ParentHashes {
	delete(dag.tips, ph)
}

// Add this block as new tip
dag.tips[block.Hash] = true

if block.Height > dag.height {
	dag.height = block.Height
}

tipCount := len(dag.tips)
dag.mu.Unlock()

// Now that this block exists (and has been replayed), any blocks that were
// queued as orphans waiting specifically on this hash as their missing
// parent can be retried. Done via a fresh top-level AddPeerBlock call
// rather than recursing — this naturally cascades: if a retried orphan
// succeeds, its own dependents get resolved the same way when ITS
// insertion reaches this point.
for _, waiting := range dag.popOrphans(block.Hash) {
	dag.AddPeerBlock(waiting)
}

fmt.Printf("[DAG] ✓ Added peer block #%d | Tips: %d\n", block.Height, tipCount)
return true
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
func (dag *BlockDAG) replayTransactions(block *Block) {
	// Fix 4: Deduplication guard — if this block has already been replayed,
	// skip it. Prevents double-credits when a block is delivered more than once.
	dag.replayedMu.Lock()
	if dag.replayedBlocks[block.Hash] {
		dag.replayedMu.Unlock()
		return // already replayed
	}
	// FIX 1: Cap the cache to prevent unbounded growth (memory leak).
	// dag.blocks is the authoritative deduplication store; this is a fast-path cache.
	if len(dag.replayedBlocks) > 50000 {
		dag.replayedBlocks = make(map[string]bool, 1000)
	}
	dag.replayedBlocks[block.Hash] = true
	dag.replayedMu.Unlock()

	for _, tx := range block.Transactions {
		wallet := strings.ToLower(strings.TrimSpace(tx.Wallet))
		switch tx.Type {

		case "register_human":
			nullifier := strings.TrimSpace(tx.Nullifier)
			commitment := strings.TrimSpace(tx.Commitment)
			if wallet == "" || nullifier == "" {
				fmt.Printf("[REPLAY] ⚠ Skipping register_human in block #%d: missing wallet or nullifier (older node version?)\n", block.Height)
				continue
			}
			if len(wallet) != 42 || wallet[:2] != "0x" {
				fmt.Printf("[REPLAY] ✗ Rejecting register_human in block #%d: malformed wallet %q\n", block.Height, wallet)
				continue
			}
			if len(nullifier) < 16 {
				fmt.Printf("[REPLAY] ✗ Rejecting register_human in block #%d: nullifier too short %q\n", block.Height, nullifier)
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
				fmt.Printf("[REPLAY] ✗ Rejecting register_human for %s (block #%d): missing ZK proof fields\n", wallet, block.Height)
				continue
			}
			if !dag.verifyZKProof(tx) {
				fmt.Printf("[REPLAY] ✗ ZK proof verification failed for %s (block #%d) — skipping\n", wallet, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ ZK proof verified for %s (block #%d)\n", wallet, block.Height)
			if !dag.state.TryClaimNullifier(nullifier, wallet) {
				continue // already registered
			}
			if err := dag.state.RegisterHuman(wallet); err != nil {
				// FIX: release the nullifier claimed two lines above on failure —
				// it used to stay claimed forever ("nullifier recorded, balance
				// NOT credited"), permanently burning that biometric for
				// everyone even though no registration ever actually completed
				// with it (e.g. wallet already human via a different nullifier).
				dag.state.ReleaseNullifier(nullifier)
				fmt.Printf("[REPLAY] ✗ RegisterHuman %s: %v (nullifier released, balance NOT credited)\n", wallet, err)
				continue
			}
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
			if err := dag.state.ApplyTransferDelta(wallet, to, tx.Amount, tx.FromDemurrageLost, tx.ToDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ Transfer %s->%s %.6f: %v (block #%d)\n", wallet, to, tx.Amount, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied transfer %.6f AEQ: %s->%s (block #%d)\n", tx.Amount, wallet, to, block.Height)

		case "swap_aeq_tusd":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_aeq_tusd in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.ApplySwapDelta(wallet, tx.Amount, tx.AmountOut, true, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_aeq_tusd %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_aeq_tusd %.6f AEQ->%.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "swap_tusd_aeq":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_tusd_aeq in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.ApplySwapDelta(wallet, tx.Amount, tx.AmountOut, false, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_tusd_aeq %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_tusd_aeq %.6f tUSD->%.6f AEQ for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "add_liquidity":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping add_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.AddLiquidityDelta(wallet, tx.Amount, tx.AmountOut, tx.LPShares, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ add_liquidity %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied add_liquidity %.6f AEQ + %.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "remove_liquidity":
			if wallet == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping remove_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.RemoveLiquidityDelta(wallet, tx.Amount, tx.FromDemurrageLost); err != nil {
				fmt.Printf("[REPLAY] ✗ remove_liquidity %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied remove_liquidity %.6f shares for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "faucet":
			if wallet == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping faucet in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.ApplyFaucetDelta(wallet, tx.Amount); err != nil {
				fmt.Printf("[REPLAY] ✗ faucet %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied faucet %.6f tUSD for %s (block #%d)\n", tx.Amount, wallet, block.Height)

		case "ubi_distribution":
			if tx.AmountPerHuman > 0 {
				dag.state.ApplyUBIDelta(tx.AmountPerHuman, block.Timestamp)
				fmt.Printf("[REPLAY] ✓ Applied UBI distribution %.6f AEQ/human (block #%d)\n", tx.AmountPerHuman, block.Height)
			} else {
				// Legacy TX from an older node version — no AmountPerHuman stored.
				fmt.Printf("[REPLAY] UBI distribution TX in block #%d — no amount_per_human field, skipping (old node format)\n", block.Height)
			}

		default:
			// empty string or other unknown types are silently ignored
		}
	}
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
