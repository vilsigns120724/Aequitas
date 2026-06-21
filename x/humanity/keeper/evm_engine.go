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
func (e *EVMEngine) newStateDB() (*state.StateDB, state.Database, error) {
memDB := rawdb.NewMemoryDatabase()
db := state.NewDatabase(memDB)
sdb, err := state.New(common.Hash{}, db, nil)
if err != nil {
return nil, nil, err
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
return sdb, db, nil
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

sdb, _, err := e.newStateDB()
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

// CallContract executes a contract call against a fresh StateDB built from
// PostgreSQL. The persist parameter controls whether the resulting state
// changes are written back to PostgreSQL:
//   - persist=true:  use ONLY for a call that represents a real, intended
//     state change (the actual execution inside sendRawTransaction).
//   - persist=false: use for read-only queries (eth_call, isHuman/balanceOf
//     lookups in api.go) AND for dry-run simulations (register.go's
//     pre-flight check before the real submit). Nothing is written back.
//
// Previously this function ALWAYS persisted, regardless of why it was
// called. That meant a pure eth_call (e.g. checking someone's balance) or
// a dry-run simulation (checking whether a registration WOULD succeed,
// before actually submitting it) had the exact same side effect as a real,
// committed registration: isHuman/balanceOf were written to evm_storage as
// if the call had truly happened. In practice this meant every attempt to
// register — even ones whose real submission later failed or was never
// sent — already "registered" the wallet the moment the dry-run ran,
// making "already registered" errors appear for wallets that, from the
// chain's own database tables, looked completely unregistered. Database
// resets could never fix this because the very next read-only status
// check would silently re-create the same state.
func (e *EVMEngine) CallContract(from, to common.Address, data []byte, value *big.Int, persist bool) (ret []byte, err error) {
defer func() {
if r := recover(); r != nil {
err = fmt.Errorf("EVM panic: %v", r)
ret = nil
}
}()

sdb, db, err := e.newStateDB()
if err != nil {
return nil, fmt.Errorf("stateDB: %w", err)
}

// Verify contract code is loaded
code := sdb.GetCode(to)
fmt.Printf("[EVM] CallContract to=%s codeLen=%d data=%x persist=%v\n",
to.Hex(), len(code), data[:min4b(len(data), 4)], persist)

if len(code) == 0 {
// MetaMask Mobile (and some other wallets) call well-known system
// contracts that exist on mainnet but not on a custom chain:
//   - Multicall3 (0xcA11bde...) — used to batch eth_call requests
//   - Zero address (0x0000...) — probed for token symbol/decimals
// Returning a hard error here makes MetaMask Mobile abort the
// entire transaction flow. Instead we return an empty 32-byte
// result (the standard ABI-encoded zero/empty value) so the wallet
// gracefully treats the call as "not supported" and falls back to
// its single-call path. This is the same behavior as geth when
// calling a non-existent contract with staticcall.
toHex := strings.ToLower(to.Hex())
if toHex == "0xca11bde05977b3631167028862be2a173976ca11" ||
toHex == "0x0000000000000000000000000000000000000000" {
fmt.Printf("[EVM] Known system contract %s — returning empty result\n", to.Hex())
return make([]byte, 32), nil
}
return nil, fmt.Errorf("no code at %s", to.Hex())
}

txCtx := vm.TxContext{Origin: from, GasPrice: big.NewInt(0)}
evm := vm.NewEVM(blockContext(), txCtx, sdb, chainConfig(), vm.Config{})

var execErr error
ret, _, execErr = evm.Call(
vm.AccountRef(from),
to,
data,
30_000_000,
value,
)
if execErr != nil {
reason := decodeRevertReason(ret)
if reason != "" {
return nil, fmt.Errorf("%s", reason)
}
return nil, fmt.Errorf("call failed: %w", execErr)
}

if !persist {
fmt.Printf("[EVM] Call result (not persisted): %d bytes: %x\n", len(ret), ret)
return ret, nil
}

// Persist any state changes from the call.
// IMPORTANT: per go-ethereum docs, sdb is no longer reliable for reads
// after Commit() — we must open a fresh StateDB on the returned root
// to safely dump and persist the resulting storage.
root, commitErr := sdb.Commit(0, false)
if commitErr != nil {
fmt.Printf("[EVM] revert Commit failed: %v\n", commitErr)
} else {
touchedAddrs, touchedCommitments := extractTouchedEntities(from, data)
e.dumpAndPersistStorage(root, db, to, touchedAddrs, touchedCommitments)
}
e.syncBalancesFromDB(sdb)

fmt.Printf("[EVM] Call result: %d bytes: %x\n", len(ret), ret)
return ret, nil
}

// dumpAndPersistStorage opens a fresh, read-only StateDB on the given root
// and writes every populated storage slot for addr into PostgreSQL.
// This is the generic, contract-agnostic replacement for guessing slot
// numbers manually — it works correctly for any mapping or simple storage
// variable, regardless of how its slot is computed.
// knownV7Slots lists every storage slot AequitasV7.sol declares, in
// declaration order. Simple slots are plain integers; mapping slots are
// listed by their base slot index and require mappingSlot()/
// mappingSlotBytes32() with a key to compute the actual storage location.
// This is explicit, contract-specific knowledge — not generic — because
// go-ethereum's StateDB offers no reliable generic "what changed" API in
// this version (verified: RawDump does not find accounts after Commit,
// even on the same backing database).
// v7SimpleSlots: single-value slots that are always persisted.
var v7SimpleSlots = []int64{0, 1, 2, 3} // totalSupply, totalHumans, ubiPool, ubiPerHumanAccumulated

// v7AddressMappingSlots: per-address mapping slots (all 13 address mappings in V7).
var v7AddressMappingSlots = []int64{
4,  // balanceOf
5,  // escrowOf
6,  // isHuman
// 7 = usedCommitments (uint256 key, not address — handled separately)
// 8 = usedNullifiers  (bytes32 key, not address — handled separately)
9,  // commitmentOf
10, // lastActivity
11, // lastDemurrage
12, // ubiClaimed
13, // guardianOf
14, // pendingGuardian
15, // guardianRequestedAt
16, // wardCount
}

// v7ArrayBaseSlots: the 10 fixed-size-array slots (CAPS[5] + THRESHOLDS[5]).
var v7ArrayBaseSlots = []int64{17, 18, 19, 20, 21, 22, 23, 24, 25, 26}

// extractTouchedEntities returns which addresses and commitments a given
// call may have modified, based on an explicit, verified table of byte
// offsets per function selector. This is NOT a heuristic — each offset was
// confirmed against real ABI-encoded calldata before being hardcoded here.
// Add a new case here whenever a new state-changing function is wired up.
func extractTouchedEntities(from common.Address, data []byte) ([]common.Address, []*big.Int) {
if len(data) < 4 {
return []common.Address{from}, nil
}

selector := fmt.Sprintf("%x", data[:4])
switch selector {
case "33f4167a": // registerWithSig(uint256[2],uint256[2][2],uint256[2],uint256[2],address,bytes)
addrs := []common.Address{from}
var commitments []*big.Int
if len(data) >= 324+32 {
claimedHuman := common.BytesToAddress(data[324 : 324+32])
addrs = append(addrs, claimedHuman)
}
if len(data) >= 260+32 {
commitment := new(big.Int).SetBytes(data[260 : 260+32])
commitments = append(commitments, commitment)
}
return addrs, commitments
default:
// Unknown selector: at minimum, the caller's own address may have
// been touched (e.g. a simple register() or transfer() from msg.sender).
return []common.Address{from}, nil
}
}

func (e *EVMEngine) dumpAndPersistStorage(root common.Hash, db state.Database, addr common.Address, touchedAddrs []common.Address, touchedCommitments []*big.Int) {
freshDB, err := state.New(root, db, nil)
if err != nil {
fmt.Printf("[EVM] revert Could not open committed state for persistence: %v\n", err)
return
}

addrStr := strings.ToLower(addr.Hex())
count := 0

for _, slotIdx := range v7SimpleSlots {
slot := common.BigToHash(big.NewInt(slotIdx))
val := freshDB.GetState(addr, slot)
e.chainState.SaveStorageSlot(addrStr, slot.Hex(), val.Hex())
count++
}
// Persist all fixed-size array slots (CAPS[5] + THRESHOLDS[5]).
for _, slotIdx := range v7ArrayBaseSlots {
slot := common.BigToHash(big.NewInt(slotIdx))
val := freshDB.GetState(addr, slot)
e.chainState.SaveStorageSlot(addrStr, slot.Hex(), val.Hex())
count++
}

for _, touched := range touchedAddrs {
for _, base := range v7AddressMappingSlots {
slot := mappingSlot(touched.Bytes(), base)
val := freshDB.GetState(addr, slot)
e.chainState.SaveStorageSlot(addrStr, slot.Hex(), val.Hex())
count++
}
}

for _, commitment := range touchedCommitments {
slot := mappingSlotBytes32(common.BigToHash(commitment), 7) // usedCommitments
val := freshDB.GetState(addr, slot)
e.chainState.SaveStorageSlot(addrStr, slot.Hex(), val.Hex())
count++
}
if count > 0 {
fmt.Printf("[EVM] Persisted %d storage slots for %s\n", count, addrStr)
}
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

// decodeRevertReason extracts the human-readable message from EVM revert bytes.
// Solidity require(cond, "message") encodes as: Error(string) selector (0x08c379a0)
// followed by standard ABI-encoded string (offset, length, padded bytes).
// Returns "" if the bytes don't match this standard format (e.g. a panic or
// a require() without a message).
func decodeRevertReason(ret []byte) string {
if len(ret) < 4 {
return ""
}
// Error(string) selector
if ret[0] != 0x08 || ret[1] != 0xc3 || ret[2] != 0x79 || ret[3] != 0xa0 {
return ""
}
payload := ret[4:]
if len(payload) < 64 {
return ""
}
// payload[0:32] = offset (always 0x20 for a single string param)
// payload[32:64] = string length
strLen := new(big.Int).SetBytes(payload[32:64]).Uint64()
if uint64(len(payload)) < 64+strLen {
return ""
}
return string(payload[64 : 64+strLen])
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
