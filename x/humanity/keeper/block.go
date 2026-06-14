package keeper

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"sync"
"time"
)

type Block struct {
Height       int64    `json:"height"`
Timestamp    int64    `json:"timestamp"`
ParentHashes []string `json:"parent_hashes"` // DAG: multiple parents
Hash         string   `json:"hash"`
Proposer     string   `json:"proposer"`
Humans       int      `json:"humans"`
IsGenesis    bool     `json:"is_genesis,omitempty"`
}

type BlockDAG struct {
blocks   map[string]*Block // hash -> block
tips     []string          // current DAG tips (blocks with no children)
mu       sync.RWMutex
keeper   *Keeper
nodeID   string
height   int64
}

func NewBlockchain(keeper *Keeper, nodeID string) *BlockDAG {
dag := &BlockDAG{
blocks: make(map[string]*Block),
tips:   make([]string, 0),
keeper: keeper,
nodeID: nodeID,
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
Humans:       dag.keeper.TotalHumans(),
IsGenesis:    true,
}
genesis.Hash = dag.calculateHash(genesis)
dag.blocks[genesis.Hash] = genesis
dag.tips = []string{genesis.Hash}
dag.height = 0
fmt.Printf("✓ Genesis Block (DAG): %s\n", genesis.Hash[:16]+"...")
}

func (dag *BlockDAG) calculateHash(b *Block) string {
data, _ := json.Marshal(map[string]interface{}{
"height":        b.Height,
"timestamp":     b.Timestamp,
"parent_hashes": b.ParentHashes,
"proposer":      b.Proposer,
"humans":        b.Humans,
})
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:])
}

func (dag *BlockDAG) ProduceBlock() *Block {
dag.mu.Lock()
defer dag.mu.Unlock()

// New block references ALL current tips (DAG merge)
parentHashes := make([]string, len(dag.tips))
copy(parentHashes, dag.tips)

// Height = max parent height + 1
maxParentHeight := int64(0)
for _, ph := range parentHashes {
if parent, ok := dag.blocks[ph]; ok {
if parent.Height > maxParentHeight {
maxParentHeight = parent.Height
}
}
}

block := &Block{
Height:       maxParentHeight + 1,
Timestamp:    time.Now().Unix(),
ParentHashes: parentHashes,
Proposer:     dag.nodeID,
Humans:       dag.keeper.TotalHumans(),
}
block.Hash = dag.calculateHash(block)

dag.blocks[block.Hash] = block
dag.tips = []string{block.Hash} // this block becomes the new tip
dag.height = block.Height

return block
}

func (dag *BlockDAG) LatestBlock() *Block {
dag.mu.RLock()
defer dag.mu.RUnlock()
if len(dag.tips) == 0 {
return nil
}
return dag.blocks[dag.tips[0]]
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
// Sort by height
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
return dag.tips
}

// AddPeerBlock adds a block received from a peer as a new DAG tip
func (dag *BlockDAG) AddPeerBlock(block *Block) {
dag.mu.Lock()
defer dag.mu.Unlock()

// Skip if already known
if _, exists := dag.blocks[block.Hash]; exists {
return
}

dag.blocks[block.Hash] = block

// Add as tip if not already referenced
dag.tips = append(dag.tips, block.Hash)

if block.Height > dag.height {
dag.height = block.Height
}

fmt.Printf("[DAG] Added peer block #%d | Tips: %d\n", block.Height, len(dag.tips))
}
