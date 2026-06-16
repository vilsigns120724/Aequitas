package keeper

import (
"encoding/json"
"fmt"
"io"
"net/http"
"strings"
"time"
)

type APIServer struct {
blockchain        *BlockDAG
p2pNode           *P2PNode
keeper            *Keeper
startTime         time.Time
proofServerStatus map[string]interface{}
state             *ChainState
}

func NewAPIServer(bc *BlockDAG, p2p *P2PNode, k *Keeper, state *ChainState) *APIServer {
s := &APIServer{
blockchain:        bc,
p2pNode:           p2p,
keeper:            k,
startTime:         time.Now(),
proofServerStatus: map[string]interface{}{},
state:             state,
}
go s.syncProofServerStatus()
return s
}

func (a *APIServer) syncProofServerStatus() {
for {
resp, err := http.Get("https://aequitas-proof-server-production.up.railway.app/health")
if err == nil {
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
var data map[string]interface{}
if json.Unmarshal(body, &data) == nil {
a.proofServerStatus = data
}
}
time.Sleep(30 * time.Second)
}
}

func (a *APIServer) Start(port int) {
mux := http.NewServeMux()
mux.HandleFunc("/", a.handleUI)
mux.HandleFunc("/api/status", a.handleStatus)
mux.HandleFunc("/api/blocks", a.handleBlocks)
mux.HandleFunc("/api/humans", a.handleHumans)
mux.HandleFunc("/api/sepolia/humans", a.handleSepoliaHumans)
mux.HandleFunc("/api/register", a.handleRegister)
mux.HandleFunc("/api/balance", a.handleBalance)
mux.HandleFunc("/registered", a.handleRegistered)
fmt.Println("── Starting EVM RPC ─────────────────────")
evmRPC := NewEVMRPCServer(a.blockchain, a.state)
mux.HandleFunc("/rpc", evmRPC.handleRPC)
if evmRPC.evm != nil {
fmt.Println("✓ EVM Engine ready")
} else {
fmt.Println("✗ EVM Engine failed")
}
addr := fmt.Sprintf(":%d", port)
fmt.Printf("✓ API Server listening on port %d\n", port)
go http.ListenAndServe(addr, mux)
}

func (a *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
latest := a.blockchain.LatestBlock()
uptime := int64(time.Since(a.startTime).Seconds())
humans := a.keeper.TotalHumans()
growth := humans * 10
if growth > 100 {
growth = 100
}
json.NewEncoder(w).Encode(map[string]interface{}{
"chain_id":     "aequitas-1",
"version":      "v0.3.0",
"height":       latest.Height,
"latest_hash":  latest.Hash,
"total_humans": a.state.TotalHumans(),
"total_supply": fmt.Sprintf("%.2f AEQ", a.state.TotalSupply()),
"node_id":      a.p2pNode.GetNodeID(),
"uptime":       uptime,
"block_time":   6,
"contract_v5":  "0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5",
"contract_v6":  "0xA76cA3bf34F2Ae5dFA0608696627e42b81180488",
		"contract_v7":  "0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78",
"bio_verifier": "0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2",
"chain_evm_id": 1926,
"index":        a.state.CalcAequitasIndex(),
"gini":         a.state.CalcGini(),
"growth":       growth,
"velocity":     50,
"phase":        a.state.CalcPhase(),
"fee_bps":      10,
})
}

func (a *APIServer) handleBlocks(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
blocks := a.blockchain.GetBlocks()
start := 0
if len(blocks) > 50 {
start = len(blocks) - 50
}
json.NewEncoder(w).Encode(blocks[start:])
}

func (a *APIServer) handleHumans(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
accounts := a.state.GetAllAccounts()
humans := []map[string]interface{}{}
for _, acc := range accounts {
if acc.IsHuman {
humans = append(humans, map[string]interface{}{
"address": acc.Address,
"balance": acc.Balance,
})
}
}
json.NewEncoder(w).Encode(map[string]interface{}{
"total":  len(humans),
"humans": humans,
})
}

func (a *APIServer) handleSepoliaHumans(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
resp, err := http.Get("https://aequitas-proof-server-production.up.railway.app/humans")
if err != nil {
json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
return
}
defer resp.Body.Close()
var data map[string]interface{}
json.NewDecoder(resp.Body).Decode(&data)
json.NewEncoder(w).Encode(data)
}

func (a *APIServer) handleBalance(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
wallet := strings.ToLower(r.URL.Query().Get("wallet"))
if wallet == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"balance": 0, "is_human": false})
return
}
balance := a.state.GetBalance(wallet)
isHuman := a.state.IsHuman(wallet)
json.NewEncoder(w).Encode(map[string]interface{}{
"wallet":   wallet,
"balance":  balance,
"is_human": isHuman,
})
}

func (a *APIServer) handleRegistered(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
wallet := r.URL.Query().Get("wallet")
fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Registered — Aequitas</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0A0E1A;color:#C9A84C;font-family:'Courier New',monospace;display:flex;align-items:center;justify-content:center;min-height:100vh;padding:20px;flex-direction:column;gap:20px;text-align:center}
.logo{font-size:2rem;font-weight:900;letter-spacing:8px;color:#C9A84C}
.box{background:#111827;border:1px solid #1E2D45;border-radius:12px;padding:32px;max-width:440px;width:100%}
.title{color:#22C55E;font-size:1.4rem;font-weight:bold;margin-bottom:8px}
.wallet{color:#6B7A99;font-size:0.7rem;margin-bottom:20px;word-break:break-all}
.divider{border-top:1px solid #1E2D45;margin:16px 0}
.sub{color:#6B7A99;font-size:0.82rem;line-height:1.9}
.hl{color:#C9A84C;font-weight:bold}
.btn{display:inline-block;margin-top:16px;padding:12px 24px;background:#C9A84C;color:#0A0E1A;border-radius:8px;text-decoration:none;font-weight:bold;font-size:0.8rem;letter-spacing:1px}
</style>
</head>
<body>
<div class="logo">AEQUITAS</div>
<div class="box">
<div class="title">🎉 Registered as Human!</div>
<div class="wallet">%s</div>
<div class="divider"></div>
<div class="sub">
<span class="hl">1,000 AEQ</span> has been credited to your wallet.<br><br>
Return to the <span class="hl">Aequitas App</span> — it will confirm your registration automatically.<br><br>
<span style="color:#4FC3F7">Money exists because people exist.</span>
</div>
<a class="btn" href="/">← VIEW EXPLORER</a>
</div>
</body>
</html>`, wallet)
}

func (a *APIServer) handleUI(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
fmt.Fprint(w, explorerHTML)
}
