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
fmt.Fprint(w, htmlPage)
}

const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>Aequitas — Proof of Humanity Chain</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{--bg:#050A14;--card:#0D1421;--card2:#111E2E;--border:#1A2D45;--green:#00E676;--blue:#4FC3F7;--gold:#FFB300;--purple:#CE93D8;--red:#EF5350;--text:#E8F4FD;--muted:#6B8CAE;--teal:#4DD0E1}
html,body{height:100%;overflow-x:hidden}
body{background:var(--bg);color:var(--text);font-family:'Courier New',monospace;min-height:100vh}

/* HEADER */
header{background:#080F1E;border-bottom:1px solid var(--border);padding:0 20px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:56px;gap:10px}
.logo-wrap{display:flex;align-items:center;gap:10px;flex-shrink:0}
.logo-icon{width:28px;height:28px;background:var(--gold);border-radius:6px;display:flex;align-items:center;justify-content:center;font-size:15px;flex-shrink:0}
.logo-text{font-size:1rem;font-weight:900;color:var(--gold);letter-spacing:4px}
.logo-sub{font-size:0.5rem;color:var(--muted);letter-spacing:2px}
.header-right{display:flex;gap:8px;align-items:center;flex-shrink:0}
.badge{display:flex;align-items:center;gap:4px;padding:4px 8px;border-radius:12px;font-size:0.6rem;letter-spacing:1px}
.badge-live{background:#00E67612;border:1px solid #00E67628;color:var(--green)}
.badge-dag{background:#4FC3F712;border:1px solid #4FC3F728;color:var(--blue)}
.pulse{width:5px;height:5px;border-radius:50%;background:var(--green);animation:pulse 2s infinite;flex-shrink:0}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:0.3}}

/* LANG */
.lang-select{background:#080F1E;color:var(--muted);border:1px solid var(--border);border-radius:5px;padding:4px 8px;cursor:pointer;font-family:monospace;font-size:0.62rem;outline:none}

/* TABS */
.tabs{background:#080F1E;border-bottom:1px solid var(--border);padding:0 20px;display:flex;overflow-x:auto;-webkit-overflow-scrolling:touch;scrollbar-width:none}
.tabs::-webkit-scrollbar{display:none}
.tab{padding:12px 14px;font-size:0.62rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:1px;text-transform:uppercase;white-space:nowrap;transition:all 0.2s;flex-shrink:0}
.tab:hover{color:var(--text)}
.tab.active{color:var(--blue);border-bottom-color:var(--blue)}
.tab-content{display:none}
.tab-content.active{display:block}

/* STATS GRID */
.hero{padding:16px 16px 0}
.section-label{font-size:0.55rem;color:var(--muted);letter-spacing:4px;text-transform:uppercase;margin-bottom:12px}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:1px;background:var(--border);border:1px solid var(--border);border-radius:10px;overflow:hidden;margin-bottom:16px}
.stat{background:var(--card);padding:16px 14px;position:relative;overflow:hidden}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:0.9rem;margin-bottom:6px}
.stat-lbl{font-size:0.55rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:5px}
.stat-val{font-size:1.6rem;font-weight:900;line-height:1;margin-bottom:3px}
.stat-sub{font-size:0.56rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:var(--green)}.c-green .stat-accent{background:var(--green)}
.c-blue .stat-val{color:var(--blue)}.c-blue .stat-accent{background:var(--blue)}
.c-gold .stat-val{color:var(--gold)}.c-gold .stat-accent{background:var(--gold)}
.c-purple .stat-val{color:var(--purple)}.c-purple .stat-accent{background:var(--purple)}
.c-teal .stat-val{color:var(--teal)}.c-teal .stat-accent{background:var(--teal)}

/* INFO BANNER */
.info-banner{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:16px;margin-bottom:16px;display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:16px}
.info-item-icon{font-size:1.2rem;margin-bottom:6px}
.info-item-title{font-size:0.68rem;color:var(--gold);font-weight:bold;margin-bottom:6px;letter-spacing:1px}
.info-item-text{font-size:0.63rem;color:var(--muted);line-height:1.8}

/* MAIN GRID */
.main-grid{display:grid;grid-template-columns:1fr 300px;gap:12px;padding:0 16px 16px}
@media(max-width:800px){.main-grid{grid-template-columns:1fr}.right-col{display:none}}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.sec-head{padding:11px 16px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:#080F1E}
.sec-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;display:flex;align-items:center;gap:6px}
.sec-dot{width:5px;height:5px;border-radius:50%;background:var(--green)}
.sec-count{font-size:0.58rem;color:var(--muted);background:var(--card2);padding:2px 7px;border-radius:8px;border:1px solid var(--border)}
.sec-desc{padding:9px 16px;font-size:0.62rem;color:var(--muted);background:#080F1E;border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:10px 16px;border-bottom:1px solid #0A1220;display:grid;grid-template-columns:56px 1fr auto;gap:8px;align-items:center;transition:background 0.15s}
.block-item:hover{background:#0D1421}
.block-item:last-child{border-bottom:none}
.block-num{font-size:0.78rem;font-weight:bold;color:var(--blue)}
.block-hash{font-size:0.63rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:4px;flex-wrap:wrap}
.block-parents{font-size:0.57rem;color:#3A5570}
.block-right{text-align:right}
.block-humans{font-size:0.65rem;color:var(--gold);margin-bottom:2px}
.block-time{font-size:0.57rem;color:var(--green)}
.badge-merge{background:#2D1B4E;color:var(--purple);font-size:0.53rem;padding:1px 4px;border-radius:3px;border:1px solid #4A2D7A}
.badge-tx{background:#0D2A1A;color:var(--green);font-size:0.53rem;padding:1px 4px;border-radius:3px;border:1px solid #1A4A2A}
.empty{padding:32px;text-align:center;color:var(--muted);font-size:0.68rem;line-height:2.2}
.right-col{display:flex;flex-direction:column;gap:10px}

/* CARDS */
.info-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:16px}
.ic-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:12px}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:7px 0;border-bottom:1px solid #0A1220}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.62rem;color:var(--muted)}
.ic-val{font-size:0.62rem;color:var(--text);text-align:right;max-width:58%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ic-val.green{color:var(--green)}.ic-val.blue{color:var(--blue)}.ic-val.gold{color:var(--gold)}.ic-val.purple{color:var(--purple)}
.mm-card{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:14px}
.mm-title{font-size:0.58rem;color:var(--blue);letter-spacing:2px;margin-bottom:10px;font-weight:bold}
.mm-row{display:flex;justify-content:space-between;padding:5px 0;border-bottom:1px solid #1A2D45}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.58rem;color:var(--muted)}
.mm-val{font-size:0.58rem;color:var(--purple)}
.mm-btn{width:100%;margin-top:10px;padding:9px;background:var(--blue);color:#050A14;border:none;border-radius:7px;cursor:pointer;font-family:monospace;font-size:0.65rem;font-weight:bold;letter-spacing:1px}
.phil-card{background:linear-gradient(135deg,#1A1200,#0D1421);border:1px solid #3A2800;border-radius:10px;padding:18px;text-align:center}
.phil-quote{font-size:0.78rem;color:var(--gold);font-style:italic;line-height:2;margin-bottom:5px}
.phil-sub{font-size:0.57rem;color:var(--muted);letter-spacing:2px}

/* HUMANS */
.humans-section{padding:16px;display:grid;grid-template-columns:1fr 280px;gap:12px}
@media(max-width:800px){.humans-section{grid-template-columns:1fr}}
.human-item{padding:11px 16px;border-bottom:1px solid #0A1220;display:flex;align-items:center;gap:10px}
.human-item:hover{background:#0D1421}
.human-item:last-child{border-bottom:none}
.human-avatar{width:34px;height:34px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.65rem;font-weight:bold;flex-shrink:0;border:2px solid}
.human-balance{font-size:0.78rem;color:var(--gold);font-weight:bold;margin-bottom:1px}
.human-addr{font-size:0.6rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.human-badge{font-size:0.55rem;padding:2px 7px;border-radius:8px;flex-shrink:0;background:#0D2A1A;color:var(--green);border:1px solid #1A4A2A}

/* INDEX */
.index-section{padding:16px;display:grid;grid-template-columns:1fr 1fr;gap:12px}
@media(max-width:700px){.index-section{grid-template-columns:1fr}}
.idx-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.idx-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:8px}
.idx-desc{font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:14px}
.idx-big{font-size:2.6rem;font-weight:900;color:var(--gold);line-height:1}
.idx-lbl{font-size:0.58rem;color:var(--muted);margin-top:3px}
.bar-bg{height:7px;background:#0D1421;border-radius:4px;overflow:hidden;margin:12px 0 5px}
.bar-fill{height:100%;border-radius:4px;background:linear-gradient(90deg,var(--green),var(--gold),var(--red));transition:width 1.5s}
.bar-labels{display:flex;justify-content:space-between;font-size:0.53rem;color:var(--muted)}
.metrics-row{display:grid;grid-template-columns:repeat(2,1fr);gap:7px;margin-top:12px}
.metric-box{background:#080F1E;border-radius:6px;padding:10px;text-align:center}
.metric-val{font-size:1.1rem;font-weight:bold;color:var(--gold)}
.metric-lbl{font-size:0.55rem;color:var(--muted);margin-top:2px}
.story-text{font-size:0.67rem;line-height:2;color:var(--muted)}
.story-text p{margin-bottom:12px}
.highlight-box{background:#080F1E;border-left:3px solid var(--gold);border-radius:0 8px 8px 0;padding:12px 16px;margin:14px 0;font-size:0.65rem;color:var(--text);line-height:1.8}

/* NETWORK */
.net-section{padding:16px;display:grid;grid-template-columns:1fr 1fr;gap:12px}
@media(max-width:700px){.net-section{grid-template-columns:1fr}}
.net-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px}
.net-title{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:12px}
.node-box{background:#080F1E;border-radius:7px;padding:12px;border:1px solid var(--border);margin-bottom:8px}
.node-status{display:flex;align-items:center;gap:5px;font-size:0.65rem;color:var(--green);margin-bottom:4px;font-weight:bold}
.node-dot{width:6px;height:6px;border-radius:50%;background:var(--green);box-shadow:0 0 5px var(--green)}
.node-url{font-size:0.57rem;color:var(--muted);word-break:break-all;margin-bottom:3px}
.node-desc{font-size:0.57rem;color:#3A5570}
.spec-table{width:100%;border-collapse:collapse}
.spec-table td{padding:7px 0;border-bottom:1px solid #0A1220;font-size:0.62rem}
.spec-table tr:last-child td{border-bottom:none}
.spec-table td:first-child{color:var(--muted);width:45%}
.spec-table td:last-child{text-align:right}
.bootstrap-box{background:#080F1E;border-radius:7px;padding:10px;font-size:0.58rem;color:var(--purple);word-break:break-all;line-height:1.7;border:1px solid var(--border)}

/* REGISTER */
.reg-section{padding:16px;max-width:600px;margin:0 auto}
.reg-hero{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:20px;margin-bottom:14px;text-align:center}
.reg-hero-title{font-size:0.95rem;font-weight:bold;color:var(--text);margin-bottom:7px}
.reg-hero-sub{font-size:0.65rem;color:var(--muted);line-height:1.8}
.app-only{background:#0D1220;border:1px solid #1A2040;border-radius:10px;padding:18px;text-align:center;margin-bottom:14px}
.app-only-icon{font-size:1.8rem;margin-bottom:7px}
.app-only-title{font-size:0.68rem;color:var(--purple);font-weight:bold;letter-spacing:2px;margin-bottom:8px}
.app-only-text{font-size:0.63rem;color:var(--muted);line-height:1.8}
.reg-steps{display:grid;grid-template-columns:repeat(4,1fr);gap:7px;margin-bottom:14px}
@media(max-width:520px){.reg-steps{grid-template-columns:repeat(2,1fr)}}
.reg-step{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:14px;text-align:center}
.step-num{width:26px;height:26px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 8px;font-weight:bold;font-size:0.7rem;color:#050A14}
.step-title{font-size:0.62rem;color:var(--text);font-weight:bold;margin-bottom:4px}
.step-desc{font-size:0.58rem;color:var(--muted);line-height:1.6}
.priv-bar{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:9px 12px;margin-bottom:12px;font-size:0.62rem;color:var(--green);text-align:center;line-height:1.7}
.reg-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px;margin-bottom:12px}
.wallet-box{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:9px;margin-bottom:9px;display:none}
.wallet-lbl{font-size:0.55rem;color:var(--muted);margin-bottom:2px;letter-spacing:1px}
.wallet-addr{font-size:0.7rem;color:var(--green);font-weight:bold}
.proof-box{background:var(--card2);border:1px solid #3A2800;border-radius:7px;padding:9px;margin-bottom:9px;display:none}
.proof-lbl{font-size:0.55rem;color:var(--gold);margin-bottom:2px;letter-spacing:1px}
.proof-val{font-size:0.62rem;color:var(--muted)}
.reg-btn{width:100%;padding:13px;border-radius:7px;border:none;cursor:pointer;font-family:monospace;font-size:0.72rem;font-weight:bold;letter-spacing:1px;transition:all 0.2s;margin-bottom:7px}
.btn-connect{background:var(--blue);color:#050A14}
.btn-connect:hover{opacity:0.87}
.btn-register{background:var(--gold);color:#050A14}
.btn-register:hover{opacity:0.87}
.reg-btn:disabled{opacity:0.3;cursor:not-allowed}
.reg-log{background:#080F1E;border-radius:7px;padding:10px;font-size:0.63rem;line-height:1.9;min-height:50px;border:1px solid var(--border)}
.reg-log .ok{color:var(--green)}.reg-log .err{color:var(--red)}.reg-log .info{color:var(--gold)}

/* PROTOCOL */
.proto-section{padding:16px;max-width:800px;margin:0 auto}

/* MOBILE */
@media(max-width:480px){
  .stats-grid{grid-template-columns:repeat(2,1fr)}
  .stat-val{font-size:1.3rem}
  header{height:50px}
  .logo-text{font-size:0.85rem;letter-spacing:2px}
  .badge-dag{display:none}
}
</style>
</head>
<body>
<header>
  <div class="logo-wrap">
    <div class="logo-icon">⚖</div>
    <div>
      <div class="logo-text">AEQUITAS</div>
      <div class="logo-sub">PROOF OF HUMANITY</div>
    </div>
  </div>
  <select class="lang-select" onchange="setLang(this.value)" id="lang-select">
    <option value="en">🌐 EN</option>
    <option value="de">🌐 DE</option>
    <option value="es">🌐 ES</option>
    <option value="ru">🌐 RU</option>
    <option value="zh">🌐 ZH</option>
    <option value="id">🌐 ID</option>
  </select>
  <div class="header-right">
    <div class="badge badge-live"><span class="pulse"></span>LIVE</div>
    <div class="badge badge-dag">● BLOCKDAG</div>
  </div>
</header>

<div class="tabs">
  <div class="tab active" onclick="showTab('register',this)">🔐 Register</div>
  <div class="tab" onclick="showTab('explorer',this)">🔍 Explorer</div>
  <div class="tab" onclick="showTab('humans',this)">👥 Humans</div>
  <div class="tab" onclick="showTab('index',this)">📊 Index</div>
  <div class="tab" onclick="showTab('network',this)">🌐 Network</div>
  <div class="tab" onclick="showTab('protocol',this)">📜 Protocol V6</div>
</div>

<!-- REGISTER -->
<div id="tab-register" class="tab-content active">
  <div class="reg-section">
    <div class="reg-hero">
      <div class="reg-hero-title">🔐 Register as a Verified Human</div>
      <div class="reg-hero-sub">Join the Aequitas network and receive 1,000 AEQ. One-time, permanent, gasless — cryptographically proving you are a unique human. No personal data stored. No gas fees. No waiting.</div>
    </div>

    <div class="app-only">
      <div class="app-only-icon">📱</div>
      <div class="app-only-title">REGISTRATION VIA ANDROID APP</div>
      <div class="app-only-text">Proof of Humanity requires biometric verification on your personal device. Your fingerprint is processed by the Hardware Secure Element — a dedicated chip that cannot be accessed remotely. The raw fingerprint data never leaves your phone. Download the Aequitas App, scan your fingerprint, connect your MetaMask wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically</strong>.</div>
    </div>

    <div class="reg-steps">
      <div class="reg-step">
        <div class="step-num">1</div>
        <div class="step-title">Biometric Scan</div>
        <div class="step-desc">Open app · scan fingerprint · processed by Hardware Secure Element · data never leaves device</div>
      </div>
      <div class="reg-step">
        <div class="step-num">2</div>
        <div class="step-title">ZKP Generation</div>
        <div class="step-desc">Groth16 proof generated · uniqueness verified · biometric hash never revealed</div>
      </div>
      <div class="reg-step">
        <div class="step-num">3</div>
        <div class="step-title">Connect Wallet</div>
        <div class="step-desc">App opens MetaMask · connect your wallet · address receives 1,000 AEQ permanently</div>
      </div>
      <div class="reg-step">
        <div class="step-num">4</div>
        <div class="step-title">1,000 AEQ</div>
        <div class="step-desc">Registered on Aequitas V6 · confirmed in next block · app notifies automatically</div>
      </div>
    </div>

    <div class="priv-bar">🔒 Hardware Secure Element · Groth16 ZKP · Data never leaves device · No gas fees · Permanent Sybil protection</div>

    <div class="reg-card">
      <div class="wallet-box" id="wallet-box">
        <div class="wallet-lbl">CONNECTED WALLET</div>
        <div class="wallet-addr" id="wallet-addr">—</div>
      </div>
      <div class="proof-box" id="proof-box">
        <div class="proof-lbl">⚡ ZK PROOF RECEIVED FROM APP</div>
        <div class="proof-val" id="proof-val">Connect wallet to register</div>
      </div>
      <button class="reg-btn btn-connect" id="btn-connect" onclick="connectWallet()">🦊 CONNECT METAMASK</button>
      <button class="reg-btn btn-register" id="btn-register" onclick="register()" disabled>🔐 REGISTER ON-CHAIN</button>
      <div class="reg-log" id="reg-status"><span class="info">// Open Aequitas Android App to generate your proof, then return here...</span></div>
    </div>

    <div class="info-card">
      <div class="ic-title">Registration Details</div>
      <div class="ic-row"><span class="ic-key">Network</span><span class="ic-val purple">Aequitas Chain (BlockDAG)</span></div>
      <div class="ic-row"><span class="ic-key">Chain ID</span><span class="ic-val blue">1926</span></div>
      <div class="ic-row"><span class="ic-key">Grant Amount</span><span class="ic-val gold">1,000 AEQ</span></div>
      <div class="ic-row"><span class="ic-key">Gas Fee</span><span class="ic-val green">FREE (gasless)</span></div>
      <div class="ic-row"><span class="ic-key">Registrations</span><span class="ic-val">Once per human · permanent · immutable</span></div>
      <div class="ic-row"><span class="ic-key">Biometric Data</span><span class="ic-val green">Never stored anywhere</span></div>
      <div class="ic-row"><span class="ic-key">Confirmation</span><span class="ic-val">Within 6 seconds (next block)</span></div>
      <div class="ic-row"><span class="ic-key">Contract V6</span><span class="ic-val" style="font-size:0.55rem">0xA76cA3...80488</span></div>
    </div>
  </div>
</div>

<!-- EXPLORER -->
<div id="tab-explorer" class="tab-content">
  <div class="hero">
    <div class="section-label">Live Chain Statistics</div>
    <div class="stats-grid">
      <div class="stat c-blue"><div class="stat-accent"></div><div class="stat-icon">🔗</div><div class="stat-lbl">Block Height</div><div class="stat-val" id="s-height">—</div><div class="stat-sub">New block every 6s · BlockDAG · Two nodes parallel</div></div>
      <div class="stat c-green"><div class="stat-accent"></div><div class="stat-icon">🧬</div><div class="stat-lbl">Verified Humans</div><div class="stat-val" id="s-humans">—</div><div class="stat-sub">Biometric ZKP · One person, one wallet, forever</div></div>
      <div class="stat c-gold"><div class="stat-accent"></div><div class="stat-icon">🪙</div><div class="stat-lbl">Total Supply</div><div class="stat-val" id="s-supply">—</div><div class="stat-sub">Always = Humans × 1,000 AEQ</div></div>
      <div class="stat c-purple"><div class="stat-accent"></div><div class="stat-icon">⚖</div><div class="stat-lbl">Aequitas Index</div><div class="stat-val" id="s-index">—</div><div class="stat-sub">0 = perfect equality · 100 = max inequality</div></div>
      <div class="stat c-teal"><div class="stat-accent"></div><div class="stat-icon">⚡</div><div class="stat-lbl">Uptime</div><div class="stat-val" id="s-uptime" style="font-size:1rem">—</div><div class="stat-sub">Node v0.3.0 · Railway + Render · PostgreSQL</div></div>
    </div>
    <div class="info-banner">
      <div><div class="info-item-icon">🧬</div><div class="info-item-title">Proof of Humanity</div><div class="info-item-text">Every AEQ holder must prove they are a unique living human through biometric verification. No bots, no corporations, no AI systems can hold AEQ. Only real humans. Your fingerprint is processed by the Hardware Secure Element — the same chip securing your banking apps. Biometric data never leaves your device.</div></div>
      <div><div class="info-item-icon">⚖</div><div class="info-item-title">Radically Fair Distribution</div><div class="info-item-text">Every verified human receives exactly 1,000 AEQ. The first person and the billionth person receive identical amounts. No pre-mine, no founder allocation, no investor round. Total supply always equals verified humans × 1,000. The most egalitarian monetary distribution system ever designed.</div></div>
      <div><div class="info-item-icon">🔗</div><div class="info-item-title">BlockDAG Architecture</div><div class="info-item-text">Multiple blocks can be produced simultaneously by different nodes and merged. Higher throughput, lower latency, better fault tolerance. Merge events are marked 🔀 in the explorer. Built from scratch in Go — inspired by IOTA and Phantom but implemented independently.</div></div>
      <div><div class="info-item-icon">⛽</div><div class="info-item-title">Truly Gasless</div><div class="info-item-text">Registration costs absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account. If you are a human with a smartphone, you can register. Transaction fees are covered by the protocol itself — making Aequitas accessible to every person on Earth.</div></div>
    </div>
  </div>
  <div class="main-grid">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span>Recent Blocks</div><div class="sec-count" id="block-count">—</div></div>
      <div class="sec-desc">🔀 MERGE = multiple parents (BlockDAG). ✅ TX = registration transaction. Block time: ~6 seconds.</div>
      <div id="blocks-list"><div class="empty">Loading blocks...</div></div>
    </div>
    <div class="right-col">
      <div class="info-card">
        <div class="ic-title">Network Info</div>
        <div class="ic-row"><span class="ic-key">Chain Name</span><span class="ic-val gold">Aequitas Chain</span></div>
        <div class="ic-row"><span class="ic-key">Chain ID</span><span class="ic-val blue">1926</span></div>
        <div class="ic-row"><span class="ic-key">Symbol</span><span class="ic-val gold">AEQ</span></div>
        <div class="ic-row"><span class="ic-key">Block Time</span><span class="ic-val">6 seconds</span></div>
        <div class="ic-row"><span class="ic-key">Consensus</span><span class="ic-val purple">BlockDAG + PoH</span></div>
        <div class="ic-row"><span class="ic-key">Active Nodes</span><span class="ic-val green">2 Online</span></div>
        <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val">Groth16</span></div>
        <div class="ic-row"><span class="ic-key">Storage</span><span class="ic-val green">PostgreSQL</span></div>
      </div>
      <div class="mm-card">
        <div class="mm-title">🦊 ADD TO METAMASK</div>
        <div class="mm-row"><span class="mm-key">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
        <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.5rem">...9fba.up.railway.app/rpc</span></div>
        <div class="mm-row"><span class="mm-key">Chain ID</span><span class="mm-val">1926</span></div>
        <div class="mm-row"><span class="mm-key">Symbol</span><span class="mm-val">AEQ</span></div>
        <div class="mm-row"><span class="mm-key">Decimals</span><span class="mm-val">18</span></div>
        <button class="mm-btn" onclick="addToMetaMask()">+ ADD AEQUITAS NETWORK</button>
      </div>
      <div class="phil-card">
        <div class="phil-quote">"Money exists because people exist.<br>Nothing more, nothing less."</div>
        <div class="phil-sub">— THE AEQUITAS PRINCIPLE —</div>
      </div>
    </div>
  </div>
</div>

<!-- HUMANS -->
<div id="tab-humans" class="tab-content">
  <div class="hero">
    <div class="section-label">Verified Humans on Aequitas Chain</div>
    <div class="info-banner">
      <div><div class="info-item-icon">🔒</div><div class="info-item-title">What is a Verified Human?</div><div class="info-item-text">A Verified Human on Aequitas is a wallet address cryptographically proven to belong to a unique living human. Verification uses biometric data processed through the Hardware Secure Element of an Android smartphone. The data is never transmitted or stored. Only a Zero-Knowledge Proof derived from it is used. Once verified, the wallet is permanently linked to that person's biometric identity.</div></div>
      <div><div class="info-item-icon">🧮</div><div class="info-item-title">Zero-Knowledge Proof System</div><div class="info-item-text">Aequitas uses the Groth16 proving system — also used by Zcash — over the BN128 elliptic curve. Your fingerprint is hashed into a field element. This hash produces a mathematical proof guaranteeing "a unique biometric hash was used" without revealing what the hash is. Proof size: ~200 bytes. Verification time: ~10ms.</div></div>
      <div><div class="info-item-icon">🛡</div><div class="info-item-title">Sybil Attack Prevention</div><div class="info-item-text">Each biometric hash is stored permanently. Attempting to register twice with the same fingerprint is immediately rejected. One human, one wallet, forever — guaranteed by cryptography, not trust. This makes Aequitas the first cryptocurrency that is mathematically immune to Sybil attacks.</div></div>
      <div><div class="info-item-icon">🌍</div><div class="info-item-title">Global Inclusion</div><div class="info-item-text">No bank account, no credit card, no existing cryptocurrency required. Just an Android smartphone with a fingerprint sensor — a device over 3 billion people already own. Registration is free, takes under 2 minutes, and grants 1,000 AEQ instantly. Financial inclusion at a scale never before achieved.</div></div>
    </div>
  </div>
  <div class="humans-section">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span>Registered Humans</div><div class="sec-count" id="human-count-badge">0</div></div>
      <div class="sec-desc">Every address verified as unique human through biometric ZKP. Each received 1,000 AEQ. Permanent, immutable, stored in PostgreSQL and on-chain.</div>
      <div id="humans-list"><div class="empty">No humans registered yet.<br><br>Download the Aequitas Android App and be the first human on the chain!</div></div>
    </div>
    <div class="right-col">
      <div class="info-card">
        <div class="ic-title">Registry Stats</div>
        <div class="ic-row"><span class="ic-key">Total Humans</span><span class="ic-val green" id="stat-humans">0</span></div>
        <div class="ic-row"><span class="ic-key">Total Supply</span><span class="ic-val gold" id="stat-supply">0 AEQ</span></div>
        <div class="ic-row"><span class="ic-key">Grant per Human</span><span class="ic-val gold">1,000 AEQ</span></div>
        <div class="ic-row"><span class="ic-key">Registration Fee</span><span class="ic-val green">FREE</span></div>
        <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val">Groth16 / BN128</span></div>
        <div class="ic-row"><span class="ic-key">Biometric Storage</span><span class="ic-val green">Never stored</span></div>
      </div>
      <div class="info-card">
        <div class="ic-title">ZKP Technical Details</div>
        <div style="font-size:0.62rem;color:var(--muted);line-height:1.8">Groth16 over BN128 elliptic curve. Proof size: ~200 bytes. Verification time: ~10ms. Circuit compiled with snarkjs and circom. Trusted setup: Hermez ceremony parameters. Proof generation on the Aequitas Proof Server (Node.js).</div>
      </div>
    </div>
  </div>
</div>

<!-- INDEX -->
<div id="tab-index" class="tab-content">
  <div class="index-section">
    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title">Aequitas Index — Real-Time Economic Equality Score</div>
      <div class="idx-desc">Calculated from the on-chain balance distribution of all verified humans. 0 = perfect equality (everyone has identical AEQ). 100 = maximum inequality (one person controls everything). The protocol uses this index to automatically trigger redistribution when inequality grows beyond safe thresholds.</div>
      <div style="display:grid;grid-template-columns:auto 1fr;gap:20px;align-items:center;margin-top:12px">
        <div><div class="idx-big" id="idx-score">—</div><div class="idx-lbl">Current Index</div></div>
        <div>
          <div class="bar-bg"><div class="bar-fill" id="idx-bar" style="width:0%"></div></div>
          <div class="bar-labels"><span>0 — Perfect Equality</span><span>50</span><span>100 — Max Inequality</span></div>
          <div style="margin-top:8px;font-size:0.63rem;color:var(--muted);background:#080F1E;padding:8px;border-radius:6px" id="idx-phase-desc">—</div>
        </div>
      </div>
      <div class="metrics-row" style="grid-template-columns:repeat(4,1fr)">
        <div class="metric-box"><div class="metric-val" id="idx-gini">—</div><div class="metric-lbl">Gini Coefficient</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-supply2">—</div><div class="metric-lbl">Total Supply</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-phase">—</div><div class="metric-lbl">Protocol Phase</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-humans2">—</div><div class="metric-lbl">Verified Humans</div></div>
      </div>
    </div>

    <div class="idx-card">
      <div class="idx-title">Redistribution Pools</div>
      <div class="idx-desc">When inequality thresholds are exceeded, AEQ is automatically redirected into these pools. Controlled entirely by protocol logic — no human has access.</div>
      <div class="metrics-row">
        <div class="metric-box"><div class="metric-val" id="pool-v">—</div><div class="metric-lbl">Velocity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-l">—</div><div class="metric-lbl">Liquidity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-u">—</div><div class="metric-lbl">UBI Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-t">—</div><div class="metric-lbl">Treasury</div></div>
      </div>
    </div>

    <div class="idx-card">
      <div class="idx-title">Protocol Phases</div>
      <div class="idx-desc">Phase transitions happen automatically based on verified humans and Gini — no governance vote required.</div>
      <table class="spec-table">
        <tr><td>Phase 0</td><td style="color:var(--green)">Bootstrap · &lt;100 humans · Cap: 50x fairShare</td></tr>
        <tr><td>Phase 1</td><td style="color:var(--blue)">Growth · 100–10,000 humans · Cap: 20x</td></tr>
        <tr><td>Phase 2</td><td style="color:var(--gold)">Stability · 10k–1M humans · Cap: 10x</td></tr>
        <tr><td>Phase 3</td><td style="color:var(--purple)">Maturity · 1M+ humans · Cap: 3x</td></tr>
      </table>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title">The Story of Aequitas — Why This Exists</div>
      <div class="story-text">
        <p>The year is 2009. Satoshi Nakamoto releases Bitcoin. For the first time, value can transfer between any two people without a bank. A genuine revolution. But something goes wrong almost immediately.</p>
        <p>Early miners accumulate millions of coins at almost zero cost. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. Bitcoin's estimated Gini exceeds 0.85 — higher than any country on Earth. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.</p>
        <p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was created to answer: <em style="color:var(--gold)">"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"</em></p>
        <p>The answer: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.</strong></p>
        <p>The technology now exists: Zero-Knowledge Proofs verify unique humanity without revealing personal information. Blockchain stores verifications permanently. Smartphones provide the biometric sensors. Aequitas assembles these into a coherent system for the first time.</p>
        <p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em> This is not a slogan — it is the mathematical foundation of the entire system. We invite every human being to join us.</p>
      </div>
    </div>
  </div>
</div>

<!-- NETWORK -->
<div id="tab-network" class="tab-content">
  <div class="net-section">
    <div class="net-card" style="grid-column:1/-1">
      <div class="net-title">Active Nodes — Current Network Topology</div>
      <div style="font-size:0.65rem;color:var(--muted);line-height:1.8;margin-bottom:12px">The Aequitas network operates on two nodes in geographically distributed cloud environments. Both participate in block production, state synchronization, and API serving. They communicate via libp2p (also used by IPFS and Ethereum 2.0) and sync blocks via HTTP. Both share the same PostgreSQL database.</div>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px">
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span>Node 1 — Railway (Primary)</div>
          <div class="node-url">aequitas-production-9fba.up.railway.app</div>
          <div class="node-desc">Primary API · Block producer · P2P bootstrap · PostgreSQL · RPC for MetaMask</div>
        </div>
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span>Node 2 — Render (Secondary)</div>
          <div class="node-url">aequitas-node-2.onrender.com</div>
          <div class="node-desc">Secondary API · Block producer · P2P peer · HTTP sync · Shared PostgreSQL</div>
        </div>
      </div>
    </div>
    <div class="net-card">
      <div class="net-title">Bootstrap Node</div>
      <div style="font-size:0.63rem;color:var(--muted);line-height:1.8;margin-bottom:10px">To run your own Aequitas node, connect to the bootstrap node using the libp2p multiaddress below. Your node will automatically discover peers and begin participating in consensus.</div>
      <div class="bootstrap-box">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
    </div>
    <div class="net-card">
      <div class="net-title">Technical Specifications</div>
      <table class="spec-table">
        <tr><td>Chain ID</td><td>1926 (0x786)</td></tr>
        <tr><td>EVM</td><td style="color:var(--green)">Yes — JSON-RPC /rpc · MetaMask</td></tr>
        <tr><td>Block Time</td><td>~6 seconds</td></tr>
        <tr><td>Consensus</td><td style="color:var(--purple)">BlockDAG + Proof of Humanity</td></tr>
        <tr><td>P2P Protocol</td><td>libp2p (Go)</td></tr>
        <tr><td>ZKP System</td><td>Groth16 / snarkjs / circom</td></tr>
        <tr><td>Curve</td><td>BN128 (alt-bn128)</td></tr>
        <tr><td>Storage</td><td style="color:var(--green)">PostgreSQL (persistent)</td></tr>
        <tr><td>Language</td><td>Go 1.24</td></tr>
        <tr><td>Source Code</td><td><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub ↗</a></td></tr>
      </table>
    </div>
    <div class="net-card">
      <div class="net-title">MetaMask Configuration</div>
      <table class="spec-table">
        <tr><td>Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
        <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.55rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
        <tr><td>Chain ID</td><td style="color:var(--gold)">1926</td></tr>
        <tr><td>Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
        <tr><td>Decimals</td><td>18</td></tr>
        <tr><td>Block Explorer</td><td style="font-size:0.55rem">aequitas-production-9fba.up.railway.app</td></tr>
      </table>
      <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:12px">+ ADD TO METAMASK</button>
    </div>
  </div>
</div>

<!-- PROTOCOL V6 -->
<div id="tab-protocol" class="tab-content">
  <div class="proto-section">
    <div class="section-label">Aequitas V6 Protocol — Technical Documentation</div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">V6 Contract Addresses</div>
      <div class="highlight-box">
        Chain: Aequitas Chain (Chain ID: 1926 · 0x786)<br>
        RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>
        BioVerifier (Groth16): 0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2<br>
        AequitasV6 (Main):     0xA76cA3bf34F2Ae5dFA0608696627e42b81180488<br>
        V5 (Sepolia legacy):   0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">1. PROOF OF ALIVE</div>
      <div class="story-text">
        <p>What happens to money when people die or disappear? In Bitcoin, millions of BTC are permanently lost. In Aequitas, money represents people — if someone disappears, their AEQ eventually returns to the community.</p>
        <div class="highlight-box">
          Year 0-2: Normal usage<br>
          Year 2: Warning 1 — Guardian can respond<br>
          Year 2+60d: Warning 2 — Guardian can respond<br>
          Year 2+120d: Warning 3 — Guardian can respond<br>
          Year 2+180d: AEQ → PERSONAL ESCROW (not UBI yet)<br>
          Year 4: If still inactive → UBI Pool → distributed equally
        </div>
        <p>Why Escrow first? Someone imprisoned for 3 years gets their Escrow back when they return, PLUS the current fairShare. Not punished for being imprisoned.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">2. GUARDIAN SYSTEM</div>
      <div class="story-text">
        <p>What if someone cannot access their device for months? In Bitcoin, their funds are frozen forever. In Aequitas, a trusted Guardian can confirm they are still alive.</p>
        <div class="highlight-box">
          • 1 Guardian per human (another verified human)<br>
          • Guardian can ONLY call confirmAlive() — zero transaction rights<br>
          • Guardian CANNOT move funds or transfer AEQ<br>
          • Max 3 wards per Guardian<br>
          • 7-day timelock on assignment (prevents forced assignment under duress)<br>
          • No circular relationships (A guards B → B cannot guard A)<br>
          • Guardian cannot have their own Guardian
        </div>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">3. DEMURRAGE — Anti-Hoarding</div>
      <div class="story-text">
        <p>1% annual fee on any balance ABOVE your fairShare. The money goes to the UBI Pool, not deleted.</p>
        <div class="highlight-box">
          Example: fairShare = 1,000 AEQ · Your balance = 3,000 AEQ<br>
          Excess: 2,000 AEQ<br>
          Monthly fee: 2,000 × 1% ÷ 12 = 1.67 AEQ → UBI Pool<br>
          Each other human gains +0.17 AEQ/month
        </div>
        <p>Historical precedent: Wörgl, Austria (1932) — demurrage currency reduced unemployment 25% in one year while the rest of Austria suffered. The Central Bank shut it down because it worked too well.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">4. WEALTH CAP</div>
      <div class="story-text">
        <p>Hard ceiling on how much AEQ any single human can hold. Excess is instantly redistributed equally to ALL active humans.</p>
        <div class="highlight-box">
          Phase 0 (1-100 humans):    50x fairShare<br>
          Phase 1 (101-1,000):       20x fairShare<br>
          Phase 2 (1,001-10,000):    10x fairShare<br>
          Phase 3 (10,001-100,000):   5x fairShare<br>
          Phase 4 (100,000+):         3x fairShare
        </div>
        <p>ALWAYS active from human #1. Previous versions had no cap in Phase 0 — a mistake. V6 fixes this. Bitcoin's top 1% controls 90%+ of supply. In Aequitas, mathematical law makes that impossible.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">5. UNIVERSAL BASIC INCOME</div>
      <div class="story-text">
        <p>UBI from protocol economics — not taxation. Requires no government, no political decision.</p>
        <div class="highlight-box">
          Sources of UBI Pool:<br>
          1. Transaction fees (0.1%) → 20% to UBI Pool<br>
          2. Wealth cap overflow → redistributed instantly to all<br>
          3. Demurrage (1% annual on excess) → UBI Pool<br>
          4. Inactive wallet escrow (after 4 years) → UBI Pool<br><br>
          Monthly: UBI Pool ÷ active humans = equal payment to all
        </div>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:12px">
      <div class="idx-title">6. NO ALGORITHMIC INFLATION</div>
      <div class="story-text">
        <div class="highlight-box">
          The ONLY event that creates new AEQ:<br>
          A new verified human registers → 1,000 AEQ created<br><br>
          Total AEQ = Verified Active Humans × 1,000<br>
          (Always true. Always verifiable. Cannot be changed.)
        </div>
        <p>Previous versions had algorithmic inflation that could be manipulated. V6 makes manipulation impossible: only human biometric registration creates new money.</p>
      </div>
    </div>
  </div>
</div>

<script>
const PROOF_SERVER = 'https://aequitas-proof-server-production.up.railway.app';
const CHAIN_ID_HEX = '0x786'; // 1926
let walletAddr = '', proofParams = null, currentLang = 'en';

function showTab(name, el) {
  document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-' + name).classList.add('active');
  el.classList.add('active');
}


const T = {
en:{
  'tab-register':'🔐 Register','tab-explorer':'🔍 Explorer','tab-humans':'👥 Humans','tab-index':'📊 Index','tab-network':'🌐 Network','tab-protocol':'📜 Protocol V6',
  'reg-title':'🔐 Register as a Verified Human',
  'reg-sub':'Join the Aequitas network and receive 1,000 AEQ. One-time, permanent, gasless. No personal data stored.',
  'app-title':'REGISTRATION VIA ANDROID APP',
  'app-text':'Proof of Humanity requires biometric verification on your personal device. Your fingerprint is processed by the Hardware Secure Element — raw data never leaves your phone. Download the app, scan your fingerprint, connect your wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically</strong>.',
  's1t':'Biometric Scan','s1d':'Open app · scan fingerprint · HSE processes · data never leaves device',
  's2t':'ZKP Generation','s2d':'Groth16 proof generated · uniqueness verified · hash never revealed',
  's3t':'Connect Wallet','s3d':'App opens MetaMask · connect wallet · address receives 1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'Registered on V6 · confirmed in next block · app notifies automatically',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Data never leaves device · No gas fees · Permanent Sybil protection',
  'conn-wallet':'CONNECTED WALLET','btn-conn':'🦊 CONNECT METAMASK','btn-reg':'🔐 REGISTER ON-CHAIN',
  'btn-add-mm':'+ ADD AEQUITAS NETWORK','reg-log-hint':'// Open Aequitas Android App to generate your proof, then return here...',
  'phil':'"Money exists because people exist.<br>Nothing more, nothing less."','phil-sub':'— THE AEQUITAS PRINCIPLE —',
  's-height':'Block Height','s-height-sub':'New block every 6s · BlockDAG · Two nodes parallel',
  's-humans':'Verified Humans','s-humans-sub':'Biometric ZKP · One person, one wallet, forever',
  's-supply':'Total Supply','s-supply-sub':'Always = Humans × 1,000 AEQ',
  's-index':'Aequitas Index','s-index-sub':'0 = perfect equality · 100 = max inequality',
  's-uptime':'Uptime','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Proof of Humanity','ib-poh-t':'Every AEQ holder must prove they are a unique living human. No bots, no corporations, no AI. Only real humans. Biometric data never leaves your device.',
  'ib-fair':'Radically Fair','ib-fair-t':'Every verified human receives exactly 1,000 AEQ. No pre-mine, no founder allocation. Total supply always equals verified humans × 1,000.',
  'ib-dag':'BlockDAG Architecture','ib-dag-t':'Multiple blocks can be produced simultaneously and merged. Higher throughput, lower latency, better fault tolerance.',
  'ib-gas':'Truly Gasless','ib-gas-t':'Registration costs absolutely nothing. No ETH, BNB, or MATIC required. No credit card, no bank account.'
},
de:{
  'tab-register':'🔐 Registrieren','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-protocol':'📜 Protokoll V6',
  'reg-title':'🔐 Als verifizierter Mensch registrieren',
  'reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte 1.000 AEQ. Einmalig, permanent, gebührenfrei. Keine persönlichen Daten.',
  'app-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP',
  'app-text':'Der Menschlichkeitsnachweis erfordert biometrische Verifizierung auf deinem Gerät. Dein Fingerabdruck wird durch das Hardware Secure Element verarbeitet — rohe Daten verlassen niemals dein Telefon. Lade die App herunter, scanne deinen Fingerabdruck, verbinde deine Wallet, und deine <strong style="color:var(--gold)">1.000 AEQ werden automatisch gewährt</strong>.',
  's1t':'Biometrischer Scan','s1d':'App öffnen · Fingerabdruck scannen · HSE verarbeitet · Daten verlassen nie das Gerät',
  's2t':'ZKP-Erzeugung','s2d':'Groth16-Beweis generiert · Einzigartigkeit verifiziert · Hash nie enthüllt',
  's3t':'Wallet verbinden','s3d':'App öffnet MetaMask · Wallet verbinden · Adresse erhält 1.000 AEQ',
  's4t':'1.000 AEQ','s4d':'Auf V6 registriert · im nächsten Block bestätigt · App benachrichtigt automatisch',
  'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Daten verlassen nie das Gerät · Keine Gasgebühren · Permanenter Sybil-Schutz',
  'conn-wallet':'VERBUNDENE WALLET','btn-conn':'🦊 METAMASK VERBINDEN','btn-reg':'🔐 ON-CHAIN REGISTRIEREN',
  'btn-add-mm':'+ AEQUITAS-NETZWERK HINZUFÜGEN','reg-log-hint':'// Öffne die Aequitas Android-App um deinen Beweis zu generieren, dann kehre hierher zurück...',
  'phil':'"Geld existiert weil Menschen existieren.<br>Nichts mehr, nichts weniger."','phil-sub':'— DAS AEQUITAS-PRINZIP —',
  's-height':'Blockhöhe','s-height-sub':'Neuer Block alle 6 Sek · BlockDAG · Zwei Nodes parallel',
  's-humans':'Verifizierte Menschen','s-humans-sub':'Biometrischer ZKP · Eine Person, eine Wallet, für immer',
  's-supply':'Gesamtmenge','s-supply-sub':'Immer = Menschen × 1.000 AEQ',
  's-index':'Aequitas-Index','s-index-sub':'0 = vollkommene Gleichheit · 100 = maximale Ungleichheit',
  's-uptime':'Betriebszeit','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Menschlichkeitsnachweis','ib-poh-t':'Jeder AEQ-Inhaber muss beweisen dass er ein einzigartiger lebender Mensch ist. Keine Bots, keine Unternehmen, keine KI. Nur echte Menschen.',
  'ib-fair':'Radikal faire Verteilung','ib-fair-t':'Jeder verifizierte Mensch erhält genau 1.000 AEQ. Keine Vorzuteilung, keine Gründeranteile. Gesamtmenge immer = Menschen × 1.000.',
  'ib-dag':'BlockDAG-Architektur','ib-dag-t':'Mehrere Blöcke können gleichzeitig produziert und zusammengeführt werden. Höherer Durchsatz, niedrigere Latenz.',
  'ib-gas':'Wirklich gebührenfrei','ib-gas-t':'Registrierung kostet absolut nichts. Kein ETH, BNB oder MATIC. Kein Bankkonto erforderlich.'
},
es:{
  'tab-register':'🔐 Registrar','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-protocol':'📜 Protocolo V6',
  'reg-title':'🔐 Regístrate como Humano Verificado',
  'reg-sub':'Únete a la red Aequitas y recibe 1,000 AEQ. Único, permanente, sin gas. Sin datos personales.',
  'app-title':'REGISTRO SOLO VÍA APP ANDROID',
  'app-text':'La Prueba de Humanidad requiere verificación biométrica en tu dispositivo. Tu huella se procesa por el Elemento Seguro de Hardware — los datos brutos nunca salen de tu teléfono. Descarga la app, escanea tu huella, conecta tu wallet, y tus <strong style="color:var(--gold)">1,000 AEQ serán otorgados automáticamente</strong>.',
  's1t':'Escaneo Biométrico','s1d':'Abrir app · escanear huella · HSE procesa · datos nunca salen del dispositivo',
  's2t':'Generación ZKP','s2d':'Prueba Groth16 generada · unicidad verificada · hash nunca revelado',
  's3t':'Conectar Wallet','s3d':'App abre MetaMask · conectar wallet · dirección recibe 1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'Registrado en V6 · confirmado en próximo bloque · app notifica automáticamente',
  'priv-bar':'🔒 Elemento Seguro de Hardware · ZKP Groth16 · Datos nunca salen del dispositivo · Sin tarifas de gas',
  'conn-wallet':'WALLET CONECTADA','btn-conn':'🦊 CONECTAR METAMASK','btn-reg':'🔐 REGISTRAR ON-CHAIN',
  'btn-add-mm':'+ AGREGAR RED AEQUITAS','reg-log-hint':'// Abre la App Android Aequitas para generar tu prueba, luego regresa aquí...',
  'phil':'"El dinero existe porque las personas existen.<br>Nada más, nada menos."','phil-sub':'— EL PRINCIPIO AEQUITAS —',
  's-height':'Altura de Bloque','s-height-sub':'Nuevo bloque cada 6s · BlockDAG · Dos nodos paralelos',
  's-humans':'Humanos Verificados','s-humans-sub':'ZKP biométrico · Una persona, una wallet, siempre',
  's-supply':'Suministro Total','s-supply-sub':'Siempre = Humanos × 1,000 AEQ',
  's-index':'Índice Aequitas','s-index-sub':'0 = igualdad perfecta · 100 = desigualdad máxima',
  's-uptime':'Tiempo Activo','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Prueba de Humanidad','ib-poh-t':'Cada titular de AEQ debe probar que es un humano único vivo. Sin bots, sin corporaciones, sin IA. Solo humanos reales.',
  'ib-fair':'Distribución Radicalmente Justa','ib-fair-t':'Cada humano verificado recibe exactamente 1,000 AEQ. Sin pre-minado, sin asignación a fundadores.',
  'ib-dag':'Arquitectura BlockDAG','ib-dag-t':'Múltiples bloques pueden producirse simultáneamente y fusionarse. Mayor rendimiento, menor latencia.',
  'ib-gas':'Verdaderamente Sin Gas','ib-gas-t':'El registro no cuesta nada. No se necesita ETH, BNB ni MATIC. Sin cuenta bancaria.'
},
ru:{
  'tab-register':'🔐 Регистрация','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-protocol':'📜 Протокол V6',
  'reg-title':'🔐 Зарегистрируйтесь как Верифицированный Человек',
  'reg-sub':'Присоединитесь к сети Aequitas и получите 1 000 AEQ. Одноразово, постоянно, бесплатно. Никаких личных данных.',
  'app-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ',
  'app-text':'Доказательство человечности требует биометрической верификации на вашем устройстве. Ваш отпечаток обрабатывается Hardware Secure Element — сырые данные никогда не покидают ваш телефон. Скачайте приложение, отсканируйте отпечаток, подключите кошелёк, и ваши <strong style="color:var(--gold)">1 000 AEQ будут начислены автоматически</strong>.',
  's1t':'Биометрический Скан','s1d':'Открыть приложение · сканировать отпечаток · HSE обрабатывает · данные не покидают устройство',
  's2t':'Генерация ZKP','s2d':'Сгенерировано доказательство Groth16 · уникальность верифицирована · хэш не раскрывается',
  's3t':'Подключить Кошелёк','s3d':'Приложение открывает MetaMask · подключить кошелёк · адрес получает 1 000 AEQ',
  's4t':'1 000 AEQ','s4d':'Зарегистрировано на V6 · подтверждено в следующем блоке · приложение уведомляет автоматически',
  'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Данные не покидают устройство · Без комиссий',
  'conn-wallet':'ПОДКЛЮЧЁННЫЙ КОШЕЛЁК','btn-conn':'🦊 ПОДКЛЮЧИТЬ METAMASK','btn-reg':'🔐 ЗАРЕГИСТРИРОВАТЬСЯ ON-CHAIN',
  'btn-add-mm':'+ ДОБАВИТЬ СЕТЬ AEQUITAS','reg-log-hint':'// Откройте приложение Aequitas для генерации доказательства, затем вернитесь...',
  'phil':'"Деньги существуют потому что существуют люди.<br>Ничего больше, ничего меньше."','phil-sub':'— ПРИНЦИП AEQUITAS —',
  's-height':'Высота Блока','s-height-sub':'Новый блок каждые 6 сек · BlockDAG · Два узла параллельно',
  's-humans':'Верифицированных Людей','s-humans-sub':'Биометрический ZKP · Один человек, один кошелёк, навсегда',
  's-supply':'Общее Предложение','s-supply-sub':'Всегда = Люди × 1 000 AEQ',
  's-index':'Индекс Aequitas','s-index-sub':'0 = полное равенство · 100 = максимальное неравенство',
  's-uptime':'Время Работы','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Доказательство Человечности','ib-poh-t':'Каждый владелец AEQ должен доказать что он уникальный живой человек. Никаких ботов, корпораций, ИИ. Только настоящие люди.',
  'ib-fair':'Радикально Справедливое Распределение','ib-fair-t':'Каждый верифицированный человек получает ровно 1 000 AEQ. Общее предложение всегда = люди × 1 000.',
  'ib-dag':'Архитектура BlockDAG','ib-dag-t':'Несколько блоков могут производиться одновременно и объединяться. Более высокая пропускная способность.',
  'ib-gas':'По-Настоящему Бесплатно','ib-gas-t':'Регистрация не стоит ничего. Не нужен ETH, BNB или MATIC. Не нужен банковский счёт.'
},
zh:{
  'tab-register':'🔐 注册','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-protocol':'📜 协议 V6',
  'reg-title':'🔐 注册为已验证人类',
  'reg-sub':'加入Aequitas网络并接收1,000 AEQ。一次性、永久、无Gas费。不存储个人数据。',
  'app-title':'仅通过ANDROID应用注册',
  'app-text':'人类证明需要在您的个人设备上进行生物特征验证。您的指纹由硬件安全元件处理——原始数据永远不会离开您的手机。下载应用，扫描指纹，连接钱包，您的<strong style="color:var(--gold)">1,000 AEQ将自动发放</strong>。',
  's1t':'生物特征扫描','s1d':'打开应用 · 扫描指纹 · HSE处理 · 数据永不离开设备',
  's2t':'ZKP生成','s2d':'生成Groth16证明 · 验证唯一性 · 哈希从不泄露',
  's3t':'连接钱包','s3d':'应用打开MetaMask · 连接钱包 · 地址接收1,000 AEQ',
  's4t':'1,000 AEQ','s4d':'在V6上注册 · 下一个区块内确认 · 应用自动通知',
  'priv-bar':'🔒 硬件安全元件 · Groth16 ZKP · 数据永不离开设备 · 无Gas费',
  'conn-wallet':'已连接钱包','btn-conn':'🦊 连接METAMASK','btn-reg':'🔐 链上注册',
  'btn-add-mm':'+ 添加AEQUITAS网络','reg-log-hint':'// 打开Aequitas Android应用生成您的证明，然后返回此处...',
  'phil':'"货币存在是因为人类存在。<br>仅此而已，不多也不少。"','phil-sub':'— AEQUITAS原则 —',
  's-height':'区块高度','s-height-sub':'每6秒新区块 · BlockDAG · 两个节点并行',
  's-humans':'已验证人类','s-humans-sub':'生物特征ZKP · 一人一钱包，永久',
  's-supply':'总供应量','s-supply-sub':'始终 = 人类 × 1,000 AEQ',
  's-index':'Aequitas指数','s-index-sub':'0 = 完全平等 · 100 = 最大不平等',
  's-uptime':'运行时间','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'人类证明','ib-poh-t':'每个AEQ持有者必须证明自己是唯一的活人。没有机器人、公司或AI。只有真实的人类。',
  'ib-fair':'根本公平的分配','ib-fair-t':'每个经过验证的人类获得恰好1,000 AEQ。总供应量始终等于已验证人类 × 1,000。',
  'ib-dag':'BlockDAG架构','ib-dag-t':'多个区块可以同时产生并合并。更高吞吐量，更低延迟。',
  'ib-gas':'真正无Gas费','ib-gas-t':'注册绝对免费。不需要ETH、BNB或MATIC。不需要银行账户。'
},
id:{
  'tab-register':'🔐 Daftar','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-protocol':'📜 Protokol V6',
  'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi',
  'reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima 1.000 AEQ. Sekali, permanen, tanpa gas. Tidak ada data pribadi.',
  'app-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID',
  'app-text':'Bukti Kemanusiaan memerlukan verifikasi biometrik di perangkat Anda. Sidik jari Anda diproses oleh Hardware Secure Element — data mentah tidak pernah meninggalkan ponsel Anda. Unduh aplikasinya, pindai sidik jari, hubungkan wallet, dan <strong style="color:var(--gold)">1.000 AEQ Anda akan diberikan otomatis</strong>.',
  's1t':'Pemindaian Biometrik','s1d':'Buka aplikasi · pindai sidik jari · HSE memproses · data tidak pernah meninggalkan perangkat',
  's2t':'Pembuatan ZKP','s2d':'Bukti Groth16 dihasilkan · keunikan diverifikasi · hash tidak pernah terungkap',
  's3t':'Hubungkan Wallet','s3d':'Aplikasi membuka MetaMask · hubungkan wallet · alamat menerima 1.000 AEQ',
  's4t':'1.000 AEQ','s4d':'Terdaftar di V6 · dikonfirmasi di blok berikutnya · aplikasi memberi tahu otomatis',
  'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Data tidak pernah meninggalkan perangkat · Tanpa biaya gas',
  'conn-wallet':'DOMPET TERHUBUNG','btn-conn':'🦊 HUBUNGKAN METAMASK','btn-reg':'🔐 DAFTAR ON-CHAIN',
  'btn-add-mm':'+ TAMBAHKAN JARINGAN AEQUITAS','reg-log-hint':'// Buka Aplikasi Android Aequitas untuk menghasilkan bukti Anda, lalu kembali ke sini...',
  'phil':'"Uang ada karena manusia ada.<br>Tidak lebih, tidak kurang."','phil-sub':'— PRINSIP AEQUITAS —',
  's-height':'Tinggi Blok','s-height-sub':'Blok baru setiap 6 detik · BlockDAG · Dua node paralel',
  's-humans':'Manusia Terverifikasi','s-humans-sub':'ZKP biometrik · Satu orang, satu wallet, selamanya',
  's-supply':'Total Pasokan','s-supply-sub':'Selalu = Manusia × 1.000 AEQ',
  's-index':'Indeks Aequitas','s-index-sub':'0 = kesetaraan sempurna · 100 = ketidaksetaraan maksimum',
  's-uptime':'Waktu Aktif','s-uptime-sub':'Node v0.3.0 · Railway + Render · PostgreSQL',
  'ib-poh':'Bukti Kemanusiaan','ib-poh-t':'Setiap pemegang AEQ harus membuktikan bahwa mereka adalah manusia unik yang hidup. Tidak ada bot, korporasi, atau AI. Hanya manusia nyata.',
  'ib-fair':'Distribusi yang Benar-Benar Adil','ib-fair-t':'Setiap manusia terverifikasi menerima tepat 1.000 AEQ. Total pasokan selalu sama dengan manusia × 1.000.',
  'ib-dag':'Arsitektur BlockDAG','ib-dag-t':'Beberapa blok dapat diproduksi secara bersamaan dan digabungkan. Throughput lebih tinggi, latensi lebih rendah.',
  'ib-gas':'Benar-Benar Tanpa Gas','ib-gas-t':'Pendaftaran tidak memerlukan biaya sama sekali. Tidak perlu ETH, BNB, atau MATIC. Tidak perlu rekening bank.'
}
};

function setLang(lang) {
  currentLang = lang;
  document.getElementById('lang-select').value = lang;
  const t = T[lang];
  if (!t) return;
  // Tabs
  const tabs = document.querySelectorAll('.tab');
  const tabKeys = ['tab-register','tab-explorer','tab-humans','tab-index','tab-network','tab-protocol'];
  tabs.forEach((el,i) => { if(tabKeys[i] && t[tabKeys[i]]) el.innerHTML = t[tabKeys[i]]; });
  // Buttons
  const bc = document.getElementById('btn-connect');
  if(bc && !bc.style.background.includes('00E676')) bc.textContent = t['btn-conn'] || bc.textContent;
  const br = document.getElementById('btn-register');
  if(br && !br.textContent.includes('ALREADY') && !br.textContent.includes('READY')) br.textContent = t['btn-reg'] || br.textContent;
  // MetaMask buttons
  document.querySelectorAll('.mm-btn').forEach(el => { el.textContent = t['btn-add-mm'] || el.textContent; });
  // Phil card
  document.querySelectorAll('.phil-quote').forEach(el => { el.innerHTML = t['phil'] || el.innerHTML; });
  document.querySelectorAll('.phil-sub').forEach(el => { el.textContent = t['phil-sub'] || el.textContent; });
  // Stat labels
  const sl = {'s-height':['stat-lbl',0]};
  const statLbls = document.querySelectorAll('.stat-lbl');
  const statSubs = document.querySelectorAll('.stat-sub');
  const statKeys = ['s-height','s-humans','s-supply','s-index','s-uptime'];
  const subKeys  = ['s-height-sub','s-humans-sub','s-supply-sub','s-index-sub','s-uptime-sub'];
  statLbls.forEach((el,i) => { if(statKeys[i] && t[statKeys[i]]) el.textContent = t[statKeys[i]]; });
  statSubs.forEach((el,i) => { if(subKeys[i] && t[subKeys[i]]) el.textContent = t[subKeys[i]]; });
  // Info banner titles/texts
  const ibTitles = document.querySelectorAll('.info-item-title');
  const ibTexts  = document.querySelectorAll('.info-item-text');
  const ibTK = ['ib-poh','ib-fair','ib-dag','ib-gas'];
  const ibXT = ['ib-poh-t','ib-fair-t','ib-dag-t','ib-gas-t'];
  ibTitles.forEach((el,i) => { if(ibTK[i] && t[ibTK[i]]) el.textContent = t[ibTK[i]]; });
  ibTexts.forEach((el,i)  => { if(ibXT[i] && t[ibXT[i]]) el.textContent = t[ibXT[i]]; });
  // Register section
  if(t['reg-title']) { const e=document.querySelector('.reg-hero-title'); if(e) e.textContent=t['reg-title']; }
  if(t['reg-sub'])   { const e=document.querySelector('.reg-hero-sub');   if(e) e.textContent=t['reg-sub']; }
  if(t['app-title']) { const e=document.querySelector('.app-only-title'); if(e) e.textContent=t['app-title']; }
  if(t['app-text'])  { const e=document.querySelector('.app-only-text');  if(e) e.innerHTML=t['app-text']; }
  if(t['priv-bar'])  { const e=document.querySelector('.priv-bar');       if(e) e.textContent=t['priv-bar']; }
  if(t['conn-wallet']){ const e=document.querySelector('.wallet-lbl');    if(e) e.textContent=t['conn-wallet']; }
  // Steps
  const sTitles = document.querySelectorAll('.step-title');
  const sDescs  = document.querySelectorAll('.step-desc');
  [['s1t','s2t','s3t','s4t'],['s1d','s2d','s3d','s4d']].forEach((keys,ki) => {
    const els = ki===0 ? sTitles : sDescs;
    keys.forEach((k,i) => { if(els[i] && t[k]) els[i].textContent=t[k]; });
  });
  // Log hint
  const rlog = document.getElementById('reg-status');
  if(rlog && rlog.innerHTML.includes('Open Aequitas') && t['reg-log-hint'])
    rlog.innerHTML='<span class="info">'+t['reg-log-hint']+'</span>';
}

function fmt(n) {
  if (n === undefined || n === null || n === '—') return '—';
  if (typeof n === 'number') return n.toLocaleString();
  return n;
}

function timeAgo(ts) {
  const d = Math.floor(Date.now() / 1000) - ts;
  if (d < 60) return d + 's ago';
  if (d < 3600) return Math.floor(d / 60) + 'm ago';
  return Math.floor(d / 3600) + 'h ago';
}

function short(h, s = 8, e = 6) { return h ? h.slice(0, s) + '...' + h.slice(-e) : '—'; }

function avatarColor(a) {
  const c = ['#4FC3F7', '#00E676', '#FFB300', '#CE93D8', '#EF5350', '#4DD0E1'];
  return c[parseInt((a || '0x00').slice(2, 4), 16) % c.length];
}

async function addToMetaMask() {
  if (!window.ethereum) { alert('MetaMask not found. Please install MetaMask.'); return; }
  try {
    await window.ethereum.request({
      method: 'wallet_addEthereumChain',
      params: [{
        chainId: CHAIN_ID_HEX,
        chainName: 'Aequitas Chain',
        nativeCurrency: { name: 'AEQ', symbol: 'AEQ', decimals: 18 },
        rpcUrls: ['https://aequitas-production-9fba.up.railway.app/rpc'],
        blockExplorerUrls: ['https://aequitas-production-9fba.up.railway.app']
      }]
    });
  } catch (e) { console.error(e); }
}

async function loadStatus() {
  try {
    const d = await (await fetch('/api/status')).json();
    document.getElementById('s-height').textContent = fmt(d.height);
    document.getElementById('s-humans').textContent = fmt(d.total_humans);
    document.getElementById('s-supply').textContent = d.total_supply || '—';
    document.getElementById('s-index').textContent = fmt(d.index);
    const up = d.uptime || 0;
    const h = Math.floor(up / 3600), m = Math.floor((up % 3600) / 60);
    document.getElementById('s-uptime').textContent = h + 'h ' + m + 'm';
    document.getElementById('idx-score').textContent = fmt(d.index);
    document.getElementById('idx-gini').textContent = typeof d.gini === 'number' ? d.gini.toFixed(3) : '—';
    document.getElementById('idx-supply2').textContent = d.total_supply || '—';
    document.getElementById('idx-phase').textContent = fmt(d.phase);
    document.getElementById('idx-humans2').textContent = fmt(d.total_humans);
    document.getElementById('stat-humans').textContent = fmt(d.total_humans);
    document.getElementById('stat-supply').textContent = d.total_supply || '—';
    if (d.index !== undefined) {
      document.getElementById('idx-bar').style.width = Math.min(d.index, 100) + '%';
      const phaseDesc = ['Phase 0: Bootstrap — Building the network · Cap: 50x fairShare', 'Phase 1: Growth — Expanding human registry · Cap: 20x fairShare', 'Phase 2: Stability — Redistribution active · Cap: 10x fairShare', 'Phase 3: Maturity — Full decentralization · Cap: 3x fairShare'];
      document.getElementById('idx-phase-desc').textContent = phaseDesc[d.phase || 0] || 'Phase ' + (d.phase || 0);
    }
  } catch (e) {}
}

async function loadBlocks() {
  try {
    const blocks = await (await fetch('/api/blocks')).json();
    const list = document.getElementById('blocks-list');
    if (!blocks || !blocks.length) { list.innerHTML = '<div class="empty">No blocks yet</div>'; return; }
    document.getElementById('block-count').textContent = blocks.length + ' blocks';
    list.innerHTML = blocks.map(b => {
      const merge = b.parent_hashes && b.parent_hashes.length > 1;
      const hasTx = b.transactions && b.transactions.length > 0;
      return '<div class="block-item"><div class="block-num">#' + b.height + '</div><div><div class="block-hash">' + short(b.hash) + (merge ? '<span class="badge-merge">🔀 MERGE</span>' : '') + (hasTx ? '<span class="badge-tx">✅ TX</span>' : '') + '</div><div class="block-parents">' + (b.parent_hashes ? b.parent_hashes.length + ' parent(s) · ' + short(b.proposer, 8, 4) : '') + '</div></div><div class="block-right"><div class="block-humans">' + (b.humans || 0) + ' humans</div><div class="block-time">' + timeAgo(b.timestamp) + '</div></div></div>';
    }).join('');
  } catch (e) {}
}

async function loadHumans() {
  try {
    const d = await (await fetch('/api/humans')).json();
    document.getElementById('human-count-badge').textContent = fmt(d.total);
    const list = document.getElementById('humans-list');
    if (!d.humans || !d.humans.length) { list.innerHTML = '<div class="empty">No humans registered yet.<br><br>Download the Aequitas Android App and be the first!</div>'; return; }
    list.innerHTML = d.humans.map(h => {
      const color = avatarColor(h.address || '0x00');
      const init = (h.address || '??').slice(2, 4).toUpperCase();
      return '<div class="human-item"><div class="human-avatar" style="background:' + color + '20;color:' + color + ';border-color:' + color + '50">' + init + '</div><div style="flex:1;min-width:0"><div class="human-balance">' + fmt(h.balance) + ' AEQ</div><div class="human-addr">' + h.address + '</div></div><div class="human-badge">✓ HUMAN</div></div>';
    }).join('');
  } catch (e) {}
}

function checkProofParams() {
  const p = new URLSearchParams(window.location.search);
  const proofId = p.get('proofId');
  const proof = p.get('proof');
  if (proofId) {
    fetch(PROOF_SERVER + '/get/' + proofId).then(r => r.json()).then(pd => {
      proofParams = pd;
      document.getElementById('proof-box').style.display = 'block';
      document.getElementById('proof-val').textContent = '✓ Proof ID: ' + proofId + ' — Connect wallet to register';
      document.querySelectorAll('.tab')[0].click();
      setTimeout(() => connectWallet(), 500);
    }).catch(e => console.error(e));
  } else if (proof) {
    try {
      proofParams = JSON.parse(decodeURIComponent(proof));
      document.getElementById('proof-box').style.display = 'block';
      document.getElementById('proof-val').textContent = '✓ Proof received — Connect wallet to register';
      document.querySelectorAll('.tab')[0].click();
      setTimeout(() => connectWallet(), 500);
    } catch (e) {}
  }
}

async function connectWallet() {
  if (!window.ethereum) return;
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    walletAddr = accounts[0];
    document.getElementById('wallet-box').style.display = 'block';
    document.getElementById('wallet-addr').textContent = walletAddr;
    const btn = document.getElementById('btn-connect');
    btn.textContent = '✓ ' + walletAddr.slice(0, 10) + '...' + walletAddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    try {
      const br = await fetch('/api/balance?wallet=' + walletAddr);
      const bd = await br.json();
      if (bd.is_human) {
        log('✓ Already registered! Balance: ' + bd.balance + ' AEQ', 'ok');
        document.getElementById('btn-register').disabled = true;
        document.getElementById('btn-register').textContent = '✓ ALREADY REGISTERED';
      } else if (proofParams) {
        document.getElementById('btn-register').disabled = false;
        document.getElementById('btn-register').textContent = '🔐 PROOF READY — CLICK TO REGISTER';
      } else {
        document.getElementById('btn-register').disabled = true;
      }
    } catch (e) {
      document.getElementById('btn-register').disabled = !proofParams;
    }
  } catch (e) {}
}

function log(msg, type) {
  const el = document.getElementById('reg-status');
  el.innerHTML += '<div><span class="' + type + '">' + msg + '</span></div>';
}

async function register() {
  if (!walletAddr || !proofParams) return;
  try {
    log('Registering on Aequitas V6...', 'info');
    document.getElementById('btn-register').disabled = true;
    const r = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wallet: walletAddr, pA: proofParams.pA, pB: proofParams.pB, pC: proofParams.pC, pubSignals: proofParams.pubSignals })
    });
    const d = await r.json();
    if (!d.success) { log('✗ ' + d.message, 'err'); document.getElementById('btn-register').disabled = false; return; }
    log('🎉 ' + d.message + ' | TX: ' + d.tx_hash, 'ok');
    setTimeout(() => { window.location.href = '/registered?success=true&wallet=' + walletAddr; }, 1500);
  } catch (e) { log('✗ ' + e.message, 'err'); document.getElementById('btn-register').disabled = false; }
}

checkProofParams();
loadStatus();
loadBlocks();
loadHumans();
setInterval(loadStatus, 6000);
setInterval(loadBlocks, 6000);
setInterval(loadHumans, 10000);

window.ethereum?.on('accountsChanged', a => {
  walletAddr = a[0] || '';
  if (walletAddr) {
    document.getElementById('wallet-box').style.display = 'block';
    document.getElementById('wallet-addr').textContent = walletAddr;
    const btn = document.getElementById('btn-connect');
    btn.textContent = '✓ ' + walletAddr.slice(0, 10) + '...' + walletAddr.slice(-4);
    btn.style.background = 'var(--green)';
    btn.style.color = '#050A14';
    fetch('/api/balance?wallet=' + walletAddr).then(r => r.json()).then(bd => {
      if (bd.is_human) {
        document.getElementById('btn-register').disabled = true;
        document.getElementById('btn-register').textContent = '✓ ALREADY REGISTERED';
        log('✓ Already registered! Balance: ' + bd.balance + ' AEQ', 'ok');
      } else {
        document.getElementById('btn-register').disabled = !proofParams;
        if (proofParams) document.getElementById('btn-register').textContent = '🔐 PROOF READY — CLICK TO REGISTER';
      }
    }).catch(() => { document.getElementById('btn-register').disabled = !proofParams; });
  }
});
</script>
</body>
</html>`
