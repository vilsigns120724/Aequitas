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
	// Literal IP: skip DNS entirely and connect directly.
	// net.LookupHost("173.x.x.x") can fail on Alpine/Docker even for valid
	// public IPs because the minimal resolver doesn't handle PTR/A lookups
	// for already-resolved addresses the same way on all platforms.
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
			return nil, fmt.Errorf("connection to private/loopback IP rejected: %s", host)
		}
		d := &net.Dialer{Timeout: 10 * time.Second}
		return d.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
	}
	// Hostname: resolve DNS and verify every IP is public (DNS-rebinding guard).
	ips, err := net.LookupHost(host)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("DNS lookup failed for %s", host)
	}
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed == nil || parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast() || parsed.IsMulticast() || parsed.IsUnspecified() {
			return nil, fmt.Errorf("DNS resolved to private/loopback address %s for host %s", ip, host)
		}
	}
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

// fetchBlocksSince fetches up to `limit` blocks with Height > minHeight from
// nodeURL, via the height-based ?min_height=&limit= pagination on /api/blocks.
func (dag *BlockDAG) fetchBlocksSince(nodeURL string, minHeight int64, limit int) ([]*Block, error) {
	resp, err := httpSyncClient.Get(fmt.Sprintf("%s/api/blocks?min_height=%d&limit=%d", nodeURL, minHeight, limit))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	var blocks []*Block
	if err := json.Unmarshal(body, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

// doSyncOnce walks nodeURL's block history forward from our own current
// height, fetching pageSize-sized pages until it catches up to the peer's
// tip. Returns false on a network/decode error (used by syncWithNode to
// back off), true otherwise — including the normal "nothing new" case.
//
// FIX: this used to page via ?offset=dag.TotalBlocks() — i.e. treat "how
// many blocks do I have" as a position into the PEER's own array. That only
// works as long as both sides accumulate the exact same number of entries
// at the exact same pace, which breaks the moment more than one validator
// produces concurrently (the normal, intended BlockDAG case): each side
// merges multi-parent siblings at its own pace, so the two nodes' local
// block COUNTS drift apart from each other even when both are otherwise
// healthy. A node whose count fell out of step with the peer's array
// position kept requesting the same already-fully-known window forever —
// confirmed in production: a node stuck ~640 blocks behind never advanced,
// growing its own isolated, never-reconciled side chain instead. HEIGHT is
// the one frontier marker that stays meaningful regardless of how many
// duplicate-height siblings either side has — paging by "give me everything
// above the highest height I've already got" can't get stuck the same way.
func (dag *BlockDAG) doSyncOnce(nodeURL string) (ok bool) {
	const pageSize = 500
	const maxPagesPerCall = 2000 // hard cap: 1,000,000 blocks per call — headroom, not unbounded
	// FIX: requesting strictly "height > my own height" misses SIBLINGS at
	// or just below that height that other validators produced. A later
	// block's parent can be one of those siblings (e.g. the peer's own
	// previous block, at a height I already consider "done" because MY
	// own block at that height got there first) — if I never fetched it,
	// AddPeerBlock rejects every subsequent block built on top of it for
	// missing a parent I don't have, forever. Confirmed in production:
	// three concurrent validators each kept seeing ONLY their own
	// single-parent chain in their own /api/blocks output — never each
	// other's blocks — because the height-exclusive fetch never pulled in
	// the sibling branches needed to resolve later parents. Re-requesting
	// a small overlap window of already-"passed" heights each cycle is
	// cheap (AddPeerBlock dedupes by hash) and guarantees sibling forks
	// from other validators get imported before something builds on them.
	const syncOverlap = 20
	totalAdded := 0
	for page := 0; page < maxPagesPerCall; page++ {
		minHeight := dag.Height() - syncOverlap
		if minHeight < 0 {
			minHeight = 0
		}
		blocks, err := dag.fetchBlocksSince(nodeURL, minHeight, pageSize)
		if err != nil {
			fmt.Printf("[HTTP-SYNC] ✗ Could not fetch page (min_height=%d) from %s: %v\n", minHeight, nodeURL, err)
			if page == 0 {
				return false // never even got a first page — treat as a failed sync attempt
			}
			break // got at least one page this call; report what we added
		}
		if len(blocks) == 0 {
			break // caught up — peer has nothing newer than our height
		}
		addedThisPage := 0
		for _, block := range blocks {
			// FIX: genesis is always created locally and AddPeerBlock always
			// rejects a peer-supplied genesis (by design — see its own
			// comment). Without this skip, every single sync cycle forever
			// re-attempts and re-logs "Rejected peer genesis", since it's
			// never marked as "exists" and so never short-circuits like a
			// normal already-known block would.
			if block.IsGenesis {
				continue
			}
			dag.mu.RLock()
			_, exists := dag.blocks[block.Hash]
			dag.mu.RUnlock()
			if !exists && dag.AddPeerBlock(block) {
				addedThisPage++
			}
		}
		totalAdded += addedThisPage
		if len(blocks) < pageSize {
			break // last page (peer's tip is within this page)
		}
		if addedThisPage == 0 {
			// Fetched a full page (peer has more beyond this window) but
			// nothing in it was new/acceptable — e.g. every block was
			// already known, or rejected for a bad signature/parent.
			// Looping again would request virtually the same overlap
			// window and get the exact same page forever.
			fmt.Printf("[HTTP-SYNC] ⚠ Page above height %d added 0 of %d blocks — stopping sync from %s for this cycle\n", minHeight, len(blocks), nodeURL)
			break
		}
	}
	if totalAdded > 0 {
		dag.mu.RLock()
		tipCount := len(dag.tips)
		dag.mu.RUnlock()
		fmt.Printf("[HTTP-SYNC] ✓ Added %d new blocks from %s | DAG tips: %d | height %d\n", totalAdded, nodeURL, tipCount, dag.Height())
	}
	return true
}

func (dag *BlockDAG) syncWithNode(nodeURL string) {
	// Try immediately on first call — no initial delay. doSyncOnce itself
	// pages through full history starting from whatever we already have
	// locally, so this one call already performs the initial catch-up.
	backoff := 6 * time.Second
	dag.doSyncOnce(nodeURL)
	ticker := time.NewTicker(backoff)
	defer ticker.Stop()
	for range ticker.C {
		if !dag.doSyncOnce(nodeURL) {
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			} // max 30s not 60s
			ticker.Reset(backoff)
			continue
		}
		// Reset backoff on success
		if backoff > 6*time.Second {
			backoff = 6 * time.Second
			ticker.Reset(backoff)
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
// Flow for secondary nodes (Railway/VPS/self-hosted):
//   1. POST /api/peers/register to the primary with our own URL
//   2. Receive current peer list, start syncing each peer
//   3. Every 30s: repeat to heartbeat + discover new peers
//
// Flow for the primary node (IS_PRIMARY_NODE=true):
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
		"url":                  selfURL,
		"signing_address":      signerAddr,
		"peer_secret":          os.Getenv("PEER_SECRET"),
		"node_operator_wallet": strings.ToLower(os.Getenv("NODE_OPERATOR_WALLET")),
	})
	resp, err := httpSyncClient.Post(
		primaryURL+"/api/peers/register", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("[PEERS] Could not reach primary %s: %v\n", primaryURL, err)
		return
	}
	defer resp.Body.Close()
	// FIX: this used to decode the response body unconditionally, regardless
	// of HTTP status. If the primary rejects registration (e.g. 403 because
	// NODE_OPERATOR_WALLET isn't a registered human yet), the body is
	// {"error":"..."} — decoding that into {Peers, Validators} silently
	// yields two empty slices with no error. The node then never learns the
	// primary's proposer address as an authorized validator, so every block
	// from the primary gets rejected by AddPeerBlock's "not an authorized
	// validator" check forever — visible only as "stuck at block 1", with no
	// indication why, since the actual rejection reason (this 403) was never
	// logged on the secondary's side at all.
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[PEERS] ✗ Registration with primary %s rejected (HTTP %d): %s — this node will NOT sync any blocks until registration succeeds (check NODE_OPERATOR_WALLET is a registered human, and PEER_SECRET/signature)\n",
			primaryURL, resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
		return
	}
	var result struct {
		Peers      []string `json:"peers"`
		Validators []string `json:"validators"`
	}
	json.Unmarshal(bodyBytes, &result)

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
		return !ip.IsLoopback() && !ip.IsPrivate() && !ip.IsLinkLocalUnicast() && !ip.IsMulticast() && !ip.IsUnspecified()
	}

	// Hostname: require HTTPS to prevent DNS-rebinding attacks.
	if u.Scheme != "https" {
		return false
	}

	// FIX 10: DNS lookup removed from isAllowedPeerURL to eliminate TOCTOU race
	// (DNS may resolve differently at connect time vs. check time, enabling
	// DNS-rebinding). The actual IP validation is authoritative in pinningDialer,
	// which resolves DNS once and pins the connection to the resolved IP.
	// String-level checks for obviously private literal IPs are still done above.
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
