package keeper

import (
"fmt"
"sync"
)

const InitialGrant = 1000 // AEQ per verified human

type BalanceStore struct {
mu       sync.RWMutex
balances map[string]float64 // address -> AEQ balance
}

func NewBalanceStore() *BalanceStore {
return &BalanceStore{
balances: make(map[string]float64),
}
}

func (b *BalanceStore) GetBalance(address string) float64 {
b.mu.RLock()
defer b.mu.RUnlock()
return b.balances[address]
}

func (b *BalanceStore) Grant(address string, amount float64) {
b.mu.Lock()
defer b.mu.Unlock()
b.balances[address] += amount
fmt.Printf("[BALANCE] ✓ Granted %.2f AEQ to %s\n", amount, address)
}

func (b *BalanceStore) Transfer(from, to string, amount float64) error {
b.mu.Lock()
defer b.mu.Unlock()
if b.balances[from] < amount {
return fmt.Errorf("insufficient balance: have %.2f, need %.2f", b.balances[from], amount)
}
b.balances[from] -= amount
b.balances[to] += amount
fmt.Printf("[BALANCE] ✓ Transfer %.2f AEQ from %s to %s\n", amount, from, to)
return nil
}

func (b *BalanceStore) GetAll() map[string]float64 {
b.mu.RLock()
defer b.mu.RUnlock()
result := make(map[string]float64)
for k, v := range b.balances {
result[k] = v
}
return result
}

func (b *BalanceStore) TotalSupply() float64 {
b.mu.RLock()
defer b.mu.RUnlock()
total := 0.0
for _, v := range b.balances {
total += v
}
return total
}
