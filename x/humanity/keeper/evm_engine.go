package keeper

import (
"encoding/hex"
"fmt"
"math/big"
"strings"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/core/rawdb"
"github.com/ethereum/go-ethereum/core/state"
"github.com/ethereum/go-ethereum/core/vm"
"github.com/ethereum/go-ethereum/params"
)

// EVMEngine wraps go-ethereum EVM for contract deployment and calls.
// Design principle: every operation gets a fresh StateDB loaded from PostgreSQL.
// This avoids all stale-trie issues at the cost of slightly more DB reads.
type EVMEngine struct {
chainState *ChainState
}

func NewEVMEngine(cs *ChainState) (*EVMEngine, error) {
cs.InitV6StateTables()
e := &EVMEngine{chainState: cs}
e.RestoreV6FromMirror()
return e, nil
}

// ─── CHAIN CONFIG ─────────────────────────────────────────────────────────────

func chainConfig() *params.ChainConfig {
shanghai := uint64(0)
return &params.ChainConfig{
ChainID:             big.NewInt(1926),
HomesteadBlock:      big.NewInt(0),
EIP150Block:         big.NewInt(0),
EIP155Block:         big.NewInt(0),
EIP158Block:         big.NewInt(0),
ByzantiumBlock:      big.NewInt(0),
ConstantinopleBlock: big.NewInt(0),
PetersburgBlock:     big.NewInt(0),
IstanbulBlock:       big.NewInt(0),
BerlinBlock:         big.NewInt(0),
LondonBlock:         big.NewInt(0),
ShanghaiTime:        &shanghai,
}
}

func blockContext() vm.BlockContext {
return vm.BlockContext{
CanTransfer: func(_ vm.StateDB, _ common.Address, _ *big.Int) bool { return true },
Transfer:    func(_ vm.StateDB, _, _ common.Address, _ *big.Int) {},
GetHash:     func(_ uint64) common.Hash { return common.Hash{} },
Coinbase:    common.Address{},
BlockNumber: big.NewInt(1),
Time:        1000000,
Difficulty:  big.NewInt(0),
GasLimit:    30_000_000,
BaseFee:     big.NewInt(0),
}
}

// ─── FRESH STATE DB ───────────────────────────────────────────────────────────

// newStateDB creates a fresh in-memory StateDB loaded from PostgreSQL.
// Called before every Deploy or Call to ensure consistent state.
func (e *EVMEngine) newStateDB() (*state.StateDB, error) {
memDB := rawdb.NewMemoryDatabase()
db := state.NewDatabase(memDB)
sdb, err := state.New(common.Hash{}, db, nil)
if err != nil {
return nil, err
}

// Load all account balances
for _, acc := range e.chainState.GetAllAccounts() {
addr := common.HexToAddress(acc.Address)
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
wei := new(big.Int).Mul(big.NewInt(int64(acc.Balance)), decimals)
sdb.SetBalance(addr, wei)
sdb.SetNonce(addr, e.chainState.LoadNonce(acc.Address))
}

// Load all contract bytecodes and storage
for _, addrStr := range e.chainState.GetAllContracts() {
addr := common.HexToAddress(addrStr)

code, err := e.chainState.LoadContract(addrStr)
if err != nil || len(code) == 0 {
continue
}
sdb.SetCode(addr, code)

// Load storage slots
if e.chainState.db != nil {
rows, err := e.chainState.db.Query(
`SELECT slot, value FROM evm_storage WHERE address = $1`, addrStr)
if err == nil {
for rows.Next() {
var slot, val string
rows.Scan(&slot, &val)
sdb.SetState(addr, common.HexToHash(slot), common.HexToHash(val))
}
rows.Close()
}
}
}

// Don't commit — keep state in dirty/pending form so EVM can read it directly
return sdb, nil
}

// ─── DEPLOY ───────────────────────────────────────────────────────────────────

