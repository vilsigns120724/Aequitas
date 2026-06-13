package keeper

import (
"encoding/json"
"fmt"

"github.com/hanoi96international-gif/aequitas-chain/x/humanity/types"
)

type Keeper struct {
humans map[string]*types.Human
commitments map[string]bool
}

func NewKeeper() *Keeper {
return &Keeper{
humans:      make(map[string]*types.Human),
commitments: make(map[string]bool),
}
}

func (k *Keeper) RegisterHuman(address, commitment string, timestamp int64) error {
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
h, exists := k.humans[address]
return exists && h.IsActive
}

func (k *Keeper) TotalHumans() int {
count := 0
for _, h := range k.humans {
if h.IsActive {
count++
}
}
return count
}

func (k *Keeper) GetAllHumans() []*types.Human {
result := make([]*types.Human, 0, len(k.humans))
for _, h := range k.humans {
result = append(result, h)
}
return result
}

func (k *Keeper) ExportState() ([]byte, error) {
return json.Marshal(k.humans)
}
