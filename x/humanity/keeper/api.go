package keeper

import (
"encoding/json"
	"strings"
"fmt"
"io"
"net/http"
"time"
)

type APIServer struct {
blockchain    *BlockDAG
p2pNode       *P2PNode
keeper        *Keeper
startTime     time.Time
proofServerStatus map[string]interface{}
	state         *ChainState
}

func NewAPIServer(bc *BlockDAG, p2p *P2PNode, k *Keeper, state *ChainState) *APIServer {
s := &APIServer{
blockchain:    bc,
p2pNode:       p2p,
keeper:        k,
startTime:     time.Now(),
proofServerStatus: map[string]interface{}{},
		state:         state,
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
	// EVM JSON-RPC
	fmt.Println("── Starting EVM RPC ─────────────────────")
	evmRPC := NewEVMRPCServer(a.blockchain, a.state)
	mux.HandleFunc("/rpc", evmRPC.handleRPC)
	if evmRPC.evm != nil { fmt.Println("✓ EVM Engine ready") } else { fmt.Println("✗ EVM Engine failed") }
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
		"contract_v6":  "0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78",
		"bio_verifier": "0x5bEAAB193a92930fA08c917d6053C66aC6350396",
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

func (a *APIServer) handleUI(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Aequitas Chain Explorer</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{--bg:#050A14;--card:#0D1421;--card2:#111E2E;--border:#1A2D45;--green:#00E676;--blue:#4FC3F7;--gold:#FFB300;--purple:#CE93D8;--red:#EF5350;--text:#E8F4FD;--muted:#6B8CAE;--teal:#4DD0E1}
body{background:var(--bg);color:var(--text);font-family:'Courier New',monospace;min-height:100vh;overflow-x:hidden}
header{background:#080F1E;border-bottom:1px solid var(--border);padding:0 24px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:60px;gap:12px}
.logo-wrap{display:flex;align-items:center;gap:12px;flex-shrink:0}
.logo-icon{width:32px;height:32px;background:var(--gold);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px}
.logo-text{font-size:1.2rem;font-weight:900;color:var(--gold);letter-spacing:5px}
.logo-sub{font-size:0.55rem;color:var(--muted);letter-spacing:3px}
.header-center{display:flex;gap:4px;flex-wrap:wrap;justify-content:center}
.lang-btn{background:#080F1E;color:var(--muted);border:1px solid var(--border);border-radius:5px;padding:4px 10px;cursor:pointer;font-family:monospace;font-size:0.65rem;letter-spacing:1px;transition:all 0.2s}
.lang-btn:hover{color:var(--text);border-color:var(--blue)}
.lang-btn.active{background:var(--blue);color:#050A14;border-color:var(--blue);font-weight:bold}
.header-right{display:flex;gap:8px;align-items:center;flex-shrink:0}
.badge{display:flex;align-items:center;gap:5px;padding:5px 10px;border-radius:16px;font-size:0.65rem;letter-spacing:1px}
.badge-live{background:#00E67612;border:1px solid #00E67628;color:var(--green)}
.badge-dag{background:#4FC3F712;border:1px solid #4FC3F728;color:var(--blue)}
.pulse{width:6px;height:6px;border-radius:50%;background:var(--green);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:0.4}}
.tabs{background:#080F1E;border-bottom:1px solid var(--border);padding:0 24px;display:flex;overflow-x:auto}
.tab{padding:13px 18px;font-size:0.68rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:1px;text-transform:uppercase;white-space:nowrap;transition:all 0.2s}
.tab:hover{color:var(--text)}
.tab.active{color:var(--blue);border-bottom-color:var(--blue)}
.tab-content{display:none}
.tab-content.active{display:block}
.hero{padding:24px 24px 0}
.section-label{font-size:0.58rem;color:var(--muted);letter-spacing:4px;text-transform:uppercase;margin-bottom:16px}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:1px;background:var(--border);border:1px solid var(--border);border-radius:10px;overflow:hidden;margin-bottom:20px}
.stat{background:var(--card);padding:20px 16px;position:relative;overflow:hidden;transition:background 0.2s}
.stat:hover{background:var(--card2)}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:1rem;margin-bottom:8px}
.stat-lbl{font-size:0.58rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:6px}
.stat-val{font-size:1.8rem;font-weight:900;line-height:1;margin-bottom:4px}
.stat-sub{font-size:0.6rem;color:var(--muted);line-height:1.6}
.c-green .stat-val{color:#00E676!important}.c-green .stat-accent{background:#00E676}
.c-blue .stat-val{color:#4FC3F7!important}.c-blue .stat-accent{background:#4FC3F7}
.c-gold .stat-val{color:#FFB300!important}.c-gold .stat-accent{background:#FFB300}
.c-purple .stat-val{color:#CE93D8!important}.c-purple .stat-accent{background:#CE93D8}
.c-teal .stat-val{color:#4DD0E1!important}.c-teal .stat-accent{background:#4DD0E1}
.info-banner{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:20px;margin-bottom:20px;display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:20px}
.info-item-icon{font-size:1.4rem;margin-bottom:8px}
.info-item-title{font-size:0.72rem;color:var(--gold);font-weight:bold;margin-bottom:8px;letter-spacing:1px}
.info-item-text{font-size:0.67rem;color:var(--muted);line-height:1.9}
.main-grid{display:grid;grid-template-columns:1fr 320px;gap:14px;padding:0 24px 24px}
@media(max-width:860px){.main-grid{grid-template-columns:1fr}}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.sec-head{padding:13px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:#080F1E}
.sec-title{font-size:0.65rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;display:flex;align-items:center;gap:7px}
.sec-dot{width:5px;height:5px;border-radius:50%;background:var(--green)}
.sec-count{font-size:0.62rem;color:var(--muted);background:var(--card2);padding:2px 8px;border-radius:8px;border:1px solid var(--border)}
.sec-desc{padding:10px 18px;font-size:0.65rem;color:var(--muted);background:#080F1E;border-bottom:1px solid var(--border);line-height:1.8}
.block-item{padding:12px 18px;border-bottom:1px solid #0A1220;display:grid;grid-template-columns:64px 1fr auto;gap:10px;align-items:center;transition:background 0.15s;cursor:pointer}
.block-item:hover{background:#0D1421}
.block-item:last-child{border-bottom:none}
.block-num{font-size:0.82rem;font-weight:bold;color:var(--blue)}
.block-hash{font-size:0.67rem;color:var(--muted);margin-bottom:2px;display:flex;align-items:center;gap:5px;flex-wrap:wrap}
.block-parents{font-size:0.6rem;color:#3A5570}
.block-right{text-align:right}
.block-humans{font-size:0.68rem;color:var(--gold);margin-bottom:2px}
.block-time{font-size:0.6rem;color:var(--green)}
.badge-merge{background:#2D1B4E;color:var(--purple);font-size:0.56rem;padding:1px 5px;border-radius:3px;border:1px solid #4A2D7A}
.badge-tx{background:#0D2A1A;color:var(--green);font-size:0.56rem;padding:1px 5px;border-radius:3px;border:1px solid #1A4A2A}
.empty{padding:40px;text-align:center;color:var(--muted);font-size:0.72rem;line-height:2.5}
.right-col{display:flex;flex-direction:column;gap:12px}
.info-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px}
.ic-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:14px}
.ic-row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid #0A1220}
.ic-row:last-child{border-bottom:none}
.ic-key{font-size:0.65rem;color:var(--muted)}
.ic-val{font-size:0.65rem;color:var(--text);text-align:right;max-width:60%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ic-val.green{color:var(--green)}.ic-val.blue{color:var(--blue)}.ic-val.gold{color:var(--gold)}.ic-val.purple{color:var(--purple)}
.mm-card{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:16px}
.mm-title{font-size:0.6rem;color:var(--blue);letter-spacing:2px;margin-bottom:12px}
.mm-row{display:flex;justify-content:space-between;padding:5px 0;border-bottom:1px solid #1A2D45}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.6rem;color:var(--muted)}
.mm-val{font-size:0.6rem;color:var(--purple)}
.mm-btn{width:100%;margin-top:10px;padding:9px;background:var(--blue);color:#050A14;border:none;border-radius:7px;cursor:pointer;font-family:monospace;font-size:0.68rem;font-weight:bold;letter-spacing:1px}
.phil-card{background:linear-gradient(135deg,#1A1200,#0D1421);border:1px solid #3A2800;border-radius:10px;padding:20px;text-align:center}
.phil-quote{font-size:0.85rem;color:var(--gold);font-style:italic;line-height:2;margin-bottom:6px}
.phil-sub{font-size:0.6rem;color:var(--muted);letter-spacing:2px}
.humans-section{padding:20px 24px 24px;display:grid;grid-template-columns:1fr 300px;gap:14px}
@media(max-width:860px){.humans-section{grid-template-columns:1fr}}
.human-item{padding:13px 18px;border-bottom:1px solid #0A1220;display:flex;align-items:center;gap:12px;transition:background 0.15s}
.human-item:hover{background:#0D1421}
.human-item:last-child{border-bottom:none}
.human-avatar{width:38px;height:38px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.72rem;font-weight:bold;flex-shrink:0;border:2px solid}
.human-balance{font-size:0.82rem;color:var(--gold);font-weight:bold;margin-bottom:2px}
.human-addr{font-size:0.65rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.human-badge{font-size:0.58rem;padding:2px 8px;border-radius:8px;flex-shrink:0;background:#0D2A1A;color:var(--green);border:1px solid #1A4A2A}
.index-section{padding:20px 24px 24px;display:grid;grid-template-columns:1fr 1fr;gap:14px}
@media(max-width:700px){.index-section{grid-column:1fr}}
.idx-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:22px}
.idx-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:10px}
.idx-desc{font-size:0.7rem;color:var(--muted);line-height:1.9;margin-bottom:16px}
.idx-big{font-size:2.8rem;font-weight:900;color:var(--gold);line-height:1}
.idx-lbl{font-size:0.62rem;color:var(--muted);margin-top:4px}
.bar-bg{height:8px;background:#0D1421;border-radius:4px;overflow:hidden;margin:14px 0 6px}
.bar-fill{height:100%;border-radius:4px;background:linear-gradient(90deg,var(--green),var(--gold),var(--red));transition:width 1.5s}
.bar-labels{display:flex;justify-content:space-between;font-size:0.56rem;color:var(--muted)}
.metrics-row{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-top:14px}
.metric-box{background:#080F1E;border-radius:7px;padding:12px;text-align:center}
.metric-val{font-size:1.2rem;font-weight:bold;color:var(--gold)}
.metric-lbl{font-size:0.58rem;color:var(--muted);margin-top:3px}
.story-text{font-size:0.7rem;line-height:2.1;color:var(--muted)}
.story-text p{margin-bottom:14px}
.story-text p:last-child{margin-bottom:0}
.highlight-box{background:#080F1E;border-left:3px solid var(--gold);border-radius:0 8px 8px 0;padding:14px 18px;margin:16px 0;font-size:0.7rem;color:var(--text);line-height:1.9}
.net-section{padding:20px 24px 24px;display:grid;grid-template-columns:1fr 1fr;gap:14px}
@media(max-width:700px){.net-section{grid-template-columns:1fr}}
.net-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.net-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:14px}
.node-box{background:#080F1E;border-radius:8px;padding:14px;border:1px solid var(--border);margin-bottom:10px}
.node-status{display:flex;align-items:center;gap:6px;font-size:0.7rem;color:var(--green);margin-bottom:5px;font-weight:bold}
.node-dot{width:7px;height:7px;border-radius:50%;background:var(--green);box-shadow:0 0 6px var(--green)}
.node-url{font-size:0.6rem;color:var(--muted);word-break:break-all;margin-bottom:4px}
.node-desc{font-size:0.6rem;color:#3A5570}
.spec-table{width:100%;border-collapse:collapse}
.spec-table td{padding:8px 0;border-bottom:1px solid #0A1220;font-size:0.65rem}
.spec-table tr:last-child td{border-bottom:none}
.spec-table td:first-child{color:var(--muted);width:45%}
.spec-table td:last-child{text-align:right}
.bootstrap-box{background:#080F1E;border-radius:7px;padding:12px;font-size:0.62rem;color:var(--purple);word-break:break-all;line-height:1.8;border:1px solid var(--border)}
.reg-section{padding:20px 24px 24px;max-width:640px;margin:0 auto}
.reg-hero{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:24px;margin-bottom:16px;text-align:center}
.reg-hero-title{font-size:1rem;font-weight:bold;color:var(--text);margin-bottom:8px}
.reg-hero-sub{font-size:0.7rem;color:var(--muted);line-height:1.9}
.app-only{background:#0D1220;border:1px solid #1A2040;border-radius:10px;padding:22px;text-align:center;margin-bottom:16px}
.app-only-icon{font-size:2rem;margin-bottom:8px}
.app-only-title{font-size:0.72rem;color:var(--purple);font-weight:bold;letter-spacing:2px;margin-bottom:10px}
.app-only-text{font-size:0.68rem;color:var(--muted);line-height:1.9}
.reg-steps{display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-bottom:16px}
@media(max-width:560px){.reg-steps{grid-template-columns:repeat(2,1fr)}}
.reg-step{background:var(--card);border:1px solid var(--border);border-radius:9px;padding:16px;text-align:center}
.step-num{width:28px;height:28px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 10px;font-weight:bold;font-size:0.75rem;color:#050A14}
.step-title{font-size:0.65rem;color:var(--text);font-weight:bold;margin-bottom:5px}
.step-desc{font-size:0.62rem;color:var(--muted);line-height:1.7}
.priv-bar{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:10px 14px;margin-bottom:14px;font-size:0.67rem;color:var(--green);text-align:center;line-height:1.8}
.reg-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;margin-bottom:14px}
.wallet-box{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:10px;margin-bottom:10px;display:none}
.wallet-lbl{font-size:0.58rem;color:var(--muted);margin-bottom:3px;letter-spacing:1px}
.wallet-addr{font-size:0.75rem;color:var(--green);font-weight:bold}
.proof-box{background:var(--card2);border:1px solid #3A2800;border-radius:7px;padding:10px;margin-bottom:10px;display:none}
.proof-lbl{font-size:0.58rem;color:var(--gold);margin-bottom:3px;letter-spacing:1px}
.proof-val{font-size:0.65rem;color:var(--muted)}
.reg-btn{width:100%;padding:14px;border-radius:7px;border:none;cursor:pointer;font-family:monospace;font-size:0.78rem;font-weight:bold;letter-spacing:1px;transition:all 0.2s;margin-bottom:8px}
.btn-connect{background:var(--blue);color:#050A14}
.btn-connect:hover{opacity:0.87}
.btn-register{background:var(--gold);color:#050A14}
.btn-register:hover{opacity:0.87}
.reg-btn:disabled{opacity:0.3;cursor:not-allowed}
.reg-log{background:#080F1E;border-radius:7px;padding:12px;font-size:0.67rem;line-height:2;min-height:55px;border:1px solid var(--border)}
.reg-log .ok{color:var(--green)}.reg-log .err{color:var(--red)}.reg-log .info{color:var(--gold)}
</style>
</head>
<body>
<header>
  <div class="logo-wrap">
    <div class="logo-icon">⚖</div>
    <div><div class="logo-text">AEQUITAS</div><div class="logo-sub">CHAIN EXPLORER</div></div>
  </div>
  <div class="header-center">
    <button onclick="setLang('en')" class="lang-btn active" id="lb-en">EN</button>
    <button onclick="setLang('de')" class="lang-btn" id="lb-de">DE</button>
    <button onclick="setLang('es')" class="lang-btn" id="lb-es">ES</button>
    <button onclick="setLang('ru')" class="lang-btn" id="lb-ru">RU</button>
    <button onclick="setLang('zh')" class="lang-btn" id="lb-zh">ZH</button>
    <button onclick="setLang('id')" class="lang-btn" id="lb-id">ID</button>
  </div>
  <div class="header-right">
    <div class="badge badge-live"><span class="pulse"></span><span data-i18n="live">LIVE</span></div>
    <div class="badge badge-dag">● BLOCKDAG</div>
  </div>
</header>

<div class="tabs">
  <div class="tab active" onclick="showTab('explorer',this)" data-i18n-tab="tab-explorer">🔍 Explorer</div>
  <div class="tab" onclick="showTab('humans',this)" data-i18n-tab="tab-humans">👥 Humans</div>
  <div class="tab" onclick="showTab('index',this)" data-i18n-tab="tab-index">📊 Index</div>
  <div class="tab" onclick="showTab('network',this)" data-i18n-tab="tab-network">🌐 Network</div>
  <div class="tab" onclick="showTab('register',this)" data-i18n-tab="tab-register">🔐 Register</div>
  <div class="tab" onclick="showTab('protocol',this)">📜 Protocol V6</div>
</div>

<!-- EXPLORER -->
<div id="tab-explorer" class="tab-content active">
  <div class="hero">
    <div class="section-label" data-i18n="live-stats">Live Chain Statistics</div>
    <div class="stats-grid">
      <div class="stat c-blue"><div class="stat-accent"></div><div class="stat-icon">🔗</div><div class="stat-lbl" data-i18n="block-height">Block Height</div><div class="stat-val" id="s-height">—</div><div class="stat-sub" data-i18n="block-height-sub">New block every 6 seconds · BlockDAG consensus · Two nodes producing blocks in parallel</div></div>
      <div class="stat c-green"><div class="stat-accent"></div><div class="stat-icon">🧬</div><div class="stat-lbl" data-i18n="verified-humans">Verified Humans</div><div class="stat-val" id="s-humans">—</div><div class="stat-sub" data-i18n="verified-humans-sub">Each wallet verified as a unique human · Biometric ZKP · One person, one wallet, forever</div></div>
      <div class="stat c-gold"><div class="stat-accent"></div><div class="stat-icon">🪙</div><div class="stat-lbl" data-i18n="total-supply">Total Supply</div><div class="stat-val" id="s-supply">—</div><div class="stat-sub" data-i18n="total-supply-sub">Always equals Humans × 1,000 AEQ · Supply grows only when humanity grows</div></div>
      <div class="stat c-purple"><div class="stat-accent"></div><div class="stat-icon">⚖</div><div class="stat-lbl" data-i18n="aeq-index">Aequitas Index</div><div class="stat-val" id="s-index">—</div><div class="stat-sub" data-i18n="aeq-index-sub">0 = perfect equality · 100 = maximum inequality · Based on real Gini coefficient</div></div>
      <div class="stat c-teal"><div class="stat-accent"></div><div class="stat-icon">⚡</div><div class="stat-lbl" data-i18n="uptime">Uptime</div><div class="stat-val" id="s-uptime" style="font-size:1.1rem">—</div><div class="stat-sub" data-i18n="uptime-sub">Node v0.3.0 · 2 active nodes · Railway + Render · PostgreSQL persistent state</div></div>
    </div>
    <div class="info-banner">
      <div>
        <div class="info-item-icon">🧬</div>
        <div class="info-item-title" data-i18n="poh-title">Proof of Humanity</div>
        <div class="info-item-text" data-i18n="poh-text">Every single AEQ holder must prove they are a unique, living human being through biometric verification. This is not optional — it is the foundation of the entire system. Without proof of humanity, no AEQ can be received. This means no bots, no duplicate accounts, no corporations, no governments, no AI systems can hold AEQ. Only real humans. The verification uses your fingerprint via your phone's Hardware Secure Element — the same chip that secures your banking apps and your phone's lock screen. Your biometric data never leaves your device under any circumstances.</div>
      </div>
      <div>
        <div class="info-item-icon">⚖</div>
        <div class="info-item-title" data-i18n="fair-title">Radically Fair Distribution</div>
        <div class="info-item-text" data-i18n="fair-text">Every verified human on Earth receives exactly 1,000 AEQ — no more, no less. The first person to register and the billionth person to register receive identical amounts. There is no pre-mine, no founder allocation, no investor round, no early adopter advantage. The total supply of AEQ is always and permanently equal to the number of verified humans multiplied by exactly 1,000. When the number of humans on the network grows, the supply grows proportionally. When no new humans register, no new AEQ is created. This is the most egalitarian monetary distribution system ever designed.</div>
      </div>
      <div>
        <div class="info-item-icon">🔗</div>
        <div class="info-item-title" data-i18n="dag-title">BlockDAG Architecture</div>
        <div class="info-item-text" data-i18n="dag-text">Aequitas does not use a traditional linear blockchain where blocks form a single chain. Instead, it uses a Directed Acyclic Graph (DAG) where multiple blocks can be produced simultaneously by different nodes and later merged. This allows for significantly higher throughput, lower latency, and better fault tolerance. When two nodes produce blocks at the same time, both are valid and are later merged into a single "merge block" — you can see these merge events marked with 🔀 in the block explorer below. This architecture is inspired by projects like IOTA and Phantom but implemented from scratch in Go.</div>
      </div>
      <div>
        <div class="info-item-icon">⛽</div>
        <div class="info-item-title" data-i18n="gasless-title">Truly Gasless Registration</div>
        <div class="info-item-text" data-i18n="gasless-text">One of the biggest barriers to cryptocurrency adoption is the requirement to already own cryptocurrency to pay for transactions. Aequitas eliminates this completely. Registration — the most important transaction on the network — costs absolutely nothing. You do not need ETH, BNB, MATIC, or any other token. You do not need a credit card. You do not need a bank account. You do not need to buy anything. If you are a human being with a smartphone, you can register. The transaction fees are covered by the protocol itself, making Aequitas truly accessible to every person on Earth regardless of their financial situation.</div>
      </div>
    </div>
  </div>
  <div class="main-grid">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="recent-blocks">Recent Blocks</span></div><div class="sec-count" id="block-count">—</div></div>
      <div class="sec-desc" data-i18n="blocks-desc">Each block is cryptographically linked to its parents via SHA-256 hashes. 🔀 MERGE = block with multiple parents (BlockDAG feature, increases throughput). ✅ TX = block containing a registration transaction (a human joined the network). Block time: 6 seconds average.</div>
      <div id="blocks-list"><div class="empty" data-i18n="loading">Loading blocks...</div></div>
    </div>
    <div class="right-col">
      <div class="info-card">
        <div class="ic-title" data-i18n="network-info">Network Info</div>
        <div class="ic-row"><span class="ic-key" data-i18n="chain-name">Chain Name</span><span class="ic-val gold">Aequitas Chain</span></div>
        <div class="ic-row"><span class="ic-key">Chain ID</span><span class="ic-val blue">9001</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="symbol">Symbol</span><span class="ic-val gold">AEQ</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="block-time">Block Time</span><span class="ic-val">6 seconds</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="consensus">Consensus</span><span class="ic-val purple">BlockDAG + PoH</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="nodes">Active Nodes</span><span class="ic-val green">2 Online</span></div>
        <div class="ic-row"><span class="ic-key">ZKP System</span><span class="ic-val">Groth16</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="storage">Storage</span><span class="ic-val green">PostgreSQL</span></div>
      </div>
      <div class="mm-card">
        <div class="mm-title" data-i18n="add-metamask">🦊 ADD TO METAMASK</div>
        <div class="mm-row"><span class="mm-key" data-i18n="network-name">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
        <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.52rem">...9fba.up.railway.app/rpc</span></div>
        <div class="mm-row"><span class="mm-key">Chain ID</span><span class="mm-val">9001</span></div>
        <div class="mm-row"><span class="mm-key" data-i18n="symbol">Symbol</span><span class="mm-val">AEQ</span></div>
        <div class="mm-row"><span class="mm-key" data-i18n="decimals">Decimals</span><span class="mm-val">18</span></div>
        <button class="mm-btn" onclick="addToMetaMask()" data-i18n="add-network">+ ADD AEQUITAS NETWORK</button>
      </div>
      <div class="phil-card">
        <div class="phil-quote" data-i18n="philosophy">"Money exists because people exist.<br>Nothing more, nothing less."</div>
        <div class="phil-sub" data-i18n="philosophy-sub">— THE AEQUITAS PRINCIPLE —</div>
      </div>
    </div>
  </div>
</div>

<!-- HUMANS -->
<div id="tab-humans" class="tab-content">
  <div class="hero">
    <div class="section-label" data-i18n="verified-humans">Verified Humans on Aequitas Chain</div>
    <div class="info-banner">
      <div>
        <div class="info-item-icon">🔒</div>
        <div class="info-item-title" data-i18n="what-is-it">What does "Verified Human" mean?</div>
        <div class="info-item-text" data-i18n="humans-what">A "Verified Human" on Aequitas is a wallet address that has been cryptographically proven to belong to a unique, living human being. This verification is performed using biometric data — specifically a fingerprint scan — processed through the Hardware Secure Element of an Android smartphone. The biometric data itself is never transmitted, stored, or accessible to anyone including the Aequitas team. Only a mathematical Zero-Knowledge Proof derived from the biometric data is used to verify uniqueness. Once a human is verified, that wallet address is permanently and irrevocably linked to that person's biometric identity on the blockchain.</div>
      </div>
      <div>
        <div class="info-item-icon">🧮</div>
        <div class="info-item-title" data-i18n="how-works">The Zero-Knowledge Proof System</div>
        <div class="info-item-text" data-i18n="humans-how">Aequitas uses the Groth16 proving system — one of the most efficient and battle-tested ZKP systems in cryptography, also used by Zcash. When you register, your fingerprint signature is hashed into a field element of the BN128 elliptic curve. This hash is then used as input to a Groth16 circuit that produces a proof: a small mathematical object (just a few hundred bytes) that cryptographically guarantees "a unique biometric hash was used" without revealing what that hash is. The proof server verifies this proof before allowing registration. This means even if someone intercepts the network traffic during registration, they learn absolutely nothing about your biometric data.</div>
      </div>
      <div>
        <div class="info-item-icon">🛡</div>
        <div class="info-item-title" data-i18n="sybil-title">Permanent Sybil Attack Prevention</div>
        <div class="info-item-text" data-i18n="sybil-text">A Sybil attack is when one entity creates multiple fake identities to gain disproportionate influence or resources. It is the fundamental weakness of almost every existing blockchain. Bitcoin is vulnerable (one entity can control mining hardware), Ethereum is vulnerable (one entity can hold many wallets), even proof-of-stake systems are vulnerable. Aequitas solves this at the identity layer: each biometric hash is stored permanently in the PostgreSQL database and on the blockchain. Attempting to register a second time with the same fingerprint is immediately rejected. There is no way to circumvent this — the mathematics make it impossible. One human, one wallet, forever, guaranteed by cryptography not by trust.</div>
      </div>
      <div>
        <div class="info-item-icon">🌍</div>
        <div class="info-item-title" data-i18n="global-title">Designed for Global Inclusion</div>
        <div class="info-item-text" data-i18n="global-text">Aequitas is explicitly designed to be accessible to every human on Earth, regardless of their economic situation, location, or access to traditional financial infrastructure. You do not need a bank account. You do not need a credit card. You do not need to buy any existing cryptocurrency. You do not need a government ID. You do not need internet access beyond a basic smartphone connection. All you need is an Android smartphone with a fingerprint sensor — a device that over 3 billion people already own. The registration is completely free, takes under 2 minutes, and grants you 1,000 AEQ instantly. This is financial inclusion at a scale never before achieved.</div>
      </div>
    </div>
  </div>
  <div class="humans-section">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="registered-humans">Registered Humans</span></div><div class="sec-count" id="human-count-badge">0</div></div>
      <div class="sec-desc" data-i18n="humans-desc">Every address listed here has been verified as a unique human through biometric Zero-Knowledge Proof. Each received exactly 1,000 AEQ upon registration. These registrations are permanent, immutable, and stored both in PostgreSQL and on the blockchain as transactions.</div>
      <div id="humans-list"><div class="empty" data-i18n="no-humans">No humans registered yet.<br><br>Download the Aequitas Android App<br>and be the first human on the chain!</div></div>
    </div>
    <div class="right-col">
      <div class="info-card">
        <div class="ic-title" data-i18n="registry-stats">Registry Stats</div>
        <div class="ic-row"><span class="ic-key" data-i18n="total-humans-stat">Total Humans</span><span class="ic-val green" id="stat-humans">0</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="total-supply">Total Supply</span><span class="ic-val gold" id="stat-supply">0 AEQ</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="grant">Grant per Human</span><span class="ic-val gold">1,000 AEQ</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="reg-fee">Registration Fee</span><span class="ic-val green" data-i18n="free">FREE</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="zkp-system">ZKP System</span><span class="ic-val">Groth16</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="storage">Biometric Storage</span><span class="ic-val green" data-i18n="never-stored">Never stored</span></div>
      </div>
      <div class="info-card">
        <div class="ic-title" data-i18n="zkp-title">ZKP Technical Details</div>
        <div style="font-size:0.65rem;color:var(--muted);line-height:1.9" data-i18n="zkp-details">The Groth16 proving system operates over the BN128 elliptic curve. Proof size: ~200 bytes. Verification time: ~10ms. The circuit was compiled using snarkjs and circom. The trusted setup used the Hermez ceremony parameters. Proof generation happens on the Aequitas Proof Server (Node.js) after receiving the biometric hash from the Android app.</div>
      </div>
    </div>
  </div>
</div>

<!-- INDEX -->
<div id="tab-index" class="tab-content">
  <div class="index-section">
    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="aeq-index-title">Aequitas Index — Real-Time Economic Equality Score</div>
      <div class="idx-desc" data-i18n="aeq-index-desc">The Aequitas Index is a composite metric that measures the economic health of the Aequitas network in real time. It is calculated directly from the on-chain balance distribution of all verified humans. A score of 0 means every verified human has exactly the same amount of AEQ — perfect equality. A score of 100 means one person controls all the AEQ and everyone else has nothing — maximum inequality. The index is calculated from the real Gini coefficient of on-chain balances, adjusted for network size and phase. The protocol uses this index to automatically trigger redistribution mechanisms when inequality grows beyond safe thresholds, with no human decision-making required.</div>
      <div style="display:grid;grid-template-columns:auto 1fr;gap:24px;align-items:center;margin-top:14px">
        <div><div class="idx-big" id="idx-score">—</div><div class="idx-lbl" data-i18n="current-index">Current Index</div></div>
        <div>
          <div class="bar-bg"><div class="bar-fill" id="idx-bar" style="width:0%"></div></div>
          <div class="bar-labels"><span data-i18n="bar-0">0 — Perfect Equality</span><span>50</span><span data-i18n="bar-100">100 — Max Inequality</span></div>
          <div style="margin-top:10px;font-size:0.67rem;color:var(--muted);background:#080F1E;padding:9px;border-radius:6px" id="idx-phase-desc">—</div>
        </div>
      </div>
      <div class="metrics-row" style="grid-template-columns:repeat(4,1fr)">
        <div class="metric-box"><div class="metric-val" id="idx-gini">—</div><div class="metric-lbl" data-i18n="gini-coeff">Gini Coefficient</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-supply2">—</div><div class="metric-lbl" data-i18n="total-supply">Total Supply</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-phase">—</div><div class="metric-lbl" data-i18n="phase">Protocol Phase</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-humans2">—</div><div class="metric-lbl" data-i18n="verified-humans">Verified Humans</div></div>
      </div>
    </div>

    <div class="idx-card">
      <div class="idx-title" data-i18n="pools-title">Redistribution Pools</div>
      <div class="idx-desc" data-i18n="pools-desc">When inequality thresholds are exceeded, AEQ is automatically redirected into these four pools. The pools are smart contract accounts controlled entirely by protocol logic — no human has access to them. Each pool serves a specific purpose in maintaining economic health.</div>
      <div class="metrics-row">
        <div class="metric-box"><div class="metric-val" id="pool-v">—</div><div class="metric-lbl" data-i18n="vel-pool">Velocity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-l">—</div><div class="metric-lbl" data-i18n="liq-pool">Liquidity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-u">—</div><div class="metric-lbl" data-i18n="uni-pool">Unity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-t">—</div><div class="metric-lbl" data-i18n="treasury">Treasury</div></div>
      </div>
      <div style="margin-top:14px;font-size:0.65rem;color:var(--muted);line-height:1.9" data-i18n="pools-details">The Velocity Pool rewards wallets that transact regularly, incentivizing economic activity. The Liquidity Pool supports market depth and price stability. The Unity Pool grants additional AEQ to newly registered humans in later phases, welcoming them with a larger starting balance. The Treasury funds protocol development, security audits, and infrastructure costs — all spending is publicly visible on-chain.</div>
    </div>

    <div class="idx-card">
      <div class="idx-title" data-i18n="phases-title">Protocol Phases</div>
      <div class="idx-desc" data-i18n="phases-desc">The Aequitas protocol is designed to evolve through four distinct phases as the network grows. Each phase unlocks new mechanisms and increases the sophistication of the economic governance system. Phase transitions happen automatically based on the number of verified humans and the Gini coefficient — no governance vote, no human decision required.</div>
      <table class="spec-table">
        <tr><td>Phase 0</td><td style="color:var(--green)" data-i18n="phase0">Bootstrap — Building the network · &lt;100 humans</td></tr>
        <tr><td>Phase 1</td><td style="color:var(--blue)" data-i18n="phase1">Growth — Expanding human registry · 100–10,000 humans</td></tr>
        <tr><td>Phase 2</td><td style="color:var(--gold)" data-i18n="phase2">Stability — Redistribution active · 10,000–1M humans</td></tr>
        <tr><td>Phase 3</td><td style="color:var(--purple)" data-i18n="phase3">Maturity — Full decentralization · 1M+ humans · Gini &lt;0.3</td></tr>
      </table>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="gini-title">The Gini Coefficient — How We Measure Inequality</div>
      <div class="idx-desc" data-i18n="gini-desc">The Gini coefficient was developed by Italian statistician Corrado Gini in 1912 and remains the world's most widely used measure of economic inequality. It is used by the World Bank, the IMF, the UN, and virtually every economist on Earth. The coefficient is calculated by comparing the actual distribution of wealth to a perfectly equal distribution. A Gini of 0 means everyone has exactly the same amount. A Gini of 1 means one person owns everything and everyone else has nothing. In practice, no country has ever achieved a Gini below 0.15 or above 0.70. Aequitas calculates the Gini coefficient in real time from the actual on-chain balances of all verified humans and uses it as the primary input for the redistribution mechanism.</div>
      <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-top:14px">
        <div class="metric-box" style="border:1px solid #1A4A2A"><div class="metric-val" style="color:var(--green)">0.00</div><div class="metric-lbl" data-i18n="gini-0">Perfect Equality</div><div style="font-size:0.6rem;color:var(--muted);margin-top:4px;line-height:1.6" data-i18n="gini-0-sub">Every person has identical wealth. Theoretically impossible to sustain in a free market, but Aequitas approaches this in Phase 0 when all humans have exactly 1,000 AEQ.</div></div>
        <div class="metric-box" style="border:1px solid #1A2D45"><div class="metric-val" style="color:var(--blue)">0.27</div><div class="metric-lbl" data-i18n="gini-1">Low Inequality</div><div style="font-size:0.6rem;color:var(--muted);margin-top:4px;line-height:1.6" data-i18n="gini-1-sub">Scandinavia average (Sweden 0.27, Denmark 0.28). Achieved through strong social safety nets, progressive taxation, and high minimum wages. The Aequitas target zone.</div></div>
        <div class="metric-box" style="border:1px solid #3A2800"><div class="metric-val" style="color:var(--gold)">0.41</div><div class="metric-lbl" data-i18n="gini-2">Moderate Inequality</div><div style="font-size:0.6rem;color:var(--muted);margin-top:4px;line-height:1.6" data-i18n="gini-2-sub">USA average (0.41). Redistribution mechanisms activate at this level in Aequitas Phase 2. Wealth is concentrated but a significant middle class exists.</div></div>
        <div class="metric-box" style="border:1px solid #4A1A1A"><div class="metric-val" style="color:var(--red)">0.63</div><div class="metric-lbl" data-i18n="gini-3">High Inequality</div><div style="font-size:0.6rem;color:var(--muted);margin-top:4px;line-height:1.6" data-i18n="gini-3-sub">South Africa (0.63), the world's most unequal country. At this level, Aequitas activates aggressive redistribution. Bitcoin's estimated Gini is above 0.85.</div></div>
      </div>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="inflation-title">The Inflation & Redistribution Mechanism — In Detail</div>
      <div class="story-text" data-i18n="inflation-text">
        <p><span style="color:var(--gold);font-weight:bold">Why does Aequitas need an inflation mechanism?</span> In a world where people trade, save, invest, and lose money, perfect equality cannot be maintained forever even with an equal starting distribution. Over time, some people will accumulate more AEQ and some will spend theirs. Without a correction mechanism, Aequitas would eventually look like any other unequal monetary system. The inflation and redistribution mechanism is the answer.</p>
        <p><strong style="color:var(--blue)">Base Inflation — The Only Truly Justified Inflation:</strong> The only source of new AEQ is human registration. When a new human is verified and joins the network, exactly 1,000 AEQ is created and sent to their wallet. This is the only inflation in Phase 0 and Phase 1. There is no other mechanism that creates AEQ. No mining rewards, no staking rewards, no protocol emissions. The supply grows if and only if the number of verified humans grows. This means AEQ's purchasing power is directly tied to human population growth — historically one of the most stable and predictable growth rates in the world.</p>
        <p><strong style="color:var(--gold)">The Wealth Cap — Preventing Extreme Concentration:</strong> In Phase 2 and beyond, a dynamic wealth cap is enforced. The cap is calculated based on total supply, number of humans, and the current Gini coefficient. When any wallet's balance exceeds this cap, the excess is not seized — instead, all new AEQ earned by that wallet (from velocity rewards, pool distributions, etc.) is redirected to the four redistribution pools until the balance falls below the cap. This is not confiscation — it is a ceiling on accumulation, applied fairly to everyone including the founders.</p>
        <p><strong style="color:var(--purple)">Dynamic Redistribution Cycles — Automatic Economic Correction:</strong> The Keeper Bot runs redistribution cycles on a schedule. During each cycle, it reads the current Gini coefficient from on-chain balances. If the Gini exceeds 0.25 (Phase 2 threshold), a small percentage of the Velocity Pool is distributed pro-rata to all verified humans. If the Gini exceeds 0.35, the distribution percentage increases. If the Gini exceeds 0.50, emergency redistribution activates with a larger percentage. The higher the inequality, the more aggressive the correction — creating a powerful automatic stabilizer built into the protocol itself.</p>
        <p><strong style="color:var(--green)">Velocity Incentives — Rewarding Economic Participation:</strong> The single biggest driver of inequality in any monetary system is hoarding — a small number of entities accumulate vast wealth and simply hold it, removing it from circulation and making it unavailable to others. The Velocity Pool directly combats hoarding by rewarding wallets that transact regularly. Every time a verified human sends or receives AEQ, they earn a small velocity score. At the end of each cycle, the Velocity Pool is distributed proportionally to velocity scores. This creates a powerful incentive to keep money moving through the economy rather than hoarding it.</p>
        <p><strong style="color:var(--teal)">Mathematical Governance — The End of Monetary Politics:</strong> Every rule described above is encoded in immutable smart contract code. No person, organization, government, or founding team can change these rules without a hard fork that requires community consensus. The Aequitas protocol is governed by mathematics, not by committee meetings, political lobbying, or the preferences of wealthy investors. This is the fundamental promise of Aequitas: a monetary system where the rules apply equally to everyone and cannot be changed by anyone with power or money.</p>
      </div>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="story-title">The Story of Aequitas — Why This Exists</div>
      <div class="story-text" data-i18n="story-text">
        <p>The year is 2009. Satoshi Nakamoto releases Bitcoin — the first successful decentralized digital currency. For the first time in history, it is possible to transfer value between any two people on Earth without the permission of any bank, government, or intermediary. It is a genuine revolution. But something goes wrong almost immediately.</p>
        <p>Early Bitcoin miners accumulate millions of coins that cost them almost nothing. As Bitcoin's price rises, these early adopters become extraordinarily wealthy. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. The cryptocurrency that was supposed to democratize finance has instead created some of the most extreme wealth concentration in human history. Bitcoin's estimated Gini coefficient exceeds 0.85 — higher than any country on Earth, higher than historical feudal societies, approaching the theoretical maximum of inequality.</p>
        <p>Ethereum, Solana, and virtually every other cryptocurrency follow the same pattern: pre-mine for founders and early investors, ICO for wealthy participants, and the rest of humanity priced out before they even hear about it. The blockchain revolution failed to deliver on its promise of financial inclusion.</p>
        <p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness," "equity," and "equality" — was created to answer one question: <em style="color:var(--gold)">"What would a cryptocurrency look like if it was designed from first principles to be fair to every human being who has ever lived or will ever live?"</em></p>
        <p>The answer turned out to be surprisingly simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being a person.</strong> Not because they were born into wealth. Not because they were early to a speculative bet. Not because they have access to expensive mining hardware. Simply because they are human.</p>
        <p>The technology to make this possible now exists. Zero-Knowledge Proofs allow us to verify that a person is a unique human being without requiring them to reveal any personal information. Blockchain technology allows us to store these verifications permanently and transparently. Smartphones — now owned by over 3 billion people — provide the biometric sensors needed for verification. The pieces are in place. Aequitas assembles them into a coherent system for the first time.</p>
        <p>The Aequitas network launched in June 2026. It is currently in Phase 0 — the bootstrap phase. Every human who registers now is helping to prove that a fairer monetary system is possible. The goal is not to replace existing currencies overnight. The goal is to demonstrate, on a live network with real cryptography and real economics, that money can be distributed fairly, that equality can be maintained through mathematical governance, and that financial inclusion can be achieved at global scale without compromising security or decentralization.</p>
        <p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em> This is not just a slogan — it is the mathematical foundation of the entire system. Every line of code, every cryptographic proof, every protocol rule flows from this single insight.</p>
      </div>
    </div>
  </div>
</div>

<!-- NETWORK -->
<div id="tab-network" class="tab-content">
  <div class="net-section">
    <div class="net-card" style="grid-column:1/-1">
      <div class="net-title" data-i18n="active-nodes">Active Nodes — Current Network Topology</div>
      <div style="font-size:0.68rem;color:var(--muted);line-height:1.9;margin-bottom:14px" data-i18n="nodes-desc">The Aequitas network currently operates on two nodes running in geographically distributed cloud environments. Both nodes participate in block production, state synchronization, and API serving. They communicate via libp2p (the same P2P networking library used by IPFS and Ethereum 2.0) and sync blocks via HTTP. Both nodes share the same PostgreSQL database for persistent state, ensuring that registered humans are visible across the entire network instantly.</div>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px">
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span><span data-i18n="node1">Node 1 — Railway (Primary)</span></div>
          <div class="node-url">aequitas-production-9fba.up.railway.app</div>
          <div class="node-desc" data-i18n="node1-desc">Primary API server · Block producer · P2P bootstrap node · PostgreSQL connection · RPC endpoint for MetaMask</div>
        </div>
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span><span data-i18n="node2">Node 2 — Render (Secondary)</span></div>
          <div class="node-url">aequitas-node-2.onrender.com</div>
          <div class="node-desc" data-i18n="node2-desc">Secondary API server · Block producer · P2P peer · HTTP block sync · Shared PostgreSQL · Redundancy node</div>
        </div>
      </div>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="bootstrap-title">Bootstrap Node — Join the Network</div>
      <div style="margin-bottom:12px;font-size:0.67rem;color:var(--muted);line-height:1.9" data-i18n="bootstrap-desc">To run your own Aequitas node and join the network, connect to the bootstrap node using the libp2p multiaddress below. Your node will automatically discover other peers, download the full block history, and begin participating in consensus. The bootstrap node runs on Railway and is available 24/7.</div>
      <div class="bootstrap-box">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
      <div style="margin-top:10px;font-size:0.62rem;color:var(--muted);line-height:1.7" data-i18n="bootstrap-howto">To run a node: clone the GitHub repository, set the DATABASE_URL environment variable, and run the binary with the bootstrap peer address. Full documentation is available on GitHub.</div>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="tech-specs">Technical Specifications</div>
      <table class="spec-table">
        <tr><td data-i18n="chain-id">Chain ID</td><td style="color:var(--blue)">9001 (0x2329)</td></tr>
        <tr><td>EVM</td><td style="color:var(--green)" data-i18n="evm-yes">Yes — JSON-RPC at /rpc · MetaMask compatible</td></tr>
        <tr><td data-i18n="block-time">Block Time</td><td>~6 seconds average</td></tr>
        <tr><td data-i18n="consensus">Consensus</td><td style="color:var(--purple)">BlockDAG + Proof of Humanity</td></tr>
        <tr><td>P2P Protocol</td><td>libp2p (Go implementation)</td></tr>
        <tr><td>ZKP System</td><td>Groth16 / snarkjs / circom</td></tr>
        <tr><td>Curve</td><td>BN128 (alt-bn128)</td></tr>
        <tr><td data-i18n="storage">State Storage</td><td style="color:var(--green)">PostgreSQL (persistent)</td></tr>
        <tr><td data-i18n="language">Language</td><td>Go 1.21</td></tr>
        <tr><td data-i18n="source">Source Code</td><td><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub ↗</a></td></tr>
      </table>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="metamask-config">MetaMask / Web3 Configuration</div>
      <div style="font-size:0.65rem;color:var(--muted);line-height:1.9;margin-bottom:12px" data-i18n="metamask-desc">Aequitas Chain is EVM-compatible and can be added to MetaMask or any other Web3 wallet that supports custom RPC networks. Use the configuration below to add Aequitas to your wallet and view your AEQ balance.</div>
      <table class="spec-table">
        <tr><td data-i18n="network-name">Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
        <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.58rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
        <tr><td>Chain ID</td><td style="color:var(--blue)">9001</td></tr>
        <tr><td data-i18n="symbol">Currency Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
        <tr><td data-i18n="decimals">Decimals</td><td>18</td></tr>
      </table>
      <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:12px" data-i18n="add-network">+ ADD TO METAMASK</button>
    </div>
    <div class="net-card" style="grid-column:1/-1">
      <div class="net-title" data-i18n="architecture-title">System Architecture Overview</div>
      <div style="font-size:0.67rem;color:var(--muted);line-height:1.9" data-i18n="architecture-desc">The Aequitas system consists of four main components working together. The Android App handles biometric scanning and proof generation entirely on-device — no sensitive data ever leaves the phone. The Proof Server (Node.js on Railway) receives the biometric hash, generates the Groth16 ZK proof, stores the hash in PostgreSQL to prevent double registration, and returns the proof. The Blockchain Nodes (Go on Railway + Render) maintain the BlockDAG, process registration transactions, manage account balances, and expose the EVM-compatible RPC endpoint. PostgreSQL (Railway managed database) stores persistent state — account balances, biometric hashes, and registered wallets — ensuring data survives node restarts and is shared across all nodes. The entire stack is open source and auditable.</div>
    </div>
  </div>
</div>


<!-- PROTOCOL V6 -->
<div id="tab-protocol" class="tab-content">
  <div style="padding:20px 24px 24px;max-width:860px;margin:0 auto">
    <div class="section-label">Aequitas V6 Protocol — Complete Technical Documentation</div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">Why V6? The Evolution of Fair Money</div>
      <div class="story-text">
        <p>AequitasV6 is the first version of the protocol to run entirely on the Aequitas Chain — a sovereign blockchain built from scratch in Go, with a real EVM execution engine powered by go-ethereum. Previous versions ran on Ethereum Sepolia testnet. V6 is deployed at <span style="color:var(--blue)">0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78</span> on Chain ID 9001.</p>
        <p>V6 introduces five new mechanisms that make Aequitas the most sophisticated fair monetary system ever designed: Proof of Alive, the Guardian System, Demurrage, an always-active Wealth Cap, and UBI from protocol economics — not taxation.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">1. PROOF OF ALIVE — Keeping the Money Supply Real</div>
      <div class="story-text">
        <p>Proof of Alive solves one of the most fundamental problems in cryptocurrency: what happens to money when people die or disappear? In Bitcoin, an estimated 3-4 million BTC are permanently lost because people died or lost their keys. This money is gone forever — it can never circulate, never benefit anyone.</p>
        <p>In Aequitas, money represents people. If a person disappears from the network, their AEQ should eventually return to the community — not disappear forever.</p>
        <div class="highlight-box">
          THE PROCESS (designed to be generous, not punitive):<br>
          Year 0-2: Normal usage expected<br>
          Year 2: Warning 1 sent on-chain — Guardian can respond<br>
          Year 2 + 60 days: Warning 2 sent — Guardian can respond<br>
          Year 2 + 120 days: Warning 3 sent — Guardian can respond<br>
          Year 2 + 180 days: AEQ moved to PERSONAL ESCROW (NOT UBI Pool yet)<br>
          Year 4: If still no activity → AEQ enters UBI Pool → distributed equally
        </div>
        <p>Why Escrow first, not immediate UBI? Someone could be in a 3-year prison sentence. Their AEQ goes to Escrow in year 2, but they return in year 3. They get their Escrow back PLUS the current fairShare. They are not punished for being imprisoned.</p>
        <p>REACTIVATION: Submit fresh biometric Proof of Alive → receive Escrow back (if still held) + current fairShare. The same biometric commitment stays blocked permanently — no double-dipping possible.</p>
        <p>Why fairShare instead of original 1,000 AEQ on return? Because fair means equal — not historical. If 1 million people have registered and fairShare is 1,200 AEQ, the returning person gets 1,200 AEQ. This is more equitable than the original 1,000 AEQ from years ago.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">2. GUARDIAN SYSTEM — Protection for the Vulnerable</div>
      <div class="story-text">
        <p>The Guardian System answers: What happens to someone's AEQ if they cannot access their device for months or years? Consider these real scenarios: a person imprisoned for 3 years, a person hospitalized in a coma, an elderly person who lost their phone, a person in a war zone without internet.</p>
        <p>In Bitcoin, these people's funds would be frozen forever. In Aequitas, a trusted Guardian can confirm they are still alive on their behalf.</p>
        <div class="highlight-box">
          GUARDIAN RULES:<br>
          • Every verified human can appoint 1 Guardian (another verified human)<br>
          • Guardian can ONLY call confirmAlive() — zero transaction rights<br>
          • Guardian CANNOT move funds, transfer AEQ, or change anything<br>
          • Maximum 3 wards per Guardian<br>
          • 7-day timelock on Guardian assignment (prevents forced assignment under duress)<br>
          • After 3 consecutive Guardian confirmations without self-activity → community review flag raised<br>
          • No circular relationships (A guards B → B cannot guard A)<br>
          • Guardian cannot have their own Guardian — prevents layered control chains
        </div>
        <p>Example: Maria is imprisoned for 2 years. Before going to prison, she appointed her sister Ana as Guardian. Every year, Ana confirms on-chain that Maria is still alive. Maria's AEQ remains safe. When Maria is released, she confirms herself and the Guardian role ends automatically.</p>
        <p>The 7-day timelock is critical in high-crime areas: if someone is forced at gunpoint to appoint a criminal as Guardian, they have 7 days to cancel the pending assignment. This window cannot be shortened by anyone.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">3. DEMURRAGE — The Anti-Hoarding Mechanism</div>
      <div class="story-text">
        <p>Demurrage is one of the oldest ideas in monetary theory, first proposed by economist Silvio Gesell in 1916. Aequitas implements it in its purest digital form: a 1% annual fee on any balance ABOVE your fairShare.</p>
        <div class="highlight-box">
          EXAMPLE WITH NUMBERS:<br>
          Total supply: 10,000 AEQ (10 humans × 1,000 AEQ each)<br>
          fairShare: 1,000 AEQ per person<br>
          Your balance: 3,000 AEQ<br>
          Excess above fairShare: 2,000 AEQ<br>
          Monthly demurrage: 2,000 × 1% ÷ 12 = 1.67 AEQ per month<br>
          That 1.67 AEQ → UBI Pool → distributed equally to all 10 humans<br>
          Your cost: 1.67 AEQ/month. Each other human gains: +0.17 AEQ/month
        </div>
        <p>What demurrage is NOT: it is not a wealth tax (only applies to excess above fairShare), not punitive (1% annual is extremely gentle), does not reduce total supply (money goes to UBI Pool, not deleted), and is not inflation (new money is only created when new humans register).</p>
        <p>Historical proof: The town of Wörgl, Austria (1932) introduced demurrage currency during the Great Depression. Within one year, unemployment dropped 25% while the rest of Austria suffered. Money circulated 12x faster because people preferred to spend rather than hoard. The Central Bank shut it down — not because it failed, but because it worked too well.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">4. WEALTH CAP — Mathematical Redistribution</div>
      <div class="story-text">
        <p>The Wealth Cap is a hard ceiling on how much AEQ any single human can hold. When a wallet exceeds the cap, the excess is INSTANTLY redistributed equally to ALL active humans — automatically, on every transfer.</p>
        <div class="highlight-box">
          PHASE-BASED CAP:<br>
          Phase 0 (1-100 humans):   cap = 50x fairShare<br>
          Phase 1 (101-1,000):      cap = 20x fairShare<br>
          Phase 2 (1,001-10,000):   cap = 10x fairShare<br>
          Phase 3 (10,001-100,000): cap =  5x fairShare<br>
          Phase 4 (100,000+):       cap =  3x fairShare
        </div>
        <p>CRITICAL V6 DESIGN DECISION: The cap is ALWAYS active from human #1. Earlier versions had no cap in Phase 0 — this was wrong. The first 100 people could accumulate unlimited AEQ, creating a permanent oligarchy. V6 fixes this: cap starts at 50x fairShare from the very beginning and tightens as the network grows.</p>
        <p>Example (Phase 4, 1M humans, fairShare = 1,000 AEQ): cap = 3,000 AEQ. The wealthiest possible human holds 3x what the newest human receives. For comparison: the top 1% of Bitcoin holders own over 90% of all Bitcoin. In Aequitas, mathematical law makes that impossible.</p>
        <p>Why does the cap decrease as the network grows? In a small network, variance is acceptable for market discovery. In a large network with millions of humans, a 3x cap is already extremely egalitarian and sufficient for economic activity.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">5. UNIVERSAL BASIC INCOME — From Protocol, Not Politics</div>
      <div class="story-text">
        <p>Aequitas implements UBI not as a political choice but as a mathematical consequence of the protocol's fairness mechanisms. It requires no taxation, no government, no political decision.</p>
        <div class="highlight-box">
          SOURCES OF THE UBI POOL:<br>
          1. Transaction fees: 0.1% of every transfer → 20% to UBI Pool<br>
          2. Wealth cap overflow → redistributed instantly and equally to all humans<br>
          3. Demurrage: 1% annual on excess balances → UBI Pool<br>
          4. Inactive wallet escrow: after 4 years of inactivity → UBI Pool<br><br>
          DISTRIBUTION: Every month, UBI Pool ÷ total active humans = equal payment to everyone
        </div>
        <p>Example: 1,000 active humans. UBI Pool has 500 AEQ from fees and demurrage. Each human receives 0.5 AEQ. This happens automatically every month, forever, with no human intervention.</p>
        <p>Why this is different from political UBI: Political UBI requires taxation, redistribution through government, and political decisions. Aequitas UBI comes from the protocol's own economic activity, is distributed equally to every verified human with no bureaucracy, happens automatically on-chain, and cannot be stopped or changed by any authority.</p>
        <p>The big picture: As the network grows and more transactions happen, the UBI Pool grows. More humans → more economic activity → larger UBI → more incentive to join → more humans. This is a positive feedback loop where economic fairness drives adoption.</p>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">6. NO ALGORITHMIC INFLATION — Money Only From Humans</div>
      <div class="story-text">
        <p>V6 removes all algorithmic inflation from previous versions. The ONLY event that creates new AEQ is: a new verified human registers. One human = 1,000 AEQ. That is the complete and entire monetary policy.</p>
        <p>Why this matters: Previous versions had algorithmic inflation (0-1.5% annual) based on velocity, activity, and Gini scores. These parameters could theoretically be manipulated. V6 makes manipulation impossible: no external parameters, no oracle, no governance vote can change the money supply. Only human beings registering their biometric identity creates new money.</p>
        <div class="highlight-box">
          TOTAL SUPPLY FORMULA (always true, always verifiable):<br>
          Total AEQ = Verified Active Humans × 1,000<br><br>
          When humans register: supply increases<br>
          When humans become inactive (after 4 years): supply remains (escrow)<br>
          When escrow releases to UBI: supply remains (just redistributed)<br>
          Demurrage: supply remains (just redistributed)<br>
          Wealth cap overflow: supply remains (just redistributed)<br>
          NO OTHER MONEY CREATION IS POSSIBLE
        </div>
      </div>
    </div>

    <div class="idx-card" style="margin-bottom:14px">
      <div class="idx-title">7. CONTRACT ADDRESSES</div>
      <div class="story-text">
        <div class="highlight-box">
          Chain: Aequitas Chain (Chain ID: 9001)<br>
          RPC: https://aequitas-production-9fba.up.railway.app/rpc<br><br>
          BioVerifier (Groth16):  0x5bEAAB193a92930fA08c917d6053C66aC6350396<br>
          AequitasV6 (Main):      0x20D271028f32577FCd07b4583A8e0E4eBBdB4F78<br><br>
          V5 (Sepolia legacy):    0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5
        </div>
      </div>
    </div>

  </div>
</div>

<!-- REGISTER -->
<div id="tab-register" class="tab-content">
  <div class="reg-section">
    <div class="reg-hero">
      <div class="reg-hero-title" data-i18n="reg-title">🔐 Register as a Verified Human</div>
      <div class="reg-hero-sub" data-i18n="reg-sub">Join the Aequitas network and receive your 1,000 AEQ. This is a one-time, permanent, gasless registration that cryptographically proves you are a unique human being. No personal data is ever collected or stored. No gas fees. No waiting. Once registered, your wallet is permanently verified on the Aequitas blockchain.</div>
    </div>
    <div class="app-only">
      <div class="app-only-icon">📱</div>
      <div class="app-only-title" data-i18n="app-only-title">REGISTRATION VIA ANDROID APP ONLY</div>
      <div class="app-only-text" data-i18n="app-only-text">Proof of Humanity requires direct biometric verification on your personal device. This is a deliberate security decision — the fingerprint data must be processed by the Hardware Secure Element (HSE) of your phone, a dedicated cryptographic chip that is isolated from the main processor and cannot be accessed remotely. The HSE is the same chip that secures your banking apps, your payment cards, and your phone's own security. When you scan your fingerprint in the Aequitas App, the signature is processed by the HSE, hashed into a mathematical object, and used to generate a Zero-Knowledge Proof — all without the raw fingerprint data ever leaving the HSE. Download the Aequitas Android App, scan your fingerprint, connect your Web3 wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically and immediately</strong>.</div>
    </div>
    <div class="reg-steps">
      <div class="reg-step">
        <div class="step-num">1</div>
        <div class="step-title" data-i18n="step1-title">Biometric Scan</div>
        <div class="step-desc" data-i18n="step1-desc">Open the Aequitas App and tap "Prove Humanity." You will be prompted to scan your fingerprint. The scan is processed entirely by your phone's Hardware Secure Element — the raw fingerprint data never leaves this chip, never reaches the app's main process, and never leaves your device under any circumstances. The HSE generates a cryptographic signature that serves as your unique biometric identifier.</div>
      </div>
      <div class="reg-step">
        <div class="step-num">2</div>
        <div class="step-title" data-i18n="step2-title">ZKP Generation</div>
        <div class="step-desc" data-i18n="step2-desc">The app derives a biometric hash from your fingerprint signature and sends it to the Aequitas Proof Server. The Proof Server checks that this hash has not been used before (preventing double registration), then generates a Groth16 Zero-Knowledge Proof over the BN128 elliptic curve. This proof mathematically guarantees "a unique biometric hash was used" without revealing anything about the hash itself. The entire process takes 2-5 seconds.</div>
      </div>
      <div class="reg-step">
        <div class="step-num">3</div>
        <div class="step-title" data-i18n="step3-title">Connect Wallet</div>
        <div class="step-desc" data-i18n="step3-desc">The app opens MetaMask (or your preferred Web3 wallet) with the Aequitas Chain pre-configured. Connect your wallet — this is the address that will receive your 1,000 AEQ. Make sure you control this wallet and have access to the private key. Once registered, the wallet address is permanently linked to your biometric identity on the blockchain. You cannot change it later without re-registering.</div>
      </div>
      <div class="reg-step">
        <div class="step-num">4</div>
        <div class="step-title" data-i18n="step4-title">1,000 AEQ Granted</div>
        <div class="step-desc" data-i18n="step4-desc">The registration transaction is submitted to the Aequitas blockchain. The protocol verifies the ZK proof, checks that the biometric hash is unique, and immediately credits exactly 1,000 AEQ to your wallet address. The transaction is recorded in the next block (within 6 seconds), stored permanently in PostgreSQL, and visible in the block explorer. The app confirms automatically when registration is complete. Total time from fingerprint scan to receiving AEQ: under 30 seconds.</div>
      </div>
    </div>
    <div class="priv-bar" data-i18n="priv-bar">🔒 Hardware Secure Element · Groth16 ZKP · Biometric data never leaves device · No personal data collected · No gas fees · Permanent Sybil protection · Immutable on-chain record</div>
    <div class="reg-card">
      <div class="wallet-box" id="wallet-box"><div class="wallet-lbl" data-i18n="connected-wallet">CONNECTED WALLET</div><div class="wallet-addr" id="wallet-addr">—</div></div>
      <div class="proof-box" id="proof-box"><div class="proof-lbl" data-i18n="proof-detected">⚡ PROOF PARAMETERS DETECTED FROM APP</div><div class="proof-val" id="proof-val">—</div></div>
      <button class="reg-btn btn-connect" id="btn-connect" onclick="connectWallet()" data-i18n="connect-btn">🦊 CONNECT METAMASK</button>
      <button class="reg-btn btn-register" id="btn-register" onclick="register()" disabled data-i18n="register-btn">🔐 REGISTER ON-CHAIN</button>
      <div class="reg-log" id="reg-status"><span class="info" data-i18n="reg-hint">// Open Aequitas Android App to generate your proof, then return here to complete registration...</span></div>
    </div>
    <div class="info-card">
      <div class="ic-title" data-i18n="reg-details">Registration Details</div>
      <div class="ic-row"><span class="ic-key" data-i18n="network">Network</span><span class="ic-val purple">Aequitas Chain (BlockDAG)</span></div>
      <div class="ic-row"><span class="ic-key">Chain ID</span><span class="ic-val gold">9001</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="grant">Grant Amount</span><span class="ic-val gold">1,000 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="reg-fee">Gas Fee</span><span class="ic-val green" data-i18n="free">FREE (gasless)</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="reg-limit">Registrations</span><span class="ic-val" data-i18n="reg-limit-val">Once per human · permanent · immutable</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="bio-data">Biometric Data</span><span class="ic-val green" data-i18n="never-stored">Never stored anywhere</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="confirmation">Confirmation Time</span><span class="ic-val" data-i18n="conf-time">Within 6 seconds (next block)</span></div>
    </div>
  </div>
</div>

<script>
const PROOF_SERVER='https://aequitas-proof-server-production.up.railway.app';
let walletAddr='',proofParams=null,currentLang='en';

const T={
en:{
'live':'LIVE','tab-explorer':'🔍 Explorer','tab-humans':'👥 Humans','tab-index':'📊 Index','tab-network':'🌐 Network','tab-register':'🔐 Register',
'live-stats':'Live Chain Statistics','block-height':'Block Height','block-height-sub':'New block every 6 seconds · BlockDAG consensus · Two nodes producing blocks in parallel','verified-humans':'Verified Humans','verified-humans-sub':'Each wallet verified as a unique human · Biometric ZKP · One person, one wallet, forever','total-supply':'Total Supply','total-supply-sub':'Always equals Humans × 1,000 AEQ · Supply grows only when humanity grows','aeq-index':'Aequitas Index','aeq-index-sub':'0 = perfect equality · 100 = maximum inequality · Based on real Gini coefficient','uptime':'Uptime','uptime-sub':'Node v0.3.0 · 2 active nodes · Railway + Render · PostgreSQL persistent state',
'poh-title':'Proof of Humanity','poh-text':'Every single AEQ holder must prove they are a unique, living human being through biometric verification. This is not optional — it is the foundation of the entire system. Without proof of humanity, no AEQ can be received. This means no bots, no duplicate accounts, no corporations, no governments, no AI systems can hold AEQ. Only real humans. The verification uses your fingerprint via your phone\'s Hardware Secure Element — the same chip that secures your banking apps. Your biometric data never leaves your device under any circumstances.',
'fair-title':'Radically Fair Distribution','fair-text':'Every verified human on Earth receives exactly 1,000 AEQ — no more, no less. The first person to register and the billionth person to register receive identical amounts. There is no pre-mine, no founder allocation, no investor round, no early adopter advantage. The total supply of AEQ is always permanently equal to verified humans × 1,000. When new humans join, new AEQ is created. When no new humans register, no new AEQ is created. This is the most egalitarian monetary distribution system ever designed.',
'dag-title':'BlockDAG Architecture','dag-text':'Aequitas uses a Directed Acyclic Graph (DAG) where multiple blocks can be produced simultaneously by different nodes and later merged. This allows significantly higher throughput, lower latency, and better fault tolerance. When two nodes produce blocks at the same time, both are valid and are later merged — you can see these merge events marked with 🔀 in the explorer. This architecture allows the network to scale without sacrificing decentralization.',
'gasless-title':'Truly Gasless Registration','gasless-text':'One of the biggest barriers to cryptocurrency adoption is requiring you to own cryptocurrency first to pay for transactions. Aequitas eliminates this completely. Registration costs absolutely nothing — no ETH, no BNB, no MATIC required. No credit card, no bank account. If you are a human being with a smartphone, you can register. Transaction fees are covered by the protocol itself, making Aequitas truly accessible to every person on Earth regardless of their financial situation.',
'recent-blocks':'Recent Blocks','blocks-desc':'Each block is cryptographically linked to its parents via SHA-256 hashes. 🔀 MERGE = block with multiple parents (BlockDAG feature). ✅ TX = block containing a registration transaction (a human joined). Block time: 6 seconds average.','loading':'Loading blocks...','network-info':'Network Info','chain-name':'Chain Name','symbol':'Symbol','block-time':'Block Time','consensus':'Consensus','nodes':'Active Nodes','storage':'State Storage','add-metamask':'🦊 ADD TO METAMASK','network-name':'Network Name','add-network':'+ ADD AEQUITAS NETWORK','philosophy':'"Money exists because people exist.<br>Nothing more, nothing less."','philosophy-sub':'— THE AEQUITAS PRINCIPLE —','decimals':'Decimals',
'what-is-it':'What does Verified Human mean?','humans-what':'A Verified Human on Aequitas is a wallet address cryptographically proven to belong to a unique, living human. Verification uses biometric data — specifically a fingerprint — processed through the Hardware Secure Element of an Android smartphone. The biometric data itself is never transmitted or stored. Only a Zero-Knowledge Proof derived from it is used. Once verified, the wallet is permanently linked to that person\'s biometric identity on the blockchain.',
'how-works':'The Zero-Knowledge Proof System','humans-how':'Aequitas uses the Groth16 proving system — also used by Zcash — one of the most efficient ZKP systems in cryptography. Your fingerprint signature is hashed into a field element of the BN128 elliptic curve. This hash is used as input to a Groth16 circuit that produces a small mathematical proof guaranteeing "a unique biometric hash was used" without revealing what the hash is. Proof size: ~200 bytes. Verification time: ~10ms.',
'sybil-title':'Permanent Sybil Attack Prevention','sybil-text':'A Sybil attack is when one entity creates multiple fake identities. It is the fundamental weakness of almost every blockchain. Bitcoin, Ethereum, and most cryptocurrencies are vulnerable. Aequitas solves this at the identity layer: each biometric hash is stored permanently. Attempting to register twice with the same fingerprint is immediately rejected. One human, one wallet, forever — guaranteed by cryptography, not trust.',
'global-title':'Designed for Global Inclusion','global-text':'Aequitas is designed to be accessible to every human on Earth. You need no bank account, no credit card, no existing cryptocurrency. Just an Android smartphone with a fingerprint sensor — a device over 3 billion people already own. Registration is completely free, takes under 2 minutes, and grants 1,000 AEQ instantly. Financial inclusion at a scale never before achieved.',
'registered-humans':'Registered Humans','humans-desc':'Every address listed here has been verified as a unique human through biometric Zero-Knowledge Proof. Each received exactly 1,000 AEQ upon registration. These registrations are permanent, immutable, and stored both in PostgreSQL and on the blockchain.','no-humans':'No humans registered yet.\n\nDownload the Aequitas Android App and be the first human on the chain!','registry-stats':'Registry Stats','total-humans-stat':'Total Humans','grant':'Grant per Human','reg-fee':'Registration Fee','free':'FREE','zkp-system':'ZKP System','never-stored':'Never stored','zkp-title':'ZKP Technical Details','zkp-details':'Groth16 over BN128 elliptic curve. Proof size: ~200 bytes. Verification time: ~10ms. Circuit compiled with snarkjs and circom. Trusted setup used the Hermez ceremony parameters. Proof generation on the Aequitas Proof Server (Node.js) after receiving the biometric hash from the Android app.',
'aeq-index-title':'Aequitas Index — Real-Time Economic Equality Score','aeq-index-desc':'The Aequitas Index measures the economic health of the Aequitas network in real time. Calculated directly from on-chain balance distribution of all verified humans. Score 0 = every human has identical AEQ (perfect equality). Score 100 = one person controls everything (maximum inequality). Calculated from the real Gini coefficient of on-chain balances, adjusted for network size and phase. The protocol uses this index to automatically trigger redistribution when inequality grows beyond safe thresholds.','current-index':'Current Index','bar-0':'0 — Perfect Equality','bar-100':'100 — Max Inequality','gini-coeff':'Gini Coefficient','phase':'Protocol Phase',
'pools-title':'Redistribution Pools','pools-desc':'When inequality thresholds are exceeded, AEQ is automatically redirected into these four pools. The pools are smart contract accounts controlled entirely by protocol logic — no human has access to them. Each pool serves a specific purpose in maintaining economic health.','vel-pool':'Velocity Pool','liq-pool':'Liquidity Pool','uni-pool':'Unity Pool','treasury':'Treasury','pools-details':'The Velocity Pool rewards wallets that transact regularly, incentivizing economic activity. The Liquidity Pool supports market depth and price stability. The Unity Pool grants additional AEQ to newly registered humans in later phases. The Treasury funds protocol development, security audits, and infrastructure costs — all spending is publicly visible on-chain.',
'phases-title':'Protocol Phases','phases-desc':'The Aequitas protocol evolves through four distinct phases as the network grows. Each phase unlocks new mechanisms and increases the sophistication of economic governance. Phase transitions happen automatically based on verified humans and Gini coefficient — no governance vote required.','phase0':'Bootstrap — Building the network · &lt;100 humans','phase1':'Growth — Expanding human registry · 100–10,000 humans','phase2':'Stability — Redistribution active · 10,000–1M humans','phase3':'Maturity — Full decentralization · 1M+ humans · Gini &lt;0.3',
'gini-title':'The Gini Coefficient — How We Measure Inequality','gini-desc':'The Gini coefficient was developed by Italian statistician Corrado Gini in 1912 and remains the world\'s most widely used measure of economic inequality. Used by the World Bank, IMF, and virtually every economist on Earth. Calculated by comparing actual wealth distribution to a perfectly equal distribution. Aequitas calculates this in real time from actual on-chain balances and uses it as the primary input for the redistribution mechanism.',
'gini-0':'Perfect Equality','gini-0-sub':'Every person has identical wealth. Aequitas approaches this in Phase 0 when all humans have exactly 1,000 AEQ.','gini-1':'Low Inequality','gini-1-sub':'Scandinavia average (Sweden 0.27, Denmark 0.28). The Aequitas target zone for long-term operation.','gini-2':'Moderate Inequality','gini-2-sub':'USA average (0.41). Redistribution mechanisms activate at this level in Aequitas Phase 2.','gini-3':'High Inequality','gini-3-sub':'South Africa (0.63). Bitcoin\'s estimated Gini exceeds 0.85 — higher than any country on Earth.',
'inflation-title':'Inflation & Redistribution Mechanism — In Detail',
'story-title':'The Story of Aequitas — Why This Exists',
'active-nodes':'Active Nodes — Current Network Topology','nodes-desc':'The Aequitas network operates on two nodes in geographically distributed cloud environments. Both nodes participate in block production, state synchronization, and API serving. They communicate via libp2p (also used by IPFS and Ethereum 2.0) and sync blocks via HTTP. Both nodes share the same PostgreSQL database, ensuring registered humans are visible across the entire network instantly.','node1':'Node 1 — Railway (Primary)','node2':'Node 2 — Render (Secondary)','node1-desc':'Primary API server · Block producer · P2P bootstrap node · PostgreSQL connection · RPC endpoint for MetaMask','node2-desc':'Secondary API server · Block producer · P2P peer · HTTP block sync · Shared PostgreSQL · Redundancy node',
'bootstrap-title':'Bootstrap Node — Join the Network','bootstrap-desc':'To run your own Aequitas node, connect to the bootstrap node using the libp2p multiaddress below. Your node will automatically discover peers, download the full block history, and begin participating in consensus. The bootstrap node runs on Railway and is available 24/7.','bootstrap-howto':'To run a node: clone the GitHub repository, set the DATABASE_URL environment variable, and run the binary with the bootstrap peer address. Full documentation available on GitHub.',
'tech-specs':'Technical Specifications','chain-id':'Chain ID','evm-yes':'Yes — JSON-RPC at /rpc · MetaMask compatible','language':'Language','source':'Source Code',
'metamask-config':'MetaMask / Web3 Configuration','metamask-desc':'Aequitas Chain is EVM-compatible and can be added to MetaMask or any Web3 wallet supporting custom RPC networks. Use the configuration below to add Aequitas to your wallet and view your AEQ balance.',
'architecture-title':'System Architecture Overview','architecture-desc':'The Aequitas system has four main components. The Android App handles biometric scanning and proof generation entirely on-device. The Proof Server (Node.js) generates Groth16 ZK proofs and stores biometric hashes. The Blockchain Nodes (Go) maintain the BlockDAG, process transactions, and expose the EVM RPC. PostgreSQL stores persistent state shared across all nodes. The entire stack is open source.',
'reg-title':'🔐 Register as a Verified Human','reg-sub':'Join the Aequitas network and receive your 1,000 AEQ. A one-time, permanent, gasless registration that cryptographically proves you are a unique human. No personal data collected or stored. No gas fees. No waiting. Once registered, your wallet is permanently verified on the Aequitas blockchain.',
'app-only-title':'REGISTRATION VIA ANDROID APP ONLY','app-only-text':'Proof of Humanity requires biometric verification on your personal device. The fingerprint must be processed by your phone\'s Hardware Secure Element (HSE) — a dedicated cryptographic chip isolated from the main processor and inaccessible remotely. The HSE is the same chip securing your banking apps and phone lock screen. When you scan in the Aequitas App, the signature is processed by the HSE, hashed, and used to generate a Zero-Knowledge Proof — all without raw fingerprint data ever leaving the HSE. Download the app, scan your fingerprint, connect your wallet, and your <strong style="color:var(--gold)">1,000 AEQ will be granted automatically and immediately</strong>.',
'step1-title':'Biometric Scan','step1-desc':'Open the Aequitas App and tap "Prove Humanity." Scan your fingerprint — processed entirely by your phone\'s Hardware Secure Element. Raw fingerprint data never reaches the app\'s main process and never leaves your device. The HSE generates a cryptographic signature as your unique biometric identifier.',
'step2-title':'ZKP Generation','step2-desc':'The app sends your biometric hash to the Proof Server. The server checks it has not been used before (preventing double registration), then generates a Groth16 Zero-Knowledge Proof over BN128. This proof guarantees "a unique biometric hash was used" without revealing anything about the hash. Takes 2-5 seconds.',
'step3-title':'Connect Wallet','step3-desc':'The app opens MetaMask with Aequitas Chain pre-configured. Connect your wallet — this address will receive your 1,000 AEQ. Once registered, this wallet is permanently linked to your biometric identity. You cannot change it without re-registering.',
'step4-title':'1,000 AEQ Granted','step4-desc':'The registration transaction is submitted to the blockchain. The protocol verifies the ZK proof, confirms the biometric hash is unique, and credits exactly 1,000 AEQ to your wallet within the next block (under 6 seconds). The transaction is stored permanently on-chain and in PostgreSQL. Total time from fingerprint to AEQ: under 30 seconds.',
'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Biometric data never leaves device · No personal data collected · No gas fees · Permanent Sybil protection · Immutable on-chain record',
'connected-wallet':'CONNECTED WALLET','proof-detected':'⚡ PROOF PARAMETERS DETECTED FROM APP','connect-btn':'🦊 CONNECT METAMASK','register-btn':'🔐 REGISTER ON-CHAIN','reg-hint':'// Open Aequitas Android App to generate your proof, then return here to complete registration...','reg-details':'Registration Details','network':'Network','reg-limit':'Registrations','reg-limit-val':'Once per human · permanent · immutable','bio-data':'Biometric Data','never-stored':'Never stored anywhere','confirmation':'Confirmation Time','conf-time':'Within 6 seconds (next block)',
'phase-desc-0':'Phase 0: Bootstrap — Building the network and onboarding early humans','phase-desc-1':'Phase 1: Growth — Expanding the human registry globally','phase-desc-2':'Phase 2: Stability — Wealth redistribution mechanisms active','phase-desc-3':'Phase 3: Maturity — Full decentralization achieved'
},
de:{
'live':'LIVE','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-register':'🔐 Registrieren',
'live-stats':'Live Chain-Statistiken','block-height':'Blockhöhe','block-height-sub':'Neuer Block alle 6 Sekunden · BlockDAG-Konsens · Zwei Nodes produzieren Blöcke parallel','verified-humans':'Verifizierte Menschen','verified-humans-sub':'Jede Wallet als einzigartiger Mensch verifiziert · Biometrischer ZKP · Eine Person, eine Wallet, für immer','total-supply':'Gesamtmenge','total-supply-sub':'Immer gleich Menschen × 1.000 AEQ · Angebot wächst nur wenn die Menschheit wächst','aeq-index':'Aequitas-Index','aeq-index-sub':'0 = vollkommene Gleichheit · 100 = maximale Ungleichheit · Basiert auf echtem Gini-Koeffizient','uptime':'Betriebszeit','uptime-sub':'Node v0.3.0 · 2 aktive Nodes · Railway + Render · PostgreSQL persistenter Zustand',
'poh-title':'Menschlichkeitsnachweis','poh-text':'Jeder einzelne AEQ-Inhaber muss beweisen, dass er ein einzigartiger, lebender Mensch ist – durch biometrische Verifizierung. Dies ist nicht optional – es ist das Fundament des gesamten Systems. Ohne Menschlichkeitsnachweis kann kein AEQ empfangen werden. Das bedeutet: keine Bots, keine Duplikat-Konten, keine Unternehmen, keine Regierungen, keine KI-Systeme können AEQ halten. Nur echte Menschen. Die Verifizierung nutzt deinen Fingerabdruck über das Hardware Secure Element deines Smartphones – denselben Chip, der deine Banking-Apps schützt. Deine biometrischen Daten verlassen niemals dein Gerät.',
'fair-title':'Radikal faire Verteilung','fair-text':'Jeder verifizierte Mensch auf der Erde erhält genau 1.000 AEQ – nicht mehr, nicht weniger. Die erste Person, die sich registriert, und die milliardste Person erhalten identische Beträge. Es gibt kein Pre-Mine, keine Gründerzuteilung, keine Investorenrunde, keinen Erstmover-Vorteil. Die Gesamtmenge an AEQ ist immer und dauerhaft gleich der Anzahl der verifizierten Menschen multipliziert mit genau 1.000. Dies ist das egalitärste Geldsystem, das je entworfen wurde.',
'dag-title':'BlockDAG-Architektur','dag-text':'Aequitas verwendet einen Directed Acyclic Graph (DAG), bei dem mehrere Blöcke gleichzeitig von verschiedenen Nodes produziert und später zusammengeführt werden können. Dies ermöglicht deutlich höheren Durchsatz, niedrigere Latenz und bessere Fehlertoleranz. Wenn zwei Nodes gleichzeitig Blöcke produzieren, sind beide gültig und werden später zusammengeführt – du kannst diese Zusammenführungsereignisse im Explorer mit 🔀 sehen.',
'gasless-title':'Wirklich gebührenfreie Registrierung','gasless-text':'Eine der größten Hürden bei der Kryptowährungsnutzung ist die Anforderung, bereits Kryptowährung zu besitzen, um Transaktionsgebühren zu zahlen. Aequitas eliminiert dies vollständig. Die Registrierung kostet absolut nichts – kein ETH, kein BNB, kein MATIC benötigt. Keine Kreditkarte, kein Bankkonto. Wenn du ein Mensch mit einem Smartphone bist, kannst du dich registrieren.',
'recent-blocks':'Aktuelle Blöcke','blocks-desc':'Jeder Block ist kryptographisch über SHA-256-Hashes mit seinen Eltern verbunden. 🔀 MERGE = Block mit mehreren Eltern (BlockDAG-Funktion). ✅ TX = Block mit Registrierungstransaktion. Blockzeit: durchschnittlich 6 Sekunden.','loading':'Blöcke werden geladen...','network-info':'Netzwerkinformationen','chain-name':'Netzwerkname','symbol':'Symbol','block-time':'Blockzeit','consensus':'Konsens','nodes':'Aktive Nodes','storage':'Zustandsspeicher','add-metamask':'🦊 ZU METAMASK HINZUFÜGEN','network-name':'Netzwerkname','add-network':'+ AEQUITAS-NETZWERK HINZUFÜGEN','philosophy':'"Geld existiert weil Menschen existieren.<br>Nichts mehr, nichts weniger."','philosophy-sub':'— DAS AEQUITAS-PRINZIP —','decimals':'Dezimalstellen',
'what-is-it':'Was bedeutet "Verifizierter Mensch"?','humans-what':'Ein "Verifizierter Mensch" bei Aequitas ist eine Wallet-Adresse, die kryptographisch bewiesen wurde, einem einzigartigen, lebenden Menschen zu gehören. Die Verifizierung verwendet biometrische Daten – einen Fingerabdruckscan – verarbeitet durch das Hardware Secure Element eines Android-Smartphones. Die biometrischen Daten selbst werden niemals übertragen oder gespeichert.','how-works':'Das Zero-Knowledge-Proof-System','humans-how':'Aequitas verwendet das Groth16-Beweissystem – auch von Zcash verwendet – eines der effizientesten ZKP-Systeme in der Kryptographie. Deine Fingerabdrucksignatur wird in ein Feldelement der BN128-elliptischen Kurve gehasht. Dieser Hash wird als Eingabe für einen Groth16-Schaltkreis verwendet, der einen kleinen mathematischen Beweis erzeugt: "Ein einzigartiger biometrischer Hash wurde verwendet" – ohne preiszugeben, was der Hash ist.',
'sybil-title':'Permanenter Schutz vor Sybil-Angriffen','sybil-text':'Ein Sybil-Angriff liegt vor, wenn eine Entität mehrere gefälschte Identitäten erstellt. Es ist die grundlegende Schwäche fast jeder Blockchain. Bitcoin, Ethereum und die meisten Kryptowährungen sind anfällig. Aequitas löst dies auf der Identitätsebene: Jeder biometrische Hash wird dauerhaft gespeichert. Der Versuch, sich zweimal mit demselben Fingerabdruck zu registrieren, wird sofort abgelehnt. Eine Person, eine Wallet, für immer – garantiert durch Kryptographie, nicht durch Vertrauen.',
'global-title':'Für globale Inklusion konzipiert','global-text':'Aequitas wurde so konzipiert, dass es für jeden Menschen auf der Erde zugänglich ist. Du brauchst kein Bankkonto, keine Kreditkarte, keine bestehende Kryptowährung. Nur ein Android-Smartphone mit einem Fingerabdrucksensor – ein Gerät, das bereits über 3 Milliarden Menschen besitzen. Die Registrierung ist vollständig kostenlos, dauert unter 2 Minuten und gewährt sofort 1.000 AEQ.',
'registered-humans':'Registrierte Menschen','humans-desc':'Jede hier aufgeführte Adresse wurde durch biometrischen Zero-Knowledge-Beweis als einzigartiger Mensch verifiziert. Jeder erhielt genau 1.000 AEQ bei der Registrierung. Diese Registrierungen sind dauerhaft, unveränderlich und sowohl in PostgreSQL als auch on-chain gespeichert.','no-humans':'Noch keine Menschen registriert.\n\nLade die Aequitas Android-App herunter und sei der erste Mensch auf der Chain!','registry-stats':'Registrierungsstatistik','total-humans-stat':'Gesamte Menschen','grant':'Zuteilung pro Mensch','reg-fee':'Registrierungsgebühr','free':'KOSTENLOS','zkp-system':'ZKP-System','never-stored':'Niemals gespeichert','zkp-title':'Technische ZKP-Details','zkp-details':'Groth16 über BN128-elliptischer Kurve. Beweisdauer: ~200 Bytes. Verifizierungszeit: ~10ms. Schaltkreis mit snarkjs und circom kompiliert. Trusted Setup verwendete die Hermez-Zeremonieparameter. Beweiserzeugung auf dem Aequitas Proof Server (Node.js).',
'aeq-index-title':'Aequitas-Index — Wirtschaftlicher Gleichheitswert in Echtzeit','aeq-index-desc':'Der Aequitas-Index misst die wirtschaftliche Gesundheit des Aequitas-Netzwerks in Echtzeit. Berechnet direkt aus der On-Chain-Bilanzverteilung aller verifizierten Menschen. Score 0 = jeder Mensch hat identisches AEQ (vollkommene Gleichheit). Score 100 = eine Person kontrolliert alles (maximale Ungleichheit). Basiert auf dem echten Gini-Koeffizienten der On-Chain-Bilanzen, angepasst für Netzwerkgröße und Phase.','current-index':'Aktueller Index','bar-0':'0 — Vollkommene Gleichheit','bar-100':'100 — Max. Ungleichheit','gini-coeff':'Gini-Koeffizient','phase':'Protokollphase',
'pools-title':'Umverteilungspools','pools-desc':'Wenn Ungleichheitsschwellenwerte überschritten werden, wird AEQ automatisch in diese vier Pools umgeleitet. Die Pools sind Smart-Contract-Konten, die ausschließlich durch Protokolllogik kontrolliert werden – kein Mensch hat Zugriff.','vel-pool':'Velocity-Pool','liq-pool':'Liquiditäts-Pool','uni-pool':'Unity-Pool','treasury':'Tresor','pools-details':'Der Velocity-Pool belohnt Wallets, die regelmäßig handeln, und schafft Anreize für wirtschaftliche Aktivität. Der Liquiditäts-Pool unterstützt Markttiefe und Preisstabilität. Der Unity-Pool gewährt neu registrierten Menschen in späteren Phasen zusätzliches AEQ. Der Tresor finanziert Protokollentwicklung, Sicherheitsaudits und Infrastrukturkosten.',
'phases-title':'Protokollphasen','phases-desc':'Das Aequitas-Protokoll entwickelt sich durch vier verschiedene Phasen, wenn das Netzwerk wächst. Jede Phase schaltet neue Mechanismen frei. Phasenübergänge erfolgen automatisch basierend auf verifizierten Menschen und Gini-Koeffizient – keine Abstimmung erforderlich.','phase0':'Bootstrap — Netzwerk aufbauen · &lt;100 Menschen','phase1':'Wachstum — Menschenregister erweitern · 100–10.000 Menschen','phase2':'Stabilität — Umverteilung aktiv · 10.000–1M Menschen','phase3':'Reife — Vollständige Dezentralisierung · 1M+ Menschen · Gini &lt;0,3',
'gini-title':'Der Gini-Koeffizient — Wie wir Ungleichheit messen','gini-desc':'Der Gini-Koeffizient wurde 1912 vom italienischen Statistiker Corrado Gini entwickelt und ist das weltweit am häufigsten verwendete Maß für wirtschaftliche Ungleichheit. Wird von der Weltbank, dem IWF und praktisch jedem Ökonomen auf der Erde verwendet. Aequitas berechnet diesen Koeffizienten in Echtzeit aus tatsächlichen On-Chain-Bilanzen.',
'gini-0':'Vollkommene Gleichheit','gini-0-sub':'Alle haben identisches Vermögen. Aequitas nähert sich diesem Wert in Phase 0, wenn alle Menschen genau 1.000 AEQ haben.','gini-1':'Geringe Ungleichheit','gini-1-sub':'Skandinavien-Durchschnitt (Schweden 0,27, Dänemark 0,28). Die Zielzone von Aequitas für den Langzeitbetrieb.','gini-2':'Moderate Ungleichheit','gini-2-sub':'USA-Durchschnitt (0,41). Umverteilungsmechanismen aktivieren auf diesem Niveau in Aequitas Phase 2.','gini-3':'Hohe Ungleichheit','gini-3-sub':'Südafrika (0,63). Bitcoins geschätzter Gini übersteigt 0,85 – höher als jedes Land der Erde.',
'inflation-title':'Inflation & Umverteilungsmechanismus — Im Detail',
'story-title':'Die Geschichte von Aequitas — Warum das existiert',
'active-nodes':'Aktive Nodes — Aktuelle Netzwerktopologie','nodes-desc':'Das Aequitas-Netzwerk betreibt zwei Nodes in geografisch verteilten Cloud-Umgebungen. Beide Nodes nehmen an der Blockproduktion, Statussynchronisation und API-Bereitstellung teil. Sie kommunizieren über libp2p (auch von IPFS und Ethereum 2.0 verwendet) und synchronisieren Blöcke über HTTP. Beide Nodes teilen dieselbe PostgreSQL-Datenbank.','node1':'Node 1 — Railway (Primär)','node2':'Node 2 — Render (Sekundär)','node1-desc':'Primärer API-Server · Blockproduzent · P2P-Bootstrap-Node · PostgreSQL-Verbindung · RPC-Endpunkt für MetaMask','node2-desc':'Sekundärer API-Server · Blockproduzent · P2P-Peer · HTTP-Block-Sync · Geteiltes PostgreSQL · Redundanz-Node',
'bootstrap-title':'Bootstrap-Node — Dem Netzwerk beitreten','bootstrap-desc':'Um deinen eigenen Aequitas-Node zu betreiben, verbinde dich mit dem Bootstrap-Node über die unten stehende libp2p-Multiadresse. Dein Node wird automatisch Peers entdecken, die vollständige Blockhistorie herunterladen und am Konsens teilnehmen.','bootstrap-howto':'Um einen Node zu betreiben: Repository klonen, DATABASE_URL-Umgebungsvariable setzen und Binary mit Bootstrap-Peer-Adresse ausführen. Vollständige Dokumentation auf GitHub verfügbar.',
'tech-specs':'Technische Spezifikationen','chain-id':'Chain-ID','evm-yes':'Ja — JSON-RPC unter /rpc · MetaMask-kompatibel','language':'Sprache','source':'Quellcode',
'metamask-config':'MetaMask / Web3-Konfiguration','metamask-desc':'Aequitas Chain ist EVM-kompatibel und kann zu MetaMask oder jeder anderen Web3-Wallet hinzugefügt werden. Verwende die untenstehende Konfiguration um Aequitas zu deiner Wallet hinzuzufügen und dein AEQ-Guthaben zu sehen.',
'architecture-title':'Systemarchitektur-Übersicht','architecture-desc':'Das Aequitas-System besteht aus vier Hauptkomponenten. Die Android-App verarbeitet biometrische Scans und Beweiserzeugung vollständig auf dem Gerät. Der Proof Server (Node.js) erzeugt Groth16-ZK-Beweise und speichert biometrische Hashes. Die Blockchain-Nodes (Go) verwalten den BlockDAG, verarbeiten Transaktionen und stellen den EVM-RPC bereit. PostgreSQL speichert dauerhaften Zustand, der über alle Nodes geteilt wird.',
'reg-title':'🔐 Als verifizierter Mensch registrieren','reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte deine 1.000 AEQ. Eine einmalige, permanente, gebührenfreie Registrierung, die kryptographisch beweist, dass du ein einzigartiger Mensch bist. Keine persönlichen Daten werden gesammelt oder gespeichert. Keine Gasgebühren. Kein Warten.',
'app-only-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP','app-only-text':'Der Menschlichkeitsnachweis erfordert biometrische Verifizierung auf deinem persönlichen Gerät. Der Fingerabdruck muss durch das Hardware Secure Element (HSE) deines Telefons verarbeitet werden – ein dedizierter kryptographischer Chip, der vom Hauptprozessor isoliert und nicht fernzugreifbar ist. Das HSE ist derselbe Chip, der deine Banking-Apps und das Sperrbildschirm deines Telefons sichert. Wenn du in der Aequitas-App scannst, wird die Signatur vom HSE verarbeitet, gehasht und zur Erzeugung eines Zero-Knowledge-Beweises verwendet – alles ohne dass die rohen Fingerabdruckdaten jemals das HSE verlassen. Lade die App herunter, scanne deinen Fingerabdruck, verbinde deine Wallet, und deine <strong style="color:var(--gold)">1.000 AEQ werden automatisch und sofort gewährt</strong>.',
'step1-title':'Biometrischer Scan','step1-desc':'Öffne die Aequitas-App und tippe auf "Menschlichkeit beweisen." Scanne deinen Fingerabdruck – vollständig durch das Hardware Secure Element deines Telefons verarbeitet. Rohe Fingerabdruckdaten erreichen niemals den Hauptprozess der App und verlassen niemals dein Gerät. Das HSE erzeugt eine kryptographische Signatur als deinen einzigartigen biometrischen Identifikator.',
'step2-title':'ZKP-Erzeugung','step2-desc':'Die App sendet deinen biometrischen Hash an den Proof Server. Der Server prüft, ob dieser Hash noch nicht verwendet wurde (verhindert Doppelregistrierung), und erzeugt dann einen Groth16-Zero-Knowledge-Beweis über BN128. Dieser Beweis garantiert mathematisch "Ein einzigartiger biometrischer Hash wurde verwendet" ohne irgendetwas über den Hash zu enthüllen. Dauert 2-5 Sekunden.',
'step3-title':'Wallet verbinden','step3-desc':'Die App öffnet MetaMask mit der vorkonfigurierten Aequitas Chain. Verbinde deine Wallet – diese Adresse erhält deine 1.000 AEQ. Nach der Registrierung ist diese Wallet dauerhaft mit deiner biometrischen Identität verknüpft. Du kannst sie nicht ohne erneute Registrierung ändern.',
'step4-title':'1.000 AEQ gewährt','step4-desc':'Die Registrierungstransaktion wird an die Blockchain übermittelt. Das Protokoll verifiziert den ZK-Beweis, bestätigt, dass der biometrische Hash einzigartig ist, und schreibt genau 1.000 AEQ deiner Wallet innerhalb des nächsten Blocks (unter 6 Sekunden) gut. Gesamtzeit vom Fingerabdruck bis zum AEQ: unter 30 Sekunden.',
'priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Biometrische Daten verlassen niemals das Gerät · Keine persönlichen Daten · Keine Gasgebühren · Permanenter Sybil-Schutz',
'connected-wallet':'VERBUNDENE WALLET','proof-detected':'⚡ BEWEISPARAMETER VON APP ERKANNT','connect-btn':'🦊 METAMASK VERBINDEN','register-btn':'🔐 ON-CHAIN REGISTRIEREN','reg-hint':'// Öffne die Aequitas Android-App um deinen Beweis zu generieren, kehre dann hierher zurück um die Registrierung abzuschließen...','reg-details':'Registrierungsdetails','network':'Netzwerk','reg-limit':'Registrierungen','reg-limit-val':'Einmalig pro Mensch · dauerhaft · unveränderlich','bio-data':'Biometrische Daten','never-stored':'Niemals gespeichert','confirmation':'Bestätigungszeit','conf-time':'Innerhalb von 6 Sekunden (nächster Block)',
'phase-desc-0':'Phase 0: Bootstrap — Netzwerk aufbauen und frühe Menschen onboarden','phase-desc-1':'Phase 1: Wachstum — Menschenregister global erweitern','phase-desc-2':'Phase 2: Stabilität — Vermögensumverteilungsmechanismen aktiv','phase-desc-3':'Phase 3: Reife — Vollständige Dezentralisierung erreicht'
},
es:{
'live':'EN VIVO','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-register':'🔐 Registrar',
'live-stats':'Estadísticas de Cadena en Vivo','block-height':'Altura de Bloque','block-height-sub':'Nuevo bloque cada 6 segundos · Consenso BlockDAG · Dos nodos produciendo bloques en paralelo','verified-humans':'Humanos Verificados','verified-humans-sub':'Cada wallet verificada como humano único · ZKP biométrico · Una persona, una wallet, para siempre','total-supply':'Suministro Total','total-supply-sub':'Siempre igual a Humanos × 1,000 AEQ · El suministro crece solo cuando crece la humanidad','aeq-index':'Índice Aequitas','aeq-index-sub':'0 = igualdad perfecta · 100 = desigualdad máxima · Basado en coeficiente Gini real','uptime':'Tiempo Activo','uptime-sub':'Node v0.3.0 · 2 nodos activos · Railway + Render · Estado persistente PostgreSQL',
'poh-title':'Prueba de Humanidad','poh-text':'Cada titular de AEQ debe probar que es un ser humano único y vivo mediante verificación biométrica. Esto no es opcional — es el fundamento de todo el sistema. Sin prueba de humanidad, no se puede recibir AEQ. Esto significa que ningún bot, cuenta duplicada, corporación, gobierno o sistema de IA puede tener AEQ. Solo humanos reales. La verificación usa tu huella digital a través del Elemento Seguro de Hardware de tu teléfono — el mismo chip que protege tus apps bancarias. Tus datos biométricos nunca salen de tu dispositivo.',
'fair-title':'Distribución Radicalmente Justa','fair-text':'Cada humano verificado en la Tierra recibe exactamente 1,000 AEQ — ni más, ni menos. La primera persona en registrarse y la persona número mil millones reciben cantidades idénticas. No hay pre-minado, asignación a fundadores, ronda de inversores ni ventaja para los primeros adoptantes. El suministro total de AEQ siempre es igual a humanos verificados × 1,000. Este es el sistema de distribución monetaria más igualitario jamás diseñado.',
'dag-title':'Arquitectura BlockDAG','dag-text':'Aequitas usa un Grafo Acíclico Dirigido (DAG) donde múltiples bloques pueden producirse simultáneamente por diferentes nodos y fusionarse después. Esto permite mayor rendimiento, menor latencia y mejor tolerancia a fallos. Cuando dos nodos producen bloques al mismo tiempo, ambos son válidos y se fusionan — puedes ver estos eventos marcados con 🔀 en el explorador.',
'gasless-title':'Registro Verdaderamente Sin Gas','gasless-text':'Una de las mayores barreras para la adopción de criptomonedas es requerir que ya tengas criptomonedas para pagar comisiones. Aequitas elimina esto completamente. El registro no cuesta absolutamente nada — no se necesita ETH, BNB ni MATIC. Sin tarjeta de crédito, sin cuenta bancaria. Si eres humano con smartphone, puedes registrarte.',
'recent-blocks':'Bloques Recientes','blocks-desc':'Cada bloque está vinculado criptográficamente a sus padres mediante hashes SHA-256. 🔀 MERGE = bloque con múltiples padres (función BlockDAG). ✅ TX = bloque con transacción de registro. Tiempo de bloque: 6 segundos promedio.','loading':'Cargando bloques...','network-info':'Información de Red','chain-name':'Nombre de Red','symbol':'Símbolo','block-time':'Tiempo de Bloque','consensus':'Consenso','nodes':'Nodos Activos','storage':'Almacenamiento','add-metamask':'🦊 AGREGAR A METAMASK','network-name':'Nombre de Red','add-network':'+ AGREGAR RED AEQUITAS','philosophy':'"El dinero existe porque las personas existen.<br>Nada más, nada menos."','philosophy-sub':'— EL PRINCIPIO AEQUITAS —','decimals':'Decimales',
'what-is-it':'¿Qué significa Humano Verificado?','humans-what':'Un Humano Verificado en Aequitas es una dirección wallet demostrada criptográficamente que pertenece a un humano único y vivo. La verificación usa datos biométricos — específicamente una huella digital — procesados por el Elemento Seguro de Hardware de un smartphone Android. Los datos biométricos nunca se transmiten ni almacenan.','how-works':'El Sistema de Prueba de Conocimiento Cero','humans-how':'Aequitas usa el sistema de demostración Groth16 — también usado por Zcash. Tu firma de huella digital se hashea en un elemento de campo de la curva elíptica BN128. Este hash se usa como entrada para un circuito Groth16 que produce una pequeña prueba matemática garantizando "se usó un hash biométrico único" sin revelar qué es el hash.',
'sybil-title':'Prevención Permanente de Ataques Sybil','sybil-text':'Un ataque Sybil es cuando una entidad crea múltiples identidades falsas. Es la debilidad fundamental de casi toda blockchain. Aequitas resuelve esto en la capa de identidad: cada hash biométrico se almacena permanentemente. Intentar registrarse dos veces con la misma huella se rechaza inmediatamente. Una persona, una wallet, para siempre — garantizado por criptografía, no por confianza.',
'global-title':'Diseñado para Inclusión Global','global-text':'Aequitas está diseñado para ser accesible a todo humano en la Tierra. No necesitas cuenta bancaria, tarjeta de crédito ni criptomoneda existente. Solo un smartphone Android con sensor de huella — un dispositivo que ya poseen más de 3 mil millones de personas. El registro es completamente gratuito, toma menos de 2 minutos y otorga 1,000 AEQ instantáneamente.',
'registered-humans':'Humanos Registrados','humans-desc':'Cada dirección aquí listada fue verificada como humano único mediante Prueba de Conocimiento Cero biométrica. Cada uno recibió exactamente 1,000 AEQ al registrarse. Estas registraciones son permanentes, inmutables y almacenadas tanto en PostgreSQL como on-chain.','no-humans':'No hay humanos registrados aún.\n\n¡Descarga la App Android Aequitas y sé el primero en la cadena!','registry-stats':'Estadísticas del Registro','total-humans-stat':'Total de Humanos','grant':'Bono por Humano','reg-fee':'Tarifa de Registro','free':'GRATIS','zkp-system':'Sistema ZKP','never-stored':'Nunca almacenado','zkp-title':'Detalles Técnicos ZKP','zkp-details':'Groth16 sobre curva elíptica BN128. Tamaño de prueba: ~200 bytes. Tiempo de verificación: ~10ms. Circuito compilado con snarkjs y circom.',
'aeq-index-title':'Índice Aequitas — Puntuación de Igualdad Económica en Tiempo Real','aeq-index-desc':'El Índice Aequitas mide la salud económica de la red en tiempo real. Calculado directamente desde la distribución de saldos on-chain de todos los humanos verificados. Puntuación 0 = cada humano tiene AEQ idéntico. Puntuación 100 = una persona controla todo. Basado en el coeficiente Gini real de saldos on-chain.','current-index':'Índice Actual','bar-0':'0 — Igualdad Perfecta','bar-100':'100 — Máx. Desigualdad','gini-coeff':'Coeficiente Gini','phase':'Fase del Protocolo',
'pools-title':'Pools de Redistribución','pools-desc':'Cuando se superan los umbrales de desigualdad, AEQ se redirige automáticamente hacia estos cuatro pools. Son cuentas de contratos inteligentes controladas completamente por lógica de protocolo.','vel-pool':'Pool Velocidad','liq-pool':'Pool Liquidez','uni-pool':'Pool Unidad','treasury':'Tesorería','pools-details':'El Pool de Velocidad recompensa wallets que transaccionan regularmente. El Pool de Liquidez apoya profundidad de mercado y estabilidad de precios. El Pool de Unidad otorga AEQ adicional a humanos recién registrados en fases posteriores. La Tesorería financia desarrollo del protocolo.',
'phases-title':'Fases del Protocolo','phases-desc':'El protocolo Aequitas evoluciona a través de cuatro fases distintas. Las transiciones de fase ocurren automáticamente basándose en humanos verificados y coeficiente Gini — no se requiere votación.','phase0':'Bootstrap — Construyendo la red · &lt;100 humanos','phase1':'Crecimiento — Expandiendo el registro · 100–10,000 humanos','phase2':'Estabilidad — Redistribución activa · 10,000–1M humanos','phase3':'Madurez — Descentralización completa · 1M+ humanos · Gini &lt;0.3',
'gini-title':'El Coeficiente Gini — Cómo Medimos la Desigualdad','gini-desc':'El coeficiente Gini fue desarrollado por el estadístico italiano Corrado Gini en 1912 y sigue siendo la medida de desigualdad económica más utilizada del mundo. Usado por el Banco Mundial, el FMI y virtualmente todos los economistas. Aequitas lo calcula en tiempo real desde saldos on-chain reales.',
'gini-0':'Igualdad Perfecta','gini-0-sub':'Todos tienen riqueza idéntica. Aequitas se aproxima a esto en Fase 0 cuando todos tienen exactamente 1,000 AEQ.','gini-1':'Baja Desigualdad','gini-1-sub':'Promedio Escandinavia (Suecia 0.27, Dinamarca 0.28). La zona objetivo de Aequitas.','gini-2':'Desigualdad Moderada','gini-2-sub':'Promedio EE.UU. (0.41). Los mecanismos de redistribución se activan en este nivel en Fase 2.','gini-3':'Alta Desigualdad','gini-3-sub':'Sudáfrica (0.63). El Gini estimado de Bitcoin supera 0.85 — más alto que cualquier país.',
'inflation-title':'Mecanismo de Inflación y Redistribución — En Detalle',
'story-title':'La Historia de Aequitas — Por Qué Existe',
'active-nodes':'Nodos Activos — Topología de Red Actual','nodes-desc':'La red Aequitas opera en dos nodos en entornos cloud distribuidos geográficamente. Ambos participan en producción de bloques, sincronización de estado y servicio de API. Se comunican via libp2p y sincronizan bloques via HTTP. Ambos comparten la misma base de datos PostgreSQL.','node1':'Nodo 1 — Railway (Primario)','node2':'Nodo 2 — Render (Secundario)','node1-desc':'Servidor API primario · Productor de bloques · Nodo bootstrap P2P · Conexión PostgreSQL · Endpoint RPC para MetaMask','node2-desc':'Servidor API secundario · Productor de bloques · Par P2P · Sincronización HTTP de bloques · PostgreSQL compartido · Nodo de redundancia',
'bootstrap-title':'Nodo Bootstrap — Unirse a la Red','bootstrap-desc':'Para ejecutar tu propio nodo Aequitas, conéctate al nodo bootstrap usando la multidirección libp2p. Tu nodo descubrirá automáticamente pares, descargará el historial completo de bloques y comenzará a participar en consenso.','bootstrap-howto':'Para ejecutar un nodo: clona el repositorio GitHub, establece la variable de entorno DATABASE_URL y ejecuta el binario con la dirección del par bootstrap.',
'tech-specs':'Especificaciones Técnicas','chain-id':'ID de Cadena','evm-yes':'Sí — JSON-RPC en /rpc · Compatible con MetaMask','language':'Lenguaje','source':'Código Fuente',
'metamask-config':'Configuración MetaMask / Web3','metamask-desc':'Aequitas Chain es compatible con EVM y puede agregarse a MetaMask o cualquier wallet Web3 que soporte redes RPC personalizadas.',
'architecture-title':'Resumen de Arquitectura del Sistema','architecture-desc':'El sistema Aequitas tiene cuatro componentes principales. La App Android maneja escaneo biométrico y generación de pruebas completamente en el dispositivo. El Servidor de Pruebas (Node.js) genera pruebas ZK Groth16 y almacena hashes biométricos. Los Nodos Blockchain (Go) mantienen el BlockDAG y exponen el RPC EVM. PostgreSQL almacena estado persistente compartido entre todos los nodos.',
'reg-title':'🔐 Regístrate como Humano Verificado','reg-sub':'Únete a la red Aequitas y recibe tus 1,000 AEQ. Un registro único, permanente y sin gas que prueba criptográficamente que eres un humano único. No se recopilan ni almacenan datos personales.',
'app-only-title':'REGISTRO SOLO VÍA APP ANDROID','app-only-text':'La Prueba de Humanidad requiere verificación biométrica en tu dispositivo personal. La huella debe procesarse por el Elemento Seguro de Hardware (HSE) de tu teléfono — un chip criptográfico dedicado aislado del procesador principal. El HSE es el mismo chip que protege tus apps bancarias. Cuando escaneas en la App Aequitas, la firma es procesada por el HSE, hasheada y usada para generar una Prueba de Conocimiento Cero — todo sin que los datos de huella brutos salgan del HSE. Descarga la app, escanea tu huella, conecta tu wallet, y tus <strong style="color:var(--gold)">1,000 AEQ serán otorgados automática e inmediatamente</strong>.',
'step1-title':'Escaneo Biométrico','step1-desc':'Abre la App Aequitas y toca "Probar Humanidad." Escanea tu huella digital — procesada completamente por el Elemento Seguro de Hardware de tu teléfono. Los datos brutos de huella nunca alcanzan el proceso principal de la app ni salen de tu dispositivo. El HSE genera una firma criptográfica como tu identificador biométrico único.',
'step2-title':'Generación de ZKP','step2-desc':'La app envía tu hash biométrico al Servidor de Pruebas. El servidor verifica que este hash no se haya usado antes (previene doble registro), luego genera una Prueba de Conocimiento Cero Groth16 sobre BN128. Esta prueba garantiza "se usó un hash biométrico único" sin revelar nada sobre el hash. Toma 2-5 segundos.',
'step3-title':'Conectar Wallet','step3-desc':'La app abre MetaMask con Aequitas Chain preconfigurada. Conecta tu wallet — esta dirección recibirá tus 1,000 AEQ. Una vez registrado, esta wallet queda permanentemente vinculada a tu identidad biométrica. No puedes cambiarla sin re-registrarte.',
'step4-title':'1,000 AEQ Otorgados','step4-desc':'La transacción de registro se envía a la blockchain. El protocolo verifica la prueba ZK, confirma que el hash biométrico es único y acredita exactamente 1,000 AEQ a tu wallet dentro del siguiente bloque (menos de 6 segundos). Tiempo total desde huella hasta AEQ: menos de 30 segundos.',
'priv-bar':'🔒 Elemento Seguro de Hardware · ZKP Groth16 · Datos biométricos nunca salen del dispositivo · Sin datos personales · Sin tarifas de gas · Protección Sybil permanente',
'connected-wallet':'WALLET CONECTADA','proof-detected':'⚡ PARÁMETROS DE PRUEBA DETECTADOS','connect-btn':'🦊 CONECTAR METAMASK','register-btn':'🔐 REGISTRAR ON-CHAIN','reg-hint':'// Abre la App Android Aequitas para generar tu prueba, luego regresa aquí para completar el registro...','reg-details':'Detalles del Registro','network':'Red','reg-limit':'Registros','reg-limit-val':'Una vez por humano · permanente · inmutable','bio-data':'Datos Biométricos','never-stored':'Nunca almacenados','confirmation':'Tiempo de Confirmación','conf-time':'Dentro de 6 segundos (próximo bloque)',
'phase-desc-0':'Fase 0: Bootstrap — Construyendo la red','phase-desc-1':'Fase 1: Crecimiento — Expandiendo el registro global','phase-desc-2':'Fase 2: Estabilidad — Mecanismos de redistribución activos','phase-desc-3':'Fase 3: Madurez — Descentralización completa alcanzada'
},
ru:{
'live':'В ЭФИРЕ','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-register':'🔐 Регистрация',
'live-stats':'Статистика цепочки в реальном времени','block-height':'Высота Блока','block-height-sub':'Новый блок каждые 6 секунд · Консенсус BlockDAG · Два узла производят блоки параллельно','verified-humans':'Верифицированных Людей','verified-humans-sub':'Каждый кошелёк верифицирован как уникальный человек · Биометрический ZKP · Один человек, один кошелёк, навсегда','total-supply':'Общее Предложение','total-supply-sub':'Всегда равно Людям × 1 000 AEQ · Предложение растёт только вместе с человечеством','aeq-index':'Индекс Aequitas','aeq-index-sub':'0 = полное равенство · 100 = максимальное неравенство · На основе реального коэффициента Джини','uptime':'Время Работы','uptime-sub':'Node v0.3.0 · 2 активных узла · Railway + Render · Постоянное состояние PostgreSQL',
'poh-title':'Доказательство Человечности','poh-text':'Каждый владелец AEQ должен доказать, что он уникальный живой человек — через биометрическую верификацию. Это не опционально — это фундамент всей системы. Без доказательства человечности AEQ не может быть получен. Это означает: никаких ботов, дублирующих аккаунтов, корпораций, правительств или систем ИИ не может держать AEQ. Только настоящие люди. Верификация использует ваш отпечаток пальца через Hardware Secure Element вашего телефона. Ваши биометрические данные никогда не покидают ваше устройство.',
'fair-title':'Радикально Справедливое Распределение','fair-text':'Каждый верифицированный человек на Земле получает ровно 1 000 AEQ — не больше и не меньше. Первый зарегистрировавшийся и миллиардный получают одинаковые суммы. Нет предварительной добычи, распределения основателям, инвесторских раундов или преимуществ первых участников. Общее предложение AEQ всегда равно верифицированным людям × 1 000. Это самая эгалитарная система денежного распределения, когда-либо созданная.',
'dag-title':'Архитектура BlockDAG','dag-text':'Aequitas использует направленный ациклический граф (DAG), где несколько блоков могут производиться одновременно разными узлами и позже объединяться. Это обеспечивает значительно более высокую пропускную способность, меньшую задержку и лучшую отказоустойчивость. Эти события слияния отмечены 🔀 в проводнике.',
'gasless-title':'По-Настоящему Бесплатная Регистрация','gasless-text':'Одним из главных барьеров для принятия криптовалюты является требование уже иметь криптовалюту для оплаты комиссий. Aequitas полностью устраняет это. Регистрация не стоит абсолютно ничего — не нужен ETH, BNB или MATIC. Нет кредитной карты, нет банковского счёта. Если ты человек со смартфоном, ты можешь зарегистрироваться.',
'recent-blocks':'Последние Блоки','blocks-desc':'Каждый блок криптографически связан со своими родителями через хэши SHA-256. 🔀 MERGE = блок с несколькими родителями (функция BlockDAG). ✅ TX = блок с транзакцией регистрации. Время блока: в среднем 6 секунд.','loading':'Загрузка блоков...','network-info':'Информация о Сети','chain-name':'Название Сети','symbol':'Символ','block-time':'Время Блока','consensus':'Консенсус','nodes':'Активные Узлы','storage':'Хранение Состояния','add-metamask':'🦊 ДОБАВИТЬ В METAMASK','network-name':'Название Сети','add-network':'+ ДОБАВИТЬ СЕТЬ AEQUITAS','philosophy':'"Деньги существуют потому что существуют люди.<br>Ничего больше, ничего меньше."','philosophy-sub':'— ПРИНЦИП AEQUITAS —','decimals':'Знаков после запятой',
'what-is-it':'Что означает Верифицированный Человек?','humans-what':'Верифицированный Человек в Aequitas — это адрес кошелька, криптографически доказанно принадлежащий уникальному живому человеку. Верификация использует биометрические данные — отпечаток пальца — обработанные через Hardware Secure Element Android-смартфона. Сами биометрические данные никогда не передаются и не хранятся.','how-works':'Система Доказательства с Нулевым Разглашением','humans-how':'Aequitas использует систему доказательства Groth16 — также используемую Zcash. Ваша подпись отпечатка хэшируется в элемент поля эллиптической кривой BN128. Этот хэш используется как вход для схемы Groth16, которая производит маленькое математическое доказательство: "был использован уникальный биометрический хэш" без раскрытия самого хэша.',
'sybil-title':'Постоянная Защита от Атак Сивиллы','sybil-text':'Атака Сивиллы — когда одна сущность создаёт множество поддельных идентичностей. Это фундаментальная слабость почти каждого блокчейна. Aequitas решает это на уровне идентичности: каждый биометрический хэш хранится постоянно. Попытка зарегистрироваться дважды с тем же отпечатком немедленно отклоняется. Один человек, один кошелёк, навсегда — гарантировано криптографией, не доверием.',
'global-title':'Создан для Глобального Включения','global-text':'Aequitas создан доступным для каждого человека на Земле. Не нужен банковский счёт, кредитная карта или существующая криптовалюта. Только Android-смартфон с сенсором отпечатка — устройство, которым уже владеют более 3 миллиардов людей. Регистрация бесплатна, занимает менее 2 минут и мгновенно даёт 1 000 AEQ.',
'registered-humans':'Зарегистрированных Людей','humans-desc':'Каждый адрес здесь верифицирован как уникальный человек через биометрическое Доказательство с Нулевым Разглашением. Каждый получил ровно 1 000 AEQ при регистрации. Эти регистрации постоянны, неизменны и хранятся как в PostgreSQL, так и on-chain.','no-humans':'Людей ещё нет.\n\nСкачай приложение Aequitas Android и стань первым человеком в цепочке!','registry-stats':'Статистика Реестра','total-humans-stat':'Всего Людей','grant':'Грант на Человека','reg-fee':'Плата за Регистрацию','free':'БЕСПЛАТНО','zkp-system':'Система ZKP','never-stored':'Никогда не хранится','zkp-title':'Технические Детали ZKP','zkp-details':'Groth16 над эллиптической кривой BN128. Размер доказательства: ~200 байт. Время верификации: ~10мс. Схема скомпилирована с snarkjs и circom.',
'aeq-index-title':'Индекс Aequitas — Оценка Экономического Равенства в Реальном Времени','aeq-index-desc':'Индекс Aequitas измеряет экономическое здоровье сети в реальном времени. Рассчитывается непосредственно из распределения балансов on-chain всех верифицированных людей. Оценка 0 = у каждого человека идентичный AEQ. Оценка 100 = одна персона всем управляет. Основан на реальном коэффициенте Джини балансов on-chain.','current-index':'Текущий Индекс','bar-0':'0 — Полное Равенство','bar-100':'100 — Макс. Неравенство','gini-coeff':'Коэффициент Джини','phase':'Фаза Протокола',
'pools-title':'Пулы Перераспределения','pools-desc':'Когда пороги неравенства превышены, AEQ автоматически перенаправляется в эти четыре пула. Пулы — это счета смарт-контрактов, контролируемые исключительно логикой протокола.','vel-pool':'Пул Скорости','liq-pool':'Пул Ликвидности','uni-pool':'Пул Единства','treasury':'Казначейство','pools-details':'Пул Скорости вознаграждает кошельки, которые регулярно совершают транзакции. Пул Ликвидности поддерживает глубину рынка и стабильность цен. Пул Единства предоставляет дополнительный AEQ новым людям в поздних фазах. Казначейство финансирует разработку протокола.',
'phases-title':'Фазы Протокола','phases-desc':'Протокол Aequitas развивается через четыре фазы по мере роста сети. Переходы между фазами происходят автоматически на основе верифицированных людей и коэффициента Джини — голосование не требуется.','phase0':'Загрузка — Построение сети · &lt;100 людей','phase1':'Рост — Расширение реестра · 100–10 000 людей','phase2':'Стабильность — Перераспределение активно · 10 000–1М людей','phase3':'Зрелость — Полная децентрализация · 1М+ людей · Джини &lt;0,3',
'gini-title':'Коэффициент Джини — Как Мы Измеряем Неравенство','gini-desc':'Коэффициент Джини был разработан итальянским статистиком Коррадо Джини в 1912 году и остаётся наиболее широко используемой мерой экономического неравенства. Используется Всемирным банком, МВФ и практически всеми экономистами. Aequitas рассчитывает его в реальном времени из фактических балансов on-chain.',
'gini-0':'Полное Равенство','gini-0-sub':'У всех одинаковое состояние. Aequitas приближается к этому в Фазе 0, когда у всех людей ровно 1 000 AEQ.','gini-1':'Низкое Неравенство','gini-1-sub':'Средн. Скандинавия (Швеция 0,27, Дания 0,28). Целевая зона Aequitas для долгосрочной работы.','gini-2':'Умеренное Неравенство','gini-2-sub':'Средн. США (0,41). Механизмы перераспределения активируются на этом уровне в Фазе 2.','gini-3':'Высокое Неравенство','gini-3-sub':'Южная Африка (0,63). Оценочный Джини Биткоина превышает 0,85 — выше любой страны Земли.',
'inflation-title':'Механизм Инфляции и Перераспределения — Подробно',
'story-title':'История Aequitas — Почему Это Существует',
'active-nodes':'Активные Узлы — Текущая Топология Сети','nodes-desc':'Сеть Aequitas работает на двух узлах в географически распределённых облачных средах. Оба участвуют в производстве блоков, синхронизации состояния и обслуживании API. Они общаются через libp2p и синхронизируют блоки через HTTP. Оба узла используют одну базу данных PostgreSQL.','node1':'Узел 1 — Railway (Основной)','node2':'Узел 2 — Render (Вторичный)','node1-desc':'Основной API-сервер · Производитель блоков · P2P-bootstrap-узел · Подключение PostgreSQL · RPC-эндпоинт для MetaMask','node2-desc':'Вторичный API-сервер · Производитель блоков · P2P-пир · HTTP-синхронизация блоков · Общий PostgreSQL · Резервный узел',
'bootstrap-title':'Bootstrap-Узел — Присоединиться к Сети','bootstrap-desc':'Чтобы запустить собственный узел Aequitas, подключитесь к bootstrap-узлу, используя мультиадрес libp2p. Ваш узел автоматически обнаружит пиров, загрузит полную историю блоков и начнёт участвовать в консенсусе.','bootstrap-howto':'Для запуска узла: клонируйте репозиторий GitHub, установите переменную среды DATABASE_URL и запустите бинарный файл с адресом bootstrap-пира.',
'tech-specs':'Технические Характеристики','chain-id':'ID Цепочки','evm-yes':'Да — JSON-RPC по /rpc · Совместим с MetaMask','language':'Язык','source':'Исходный Код',
'metamask-config':'Настройка MetaMask / Web3','metamask-desc':'Aequitas Chain совместима с EVM и может быть добавлена в MetaMask или любой Web3-кошелёк, поддерживающий пользовательские RPC-сети.',
'architecture-title':'Обзор Системной Архитектуры','architecture-desc':'Система Aequitas имеет четыре основных компонента. Android-приложение обрабатывает биометрическое сканирование и генерацию доказательств полностью на устройстве. Сервер Доказательств (Node.js) генерирует ZK-доказательства Groth16 и хранит биометрические хэши. Узлы Блокчейна (Go) поддерживают BlockDAG и предоставляют EVM RPC. PostgreSQL хранит постоянное состояние, общее для всех узлов.',
'reg-title':'🔐 Зарегистрируйтесь как Верифицированный Человек','reg-sub':'Присоединитесь к сети Aequitas и получите 1 000 AEQ. Одноразовая, постоянная, бесплатная регистрация, криптографически доказывающая, что вы уникальный человек. Никакие личные данные не собираются и не хранятся.',
'app-only-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ','app-only-text':'Доказательство человечности требует биометрической верификации на вашем личном устройстве. Отпечаток должен обрабатываться Hardware Secure Element (HSE) вашего телефона — выделенным криптографическим чипом, изолированным от основного процессора. HSE — тот же чип, который защищает ваши банковские приложения. Когда вы сканируете в приложении Aequitas, подпись обрабатывается HSE, хэшируется и используется для генерации Доказательства с Нулевым Разглашением — всё без выхода сырых данных отпечатка за пределы HSE. Скачайте приложение, отсканируйте отпечаток, подключите кошелёк, и ваши <strong style="color:var(--gold)">1 000 AEQ будут начислены автоматически и немедленно</strong>.',
'step1-title':'Биометрический Скан','step1-desc':'Откройте приложение Aequitas и нажмите "Доказать Человечность." Отсканируйте отпечаток пальца — полностью обрабатывается Hardware Secure Element. Сырые данные отпечатка никогда не попадают в основной процесс приложения и не покидают ваше устройство. HSE генерирует криптографическую подпись как ваш уникальный биометрический идентификатор.',
'step2-title':'Генерация ZKP','step2-desc':'Приложение отправляет ваш биометрический хэш на Сервер Доказательств. Сервер проверяет, что этот хэш не использовался ранее, затем генерирует Groth16 Доказательство с Нулевым Разглашением над BN128. Это доказательство гарантирует "был использован уникальный биометрический хэш" без раскрытия чего-либо о хэше. Занимает 2-5 секунд.',
'step3-title':'Подключить Кошелёк','step3-desc':'Приложение открывает MetaMask с предварительно настроенной Aequitas Chain. Подключите кошелёк — этот адрес получит 1 000 AEQ. После регистрации кошелёк навсегда связан с вашей биометрической идентичностью.',
'step4-title':'1 000 AEQ Начислено','step4-desc':'Транзакция регистрации отправляется в блокчейн. Протокол верифицирует ZK-доказательство, подтверждает уникальность биометрического хэша и зачисляет ровно 1 000 AEQ на ваш кошелёк в течение следующего блока (менее 6 секунд). Общее время от отпечатка до AEQ: менее 30 секунд.',
'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Биометрические данные никогда не покидают устройство · Без личных данных · Без комиссий · Постоянная защита',
'connected-wallet':'ПОДКЛЮЧЁННЫЙ КОШЕЛЁК','proof-detected':'⚡ ПАРАМЕТРЫ ДОКАЗАТЕЛЬСТВА ОБНАРУЖЕНЫ','connect-btn':'🦊 ПОДКЛЮЧИТЬ METAMASK','register-btn':'🔐 ЗАРЕГИСТРИРОВАТЬСЯ ON-CHAIN','reg-hint':'// Откройте Android-приложение Aequitas для генерации доказательства, затем вернитесь сюда...','reg-details':'Детали Регистрации','network':'Сеть','reg-limit':'Регистрации','reg-limit-val':'Один раз на человека · постоянно · неизменно','bio-data':'Биометрические Данные','never-stored':'Никогда не хранятся','confirmation':'Время Подтверждения','conf-time':'В течение 6 секунд (следующий блок)',
'phase-desc-0':'Фаза 0: Загрузка — Построение сети','phase-desc-1':'Фаза 1: Рост — Глобальное расширение реестра','phase-desc-2':'Фаза 2: Стабильность — Механизмы перераспределения активны','phase-desc-3':'Фаза 3: Зрелость — Полная децентрализация достигнута'
},
zh:{
'live':'直播','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-register':'🔐 注册',
'live-stats':'链上实时统计','block-height':'区块高度','block-height-sub':'每6秒新区块 · BlockDAG共识 · 两个节点并行生产区块','verified-humans':'已验证人类','verified-humans-sub':'每个钱包验证为唯一人类 · 生物特征ZKP · 一人一钱包，永久','total-supply':'总供应量','total-supply-sub':'始终等于人类 × 1,000 AEQ · 供应量仅随人类增长','aeq-index':'Aequitas指数','aeq-index-sub':'0 = 完全平等 · 100 = 最大不平等 · 基于真实基尼系数','uptime':'运行时间','uptime-sub':'Node v0.3.0 · 2个活跃节点 · Railway + Render · PostgreSQL持久状态',
'poh-title':'人类证明','poh-text':'每个AEQ持有者必须通过生物特征验证证明自己是唯一的活人。这不是可选的——这是整个系统的基础。没有人类证明，就无法接收AEQ。这意味着没有机器人、重复账户、企业、政府或AI系统可以持有AEQ。只有真实的人类。验证使用您手机硬件安全元件中的指纹——保护您银行应用程序的同一芯片。您的生物特征数据在任何情况下都不会离开您的设备。',
'fair-title':'根本公平的分配','fair-text':'地球上每个经过验证的人类获得恰好1,000 AEQ——不多不少。第一个注册的人和第十亿个注册的人获得相同的金额。没有预挖，没有创始人分配，没有投资者轮次，没有早期采用者优势。AEQ的总供应量始终永久等于已验证人类 × 1,000。这是有史以来设计最平等的货币分配系统。',
'dag-title':'BlockDAG架构','dag-text':'Aequitas使用有向无环图（DAG），不同节点可以同时生产多个区块，之后进行合并。这允许更高的吞吐量、更低的延迟和更好的容错性。当两个节点同时生产区块时，两者都有效并在之后合并——您可以在浏览器中看到这些合并事件标有🔀。',
'gasless-title':'真正无Gas费注册','gasless-text':'加密货币采用的最大障碍之一是要求您已经拥有加密货币来支付交易费用。Aequitas完全消除了这一点。注册绝对免费——不需要ETH、BNB或MATIC。不需要信用卡，不需要银行账户。如果您是拥有智能手机的人类，您就可以注册。',
'recent-blocks':'最近区块','blocks-desc':'每个区块通过SHA-256哈希与其父区块进行加密链接。🔀 MERGE = 具有多个父区块的区块（BlockDAG功能）。✅ TX = 包含注册交易的区块。区块时间：平均6秒。','loading':'加载区块中...','network-info':'网络信息','chain-name':'网络名称','symbol':'符号','block-time':'出块时间','consensus':'共识','nodes':'活跃节点','storage':'状态存储','add-metamask':'🦊 添加到METAMASK','network-name':'网络名称','add-network':'+ 添加AEQUITAS网络','philosophy':'"货币存在是因为人类存在。<br>仅此而已，不多也不少。"','philosophy-sub':'— AEQUITAS原则 —','decimals':'小数位',
'what-is-it':'已验证人类是什么意思？','humans-what':'Aequitas上的已验证人类是一个加密证明属于独特活人的钱包地址。验证使用生物特征数据——具体是指纹扫描——通过Android智能手机的硬件安全元件处理。生物特征数据本身从不传输或存储。','how-works':'零知识证明系统','humans-how':'Aequitas使用Groth16证明系统——也被Zcash使用——密码学中最高效的ZKP系统之一。您的指纹签名被哈希为BN128椭圆曲线的域元素。该哈希用作Groth16电路的输入，产生一个小数学证明："使用了唯一的生物特征哈希"，而不揭示哈希是什么。',
'sybil-title':'永久防止女巫攻击','sybil-text':'女巫攻击是指一个实体创建多个虚假身份。这是几乎所有区块链的根本弱点。Aequitas在身份层解决了这个问题：每个生物特征哈希都永久存储。尝试用同一指纹注册两次会立即被拒绝。一人一钱包，永久——由密码学而非信任保证。',
'global-title':'为全球包容而设计','global-text':'Aequitas被设计为对地球上每个人都可访问。您不需要银行账户、信用卡或任何现有加密货币。只需一部带指纹传感器的Android智能手机——超过30亿人已经拥有的设备。注册完全免费，不到2分钟，立即获得1,000 AEQ。',
'registered-humans':'已注册人类','humans-desc':'此处列出的每个地址都通过生物特征零知识证明被验证为唯一人类。每人在注册时获得恰好1,000 AEQ。这些注册是永久、不可变的，存储在PostgreSQL和区块链上。','no-humans':'还没有人类注册。\n\n下载Aequitas Android应用成为链上第一个人！','registry-stats':'注册统计','total-humans-stat':'总人类数','grant':'每人补助','reg-fee':'注册费','free':'免费','zkp-system':'ZKP系统','never-stored':'从不存储','zkp-title':'ZKP技术细节','zkp-details':'BN128椭圆曲线上的Groth16。证明大小：~200字节。验证时间：~10毫秒。使用snarkjs和circom编译的电路。',
'aeq-index-title':'Aequitas指数 — 实时经济平等分数','aeq-index-desc':'Aequitas指数实时衡量网络的经济健康状况。直接从所有已验证人类的链上余额分布计算。分数0 = 每个人类有相同的AEQ（完全平等）。分数100 = 一个人控制一切（最大不平等）。基于链上余额的真实基尼系数计算。','current-index':'当前指数','bar-0':'0 — 完全平等','bar-100':'100 — 最大不平等','gini-coeff':'基尼系数','phase':'协议阶段',
'pools-title':'再分配池','pools-desc':'当不平等阈值被超过时，AEQ自动重定向到这四个池。这些池是完全由协议逻辑控制的智能合约账户——没有人类可以访问它们。','vel-pool':'速度池','liq-pool':'流动性池','uni-pool':'团结池','treasury':'国库','pools-details':'速度池奖励定期交易的钱包，激励经济活动。流动性池支持市场深度和价格稳定。团结池在后期阶段向新注册的人类授予额外AEQ。国库资助协议开发、安全审计和基础设施费用。',
'phases-title':'协议阶段','phases-desc':'随着网络增长，Aequitas协议经历四个不同阶段。阶段转换基于已验证人类和基尼系数自动发生——不需要投票。','phase0':'引导期 — 建立网络 · &lt;100人','phase1':'增长期 — 扩展注册 · 100-10,000人','phase2':'稳定期 — 再分配激活 · 10,000-100万人','phase3':'成熟期 — 完全去中心化 · 100万+人 · 基尼&lt;0.3',
'gini-title':'基尼系数 — 我们如何衡量不平等','gini-desc':'基尼系数由意大利统计学家科拉多·基尼于1912年开发，至今仍是世界上最广泛使用的经济不平等衡量标准。被世界银行、国际货币基金组织和几乎所有经济学家使用。Aequitas从实际链上余额实时计算这一指标。',
'gini-0':'完全平等','gini-0-sub':'每个人拥有相同财富。Aequitas在第0阶段接近这一点，那时所有人都有恰好1,000 AEQ。','gini-1':'低不平等','gini-1-sub':'斯堪的纳维亚平均（瑞典0.27，丹麦0.28）。Aequitas长期运营的目标区。','gini-2':'中度不平等','gini-2-sub':'美国平均（0.41）。再分配机制在第2阶段的这个水平激活。','gini-3':'高不平等','gini-3-sub':'南非（0.63）。比特币估计基尼超过0.85——高于地球上任何国家。',
'inflation-title':'通胀与再分配机制 — 详细说明',
'story-title':'Aequitas的故事 — 为什么存在',
'active-nodes':'活跃节点 — 当前网络拓扑','nodes-desc':'Aequitas网络在地理分布的云环境中运行两个节点。两者都参与区块生产、状态同步和API服务。它们通过libp2p（IPFS和以太坊2.0也使用）通信，通过HTTP同步区块。两个节点共享同一个PostgreSQL数据库。','node1':'节点1 — Railway（主要）','node2':'节点2 — Render（次要）','node1-desc':'主要API服务器 · 区块生产者 · P2P引导节点 · PostgreSQL连接 · MetaMask的RPC端点','node2-desc':'次要API服务器 · 区块生产者 · P2P对等节点 · HTTP区块同步 · 共享PostgreSQL · 冗余节点',
'bootstrap-title':'引导节点 — 加入网络','bootstrap-desc':'要运行您自己的Aequitas节点，请使用下面的libp2p多地址连接到引导节点。您的节点将自动发现对等节点，下载完整的区块历史，并开始参与共识。','bootstrap-howto':'运行节点：克隆GitHub仓库，设置DATABASE_URL环境变量，并使用引导对等地址运行二进制文件。',
'tech-specs':'技术规格','chain-id':'链ID','evm-yes':'是 — /rpc处的JSON-RPC · MetaMask兼容','language':'语言','source':'源代码',
'metamask-config':'MetaMask / Web3配置','metamask-desc':'Aequitas链兼容EVM，可以添加到MetaMask或任何支持自定义RPC网络的Web3钱包。使用下面的配置将Aequitas添加到您的钱包并查看您的AEQ余额。',
'architecture-title':'系统架构概述','architecture-desc':'Aequitas系统有四个主要组件。Android应用完全在设备上处理生物特征扫描和证明生成。证明服务器（Node.js）生成Groth16 ZK证明并存储生物特征哈希。区块链节点（Go）维护BlockDAG并公开EVM RPC。PostgreSQL存储所有节点共享的持久状态。整个技术栈是开源的。',
'reg-title':'🔐 注册为已验证人类','reg-sub':'加入Aequitas网络并接收您的1,000 AEQ。一次性、永久、无Gas费的注册，加密证明您是独特的人类。不收集或存储任何个人数据。无Gas费。无等待。',
'app-only-title':'仅通过ANDROID应用注册','app-only-text':'人类证明需要在您的个人设备上进行生物特征验证。指纹必须由您手机的硬件安全元件（HSE）处理——这是一个专用加密芯片，与主处理器隔离，无法远程访问。HSE是保护您银行应用程序和手机锁屏的同一芯片。当您在Aequitas应用中扫描时，签名由HSE处理、哈希并用于生成零知识证明——所有这些都不会让原始指纹数据离开HSE。下载应用，扫描指纹，连接钱包，您的<strong style="color:var(--gold)">1,000 AEQ将自动立即发放</strong>。',
'step1-title':'生物特征扫描','step1-desc':'打开Aequitas应用并点击"证明人类"。扫描您的指纹——完全由您手机的硬件安全元件处理。原始指纹数据永远不会到达应用的主进程，也不会以任何方式离开您的设备。HSE生成加密签名作为您独特的生物特征标识符。',
'step2-title':'ZKP生成','step2-desc':'应用将您的生物特征哈希发送到证明服务器。服务器检查此哈希之前未被使用（防止双重注册），然后在BN128上生成Groth16零知识证明。这个证明在不透露任何关于哈希的信息的情况下，数学上保证"使用了唯一的生物特征哈希"。需要2-5秒。',
'step3-title':'连接钱包','step3-desc':'应用打开MetaMask并预配置Aequitas链。连接您的钱包——此地址将接收您的1,000 AEQ。注册后，此钱包将永久与您的生物特征身份绑定。不重新注册无法更改。',
'step4-title':'1,000 AEQ已发放','step4-desc':'注册交易提交到区块链。协议验证ZK证明，确认生物特征哈希唯一，并在下一个区块（6秒内）将恰好1,000 AEQ存入您的钱包。从指纹到AEQ的总时间：不到30秒。',
'priv-bar':'🔒 硬件安全元件 · Groth16 ZKP · 生物特征数据永不离开设备 · 不收集个人数据 · 无Gas费 · 永久女巫攻击防护',
'connected-wallet':'已连接钱包','proof-detected':'⚡ 已从应用检测到证明参数','connect-btn':'🦊 连接METAMASK','register-btn':'🔐 链上注册','reg-hint':'// 打开Aequitas Android应用生成您的证明，然后返回此处完成注册...','reg-details':'注册详情','network':'网络','reg-limit':'注册次数','reg-limit-val':'每人一次 · 永久 · 不可变','bio-data':'生物特征数据','never-stored':'从不存储','confirmation':'确认时间','conf-time':'在6秒内（下一个区块）',
'phase-desc-0':'阶段0：引导期 — 建立网络','phase-desc-1':'阶段1：增长期 — 全球扩展人类注册','phase-desc-2':'阶段2：稳定期 — 财富再分配机制激活','phase-desc-3':'阶段3：成熟期 — 完全去中心化实现'
},
id:{
'live':'SIARAN LANGSUNG','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-register':'🔐 Daftar',
'live-stats':'Statistik Rantai Langsung','block-height':'Tinggi Blok','block-height-sub':'Blok baru setiap 6 detik · Konsensus BlockDAG · Dua node memproduksi blok secara paralel','verified-humans':'Manusia Terverifikasi','verified-humans-sub':'Setiap dompet diverifikasi sebagai manusia unik · ZKP biometrik · Satu orang, satu dompet, selamanya','total-supply':'Total Pasokan','total-supply-sub':'Selalu sama dengan Manusia × 1.000 AEQ · Pasokan hanya tumbuh ketika kemanusiaan tumbuh','aeq-index':'Indeks Aequitas','aeq-index-sub':'0 = kesetaraan sempurna · 100 = ketidaksetaraan maksimum · Berdasarkan koefisien Gini nyata','uptime':'Waktu Aktif','uptime-sub':'Node v0.3.0 · 2 node aktif · Railway + Render · Status persisten PostgreSQL',
'poh-title':'Bukti Kemanusiaan','poh-text':'Setiap pemegang AEQ harus membuktikan bahwa mereka adalah manusia unik yang hidup melalui verifikasi biometrik. Ini bukan opsional — ini adalah fondasi seluruh sistem. Tanpa bukti kemanusiaan, tidak ada AEQ yang dapat diterima. Ini berarti tidak ada bot, akun duplikat, korporasi, pemerintah, atau sistem AI yang dapat memegang AEQ. Hanya manusia nyata. Verifikasi menggunakan sidik jari Anda melalui Hardware Secure Element ponsel Anda. Data biometrik Anda tidak pernah meninggalkan perangkat Anda.',
'fair-title':'Distribusi yang Benar-Benar Adil','fair-text':'Setiap manusia terverifikasi di Bumi menerima tepat 1.000 AEQ — tidak lebih, tidak kurang. Orang pertama yang mendaftar dan orang ke satu miliar menerima jumlah yang sama. Tidak ada pre-mine, alokasi pendiri, putaran investor, atau keunggulan pengadopsi awal. Total pasokan AEQ selalu sama dengan manusia terverifikasi × 1.000. Ini adalah sistem distribusi moneter paling egalitarian yang pernah dirancang.',
'dag-title':'Arsitektur BlockDAG','dag-text':'Aequitas menggunakan Graf Asiklik Terarah (DAG) di mana beberapa blok dapat diproduksi secara bersamaan oleh node berbeda dan kemudian digabungkan. Ini memungkinkan throughput lebih tinggi, latensi lebih rendah, dan toleransi kesalahan lebih baik. Ketika dua node memproduksi blok secara bersamaan, keduanya valid dan kemudian digabungkan — Anda dapat melihat peristiwa penggabungan ini ditandai dengan 🔀 di penjelajah.',
'gasless-title':'Pendaftaran Tanpa Gas yang Sesungguhnya','gasless-text':'Salah satu hambatan terbesar untuk adopsi cryptocurrency adalah mengharuskan Anda sudah memiliki cryptocurrency untuk membayar biaya transaksi. Aequitas menghilangkan ini sepenuhnya. Pendaftaran tidak memerlukan biaya sama sekali — tidak perlu ETH, BNB, atau MATIC. Tidak perlu kartu kredit, tidak perlu rekening bank. Jika Anda adalah manusia dengan smartphone, Anda dapat mendaftar.',
'recent-blocks':'Blok Terbaru','blocks-desc':'Setiap blok terhubung secara kriptografis ke induknya melalui hash SHA-256. 🔀 MERGE = blok dengan beberapa induk (fitur BlockDAG). ✅ TX = blok berisi transaksi registrasi. Waktu blok: rata-rata 6 detik.','loading':'Memuat blok...','network-info':'Informasi Jaringan','chain-name':'Nama Jaringan','symbol':'Simbol','block-time':'Waktu Blok','consensus':'Konsensus','nodes':'Node Aktif','storage':'Penyimpanan Status','add-metamask':'🦊 TAMBAHKAN KE METAMASK','network-name':'Nama Jaringan','add-network':'+ TAMBAHKAN JARINGAN AEQUITAS','philosophy':'"Uang ada karena manusia ada.<br>Tidak lebih, tidak kurang."','philosophy-sub':'— PRINSIP AEQUITAS —','decimals':'Desimal',
'what-is-it':'Apa artinya Manusia Terverifikasi?','humans-what':'Manusia Terverifikasi di Aequitas adalah alamat dompet yang terbukti secara kriptografis milik manusia unik yang hidup. Verifikasi menggunakan data biometrik — khususnya pemindaian sidik jari — diproses melalui Hardware Secure Element smartphone Android. Data biometrik itu sendiri tidak pernah ditransmisikan atau disimpan.','how-works':'Sistem Bukti Zero-Knowledge','humans-how':'Aequitas menggunakan sistem pembuktian Groth16 — juga digunakan Zcash. Tanda tangan sidik jari Anda di-hash menjadi elemen medan kurva eliptik BN128. Hash ini digunakan sebagai input untuk sirkuit Groth16 yang menghasilkan bukti matematika kecil yang menjamin "hash biometrik unik digunakan" tanpa mengungkapkan apa hash tersebut.',
'sybil-title':'Pencegahan Permanen Serangan Sybil','sybil-text':'Serangan Sybil adalah ketika satu entitas membuat banyak identitas palsu. Ini adalah kelemahan mendasar hampir setiap blockchain. Aequitas menyelesaikan ini di lapisan identitas: setiap hash biometrik disimpan secara permanen. Mencoba mendaftar dua kali dengan sidik jari yang sama langsung ditolak. Satu orang, satu dompet, selamanya — dijamin oleh kriptografi, bukan kepercayaan.',
'global-title':'Dirancang untuk Inklusi Global','global-text':'Aequitas dirancang untuk dapat diakses oleh setiap manusia di Bumi. Anda tidak memerlukan rekening bank, kartu kredit, atau cryptocurrency yang ada. Hanya smartphone Android dengan sensor sidik jari — perangkat yang sudah dimiliki lebih dari 3 miliar orang. Pendaftaran sepenuhnya gratis, membutuhkan waktu kurang dari 2 menit, dan memberikan 1.000 AEQ secara instan.',
'registered-humans':'Manusia Terdaftar','humans-desc':'Setiap alamat yang terdaftar di sini telah diverifikasi sebagai manusia unik melalui Bukti Zero-Knowledge biometrik. Masing-masing menerima tepat 1.000 AEQ saat pendaftaran. Registrasi ini permanen, tidak dapat diubah, dan disimpan baik di PostgreSQL maupun on-chain.','no-humans':'Belum ada manusia terdaftar.\n\nUnduh Aplikasi Android Aequitas dan jadilah yang pertama di rantai!','registry-stats':'Statistik Registri','total-humans-stat':'Total Manusia','grant':'Hibah per Manusia','reg-fee':'Biaya Pendaftaran','free':'GRATIS','zkp-system':'Sistem ZKP','never-stored':'Tidak pernah disimpan','zkp-title':'Detail Teknis ZKP','zkp-details':'Groth16 atas kurva eliptik BN128. Ukuran bukti: ~200 byte. Waktu verifikasi: ~10ms. Sirkuit dikompilasi dengan snarkjs dan circom.',
'aeq-index-title':'Indeks Aequitas — Skor Kesetaraan Ekonomi Real-Time','aeq-index-desc':'Indeks Aequitas mengukur kesehatan ekonomi jaringan secara real-time. Dihitung langsung dari distribusi saldo on-chain semua manusia terverifikasi. Skor 0 = setiap manusia memiliki AEQ yang sama. Skor 100 = satu orang mengendalikan segalanya. Berdasarkan koefisien Gini nyata dari saldo on-chain.','current-index':'Indeks Saat Ini','bar-0':'0 — Kesetaraan Sempurna','bar-100':'100 — Ketidaksetaraan Maks.','gini-coeff':'Koefisien Gini','phase':'Fase Protokol',
'pools-title':'Pool Redistribusi','pools-desc':'Ketika ambang ketidaksetaraan terlampaui, AEQ secara otomatis diarahkan ke empat pool ini. Pool adalah akun kontrak pintar yang sepenuhnya dikendalikan oleh logika protokol — tidak ada manusia yang memiliki akses.','vel-pool':'Pool Kecepatan','liq-pool':'Pool Likuiditas','uni-pool':'Pool Kesatuan','treasury':'Perbendaharaan','pools-details':'Pool Kecepatan memberi penghargaan kepada dompet yang bertransaksi secara teratur, mendorong aktivitas ekonomi. Pool Likuiditas mendukung kedalaman pasar dan stabilitas harga. Pool Kesatuan memberikan AEQ tambahan kepada manusia yang baru terdaftar di fase selanjutnya. Perbendaharaan mendanai pengembangan protokol.',
'phases-title':'Fase Protokol','phases-desc':'Protokol Aequitas berkembang melalui empat fase berbeda seiring pertumbuhan jaringan. Transisi fase terjadi secara otomatis berdasarkan manusia terverifikasi dan koefisien Gini — tidak diperlukan pemungutan suara.','phase0':'Bootstrap — Membangun jaringan · &lt;100 manusia','phase1':'Pertumbuhan — Memperluas registri · 100-10.000 manusia','phase2':'Stabilitas — Redistribusi aktif · 10.000-1M manusia','phase3':'Kedewasaan — Desentralisasi penuh · 1M+ manusia · Gini &lt;0,3',
'gini-title':'Koefisien Gini — Bagaimana Kami Mengukur Ketidaksetaraan','gini-desc':'Koefisien Gini dikembangkan oleh statistikawan Italia Corrado Gini pada tahun 1912 dan tetap menjadi ukuran ketidaksetaraan ekonomi yang paling banyak digunakan di dunia. Digunakan oleh Bank Dunia, IMF, dan hampir semua ekonom. Aequitas menghitung ini secara real-time dari saldo on-chain aktual.',
'gini-0':'Kesetaraan Sempurna','gini-0-sub':'Setiap orang memiliki kekayaan yang sama. Aequitas mendekati ini di Fase 0 ketika semua manusia memiliki tepat 1.000 AEQ.','gini-1':'Ketidaksetaraan Rendah','gini-1-sub':'Rata-rata Skandinavia (Swedia 0,27, Denmark 0,28). Zona target Aequitas untuk operasi jangka panjang.','gini-2':'Ketidaksetaraan Sedang','gini-2-sub':'Rata-rata AS (0,41). Mekanisme redistribusi aktif pada tingkat ini di Fase 2.','gini-3':'Ketidaksetaraan Tinggi','gini-3-sub':'Afrika Selatan (0,63). Perkiraan Gini Bitcoin melebihi 0,85 — lebih tinggi dari negara mana pun.',
'inflation-title':'Mekanisme Inflasi & Redistribusi — Secara Rinci',
'story-title':'Kisah Aequitas — Mengapa Ini Ada',
'active-nodes':'Node Aktif — Topologi Jaringan Saat Ini','nodes-desc':'Jaringan Aequitas beroperasi pada dua node di lingkungan cloud yang didistribusikan secara geografis. Keduanya berpartisipasi dalam produksi blok, sinkronisasi status, dan layanan API. Mereka berkomunikasi melalui libp2p dan menyinkronkan blok melalui HTTP. Kedua node berbagi database PostgreSQL yang sama.','node1':'Node 1 — Railway (Utama)','node2':'Node 2 — Render (Sekunder)','node1-desc':'Server API utama · Produsen blok · Node bootstrap P2P · Koneksi PostgreSQL · Endpoint RPC untuk MetaMask','node2-desc':'Server API sekunder · Produsen blok · Peer P2P · Sinkronisasi blok HTTP · PostgreSQL bersama · Node redundansi',
'bootstrap-title':'Node Bootstrap — Bergabung dengan Jaringan','bootstrap-desc':'Untuk menjalankan node Aequitas Anda sendiri, hubungkan ke node bootstrap menggunakan alamat libp2p multi di bawah. Node Anda akan secara otomatis menemukan peer, mengunduh riwayat blok lengkap, dan mulai berpartisipasi dalam konsensus.','bootstrap-howto':'Untuk menjalankan node: klon repositori GitHub, atur variabel lingkungan DATABASE_URL, dan jalankan biner dengan alamat peer bootstrap.',
'tech-specs':'Spesifikasi Teknis','chain-id':'ID Rantai','evm-yes':'Ya — JSON-RPC di /rpc · Kompatibel MetaMask','language':'Bahasa','source':'Kode Sumber',
'metamask-config':'Konfigurasi MetaMask / Web3','metamask-desc':'Aequitas Chain kompatibel dengan EVM dan dapat ditambahkan ke MetaMask atau dompet Web3 mana pun yang mendukung jaringan RPC kustom.',
'architecture-title':'Ikhtisar Arsitektur Sistem','architecture-desc':'Sistem Aequitas memiliki empat komponen utama. Aplikasi Android menangani pemindaian biometrik dan pembuatan bukti sepenuhnya di perangkat. Server Bukti (Node.js) menghasilkan bukti ZK Groth16 dan menyimpan hash biometrik. Node Blockchain (Go) mempertahankan BlockDAG dan mengekspos EVM RPC. PostgreSQL menyimpan status persisten yang dibagikan ke semua node.',
'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi','reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima 1.000 AEQ Anda. Pendaftaran sekali, permanen, tanpa gas yang membuktikan secara kriptografis bahwa Anda adalah manusia unik. Tidak ada data pribadi yang dikumpulkan atau disimpan.',
'app-only-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID','app-only-text':'Bukti Kemanusiaan memerlukan verifikasi biometrik di perangkat pribadi Anda. Sidik jari harus diproses oleh Hardware Secure Element (HSE) ponsel Anda — chip kriptografi khusus yang terisolasi dari prosesor utama dan tidak dapat diakses dari jarak jauh. HSE adalah chip yang sama yang mengamankan aplikasi perbankan Anda. Ketika Anda memindai di Aplikasi Aequitas, tanda tangan diproses oleh HSE, di-hash, dan digunakan untuk menghasilkan Bukti Zero-Knowledge — semua tanpa data sidik jari mentah pernah keluar dari HSE. Unduh aplikasinya, pindai sidik jari Anda, hubungkan dompet Anda, dan <strong style="color:var(--gold)">1.000 AEQ Anda akan diberikan secara otomatis dan segera</strong>.',
'step1-title':'Pemindaian Biometrik','step1-desc':'Buka Aplikasi Aequitas dan ketuk "Buktikan Kemanusiaan." Pindai sidik jari Anda — diproses sepenuhnya oleh Hardware Secure Element ponsel Anda. Data sidik jari mentah tidak pernah mencapai proses utama aplikasi dan tidak pernah meninggalkan perangkat Anda. HSE menghasilkan tanda tangan kriptografis sebagai pengenal biometrik unik Anda.',
'step2-title':'Pembuatan ZKP','step2-desc':'Aplikasi mengirim hash biometrik Anda ke Server Bukti. Server memeriksa bahwa hash ini belum pernah digunakan sebelumnya (mencegah registrasi ganda), kemudian menghasilkan Bukti Zero-Knowledge Groth16 atas BN128. Bukti ini secara matematis menjamin "hash biometrik unik digunakan" tanpa mengungkapkan apa pun tentang hash. Membutuhkan 2-5 detik.',
'step3-title':'Hubungkan Dompet','step3-desc':'Aplikasi membuka MetaMask dengan Aequitas Chain yang telah dikonfigurasi sebelumnya. Hubungkan dompet Anda — alamat ini akan menerima 1.000 AEQ Anda. Setelah terdaftar, dompet ini secara permanen terhubung ke identitas biometrik Anda. Anda tidak dapat mengubahnya tanpa mendaftar ulang.',
'step4-title':'1.000 AEQ Diberikan','step4-desc':'Transaksi registrasi dikirim ke blockchain. Protokol memverifikasi bukti ZK, mengkonfirmasi hash biometrik unik, dan mengkreditkan tepat 1.000 AEQ ke dompet Anda dalam blok berikutnya (dalam 6 detik). Total waktu dari sidik jari ke AEQ: kurang dari 30 detik.',
'priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Data biometrik tidak pernah meninggalkan perangkat · Tidak ada data pribadi · Tanpa biaya gas · Perlindungan Sybil permanen',
'connected-wallet':'DOMPET TERHUBUNG','proof-detected':'⚡ PARAMETER BUKTI TERDETEKSI DARI APLIKASI','connect-btn':'🦊 HUBUNGKAN METAMASK','register-btn':'🔐 DAFTAR ON-CHAIN','reg-hint':'// Buka Aplikasi Android Aequitas untuk menghasilkan bukti Anda, lalu kembali ke sini untuk menyelesaikan pendaftaran...','reg-details':'Detail Pendaftaran','network':'Jaringan','reg-limit':'Pendaftaran','reg-limit-val':'Sekali per manusia · permanen · tidak dapat diubah','bio-data':'Data Biometrik','never-stored':'Tidak pernah disimpan','confirmation':'Waktu Konfirmasi','conf-time':'Dalam 6 detik (blok berikutnya)',
'phase-desc-0':'Fase 0: Bootstrap — Membangun jaringan','phase-desc-1':'Fase 1: Pertumbuhan — Memperluas registri manusia secara global','phase-desc-2':'Fase 2: Stabilitas — Mekanisme redistribusi kekayaan aktif','phase-desc-3':'Fase 3: Kedewasaan — Desentralisasi penuh tercapai'
}
};

const STORY={
en:'<p>The year is 2009. Satoshi Nakamoto releases Bitcoin — the first successful decentralized digital currency. For the first time in history, it is possible to transfer value between any two people on Earth without the permission of any bank, government, or intermediary. It is a genuine revolution. But something goes wrong almost immediately.</p><p>Early Bitcoin miners accumulate millions of coins that cost them almost nothing. As Bitcoin\'s price rises, these early adopters become extraordinarily wealthy. By 2021, the top 1% of Bitcoin addresses control over 90% of all Bitcoin. The cryptocurrency that was supposed to democratize finance has instead created some of the most extreme wealth concentration in human history. Bitcoin\'s estimated Gini coefficient exceeds 0.85 — higher than any country on Earth, higher than historical feudal societies, approaching the theoretical maximum of inequality.</p><p>Ethereum, Solana, and virtually every other cryptocurrency follow the same pattern: pre-mine for founders and early investors, ICO for wealthy participants, and the rest of humanity priced out before they even hear about it. The blockchain revolution failed to deliver on its promise of financial inclusion.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness," "equity," and "equality" — was created to answer one question: <em style="color:var(--gold)">"What would a cryptocurrency look like if it was designed from first principles to be fair to every human being who has ever lived or will ever live?"</em></p><p>The answer turned out to be surprisingly simple: <strong style="color:var(--text)">Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being a person.</strong> Not because they were born into wealth. Not because they were early to a speculative bet. Not because they have access to expensive mining hardware. Simply because they are human.</p><p>The technology to make this possible now exists. Zero-Knowledge Proofs allow us to verify that a person is a unique human being without requiring them to reveal any personal information. Blockchain technology allows us to store these verifications permanently and transparently. Smartphones — now owned by over 3 billion people — provide the biometric sensors needed for verification. The pieces are in place. Aequitas assembles them into a coherent system for the first time.</p><p>The Aequitas network launched in June 2026. It is currently in Phase 0 — the bootstrap phase. Every human who registers now is helping to prove that a fairer monetary system is possible. The goal is not to replace existing currencies overnight. The goal is to demonstrate, on a live network with real cryptography and real economics, that money can be distributed fairly, that equality can be maintained through mathematical governance, and that financial inclusion can be achieved at global scale without compromising security or decentralization.</p><p><em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em> This is not just a slogan — it is the mathematical foundation of the entire system. Every line of code, every cryptographic proof, every protocol rule flows from this single insight. We invite every human being on Earth to join us.</p>',
de:'<p>Das Jahr ist 2009. Satoshi Nakamoto veröffentlicht Bitcoin — die erste erfolgreiche dezentralisierte digitale Währung. Zum ersten Mal in der Geschichte ist es möglich, Wert zwischen zwei beliebigen Menschen auf der Erde zu übertragen ohne die Erlaubnis einer Bank, Regierung oder eines Vermittlers. Es ist eine echte Revolution. Aber fast sofort geht etwas schief.</p><p>Frühe Bitcoin-Miner häufen Millionen von Coins an, die sie fast nichts kosten. Als der Bitcoin-Preis steigt, werden diese frühen Adopter außerordentlich reich. Bis 2021 kontrolliert das oberste 1% der Bitcoin-Adressen über 90% aller Bitcoins. Die Kryptowährung, die den Finanzmarkt demokratisieren sollte, hat stattdessen einige der extremsten Vermögenskonzentrationen in der menschlichen Geschichte geschaffen. Bitcoins geschätzter Gini-Koeffizient übersteigt 0,85 — höher als jedes Land auf der Erde.</p><p>Ethereum, Solana und praktisch jede andere Kryptowährung folgen demselben Muster. Die Blockchain-Revolution hat ihr Versprechen der finanziellen Inklusion nicht eingelöst.</p><p><span style="color:var(--gold)">Aequitas</span> — Lateinisch für "Fairness," "Gerechtigkeit" und "Gleichheit" — wurde geschaffen um eine Frage zu beantworten: <em style="color:var(--gold)">"Wie würde eine Kryptowährung aussehen, die von Grund auf fair für jeden Menschen konzipiert wurde?"</em></p><p>Die Antwort war überraschend einfach: <strong style="color:var(--text)">Geld existiert weil Menschen existieren. Daher sollte jeder Mensch einfach aufgrund seiner Menschlichkeit einen gleichen Anteil am Geld haben.</strong> Nicht weil er in Wohlstand geboren wurde. Nicht weil er früh bei einer spekulativen Wette dabei war. Sondern einfach weil er ein Mensch ist.</p><p>Das Aequitas-Netzwerk wurde im Juni 2026 gestartet. Es befindet sich derzeit in Phase 0. Jeder Mensch, der sich jetzt registriert, hilft zu beweisen, dass ein gerechteres Geldsystem möglich ist. Das Ziel ist es zu demonstrieren, dass Geld fair verteilt werden kann, dass Gleichheit durch mathematische Governance aufrechterhalten werden kann, und dass finanzielle Inklusion im globalen Maßstab ohne Kompromisse bei Sicherheit oder Dezentralisierung erreicht werden kann.</p><p><em style="color:var(--gold)">"Geld existiert weil Menschen existieren. Nichts mehr, nichts weniger."</em> Dies ist nicht nur ein Slogan — es ist das mathematische Fundament des gesamten Systems. Wir laden jeden Menschen auf der Erde ein, sich uns anzuschließen.</p>',
es:'<p>El año es 2009. Satoshi Nakamoto lanza Bitcoin — la primera moneda digital descentralizada exitosa. Por primera vez en la historia, es posible transferir valor entre dos personas cualesquiera en la Tierra sin el permiso de ningún banco, gobierno o intermediario. Es una revolución genuina. Pero algo sale mal casi de inmediato.</p><p>Los primeros mineros de Bitcoin acumulan millones de monedas que les cuestan casi nada. A medida que sube el precio de Bitcoin, estos primeros adoptantes se vuelven extraordinariamente ricos. Para 2021, el 1% superior de las direcciones de Bitcoin controla más del 90% de todo el Bitcoin. La criptomoneda que se suponía democratizaría las finanzas ha creado en cambio algunas de las concentraciones de riqueza más extremas en la historia humana.</p><p><span style="color:var(--gold)">Aequitas</span> — Latín para "justicia," "equidad" e "igualdad" — fue creado para responder una pregunta: <em style="color:var(--gold)">"¿Cómo sería una criptomoneda diseñada desde sus principios para ser justa con todo ser humano?"</em></p><p>La respuesta resultó ser sorprendentemente simple: <strong style="color:var(--text)">El dinero existe porque las personas existen. Por lo tanto, cada persona debería tener una parte igual del dinero simplemente por ser persona.</strong></p><p>La red Aequitas se lanzó en junio de 2026. Actualmente está en Fase 0. El objetivo es demostrar, en una red en vivo con criptografía real y economía real, que el dinero puede distribuirse equitativamente y que la inclusión financiera puede lograrse a escala global.</p><p><em style="color:var(--gold)">"El dinero existe porque las personas existen. Nada más, nada menos."</em> Invitamos a todo ser humano en la Tierra a unirse a nosotros.</p>',
ru:'<p>2009 год. Сатоши Накамото выпускает Биткоин — первую успешную децентрализованную цифровую валюту. Впервые в истории стало возможным передавать ценность между любыми двумя людьми на Земле без разрешения какого-либо банка, правительства или посредника. Это подлинная революция. Но почти сразу что-то идёт не так.</p><p>Ранние майнеры Bitcoin накапливают миллионы монет, которые стоят им почти ничего. По мере роста цены Bitcoin эти ранние участники становятся невероятно богатыми. К 2021 году верхний 1% адресов Bitcoin контролирует более 90% всех Bitcoin. Криптовалюта, которая должна была демократизировать финансы, вместо этого создала одну из самых крайних концентраций богатства в истории человечества.</p><p><span style="color:var(--gold)">Aequitas</span> — латинское слово "справедливость," "равноправие" и "равенство" — был создан для ответа на вопрос: <em style="color:var(--gold)">"Как выглядела бы криптовалюта, разработанная с нуля для справедливости к каждому человеческому существу?"</em></p><p>Ответ оказался удивительно простым: <strong style="color:var(--text)">Деньги существуют потому что существуют люди. Поэтому каждый человек должен иметь равную долю денег просто в силу того, что он является человеком.</strong></p><p>Сеть Aequitas запустилась в июне 2026 года. В настоящее время она находится в Фазе 0. Цель — продемонстрировать, что деньги могут распределяться справедливо и что финансовая инклюзия может быть достигнута в глобальном масштабе.</p><p><em style="color:var(--gold)">"Деньги существуют потому что существуют люди. Ничего больше, ничего меньше."</em> Мы приглашаем каждого человека на Земле присоединиться к нам.</p>',
zh:'<p>2009年。中本聪发布比特币——第一种成功的去中心化数字货币。有史以来第一次，可以在地球上任意两个人之间转移价值，无需任何银行、政府或中介的许可。这是一场真正的革命。但几乎立即就出现了问题。</p><p>早期比特币矿工以几乎为零的成本积累了数百万枚比特币。随着比特币价格上涨，这些早期采用者变得极为富有。到2021年，比特币地址前1%控制了超过90%的所有比特币。这种本应使金融民主化的加密货币，反而创造了人类历史上一些最极端的财富集中。比特币的估计基尼系数超过0.85——高于地球上任何国家。</p><p><span style="color:var(--gold)">Aequitas</span>——拉丁语"公平"、"公正"和"平等"——被创建来回答一个问题：<em style="color:var(--gold)">"如果一种加密货币从第一原则出发设计，对每个人都公平，它会是什么样子？"</em></p><p>答案出乎意料地简单：<strong style="color:var(--text)">货币存在是因为人类存在。因此，每个人仅凭其是人类这一事实，就应该拥有等额的货币。</strong></p><p>Aequitas网络于2026年6月启动。目前处于第0阶段。我们的目标是证明，货币可以公平分配，平等可以通过数学治理来维持，金融普惠可以在全球范围内实现。</p><p><em style="color:var(--gold)">"货币存在是因为人类存在。仅此而已，不多也不少。"</em>我们邀请地球上每一个人类加入我们。</p>',
id:'<p>Tahun 2009. Satoshi Nakamoto merilis Bitcoin — mata uang digital terdesentralisasi pertama yang sukses. Untuk pertama kalinya dalam sejarah, dimungkinkan untuk mentransfer nilai antara dua orang mana pun di Bumi tanpa izin bank, pemerintah, atau perantara mana pun. Ini adalah revolusi sejati. Tetapi sesuatu segera berjalan salah.</p><p>Penambang Bitcoin awal mengumpulkan jutaan koin yang hampir tidak memerlukan biaya. Seiring naiknya harga Bitcoin, para pengadopsi awal ini menjadi sangat kaya. Pada 2021, 1% teratas alamat Bitcoin mengendalikan lebih dari 90% semua Bitcoin. Cryptocurrency yang seharusnya mendemokratisasi keuangan malah menciptakan beberapa konsentrasi kekayaan paling ekstrem dalam sejarah manusia.</p><p><span style="color:var(--gold)">Aequitas</span> — bahasa Latin untuk "keadilan," "ekuitas," dan "kesetaraan" — diciptakan untuk menjawab satu pertanyaan: <em style="color:var(--gold)">"Seperti apa cryptocurrency jika dirancang dari prinsip pertama untuk adil bagi setiap manusia?"</em></p><p>Jawabannya ternyata sangat sederhana: <strong style="color:var(--text)">Uang ada karena manusia ada. Oleh karena itu, setiap orang harus memiliki bagian yang sama dari uang hanya karena menjadi manusia.</strong></p><p>Jaringan Aequitas diluncurkan pada Juni 2026. Saat ini berada di Fase 0. Tujuannya adalah menunjukkan bahwa uang dapat didistribusikan secara adil dan inklusi keuangan dapat dicapai pada skala global.</p><p><em style="color:var(--gold)">"Uang ada karena manusia ada. Tidak lebih, tidak kurang."</em> Kami mengundang setiap manusia di Bumi untuk bergabung dengan kami.</p>'
};

const INFLATION={
en:'<p><span style="color:var(--gold);font-weight:bold">Why does Aequitas need an inflation mechanism?</span> In a world where people trade, save, invest, and lose money, perfect equality cannot be maintained forever even with an equal starting distribution. Over time, some people will accumulate more AEQ and some will spend theirs. Without a correction mechanism, Aequitas would eventually look like any other unequal monetary system. The inflation and redistribution mechanism is the answer to this fundamental challenge.</p><p><strong style="color:var(--blue)">Base Inflation — The Only Truly Justified Inflation:</strong> The only source of new AEQ is human registration. When a new human is verified and joins the network, exactly 1,000 AEQ is created and sent to their wallet. This is the only inflation in Phase 0 and Phase 1. There is no other mechanism that creates AEQ — no mining rewards, no staking rewards, no protocol emissions, no quantitative easing. The supply grows if and only if the number of verified humans grows. This means AEQ\'s purchasing power is directly tied to human population growth — historically one of the most stable and predictable growth rates in the world.</p><p><strong style="color:var(--gold)">The Wealth Cap — Preventing Extreme Concentration:</strong> In Phase 2 and beyond, a dynamic wealth cap is enforced. The cap is calculated as a multiple of the mean balance (total supply divided by number of humans), adjusted by the current Gini coefficient. A higher Gini means a lower cap, creating a more aggressive ceiling during times of high inequality. When any wallet\'s balance exceeds this cap, the excess is not seized — instead, all new AEQ earned by that wallet (from velocity rewards, pool distributions, etc.) is redirected to the four redistribution pools until the balance falls below the cap. This is not confiscation — it is a ceiling on accumulation, applied fairly to everyone including the protocol\'s own treasury and the founding team\'s wallets.</p><p><strong style="color:var(--purple)">Dynamic Redistribution Cycles — Automatic Economic Correction:</strong> The Keeper Bot (a Node.js process running alongside the blockchain nodes) executes redistribution cycles on a scheduled basis. During each cycle, it reads the current Gini coefficient calculated from all on-chain balances. If the Gini exceeds 0.25 (the Phase 2 activation threshold), a percentage of the Velocity Pool proportional to the excess above the threshold is distributed equally to all verified humans. If the Gini exceeds 0.35, this percentage increases. If the Gini exceeds 0.50, emergency redistribution activates with a significantly larger percentage transfer. The mathematics create a powerful automatic stabilizer: the more unequal the distribution becomes, the stronger the force pushing it back toward equality.</p><p><strong style="color:var(--green)">Velocity Incentives — Rewarding Economic Participation:</strong> The single biggest driver of inequality in any monetary system is hoarding — a small number of entities accumulate vast wealth and simply hold it, removing it from circulation and making it unavailable to others. The Velocity Pool directly combats this by rewarding wallets that transact regularly. Every time a verified human sends or receives AEQ, their velocity score increases. At the end of each redistribution cycle, the Velocity Pool is distributed proportionally to velocity scores across all wallets that participated in the economy. This creates a powerful incentive to keep money moving through the economy rather than sitting in a wallet accumulating dust.</p><p><strong style="color:var(--teal)">Mathematical Governance — The End of Monetary Politics:</strong> Every rule described above — the wealth cap formula, the redistribution percentages, the Gini thresholds, the pool allocation ratios — is encoded in immutable smart contract code deployed on the Aequitas blockchain. No person, organization, government, or founding team can change these rules without a hard fork that requires the consensus of the entire network. There is no governance token that wealthy holders can vote with. There is no committee that can meet and change the rules. There is no central bank. The Aequitas protocol is governed by mathematics, not by committee meetings, political lobbying, or the preferences of wealthy investors. This is the fundamental promise of Aequitas: a monetary system where the rules apply equally to everyone and cannot be changed by anyone with power or money.</p>',
de:'<p><span style="color:var(--gold);font-weight:bold">Warum braucht Aequitas einen Inflationsmechanismus?</span> In einer Welt, in der Menschen handeln, sparen, investieren und Geld verlieren, kann vollkommene Gleichheit nicht für immer aufrechterhalten werden, selbst mit einer gleichen Ausgangsverteilung. Mit der Zeit werden manche Menschen mehr AEQ ansammeln und andere ihres ausgeben. Ohne Korrektionsmechanismus würde Aequitas schließlich wie jedes andere ungleiche Geldsystem aussehen.</p><p><strong style="color:var(--blue)">Basisinflation — Die einzige wirklich gerechtfertigte Inflation:</strong> Die einzige Quelle für neues AEQ ist die Menschenregistrierung. Wenn ein neuer Mensch verifiziert wird und dem Netzwerk beitritt, werden genau 1.000 AEQ erstellt und an seine Wallet gesendet. Dies ist die einzige Inflation in Phase 0 und Phase 1. Es gibt keinen anderen Mechanismus, der AEQ erstellt — keine Mining-Belohnungen, keine Staking-Belohnungen, kein quantitatives Easing.</p><p><strong style="color:var(--gold)">Vermögensobergrenze — Extreme Konzentration verhindern:</strong> In Phase 2 und darüber hinaus wird eine dynamische Vermögensobergrenze durchgesetzt. Die Obergrenze wird als Vielfaches des Durchschnittsguthabens berechnet, angepasst durch den aktuellen Gini-Koeffizienten. Wenn das Guthaben einer Wallet diese Obergrenze überschreitet, wird das Überschüssige nicht beschlagnahmt — stattdessen wird alles neue AEQ, das diese Wallet verdient, in die vier Umverteilungspools umgeleitet.</p><p><strong style="color:var(--purple)">Dynamische Umverteilungszyklen — Automatische wirtschaftliche Korrektur:</strong> Der Keeper Bot führt Umverteilungszyklen planmäßig aus. Während jedes Zyklus liest er den aktuellen Gini-Koeffizienten. Wenn der Gini 0,25 übersteigt, wird ein Teil des Velocity-Pools gleichmäßig an alle verifizierten Menschen verteilt. Je höher die Ungleichheit, desto aggressiver die Korrektur.</p><p><strong style="color:var(--green)">Velocity-Anreize — Wirtschaftliche Beteiligung belohnen:</strong> Der größte Treiber von Ungleichheit in jedem Geldsystem ist das Horten. Der Velocity-Pool bekämpft dies direkt, indem er Wallets belohnt, die regelmäßig handeln. Jedes Mal, wenn ein verifizierter Mensch AEQ sendet oder empfängt, steigt sein Velocity-Score.</p><p><strong style="color:var(--teal)">Mathematische Steuerung — Das Ende der monetären Politik:</strong> Jede oben beschriebene Regel ist in unveränderlichem Smart-Contract-Code kodiert. Keine Person, Organisation, Regierung oder das Gründerteam kann diese Regeln ändern. Das Aequitas-Protokoll wird durch Mathematik regiert, nicht durch Ausschusssitzungen oder politisches Lobbying.</p>',
es:'<p><span style="color:var(--gold);font-weight:bold">¿Por qué Aequitas necesita un mecanismo de inflación?</span> En un mundo donde las personas comercian, ahorran, invierten y pierden dinero, la igualdad perfecta no puede mantenerse para siempre incluso con una distribución inicial igual. Con el tiempo, algunas personas acumularán más AEQ. Sin un mecanismo de corrección, Aequitas eventualmente parecería cualquier otro sistema monetario desigual.</p><p><strong style="color:var(--blue)">Inflación Base — La única inflación verdaderamente justificada:</strong> La única fuente de nuevo AEQ es el registro humano. Cuando un nuevo humano es verificado, se crean exactamente 1,000 AEQ. No hay minería, no hay staking, no hay emisión arbitraria. El suministro crece si y solo si el número de humanos verificados crece.</p><p><strong style="color:var(--gold)">Límite de Riqueza — Prevenir la concentración extrema:</strong> En la Fase 2, se aplica un límite de riqueza dinámico. Cuando el saldo de cualquier wallet supera este límite, todo el nuevo AEQ ganado se redirige a los cuatro pools de redistribución hasta que el saldo caiga por debajo del límite.</p><p><strong style="color:var(--purple)">Ciclos de Redistribución Dinámica — Corrección económica automática:</strong> El Keeper Bot ejecuta ciclos de redistribución programados. Si el Gini supera 0.25, un porcentaje del Pool de Velocidad se distribuye igualmente a todos los humanos verificados. Cuanto mayor sea la desigualdad, más agresiva será la corrección.</p><p><strong style="color:var(--green)">Incentivos de Velocidad — Recompensando la participación económica:</strong> El Pool de Velocidad combate directamente el acaparamiento recompensando las wallets que transaccionan regularmente.</p><p><strong style="color:var(--teal)">Gobernanza Matemática:</strong> Cada regla está codificada en código de contrato inteligente inmutable. Nadie puede cambiar estas reglas. Las matemáticas, no la política, gobiernan el suministro monetario de Aequitas.</p>',
ru:'<p><span style="color:var(--gold);font-weight:bold">Почему Aequitas нуждается в механизме инфляции?</span> В мире, где люди торгуют, сберегают, инвестируют и теряют деньги, совершенное равенство не может поддерживаться вечно. Без механизма коррекции Aequitas в конечном итоге стал бы похожим на любую другую неравную денежную систему.</p><p><strong style="color:var(--blue)">Базовая инфляция:</strong> Единственный источник нового AEQ — регистрация людей. Когда новый человек верифицирован, создаётся ровно 1 000 AEQ. Нет майнинга, нет стейкинга, нет произвольной эмиссии. Предложение растёт если и только если растёт число верифицированных людей.</p><p><strong style="color:var(--gold)">Ограничение богатства:</strong> В Фазе 2 применяется динамический предел богатства. Когда баланс кошелька превышает его, весь новый AEQ перенаправляется в четыре пула перераспределения.</p><p><strong style="color:var(--purple)">Динамические циклы перераспределения:</strong> Keeper Bot выполняет запланированные циклы перераспределения. Если коэффициент Джини превышает 0,25, часть Пула Скорости равномерно распределяется между всеми верифицированными людьми. Чем выше неравенство, тем агрессивнее коррекция.</p><p><strong style="color:var(--green)">Стимулы Скорости:</strong> Пул Скорости напрямую борется с накоплением, вознаграждая кошельки, которые регулярно совершают транзакции.</p><p><strong style="color:var(--teal)">Математическое управление:</strong> Каждое правило закодировано в неизменяемом коде смарт-контракта. Никто не может изменить эти правила. Математика, а не политика, управляет денежным предложением Aequitas.</p>',
zh:'<p><span style="color:var(--gold);font-weight:bold">为什么Aequitas需要通胀机制？</span>在一个人们交易、储蓄、投资和亏钱的世界里，即使是等额的起始分配也无法永远维持完全平等。随着时间推移，一些人会积累更多AEQ。没有纠正机制，Aequitas最终会看起来像任何其他不平等的货币系统。</p><p><strong style="color:var(--blue)">基础通胀——唯一真正合理的通胀：</strong>新AEQ的唯一来源是人类注册。当一个新人类被验证并加入网络时，恰好创建1,000 AEQ。没有挖矿奖励，没有质押奖励，没有任意发行。供应量增长当且仅当已验证人类数量增长时。</p><p><strong style="color:var(--gold)">财富上限——防止极端集中：</strong>在第2阶段，执行动态财富上限。当任何钱包的余额超过此上限时，该钱包获得的所有新AEQ都被重定向到四个再分配池，直到余额降至上限以下。</p><p><strong style="color:var(--purple)">动态再分配周期——自动经济纠正：</strong>如果基尼系数超过0.25，速度池的一个百分比等额分配给所有已验证人类。不平等程度越高，纠正力度越大。</p><p><strong style="color:var(--green)">速度激励——奖励经济参与：</strong>速度池通过奖励定期交易的钱包直接对抗囤积行为。</p><p><strong style="color:var(--teal)">数学治理：</strong>每条规则都编码在不可变的智能合约代码中。没有人能改变这些规则。数学而非政治治理着Aequitas的货币供应。</p>',
id:'<p><span style="color:var(--gold);font-weight:bold">Mengapa Aequitas memerlukan mekanisme inflasi?</span>Di dunia di mana orang berdagang, menabung, berinvestasi, dan kehilangan uang, kesetaraan sempurna tidak dapat dipertahankan selamanya bahkan dengan distribusi awal yang sama. Tanpa mekanisme koreksi, Aequitas akhirnya akan terlihat seperti sistem moneter tidak setara lainnya.</p><p><strong style="color:var(--blue)">Inflasi Dasar:</strong>Satu-satunya sumber AEQ baru adalah registrasi manusia. Ketika manusia baru diverifikasi, tepat 1.000 AEQ dibuat. Tidak ada hadiah penambangan, tidak ada hadiah staking, tidak ada emisi sembarangan. Pasokan tumbuh jika dan hanya jika jumlah manusia terverifikasi tumbuh.</p><p><strong style="color:var(--gold)">Batas Kekayaan:</strong>Ketika saldo dompet mana pun melebihi batas dinamis, semua AEQ baru yang diperoleh dompet tersebut diarahkan ke empat pool redistribusi.</p><p><strong style="color:var(--purple)">Siklus Redistribusi Dinamis:</strong>Jika koefisien Gini melebihi 0,25, sebagian Pool Kecepatan didistribusikan secara merata ke semua manusia terverifikasi. Semakin tinggi ketidaksetaraan, semakin agresif koreksinya.</p><p><strong style="color:var(--green)">Insentif Kecepatan:</strong>Pool Kecepatan secara langsung memerangi penimbunan dengan memberi penghargaan kepada dompet yang bertransaksi secara teratur.</p><p><strong style="color:var(--teal)">Tata Kelola Matematis:</strong>Setiap aturan dikodekan dalam kode kontrak pintar yang tidak dapat diubah. Tidak ada yang bisa mengubah aturan-aturan ini. Matematika mengatur pasokan uang Aequitas.</p>'
};

function setLang(lang){
  currentLang=lang;
  document.querySelectorAll('.lang-btn').forEach(b=>b.classList.remove('active'));
  const lb=document.getElementById('lb-'+lang);
  if(lb)lb.classList.add('active');
  const t=T[lang];
  if(!t)return;
  document.querySelectorAll('[data-i18n]').forEach(el=>{
    const key=el.getAttribute('data-i18n');
    if(t[key]!==undefined)el.innerHTML=t[key];
  });
  document.querySelectorAll('[data-i18n-tab]').forEach(el=>{
    const key=el.getAttribute('data-i18n-tab');
    if(t[key]!==undefined)el.innerHTML=t[key];
  });
  const story=document.querySelector('[data-i18n="story-text"]');
  if(story&&STORY[lang])story.innerHTML=STORY[lang];
  const inflation=document.querySelector('[data-i18n="inflation-text"]');
  if(inflation&&INFLATION[lang])inflation.innerHTML=INFLATION[lang];
  const phasedesc=document.getElementById('idx-phase-desc');
  if(phasedesc){
    const phase=parseInt(document.getElementById('idx-phase').textContent)||0;
    phasedesc.textContent=t['phase-desc-'+phase]||'Phase '+phase;
  }
}

function showTab(name,el){
  document.querySelectorAll('.tab-content').forEach(t=>t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t=>t.classList.remove('active'));
  document.getElementById('tab-'+name).classList.add('active');
  el.classList.add('active');
}

function fmt(n){if(n===undefined||n===null||n==='—')return'—';if(typeof n==='number')return n.toLocaleString();return n}
function timeAgo(ts){const d=Math.floor(Date.now()/1000)-ts;if(d<60)return d+'s ago';if(d<3600)return Math.floor(d/60)+'m ago';return Math.floor(d/3600)+'h ago'}
function short(h,s=8,e=6){return h?h.slice(0,s)+'...'+h.slice(-e):'—'}
function avatarColor(a){const c=['#4FC3F7','#00E676','#FFB300','#CE93D8','#EF5350','#4DD0E1'];return c[parseInt((a||'0x00').slice(2,4),16)%c.length]}

async function addToMetaMask(){
  if(!window.ethereum){alert('MetaMask not found');return}
  try{await window.ethereum.request({method:'wallet_addEthereumChain',params:[{chainId:'0x2329',chainName:'Aequitas Chain',nativeCurrency:{name:'AEQ',symbol:'AEQ',decimals:18},rpcUrls:['https://aequitas-production-9fba.up.railway.app/rpc'],blockExplorerUrls:['https://aequitas-production-9fba.up.railway.app']}]})}catch(e){}
}

async function loadStatus(){
  try{
    const d=await(await fetch('/api/status')).json();
    document.getElementById('s-height').textContent=fmt(d.height);
    document.getElementById('s-humans').textContent=fmt(d.total_humans);
    document.getElementById('s-supply').textContent=d.total_supply||'—';
    document.getElementById('s-index').textContent=fmt(d.index);
    document.getElementById('s-uptime').textContent=d.uptime||'—';
    document.getElementById('idx-score').textContent=fmt(d.index);
    document.getElementById('idx-gini').textContent=typeof d.gini==='number'?d.gini.toFixed(3):d.gini;
    document.getElementById('idx-supply2').textContent=d.total_supply||'—';
    document.getElementById('idx-phase').textContent=fmt(d.phase);
    document.getElementById('idx-humans2').textContent=fmt(d.total_humans);
    document.getElementById('stat-humans').textContent=fmt(d.total_humans);
    document.getElementById('stat-supply').textContent=d.total_supply||'—';
    if(d.index!==undefined){
      document.getElementById('idx-bar').style.width=Math.min(d.index,100)+'%';
      const t=T[currentLang]||T.en;
      document.getElementById('idx-phase-desc').textContent=t['phase-desc-'+(d.phase||0)]||'Phase '+(d.phase||0);
    }
  }catch(e){}
}

async function loadBlocks(){
  try{
    const blocks=await(await fetch('/api/blocks?limit=50')).json();
    const list=document.getElementById('blocks-list');
    if(!blocks||!blocks.length){list.innerHTML='<div class="empty">No blocks yet</div>';return}
    document.getElementById('block-count').textContent=blocks.length+' blocks';
    list.innerHTML=blocks.map(b=>{
      const merge=b.parent_hashes&&b.parent_hashes.length>1;
      const hasTx=b.transactions&&b.transactions.length>0;
      return ` + "`" + `<div class="block-item"><div class="block-num">#${b.height}</div><div><div class="block-hash">${short(b.hash)}${merge?'<span class="badge-merge">🔀 MERGE</span>':''}${hasTx?'<span class="badge-tx">✅ TX</span>':''}</div><div class="block-parents">${b.parent_hashes?b.parent_hashes.length+' parent(s) · '+short(b.proposer,8,4):''}</div></div><div class="block-right"><div class="block-humans">${b.humans||0} humans</div><div class="block-time">${timeAgo(b.timestamp)}</div></div></div>` + "`" + `;
    }).join('');
  }catch(e){}
}

async function loadHumans(){
  try{
    const d=await(await fetch('/api/humans')).json();
    document.getElementById('human-count-badge').textContent=fmt(d.total);
    const list=document.getElementById('humans-list');
    if(!d.humans||!d.humans.length){list.innerHTML='<div class="empty">No humans registered yet.<br><br>Download the Aequitas Android App and be the first!</div>';return}
    list.innerHTML=d.humans.map(h=>{
      const color=avatarColor(h.address||'0x00');
      const init=(h.address||'??').slice(2,4).toUpperCase();
      return ` + "`" + `<div class="human-item"><div class="human-avatar" style="background:${color}20;color:${color};border-color:${color}50">${init}</div><div style="flex:1;min-width:0"><div class="human-balance">${fmt(h.balance)} AEQ</div><div class="human-addr">${h.address}</div></div><div class="human-badge">✓ HUMAN</div></div>` + "`" + `;
    }).join('');
  }catch(e){}
}

function checkProofParams(){
  const p=new URLSearchParams(window.location.search);
  const proof=p.get('proof');
  if(proof){
    try{
      proofParams=JSON.parse(decodeURIComponent(proof));
      document.getElementById('proof-box').style.display='block';
      document.getElementById('proof-val').textContent='bio: '+proofParams.bio.slice(0,15)+'... | salt: '+proofParams.salt.slice(0,10)+'...';
      document.querySelectorAll('.tab')[4].click();
    }catch(e){}
  }
}

async function connectWallet(){
  if(!window.ethereum)return;
  try{
    await addToMetaMask();
    const accounts=await window.ethereum.request({method:'eth_requestAccounts'});
    walletAddr=accounts[0];
    document.getElementById('wallet-box').style.display='block';
    document.getElementById('wallet-addr').textContent=walletAddr;
    document.getElementById('btn-register').disabled=!proofParams;
    const btn=document.getElementById('btn-connect');
    btn.textContent='✓ '+walletAddr.slice(0,10)+'...'+walletAddr.slice(-4);
    btn.style.background='var(--green)';btn.style.color='#050A14';
  }catch(e){}
}

function log(msg,type){const el=document.getElementById('reg-status');el.innerHTML+=` + "`" + `<div><span class="${type}">${msg}</span></div>` + "`" + `}

async function register(){
  if(!walletAddr||!proofParams)return;
  try{
    log('⏳ Step 1/2: Generating ZK proof...','info');
    document.getElementById('btn-register').disabled=true;
    const pr=await fetch(PROOF_SERVER+'/prove',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({bio:proofParams.bio,salt:proofParams.salt,wallet:walletAddr})});
    const pd=await pr.json();
    if(!pr.ok){log('✗ '+(pd.error||'Proof failed'),'err');document.getElementById('btn-register').disabled=false;return}
    log('✓ ZK Proof generated! Step 2/2: Registering on chain...','ok');
    const r=await fetch('/api/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({bio:proofParams.bio,salt:proofParams.salt,wallet:walletAddr})});
    const d=await r.json();
    if(!d.success){log('✗ '+d.message,'err');document.getElementById('btn-register').disabled=false;return}
    log('🎉 '+d.message+' | TX: '+d.tx_hash,'ok');
    setTimeout(()=>{window.location.href='/registered?success=true&wallet='+walletAddr},1500);
  }catch(e){log('✗ '+e.message,'err');document.getElementById('btn-register').disabled=false}
}

checkProofParams();
loadStatus();loadBlocks();loadHumans();
setInterval(loadStatus,6000);setInterval(loadBlocks,6000);setInterval(loadHumans,10000);
window.ethereum?.on('accountsChanged',a=>{
  walletAddr=a[0]||'';
  if(walletAddr){
    document.getElementById('wallet-box').style.display='block';
    document.getElementById('wallet-addr').textContent=walletAddr;
    document.getElementById('btn-register').disabled=!proofParams;
    const btn=document.getElementById('btn-connect');
    btn.textContent='✓ '+walletAddr.slice(0,10)+'...'+walletAddr.slice(-4);
    btn.style.background='var(--green)';btn.style.color='#050A14';
  }
});
</script>
</body>
</html>

`)
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
<style>
body{background:#0A0E1A;color:#C9A84C;font-family:monospace;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;flex-direction:column;gap:16px;text-align:center;padding:20px}
.logo{font-size:2rem;font-weight:bold;letter-spacing:8px}
.msg{color:#22C55E;font-size:1.4rem;font-weight:bold}
.wallet{color:#6B7A99;font-size:0.75rem;margin-top:4px;word-break:break-all}
.sub{color:#6B7A99;font-size:0.85rem;line-height:1.8;max-width:400px}
.highlight{color:#C9A84C}
.box{background:#111827;border:1px solid #1E2D45;border-radius:12px;padding:24px;max-width:420px;width:100%%}
</style>
</head>
<body>
<div class="logo">AEQUITAS</div>
<div class="box">
<div class="msg">🎉 Registered as Human!</div>
<div class="wallet">%s</div>
<div style="margin:16px 0;border-top:1px solid #1E2D45"></div>
<div class="sub">
<span class="highlight">1,000 AEQ</span> has been credited to your wallet.<br><br>
Return to the <span class="highlight">Aequitas App</span> to see your registration status.<br><br>
The app will confirm automatically.
</div>
</div>
</body>
</html>`, wallet)
}
