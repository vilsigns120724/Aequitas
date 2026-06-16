package keeper

import (
"encoding/hex"
"fmt"
"math/big"
"strings"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/crypto"
)

const V6_CONTRACT_ADDR = "0xA76cA3bf34F2Ae5dFA0608696627e42b81180488"
const V7_CONTRACT_ADDR = "0xc7553c86D2EfE4771d35880210FA4f93AC0Ea491"

// MirrorV6Registration mirrors a V6 registration to PostgreSQL
func (e *EVMEngine) MirrorV6Registration(wallet, commitment string) {
e.chainState.SaveV6Human(wallet, commitment)
e.chainState.SaveV6Commitment(commitment, wallet)

decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
grant := new(big.Int).Mul(big.NewInt(1000), decimals)
e.chainState.SaveV6Balance(wallet, hex.EncodeToString(grant.Bytes()))

humans := e.chainState.GetAllV6Humans()
e.chainState.SaveV6State("totalHumans", fmt.Sprintf("%x", len(humans)))

fmt.Printf("[V6] Mirrored registration: %s\n", wallet)
}

// RestoreV6FromMirror restores V6 EVM state from PostgreSQL to evm_storage table
// so that CallContract can read it via newStateDB()
func (e *EVMEngine) RestoreV6FromMirror() {
contractAddr := common.HexToAddress(V6_CONTRACT_ADDR)
addrStr := strings.ToLower(contractAddr.Hex())

humans := e.chainState.GetAllV6Humans()
if len(humans) == 0 {
return
}

fmt.Printf("[V6] Restoring %d registrations to storage...\n", len(humans))

for _, human := range humans {
walletAddr := common.HexToAddress(human["address"])
commitment := human["commitment"]

// isHuman[wallet] = true (slot 3)
isHumanSlot := mappingSlot(walletAddr.Bytes(), 3)
e.chainState.SaveStorageSlot(addrStr, isHumanSlot.Hex(), common.HexToHash("0x01").Hex())

// commitmentOf[wallet] (slot 14)
if commitment != "" {
commitmentSlot := mappingSlot(walletAddr.Bytes(), 14)
commitBig := new(big.Int)
commitBig.SetString(commitment, 16)
e.chainState.SaveStorageSlot(addrStr, commitmentSlot.Hex(), common.BigToHash(commitBig).Hex())

// usedCommitments[commitment] = true (slot 2)
commitHash := common.BigToHash(commitBig)
usedSlot := mappingSlotBytes32(commitHash, 2)
e.chainState.SaveStorageSlot(addrStr, usedSlot.Hex(), common.HexToHash("0x01").Hex())
}

// balanceOf[wallet] (slot 1)
balWeiHex := e.chainState.LoadV6Balance(human["address"])
if balWeiHex != "" {
balBig := new(big.Int)
balBig.SetString(balWeiHex, 16)
balSlot := mappingSlot(walletAddr.Bytes(), 1)
e.chainState.SaveStorageSlot(addrStr, balSlot.Hex(), common.BigToHash(balBig).Hex())
}
}

// totalHumans (slot 9)
totalHumansHex := e.chainState.LoadV6State("totalHumans")
if totalHumansHex != "" {
n := new(big.Int)
n.SetString(totalHumansHex, 16)
slot9 := common.BigToHash(big.NewInt(9))
e.chainState.SaveStorageSlot(addrStr, slot9.Hex(), common.BigToHash(n).Hex())
}

fmt.Printf("[V6] ✓ Storage restored for %d humans\n", len(humans))
}

func mappingSlot(key []byte, slot int64) common.Hash {
paddedKey := make([]byte, 32)
copy(paddedKey[32-len(key):], key)
slotBytes := common.BigToHash(big.NewInt(slot)).Bytes()
data := append(paddedKey, slotBytes...)
return common.BytesToHash(crypto.Keccak256(data))
}

func mappingSlotBytes32(key common.Hash, slot int64) common.Hash {
slotBytes := common.BigToHash(big.NewInt(slot)).Bytes()
data := append(key.Bytes(), slotBytes...)
return common.BytesToHash(crypto.Keccak256(data))
}
