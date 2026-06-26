package main

import (
"context"
"encoding/json"
"fmt"
"net"
"net/url"
"os"
"os/signal"
"syscall"
"time"
_ "time/tzdata" // embed IANA timezone DB so Europe/Berlin works on Alpine without system tzdata

"github.com/hanoi96international-gif/aequitas-chain/x/humanity/keeper"
)

const (
VERSION       = "v0.3.0"
// NOTE: the actually-active contract addresses (V6, V7, bio verifier) live
// in x/humanity/keeper/evm_v6mirror.go (V6_CONTRACT_ADDR, V7_CONTRACT_ADDR,
// BIO_VERIFIER_ADDR) — that is the single source of truth. Do not redeclare
// addresses here; a previous version of this file had a stale CONTRACT_V6
// and BIO_VERIFIER value that didn't match what was actually deployed and
// was never even referenced anywhere in this file.
PROOF_SERVER  = "https://aequitas-proof-server-production.up.railway.app"
INITIAL_GRANT = 1000
CHAIN_ID      = "aequitas-1"
BLOCK_TIME    = 6 * time.Second
API_PORT      = 8080
)

type Genesis struct {
ChainID     string      `json:"chain_id"`
GenesisTime string      `json:"genesis_time"`
AppState    interface{} `json:"app_state"`
}

func loadGenesis() (*Genesis, error) {
data, err := os.ReadFile("genesis.json")
if err != nil {
return nil, err
}
var genesis Genesis
err = json.Unmarshal(data, &genesis)
return &genesis, err
}

func main() {
fmt.Println("╔════════════════════════════════════════╗")
fmt.Println("║         AEQUITAS CHAIN NODE            ║")
fmt.Println("║      Proof of Humanity Consensus       ║")
fmt.Println("╚════════════════════════════════════════╝")
fmt.Println()
fmt.Printf("Version:       %s\n", VERSION)
fmt.Printf("Chain ID:      %s\n", CHAIN_ID)
fmt.Printf("Block Time:    %s\n", BLOCK_TIME)
fmt.Println()

fmt.Println("── Loading Genesis Block ────────────────")
genesis, err := loadGenesis()
if err != nil {
fmt.Printf("✗ Genesis error: %v\n", err)
} else {
fmt.Printf("✓ Chain ID: %s\n", genesis.ChainID)
fmt.Printf("✓ Genesis Time: %s\n", genesis.GenesisTime)
}
fmt.Println()

humanKeeper := keeper.NewKeeper()

fmt.Println()

fmt.Println("── Initializing Blockchain ──────────────")
p2pNode, err := keeper.NewP2PNode(humanKeeper)
if err != nil {
fmt.Printf("✗ P2P Error: %v\n", err)
return
}

chainState := keeper.NewChainState("/tmp/aequitas_state.json")
bc := keeper.NewBlockchain(humanKeeper, p2pNode.GetNodeID(), chainState)
// Load individually-registered validator keys from DB into the DAG's
// authorized set so they survive node restarts without re-registration.
chainState.LoadValidatorKeysIntoDAG(bc)
fmt.Println()

p2pNode.SetDAG(bc)

	// Bootstrap from a peer snapshot if this is a fresh node (no humans in DB).
	// Set BOOTSTRAP_SNAPSHOT_URL to the primary node's /api/snapshot endpoint.
	// Set SNAPSHOT_TOKEN to match the primary node's SNAPSHOT_TOKEN env var.
	// Set BOOTSTRAP_SIGNER to the primary node's signing address (0x...) to
	// verify the snapshot's ECDSA signature before importing.
	// After startup, ongoing state sync happens via block TX replay in AddPeerBlock.
	if bootstrapURL := os.Getenv("BOOTSTRAP_SNAPSHOT_URL"); bootstrapURL != "" && chainState.TotalHumans() == 0 {
		// FIX 15: Validate URL scheme and host before fetching to prevent SSRF.
		parsedBootstrap, urlErr := url.Parse(bootstrapURL)
		if urlErr != nil || (parsedBootstrap.Scheme != "https" && parsedBootstrap.Scheme != "http") {
			fmt.Printf("[BOOTSTRAP] ✗ Refused: BOOTSTRAP_SNAPSHOT_URL must be an http or https URL (got %q)\n", bootstrapURL)
		} else if host := parsedBootstrap.Hostname(); isRFC1918OrLoopback(host) {
			fmt.Printf("[BOOTSTRAP] ✗ Refused: BOOTSTRAP_SNAPSHOT_URL must not point to a private/loopback address (got %q)\n", host)
		} else {
			expectedSigner := os.Getenv("BOOTSTRAP_SIGNER")
			if expectedSigner == "" {
				fmt.Println("[BOOTSTRAP] ✗ Refused: BOOTSTRAP_SIGNER must be set to the primary node's signing address (0x...)")
			} else {
				fmt.Printf("[BOOTSTRAP] Fresh node — importing state from %s\n", bootstrapURL)
				if err := chainState.ImportSnapshotFromURL(bootstrapURL, expectedSigner); err != nil {
					fmt.Printf("[BOOTSTRAP] ✗ Import failed: %v\n", err)
				}
			}
		}
	}

	// Save price snapshots every 30 seconds so the chart interval buttons
	// (1m/5m/30m/1h/4h) show meaningful historical data even without swaps.
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			chainState.SavePriceSnapshot()
		}
	}()

	// HTTP Block Sync between nodes.
	// SELF_URL must be set to this node's own public URL so the sync loop
	// can exclude it from the peer list — without this, a node would try to
	// sync from itself and generate spurious hash-mismatch rejections.
	selfURL := os.Getenv("SELF_URL")
	if selfURL == "" {
		selfURL = "https://aequitas.digital" // default; override with SELF_URL env var
	}
	bc.StartHTTPBlockSync(selfURL)
	p2pNode.Start()
	bc.ReconstructState(chainState)

