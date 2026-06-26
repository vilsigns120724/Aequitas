package keeper

import (
"encoding/json"
"fmt"
"sync"

"github.com/hanoi96international-gif/aequitas-chain/x/humanity/types"
)

type Keeper struct {
mu          sync.RWMutex
humans      map[string]*types.Human
commitments map[string]bool
}

func NewKeeper() *Keeper {
return &Keeper{
humans:      make(map[string]*types.Human),
commitments: make(map[string]bool),
}
}

func (k *Keeper) RegisterHuman(address, commitment string, timestamp int64) error {
k.mu.Lock()
defer k.mu.Unlock()
if _, exists := k.humans[address]; exists {
return fmt.Errorf("address already registered")
}
if k.commitments[commitment] {
return fmt.Errorf("commitment already used")
}
k.humans[address] = &types.Human{
Address:      address,
Commitment:   commitment,
RegisteredAt: timestamp,
IsActive:     true,
}
k.commitments[commitment] = true
return nil
}

func (k *Keeper) IsHuman(address string) bool {
k.mu.RLock()
defer k.mu.RUnlock()
h, exists := k.humans[address]
return exists && h.IsActive
}

func (k *Keeper) TotalHumans() int {
k.mu.RLock()
defer k.mu.RUnlock()
count := 0
for _, h := range k.humans {
if h.IsActive {
count++
}
}
return count
}

// GetAllHumans returns copies of every registered Human, not pointers into
// the live map — a caller mutating a returned *Human used to mutate the
// keeper's own state unsynchronized with k.mu, defeating the lock's purpose.
func (k *Keeper) GetAllHumans() []*types.Human {
k.mu.RLock()
defer k.mu.RUnlock()
result := make([]*types.Human, 0, len(k.humans))
for _, h := range k.humans {
hCopy := *h
result = append(result, &hCopy)
}
return result
}

func (k *Keeper) ExportState() ([]byte, error) {
k.mu.RLock()
defer k.mu.RUnlock()
return json.Marshal(k.humans)
}