func (e *EVMEngine) DeployContract(from common.Address, bytecode []byte, value *big.Int) (contractAddr common.Address, ret []byte, err error) {
defer func() {
if r := recover(); r != nil {
err = fmt.Errorf("EVM panic: %v", r)
contractAddr = common.Address{}
ret = nil
}
}()

sdb, err := e.newStateDB()
if err != nil {
return common.Address{}, nil, fmt.Errorf("stateDB: %w", err)
}

nonce := e.chainState.LoadNonce(strings.ToLower(from.Hex()))
sdb.SetNonce(from, nonce)

txCtx := vm.TxContext{Origin: from, GasPrice: big.NewInt(0)}
evm := vm.NewEVM(blockContext(), txCtx, sdb, chainConfig(), vm.Config{})

_, contractAddr, _, err = evm.Create(
vm.AccountRef(from),
bytecode,
30_000_000,
value,
)
if err != nil {
return common.Address{}, nil, fmt.Errorf("deploy failed: %w", err)
}

// Commit to get runtime code into trie
sdb.Commit(0, false)

// Read runtime bytecode directly from stateDB after commit
runtimeCode := sdb.GetCode(contractAddr)
if len(runtimeCode) == 0 {
return common.Address{}, nil, fmt.Errorf("deploy succeeded but no runtime code found")
}

addrStr := strings.ToLower(contractAddr.Hex())
fromStr := strings.ToLower(from.Hex())

// Persist to PostgreSQL
e.chainState.SaveContract(addrStr, runtimeCode, fromStr)
e.chainState.SaveNonce(fromStr, nonce+1)

// Persist ALL storage slots from stateDB to PostgreSQL
// We iterate known slots by checking the stateDB journal
// For immutable contracts, storage is set during constructor
// We save slots 0-99 to catch all constructor-set values
for i := int64(0); i < 100; i++ {
slot := common.BigToHash(big.NewInt(i))
val := sdb.GetState(contractAddr, slot)
if val != (common.Hash{}) {
e.chainState.SaveStorageSlot(addrStr, slot.Hex(), val.Hex())
fmt.Printf("[EVM] Storage slot %d = %s\n", i, val.Hex())
}
}

fmt.Printf("[EVM] ✓ Deployed %s (%d bytes)\n", contractAddr.Hex(), len(runtimeCode))
return contractAddr, runtimeCode, nil
}

// ─── CALL ─────────────────────────────────────────────────────────────────────

func (e *EVMEngine) CallContract(from, to common.Address, data []byte, value *big.Int) (ret []byte, err error) {
defer func() {
if r := recover(); r != nil {
err = fmt.Errorf("EVM panic: %v", r)
ret = nil
}
}()

sdb, err := e.newStateDB()
if err != nil {
return nil, fmt.Errorf("stateDB: %w", err)
}

// Verify contract code is loaded
code := sdb.GetCode(to)
fmt.Printf("[EVM] CallContract to=%s codeLen=%d data=%x\n",
to.Hex(), len(code), data[:min4b(len(data), 4)])

if len(code) == 0 {
return nil, fmt.Errorf("no code at %s", to.Hex())
}

txCtx := vm.TxContext{Origin: from, GasPrice: big.NewInt(0)}
evm := vm.NewEVM(blockContext(), txCtx, sdb, chainConfig(), vm.Config{})

ret, _, err = evm.Call(
vm.AccountRef(from),
to,
data,
30_000_000,
value,
)
if err != nil {
return nil, fmt.Errorf("call failed: %w", err)
}

// Persist any state changes from the call
sdb.Commit(0, false)
e.persistStorageFromDB(sdb, to)
e.syncBalancesFromDB(sdb)

fmt.Printf("[EVM] Call result: %d bytes: %x\n", len(ret), ret)
return ret, nil
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func (e *EVMEngine) GetCode(addr common.Address) []byte {
code, _ := e.chainState.LoadContract(strings.ToLower(addr.Hex()))
return code
}

func (e *EVMEngine) SetCode(addr common.Address, code []byte) {
// No-op: code is always loaded fresh from DB
}

// persistStorageFromDB reads dirty storage slots from stateDB and saves to PostgreSQL
func (e *EVMEngine) persistStorageFromDB(sdb *state.StateDB, addr common.Address) {
e.PersistContractStorage(addr)
}

// syncBalancesFromDB updates PostgreSQL balances from stateDB after a state-changing call
func (e *EVMEngine) syncBalancesFromDB(sdb *state.StateDB) {
accounts := e.chainState.GetAllAccounts()
decimals := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
for _, acc := range accounts {
addr := common.HexToAddress(acc.Address)
wei := sdb.GetBalance(addr)
if wei == nil || wei.Sign() == 0 {
continue
}
balAEQ, _ := new(big.Float).Quo(new(big.Float).SetInt(wei), decimals).Float64()
if balAEQ > 0 && balAEQ != acc.Balance {
e.chainState.SetBalance(acc.Address, balAEQ)
}
}
}

func (e *EVMEngine) LoadContractStorage(addr common.Address) {
// No-op: storage loaded in newStateDB()
}

func HexToBytecode(hexStr string) ([]byte, error) {
hexStr = strings.TrimPrefix(hexStr, "0x")
return hex.DecodeString(hexStr)
}

func min4b(a, b int) int {
if a < b {
return a
}
return b
}
