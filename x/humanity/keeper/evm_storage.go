package keeper

import (
"database/sql"
"encoding/hex"
"fmt"
"math/big"

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
}
}

return stateDB, nil
}
