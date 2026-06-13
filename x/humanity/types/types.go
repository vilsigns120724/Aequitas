package types

const (
    ModuleName = "humanity"
    StoreKey   = ModuleName
    RouterKey  = ModuleName
)

// Human represents a verified human validator
type Human struct {
    Address   string `json:"address"`
    Commitment string `json:"commitment"`
    RegisteredAt int64 `json:"registered_at"`
    IsActive  bool   `json:"is_active"`
}
