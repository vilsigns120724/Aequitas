package keeper

import (
"encoding/hex"
"encoding/json"
"fmt"
"io"
"math/big"
"net/http"
"strings"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/crypto"
)

type RegisterRequest struct {
Wallet     string     `json:"wallet"`
PubSignals []string   `json:"pubSignals"`
PA         []string   `json:"pA"`
PB         [][]string `json:"pB"`
PC         []string   `json:"pC"`
}

type RegisterResponse struct {
Success bool    `json:"success"`
Message string  `json:"message"`
Balance float64 `json:"balance"`
TxHash  string  `json:"tx_hash"`
}

func (a *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

if r.Method == "OPTIONS" {
w.WriteHeader(200)
return
}
if r.Method != "POST" {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "POST required"})
return
}

body, _ := io.ReadAll(r.Body)
var req RegisterRequest
if err := json.Unmarshal(body, &req); err != nil {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "invalid request"})
return
}

wallet := strings.ToLower(req.Wallet)
if wallet == "" {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "wallet required"})
return
}

// Check if already registered
if a.state.IsHuman(wallet) {
json.NewEncoder(w).Encode(RegisterResponse{
Success: false,
Message: "already registered on Aequitas Chain",
Balance: a.state.GetBalance(wallet),
})
return
}

fmt.Printf("[REGISTER] Registering wallet: %s\n", wallet)

// If ZKP proof is provided, call V6 contract
evmRPC := NewEVMRPCServer(a.blockchain, a.state)
if evmRPC.evm != nil && len(req.PubSignals) >= 2 {
txHash, err := a.registerOnV6(evmRPC.evm, wallet, req)
if err != nil {
fmt.Printf("[REGISTER] V6 error: %v - falling back to BlockDAG\n", err)
} else {
// Mirror to V6 state
commitment := req.PubSignals[1]
evmRPC.evm.MirrorV6Registration(wallet, commitment)

// Also register in BlockDAG for consistency
a.state.RegisterHuman(wallet)
a.blockchain.AddTransaction(Transaction{
Type:   "register_human_v6",
Wallet: wallet,
Amount: 1000,
TxHash: txHash,
})

fmt.Printf("[REGISTER] ✓ Human registered on V6: %s\n", wallet)
json.NewEncoder(w).Encode(RegisterResponse{
Success: true,
Message: "Registered as human on Aequitas V6! 1,000 AEQ granted.",
Balance: 1000,
TxHash:  txHash,
})
return
}
}

// Fallback: register only in BlockDAG (gasless, no ZKP required)
a.state.RegisterHuman(wallet)
txHash := fmt.Sprintf("0x%064x", len(wallet))
a.blockchain.AddTransaction(Transaction{
Type:   "register_human",
Wallet: wallet,
Amount: 1000,
TxHash: txHash,
})

fmt.Printf("[REGISTER] ✓ Human registered (BlockDAG): %s | 1000 AEQ\n", wallet)
json.NewEncoder(w).Encode(RegisterResponse{
Success: true,
Message: "Registered as human! 1,000 AEQ granted.",
Balance: 1000,
TxHash:  txHash,
})
}

// registerOnV6 calls registerHuman() on the V6 EVM contract
func (a *APIServer) registerOnV6(evm *EVMEngine, wallet string, req RegisterRequest) (string, error) {
if len(req.PA) < 2 || len(req.PB) < 2 || len(req.PC) < 2 || len(req.PubSignals) < 2 {
return "", fmt.Errorf("incomplete ZKP proof")
}

// Build calldata for registerHuman(uint[2],uint[2][2],uint[2],uint[2])
// Function selector: keccak256("registerHuman(uint256[2],uint256[2][2],uint256[2],uint256[2])")[:4]
selector := crypto.Keccak256([]byte("registerHuman(uint256[2],uint256[2][2],uint256[2],uint256[2])"))[0:4]

// Encode parameters
calldata := selector
calldata = append(calldata, encodeUint256Array2(req.PA)...)
calldata = append(calldata, encodeUint256Array2x2(req.PB)...)
calldata = append(calldata, encodeUint256Array2(req.PC)...)
calldata = append(calldata, encodeUint256Array2(req.PubSignals)...)

from := common.HexToAddress(wallet)
to := common.HexToAddress(V6_CONTRACT_ADDR)

result, err := evm.CallContract(from, to, calldata, big.NewInt(0))
if err != nil {
return "", fmt.Errorf("V6 call failed: %v", err)
}

// Generate TX hash
txData := append(calldata, from.Bytes()...)
txHash := "0x" + hex.EncodeToString(crypto.Keccak256(txData))

_ = result
_ = types.Transaction{}

return txHash, nil
}

func encodeUint256Array2(values []string) []byte {
result := make([]byte, 64)
for i := 0; i < 2 && i < len(values); i++ {
n := new(big.Int)
n.SetString(values[i], 10)
b := common.BigToHash(n).Bytes()
copy(result[i*32:], b)
}
return result
}

func encodeUint256Array2x2(values [][]string) []byte {
result := make([]byte, 128)
for i := 0; i < 2 && i < len(values); i++ {
for j := 0; j < 2 && j < len(values[i]); j++ {
n := new(big.Int)
n.SetString(values[i][j], 10)
b := common.BigToHash(n).Bytes()
copy(result[(i*2+j)*32:], b)
}
}
return result
}