// Humans register natively via the V7 contract (see register.go) and are
// reconstructed from blockchain transactions above via ReconstructState.
// An old Sepolia-polling sync (humanKeeper.StartSync(), in the now-removed
// sync.go) used to inject placeholder "sepolia_human_N" entries into the
// keeper on every tick — exactly the fake registrations that were
// deliberately removed earlier in this project. That code path is gone
// now, not just disabled, so it can't be accidentally re-enabled.

multiaddr := p2pNode.GetMultiaddr()
fmt.Println("── Share this address to join network ───")
fmt.Printf("%s\n", multiaddr)
fmt.Println()

fmt.Println("── Starting API Server ──────────────────")
api := keeper.NewAPIServer(bc, p2pNode, humanKeeper, chainState)
go api.Start(API_PORT)

fmt.Println()

fmt.Println("── Starting Block Production ────────────")
go func() {
ticker := time.NewTicker(BLOCK_TIME)
for range ticker.C {
block := bc.ProduceBlock()
			p2pNode.BroadcastBlock(block)
fmt.Printf("[Block #%d] Hash: %s... | Humans: %d | Time: %s\n",
block.Height,
block.Hash[:16],
block.Humans,
time.Unix(block.Timestamp, 0).Format("15:04:05"),
)
}
}()

fmt.Println("── Starting Daily Pool Distributions ───")
// IS_PRIMARY_NODE is no longer required. Every node schedules distributions
// at 20:00 Berlin but uses a PostgreSQL CAS lock (TryLockDistribution) to
// ensure only ONE node actually executes — the first to atomically claim
// the lock wins; others see 0 rows updated and skip. This eliminates the
// IS_PRIMARY_NODE env-var footgun where any operator could accidentally (or
// maliciously) trigger double distributions by setting IS_PRIMARY_NODE=true.
// FIX 11: Create a cancellable context for the distribution goroutine so
// it can be cleanly stopped on node shutdown.
distCtx, distCancel := context.WithCancel(context.Background())
_ = distCancel // called in shutdown handler below
// FIX 6: distDone is closed when the goroutine exits so the shutdown handler
// can wait for any in-progress distribution to finish before the process exits.
distDone := make(chan struct{})
go func(ctx context.Context) {
	defer close(distDone)
berlin, err := time.LoadLocation("Europe/Berlin")
if err != nil {
berlin = time.FixedZone("CET", 2*60*60) // CEST fallback (summer, UTC+2)
}

nextDaily20 := func(after time.Time) time.Time {
t := after.In(berlin)
candidate := time.Date(t.Year(), t.Month(), t.Day(), 20, 0, 0, 0, berlin)
if !after.Before(candidate) {
// Use AddDate to get exactly 20:00 next day regardless of DST transitions.
// Add(24h) would be 1h off on the two DST changeover nights per year.
candidate = time.Date(t.Year(), t.Month(), t.Day()+1, 20, 0, 0, 0, berlin)
}
return candidate
}

lastAt := chainState.GetLastUBIAt()
var firstTarget time.Time
// FIX 10: Apply the same 1-hour guard for fresh DBs (lastAt == 0) to prevent
// the distribution from firing immediately on a brand-new node startup.
if lastAt == 0 || time.Since(time.Unix(lastAt, 0)) < time.Hour {
firstTarget = nextDaily20(time.Now().Add(time.Hour))
} else {
firstTarget = nextDaily20(time.Now())
}

firstDelay := time.Until(firstTarget)
chainState.SetNextUBIAt(firstTarget.Unix())
fmt.Printf("[POOLS] Next distribution at %s Berlin time (in %s)\n",
firstTarget.In(berlin).Format("02.01. 15:04:05"), firstDelay.Round(time.Minute))

for {
// FIX 11: Use a select so the goroutine unblocks on shutdown.
select {
case <-ctx.Done():
return
case <-time.After(time.Until(firstTarget)):
}
// Distributed lock: only the node that atomically claims last_ubi_at proceeds.
if chainState.TryLockDistribution() {
// Read UBI pool balance and human count BEFORE distribution so the
// TX carries the exact per-human share for secondary node replay.
preUBIBalance := chainState.GetBalance("0x4a9b8f99f0d8cff0e510fef502100571203b054a")
totalHumans := chainState.TotalHumans()
chainState.DistributeUBIPool()
chainState.DistributeValidatorsPool()
chainState.DistributeLPPool()
// Move inactive wallets to escrow, then release long-sitting escrow
// balances to the UBI pool for the next distribution cycle.
chainState.CheckAndMoveToEscrow()
chainState.ReleaseEscrowToUBI()
// Emit a ubi_distribution TX so secondary nodes can replay the exact
// per-human UBI credit via ApplyUBIDelta instead of running the
// distribution logic independently (which diverges if state differs).
var amountPerHuman float64
if totalHumans > 0 { amountPerHuman = preUBIBalance / float64(totalHumans) }
bc.AddTransaction(keeper.Transaction{
	Type:           "ubi_distribution",
	Wallet:         "0x0000000000000000000000000000000000000000",
	Amount:         preUBIBalance,
	AmountPerHuman: amountPerHuman,
})
fmt.Printf("[POOLS] ✓ Distribution done at %s\n", time.Now().In(berlin).Format("02.01. 15:04:05"))
} else {
fmt.Printf("[POOLS] Another node ran distribution first — skipping\n")
}
firstTarget = nextDaily20(time.Now())
chainState.SetNextUBIAt(firstTarget.Unix())
fmt.Printf("[POOLS] Next distribution at %s Berlin time\n", firstTarget.In(berlin).Format("02.01. 15:04:05"))
}
}(distCtx)

