package keeper

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"os"
"sync"
)

type AccountState struct {
Address  string  `json:"address"`
Balance  float64 `json:"balance"`
IsHuman  bool    `json:"is_human"`
}

type ChainState struct {
mu       sync.RWMutex
accounts map[string]*AccountState
dataFile string
}

func NewChainState(dataFile string) *ChainState {
cs := &ChainState{
accounts: make(map[string]*AccountState),
dataFile: dataFile,
}
cs.load()
return cs
}

func (cs *ChainState) load() {
data, err := os.ReadFile(cs.dataFile)
if err != nil {
fmt.Println("✓ Starting with fresh chain state")
return
}
var accounts map[string]*AccountState
if err := json.Unmarshal(data, &accounts); err != nil {
fmt.Println("⚠ Could not load state, starting fresh")
return
}
cs.accounts = accounts
fmt.Printf("✓ Loaded chain state: %d accounts\n", len(accounts))
}

func (cs *ChainState) save() {
cs.mu.RLock()
data, _ := json.Marshal(cs.accounts)
cs.mu.RUnlock()
os.WriteFile(cs.dataFile, data, 0644)
}

func (cs *ChainState) GetBalance(address string) float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
if acc, ok := cs.accounts[address]; ok {
return acc.Balance
}
return 0
}

func (cs *ChainState) IsHuman(address string) bool {
cs.mu.RLock()
defer cs.mu.RUnlock()
if acc, ok := cs.accounts[address]; ok {
return acc.IsHuman
}
return false
}

func (cs *ChainState) RegisterHuman(address string) error {
cs.mu.Lock()
defer cs.mu.Unlock()

if acc, ok := cs.accounts[address]; ok && acc.IsHuman {
return fmt.Errorf("already registered")
}

if _, ok := cs.accounts[address]; !ok {
cs.accounts[address] = &AccountState{Address: address}
}

cs.accounts[address].IsHuman = true
cs.accounts[address].Balance += 1000
cs.save()

fmt.Printf("[STATE] ✓ Human registered: %s | Balance: %.2f AEQ\n",
address, cs.accounts[address].Balance)
return nil
}

func (cs *ChainState) Transfer(from, to string, amount float64) error {
cs.mu.Lock()
defer cs.mu.Unlock()

fromAcc, ok := cs.accounts[from]
if !ok || fromAcc.Balance < amount {
return fmt.Errorf("insufficient balance")
}

fromAcc.Balance -= amount

if _, ok := cs.accounts[to]; !ok {
cs.accounts[to] = &AccountState{Address: to}
}
cs.accounts[to].Balance += amount
cs.save()

fmt.Printf("[STATE] ✓ Transfer %.2f AEQ: %s → %s\n", amount, from, to)
return nil
}

func (cs *ChainState) TotalSupply() float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
total := 0.0
for _, acc := range cs.accounts {
total += acc.Balance
}
return total
}

func (cs *ChainState) TotalHumans() int {
cs.mu.RLock()
defer cs.mu.RUnlock()
count := 0
for _, acc := range cs.accounts {
if acc.IsHuman {
count++
}
}
return count
}

func (cs *ChainState) GetAllAccounts() []*AccountState {
cs.mu.RLock()
defer cs.mu.RUnlock()
result := make([]*AccountState, 0, len(cs.accounts))
for _, acc := range cs.accounts {
result = append(result, acc)
}
return result
}

func (cs *ChainState) StateRoot() string {
cs.mu.RLock()
data, _ := json.Marshal(cs.accounts)
cs.mu.RUnlock()
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:])
}
