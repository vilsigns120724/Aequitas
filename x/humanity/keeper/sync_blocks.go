package keeper

import (
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

var KnownNodes = []string{
"https://aequitas-node-2.onrender.com",
}

func (dag *BlockDAG) StartHTTPBlockSync(selfURL string) {
fmt.Println("── Starting HTTP Block Sync ─────────────")
for _, nodeURL := range KnownNodes {
if nodeURL == selfURL {
continue
}
fmt.Printf("✓ Syncing with: %s\n", nodeURL)
go dag.syncWithNode(nodeURL)
}
}

func (dag *BlockDAG) syncWithNode(nodeURL string) {
ticker := time.NewTicker(6 * time.Second)
for range ticker.C {
resp, err := http.Get(nodeURL + "/api/blocks")
if err != nil {
continue
}
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()

var blocks []*Block
if err := json.Unmarshal(body, &blocks); err != nil {
continue
}

added := 0
for _, block := range blocks {
dag.mu.RLock()
_, exists := dag.blocks[block.Hash]
dag.mu.RUnlock()
if !exists {
dag.AddPeerBlock(block)
added++
}
}
if added > 0 {
fmt.Printf("[HTTP-SYNC] ✓ Added %d new blocks from %s | DAG tips: %d\n",
added, nodeURL, len(dag.tips))
}
}
}
