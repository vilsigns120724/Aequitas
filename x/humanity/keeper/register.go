package keeper

import (
"encoding/json"
"fmt"
"io"
"math/big"
"net/http"
"os"
"strings"

"github.com/ethereum/go-ethereum/accounts/abi"
"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/crypto"
)

// V7 ABI fragment for registerWithSig — used only to encode calldata correctly,
// including the dynamic `bytes signature` parameter (offset/length tail encoding).
const registerWithSigABI = `[{
"name": "registerWithSig",
"type": "function",
"inputs": [
{"name": "pA", "type": "uint256[2]"},
{"name": "pB", "type": "uint256[2][2]"},
{"name": "pC", "type": "uint256[2]"},
{"name": "pubSignals", "type": "uint256[2]"},
{"name": "claimedHuman", "type": "address"},
{"name": "signature", "type": "bytes"}
]
}]`

type RegisterRequest struct {
Wallet     string     `json:"wallet"`
PubSignals []string   `json:"pubSignals"`
PA         []string   `json:"pA"`
PB         [][]string `json:"pB"`
PC         []string   `json:"pC"`
Signature  string     `json:"signature"` // hex-encoded, 65 bytes, from personal_sign
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
if len(req.PA) < 2 || len(req.PB) < 2 || len(req.PC) < 2 || len(req.PubSignals) < 2 {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "incomplete ZK proof"})
return
}
if req.Signature == "" {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "signature required"})
return
}

fmt.Printf("[REGISTER] Relaying registerWithSig for: %s\n", wallet)

evmRPC := NewEVMRPCServer(a.blockchain, a.state)
if evmRPC.evm == nil {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "EVM engine unavailable"})
return
}

txHash, err := a.registerOnV7(evmRPC, wallet, req)
if err != nil {
fmt.Printf("[REGISTER] V7 registerWithSig failed: %v\n", err)
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: err.Error()})
return
}

fmt.Printf("[REGISTER] ✓ Relayed registerWithSig for %s, tx=%s\n", wallet, txHash)
json.NewEncoder(w).Encode(RegisterResponse{
Success: true,
Message: "Registered as human on Aequitas V7! 1,000 AEQ granted.",
Balance: 1000,
TxHash:  txHash,
})
}

// registerOnV7 builds, signs (with the relayer key), and submits a registerWithSig
// transaction on the V7 contract. claimedHuman in the contract call is the user's
// own wallet — verified via the signature, not via msg.sender (which is the relayer).
func (a *APIServer) registerOnV7(evmRPC *EVMRPCServer, wallet string, req RegisterRequest) (string, error) {
relayerPK := os.Getenv("RELAYER_PRIVATE_KEY")
if relayerPK == "" {
return "", fmt.Errorf("server misconfigured: RELAYER_PRIVATE_KEY not set")
}
relayerPK = strings.TrimPrefix(relayerPK, "0x")

relayerKey, err := crypto.HexToECDSA(relayerPK)
if err != nil {
return "", fmt.Errorf("invalid relayer key: %w", err)
}
relayerAddr := crypto.PubkeyToAddress(relayerKey.PublicKey)

parsedABI, err := abi.JSON(strings.NewReader(registerWithSigABI))
if err != nil {
return "", fmt.Errorf("abi parse failed: %w", err)
}

pA, err := parseUint2(req.PA)
if err != nil {
return "", fmt.Errorf("invalid pA: %w", err)
}
pB, err := parseUint2x2(req.PB)
if err != nil {
return "", fmt.Errorf("invalid pB: %w", err)
}
pC, err := parseUint2(req.PC)
if err != nil {
return "", fmt.Errorf("invalid pC: %w", err)
}
pubSignals, err := parseUint2(req.PubSignals)
if err != nil {
return "", fmt.Errorf("invalid pubSignals: %w", err)
}

sigHex := strings.TrimPrefix(req.Signature, "0x")
sigBytes := common.Hex2Bytes(sigHex)
if len(sigBytes) != 65 {
return "", fmt.Errorf("signature must be 65 bytes, got %d", len(sigBytes))
}

claimedHuman := common.HexToAddress(wallet)

calldata, err := parsedABI.Pack("registerWithSig", pA, pB, pC, pubSignals, claimedHuman, sigBytes)
if err != nil {
return "", fmt.Errorf("encoding failed: %w", err)
}

to := common.HexToAddress(V7_CONTRACT_ADDR)
relayerAddrStr := strings.ToLower(relayerAddr.Hex())
nonce := a.state.LoadNonce(relayerAddrStr)

tx := types.NewTx(&types.LegacyTx{
Nonce:    nonce,
To:       &to,
Value:    big.NewInt(0),
Gas:      6_000_000,
GasPrice: big.NewInt(0),
Data:     calldata,
})

signer := types.NewEIP155Signer(big.NewInt(1926))
signedTx, err := types.SignTx(tx, signer, relayerKey)
if err != nil {
return "", fmt.Errorf("relayer signing failed: %w", err)
}

rawBytes, err := signedTx.MarshalBinary()
if err != nil {
return "", fmt.Errorf("tx encoding failed: %w", err)
}

rawHex := "0x" + common.Bytes2Hex(rawBytes)
result, rpcErr := evmRPC.sendRawTransaction([]json.RawMessage{
mustMarshal(rawHex),
})
if rpcErr != nil {
return "", fmt.Errorf("submission failed: %s", rpcErr.Message)
}

txHash, ok := result.(string)
if !ok {
return "", fmt.Errorf("unexpected response from relay")
}
return txHash, nil
}

func parseUint2(values []string) ([2]*big.Int, error) {
var out [2]*big.Int
for i := 0; i < 2; i++ {
n := new(big.Int)
_, ok := n.SetString(values[i], 10)
if !ok {
return out, fmt.Errorf("invalid number at index %d: %s", i, values[i])
}
out[i] = n
}
return out, nil
}

func parseUint2x2(values [][]string) ([2][2]*big.Int, error) {
var out [2][2]*big.Int
for i := 0; i < 2; i++ {
if len(values[i]) < 2 {
return out, fmt.Errorf("row %d has fewer than 2 elements", i)
}
for j := 0; j < 2; j++ {
n := new(big.Int)
_, ok := n.SetString(values[i][j], 10)
if !ok {
return out, fmt.Errorf("invalid number at [%d][%d]: %s", i, j, values[i][j])
}
out[i][j] = n
}
}
return out, nil
}

func mustMarshal(v interface{}) json.RawMessage {
b, _ := json.Marshal(v)
return b
}
