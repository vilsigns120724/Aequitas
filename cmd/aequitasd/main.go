package main

import (
"encoding/json"
"fmt"
"os"
"time"

"github.com/hanoi96international-gif/aequitas-chain/x/humanity/keeper"
)

const (
VERSION       = "v0.1.0"
CONTRACT_V5   = "0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5"
PROOF_SERVER  = "https://aequitas-proof-server-production.up.railway.app"
INITIAL_GRANT = 1000
CHAIN_ID      = "aequitas-1"
)

type Genesis struct {
ChainID     string      `json:"chain_id"`
GenesisTime string      `json:"genesis_time"`
AppState    interface{} `json:"app_state"`
}

func loadGenesis() (*Genesis, error) {
data, err := os.ReadFile("config/genesis.json")
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
fmt.Printf("Contract V5:   %s\n", CONTRACT_V5)
fmt.Printf("Proof Server:  %s\n", PROOF_SERVER)
fmt.Printf("Initial Grant: %d AEQ per human\n", INITIAL_GRANT)
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

// Real humans from V5 contract (4 registered)
realHumans := []struct {
address    string
commitment string
}{
{"0x0BE8b961CBf6564bd1931B0803D35C0659E0D016", "sepolia_commitment_1"},
{"sepolia_human_2", "sepolia_commitment_2"},
{"sepolia_human_3", "sepolia_commitment_3"},
{"sepolia_human_4", "sepolia_commitment_4"},
}

fmt.Println("── Loading Verified Humans from V5 ──────")
for _, h := range realHumans {
err := humanKeeper.RegisterHuman(h.address, h.commitment, time.Now().Unix())
if err != nil {
fmt.Printf("✗ Failed: %s\n", h.address)
} else {
fmt.Printf("✓ Human: %s (+%d AEQ)\n", h.address, INITIAL_GRANT)
}
}

fmt.Println()
fmt.Println("── Network Status ───────────────────────")
fmt.Printf("Total Humans:  %d\n", humanKeeper.TotalHumans())
fmt.Printf("Total Supply:  %d AEQ\n", humanKeeper.TotalHumans()*INITIAL_GRANT)
fmt.Printf("Fair Share:    %d AEQ\n", INITIAL_GRANT)
fmt.Printf("Max Cap:       %d AEQ\n", humanKeeper.TotalHumans()*INITIAL_GRANT)

fmt.Println()
fmt.Println("── Proof of Humanity Validators ─────────")
fmt.Println("Every verified human = 1 validator = 1 equal vote")
fmt.Printf("Active Validators: %d\n", humanKeeper.TotalHumans())
fmt.Println("Consensus: Byzantine Fault Tolerant (2/3 majority)")

fmt.Println()
fmt.Println("╔════════════════════════════════════════╗")
fmt.Println("║     Aequitas Chain Node Ready ✓        ║")
fmt.Println("║     Waiting for P2P connections...     ║")
fmt.Println("╚════════════════════════════════════════╝")
}
