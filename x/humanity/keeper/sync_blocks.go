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
// Configured exclusively via the PEER_NODES environment variable
// (comma-separated URLs). No hardcoded defaults — only nodes explicitly
// added here participate in block sync, so an outdated peer running a
// different hash format cannot pollute the DAG with rejected blocks.
var KnownNodes = loadKnownNodes()

func loadKnownNodes() []string {
extra := os.Getenv("PEER_NODES")
if extra == "" {
return nil
}
seen := make(map[string]bool)
nodes := make([]string, 0, 4)
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

var httpSyncClient = &http.Client{Timeout: 30 * time.Second}

func (dag *BlockDAG) syncWithNode(nodeURL string) {
ticker := time.NewTicker(6 * time.Second)
for range ticker.C {
resp, err := httpSyncClient.Get(nodeURL + "/api/blocks")
if err != nil {
continue
}
body, _ := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10 MB response cap
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
