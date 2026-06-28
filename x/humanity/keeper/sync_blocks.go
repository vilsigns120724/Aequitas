package keeper

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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
	// FIX: this used to decode resp.Body unconditionally regardless of HTTP
	// status — a 500/403/429/HTML error page got handed to json.Unmarshal
	// and surfaced as an opaque decode error indistinguishable from "peer
	// sent malformed JSON". Checking the status explicitly means operators
	// see the real cause (e.g. "peer returned 503") instead of a generic
	// decode failure during exactly the kind of outage/drift situation
	// where that distinction matters most.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("peer returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	var blocks []*Block
	if err := json.Unmarshal(body, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

// fetchBlockByHash fetches exactly one block by hash from nodeURL's
// /api/block endpoint, or nil if the peer genuinely doesn't have it (404).
// Any OTHER non-200 (500/403/429/...) is returned as an error rather than
// silently treated the same as "not found" — those mean something is wrong
// with the peer or the request, not that the block doesn't exist, and
// fetchMissingAncestors' caller should be able to tell the difference in
// its logs instead of just "block not found" for every failure mode.
func (dag *BlockDAG) fetchBlockByHash(nodeURL, hash string) (*Block, error) {
	resp, err := httpSyncClient.Get(fmt.Sprintf("%s/api/block?hash=%s", nodeURL, hash))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("peer returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var block Block
	if err := json.Unmarshal(body, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// fetchMissingAncestors resolves orphaned blocks by walking backward one
// specific hash at a time, instead of waiting for them to fall inside
// doSyncOnce's height-windowed (?min_height=) pagination.
//
// FIX: that height window only ever looks near THIS node's own current
// frontier (dag.Height()-syncOverlap). Once a node's chain has drifted from
// a peer's by more than the overlap window — which can start from any
// transient gap, however brief — the actual common-ancestor blocks needed
// to bridge the two chains permanently fall outside that window and are
// never fetched again: every later block from that peer queues as an
// orphan whose missing parent doSyncOnce will never ask for. Confirmed in
// production: cd20 and a VPS secondary both briefly merged with the primary
// right after first connecting (small gap, within the 20-block overlap),
// then permanently regressed to fully isolated single-parent chains once
// that gap — for whatever transient reason — exceeded 20 blocks. This walks
// directly from "what hash is missing" instead of "what height window might
// contain it", so it has no such ceiling: each resolved ancestor's own
// AddPeerBlock call may reveal a further ancestor still missing, which gets
// picked up on the next call here (capped at maxAncestorFetchPerCycle per
// call to bound a single cycle's work, not the total depth reachable across
// repeated calls). Re-snapshots the orphan set after each batch so a chain
// of N missing ancestors in a row gets walked all the way back to a known
// block within a single call, not one hop per call.
// triggerOrphanResolve runs fetchMissingAncestors against every peer this
// node is currently syncing with, in parallel, right now — instead of
// waiting for each peer's own up-to-6s syncWithNode ticker to come around.
//
// Coordination: at most one resolve pass runs at a time (orphanResolveInFlight),
// since concurrent passes would just duplicate the same peer requests. If a
// new orphan triggers this while a pass is already running, that arrival is
// recorded (orphanResolveAgain) and the in-flight pass loops once more
// immediately after finishing — so a burst of orphans arriving faster than
// one pass can complete still gets a fresh attempt covering all of them,
// rather than silently relying on the next periodic tick.
func (dag *BlockDAG) triggerOrphanResolve() {
	dag.orphanResolveMu.Lock()
	if dag.orphanResolveInFlight {
		dag.orphanResolveAgain = true
		dag.orphanResolveMu.Unlock()
		return
	}
	dag.orphanResolveInFlight = true
	dag.orphanResolveMu.Unlock()

	for {
		dag.syncPeerMu.Lock()
		peers := make([]string, 0, len(dag.activeSyncPeers))
		for p := range dag.activeSyncPeers {
			peers = append(peers, p)
		}
		dag.syncPeerMu.Unlock()

		var wg sync.WaitGroup
		for _, peerURL := range peers {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				dag.fetchMissingAncestors(p)
			}(peerURL)
		}
		wg.Wait()

		dag.orphanResolveMu.Lock()
		if !dag.orphanResolveAgain {
			dag.orphanResolveInFlight = false
			dag.orphanResolveMu.Unlock()
			return
		}
		dag.orphanResolveAgain = false
		dag.orphanResolveMu.Unlock()
		// loop again — another orphan arrived mid-pass
	}
}

func (dag *BlockDAG) fetchMissingAncestors(nodeURL string) {
	const maxAncestorFetchPerCycle = 200
	fetched := 0
	for fetched < maxAncestorFetchPerCycle {
		hashes := dag.MissingParentHashes()
		if len(hashes) == 0 {
			return
		}
		progressed := false
		for _, hash := range hashes {
			if fetched >= maxAncestorFetchPerCycle {
				break
			}
			if dag.GetBlockByHash(hash) != nil {
				continue
			}
			if !dag.shouldAttemptFetch(hash) {
				// Tried this exact hash too recently (orphanFetchCooldown) —
				// skip it this pass instead of re-hitting every peer for a
				// hash that just failed moments ago. See orphanFetchCooldown.
				continue
			}
			block, err := dag.fetchBlockByHash(nodeURL, hash)
			if err != nil {
				fmt.Printf("[HTTP-SYNC] ✗ Could not fetch missing ancestor %s from %s: %v\n", hash[:min(16, len(hash))], nodeURL, err)
				continue
			}
			if block == nil {
				continue // peer genuinely doesn't have this block (404)
			}
			fetched++
			progressed = true
			dag.AddPeerBlock(block)
		}
		if !progressed {
			return // peer doesn't have any of these either — stop for this cycle
		}
	}
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

	// Resolve any orphans (this cycle's or earlier ones) by fetching their
	// specific missing-parent hash directly — see fetchMissingAncestors for
	// why the height-windowed pagination above can't reach them once the
	// gap exceeds syncOverlap.
	dag.fetchMissingAncestors(nodeURL)
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
//   PRIMARY_NODE_URL  — primary node to register with (omit on the primary itself).
//                        Daily pool distribution (UBI/validator/LP/escrow, see
//                        main.go) is gated separately by DISTRIBUTION_ENABLED,
//                        NOT by this variable — a node missing PRIMARY_NODE_URL
//                        used to silently self-identify as the distribution
//                        authority, which is exactly the duplicate-distribution
//                        failure class this whole mechanism exists to prevent.
//                        Set DISTRIBUTION_ENABLED=true on exactly one node.
//   PEER_NODES        — comma-separated static peer list (optional fallback)
//
// Flow for secondary nodes (Railway/VPS/self-hosted):
//   1. POST /api/peers/register to the primary with our own URL
//   2. Receive current peer list, start syncing each peer
//   3. Every 30s: repeat to heartbeat + discover new peers
//
// Flow for the primary node (IS_PRIMARY_NODE=true):
//   - Accepts registrations, serves peer list — no outbound registration needed
// NormalizeNodeURL prepends "https://" to rawURL if it has no http(s) scheme.
// Several hosting providers' "public domain" variables (e.g. Railway's
// RAILWAY_PUBLIC_DOMAIN) are bare hostnames with no scheme; SELF_URL/
// PRIMARY_NODE_URL set from those would otherwise fail isAllowedPeerURL's
// "must be public HTTPS" check and silently break peer registration.
func NormalizeNodeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return rawURL
	}
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	return "https://" + rawURL
}

func (dag *BlockDAG) StartPeerDiscovery(selfURL string) {
	selfURL = strings.TrimRight(NormalizeNodeURL(selfURL), "/")
	primaryURL := strings.TrimRight(NormalizeNodeURL(os.Getenv("PRIMARY_NODE_URL")), "/")

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
// fetchAndSignPeerChallenge implements the same manual flow the node
// operator UI already documents (GET a challenge, sign it, send the
// signature back) automatically: fetches a fresh challenge for signerAddr
// from primaryURL and signs it with signingKey using the personal_sign /
// "Ethereum Signed Message" scheme VerifyPeerChallenge expects (see
// block.go). Returns "" (not an error) on any failure — signing is a
// best-effort upgrade over PEER_SECRET, not a hard requirement, since a
// node with PEER_SECRET configured should keep working exactly as before
// even if the challenge round-trip fails for some transient reason.
func fetchAndSignPeerChallenge(primaryURL, signerAddr string, signingKey *ecdsa.PrivateKey) string {
	if signingKey == nil || signerAddr == "" {
		return ""
	}
	resp, err := httpSyncClient.Get(primaryURL + "/api/peers/challenge?address=" + signerAddr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	var result struct {
		Challenge string `json:"challenge"`
	}
	if err := json.Unmarshal(body, &result); err != nil || result.Challenge == "" {
		return ""
	}
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(result.Challenge), result.Challenge)
	hash := crypto.Keccak256Hash([]byte(msg))
	sig, err := crypto.Sign(hash.Bytes(), signingKey)
	if err != nil {
		return ""
	}
	return "0x" + hex.EncodeToString(sig)
}

// registerAndDiscover POSTs our URL and signing address to the primary's
// /api/peers/register, automatically proving private-key ownership via a
// signed challenge (see fetchAndSignPeerChallenge) rather than relying on
// PEER_SECRET alone.
//
// FIX (decentralization): this used to send peer_secret but never a
// signature, even though api.go's handlePeerRegister already accepts EITHER
// a matching PEER_SECRET OR a valid challenge-response signature
// (secretOK || sigOK) — the signature path existed only for the
// manual/operator-documented flow, never wired into the automatic one. In
// practice this meant a single shared secret was the ONLY thing that
// determined whether a new node could join: leaking it lets anyone register
// as a peer, and rotating it (e.g. after a leak) breaks every legitimate
// node's auto-join until each one is individually updated. Every node that
// can sign its own blocks (RELAYER_PRIVATE_KEY, required for block
// production anyway) can now prove ownership of its signing address the
// same way the manual flow always could, making PEER_SECRET an optional
// bootstrap fallback instead of the only practical path.
func (dag *BlockDAG) registerAndDiscover(selfURL, primaryURL string) {
	signerAddr := ""
	if dag.signingKey != nil {
		signerAddr = strings.ToLower(crypto.PubkeyToAddress(dag.signingKey.PublicKey).Hex())
	}
	signature := fetchAndSignPeerChallenge(primaryURL, signerAddr, dag.signingKey)

	// Resolve the operator binding signature. Prefer the explicit env var
	// (set manually via /node-binding for cases where the operator wallet
	// is separate from the RELAYER key). Auto-sign when the two coincide:
	// if NODE_OPERATOR_WALLET matches the RELAYER address (same private
	// key), the node can produce the EIP-191 binding proof itself without
	// the operator doing anything out-of-band. This is the common
	// single-key deployment pattern.
	operatorBindingSig := os.Getenv("NODE_OPERATOR_BINDING_SIGNATURE")
	if operatorBindingSig == "" && dag.signingKey != nil && signerAddr != "" {
		nodeWallet := strings.ToLower(strings.TrimSpace(os.Getenv("NODE_OPERATOR_WALLET")))
		if nodeWallet == "" {
			nodeWallet = signerAddr
		}
		if nodeWallet == signerAddr {
			bindingMsg := "Aequitas: authorize validator " + signerAddr
			msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(bindingMsg), bindingMsg)
			hash := crypto.Keccak256Hash([]byte(msg))
			if sig, err := crypto.Sign(hash.Bytes(), dag.signingKey); err == nil {
				operatorBindingSig = "0x" + hex.EncodeToString(sig)
			}
		}
	}

	body, _ := json.Marshal(map[string]string{
		"url":                        selfURL,
		"signing_address":            signerAddr,
		"signature":                  signature,
		"peer_secret":                os.Getenv("PEER_SECRET"),
		"node_operator_wallet":       strings.ToLower(os.Getenv("NODE_OPERATOR_WALLET")),
		"operator_binding_signature": operatorBindingSig,
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
