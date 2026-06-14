package keeper

import (
"crypto/sha256"
"database/sql"
"encoding/hex"
"encoding/json"
"fmt"
"os"
"sync"
	"strings"

_ "github.com/lib/pq"
)

type AccountState struct {
Address string  `json:"address"`
Balance float64 `json:"balance"`
IsHuman bool    `json:"is_human"`
}

type ChainState struct {
mu       sync.RWMutex
accounts map[string]*AccountState
db       *sql.DB
useDB    bool
}

func NewChainState(dataFile string) *ChainState {
cs := &ChainState{
accounts: make(map[string]*AccountState),
}

// Try PostgreSQL first
dbURL := os.Getenv("DATABASE_URL")
if dbURL != "" {
// Add sslmode if not present
if !strings.Contains(dbURL, "sslmode") {
if strings.Contains(dbURL, "?") {
dbURL += "&sslmode=disable"
} else {
dbURL += "?sslmode=disable"
}
}
db, err := sql.Open("postgres", dbURL)
if err == nil {
err = db.Ping()
if err == nil {
cs.db = db
cs.useDB = true
cs.initDB()
cs.loadFromDB()
fmt.Println("✓ ChainState using PostgreSQL")
return cs
}
}
fmt.Printf("⚠ PostgreSQL failed: %v - using file\n", err)
}

// Fallback to file
cs.useDB = false
if os.Getenv("RESET_STATE") == "true" {
fmt.Println("✓ RESET_STATE=true — starting fresh")
os.Remove(dataFile)
} else {
cs.loadFromFile(dataFile)
}
return cs
}

func (cs *ChainState) initDB() {
cs.db.Exec(`CREATE TABLE IF NOT EXISTS chain_accounts (
address TEXT PRIMARY KEY,
balance FLOAT NOT NULL DEFAULT 0,
is_human BOOLEAN NOT NULL DEFAULT false
)`)
}

func (cs *ChainState) loadFromDB() {
rows, err := cs.db.Query("SELECT address, balance, is_human FROM chain_accounts")
if err != nil {
fmt.Printf("⚠ Could not load from DB: %v\n", err)
return
}
defer rows.Close()
count := 0
for rows.Next() {
acc := &AccountState{}
rows.Scan(&acc.Address, &acc.Balance, &acc.IsHuman)
cs.accounts[acc.Address] = acc
count++
}
fmt.Printf("✓ Loaded %d accounts from PostgreSQL\n", count)
}

func (cs *ChainState) loadFromFile(dataFile string) {
data, err := os.ReadFile(dataFile)
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
if cs.useDB {
return // DB saves immediately in RegisterHuman/Transfer
}
data, _ := json.Marshal(cs.accounts)
os.WriteFile("/tmp/aequitas_state.json", data, 0644)
}

func (cs *ChainState) saveAccountToDB(acc *AccountState) {
if !cs.useDB {
return
}
cs.db.Exec(`INSERT INTO chain_accounts (address, balance, is_human) VALUES ($1, $2, $3)
ON CONFLICT (address) DO UPDATE SET balance = $2, is_human = $3`,
acc.Address, acc.Balance, acc.IsHuman)
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
cs.saveAccountToDB(cs.accounts[address])
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
cs.saveAccountToDB(fromAcc)

if _, ok := cs.accounts[to]; !ok {
cs.accounts[to] = &AccountState{Address: to}
}
cs.accounts[to].Balance += amount
cs.saveAccountToDB(cs.accounts[to])
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
