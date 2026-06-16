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

func (cs *ChainState) SaveContract(address string, bytecode []byte, deployer string) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_contracts (address, bytecode, deployer) VALUES ($1, $2, $3) ON CONFLICT (address) DO UPDATE SET bytecode = $2`,
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
err := cs.db.QueryRow(`SELECT bytecode FROM evm_contracts WHERE address = $1`, address).Scan(&bytecodeHex)
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

func (cs *ChainState) SaveNonce(address string, nonce uint64) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_nonces (address, nonce) VALUES ($1, $2) ON CONFLICT (address) DO UPDATE SET nonce = $2`,
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

func (cs *ChainState) SaveStorageSlot(address, slot, value string) error {
if cs.db == nil {
return nil
}
_, err := cs.db.Exec(
`INSERT INTO evm_storage (address, slot, value) VALUES ($1, $2, $3) ON CONFLICT (address, slot) DO UPDATE SET value = $3`,
address, slot, value,
)
return err
}

func (cs *ChainState) LoadStorageSlot(address, slot string) (string, error) {
if cs.db == nil {
return "", nil
}
var value string
err := cs.db.QueryRow(`SELECT value FROM evm_storage WHERE address = $1 AND slot = $2`, address, slot).Scan(&value)
if err == sql.ErrNoRows {
return "", nil
}
return value, err
}

func (e *EVMEngine) PersistContractStorage(contractAddr common.Address) {
fmt.Printf("[EVM] Contract %s active in session\n", strings.ToLower(contractAddr.Hex()))
}

func (e *EVMEngine) LoadContractStorage(contractAddr common.Address) {
addrStr := strings.ToLower(contractAddr.Hex())
rows, err := e.chainState.db.Query(`SELECT slot, value FROM evm_storage WHERE address = $1`, addrStr)
if err != nil {
return
}
defer rows.Close()
count := 0
for rows.Next() {
var slot, value string
rows.Scan(&slot, &value)
e.stateDB.SetState(contractAddr, common.HexToHash(slot), common.HexToHash(value))
count++
}
if count > 0 {
fmt.Printf("[EVM] Loaded %d storage slots for %s\n", count, addrStr)
}
}

func NewPersistentStateDB(cs *ChainState) (*state.StateDB, error) {
memDB := rawdb.NewMemoryDatabase()
stateDB, err := state.New(common.Hash{}, state.NewDatabase(memDB), nil)
if err != nil {
return nil, err
}

accounts := cs.GetAllAccounts()
for _, acc := range accounts {
addr := common.HexToAddress(acc.Address)
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
balanceWei := new(big.Int).Mul(big.NewInt(int64(acc.Balance)), decimals)
stateDB.SetBalance(addr, balanceWei)
nonce := cs.LoadNonce(acc.Address)
stateDB.SetNonce(addr, nonce)
}

contracts := cs.GetAllContracts()
for _, addrStr := range contracts {
addr := common.HexToAddress(addrStr)
bytecode, err := cs.LoadContract(addrStr)
if err == nil && bytecode != nil {
stateDB.SetCode(addr, bytecode)
fmt.Printf("[EVM] Loaded contract: %s (%d bytes)\n", addrStr, len(bytecode))

// Load storage slots
if cs.db != nil {
rows, err := cs.db.Query(`SELECT slot, value FROM evm_storage WHERE address = $1`, addrStr)
if err == nil {
for rows.Next() {
var slot, value string
rows.Scan(&slot, &value)
stateDB.SetState(addr, common.HexToHash(slot), common.HexToHash(value))
}
rows.Close()
}
}
}
}

// Commit loaded state into trie so it persists across calls
stateDB.Commit(0, false)
return stateDB, nil
}
