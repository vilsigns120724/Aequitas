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
fmt.Println()

p2pNode.SetDAG(bc)

	// HTTP Block Sync between nodes
	bc.StartHTTPBlockSync("https://aequitas-production-9fba.up.railway.app")
	p2pNode.Start()
	// Reconstruct state from blockchain
	time.Sleep(10 * time.Second)
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

fmt.Println("── Starting Daily UBI Distribution ──────")
// Runs independently of block production — UBI payout is a once-a-day
// event, not something tied to the ~6-second block cadence. Every node
// running this same loop would double/triple-pay the pool if more than
// one node's distribution call landed in the same day; for the current
// single-node deployment this is fine, but multi-node deployments would
// need this gated to a single designated node (not yet implemented).
go func() {
ticker := time.NewTicker(24 * time.Hour)
for range ticker.C {
chainState.DistributeUBIPool()
}
}()

fmt.Println("╔════════════════════════════════════════╗")
fmt.Println("║     Aequitas Node Running ✓            ║")
fmt.Println("║     Producing blocks every 6 seconds   ║")
fmt.Println("╚════════════════════════════════════════╝")

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
fmt.Println("\nNode stopped.")
}
