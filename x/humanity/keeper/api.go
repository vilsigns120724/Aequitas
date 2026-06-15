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
sepoliaStatus map[string]interface{}
	state         *ChainState
}

func NewAPIServer(bc *BlockDAG, p2p *P2PNode, k *Keeper, state *ChainState) *APIServer {
s := &APIServer{
blockchain:    bc,
p2pNode:       p2p,
keeper:        k,
startTime:     time.Now(),
sepoliaStatus: map[string]interface{}{},
		state:         state,
}
go s.syncSepoliaStatus()
return s
}

func (a *APIServer) syncSepoliaStatus() {
for {
resp, err := http.Get("https://aequitas-proof-server-production.up.railway.app/health")
if err == nil {
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
var data map[string]interface{}
if json.Unmarshal(body, &data) == nil {
a.sepoliaStatus = data
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
	evmRPC := NewEVMRPCServer(a.blockchain, a.state)
	mux.HandleFunc("/rpc", evmRPC.handleRPC)
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
.stat-sub{font-size:0.6rem;color:var(--muted);line-height:1.5}
.c-green .stat-val{color:#00E676!important}.c-green .stat-accent{background:#00E676}
.c-blue .stat-val{color:#4FC3F7!important}.c-blue .stat-accent{background:#4FC3F7}
.c-gold .stat-val{color:#FFB300!important}.c-gold .stat-accent{background:#FFB300}
.c-purple .stat-val{color:#CE93D8!important}.c-purple .stat-accent{background:#CE93D8}
.c-teal .stat-val{color:#4DD0E1!important}.c-teal .stat-accent{background:#4DD0E1}
.info-banner{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:20px;margin-bottom:20px;display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:16px}
.info-item-icon{font-size:1.2rem;margin-bottom:6px}
.info-item-title{font-size:0.68rem;color:var(--gold);font-weight:bold;margin-bottom:4px}
.info-item-text{font-size:0.65rem;color:var(--muted);line-height:1.7}
.main-grid{display:grid;grid-template-columns:1fr 320px;gap:14px;padding:0 24px 24px}
@media(max-width:860px){.main-grid{grid-template-columns:1fr}}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.sec-head{padding:13px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:#080F1E}
.sec-title{font-size:0.65rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;display:flex;align-items:center;gap:7px}
.sec-dot{width:5px;height:5px;border-radius:50%;background:var(--green)}
.sec-count{font-size:0.62rem;color:var(--muted);background:var(--card2);padding:2px 8px;border-radius:8px;border:1px solid var(--border)}
.sec-desc{padding:8px 18px;font-size:0.62rem;color:var(--muted);background:#080F1E;border-bottom:1px solid var(--border);line-height:1.7}
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
.ic-val{font-size:0.65rem;color:var(--text);text-align:right;max-width:55%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ic-val.green{color:var(--green)}.ic-val.blue{color:var(--blue)}.ic-val.gold{color:var(--gold)}.ic-val.purple{color:var(--purple)}
.mm-card{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:16px}
.mm-title{font-size:0.6rem;color:var(--blue);letter-spacing:2px;margin-bottom:12px}
.mm-row{display:flex;justify-content:space-between;padding:5px 0;border-bottom:1px solid #1A2D45}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.6rem;color:var(--muted)}
.mm-val{font-size:0.6rem;color:var(--purple)}
.mm-btn{width:100%;margin-top:10px;padding:9px;background:var(--blue);color:#050A14;border:none;border-radius:7px;cursor:pointer;font-family:monospace;font-size:0.68rem;font-weight:bold;letter-spacing:1px}
.phil-card{background:linear-gradient(135deg,#1A1200,#0D1421);border:1px solid #3A2800;border-radius:10px;padding:18px;text-align:center}
.phil-quote{font-size:0.82rem;color:var(--gold);font-style:italic;line-height:1.9;margin-bottom:5px}
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
@media(max-width:700px){.index-section{grid-template-columns:1fr}}
.idx-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.idx-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:8px}
.idx-desc{font-size:0.68rem;color:var(--muted);line-height:1.8;margin-bottom:16px}
.idx-big{font-size:2.8rem;font-weight:900;color:var(--gold);line-height:1}
.idx-lbl{font-size:0.62rem;color:var(--muted);margin-top:4px}
.bar-bg{height:8px;background:#0D1421;border-radius:4px;overflow:hidden;margin:14px 0 6px}
.bar-fill{height:100%;border-radius:4px;background:linear-gradient(90deg,var(--green),var(--gold),var(--red));transition:width 1.5s}
.bar-labels{display:flex;justify-content:space-between;font-size:0.56rem;color:var(--muted)}
.metrics-row{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin-top:14px}
.metric-box{background:#080F1E;border-radius:7px;padding:10px;text-align:center}
.metric-val{font-size:1.2rem;font-weight:bold;color:var(--gold)}
.metric-lbl{font-size:0.56rem;color:var(--muted);margin-top:3px}
.story-text{font-size:0.7rem;line-height:2}
.story-text p{margin-bottom:12px}
.net-section{padding:20px 24px 24px;display:grid;grid-template-columns:1fr 1fr;gap:14px}
@media(max-width:700px){.net-section{grid-template-columns:1fr}}
.net-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;margin-bottom:0}
.net-title{font-size:0.62rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:14px}
.node-box{background:#080F1E;border-radius:8px;padding:14px;border:1px solid var(--border);margin-bottom:10px}
.node-status{display:flex;align-items:center;gap:6px;font-size:0.7rem;color:var(--green);margin-bottom:5px;font-weight:bold}
.node-dot{width:7px;height:7px;border-radius:50%;background:var(--green);box-shadow:0 0 6px var(--green)}
.node-url{font-size:0.6rem;color:var(--muted);word-break:break-all}
.spec-table{width:100%;border-collapse:collapse}
.spec-table td{padding:8px 0;border-bottom:1px solid #0A1220;font-size:0.65rem}
.spec-table tr:last-child td{border-bottom:none}
.spec-table td:first-child{color:var(--muted);width:45%}
.spec-table td:last-child{text-align:right}
.bootstrap-box{background:#080F1E;border-radius:7px;padding:12px;font-size:0.62rem;color:var(--purple);word-break:break-all;line-height:1.8;border:1px solid var(--border)}
.reg-section{padding:20px 24px 24px;max-width:620px;margin:0 auto}
.reg-hero{background:#0D1E3A;border:1px solid #1A3A5C;border-radius:10px;padding:22px;margin-bottom:16px;text-align:center}
.reg-hero-title{font-size:1rem;font-weight:bold;color:var(--text);margin-bottom:6px}
.reg-hero-sub{font-size:0.7rem;color:var(--muted);line-height:1.8}
.app-only{background:#0D1220;border:1px solid #1A2040;border-radius:10px;padding:20px;text-align:center;margin-bottom:16px}
.app-only-icon{font-size:2rem;margin-bottom:8px}
.app-only-title{font-size:0.72rem;color:var(--purple);font-weight:bold;letter-spacing:2px;margin-bottom:8px}
.app-only-text{font-size:0.68rem;color:var(--muted);line-height:1.8}
.reg-steps{display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-bottom:16px}
@media(max-width:560px){.reg-steps{grid-template-columns:repeat(2,1fr)}}
.reg-step{background:var(--card);border:1px solid var(--border);border-radius:9px;padding:14px;text-align:center}
.step-num{width:26px;height:26px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 8px;font-weight:bold;font-size:0.75rem;color:#050A14}
.step-title{font-size:0.62rem;color:var(--text);font-weight:bold;margin-bottom:3px}
.step-desc{font-size:0.58rem;color:var(--muted);line-height:1.5}
.priv-bar{background:#0D1A0D;border:1px solid #1A3020;border-radius:7px;padding:9px 14px;margin-bottom:14px;font-size:0.67rem;color:var(--green);text-align:center}
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
</div>

<!-- EXPLORER -->
<div id="tab-explorer" class="tab-content active">
  <div class="hero">
    <div class="section-label" data-i18n="live-stats">Live Chain Statistics</div>
    <div class="stats-grid">
      <div class="stat c-blue"><div class="stat-accent"></div><div class="stat-icon">🔗</div><div class="stat-lbl" data-i18n="block-height">Block Height</div><div class="stat-val" id="s-height">—</div><div class="stat-sub" data-i18n="block-height-sub">New block every 6s · BlockDAG</div></div>
      <div class="stat c-green"><div class="stat-accent"></div><div class="stat-icon">🧬</div><div class="stat-lbl" data-i18n="verified-humans">Verified Humans</div><div class="stat-val" id="s-humans">—</div><div class="stat-sub" data-i18n="verified-humans-sub">Proof of Humanity · One person one wallet</div></div>
      <div class="stat c-gold"><div class="stat-accent"></div><div class="stat-icon">🪙</div><div class="stat-lbl" data-i18n="total-supply">Total Supply</div><div class="stat-val" id="s-supply">—</div><div class="stat-sub" data-i18n="total-supply-sub">Humans × 1,000 AEQ</div></div>
      <div class="stat c-purple"><div class="stat-accent"></div><div class="stat-icon">⚖</div><div class="stat-lbl" data-i18n="aeq-index">Aequitas Index</div><div class="stat-val" id="s-index">—</div><div class="stat-sub" data-i18n="aeq-index-sub">0 = equal · 100 = unequal</div></div>
      <div class="stat c-teal"><div class="stat-accent"></div><div class="stat-icon">⚡</div><div class="stat-lbl" data-i18n="uptime">Uptime</div><div class="stat-val" id="s-uptime" style="font-size:1.1rem">—</div><div class="stat-sub">Node v0.3.0 · 2 nodes</div></div>
    </div>
    <div class="info-banner">
      <div><div class="info-item-icon">🧬</div><div class="info-item-title" data-i18n="poh-title">Proof of Humanity</div><div class="info-item-text" data-i18n="poh-text">Every AEQ holder proves they are a unique human via biometric verification. No bots, no duplicates.</div></div>
      <div><div class="info-item-icon">⚖</div><div class="info-item-title" data-i18n="fair-title">Fair Distribution</div><div class="info-item-text" data-i18n="fair-text">Every verified human receives exactly 1,000 AEQ. No pre-mine, no investor allocation.</div></div>
      <div><div class="info-item-icon">🔗</div><div class="info-item-title" data-i18n="dag-title">BlockDAG Chain</div><div class="info-item-text" data-i18n="dag-text">A Directed Acyclic Graph allows parallel block production and higher throughput than traditional blockchains.</div></div>
      <div><div class="info-item-icon">⛽</div><div class="info-item-title" data-i18n="gasless-title">Gasless</div><div class="info-item-text" data-i18n="gasless-text">Registration is completely free. No ETH needed. If you are human, you can register.</div></div>
    </div>
  </div>
  <div class="main-grid">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="recent-blocks">Recent Blocks</span></div><div class="sec-count" id="block-count">—</div></div>
      <div class="sec-desc" data-i18n="blocks-desc">🔀 = BlockDAG merge · ✅ TX = registration transaction · Block time: 6 seconds</div>
      <div id="blocks-list"><div class="empty" data-i18n="loading">Loading...</div></div>
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
      </div>
      <div class="mm-card">
        <div class="mm-title" data-i18n="add-metamask">🦊 ADD TO METAMASK</div>
        <div class="mm-row"><span class="mm-key" data-i18n="network-name">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
        <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.55rem">...9fba.up.railway.app/rpc</span></div>
        <div class="mm-row"><span class="mm-key">Chain ID</span><span class="mm-val">9001</span></div>
        <div class="mm-row"><span class="mm-key" data-i18n="symbol">Symbol</span><span class="mm-val">AEQ</span></div>
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
    <div class="section-label" data-i18n="verified-humans">Verified Humans</div>
    <div class="info-banner">
      <div><div class="info-item-icon">🔒</div><div class="info-item-title" data-i18n="what-is-it">What is it?</div><div class="info-item-text" data-i18n="humans-what">Each address has been verified as a unique human using biometric data. The actual data never leaves the device.</div></div>
      <div><div class="info-item-icon">🧮</div><div class="info-item-title" data-i18n="how-works">How it works</div><div class="info-item-text" data-i18n="humans-how">A Groth16 Zero-Knowledge Proof is generated from your biometric hash. This proves you are human without revealing any personal data.</div></div>
      <div><div class="info-item-icon">🛡</div><div class="info-item-title" data-i18n="sybil-title">Sybil Protection</div><div class="info-item-text" data-i18n="sybil-text">Each biometric hash is stored permanently. One fingerprint = one registration = one wallet = 1,000 AEQ. This can never be circumvented.</div></div>
      <div><div class="info-item-icon">🌍</div><div class="info-item-title" data-i18n="global-title">Global Access</div><div class="info-item-text" data-i18n="global-text">Anyone with a smartphone and a fingerprint can register. No bank account, no credit card, no ETH required.</div></div>
    </div>
  </div>
  <div class="humans-section">
    <div class="section">
      <div class="sec-head"><div class="sec-title"><span class="sec-dot"></span><span data-i18n="registered-humans">Registered Humans</span></div><div class="sec-count" id="human-count-badge">0</div></div>
      <div class="sec-desc" data-i18n="humans-desc">All verified humans on the Aequitas Chain. Each received 1,000 AEQ upon registration. Permanent and non-transferable.</div>
      <div id="humans-list"><div class="empty" data-i18n="no-humans">No humans registered yet. Download the Aequitas App to be first!</div></div>
    </div>
    <div class="right-col">
      <div class="info-card">
        <div class="ic-title" data-i18n="registry-stats">Registry Stats</div>
        <div class="ic-row"><span class="ic-key" data-i18n="total-humans-stat">Total Humans</span><span class="ic-val green" id="stat-humans">0</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="total-supply">Total Supply</span><span class="ic-val gold" id="stat-supply">0 AEQ</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="grant">Grant per Human</span><span class="ic-val gold">1,000 AEQ</span></div>
        <div class="ic-row"><span class="ic-key" data-i18n="reg-fee">Registration Fee</span><span class="ic-val green" data-i18n="free">FREE</span></div>
      </div>
    </div>
  </div>
</div>

<!-- INDEX -->
<div id="tab-index" class="tab-content">
  <div class="index-section">
    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="aeq-index-title">Aequitas Index — Economic Equality Score</div>
      <div class="idx-desc" data-i18n="aeq-index-desc">The Aequitas Index measures economic equality on a scale from 0 (perfect equality) to 100 (maximum inequality). It combines the Gini coefficient with network growth metrics. The protocol automatically activates redistribution mechanisms when inequality exceeds thresholds.</div>
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
      <div class="idx-desc" data-i18n="pools-desc">When inequality thresholds are exceeded, AEQ flows into these pools automatically.</div>
      <div class="metrics-row">
        <div class="metric-box"><div class="metric-val" id="pool-v">—</div><div class="metric-lbl" data-i18n="vel-pool">Velocity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-l">—</div><div class="metric-lbl" data-i18n="liq-pool">Liquidity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-u">—</div><div class="metric-lbl" data-i18n="uni-pool">Unity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-t">—</div><div class="metric-lbl" data-i18n="treasury">Treasury</div></div>
      </div>
    </div>

    <div class="idx-card">
      <div class="idx-title" data-i18n="phases-title">Protocol Phases</div>
      <div class="idx-desc" data-i18n="phases-desc">The Aequitas protocol evolves through phases as the network grows.</div>
      <table class="spec-table">
        <tr><td>Phase 0</td><td style="color:var(--green)" data-i18n="phase0">Bootstrap — Building the network</td></tr>
        <tr><td>Phase 1</td><td style="color:var(--blue)" data-i18n="phase1">Growth — Expanding human registry</td></tr>
        <tr><td>Phase 2</td><td style="color:var(--gold)" data-i18n="phase2">Stability — Redistribution active</td></tr>
        <tr><td>Phase 3</td><td style="color:var(--purple)" data-i18n="phase3">Maturity — Full decentralization</td></tr>
      </table>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="inflation-title">Inflation & Redistribution Mechanism</div>
      <div class="story-text" data-i18n="inflation-text">
        <p><span style="color:var(--gold)">Aequitas uses a dynamic inflation model</span> that is fundamentally different from traditional cryptocurrencies. Instead of a fixed supply cap (like Bitcoin) or uncontrolled inflation (like fiat currencies), Aequitas ties its supply directly to humanity.</p>
        <p><strong style="color:var(--blue)">Base Inflation:</strong> Every new verified human registration creates exactly 1,000 AEQ. This is the only form of inflation in Phase 0 and Phase 1. The supply grows only when humanity grows — never arbitrarily.</p>
        <p><strong style="color:var(--gold)">Wealth Cap:</strong> When a single wallet exceeds a dynamically calculated wealth cap, excess AEQ is automatically redistributed to the four pools: Velocity (transaction incentives), Liquidity (market stability), Unity (new registrations), and Treasury (protocol development).</p>
        <p><strong style="color:var(--purple)">Dynamic Redistribution:</strong> In Phase 2 and beyond, the protocol runs automatic cycles that analyze on-chain wealth distribution using the Gini coefficient. If the Gini exceeds 0.35, redistribution mechanisms activate. The higher the inequality, the more aggressive the redistribution.</p>
        <p><strong style="color:var(--green)">Velocity Incentives:</strong> The Velocity Pool rewards active economic participation. Wallets that transact regularly receive small AEQ bonuses, encouraging money to flow through the economy rather than being hoarded.</p>
        <p><strong style="color:var(--teal)">Mathematical Governance:</strong> All these mechanisms are encoded in the smart contract. No human can override them. Mathematics, not politics, governs the Aequitas money supply.</p>
      </div>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="gini-title">The Gini Coefficient — Measuring Inequality</div>
      <div class="idx-desc" data-i18n="gini-desc">The Gini coefficient is the most widely used measure of economic inequality. It ranges from 0 (perfect equality) to 1 (one person owns everything). Aequitas continuously monitors this metric to trigger automatic redistribution.</div>
      <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:8px;margin-top:14px">
        <div class="metric-box" style="border:1px solid #1A4A2A"><div class="metric-val" style="color:var(--green)">0.00</div><div class="metric-lbl" data-i18n="gini-0">Perfect Equality</div><div style="font-size:0.58rem;color:var(--muted);margin-top:3px" data-i18n="gini-0-sub">Everyone equal</div></div>
        <div class="metric-box" style="border:1px solid #1A2D45"><div class="metric-val" style="color:var(--blue)">0.27</div><div class="metric-lbl" data-i18n="gini-1">Low Inequality</div><div style="font-size:0.58rem;color:var(--muted);margin-top:3px" data-i18n="gini-1-sub">Scandinavia avg.</div></div>
        <div class="metric-box" style="border:1px solid #3A2800"><div class="metric-val" style="color:var(--gold)">0.41</div><div class="metric-lbl" data-i18n="gini-2">Moderate</div><div style="font-size:0.58rem;color:var(--muted);margin-top:3px" data-i18n="gini-2-sub">USA average</div></div>
        <div class="metric-box" style="border:1px solid #4A1A1A"><div class="metric-val" style="color:var(--red)">0.63</div><div class="metric-lbl" data-i18n="gini-3">High Inequality</div><div style="font-size:0.58rem;color:var(--muted);margin-top:3px" data-i18n="gini-3-sub">South Africa</div></div>
      </div>
    </div>

    <div class="idx-card" style="grid-column:1/-1">
      <div class="idx-title" data-i18n="story-title">The Story of Aequitas</div>
      <div class="story-text" data-i18n="story-text">
        <p>The global financial system was not designed for equality. Today, the richest 1% own more wealth than the bottom 50% of humanity combined. Traditional cryptocurrencies like Bitcoin replicated this problem — early adopters accumulated vast wealth while latecomers were priced out.</p>
        <p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was founded on a single principle: <em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p>
        <p>Every verified human receives exactly 1,000 AEQ. The total supply always equals verified humans × 1,000. The money supply grows with humanity itself — not with speculation, not with mining, not with printing.</p>
        <p>The verification system uses <span style="color:var(--blue)">Groth16 Zero-Knowledge Proofs</span> — your fingerprint never leaves your device. Only a mathematical proof of humanity is transmitted. One person, one wallet, forever.</p>
        <p>The <span style="color:var(--purple)">Aequitas Index</span> continuously monitors economic equality. When wealth concentration exceeds safe thresholds, the protocol automatically redistributes through smart contract mechanisms — no human intervention required. Mathematics, not politics, governs the money supply.</p>
      </div>
    </div>
  </div>
</div>

<!-- NETWORK -->
<div id="tab-network" class="tab-content">
  <div class="net-section">
    <div class="net-card" style="grid-column:1/-1">
      <div class="net-title" data-i18n="active-nodes">Active Nodes</div>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:10px">
        <div class="node-box"><div class="node-status"><span class="node-dot"></span><span data-i18n="node1">Node 1 — Railway (Primary)</span></div><div class="node-url">aequitas-production-9fba.up.railway.app</div></div>
        <div class="node-box"><div class="node-status"><span class="node-dot"></span><span data-i18n="node2">Node 2 — Render (Secondary)</span></div><div class="node-url">aequitas-node-2.onrender.com</div></div>
      </div>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="bootstrap-title">Bootstrap Node Address</div>
      <div style="margin-bottom:10px;font-size:0.65rem;color:var(--muted)" data-i18n="bootstrap-desc">Connect to this address to join the Aequitas P2P network using libp2p:</div>
      <div class="bootstrap-box">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="tech-specs">Technical Specifications</div>
      <table class="spec-table">
        <tr><td data-i18n="chain-id">Chain ID</td><td style="color:var(--blue)">9001</td></tr>
        <tr><td>EVM</td><td style="color:var(--green)" data-i18n="evm-yes">Yes (JSON-RPC at /rpc)</td></tr>
        <tr><td data-i18n="block-time">Block Time</td><td>6 seconds</td></tr>
        <tr><td data-i18n="consensus">Consensus</td><td style="color:var(--purple)">BlockDAG + PoH</td></tr>
        <tr><td>P2P</td><td>libp2p (Go)</td></tr>
        <tr><td>ZKP</td><td>Groth16 (snarkjs)</td></tr>
        <tr><td data-i18n="storage">State Storage</td><td style="color:var(--green)">PostgreSQL</td></tr>
        <tr><td data-i18n="source">Source Code</td><td><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub ↗</a></td></tr>
      </table>
    </div>
    <div class="net-card">
      <div class="net-title" data-i18n="metamask-config">MetaMask Configuration</div>
      <table class="spec-table">
        <tr><td data-i18n="network-name">Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
        <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.58rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
        <tr><td>Chain ID</td><td style="color:var(--blue)">9001</td></tr>
        <tr><td data-i18n="symbol">Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
        <tr><td data-i18n="decimals">Decimals</td><td>18</td></tr>
      </table>
      <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:12px" data-i18n="add-network">+ ADD TO METAMASK</button>
    </div>
  </div>
</div>

<!-- REGISTER -->
<div id="tab-register" class="tab-content">
  <div class="reg-section">
    <div class="reg-hero">
      <div class="reg-hero-title" data-i18n="reg-title">🔐 Register as a Verified Human</div>
      <div class="reg-hero-sub" data-i18n="reg-sub">Join the Aequitas network and receive your 1,000 AEQ. Registration requires biometric verification via the Android app. No gas fees. Permanent.</div>
    </div>
    <div class="app-only">
      <div class="app-only-icon">📱</div>
      <div class="app-only-title" data-i18n="app-only-title">REGISTRATION VIA ANDROID APP ONLY</div>
      <div class="app-only-text" data-i18n="app-only-text">Proof of Humanity requires biometric verification on your device. Download the Aequitas App, scan your fingerprint, and your 1,000 AEQ will be granted automatically. Your biometric data <strong style="color:var(--gold)">never leaves your device</strong>.</div>
    </div>
    <div class="reg-steps">
      <div class="reg-step"><div class="step-num">1</div><div class="step-title" data-i18n="step1-title">Biometric Scan</div><div class="step-desc" data-i18n="step1-desc">Fingerprint via Hardware Secure Element — stays on device</div></div>
      <div class="reg-step"><div class="step-num">2</div><div class="step-title" data-i18n="step2-title">ZKP Generated</div><div class="step-desc" data-i18n="step2-desc">Groth16 Zero-Knowledge Proof — proves humanity without revealing data</div></div>
      <div class="reg-step"><div class="step-num">3</div><div class="step-title" data-i18n="step3-title">Connect Wallet</div><div class="step-desc" data-i18n="step3-desc">Connect MetaMask or any Web3 wallet to receive your AEQ</div></div>
      <div class="reg-step"><div class="step-num">4</div><div class="step-title" data-i18n="step4-title">1,000 AEQ</div><div class="step-desc" data-i18n="step4-desc">Instantly credited. No gas fees. No waiting. Permanent.</div></div>
    </div>
    <div class="priv-bar" data-i18n="priv-bar">🔒 Hardware Secure Element · Groth16 ZKP · No gas fees · Permanent Sybil protection</div>
    <div class="reg-card">
      <div class="wallet-box" id="wallet-box"><div class="wallet-lbl" data-i18n="connected-wallet">CONNECTED WALLET</div><div class="wallet-addr" id="wallet-addr">—</div></div>
      <div class="proof-box" id="proof-box"><div class="proof-lbl" data-i18n="proof-detected">⚡ PROOF PARAMETERS DETECTED FROM APP</div><div class="proof-val" id="proof-val">—</div></div>
      <button class="reg-btn btn-connect" id="btn-connect" onclick="connectWallet()" data-i18n="connect-btn">🦊 CONNECT METAMASK</button>
      <button class="reg-btn btn-register" id="btn-register" onclick="register()" disabled data-i18n="register-btn">🔐 REGISTER ON-CHAIN</button>
      <div class="reg-log" id="reg-status"><span class="info" data-i18n="reg-hint">// Open Aequitas Android App to generate your proof...</span></div>
    </div>
    <div class="info-card">
      <div class="ic-title" data-i18n="reg-details">Registration Details</div>
      <div class="ic-row"><span class="ic-key" data-i18n="network">Network</span><span class="ic-val purple">Aequitas Chain (BlockDAG)</span></div>
      <div class="ic-row"><span class="ic-key">Chain ID</span><span class="ic-val gold">9001</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="grant">Grant Amount</span><span class="ic-val gold">1,000 AEQ</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="reg-fee">Gas Fee</span><span class="ic-val green" data-i18n="free">FREE</span></div>
      <div class="ic-row"><span class="ic-key" data-i18n="reg-limit">Registrations</span><span class="ic-val" data-i18n="reg-limit-val">Once per human · permanent</span></div>
    </div>
  </div>
</div>

<script>
const PROOF_SERVER='https://aequitas-proof-server-production.up.railway.app';
let walletAddr='',proofParams=null;
let currentLang='en';

const T={
en:{
'live':'LIVE','tab-explorer':'🔍 Explorer','tab-humans':'👥 Humans','tab-index':'📊 Index','tab-network':'🌐 Network','tab-register':'🔐 Register',
'live-stats':'Live Chain Statistics','block-height':'Block Height','block-height-sub':'New block every 6s · BlockDAG','verified-humans':'Verified Humans','verified-humans-sub':'Proof of Humanity · One person one wallet','total-supply':'Total Supply','total-supply-sub':'Humans × 1,000 AEQ','aeq-index':'Aequitas Index','aeq-index-sub':'0 = equal · 100 = unequal','uptime':'Uptime',
'poh-title':'Proof of Humanity','poh-text':'Every AEQ holder proves they are a unique human via biometric verification. No bots, no duplicates.','fair-title':'Fair Distribution','fair-text':'Every verified human receives exactly 1,000 AEQ. No pre-mine, no investor allocation.','dag-title':'BlockDAG Chain','dag-text':'A Directed Acyclic Graph allows parallel block production and higher throughput than traditional blockchains.','gasless-title':'Gasless','gasless-text':'Registration is completely free. No ETH needed. If you are human, you can register.',
'recent-blocks':'Recent Blocks','blocks-desc':'🔀 = BlockDAG merge · ✅ TX = registration transaction · Block time: 6 seconds','loading':'Loading...','network-info':'Network Info','chain-name':'Chain Name','symbol':'Symbol','block-time':'Block Time','consensus':'Consensus','nodes':'Active Nodes','add-metamask':'🦊 ADD TO METAMASK','network-name':'Network Name','add-network':'+ ADD AEQUITAS NETWORK','philosophy':'"Money exists because people exist.<br>Nothing more, nothing less."','philosophy-sub':'— THE AEQUITAS PRINCIPLE —',
'what-is-it':'What is it?','humans-what':'Each address has been verified as a unique human using biometric data. The actual data never leaves the device.','how-works':'How it works','humans-how':'A Groth16 Zero-Knowledge Proof is generated from your biometric hash. This proves you are human without revealing any personal data.','sybil-title':'Sybil Protection','sybil-text':'Each biometric hash is stored permanently. One fingerprint = one registration = one wallet = 1,000 AEQ.','global-title':'Global Access','global-text':'Anyone with a smartphone and a fingerprint can register. No bank account, no credit card, no ETH required.',
'registered-humans':'Registered Humans','humans-desc':'All verified humans on the Aequitas Chain. Each received 1,000 AEQ upon registration.','no-humans':'No humans registered yet. Download the Aequitas App to be first!','registry-stats':'Registry Stats','total-humans-stat':'Total Humans','grant':'Grant per Human','reg-fee':'Gas Fee','free':'FREE',
'aeq-index-title':'Aequitas Index — Economic Equality Score','aeq-index-desc':'The Aequitas Index measures economic equality on a scale from 0 (perfect equality) to 100 (maximum inequality). It combines the Gini coefficient with network growth metrics.','current-index':'Current Index','bar-0':'0 — Perfect Equality','bar-100':'100 — Max Inequality','gini-coeff':'Gini Coefficient','phase':'Protocol Phase',
'pools-title':'Redistribution Pools','pools-desc':'When inequality thresholds are exceeded, AEQ flows into these pools automatically.','vel-pool':'Velocity Pool','liq-pool':'Liquidity Pool','uni-pool':'Unity Pool','treasury':'Treasury',
'phases-title':'Protocol Phases','phases-desc':'The Aequitas protocol evolves through phases as the network grows.','phase0':'Bootstrap — Building the network','phase1':'Growth — Expanding human registry','phase2':'Stability — Redistribution active','phase3':'Maturity — Full decentralization',
'inflation-title':'Inflation & Redistribution Mechanism',
'gini-title':'The Gini Coefficient — Measuring Inequality','gini-desc':'The Gini coefficient is the most widely used measure of economic inequality. It ranges from 0 (perfect equality) to 1 (one person owns everything).','gini-0':'Perfect Equality','gini-0-sub':'Everyone equal','gini-1':'Low Inequality','gini-1-sub':'Scandinavia avg.','gini-2':'Moderate','gini-2-sub':'USA average','gini-3':'High Inequality','gini-3-sub':'South Africa',
'story-title':'The Story of Aequitas',
'active-nodes':'Active Nodes','node1':'Node 1 — Railway (Primary)','node2':'Node 2 — Render (Secondary)','bootstrap-title':'Bootstrap Node Address','bootstrap-desc':'Connect to this address to join the Aequitas P2P network using libp2p:','tech-specs':'Technical Specifications','chain-id':'Chain ID','evm-yes':'Yes (JSON-RPC at /rpc)','storage':'State Storage','source':'Source Code','metamask-config':'MetaMask Configuration','decimals':'Decimals',
'reg-title':'🔐 Register as a Verified Human','reg-sub':'Join the Aequitas network and receive your 1,000 AEQ. Registration requires biometric verification via the Android app. No gas fees. Permanent.','app-only-title':'REGISTRATION VIA ANDROID APP ONLY','app-only-text':'Proof of Humanity requires biometric verification on your device. Download the Aequitas App, scan your fingerprint, and your 1,000 AEQ will be granted automatically. Your biometric data <strong style="color:var(--gold)">never leaves your device</strong>.','step1-title':'Biometric Scan','step1-desc':'Fingerprint via Hardware Secure Element — stays on device','step2-title':'ZKP Generated','step2-desc':'Groth16 Zero-Knowledge Proof — proves humanity without revealing data','step3-title':'Connect Wallet','step3-desc':'Connect MetaMask or any Web3 wallet to receive your AEQ','step4-title':'1,000 AEQ','step4-desc':'Instantly credited. No gas fees. No waiting. Permanent.','priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · No gas fees · Permanent Sybil protection','connected-wallet':'CONNECTED WALLET','proof-detected':'⚡ PROOF PARAMETERS DETECTED FROM APP','connect-btn':'🦊 CONNECT METAMASK','register-btn':'🔐 REGISTER ON-CHAIN','reg-hint':'// Open Aequitas Android App to generate your proof...','reg-details':'Registration Details','network':'Network','reg-limit':'Registrations','reg-limit-val':'Once per human · permanent',
'phase-desc-0':'Phase 0: Bootstrap — Building the network and onboarding early humans','phase-desc-1':'Phase 1: Growth — Expanding the human registry globally','phase-desc-2':'Phase 2: Stability — Wealth redistribution mechanisms active','phase-desc-3':'Phase 3: Maturity — Full decentralization achieved'
},
de:{
'live':'LIVE','tab-explorer':'🔍 Explorer','tab-humans':'👥 Menschen','tab-index':'📊 Index','tab-network':'🌐 Netzwerk','tab-register':'🔐 Registrieren',
'live-stats':'Live Chain-Statistiken','block-height':'Blockhöhe','block-height-sub':'Neuer Block alle 6s · BlockDAG','verified-humans':'Verifizierte Menschen','verified-humans-sub':'Menschlichkeitsnachweis · Eine Person eine Wallet','total-supply':'Gesamtmenge','total-supply-sub':'Menschen × 1.000 AEQ','aeq-index':'Aequitas-Index','aeq-index-sub':'0 = gleich · 100 = ungleich','uptime':'Betriebszeit',
'poh-title':'Menschlichkeitsnachweis','poh-text':'Jeder AEQ-Inhaber beweist seine Einzigartigkeit per Biometrie. Keine Bots, keine Duplikate.','fair-title':'Faire Verteilung','fair-text':'Jeder verifizierte Mensch erhält genau 1.000 AEQ. Kein Pre-Mine, keine Investorenzuteilung.','dag-title':'BlockDAG-Chain','dag-text':'Ein gerichteter azyklischer Graph ermöglicht parallele Blockproduktion und höheren Durchsatz.','gasless-title':'Gebührenfrei','gasless-text':'Die Registrierung ist völlig kostenlos. Kein ETH nötig. Wenn du ein Mensch bist, kannst du dich registrieren.',
'recent-blocks':'Aktuelle Blöcke','blocks-desc':'🔀 = BlockDAG-Zusammenführung · ✅ TX = Registrierungstransaktion · Blockzeit: 6 Sekunden','loading':'Laden...','network-info':'Netzwerkinformationen','chain-name':'Netzwerkname','symbol':'Symbol','block-time':'Blockzeit','consensus':'Konsens','nodes':'Aktive Nodes','add-metamask':'🦊 ZU METAMASK HINZUFÜGEN','network-name':'Netzwerkname','add-network':'+ AEQUITAS-NETZWERK HINZUFÜGEN','philosophy':'"Geld existiert weil Menschen existieren.<br>Nichts mehr, nichts weniger."','philosophy-sub':'— DAS AEQUITAS-PRINZIP —',
'what-is-it':'Was ist das?','humans-what':'Jede Adresse wurde durch biometrische Daten als einzigartiger Mensch verifiziert. Die Daten verlassen niemals das Gerät.','how-works':'Wie es funktioniert','humans-how':'Ein Groth16-Zero-Knowledge-Beweis wird aus deinem biometrischen Hash generiert. Er beweist deine Menschlichkeit ohne persönliche Daten preiszugeben.','sybil-title':'Sybil-Schutz','sybil-text':'Jeder biometrische Hash wird dauerhaft gespeichert. Ein Fingerabdruck = eine Registrierung = eine Wallet = 1.000 AEQ.','global-title':'Globaler Zugang','global-text':'Jeder mit einem Smartphone und einem Fingerabdruck kann sich registrieren. Kein Bankkonto, keine Kreditkarte, kein ETH nötig.',
'registered-humans':'Registrierte Menschen','humans-desc':'Alle verifizierten Menschen auf der Aequitas Chain. Jeder erhielt 1.000 AEQ bei der Registrierung.','no-humans':'Noch keine Menschen registriert. Lade die Aequitas-App herunter um der Erste zu sein!','registry-stats':'Registrierungsstatistik','total-humans-stat':'Gesamte Menschen','grant':'Zuteilung pro Mensch','reg-fee':'Gasgebühr','free':'KOSTENLOS',
'aeq-index-title':'Aequitas-Index — Wirtschaftlicher Gleichheitswert','aeq-index-desc':'Der Aequitas-Index misst wirtschaftliche Gleichheit auf einer Skala von 0 (vollkommene Gleichheit) bis 100 (maximale Ungleichheit).','current-index':'Aktueller Index','bar-0':'0 — Vollkommene Gleichheit','bar-100':'100 — Max. Ungleichheit','gini-coeff':'Gini-Koeffizient','phase':'Protokollphase',
'pools-title':'Umverteilungspools','pools-desc':'Wenn Ungleichheitsschwellenwerte überschritten werden, fließt AEQ automatisch in diese Pools.','vel-pool':'Velocity-Pool','liq-pool':'Liquiditäts-Pool','uni-pool':'Unity-Pool','treasury':'Tresor',
'phases-title':'Protokollphasen','phases-desc':'Das Aequitas-Protokoll entwickelt sich durch Phasen wenn das Netzwerk wächst.','phase0':'Bootstrap — Netzwerk aufbauen','phase1':'Wachstum — Menschenregister erweitern','phase2':'Stabilität — Umverteilung aktiv','phase3':'Reife — Vollständige Dezentralisierung',
'inflation-title':'Inflation & Umverteilungsmechanismus',
'gini-title':'Der Gini-Koeffizient — Ungleichheit messen','gini-desc':'Der Gini-Koeffizient ist das am häufigsten verwendete Maß für wirtschaftliche Ungleichheit. Von 0 (vollkommene Gleichheit) bis 1 (eine Person besitzt alles).','gini-0':'Vollkommene Gleichheit','gini-0-sub':'Alle gleich','gini-1':'Geringe Ungleichheit','gini-1-sub':'Skandinavien','gini-2':'Moderat','gini-2-sub':'USA-Durchschnitt','gini-3':'Hohe Ungleichheit','gini-3-sub':'Südafrika',
'story-title':'Die Geschichte von Aequitas',
'active-nodes':'Aktive Nodes','node1':'Node 1 — Railway (Primär)','node2':'Node 2 — Render (Sekundär)','bootstrap-title':'Bootstrap-Node-Adresse','bootstrap-desc':'Verbinde dich mit dieser Adresse um dem Aequitas P2P-Netzwerk beizutreten:','tech-specs':'Technische Spezifikationen','chain-id':'Chain-ID','evm-yes':'Ja (JSON-RPC unter /rpc)','storage':'Zustandsspeicher','source':'Quellcode','metamask-config':'MetaMask-Konfiguration','decimals':'Dezimalstellen',
'reg-title':'🔐 Als verifizierter Mensch registrieren','reg-sub':'Tritt dem Aequitas-Netzwerk bei und erhalte deine 1.000 AEQ. Die Registrierung erfordert biometrische Verifizierung über die Android-App. Keine Gasgebühren.','app-only-title':'REGISTRIERUNG NUR ÜBER ANDROID-APP','app-only-text':'Der Menschlichkeitsnachweis erfordert biometrische Verifizierung auf deinem Gerät. Lade die Aequitas-App herunter, scanne deinen Fingerabdruck, und deine 1.000 AEQ werden automatisch gutgeschrieben. Deine biometrischen Daten <strong style="color:var(--gold)">verlassen niemals dein Gerät</strong>.','step1-title':'Biometrischer Scan','step1-desc':'Fingerabdruck via Hardware Secure Element — bleibt auf dem Gerät','step2-title':'ZKP generiert','step2-desc':'Groth16 Zero-Knowledge-Beweis — beweist Menschlichkeit ohne Daten preiszugeben','step3-title':'Wallet verbinden','step3-desc':'Verbinde MetaMask oder eine Web3-Wallet um AEQ zu erhalten','step4-title':'1.000 AEQ','step4-desc':'Sofort gutgeschrieben. Keine Gasgebühren. Dauerhaft.','priv-bar':'🔒 Hardware Secure Element · Groth16 ZKP · Keine Gasgebühren · Permanenter Sybil-Schutz','connected-wallet':'VERBUNDENE WALLET','proof-detected':'⚡ BEWEISPARAMETER VON APP ERKANNT','connect-btn':'🦊 METAMASK VERBINDEN','register-btn':'🔐 ON-CHAIN REGISTRIEREN','reg-hint':'// Öffne die Aequitas Android-App um deinen Beweis zu generieren...','reg-details':'Registrierungsdetails','network':'Netzwerk','reg-limit':'Registrierungen','reg-limit-val':'Einmalig pro Mensch · dauerhaft',
'phase-desc-0':'Phase 0: Bootstrap — Netzwerk aufbauen und erste Menschen onboarden','phase-desc-1':'Phase 1: Wachstum — Menschenregister global erweitern','phase-desc-2':'Phase 2: Stabilität — Vermögensumverteilungsmechanismen aktiv','phase-desc-3':'Phase 3: Reife — Vollständige Dezentralisierung erreicht'
},
es:{
'live':'EN VIVO','tab-explorer':'🔍 Explorador','tab-humans':'👥 Humanos','tab-index':'📊 Índice','tab-network':'🌐 Red','tab-register':'🔐 Registrar',
'live-stats':'Estadísticas en Vivo','block-height':'Altura de Bloque','block-height-sub':'Nuevo bloque cada 6s · BlockDAG','verified-humans':'Humanos Verificados','verified-humans-sub':'Prueba de Humanidad · Una persona una wallet','total-supply':'Suministro Total','total-supply-sub':'Humanos × 1,000 AEQ','aeq-index':'Índice Aequitas','aeq-index-sub':'0 = igual · 100 = desigual','uptime':'Tiempo activo',
'poh-title':'Prueba de Humanidad','poh-text':'Cada titular de AEQ prueba ser un humano único mediante verificación biométrica. Sin bots, sin duplicados.','fair-title':'Distribución Justa','fair-text':'Cada humano verificado recibe exactamente 1,000 AEQ. Sin pre-minado, sin asignación a inversores.','dag-title':'Cadena BlockDAG','dag-text':'Un Grafo Acíclico Dirigido permite producción paralela de bloques y mayor rendimiento.','gasless-title':'Sin Gas','gasless-text':'El registro es completamente gratuito. No se necesita ETH. Si eres humano, puedes registrarte.',
'recent-blocks':'Bloques Recientes','blocks-desc':'🔀 = fusión BlockDAG · ✅ TX = transacción de registro · Tiempo de bloque: 6 segundos','loading':'Cargando...','network-info':'Información de Red','chain-name':'Nombre de Red','symbol':'Símbolo','block-time':'Tiempo de Bloque','consensus':'Consenso','nodes':'Nodos Activos','add-metamask':'🦊 AGREGAR A METAMASK','network-name':'Nombre de Red','add-network':'+ AGREGAR RED AEQUITAS','philosophy':'"El dinero existe porque las personas existen.<br>Nada más, nada menos."','philosophy-sub':'— EL PRINCIPIO AEQUITAS —',
'what-is-it':'¿Qué es?','humans-what':'Cada dirección ha sido verificada como un humano único usando datos biométricos. Los datos nunca salen del dispositivo.','how-works':'Cómo funciona','humans-how':'Se genera una Prueba de Conocimiento Cero Groth16 desde tu hash biométrico. Esto prueba que eres humano sin revelar datos personales.','sybil-title':'Protección Sybil','sybil-text':'Cada hash biométrico se almacena permanentemente. Una huella = un registro = una wallet = 1,000 AEQ.','global-title':'Acceso Global','global-text':'Cualquiera con un smartphone y huella digital puede registrarse. Sin cuenta bancaria, sin tarjeta, sin ETH.',
'registered-humans':'Humanos Registrados','humans-desc':'Todos los humanos verificados en la Cadena Aequitas. Cada uno recibió 1,000 AEQ al registrarse.','no-humans':'No hay humanos registrados aún. ¡Descarga la App Aequitas para ser el primero!','registry-stats':'Estadísticas del Registro','total-humans-stat':'Total de Humanos','grant':'Bono por Humano','reg-fee':'Tarifa de Gas','free':'GRATIS',
'aeq-index-title':'Índice Aequitas — Puntuación de Igualdad Económica','aeq-index-desc':'El Índice Aequitas mide la igualdad económica en una escala del 0 al 100.','current-index':'Índice Actual','bar-0':'0 — Igualdad Perfecta','bar-100':'100 — Máx. Desigualdad','gini-coeff':'Coeficiente Gini','phase':'Fase del Protocolo',
'pools-title':'Pools de Redistribución','pools-desc':'Cuando se superan los umbrales de desigualdad, AEQ fluye automáticamente hacia estos pools.','vel-pool':'Pool Velocidad','liq-pool':'Pool Liquidez','uni-pool':'Pool Unidad','treasury':'Tesorería',
'phases-title':'Fases del Protocolo','phases-desc':'El protocolo Aequitas evoluciona a través de fases a medida que crece la red.','phase0':'Bootstrap — Construyendo la red','phase1':'Crecimiento — Expandiendo el registro','phase2':'Estabilidad — Redistribución activa','phase3':'Madurez — Descentralización completa',
'inflation-title':'Mecanismo de Inflación y Redistribución',
'gini-title':'El Coeficiente Gini — Midiendo la Desigualdad','gini-desc':'El coeficiente Gini va de 0 (igualdad perfecta) a 1 (una persona lo posee todo).','gini-0':'Igualdad Perfecta','gini-0-sub':'Todos iguales','gini-1':'Baja Desigualdad','gini-1-sub':'Promedio Escandinavia','gini-2':'Moderado','gini-2-sub':'Promedio EE.UU.','gini-3':'Alta Desigualdad','gini-3-sub':'Sudáfrica',
'story-title':'La Historia de Aequitas',
'active-nodes':'Nodos Activos','node1':'Nodo 1 — Railway (Primario)','node2':'Nodo 2 — Render (Secundario)','bootstrap-title':'Dirección del Nodo Bootstrap','bootstrap-desc':'Conéctate a esta dirección para unirte a la red P2P de Aequitas:','tech-specs':'Especificaciones Técnicas','chain-id':'ID de Cadena','evm-yes':'Sí (JSON-RPC en /rpc)','storage':'Almacenamiento de Estado','source':'Código Fuente','metamask-config':'Configuración de MetaMask','decimals':'Decimales',
'reg-title':'🔐 Regístrate como Humano Verificado','reg-sub':'Únete a la red Aequitas y recibe tus 1,000 AEQ. El registro requiere verificación biométrica mediante la app Android.','app-only-title':'REGISTRO SOLO VÍA APP ANDROID','app-only-text':'La Prueba de Humanidad requiere verificación biométrica en tu dispositivo. Descarga la App Aequitas, escanea tu huella, y tus 1,000 AEQ se acreditarán automáticamente. Tus datos biométricos <strong style="color:var(--gold)">nunca salen de tu dispositivo</strong>.','step1-title':'Escaneo Biométrico','step1-desc':'Huella via Hardware Secure Element — permanece en el dispositivo','step2-title':'ZKP Generado','step2-desc':'Prueba de Conocimiento Cero Groth16 — prueba humanidad sin revelar datos','step3-title':'Conectar Wallet','step3-desc':'Conecta MetaMask o cualquier wallet Web3 para recibir AEQ','step4-title':'1,000 AEQ','step4-desc':'Acreditado instantáneamente. Sin tarifas. Sin espera. Permanente.','priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Sin tarifas de gas · Protección Sybil permanente','connected-wallet':'WALLET CONECTADA','proof-detected':'⚡ PARÁMETROS DE PRUEBA DETECTADOS','connect-btn':'🦊 CONECTAR METAMASK','register-btn':'🔐 REGISTRAR ON-CHAIN','reg-hint':'// Abre la App Android Aequitas para generar tu prueba...','reg-details':'Detalles del Registro','network':'Red','reg-limit':'Registros','reg-limit-val':'Una vez por humano · permanente',
'phase-desc-0':'Fase 0: Bootstrap — Construyendo la red','phase-desc-1':'Fase 1: Crecimiento — Expandiendo el registro de humanos','phase-desc-2':'Fase 2: Estabilidad — Mecanismos de redistribución activos','phase-desc-3':'Fase 3: Madurez — Descentralización completa alcanzada'
},
ru:{
'live':'В ЭФИРЕ','tab-explorer':'🔍 Проводник','tab-humans':'👥 Люди','tab-index':'📊 Индекс','tab-network':'🌐 Сеть','tab-register':'🔐 Регистрация',
'live-stats':'Статистика цепочки в реальном времени','block-height':'Высота Блока','block-height-sub':'Новый блок каждые 6с · BlockDAG','verified-humans':'Верифицированных Людей','verified-humans-sub':'Доказательство человечности · Один человек одна кошелёк','total-supply':'Общее Предложение','total-supply-sub':'Люди × 1 000 AEQ','aeq-index':'Индекс Aequitas','aeq-index-sub':'0 = равенство · 100 = неравенство','uptime':'Время работы',
'poh-title':'Доказательство Человечности','poh-text':'Каждый владелец AEQ доказывает, что он уникальный человек через биометрическую верификацию. Без ботов, без дублей.','fair-title':'Справедливое Распределение','fair-text':'Каждый верифицированный человек получает ровно 1 000 AEQ. Без предварительной добычи.','dag-title':'BlockDAG Цепочка','dag-text':'Направленный ациклический граф позволяет параллельное производство блоков и более высокую пропускную способность.','gasless-title':'Без Комиссий','gasless-text':'Регистрация абсолютно бесплатна. Не нужен ETH. Если ты человек, ты можешь зарегистрироваться.',
'recent-blocks':'Последние Блоки','blocks-desc':'🔀 = слияние BlockDAG · ✅ TX = транзакция регистрации · Время блока: 6 секунд','loading':'Загрузка...','network-info':'Информация о Сети','chain-name':'Название Сети','symbol':'Символ','block-time':'Время Блока','consensus':'Консенсус','nodes':'Активные Ноды','add-metamask':'🦊 ДОБАВИТЬ В METAMASK','network-name':'Название Сети','add-network':'+ ДОБАВИТЬ СЕТЬ AEQUITAS','philosophy':'"Деньги существуют потому что существуют люди.<br>Ничего больше, ничего меньше."','philosophy-sub':'— ПРИНЦИП AEQUITAS —',
'what-is-it':'Что это?','humans-what':'Каждый адрес верифицирован как уникальный человек с помощью биометрических данных. Данные никогда не покидают устройство.','how-works':'Как это работает','humans-how':'Доказательство с нулевым разглашением Groth16 генерируется из биометрического хэша. Это доказывает человечность без раскрытия личных данных.','sybil-title':'Защита от Сибилл','sybil-text':'Каждый биометрический хэш хранится постоянно. Один отпечаток = одна регистрация = один кошелёк = 1 000 AEQ.','global-title':'Глобальный Доступ','global-text':'Любой со смартфоном и отпечатком пальца может зарегистрироваться. Без банковского счёта, без кредитной карты, без ETH.',
'registered-humans':'Зарегистрированных Людей','humans-desc':'Все верифицированные люди в цепочке Aequitas. Каждый получил 1 000 AEQ при регистрации.','no-humans':'Людей ещё нет. Скачай приложение Aequitas чтобы стать первым!','registry-stats':'Статистика Реестра','total-humans-stat':'Всего Людей','grant':'Грант на Человека','reg-fee':'Комиссия Газа','free':'БЕСПЛАТНО',
'aeq-index-title':'Индекс Aequitas — Оценка Экономического Равенства','aeq-index-desc':'Индекс Aequitas измеряет экономическое равенство по шкале от 0 до 100.','current-index':'Текущий Индекс','bar-0':'0 — Полное Равенство','bar-100':'100 — Макс. Неравенство','gini-coeff':'Коэффициент Джини','phase':'Фаза Протокола',
'pools-title':'Пулы Перераспределения','pools-desc':'Когда пороги неравенства превышены, AEQ автоматически поступает в эти пулы.','vel-pool':'Пул Скорости','liq-pool':'Пул Ликвидности','uni-pool':'Пул Единства','treasury':'Казначейство',
'phases-title':'Фазы Протокола','phases-desc':'Протокол Aequitas развивается через фазы по мере роста сети.','phase0':'Загрузка — Построение сети','phase1':'Рост — Расширение реестра','phase2':'Стабильность — Перераспределение активно','phase3':'Зрелость — Полная децентрализация',
'inflation-title':'Механизм Инфляции и Перераспределения',
'gini-title':'Коэффициент Джини — Измерение Неравенства','gini-desc':'Коэффициент Джини от 0 (полное равенство) до 1 (один человек владеет всем).','gini-0':'Полное Равенство','gini-0-sub':'Все равны','gini-1':'Низкое Неравенство','gini-1-sub':'Скандинавия','gini-2':'Умеренное','gini-2-sub':'Средн. США','gini-3':'Высокое Неравенство','gini-3-sub':'ЮАР',
'story-title':'История Aequitas',
'active-nodes':'Активные Ноды','node1':'Нода 1 — Railway (Основная)','node2':'Нода 2 — Render (Вторичная)','bootstrap-title':'Адрес Bootstrap-Ноды','bootstrap-desc':'Подключись к этому адресу чтобы присоединиться к P2P-сети Aequitas:','tech-specs':'Технические Характеристики','chain-id':'ID Цепочки','evm-yes':'Да (JSON-RPC по /rpc)','storage':'Хранение Состояния','source':'Исходный Код','metamask-config':'Настройка MetaMask','decimals':'Знаков после запятой',
'reg-title':'🔐 Зарегистрируйся как Верифицированный Человек','reg-sub':'Присоединись к сети Aequitas и получи свои 1 000 AEQ. Регистрация требует биометрической верификации через Android-приложение.','app-only-title':'РЕГИСТРАЦИЯ ТОЛЬКО ЧЕРЕЗ ANDROID-ПРИЛОЖЕНИЕ','app-only-text':'Доказательство человечности требует биометрической верификации на твоём устройстве. Скачай приложение Aequitas, отсканируй отпечаток, и твои 1 000 AEQ будут начислены автоматически. Биометрические данные <strong style="color:var(--gold)">никогда не покидают устройство</strong>.','step1-title':'Биометрический Скан','step1-desc':'Отпечаток через Hardware Secure Element — остаётся на устройстве','step2-title':'ZKP Сгенерировано','step2-desc':'Доказательство Groth16 — доказывает человечность без раскрытия данных','step3-title':'Подключить Кошелёк','step3-desc':'Подключи MetaMask или любой Web3-кошелёк для получения AEQ','step4-title':'1 000 AEQ','step4-desc':'Начислено мгновенно. Без комиссий. Без ожидания. Постоянно.','priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Без комиссий · Постоянная защита от Сибилл','connected-wallet':'ПОДКЛЮЧЁННЫЙ КОШЕЛЁК','proof-detected':'⚡ ПАРАМЕТРЫ ДОКАЗАТЕЛЬСТВА ОБНАРУЖЕНЫ','connect-btn':'🦊 ПОДКЛЮЧИТЬ METAMASK','register-btn':'🔐 ЗАРЕГИСТРИРОВАТЬСЯ ON-CHAIN','reg-hint':'// Открой Android-приложение Aequitas для генерации доказательства...','reg-details':'Детали Регистрации','network':'Сеть','reg-limit':'Регистрации','reg-limit-val':'Один раз на человека · постоянно',
'phase-desc-0':'Фаза 0: Загрузка — Построение сети','phase-desc-1':'Фаза 1: Рост — Глобальное расширение реестра','phase-desc-2':'Фаза 2: Стабильность — Механизмы перераспределения активны','phase-desc-3':'Фаза 3: Зрелость — Полная децентрализация достигнута'
},
zh:{
'live':'直播','tab-explorer':'🔍 浏览器','tab-humans':'👥 人类','tab-index':'📊 指数','tab-network':'🌐 网络','tab-register':'🔐 注册',
'live-stats':'链上实时统计','block-height':'区块高度','block-height-sub':'每6秒新区块 · BlockDAG','verified-humans':'已验证人类','verified-humans-sub':'人类证明 · 一人一钱包','total-supply':'总供应量','total-supply-sub':'人类 × 1,000 AEQ','aeq-index':'Aequitas指数','aeq-index-sub':'0 = 平等 · 100 = 不平等','uptime':'运行时间',
'poh-title':'人类证明','poh-text':'每个AEQ持有者通过生物特征验证证明自己是独特的人类。无机器人，无重复。','fair-title':'公平分配','fair-text':'每个经过验证的人类获得恰好1,000 AEQ。无预挖，无投资者分配。','dag-title':'BlockDAG链','dag-text':'有向无环图允许并行区块生产和比传统区块链更高的吞吐量。','gasless-title':'无Gas费','gasless-text':'注册完全免费。不需要ETH。如果你是人类，你就可以注册。',
'recent-blocks':'最近区块','blocks-desc':'🔀 = BlockDAG合并 · ✅ TX = 注册交易 · 出块时间：6秒','loading':'加载中...','network-info':'网络信息','chain-name':'网络名称','symbol':'符号','block-time':'出块时间','consensus':'共识','nodes':'活跃节点','add-metamask':'🦊 添加到METAMASK','network-name':'网络名称','add-network':'+ 添加AEQUITAS网络','philosophy':'"货币存在是因为人类存在。<br>仅此而已。"','philosophy-sub':'— AEQUITAS原则 —',
'what-is-it':'这是什么？','humans-what':'每个地址都通过生物特征数据被验证为独特的人类。实际数据永远不会离开设备。','how-works':'工作原理','humans-how':'从您的生物特征哈希生成Groth16零知识证明。这证明您是人类而不透露任何个人数据。','sybil-title':'女巫攻击防护','sybil-text':'每个生物特征哈希都永久存储。一个指纹 = 一次注册 = 一个钱包 = 1,000 AEQ。','global-title':'全球访问','global-text':'任何拥有智能手机和指纹的人都可以注册。无需银行账户、信用卡或ETH。',
'registered-humans':'已注册人类','humans-desc':'Aequitas链上所有经过验证的人类。每人在注册时获得1,000 AEQ。','no-humans':'还没有人类注册。下载Aequitas应用成为第一个！','registry-stats':'注册统计','total-humans-stat':'总人类数','grant':'每人补助','reg-fee':'Gas费','free':'免费',
'aeq-index-title':'Aequitas指数 — 经济平等分数','aeq-index-desc':'Aequitas指数在0到100的范围内衡量经济平等。','current-index':'当前指数','bar-0':'0 — 完全平等','bar-100':'100 — 最大不平等','gini-coeff':'基尼系数','phase':'协议阶段',
'pools-title':'再分配池','pools-desc':'当不平等阈值被超过时，AEQ自动流入这些池。','vel-pool':'速度池','liq-pool':'流动性池','uni-pool':'团结池','treasury':'国库',
'phases-title':'协议阶段','phases-desc':'随着网络增长，Aequitas协议通过各个阶段演进。','phase0':'引导期 — 建立网络','phase1':'增长期 — 扩展注册','phase2':'稳定期 — 再分配激活','phase3':'成熟期 — 完全去中心化',
'inflation-title':'通胀与再分配机制',
'gini-title':'基尼系数 — 衡量不平等','gini-desc':'基尼系数从0（完全平等）到1（一人拥有一切）。','gini-0':'完全平等','gini-0-sub':'所有人相等','gini-1':'低不平等','gini-1-sub':'斯堪的纳维亚','gini-2':'中等','gini-2-sub':'美国平均','gini-3':'高不平等','gini-3-sub':'南非',
'story-title':'Aequitas的故事',
'active-nodes':'活跃节点','node1':'节点1 — Railway（主要）','node2':'节点2 — Render（次要）','bootstrap-title':'引导节点地址','bootstrap-desc':'连接到此地址以使用libp2p加入Aequitas P2P网络：','tech-specs':'技术规格','chain-id':'链ID','evm-yes':'是（/rpc的JSON-RPC）','storage':'状态存储','source':'源代码','metamask-config':'MetaMask配置','decimals':'小数位',
'reg-title':'🔐 注册为已验证人类','reg-sub':'加入Aequitas网络并获得您的1,000 AEQ。注册需要通过Android应用进行生物特征验证。','app-only-title':'仅通过ANDROID应用注册','app-only-text':'人类证明需要在您的设备上进行生物特征验证。下载Aequitas应用，扫描您的指纹，您的1,000 AEQ将自动发放。您的生物特征数据<strong style="color:var(--gold)">永远不会离开您的设备</strong>。','step1-title':'生物特征扫描','step1-desc':'通过硬件安全元件的指纹 — 留在设备上','step2-title':'ZKP已生成','step2-desc':'Groth16零知识证明 — 无需透露数据即可证明人类身份','step3-title':'连接钱包','step3-desc':'连接MetaMask或任何Web3钱包以接收AEQ','step4-title':'1,000 AEQ','step4-desc':'即时发放。无Gas费。无等待。永久。','priv-bar':'🔒 硬件安全元件 · Groth16 ZKP · 无Gas费 · 永久女巫攻击防护','connected-wallet':'已连接钱包','proof-detected':'⚡ 已检测到来自应用的证明参数','connect-btn':'🦊 连接METAMASK','register-btn':'🔐 链上注册','reg-hint':'// 打开Aequitas Android应用生成您的证明...','reg-details':'注册详情','network':'网络','reg-limit':'注册次数','reg-limit-val':'每人一次 · 永久',
'phase-desc-0':'阶段0：引导期 — 建立网络','phase-desc-1':'阶段1：增长期 — 全球扩展人类注册','phase-desc-2':'阶段2：稳定期 — 财富再分配机制激活','phase-desc-3':'阶段3：成熟期 — 完全去中心化实现'
},
id:{
'live':'SIARAN LANGSUNG','tab-explorer':'🔍 Penjelajah','tab-humans':'👥 Manusia','tab-index':'📊 Indeks','tab-network':'🌐 Jaringan','tab-register':'🔐 Daftar',
'live-stats':'Statistik Rantai Langsung','block-height':'Tinggi Blok','block-height-sub':'Blok baru setiap 6d · BlockDAG','verified-humans':'Manusia Terverifikasi','verified-humans-sub':'Bukti Kemanusiaan · Satu orang satu dompet','total-supply':'Total Pasokan','total-supply-sub':'Manusia × 1.000 AEQ','aeq-index':'Indeks Aequitas','aeq-index-sub':'0 = sama · 100 = tidak sama','uptime':'Waktu Aktif',
'poh-title':'Bukti Kemanusiaan','poh-text':'Setiap pemegang AEQ membuktikan dirinya adalah manusia unik melalui verifikasi biometrik. Tanpa bot, tanpa duplikat.','fair-title':'Distribusi Adil','fair-text':'Setiap manusia yang terverifikasi menerima tepat 1.000 AEQ. Tanpa pre-mine, tanpa alokasi investor.','dag-title':'Rantai BlockDAG','dag-text':'Graf Asiklik Terarah memungkinkan produksi blok paralel dan throughput lebih tinggi.','gasless-title':'Tanpa Gas','gasless-text':'Pendaftaran sepenuhnya gratis. Tidak perlu ETH. Jika kamu manusia, kamu bisa mendaftar.',
'recent-blocks':'Blok Terbaru','blocks-desc':'🔀 = penggabungan BlockDAG · ✅ TX = transaksi pendaftaran · Waktu blok: 6 detik','loading':'Memuat...','network-info':'Informasi Jaringan','chain-name':'Nama Jaringan','symbol':'Simbol','block-time':'Waktu Blok','consensus':'Konsensus','nodes':'Node Aktif','add-metamask':'🦊 TAMBAHKAN KE METAMASK','network-name':'Nama Jaringan','add-network':'+ TAMBAHKAN JARINGAN AEQUITAS','philosophy':'"Uang ada karena manusia ada.<br>Tidak lebih, tidak kurang."','philosophy-sub':'— PRINSIP AEQUITAS —',
'what-is-it':'Apa ini?','humans-what':'Setiap alamat telah diverifikasi sebagai manusia unik menggunakan data biometrik. Data aktual tidak pernah meninggalkan perangkat.','how-works':'Cara kerjanya','humans-how':'Bukti Zero-Knowledge Groth16 dihasilkan dari hash biometrik Anda. Ini membuktikan Anda adalah manusia tanpa mengungkapkan data pribadi.','sybil-title':'Perlindungan Sybil','sybil-text':'Setiap hash biometrik disimpan secara permanen. Satu sidik jari = satu pendaftaran = satu dompet = 1.000 AEQ.','global-title':'Akses Global','global-text':'Siapa pun dengan smartphone dan sidik jari bisa mendaftar. Tanpa rekening bank, tanpa kartu kredit, tanpa ETH.',
'registered-humans':'Manusia Terdaftar','humans-desc':'Semua manusia terverifikasi di Rantai Aequitas. Masing-masing menerima 1.000 AEQ saat pendaftaran.','no-humans':'Belum ada manusia terdaftar. Unduh Aplikasi Aequitas untuk menjadi yang pertama!','registry-stats':'Statistik Registri','total-humans-stat':'Total Manusia','grant':'Hibah per Manusia','reg-fee':'Biaya Gas','free':'GRATIS',
'aeq-index-title':'Indeks Aequitas — Skor Kesetaraan Ekonomi','aeq-index-desc':'Indeks Aequitas mengukur kesetaraan ekonomi pada skala 0 hingga 100.','current-index':'Indeks Saat Ini','bar-0':'0 — Kesetaraan Sempurna','bar-100':'100 — Ketidaksetaraan Maks.','gini-coeff':'Koefisien Gini','phase':'Fase Protokol',
'pools-title':'Pool Redistribusi','pools-desc':'Ketika ambang ketidaksetaraan terlampaui, AEQ mengalir ke pool ini secara otomatis.','vel-pool':'Pool Kecepatan','liq-pool':'Pool Likuiditas','uni-pool':'Pool Kesatuan','treasury':'Perbendaharaan',
'phases-title':'Fase Protokol','phases-desc':'Protokol Aequitas berkembang melalui fase seiring pertumbuhan jaringan.','phase0':'Bootstrap — Membangun jaringan','phase1':'Pertumbuhan — Memperluas registri','phase2':'Stabilitas — Redistribusi aktif','phase3':'Kedewasaan — Desentralisasi penuh',
'inflation-title':'Mekanisme Inflasi & Redistribusi',
'gini-title':'Koefisien Gini — Mengukur Ketidaksetaraan','gini-desc':'Koefisien Gini dari 0 (kesetaraan sempurna) hingga 1 (satu orang memiliki segalanya).','gini-0':'Kesetaraan Sempurna','gini-0-sub':'Semua sama','gini-1':'Ketidaksetaraan Rendah','gini-1-sub':'Rata-rata Skandinavia','gini-2':'Sedang','gini-2-sub':'Rata-rata AS','gini-3':'Ketidaksetaraan Tinggi','gini-3-sub':'Afrika Selatan',
'story-title':'Kisah Aequitas',
'active-nodes':'Node Aktif','node1':'Node 1 — Railway (Utama)','node2':'Node 2 — Render (Sekunder)','bootstrap-title':'Alamat Node Bootstrap','bootstrap-desc':'Hubungkan ke alamat ini untuk bergabung dengan jaringan P2P Aequitas:','tech-specs':'Spesifikasi Teknis','chain-id':'ID Rantai','evm-yes':'Ya (JSON-RPC di /rpc)','storage':'Penyimpanan Status','source':'Kode Sumber','metamask-config':'Konfigurasi MetaMask','decimals':'Desimal',
'reg-title':'🔐 Daftar sebagai Manusia Terverifikasi','reg-sub':'Bergabunglah dengan jaringan Aequitas dan terima 1.000 AEQ Anda. Pendaftaran memerlukan verifikasi biometrik melalui aplikasi Android.','app-only-title':'PENDAFTARAN HANYA MELALUI APLIKASI ANDROID','app-only-text':'Bukti Kemanusiaan memerlukan verifikasi biometrik di perangkat Anda. Unduh Aplikasi Aequitas, pindai sidik jari Anda, dan 1.000 AEQ Anda akan diberikan secara otomatis. Data biometrik Anda <strong style="color:var(--gold)">tidak pernah meninggalkan perangkat Anda</strong>.','step1-title':'Pemindaian Biometrik','step1-desc':'Sidik jari melalui Hardware Secure Element — tetap di perangkat','step2-title':'ZKP Dihasilkan','step2-desc':'Bukti Zero-Knowledge Groth16 — membuktikan kemanusiaan tanpa mengungkapkan data','step3-title':'Hubungkan Dompet','step3-desc':'Hubungkan MetaMask atau dompet Web3 mana pun untuk menerima AEQ','step4-title':'1.000 AEQ','step4-desc':'Dikreditkan seketika. Tanpa biaya gas. Tanpa menunggu. Permanen.','priv-bar':'🔒 Hardware Secure Element · ZKP Groth16 · Tanpa biaya gas · Perlindungan Sybil permanen','connected-wallet':'DOMPET TERHUBUNG','proof-detected':'⚡ PARAMETER BUKTI TERDETEKSI DARI APLIKASI','connect-btn':'🦊 HUBUNGKAN METAMASK','register-btn':'🔐 DAFTAR ON-CHAIN','reg-hint':'// Buka Aplikasi Android Aequitas untuk menghasilkan bukti Anda...','reg-details':'Detail Pendaftaran','network':'Jaringan','reg-limit':'Pendaftaran','reg-limit-val':'Sekali per manusia · permanen',
'phase-desc-0':'Fase 0: Bootstrap — Membangun jaringan','phase-desc-1':'Fase 1: Pertumbuhan — Memperluas registri manusia secara global','phase-desc-2':'Fase 2: Stabilitas — Mekanisme redistribusi kekayaan aktif','phase-desc-3':'Fase 3: Kedewasaan — Desentralisasi penuh tercapai'
}
};

const STORY={
en:'<p>The global financial system was not designed for equality. Today, the richest 1% own more wealth than the bottom 50% of humanity combined. Traditional cryptocurrencies like Bitcoin replicated this problem — early adopters accumulated vast wealth while latecomers were priced out.</p><p><span style="color:var(--gold)">Aequitas</span> — Latin for "fairness" and "equality" — was founded on a single principle: <em style="color:var(--gold)">"Money exists because people exist. Nothing more, nothing less."</em></p><p>Every verified human receives exactly 1,000 AEQ. The total supply always equals verified humans × 1,000. The money supply grows with humanity itself — not with speculation, not with mining, not with printing.</p><p>The verification system uses <span style="color:var(--blue)">Groth16 Zero-Knowledge Proofs</span> — your fingerprint never leaves your device. Only a mathematical proof of humanity is transmitted. One person, one wallet, forever.</p><p>The <span style="color:var(--purple)">Aequitas Index</span> continuously monitors economic equality. When wealth concentration exceeds safe thresholds, the protocol automatically redistributes — no human intervention required. Mathematics, not politics, governs the money supply.</p>',
de:'<p>Das globale Finanzsystem wurde nicht für Gleichheit entworfen. Heute besitzt das reichste 1% mehr Vermögen als die ärmsten 50% der Menschheit zusammen. Traditionelle Kryptowährungen haben dieses Problem repliziert.</p><p><span style="color:var(--gold)">Aequitas</span> — Lateinisch für "Fairness" und "Gleichheit" — wurde auf einem Prinzip gegründet: <em style="color:var(--gold)">"Geld existiert weil Menschen existieren. Nichts mehr, nichts weniger."</em></p><p>Jeder verifizierte Mensch erhält genau 1.000 AEQ. Die Gesamtmenge entspricht immer verifizierten Menschen × 1.000. Das Geldangebot wächst mit der Menschheit selbst.</p><p>Das Verifizierungssystem verwendet <span style="color:var(--blue)">Groth16 Zero-Knowledge-Beweise</span> — dein Fingerabdruck verlässt niemals dein Gerät. Eine Person, eine Wallet, für immer.</p><p>Der <span style="color:var(--purple)">Aequitas-Index</span> überwacht kontinuierlich die wirtschaftliche Gleichheit. Mathematik, nicht Politik, regiert das Geldangebot.</p>',
es:'<p>El sistema financiero global no fue diseñado para la igualdad. Hoy, el 1% más rico posee más riqueza que el 50% más pobre combinado. Las criptomonedas tradicionales replicaron este problema.</p><p><span style="color:var(--gold)">Aequitas</span> — Latín para "justicia" — se fundó en un principio: <em style="color:var(--gold)">"El dinero existe porque las personas existen. Nada más, nada menos."</em></p><p>Cada humano verificado recibe exactamente 1,000 AEQ. El suministro total siempre equivale a humanos verificados × 1,000.</p><p>El <span style="color:var(--purple)">Índice Aequitas</span> monitorea continuamente la igualdad económica. Las matemáticas, no la política, gobiernan el suministro monetario.</p>',
ru:'<p>Глобальная финансовая система не была создана для равенства. Сегодня богатейший 1% владеет большим состоянием, чем беднейшие 50% человечества вместе взятые.</p><p><span style="color:var(--gold)">Aequitas</span> — латинское слово "справедливость" — основан на принципе: <em style="color:var(--gold)">"Деньги существуют потому что существуют люди. Ничего больше, ничего меньше."</em></p><p>Каждый верифицированный человек получает ровно 1 000 AEQ. Общее предложение всегда равно верифицированным людям × 1 000.</p><p>Математика, а не политика, управляет денежным предложением Aequitas.</p>',
zh:'<p>全球金融体系并非为平等而设计。今天，最富有的1%拥有的财富超过最贫穷的50%人类的总和。</p><p><span style="color:var(--gold)">Aequitas</span>——拉丁语"公平"——建立在一个原则上：<em style="color:var(--gold)">"货币存在是因为人类存在。仅此而已。"</em></p><p>每个经过验证的人类获得恰好1,000 AEQ。总供应量始终等于已验证人类 × 1,000。</p><p>数学而非政治治理着Aequitas的货币供应。</p>',
id:'<p>Sistem keuangan global tidak dirancang untuk kesetaraan. Hari ini, 1% terkaya memiliki lebih banyak kekayaan dari 50% terbawah umat manusia.</p><p><span style="color:var(--gold)">Aequitas</span> — bahasa Latin untuk "keadilan" — didirikan pada prinsip: <em style="color:var(--gold)">"Uang ada karena manusia ada. Tidak lebih, tidak kurang."</em></p><p>Setiap manusia yang terverifikasi menerima tepat 1.000 AEQ. Total pasokan selalu sama dengan manusia terverifikasi × 1.000.</p><p>Matematika, bukan politik, mengatur pasokan uang Aequitas.</p>'
};

const INFLATION={
en:'<p><span style="color:var(--gold)">Aequitas uses a dynamic inflation model</span> fundamentally different from traditional cryptocurrencies. Instead of a fixed supply cap (like Bitcoin) or uncontrolled inflation (like fiat currencies), Aequitas ties its supply directly to humanity.</p><p><strong style="color:var(--blue)">Base Inflation:</strong> Every new verified human registration creates exactly 1,000 AEQ. This is the only form of inflation in Phase 0 and Phase 1. The supply grows only when humanity grows — never arbitrarily.</p><p><strong style="color:var(--gold)">Wealth Cap:</strong> When a single wallet exceeds a dynamically calculated wealth cap, excess AEQ is automatically redistributed to the four pools: Velocity (transaction incentives), Liquidity (market stability), Unity (new registrations), and Treasury (protocol development).</p><p><strong style="color:var(--purple)">Dynamic Redistribution:</strong> In Phase 2 and beyond, automatic cycles analyze on-chain wealth distribution using the Gini coefficient. If the Gini exceeds 0.35, redistribution activates. The higher the inequality, the more aggressive the redistribution.</p><p><strong style="color:var(--green)">Velocity Incentives:</strong> The Velocity Pool rewards active economic participation. Regular transactions receive small AEQ bonuses, encouraging money to flow rather than being hoarded.</p><p><strong style="color:var(--teal)">Mathematical Governance:</strong> All mechanisms are encoded in the smart contract. No human can override them. Mathematics, not politics, governs the money supply.</p>',
de:'<p><span style="color:var(--gold)">Aequitas verwendet ein dynamisches Inflationsmodell</span>, das sich grundlegend von traditionellen Kryptowährungen unterscheidet. Statt eines festen Versorgungsmaximums oder unkontrollierter Inflation knüpft Aequitas sein Angebot direkt an die Menschheit.</p><p><strong style="color:var(--blue)">Basisinflation:</strong> Jede neue verifizierte Menschenregistrierung erzeugt genau 1.000 AEQ. Das Angebot wächst nur wenn die Menschheit wächst — niemals willkürlich.</p><p><strong style="color:var(--gold)">Vermögensobergrenze:</strong> Wenn eine einzelne Wallet eine dynamisch berechnete Obergrenze überschreitet, wird überschüssiges AEQ automatisch in die vier Pools umverteilt: Velocity, Liquidität, Unity und Tresor.</p><p><strong style="color:var(--purple)">Dynamische Umverteilung:</strong> In Phase 2 analysieren automatische Zyklen die On-Chain-Vermögensverteilung mittels Gini-Koeffizient. Übersteigt dieser 0,35, aktivieren sich Umverteilungsmechanismen.</p><p><strong style="color:var(--teal)">Mathematische Steuerung:</strong> Alle Mechanismen sind im Smart Contract kodiert. Kein Mensch kann sie außer Kraft setzen. Mathematik, nicht Politik, regiert das Geldangebot.</p>',
es:'<p><span style="color:var(--gold)">Aequitas usa un modelo de inflación dinámico</span> fundamentalmente diferente. En lugar de un suministro máximo fijo o inflación descontrolada, Aequitas vincula su suministro directamente a la humanidad.</p><p><strong style="color:var(--blue)">Inflación Base:</strong> Cada nuevo registro humano verificado crea exactamente 1,000 AEQ. El suministro solo crece cuando la humanidad crece.</p><p><strong style="color:var(--gold)">Límite de Riqueza:</strong> Cuando una sola wallet supera un límite calculado dinámicamente, el exceso de AEQ se redistribuye automáticamente a los cuatro pools.</p><p><strong style="color:var(--purple)">Redistribución Dinámica:</strong> En la Fase 2, ciclos automáticos analizan la distribución de riqueza usando el coeficiente Gini. Si supera 0,35, se activan los mecanismos de redistribución.</p><p><strong style="color:var(--teal)">Gobernanza Matemática:</strong> Todos los mecanismos están codificados en el contrato inteligente. Nadie puede anularlos. Las matemáticas, no la política, gobiernan el suministro.</p>',
ru:'<p><span style="color:var(--gold)">Aequitas использует динамическую модель инфляции</span>, принципиально отличающуюся от традиционных криптовалют. Вместо фиксированного максимума предложения Aequitas привязывает предложение напрямую к человечеству.</p><p><strong style="color:var(--blue)">Базовая инфляция:</strong> Каждая новая регистрация человека создаёт ровно 1 000 AEQ. Предложение растёт только вместе с человечеством.</p><p><strong style="color:var(--gold)">Ограничение богатства:</strong> Когда один кошелёк превышает расчётный предел, избыточный AEQ автоматически перераспределяется в четыре пула.</p><p><strong style="color:var(--purple)">Динамическое перераспределение:</strong> В Фазе 2 автоматические циклы анализируют распределение богатства с помощью коэффициента Джини. Если он превышает 0,35, активируются механизмы перераспределения.</p><p><strong style="color:var(--teal)">Математическое управление:</strong> Все механизмы закодированы в смарт-контракте. Никто не может их отменить. Математика управляет денежным предложением.</p>',
zh:'<p><span style="color:var(--gold)">Aequitas使用动态通胀模型</span>，与传统加密货币根本不同。Aequitas将其供应量直接与人类挂钩，而不是固定供应上限。</p><p><strong style="color:var(--blue)">基础通胀：</strong>每个新的已验证人类注册创建恰好1,000 AEQ。供应量只随人类增长而增长。</p><p><strong style="color:var(--gold)">财富上限：</strong>当单个钱包超过动态计算的财富上限时，多余的AEQ自动重新分配到四个池中。</p><p><strong style="color:var(--purple)">动态再分配：</strong>在第2阶段，自动周期使用基尼系数分析链上财富分布。如果超过0.35，再分配机制激活。</p><p><strong style="color:var(--teal)">数学治理：</strong>所有机制都编码在智能合约中。没有人可以推翻它们。数学而非政治治理货币供应。</p>',
id:'<p><span style="color:var(--gold)">Aequitas menggunakan model inflasi dinamis</span> yang berbeda secara fundamental. Alih-alih batas pasokan tetap, Aequitas mengikat pasokannya langsung ke kemanusiaan.</p><p><strong style="color:var(--blue)">Inflasi Dasar:</strong> Setiap pendaftaran manusia terverifikasi baru menciptakan tepat 1.000 AEQ. Pasokan hanya tumbuh ketika kemanusiaan tumbuh.</p><p><strong style="color:var(--gold)">Batas Kekayaan:</strong> Ketika satu dompet melebihi batas yang dihitung secara dinamis, kelebihan AEQ secara otomatis didistribusikan ulang ke empat pool.</p><p><strong style="color:var(--purple)">Redistribusi Dinamis:</strong> Pada Fase 2, siklus otomatis menganalisis distribusi kekayaan menggunakan koefisien Gini. Jika melebihi 0,35, mekanisme redistribusi diaktifkan.</p><p><strong style="color:var(--teal)">Tata Kelola Matematis:</strong> Semua mekanisme dikodekan dalam kontrak pintar. Tidak ada yang bisa mengesampingkannya. Matematika mengatur pasokan uang.</p>'
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
    document.getElementById('idx-gini').textContent=fmt(d.gini);
    document.getElementById('idx-supply2').textContent=d.total_supply||'—';
    document.getElementById('idx-phase').textContent=fmt(d.phase);
    document.getElementById('idx-humans2').textContent=fmt(d.total_humans);
    document.getElementById('stat-humans').textContent=fmt(d.total_humans);
    document.getElementById('stat-supply').textContent=d.total_supply||'—';
    if(d.index!==undefined){
      document.getElementById('idx-bar').style.width=Math.min(d.index,100)+'%';
      const phases=T[currentLang]||T.en;
      document.getElementById('idx-phase-desc').textContent=phases['phase-desc-'+(d.phase||0)]||'Phase '+(d.phase||0);
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
      return ` + "`" + `<div class="block-item"><div class="block-num">#${b.height}</div><div><div class="block-hash">${short(b.hash)}${merge?'<span class="badge-merge">🔀 MERGE</span>':''}${hasTx?'<span class="badge-tx">✅ TX</span>':''}</div><div class="block-parents">${b.parent_hashes?b.parent_hashes.length+' parent(s)':''}</div></div><div class="block-right"><div class="block-humans">${b.humans||0} humans</div><div class="block-time">${timeAgo(b.timestamp)}</div></div></div>` + "`" + `;
    }).join('');
  }catch(e){}
}

async function loadHumans(){
  try{
    const d=await(await fetch('/api/humans')).json();
    document.getElementById('human-count-badge').textContent=fmt(d.total);
    const list=document.getElementById('humans-list');
    if(!d.humans||!d.humans.length){list.innerHTML='<div class="empty" data-i18n="no-humans">No humans registered yet.</div>';return}
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
  if(!window.ethereum){return}
  try{
    await addToMetaMask();
    const accounts=await window.ethereum.request({method:'eth_requestAccounts'});
    walletAddr=accounts[0];
    document.getElementById('wallet-box').style.display='block';
    document.getElementById('wallet-addr').textContent=walletAddr;
    document.getElementById('btn-register').disabled=!proofParams;
    const btn=document.getElementById('btn-connect');
    btn.textContent='✓ '+walletAddr.slice(0,10)+'...'+walletAddr.slice(-4);
    btn.style.background='var(--green)';
    btn.style.color='#050A14';
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
setInterval(loadStatus,6000);
setInterval(loadBlocks,6000);
setInterval(loadHumans,10000);
window.ethereum?.on('accountsChanged',a=>{walletAddr=a[0]||'';if(walletAddr){document.getElementById('wallet-box').style.display='block';document.getElementById('wallet-addr').textContent=walletAddr;document.getElementById('btn-register').disabled=!proofParams;const btn=document.getElementById('btn-connect');btn.textContent='✓ '+walletAddr.slice(0,10)+'...'+walletAddr.slice(-4);btn.style.background='var(--green)';btn.style.color='#050A14';}});
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
