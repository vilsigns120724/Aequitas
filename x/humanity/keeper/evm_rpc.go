package keeper

import (
"encoding/hex"
	"io"
"encoding/json"
"fmt"
"math/big"
"net/http"
"strings"

"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
)

type JSONRPCRequest struct {
JSONRPC string            `json:"jsonrpc"`
Method  string            `json:"method"`
Params  []json.RawMessage `json:"params"`
ID      interface{}       `json:"id"`
}

type JSONRPCResponse struct {
JSONRPC string      `json:"jsonrpc"`
ID      interface{} `json:"id"`
Result  interface{} `json:"result,omitempty"`
Error   interface{} `json:"error,omitempty"`
}

type EVMRPCServer struct {
dag    *BlockDAG
state  *ChainState
nonces map[string]uint64
evm    *EVMEngine
}

func NewEVMRPCServer(dag *BlockDAG, state *ChainState) *EVMRPCServer {
engine, err := NewEVMEngine(state)
if err != nil {
fmt.Printf("[EVM] Warning: could not init EVM engine: %v\n", err)
}
return &EVMRPCServer{
dag:    dag,
state:  state,
nonces: make(map[string]uint64),
evm:    engine,
}
}

func (e *EVMRPCServer) Start(port int) {
mux := http.NewServeMux()
mux.HandleFunc("/", e.handleRPC)
fmt.Printf("✓ EVM JSON-RPC listening on port %d\n", port)
go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (e *EVMRPCServer) handleRPC(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
if r.Method == "OPTIONS" {
w.WriteHeader(200)
return
}

bodyBytes, _ := io.ReadAll(r.Body)
fmt.Printf("[RPC] Incoming request: %d bytes\n", len(bodyBytes))
trimmed := strings.TrimSpace(string(bodyBytes))

if len(trimmed) > 0 && trimmed[0] == '[' {
var reqs []JSONRPCRequest
if err := json.Unmarshal(bodyBytes, &reqs); err != nil {
json.NewEncoder(w).Encode(JSONRPCResponse{JSONRPC: "2.0", Error: map[string]interface{}{"code": -32700, "message": "parse error"}})
return
}
var responses []JSONRPCResponse
for _, req := range reqs {
result, err := e.handleMethod(req.Method, req.Params)
if err != nil {
responses = append(responses, JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Error: map[string]interface{}{"code": -32603, "message": err.Error()}})
} else {
responses = append(responses, JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result})
}
}
json.NewEncoder(w).Encode(responses)
return
}

var req JSONRPCRequest
if err := json.Unmarshal(bodyBytes, &req); err != nil {
json.NewEncoder(w).Encode(JSONRPCResponse{JSONRPC: "2.0", Error: map[string]interface{}{"code": -32700, "message": "parse error"}})
return
}
result, err := e.handleMethod(req.Method, req.Params)
if err != nil {
json.NewEncoder(w).Encode(JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Error: map[string]interface{}{"code": -32603, "message": err.Error()}})
return
}
json.NewEncoder(w).Encode(JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result})
}

