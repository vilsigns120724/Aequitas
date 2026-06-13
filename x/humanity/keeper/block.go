package keeper

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"time"
)

type Block struct {
Height    int64  `json:"height"`
Timestamp int64  `json:"timestamp"`
PrevHash  string `json:"prev_hash"`
Hash      string `json:"hash"`
Proposer  string `json:"proposer"`
Humans    int    `json:"humans"`
}

type Blockchain struct {
blocks []*Block
keeper *Keeper
nodeID string
}

func NewBlockchain(keeper *Keeper, nodeID string) *Blockchain {
bc := &Blockchain{
blocks: make([]*Block, 0),
keeper: keeper,
nodeID: nodeID,
}
bc.createGenesisBlock()
return bc
}

func (bc *Blockchain) createGenesisBlock() {
genesis := &Block{
Height:    0,
Timestamp: time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC).Unix(),
PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
Proposer:  "genesis",
Humans:    bc.keeper.TotalHumans(),
}
genesis.Hash = bc.calculateHash(genesis)
bc.blocks = append(bc.blocks, genesis)
fmt.Printf("✓ Genesis Block: %s\n", genesis.Hash[:16]+"...")
}

func (bc *Blockchain) calculateHash(b *Block) string {
data, _ := json.Marshal(map[string]interface{}{
"height":    b.Height,
"timestamp": b.Timestamp,
"prev_hash": b.PrevHash,
"proposer":  b.Proposer,
"humans":    b.Humans,
})
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) ProduceBlock() *Block {
prev := bc.blocks[len(bc.blocks)-1]
block := &Block{
Height:    prev.Height + 1,
Timestamp: time.Now().Unix(),
PrevHash:  prev.Hash,
Proposer:  bc.nodeID,
Humans:    bc.keeper.TotalHumans(),
}
block.Hash = bc.calculateHash(block)
bc.blocks = append(bc.blocks, block)
return block
}

func (bc *Blockchain) LatestBlock() *Block {
return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) Height() int64 {
return bc.blocks[len(bc.blocks)-1].Height
}

func (bc *Blockchain) GetBlocks() []*Block {
return bc.blocks
}
