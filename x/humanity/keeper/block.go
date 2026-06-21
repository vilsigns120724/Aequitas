package keeper

import (
"crypto/ecdsa"
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"os"
"strings"
"sync"
"time"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	Type   string  `json:"type"`
	Wallet string  `json:"wallet"`
	Amount float64 `json:"amount,omitempty"`
	TxHash string  `json:"tx_hash"`
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

type BlockDAG struct {
blocks               map[string]*Block
tips                 map[string]bool
mu                   sync.RWMutex
keeper               *Keeper
state                *ChainState
nodeID               string
height               int64
pendingTxs           []Transaction
txMu                 sync.Mutex
signingKey           *ecdsa.PrivateKey
authorizedValidators map[string]bool // Ethereum addresses allowed to propose blocks
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

func (dag *BlockDAG) AddTransaction(tx Transaction) {
dag.txMu.Lock()
defer dag.txMu.Unlock()
dag.pendingTxs = append(dag.pendingTxs, tx)
}

func NewBlockchain(keeper *Keeper, nodeID string, state *ChainState) *BlockDAG {
dag := &BlockDAG{
blocks:               make(map[string]*Block),
tips:                 make(map[string]bool),
keeper:               keeper,
state:                state,
nodeID:               nodeID,
authorizedValidators: loadAuthorizedValidators(),
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
Timestamp:    time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC).Unix(),
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
txData, _ := json.Marshal(b.Transactions)
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

func (dag *BlockDAG) AddPeerBlock(block *Block) {
dag.mu.Lock()
defer dag.mu.Unlock()

// Skip if already known
if _, exists := dag.blocks[block.Hash]; exists {
return
}

// Integrity check 1: recompute hash from block fields.
expectedHash := dag.calculateHash(block)
if expectedHash != block.Hash {
fmt.Printf("[DAG] ✗ Rejected peer block #%d: hash mismatch (claimed %s..., computed %s...)\n",
block.Height, block.Hash[:min(16, len(block.Hash))], expectedHash[:16])
return
}

// Integrity check 2: all non-genesis blocks must carry a valid ECDSA
// signature from the proposer. Unsigned blocks are rejected — this is the
// primary consensus enforcement mechanism.
if !block.IsGenesis && block.Signature == "" {
	fmt.Printf("[DAG] ✗ Rejected peer block #%d from %s: missing signature\n",
		block.Height, block.Proposer)
	return
}
if block.Signature != "" && !block.IsGenesis {
	sigBytes, sigErr := hex.DecodeString(block.Signature)
	if sigErr != nil || len(sigBytes) != 65 {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: malformed signature\n", block.Height)
		return
	}
	hashBytes := common.HexToHash(block.Hash)
	pubkeyBytes, recErr := crypto.Ecrecover(hashBytes[:], sigBytes)
	if recErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: signature recovery failed: %v\n", block.Height, recErr)
		return
	}
	pubkey, parseErr := crypto.UnmarshalPubkey(pubkeyBytes)
	if parseErr != nil {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: invalid public key: %v\n", block.Height, parseErr)
		return
	}
	recoveredAddr := strings.ToLower(crypto.PubkeyToAddress(*pubkey).Hex())
	proposer := strings.ToLower(block.Proposer)
	// Proposer must be the Ethereum address that produced the signature.
	// Blocks where the proposer field does not match the recovered signing
	// address are unconditionally rejected — no libp2p-nodeID exemption.
	if recoveredAddr != proposer {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: signature mismatch (signer %s, proposer %s)\n",
			block.Height, recoveredAddr, proposer)
		return
	}
	// Proposer must be in the authorized validator set. Without this check
	// anyone can generate an Ethereum key, sign a block, and feed it in.
	if !dag.authorizedValidators[proposer] {
		fmt.Printf("[DAG] ✗ Rejected peer block #%d: proposer %s is not an authorized validator\n",
			block.Height, proposer)
		return
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

fmt.Printf("[DAG] ✓ Added peer block #%d | Tips: %d\n", block.Height, len(dag.tips))
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
for i := 0; i < len(result)-1; i++ {
for j := i + 1; j < len(result); j++ {
if result[i].Height > result[j].Height {
result[i], result[j] = result[j], result[i]
}
}
}
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

func (dag *BlockDAG) ReconstructState(state *ChainState) {
dag.mu.RLock()
defer dag.mu.RUnlock()

count := 0
for _, block := range dag.blocks {
	// Only replay transactions from authorized validators. Unsigned/unknown
	// proposers could have injected register_human entries via forged blocks.
	if !block.IsGenesis {
		proposer := strings.ToLower(block.Proposer)
		if !dag.authorizedValidators[proposer] {
			continue
		}
	}
	for _, tx := range block.Transactions {
		if tx.Type == "register_human" && tx.Wallet != "" {
			if !state.IsHuman(tx.Wallet) {
				state.RegisterHuman(tx.Wallet)
				count++
			}
		}
	}
}
if count > 0 {
fmt.Printf("[CHAIN] ✓ Reconstructed %d registrations from blockchain\n", count)
}
}
