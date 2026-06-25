package keeper

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net"
"net/http"
"net/url"
"os"
"sort"
"strings"
"sync"
"time"

"github.com/ethereum/go-ethereum/crypto"
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

// pinningDialer resolves the hostname once, verifies all IPs are public,
// then connects directly to the first IP — bypassing any subsequent DNS
// re-resolution that DNS-rebinding attacks rely on.
func pinningDialer(ctx context.Context, network, addr string) (net.Conn, error) {
host, port, err := net.SplitHostPort(addr)
if err != nil {
return nil, err
}
ips, err := net.LookupHost(host)
if err != nil || len(ips) == 0 {
return nil, fmt.Errorf("DNS lookup failed for %s", host)
}
for _, ip := range ips {
parsed := net.ParseIP(ip)
if parsed == nil || parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast() {
return nil, fmt.Errorf("DNS resolved to private/loopback address %s for host %s", ip, host)
}
}
// Pin to first resolved IP so no second lookup can redirect to a private address.
d := &net.Dialer{Timeout: 10 * time.Second}
return d.DialContext(ctx, network, net.JoinHostPort(ips[0], port))
}

var httpSyncClient = &http.Client{
Timeout: 30 * time.Second,
// Never follow redirects — a public URL could redirect internally after our check.
CheckRedirect: func(req *http.Request, via []*http.Request) error {
return http.ErrUseLastResponse
},
Transport: &http.Transport{
DialContext: pinningDialer,
},
}

const maxSyncPeers = 20

