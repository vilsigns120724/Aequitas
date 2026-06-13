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
VERSION       = "v0.2.0"
CONTRACT_V5   = "0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5"
PROOF_SERVER  = "https://aequitas-proof-server-production.up.railway.app"
INITIAL_GRANT = 1000
CHAIN_ID      = "aequitas-1"
BLOCK_TIME    = 6 * time.Second
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

// Load Genesis
fmt.Println("── Loading Genesis Block ────────────────")
genesis, err := loadGenesis()
if err != nil {
fmt.Printf("✗ Genesis error: %v\n", err)
} else {
fmt.Printf("✓ Chain ID: %s\n", genesis.ChainID)
fmt.Printf("✓ Genesis Time: %s\n", genesis.GenesisTime)
}
fmt.Println()

// Initialize Humanity Keeper
humanKeeper := keeper.NewKeeper()

realHumans := []struct {
address    string
commitment string
}{
{"0x0BE8b961CBf6564bd1931B0803D35C0659E0D016", "sepolia_commitment_1"},
{"sepolia_human_2", "sepolia_commitment_2"},
{"sepolia_human_3", "sepolia_commitment_3"},
{"sepolia_human_4", "sepolia_commitment_4"},
}

fmt.Println("── Loading Verified Humans ──────────────")
for _, h := range realHumans {
err := humanKeeper.RegisterHuman(h.address, h.commitment, time.Now().Unix())
if err == nil {
fmt.Printf("✓ Human: %s (+%d AEQ)\n", h.address, INITIAL_GRANT)
}
}
fmt.Println()
fmt.Printf("Total Humans:  %d\n", humanKeeper.TotalHumans())
fmt.Printf("Total Supply:  %d AEQ\n", humanKeeper.TotalHumans()*INITIAL_GRANT)
fmt.Println()

// Initialize Blockchain
fmt.Println("── Initializing Blockchain ──────────────")
p2pNode, err := keeper.NewP2PNode(humanKeeper)
if err != nil {
fmt.Printf("✗ P2P Error: %v\n", err)
return
}

bc := keeper.NewBlockchain(humanKeeper, p2pNode.GetNodeID())
fmt.Println()

// Start P2P
p2pNode.Start()

multiaddr := p2pNode.GetMultiaddr()
fmt.Println("── Share this address to join network ───")
fmt.Printf("%s\n", multiaddr)
fmt.Println()

// Start block production
fmt.Println("── Starting Block Production ────────────")
go func() {
ticker := time.NewTicker(BLOCK_TIME)
for range ticker.C {
block := bc.ProduceBlock()
fmt.Printf("[Block #%d] Hash: %s... | Humans: %d | Time: %s\n",
block.Height,
block.Hash[:16],
block.Humans,
time.Unix(block.Timestamp, 0).Format("15:04:05"),
)
}
}()

fmt.Println("╔════════════════════════════════════════╗")
fmt.Println("║     Aequitas Node Running ✓            ║")
fmt.Println("║     Producing blocks every 6 seconds   ║")
fmt.Println("║     Press Ctrl+C to stop               ║")
fmt.Println("╚════════════════════════════════════════╝")

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
fmt.Println("\nNode stopped.")
}