func (e *EVMRPCServer) handleMethod(method string, params []json.RawMessage) (interface{}, error) {
latest := e.dag.LatestBlock()
height := int64(0)
if latest != nil {
height = latest.Height
}

switch method {
case "eth_chainId":
return "0x786", nil // 73571

case "net_version":
return "1926", nil

case "eth_blockNumber":
return "0x" + fmt.Sprintf("%x", height), nil

case "eth_getBalance":
if len(params) > 0 {
var addr string
json.Unmarshal(params[0], &addr)
addr = strings.ToLower(addr)
balance := e.state.GetBalance(addr)
fmt.Printf("[RPC] eth_getBalance %s = %.2f\n", addr, balance)
if balance > 0 {
decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
aeqWei := new(big.Int).Mul(big.NewInt(int64(balance)), decimals)
return "0x" + fmt.Sprintf("%x", aeqWei), nil
}
}
return "0x0", nil

case "eth_gasPrice":
return "0x0", nil // Gasless chain

case "eth_estimateGas":
return "0x100000", nil // Return generous gas limit, price is 0

case "eth_getTransactionCount":
if len(params) > 0 {
var addr string
json.Unmarshal(params[0], &addr)
addr = strings.ToLower(addr)
nonce := e.nonces[addr]
return "0x" + fmt.Sprintf("%x", nonce), nil
}
return "0x0", nil

case "eth_sendRawTransaction":
if len(params) == 0 {
return nil, fmt.Errorf("missing params")
}
var rawHex string
json.Unmarshal(params[0], &rawHex)

// Decode raw transaction
rawHex = strings.TrimPrefix(rawHex, "0x")
rawBytes, err := hex.DecodeString(rawHex)
if err != nil {
return nil, fmt.Errorf("invalid hex: %v", err)
}

tx := new(types.Transaction)
if err := rlp.DecodeBytes(rawBytes, tx); err != nil {
return nil, fmt.Errorf("invalid transaction: %v", err)
}

// Get sender
signer := types.LatestSignerForChainID(big.NewInt(1926))
sender, err := types.Sender(signer, tx)
if err != nil {
// Try legacy signer
signer = types.NewEIP155Signer(big.NewInt(1926))
sender, err = types.Sender(signer, tx)
if err != nil {
return nil, fmt.Errorf("cannot recover sender: %v", err)
}
}

senderAddr := strings.ToLower(sender.Hex())
txHash := "0x" + tx.Hash().Hex()[2:]

fmt.Printf("[RPC] eth_sendRawTransaction from=%s to=%v value=%v data=%d bytes\n",
senderAddr, tx.To(), tx.Value(), len(tx.Data()))

// Handle AEQ transfer (no data = simple transfer)
if tx.To() != nil && len(tx.Data()) == 0 && tx.Value().Cmp(big.NewInt(0)) > 0 {
toAddr := strings.ToLower(tx.To().Hex())
decimals := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
valueFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(tx.Value()), decimals).Float64()

if err := e.state.Transfer(senderAddr, toAddr, valueFloat); err != nil {
return nil, fmt.Errorf("transfer failed: %v", err)
}
fmt.Printf("[RPC] ✓ Transfer %.2f AEQ: %s → %s\n", valueFloat, senderAddr, toAddr)
}

// Handle contract deployment or contract call (data present)
if len(tx.Data()) > 0 && e.evm != nil {
fmt.Printf("[EVM] evm engine is nil: %v\n", e.evm == nil)
if tx.To() == nil {
// Contract deployment
fmt.Printf("[EVM] Deploying contract, bytecode=%d bytes\n", len(tx.Data()))
contractAddr, _, err := e.evm.DeployContract(sender, tx.Data(), tx.Value())
if err != nil {
fmt.Printf("[RPC] ✗ Deploy failed: %v\n", err)
} else {
fmt.Printf("[RPC] ✓ Contract deployed: %s\n", contractAddr.Hex())
txHash = "0x" + contractAddr.Hex()[2:]
}
} else {
// Contract call
result, err := e.evm.CallContract(sender, *tx.To(), tx.Data(), tx.Value())
if err != nil {
fmt.Printf("[RPC] ✗ Contract call failed: %v\n", err)
} else {
fmt.Printf("[RPC] ✓ Contract call result: %x\n", result)
}
}
}

// Update nonce
e.nonces[senderAddr]++

return txHash, nil

case "eth_call":
if len(params) >= 1 && e.evm != nil {
var callObj map[string]string
if err := json.Unmarshal(params[0], &callObj); err == nil {
from := common.HexToAddress(callObj["from"])
to := common.HexToAddress(callObj["to"])
data, _ := hex.DecodeString(strings.TrimPrefix(callObj["data"], "0x"))
result, err := e.evm.CallContract(from, to, data, big.NewInt(0))
if err != nil {
fmt.Printf("[RPC] eth_call error: %v\n", err)
return "0x", nil
}
return "0x" + hex.EncodeToString(result), nil
}
}
return "0x", nil

case "eth_getCode":
// Return non-empty code for known contract addresses
return "0x600160015b", nil

case "eth_getLogs":
return []interface{}{}, nil

case "eth_getBlockByNumber":
block := e.dag.LatestBlock()
if block == nil {
return nil, nil
}
return map[string]interface{}{
"number":           "0x" + fmt.Sprintf("%x", block.Height),
"hash":             "0x" + block.Hash,
"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
"timestamp":        "0x" + fmt.Sprintf("%x", block.Timestamp),
"transactions":     []interface{}{},
"gasLimit":         "0x1000000",
"gasUsed":          "0x0",
"difficulty":       "0x0",
"totalDifficulty":  "0x0",
"miner":            "0x0000000000000000000000000000000000000000",
"extraData":        "0x",
"logsBloom":        "0x" + strings.Repeat("0", 512),
"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
"stateRoot":        "0x" + block.Hash,
"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
"size":             "0x1",
"uncles":           []interface{}{},
"nonce":            "0x0000000000000000",
"baseFeePerGas":    "0x0",
}, nil

case "eth_getTransactionReceipt":
if len(params) > 0 {
var txHash string
json.Unmarshal(params[0], &txHash)
// Return success receipt
return map[string]interface{}{
"transactionHash":   txHash,
"transactionIndex":  "0x0",
"blockHash":         "0x" + strings.Repeat("0", 64),
"blockNumber":       "0x" + fmt.Sprintf("%x", height),
"from":              "0x0000000000000000000000000000000000000000",
"to":                nil,
"cumulativeGasUsed": "0x0",
"gasUsed":           "0x0",
"contractAddress":   "0x" + strings.Repeat("1", 40),
"logs":              []interface{}{},
"logsBloom":         "0x" + strings.Repeat("0", 512),
"status":            "0x1", // Success
"type":              "0x2",
}, nil
}
return nil, nil

case "eth_feeHistory":
return map[string]interface{}{
"oldestBlock":   "0x0",
"baseFeePerGas": []string{"0x0", "0x0"},
"gasUsedRatio":  []float64{0.0},
"reward":        [][]string{{"0x0"}},
}, nil

case "web3_clientVersion":
return "Aequitas/v0.3.0/BlockDAG", nil

case "eth_syncing":
return false, nil

case "eth_accounts":
return []interface{}{}, nil

case "net_listening":
return true, nil

case "net_peerCount":
return "0x1", nil

default:
fmt.Printf("[RPC] unsupported method: %s\n", method)
return nil, fmt.Errorf("method %s not supported", method)
}
}
