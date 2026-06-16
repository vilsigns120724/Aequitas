package keeper

import (
"database/sql"
"encoding/hex"
"fmt"
"math/big"
"strings"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/core/rawdb"
"github.com/ethereum/go-ethereum/core/state"
)

// ─── CONTRACT STORAGE ─────────────────────────────────────────────────────────

func (cs *ChainState) SaveContract(address string, bytecode []byte, deployer string) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_contracts (address, bytecode, deployer) VALUES ($1, $2, $3)
 ON CONFLICT (address) DO UPDATE SET bytecode = $2`,
address, hex.EncodeToString(bytecode), deployer,
)
if err != nil {
fmt.Printf("[EVM] Error saving contract: %v\n", err)
}
return err
}

func (cs *ChainState) LoadContract(address string) ([]byte, error) {
if cs.db == nil {
return nil, nil
}
var bytecodeHex string
err := cs.db.QueryRow(
`SELECT bytecode FROM evm_contracts WHERE address = $1`, address,
).Scan(&bytecodeHex)
if err == sql.ErrNoRows {
return nil, nil
}
if err != nil {
return nil, err
}
return hex.DecodeString(bytecodeHex)
}

func (cs *ChainState) GetAllContracts() []string {
if cs.db == nil {
return nil
}
rows, err := cs.db.Query(`SELECT address FROM evm_contracts`)
if err != nil {
return nil
}
defer rows.Close()
var addrs []string
for rows.Next() {
var addr string
rows.Scan(&addr)
addrs = append(addrs, addr)
}
return addrs
}

// ─── NONCE STORAGE ────────────────────────────────────────────────────────────

func (cs *ChainState) SaveNonce(address string, nonce uint64) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_nonces (address, nonce) VALUES ($1, $2)
 ON CONFLICT (address) DO UPDATE SET nonce = $2`,
address, nonce,
)
return err
}

func (cs *ChainState) LoadNonce(address string) uint64 {
if cs.db == nil {
return 0
}
var nonce uint64
cs.db.QueryRow(`SELECT nonce FROM evm_nonces WHERE address = $1`, address).Scan(&nonce)
return nonce
}

// ─── CONTRACT STORAGE SLOTS ───────────────────────────────────────────────────

func (cs *ChainState) SaveStorageSlot(address, slot, value string) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_storage (address, slot, value) VALUES ($1, $2, $3)
 ON CONFLICT (address, slot) DO UPDATE SET value = $3`,
address, slot, value,
)
return err
}

func (cs *ChainState) LoadStorageSlot(address, slot string) (string, error) {
if cs.db == nil {
return "", nil
}
var value string
err := cs.db.QueryRow(
`SELECT value FROM evm_storage WHERE address = $1 AND slot = $2`,
address, slot,
).Scan(&value)
if err == sql.ErrNoRows {
return "", nil
}
return value, err
}

// ─── EVM ENGINE HELPERS ───────────────────────────────────────────────────────

// PersistContractStorage reads storage slots from a stateDB and saves to PostgreSQL.
// Since we no longer have a persistent stateDB, this is a no-op log.
func (e *EVMEngine) PersistContractStorage(contractAddr common.Address) {
fmt.Printf("[EVM] Contract %s active in session\n", strings.ToLower(contractAddr.Hex()))
}

// NewPersistentStateDB creates a StateDB loaded from PostgreSQL.
// Used by tests and legacy code. For production use EVMEngine.newStateDB().
func NewPersistentStateDB(cs *ChainState) (*state.StateDB, error) {
memDB := rawdb.NewMemoryDatabase()
sdb, err := state.New(common.Hash{}, state.NewDatabase(memDB), nil)
if err != nil {
return nil, err
}

for _, acc := range cs.GetAllAccounts() {
addr := common.HexToAddress(acc.Address)
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
wei := new(big.Int).Mul(big.NewInt(int64(acc.Balance)), decimals)
sdb.SetBalance(addr, wei)
sdb.SetNonce(addr, cs.LoadNonce(acc.Address))
}

for _, addrStr := range cs.GetAllContracts() {
addr := common.HexToAddress(addrStr)
code, err := cs.LoadContract(addrStr)
if err != nil || len(code) == 0 {
continue
}
sdb.SetCode(addr, code)
fmt.Printf("[EVM] Loaded contract: %s (%d bytes)\n", addrStr, len(code))

if cs.db != nil {
rows, err := cs.db.Query(
`SELECT slot, value FROM evm_storage WHERE address = $1`, addrStr)
if err == nil {
for rows.Next() {
var slot, value string
rows.Scan(&slot, &value)
sdb.SetState(addr, common.HexToHash(slot), common.HexToHash(value))
}
rows.Close()
}
}
}

sdb.Commit(0, false)
return sdb, nil
}
