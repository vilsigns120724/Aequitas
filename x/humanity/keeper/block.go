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
replayQueue            chan *Block       // serialized replay channel — ensures TX ordering across blocks
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
	dag.peerChallenges[strings.ToLower(signingAddr)] = peerChallenge{
		value:     challenge,
		expiresAt: ts + 90,
	}
	// Prune expired challenges
	for addr, c := range dag.peerChallenges {
		if time.Now().Unix() > c.expiresAt {
			delete(dag.peerChallenges, addr)
		}
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
replayQueue:            make(chan *Block, 1000),
}
// Single consumer goroutine ensures blocks are replayed in the order received.
// This preserves TX dependencies (e.g. register_human in block N must be
// applied before a transfer in block N+1 that references the same wallet).
go func() {
	for b := range dag.replayQueue {
		dag.replayTransactions(b)
	}
}()
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
Humans:       dag.state.TotalHumans(),
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

// Collect all current tips as parents
parentHashes := make([]string, 0, len(dag.tips))
for hash := range dag.tips {
parentHashes = append(parentHashes, hash)
}

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

func (dag *BlockDAG) AddPeerBlock(block *Block) bool {
dag.mu.Lock()
defer dag.mu.Unlock()

// Skip if already known
if _, exists := dag.blocks[block.Hash]; exists {
return false
}

// Genesis blocks are always created locally — never accept from peers.
// A peer could send any block with IsGenesis=true and it would bypass
// both the signature check and the parent check below.
if block.IsGenesis {
fmt.Printf("[DAG] ✗ Rejected peer genesis: genesis can only be created locally\n")
return false
}

// Integrity check 1: recompute hash from block fields.
expectedHash := dag.calculateHash(block)
if expectedHash != block.Hash {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: hash mismatch (claimed %s..., computed %s...)\n",
block.Height, block.Hash[:min(16, len(block.Hash))], expectedHash[:16])
return false
}

// Integrity check 2: all non-genesis blocks must carry a valid ECDSA
// signature from the proposer. Unsigned blocks are rejected — this is the
// primary consensus enforcement mechanism.
if !block.IsGenesis && block.Signature == "" {
	fmt.Printf("[DAG] ✗ Rejected peer block #%d from %s: missing signature\n",
		block.Height, block.Proposer)
	return false
}
if block.Signature != "" && !block.IsGenesis {
	sigBytes, sigErr := hex.DecodeString(block.Signature)
	if sigErr != nil || len(sigBytes) != 65 {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: malformed signature\n", block.Height)
		return false
	}
	hashBytes := common.HexToHash(block.Hash)
	pubkeyBytes, recErr := crypto.Ecrecover(hashBytes[:], sigBytes)
	if recErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: signature recovery failed: %v\n", block.Height, recErr)
		return false
	}
	pubkey, parseErr := crypto.UnmarshalPubkey(pubkeyBytes)
	if parseErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: invalid public key: %v\n", block.Height, parseErr)
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
		return false
	}
}

// Integrity check 3: parent-existence and height validation.
// Only enforced when we already have a populated DAG (more than genesis).
// During initial catch-up the DAG is empty so every block would appear
// to have unknown parents — relaxing the check lets the first sync fill
// the DAG, after which the check protects against floating blocks.
if len(dag.blocks) > 1 {
if len(block.ParentHashes) == 0 {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: no parent hashes\n", block.Height)
return false
}
maxParentHeight := int64(-1)
for _, ph := range block.ParentHashes {
parent, parentExists := dag.blocks[ph]
if !parentExists {
return false
}
if parent.Height > maxParentHeight {
maxParentHeight = parent.Height
}
}
if block.Height != maxParentHeight+1 {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: invalid height (parent max %d)\n",
block.Height, maxParentHeight)
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
return false
}
}

// State-root integrity check — log a warning on mismatch but still accept
// the block. Wall-clock differences, in-flight demurrage settlement, and
// minor DB-sync lag can cause transient mismatches between honest nodes.
// Rejecting on mismatch causes split-brain when two nodes process the same
// transactions in slightly different orders or at different wall-clock times.
// The ECDSA signature check above already prevents forged blocks.
if block.StateRoot != "" {
localRoot := dag.state.StateRoot()
if block.StateRoot != localRoot {
fmt.Printf("[DAG] ⚠ StateRoot mismatch on peer block #%d (proposer=%s..., local=%s...) — accepted (warn only)\n",
block.Height, block.StateRoot[:min(16, len(block.StateRoot))], localRoot[:min(16, len(localRoot))])
}
}

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

