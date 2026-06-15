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

type EVMEngine struct {
chainState *ChainState
stateDB    *state.StateDB
contracts  map[common.Address][]byte
}

func NewEVMEngine(cs *ChainState) (*EVMEngine, error) {
memDB := rawdb.NewMemoryDatabase()
stateDB, err := state.New(common.Hash{}, state.NewDatabase(memDB), nil)
if err != nil {
return nil, fmt.Errorf("failed to create stateDB: %w", err)
}

engine := &EVMEngine{
chainState: cs,
stateDB:    stateDB,
contracts:  make(map[common.Address][]byte),
}

engine.syncFromChainState()
return engine, nil
}

func (e *EVMEngine) syncFromChainState() {
accounts := e.chainState.GetAllAccounts()
for _, acc := range accounts {
addr := common.HexToAddress(acc.Address)
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
balanceWei := new(big.Int).Mul(big.NewInt(int64(acc.Balance)), decimals)
e.stateDB.SetBalance(addr, balanceWei)
}
}

func (e *EVMEngine) DeployContract(from common.Address, bytecode []byte, value *big.Int) (common.Address, []byte, error) {
shanghai := uint64(0)
chainConfig := &params.ChainConfig{
	ChainID: big.NewInt(9001),
	HomesteadBlock: big.NewInt(0),
	EIP150Block: big.NewInt(0),
	EIP155Block: big.NewInt(0),
	EIP158Block: big.NewInt(0),
	ByzantiumBlock: big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	PetersburgBlock: big.NewInt(0),
	IstanbulBlock: big.NewInt(0),
	BerlinBlock: big.NewInt(0),
	LondonBlock: big.NewInt(0),
	ShanghaiTime: &shanghai,
}

blockCtx := vm.BlockContext{
CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool { return true },
Transfer:    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {},
GetHash:     func(n uint64) common.Hash { return common.Hash{} },
Coinbase:    common.Address{},
BlockNumber: big.NewInt(1),
Time:        1000000,
Difficulty:  big.NewInt(0),
GasLimit:    10000000,
BaseFee:     big.NewInt(0),
}

txCtx := vm.TxContext{
Origin:   from,
GasPrice: big.NewInt(0),
}

evm := vm.NewEVM(blockCtx, txCtx, e.stateDB, chainConfig, vm.Config{})

nonce := e.stateDB.GetNonce(from)
ret, contractAddr, _, err := evm.Create(
vm.AccountRef(from),
bytecode,
10000000,
value,
)
if err != nil {
return common.Address{}, nil, fmt.Errorf("deployment failed: %w", err)
}

e.contracts[contractAddr] = ret
e.stateDB.SetNonce(from, nonce+1)
e.stateDB.Commit(1, false)

fmt.Printf("[EVM] ✓ Contract deployed at %s\n", contractAddr.Hex())
return contractAddr, ret, nil
}

func (e *EVMEngine) CallContract(from, to common.Address, data []byte, value *big.Int) ([]byte, error) {
shanghai2 := uint64(0)
chainConfig2 := &params.ChainConfig{
	ChainID: big.NewInt(9001),
	HomesteadBlock: big.NewInt(0),
	EIP150Block: big.NewInt(0),
	EIP155Block: big.NewInt(0),
	EIP158Block: big.NewInt(0),
	ByzantiumBlock: big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	PetersburgBlock: big.NewInt(0),
	IstanbulBlock: big.NewInt(0),
	BerlinBlock: big.NewInt(0),
	LondonBlock: big.NewInt(0),
	ShanghaiTime: &shanghai2,
}

blockCtx := vm.BlockContext{
CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool { return true },
Transfer:    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {},
GetHash:     func(n uint64) common.Hash { return common.Hash{} },
Coinbase:    common.Address{},
BlockNumber: big.NewInt(1),
Time:        1000000,
Difficulty:  big.NewInt(0),
GasLimit:    10000000,
BaseFee:     big.NewInt(0),
}

txCtx := vm.TxContext{
Origin:   from,
GasPrice: big.NewInt(0),
}

evm := vm.NewEVM(blockCtx, txCtx, e.stateDB, chainConfig2, vm.Config{})

ret, _, err := evm.Call(
vm.AccountRef(from),
to,
data,
10000000,
value,
)
if err != nil {
return nil, fmt.Errorf("call failed: %w", err)
}

e.syncToChainState()
return ret, nil
}

func (e *EVMEngine) syncToChainState() {
accounts := e.chainState.GetAllAccounts()
for _, acc := range accounts {
addr := common.HexToAddress(acc.Address)
balanceWei := e.stateDB.GetBalance(addr)
if balanceWei != nil {
decimals := new(big.Float).SetInt(
new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
)
balanceAEQ, _ := new(big.Float).Quo(
new(big.Float).SetInt(balanceWei),
decimals,
).Float64()
if balanceAEQ != acc.Balance {
e.chainState.SetBalance(acc.Address, balanceAEQ)
}
}
}
}

func HexToBytecode(hexStr string) ([]byte, error) {
hexStr = strings.TrimPrefix(hexStr, "0x")
return hex.DecodeString(hexStr)
}