// startSyncForPeer starts a long-running syncWithNode goroutine for peerURL.
// No-op if already syncing that URL or if the peer cap is reached.
func (dag *BlockDAG) startSyncForPeer(peerURL string) {
peerURL = strings.TrimRight(peerURL, "/")
if !isAllowedPeerURL(peerURL) {
fmt.Printf("[PEERS] Rejected peer URL (must be public HTTPS): %s\n", peerURL)
return
}
dag.syncPeerMu.Lock()
already := dag.activeSyncPeers[peerURL]
tooMany := len(dag.activeSyncPeers) >= maxSyncPeers
if !already && !tooMany {
dag.activeSyncPeers[peerURL] = true
}
dag.syncPeerMu.Unlock()
if already || tooMany {
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
// P3-1: exponential backoff on error — doubles up to 60s, resets on success.
backoff := 6 * time.Second
ticker := time.NewTicker(backoff)
defer ticker.Stop()
for range ticker.C {
resp, err := httpSyncClient.Get(nodeURL + "/api/blocks")
if err != nil {
backoff *= 2
if backoff > 60*time.Second { backoff = 60 * time.Second }
ticker.Reset(backoff)
continue
}
body, _ := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
resp.Body.Close()

var blocks []*Block
if err := json.Unmarshal(body, &blocks); err != nil {
backoff *= 2
if backoff > 60*time.Second { backoff = 60 * time.Second }
ticker.Reset(backoff)
continue
}
// P2-AUDIT: Reject unreasonably large block lists from peers. A legitimate
// node syncing 6-second blocks produces at most ~14400 blocks/day; a
// response with more than 50000 blocks is either a buggy peer or an attack.
const maxBlocksPerSync = 50000
if len(blocks) > maxBlocksPerSync {
fmt.Printf("[HTTP-SYNC] Peer %s returned %d blocks (max %d) -- skipping\n", nodeURL, len(blocks), maxBlocksPerSync)
backoff *= 2
if backoff > 60*time.Second { backoff = 60 * time.Second }
ticker.Reset(backoff)
continue
}
// Reset backoff on success
if backoff > 6*time.Second { backoff = 6 * time.Second; ticker.Reset(backoff) }

sort.Slice(blocks, func(i, j int) bool { return blocks[i].Height < blocks[j].Height })
added := 0
for _, block := range blocks {
dag.mu.RLock()
_, exists := dag.blocks[block.Hash]
dag.mu.RUnlock()
if !exists && dag.AddPeerBlock(block) { added++ }
}
if added > 0 {
	// P2-FIX: dag.tips is a map protected by dag.mu; reading it
	// without any lock is a data race. Snapshot count under RLock.
	dag.mu.RLock()
	tipCount := len(dag.tips)
	dag.mu.RUnlock()
fmt.Printf("[HTTP-SYNC] ✓ Added %d new blocks from %s | DAG tips: %d\n", added, nodeURL, tipCount)
}
} // end for range ticker.C
} // end syncWithNode

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
// The primary never includes itself in its own peer list (/api/peers
// only contains registered secondary nodes). Start syncing from it
// directly so secondary nodes always receive primary blocks.
dag.startSyncForPeer(primaryURL)
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

// registerAndDiscover POSTs our URL and signing address to the primary's
// /api/peers/register. The primary adds our signing address to its authorized
// validator set so our blocks are accepted without any manual configuration.
// We receive the peer list and the current authorized validator addresses back.
func (dag *BlockDAG) registerAndDiscover(selfURL, primaryURL string) {
signerAddr := ""
if dag.signingKey != nil {
signerAddr = strings.ToLower(crypto.PubkeyToAddress(dag.signingKey.PublicKey).Hex())
}
body, _ := json.Marshal(map[string]string{
"url":             selfURL,
"signing_address": signerAddr,
"peer_secret":     os.Getenv("PEER_SECRET"),
})
resp, err := httpSyncClient.Post(
primaryURL+"/api/peers/register", "application/json", bytes.NewReader(body))
if err != nil {
fmt.Printf("[PEERS] Could not reach primary %s: %v\n", primaryURL, err)
return
}
defer resp.Body.Close()
var result struct {
Peers      []string `json:"peers"`
Validators []string `json:"validators"`
}
json.NewDecoder(resp.Body).Decode(&result)

// Add newly discovered authorized validators to our local set so we
// accept blocks from them without requiring AUTHORIZED_VALIDATORS env var.
dag.mu.Lock()
for _, addr := range result.Validators {
addr = strings.ToLower(strings.TrimSpace(addr))
if addr != "" && !dag.authorizedValidators[addr] {
dag.authorizedValidators[addr] = true
fmt.Printf("[PEERS] Auto-authorized validator: %s\n", addr)
}
}
dag.mu.Unlock()

for _, peer := range result.Peers {
peer = strings.TrimRight(peer, "/")
if peer == selfURL {
continue
}
GlobalPeerRegistry.Register(peer)
dag.startSyncForPeer(peer)
}
}

// isAllowedPeerURL returns true for URLs pointing to public IP addresses.
// HTTPS is preferred; HTTP is accepted for literal public IP addresses
// (e.g. http://173.249.37.118:8080) so VPS nodes without a domain can
// participate. HTTP with a hostname is still rejected (DNS-rebinding risk).
func isAllowedPeerURL(rawURL string) bool {
u, err := url.Parse(rawURL)
if err != nil || u.Host == "" {
return false
}
if u.Scheme != "https" && u.Scheme != "http" {
return false
}
host := u.Hostname()
if host == "" || host == "0.0.0.0" || host == "[::]" {
return false
}

// Literal IP: allow HTTP or HTTPS as long as the IP is public.
if ip := net.ParseIP(host); ip != nil {
return !ip.IsLoopback() && !ip.IsPrivate() && !ip.IsLinkLocalUnicast()
}

// Hostname: require HTTPS to prevent DNS-rebinding attacks.
if u.Scheme != "https" {
return false
}

// For hostnames: resolve DNS and verify every returned IP is public.
addrs, err := net.LookupHost(host)
if err != nil || len(addrs) == 0 {
return false
}
for _, addr := range addrs {
ip := net.ParseIP(addr)
if ip == nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
return false
}
}
return true
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