// Register this node's operator wallet so it participates in validator
// pool distributions. Any node that sets NODE_OPERATOR_WALLET gets
// included automatically — no code change needed when new nodes join.
if wallet := os.Getenv("NODE_OPERATOR_WALLET"); wallet != "" {
if warn := chainState.ValidateNodeOperatorWallet(wallet); warn != "" {
fmt.Printf("[NODE] ERROR: %s\n", warn)
fmt.Println("[NODE] NODE_OPERATOR_WALLET rejected — complete biometric registration first")
// Do NOT register — only verified humans receive validator rewards
} else {
chainState.RegisterNode(wallet)
}
} else {
fmt.Println("[NODE] NODE_OPERATOR_WALLET not set — this node won't receive validator rewards")
}

fmt.Println("╔════════════════════════════════════════╗")
fmt.Println("║     Aequitas Node Running ✓            ║")
fmt.Println("║     Producing blocks every 6 seconds   ║")
fmt.Println("╚════════════════════════════════════════╝")

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
// FIX 11: Signal the distribution goroutine to stop cleanly.
distCancel()
// FIX 6: Wait for the distribution goroutine to finish any in-progress work
// before exiting. A 10-second timeout prevents hanging indefinitely.
select {
case <-distDone:
case <-time.After(10 * time.Second):
	fmt.Println("[WARN] Distribution goroutine did not stop in 10 seconds — forcing exit")
}
fmt.Println("Node stopped.")
}

// isRFC1918OrLoopback returns true if host is a loopback address or an
// RFC 1918 private-network address. Used to block SSRF via BOOTSTRAP_SNAPSHOT_URL.
func isRFC1918OrLoopback(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		// Not a bare IP - hostname; resolve it for the check.
		addrs, err := net.LookupHost(host)
		if err != nil || len(addrs) == 0 {
			return false // cannot resolve; let the HTTP client fail naturally
		}
		ip = net.ParseIP(addrs[0])
		if ip == nil {
			return false
		}
	}
	// FIX 9: Normalise IPv6-mapped IPv4 addresses (e.g. ::ffff:10.0.0.1) to their
	// IPv4 form so the RFC-1918 CIDR checks below match them correctly.
	// net.IP.To4() returns nil for a pure IPv6 address and a 4-byte slice for
	// both native IPv4 and IPv6-mapped IPv4 — exactly what we want here.
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	if ip.IsLoopback() {
		return true
	}
	private := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16", // link-local
	}
	for _, cidr := range private {
		_, network, _ := net.ParseCIDR(cidr)
		if network != nil && network.Contains(ip) {
			return true
		}
	}
	return false
}
