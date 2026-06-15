package keeper

import (
"encoding/hex"
"fmt"
"math/big"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/crypto"
)

const V6_CONTRACT_ADDR = "0x371C577B1e2c49A07123B32F556bCcdf79317A0C"

// After each registerHuman TX on V6, mirror the state to PostgreSQL
func (e *EVMEngine) MirrorV6Registration(wallet, commitment string) {
e.chainState.SaveV6Human(wallet, commitment)
e.chainState.SaveV6Commitment(commitment, wallet)

// Also save 1000 AEQ balance in Wei
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
grant := new(big.Int).Mul(big.NewInt(1000), decimals)
e.chainState.SaveV6Balance(wallet, hex.EncodeToString(grant.Bytes()))

// Update totalHumans counter
humans := e.chainState.GetAllV6Humans()
e.chainState.SaveV6State("totalHumans", fmt.Sprintf("%x", len(humans)))

fmt.Printf("[V6] Mirrored registration: %s\n", wallet)
}

// At node start, restore V6 EVM state from PostgreSQL mirror
func (e *EVMEngine) RestoreV6FromMirror() {
contractAddr := common.HexToAddress(V6_CONTRACT_ADDR)

humans := e.chainState.GetAllV6Humans()
if len(humans) == 0 {
return
}

fmt.Printf("[V6] Restoring %d registrations to EVM...\n", len(humans))

for _, human := range humans {
walletAddr := common.HexToAddress(human["address"])
commitment := human["commitment"]

// Restore isHuman[wallet] = true
isHumanSlot := mappingSlot(walletAddr.Bytes(), 3)
e.stateDB.SetState(contractAddr, isHumanSlot, common.HexToHash("0x01"))

// Restore commitmentOf[wallet] 
if commitment != "" {
commitmentSlot := mappingSlot(walletAddr.Bytes(), 14)
commitBig := new(big.Int)
commitBig.SetString(commitment, 16)
e.stateDB.SetState(contractAddr, commitmentSlot, common.BigToHash(commitBig))
}

// Restore usedCommitments[commitment] = true
if commitment != "" {
commitBig := new(big.Int)
commitBig.SetString(commitment, 16)
commitHash := common.BigToHash(commitBig)
usedSlot := mappingSlotBytes32(commitHash, 2)
e.stateDB.SetState(contractAddr, usedSlot, common.HexToHash("0x01"))
}

// Restore balanceOf[wallet] = 1000 AEQ in Wei
balWeiHex := e.chainState.LoadV6Balance(human["address"])
if balWeiHex != "" {
balBig := new(big.Int)
balBig.SetString(balWeiHex, 16)
balSlot := mappingSlot(walletAddr.Bytes(), 1)
e.stateDB.SetState(contractAddr, balSlot, common.BigToHash(balBig))
}
}

// Restore totalHumans (storage slot 9 in V6)
totalHumansHex := e.chainState.LoadV6State("totalHumans")
if totalHumansHex != "" {
n := new(big.Int)
n.SetString(totalHumansHex, 16)
slot9 := common.BigToHash(big.NewInt(9))
e.stateDB.SetState(contractAddr, slot9, common.BigToHash(n))
}

e.stateDB.Commit(1, false)
fmt.Printf("[V6] ✓ EVM state restored for %d humans\n", len(humans))
}

// mappingSlot calculates keccak256(abi.encode(key, slot))
func mappingSlot(key []byte, slot int64) common.Hash {
// Pad key to 32 bytes
paddedKey := make([]byte, 32)
copy(paddedKey[32-len(key):], key)
// Encode slot as 32 bytes
slotBytes := common.BigToHash(big.NewInt(slot)).Bytes()
// Concatenate and hash
data := append(paddedKey, slotBytes...)
return common.BytesToHash(crypto.Keccak256(data))
}

func mappingSlotBytes32(key common.Hash, slot int64) common.Hash {
slotBytes := common.BigToHash(big.NewInt(slot)).Bytes()
data := append(key.Bytes(), slotBytes...)
return common.BytesToHash(crypto.Keccak256(data))
}