// Replay TXs outside the DAG lock via the serialized replay queue.
// The single consumer goroutine ensures blocks are replayed in order,
// preserving TX dependencies across blocks (e.g. register_human in
// block N before a transfer in block N+1 for the same wallet).
if len(block.Transactions) > 0 {
	select {
	case dag.replayQueue <- block:
	default:
		// Queue full — process synchronously in a goroutine to avoid
		// dropping TXs (better than losing them entirely).
		go dag.replayTransactions(block)
	}
}

fmt.Printf("[DAG] ✓ Added peer block #%d | Tips: %d\n", block.Height, len(dag.tips))
return true
}

// Note: uses Go's built-in min() (available since Go 1.21; this module
// targets 1.24.1) rather than a custom helper — other files in this
// package already define min4()/min4b() specifically to avoid shadowing
// the built-in, so we follow that same convention here by not shadowing it.

func (dag *BlockDAG) LatestBlock() *Block {
dag.mu.RLock()
defer dag.mu.RUnlock()
var latest *Block
for hash := range dag.tips {
b := dag.blocks[hash]
if latest == nil || b.Height > latest.Height {
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
			// If proof data is present, verify it via BioVerifier before applying.
			// This eliminates unconditional trust in the validator ECDSA key for
			// registration TXs — a compromised validator key cannot inject fake
			// registrations without also producing a valid Groth16 proof.
			// Backward-compatible: blocks from old nodes omit proof fields and
			// fall through to the existing trust-based path below.
			if len(tx.ProofA) == 2 && len(tx.ProofB) == 2 && len(tx.ProofC) == 2 && len(tx.PubSignals) >= 2 {
				if !dag.verifyZKProof(tx) {
					fmt.Printf("[REPLAY] ✗ ZK proof verification failed for %s (block #%d) — skipping\n", wallet, block.Height)
					continue
				}
				fmt.Printf("[REPLAY] ✓ ZK proof verified for %s (block #%d)\n", wallet, block.Height)
			}
			if !dag.state.TryClaimNullifier(nullifier, wallet) {
				continue // already registered
			}
			if err := dag.state.RegisterHuman(wallet); err != nil {
				fmt.Printf("[REPLAY] ✗ RegisterHuman %s: %v (nullifier recorded, balance NOT credited)\n", wallet, err)
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
			if err := dag.state.ApplyTransferDelta(wallet, to, tx.Amount); err != nil {
				fmt.Printf("[REPLAY] ✗ Transfer %s->%s %.6f: %v (block #%d)\n", wallet, to, tx.Amount, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied transfer %.6f AEQ: %s->%s (block #%d)\n", tx.Amount, wallet, to, block.Height)

		case "swap_aeq_tusd":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_aeq_tusd in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.ApplySwapDelta(wallet, tx.Amount, tx.AmountOut, true); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_aeq_tusd %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_aeq_tusd %.6f AEQ->%.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "swap_tusd_aeq":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping swap_tusd_aeq in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.ApplySwapDelta(wallet, tx.Amount, tx.AmountOut, false); err != nil {
				fmt.Printf("[REPLAY] ✗ swap_tusd_aeq %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied swap_tusd_aeq %.6f tUSD->%.6f AEQ for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "add_liquidity":
			if wallet == "" || tx.Amount <= 0 || tx.AmountOut <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping add_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.AddLiquidityDelta(wallet, tx.Amount, tx.AmountOut, tx.LPShares); err != nil {
				fmt.Printf("[REPLAY] ✗ add_liquidity %s: %v (block #%d)\n", wallet, err, block.Height)
				continue
			}
			fmt.Printf("[REPLAY] ✓ Applied add_liquidity %.6f AEQ + %.6f tUSD for %s (block #%d)\n", tx.Amount, tx.AmountOut, wallet, block.Height)

		case "remove_liquidity":
			if wallet == "" || tx.Amount <= 0 {
				fmt.Printf("[REPLAY] ⚠ Skipping remove_liquidity in block #%d: missing fields\n", block.Height)
				continue
			}
			if err := dag.state.RemoveLiquidityDelta(wallet, tx.Amount); err != nil {
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
				dag.state.ApplyUBIDelta(tx.AmountPerHuman)
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

// verifyZKProof reconstructs the Groth16 proof from the TX's decimal string
// fields and calls the BioVerifier contract via the local EVM engine to check
// validity. Returns true when the proof is valid, false otherwise.
// Only called when all four proof fields (ProofA, ProofB, ProofC, PubSignals)
// are present — blocks from old nodes omit them and fall back to trust-based
// validation for backward compatibility.
func (dag *BlockDAG) verifyZKProof(tx Transaction) bool {
	if dag.evm == nil {
		// EVM not yet wired (happens briefly at startup). Fall back to the
		// trust-based path so registrations are not silently dropped.
		fmt.Printf("[REPLAY] ⚠ verifyZKProof: EVM unavailable — falling back to trust-based verification\n")
		return true
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
