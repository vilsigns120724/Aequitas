package main

import (
"encoding/json"
"fmt"
"os"
"os/signal"
"syscall"
"time"

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
	if bootstrapURL := os.Getenv("BOOTSTRAP_SNAPSHOT_URL"); bootstrapURL != "" && chainState.TotalHumans() == 0 {
		expectedSigner := os.Getenv("BOOTSTRAP_SIGNER")
		// BOOTSTRAP_SIGNER is always required — it verifies the cryptographic
		// integrity of the snapshot payload regardless of HTTP token auth.
		// A token-only import would accept any syntactically valid snapshot.
		if expectedSigner == "" {
			fmt.Println("[BOOTSTRAP] ✗ Refused: BOOTSTRAP_SIGNER must be set to the primary node's signing address (0x...)")
		} else {
			fmt.Printf("[BOOTSTRAP] Fresh node — importing state from %s\n", bootstrapURL)
			if err := chainState.ImportSnapshotFromURL(bootstrapURL, expectedSigner); err != nil {
				fmt.Printf("[BOOTSTRAP] ✗ Import failed: %v\n", err)
			}
		}
	}

	// HTTP Block Sync between nodes.
	// SELF_URL must be set to this node's own public URL so the sync loop
	// can exclude it from the peer list — without this, a node would try to
	// sync from itself and generate spurious hash-mismatch rejections.
	selfURL := os.Getenv("SELF_URL")
	if selfURL == "" {
		selfURL = "https://aequitas-production-9fba.up.railway.app" // legacy default for Railway
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
if os.Getenv("IS_PRIMARY_NODE") == "true" {
fmt.Println("[POOLS] Primary node — daily distributions enabled (UBI + Validators + LP)")
go func() {
// Calculate how long until the NEXT distribution, accounting for
// any time that has already elapsed since the last one. This ensures
// a restart does not reset the 24h clock and delay payouts.
lastAt := chainState.GetLastUBIAt()
var firstDelay time.Duration
if lastAt == 0 {
firstDelay = 24 * time.Hour // never distributed yet
} else {
nextRun := time.Unix(lastAt, 0).Add(24 * time.Hour)
firstDelay = time.Until(nextRun)
if firstDelay < 0 {
firstDelay = 0 // overdue — run immediately
}
}
fmt.Printf("[POOLS] Next distribution in %s\n", firstDelay.Round(time.Minute))
time.Sleep(firstDelay)
chainState.DistributeUBIPool()
chainState.DistributeValidatorsPool()
chainState.DistributeLPPool()
// Then every 24 h precisely
ticker := time.NewTicker(24 * time.Hour)
for range ticker.C {
chainState.DistributeUBIPool()
chainState.DistributeValidatorsPool()
chainState.DistributeLPPool()
}
}()
} else {
fmt.Println("[POOLS] Not primary node — skipping daily distributions")
}

// Register this node's operator wallet so it participates in validator
// pool distributions. Any node that sets NODE_OPERATOR_WALLET gets
// included automatically — no code change needed when new nodes join.
if wallet := os.Getenv("NODE_OPERATOR_WALLET"); wallet != "" {
if warn := chainState.ValidateNodeOperatorWallet(wallet); warn != "" {
fmt.Printf("[NODE] Warning: %s\n", warn)
}
chainState.RegisterNode(wallet)
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
fmt.Println("\nNode stopped.")
}
