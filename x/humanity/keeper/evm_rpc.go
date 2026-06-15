package keeper

import (
"encoding/json"
"fmt"
	"math/big"
"net/http"
"strings"
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
dag *BlockDAG
state *ChainState
}

func NewEVMRPCServer(dag *BlockDAG, state *ChainState) *EVMRPCServer {
return &EVMRPCServer{dag: dag, state: state}
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

var req JSONRPCRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
json.NewEncoder(w).Encode(JSONRPCResponse{
JSONRPC: "2.0",
Error:   map[string]interface{}{"code": -32700, "message": "parse error"},
})
return
}

result, err := e.handleMethod(req.Method, req.Params)
if err != nil {
json.NewEncoder(w).Encode(JSONRPCResponse{
JSONRPC: "2.0",
ID:      req.ID,
Error:   map[string]interface{}{"code": -32603, "message": err.Error()},
})
return
}

json.NewEncoder(w).Encode(JSONRPCResponse{
JSONRPC: "2.0",
ID:      req.ID,
Result:  result,
})
}

func (e *EVMRPCServer) handleMethod(method string, params []json.RawMessage) (interface{}, error) {
latest := e.dag.LatestBlock()
height := int64(0)
if latest != nil {
height = latest.Height
}

switch method {
case "eth_chainId":
return "0x" + fmt.Sprintf("%x", 9001), nil

case "net_version":
return "9001", nil

case "eth_blockNumber":
return "0x" + fmt.Sprintf("%x", height), nil

case "eth_getBlockByNumber":
block := e.dag.LatestBlock()
if block == nil {
return nil, nil
}
return map[string]interface{}{
"number":           "0x" + fmt.Sprintf("%x", block.Height),
"hash":             "0x" + block.Hash,
"parentHash":       "0x" + strings.Join(block.ParentHashes, ","),
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

case "eth_getBalance":
if len(params) > 0 {
var addr string
json.Unmarshal(params[0], &addr)
addr = strings.ToLower(addr)
if e.state == nil {
fmt.Println("[RPC] state is nil!")
return "0x0", nil
}
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
return "0x0", nil // Zero gas - gasless chain

case "eth_estimateGas":
return "0x0", nil // Zero gas - gasless chain

case "eth_getTransactionCount":
return "0x0", nil

case "eth_call":
return "0x", nil

case "eth_sendRawTransaction":
return "0x" + strings.Repeat("0", 64), nil

case "eth_getCode":
return "0x", nil

case "eth_getLogs":
return []interface{}{}, nil

case "web3_clientVersion":
return "Aequitas/v0.1.0/BlockDAG", nil

case "eth_syncing":
return false, nil

case "eth_accounts":
return []interface{}{}, nil

case "net_listening":
return true, nil

case "net_peerCount":
return "0x1", nil

default:
return nil, fmt.Errorf("method %s not supported", method)
}
}
