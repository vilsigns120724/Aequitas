package keeper

import (
"encoding/json"
"fmt"
"net/http"
"time"
)

const (
ProofServer    = "https://aequitas-proof-server-production.up.railway.app"
SyncInterval   = 60 * time.Second
)

type ChainStatus struct {
Registrations int     `json:"registrations"`
Supply        string  `json:"supply"`
Gini          int     `json:"gini"`
Index         int     `json:"index"`
Phase         int     `json:"phase"`
}

func (k *Keeper) StartSync() {
fmt.Println("── Starting Sepolia Sync ────────────────")
go func() {
ticker := time.NewTicker(SyncInterval)
for range ticker.C {
k.syncFromSepolia()
}
}()
}

func (k *Keeper) syncFromSepolia() {
resp, err := http.Get(ProofServer + "/health")
if err != nil {
fmt.Printf("[SYNC] Error: %v\n", err)
return
}
defer resp.Body.Close()

var status ChainStatus
if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
fmt.Printf("[SYNC] Parse error: %v\n", err)
return
}

currentHumans := k.TotalHumans()
if status.Registrations > currentHumans {
fmt.Printf("[SYNC] Sepolia: %d humans | Layer1: %d humans — syncing...\n",
status.Registrations, currentHumans)
k.syncHumans(status.Registrations)
}
}

func (k *Keeper) syncHumans(target int) {
current := k.TotalHumans()
for i := current; i < target; i++ {
addr := fmt.Sprintf("sepolia_human_%d", i+1)
commitment := fmt.Sprintf("sepolia_commitment_%d", i+1)
err := k.RegisterHuman(addr, commitment, time.Now().Unix())
if err == nil {
fmt.Printf("[SYNC] ✓ Synced human #%d\n", i+1)
}
}
}
