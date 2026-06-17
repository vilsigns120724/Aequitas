package keeper

import (
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"strings"
"time"
)

// KnownNodes lists peer node URLs to sync blocks with over HTTP.
// Defaults to the original hardcoded list, but can be overridden/extended
// via the PEER_NODES environment variable (comma-separated URLs), so a new
// node can join the network by setting an env var on Railway/Render
// instead of requiring a code change and redeploy of every existing node.
var KnownNodes = loadKnownNodes()

func loadKnownNodes() []string {
defaults := []string{
"https://aequitas-node-2.onrender.com",
}
extra := os.Getenv("PEER_NODES")
if extra == "" {
return defaults
}
seen := make(map[string]bool)
nodes := make([]string, 0, len(defaults)+4)
for _, n := range defaults {
if !seen[n] {
seen[n] = true
nodes = append(nodes, n)
}
}
for _, n := range strings.Split(extra, ",") {
n = strings.TrimSpace(n)
if n != "" && !seen[n] {
seen[n] = true
nodes = append(nodes, n)
}
}
return nodes
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
