package keeper

import (
"bytes"
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"strings"
"sync"
"time"
)

// ─── PEER REGISTRY ───────────────────────────────────────────────────────────

// PeerRegistry tracks known peer nodes with heartbeat timestamps.
// The primary node (IS_PRIMARY_NODE=true) collects registrations from
// secondary nodes; secondary nodes query it to discover all peers.
var GlobalPeerRegistry = &PeerRegistry{peers: make(map[string]time.Time)}

type PeerRegistry struct {
mu    sync.RWMutex
peers map[string]time.Time // URL → last heartbeat
}

func (pr *PeerRegistry) Register(url string) {
if url == "" {
return
}
url = strings.TrimRight(url, "/")
pr.mu.Lock()
pr.peers[url] = time.Now()
pr.mu.Unlock()
}

// ActivePeers returns peers that sent a heartbeat in the last 5 minutes,
// excluding selfURL so a node never syncs with itself.
func (pr *PeerRegistry) ActivePeers(selfURL string) []string {
pr.mu.RLock()
defer pr.mu.RUnlock()
self := strings.TrimRight(selfURL, "/")
var result []string
for url, lastSeen := range pr.peers {
if time.Since(lastSeen) < 5*time.Minute && url != self {
result = append(result, url)
}
}
return result
}

func (pr *PeerRegistry) AllPeers() []string {
pr.mu.RLock()
defer pr.mu.RUnlock()
result := make([]string, 0, len(pr.peers))
for url := range pr.peers {
result = append(result, url)
}
return result
}

// ─── BLOCK SYNC ──────────────────────────────────────────────────────────────

var httpSyncClient = &http.Client{Timeout: 30 * time.Second}

// startSyncForPeer starts a long-running syncWithNode goroutine for peerURL.
// If a goroutine is already running for that URL, this is a no-op.
func (dag *BlockDAG) startSyncForPeer(peerURL string) {
peerURL = strings.TrimRight(peerURL, "/")
dag.syncPeerMu.Lock()
already := dag.activeSyncPeers[peerURL]
if !already {
dag.activeSyncPeers[peerURL] = true
}
dag.syncPeerMu.Unlock()
if already {
return
}
go func() {
defer func() {
dag.syncPeerMu.Lock()
delete(dag.activeSyncPeers, peerURL)
dag.syncPeerMu.Unlock()
}()
dag.syncWithNode(peerURL)
}()
}

func (dag *BlockDAG) syncWithNode(nodeURL string) {
ticker := time.NewTicker(6 * time.Second)
for range ticker.C {
resp, err := httpSyncClient.Get(nodeURL + "/api/blocks")
if err != nil {
continue
}
body, _ := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
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

// ─── PEER DISCOVERY ──────────────────────────────────────────────────────────

// StartPeerDiscovery handles automatic peer registration and discovery.
//
// Environment variables:
//   SELF_URL          — this node's own public URL (required for registration)
//   PRIMARY_NODE_URL  — primary node to register with (omit on the primary itself)
//   PEER_NODES        — comma-separated static peer list (optional fallback)
//
// Flow for secondary nodes (e.g. Render):
//   1. POST /api/peers/register to the primary with our own URL
//   2. Receive current peer list, start syncing each peer
//   3. Every 30s: repeat to heartbeat + discover new peers
//
// Flow for the primary node (Railway, IS_PRIMARY_NODE=true):
//   - Accepts registrations, serves peer list — no outbound registration needed
func (dag *BlockDAG) StartPeerDiscovery(selfURL string) {
selfURL = strings.TrimRight(selfURL, "/")
primaryURL := strings.TrimRight(os.Getenv("PRIMARY_NODE_URL"), "/")

fmt.Println("── Starting Peer Discovery ──────────────")
if selfURL == "" {
fmt.Println("[PEERS] SELF_URL not set — no peer sync (isolated node)")
return
}
fmt.Printf("[PEERS] Self: %s\n", selfURL)

// Seed from explicit PEER_NODES (backwards compat + manual override)
for _, peer := range staticPeers(selfURL) {
GlobalPeerRegistry.Register(peer)
dag.startSyncForPeer(peer)
fmt.Printf("[PEERS] Static peer: %s\n", peer)
}

if primaryURL != "" && primaryURL != selfURL {
fmt.Printf("[PEERS] Primary: %s\n", primaryURL)
dag.registerAndDiscover(selfURL, primaryURL)
} else {
fmt.Println("[PEERS] Primary node — accepting registrations from peers")
}

go func() {
ticker := time.NewTicker(30 * time.Second)
for range ticker.C {
if primaryURL != "" && primaryURL != selfURL {
dag.registerAndDiscover(selfURL, primaryURL)
}
}
}()
}

// registerAndDiscover POSTs our URL to /api/peers/register on the primary,
// gets back the full peer list, and starts sync goroutines for any new peers.
func (dag *BlockDAG) registerAndDiscover(selfURL, primaryURL string) {
body, _ := json.Marshal(map[string]string{"url": selfURL})
resp, err := httpSyncClient.Post(
primaryURL+"/api/peers/register", "application/json", bytes.NewReader(body))
if err != nil {
fmt.Printf("[PEERS] Could not reach primary %s: %v\n", primaryURL, err)
return
}
defer resp.Body.Close()
var result struct {
Peers []string `json:"peers"`
}
json.NewDecoder(resp.Body).Decode(&result)
for _, peer := range result.Peers {
peer = strings.TrimRight(peer, "/")
if peer == selfURL {
continue
}
GlobalPeerRegistry.Register(peer)
dag.startSyncForPeer(peer)
}
}

// staticPeers reads the PEER_NODES env var for backwards compatibility.
func staticPeers(selfURL string) []string {
raw := os.Getenv("PEER_NODES")
if raw == "" {
return nil
}
var out []string
for _, p := range strings.Split(raw, ",") {
p = strings.TrimSpace(strings.TrimRight(p, "/"))
if p != "" && p != selfURL {
out = append(out, p)
}
}
return out
}

// StartHTTPBlockSync is an alias kept for call-site compatibility.
func (dag *BlockDAG) StartHTTPBlockSync(selfURL string) {
dag.StartPeerDiscovery(selfURL)
}
