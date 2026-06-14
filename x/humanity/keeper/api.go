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
"index":        65,
"gini":         0,
"growth":       growth,
"velocity":     50,
"phase":        0,
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
header{background:linear-gradient(180deg,#080F1E 0%,#050A14 100%);border-bottom:1px solid #1A2D45;padding:0 32px;position:sticky;top:0;z-index:100;display:flex;align-items:center;justify-content:space-between;height:64px}
.logo-wrap{display:flex;align-items:center;gap:16px}
.logo-icon{width:36px;height:36px;background:linear-gradient(135deg,#FFB300,#FF6F00);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:20px;box-shadow:0 0 20px #FFB30040}
.logo-text{font-size:1.3rem;font-weight:900;color:var(--gold);letter-spacing:6px}
.logo-sub{font-size:0.6rem;color:var(--muted);letter-spacing:4px;margin-top:1px}
.header-right{display:flex;gap:10px;align-items:center}
.badge{display:flex;align-items:center;gap:6px;padding:6px 14px;border-radius:20px;font-size:0.7rem;letter-spacing:1px}
.badge-live{background:#00E67615;border:1px solid #00E67630;color:var(--green)}
.badge-dag{background:#4FC3F715;border:1px solid #4FC3F730;color:var(--blue)}
.pulse{width:7px;height:7px;border-radius:50%;background:var(--green);animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.5;transform:scale(0.8)}}
.tabs{background:#080F1E;border-bottom:1px solid var(--border);padding:0 32px;display:flex;gap:0;overflow-x:auto}
.tab{padding:14px 22px;font-size:0.72rem;color:var(--muted);cursor:pointer;border-bottom:2px solid transparent;letter-spacing:1.5px;text-transform:uppercase;white-space:nowrap;transition:all 0.2s}
.tab:hover{color:var(--text)}
.tab.active{color:var(--blue);border-bottom-color:var(--blue)}
.tab-content{display:none}
.tab-content.active{display:block}
.hero{padding:32px 32px 0;background:radial-gradient(ellipse at 50% 0%,#0D1E3A 0%,transparent 70%)}
.hero-label{font-size:0.6rem;color:var(--muted);letter-spacing:4px;text-transform:uppercase;margin-bottom:20px;display:flex;align-items:center;gap:8px}
.hero-label::before{content:'';display:inline-block;width:20px;height:1px;background:var(--muted)}
.hero-label::after{content:'';display:inline-block;width:20px;height:1px;background:var(--muted)}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:1px;background:var(--border);border:1px solid var(--border);border-radius:12px;overflow:hidden;margin-bottom:28px}
.stat{background:var(--card);padding:24px 20px;position:relative;overflow:hidden;transition:background 0.2s}
.stat:hover{background:var(--card2)}
.stat-accent{position:absolute;top:0;left:0;right:0;height:2px}
.stat-icon{font-size:1.2rem;margin-bottom:10px}
.stat-label{font-size:0.6rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:8px}
.stat-value{font-size:2rem;font-weight:900;line-height:1;margin-bottom:6px}
.stat-sub{font-size:0.65rem;color:var(--muted);line-height:1.6}
.c-green .stat-value{color:#00E676!important}.c-green .stat-accent{background:#00E676}
.c-blue .stat-value{color:#4FC3F7!important}.c-blue .stat-accent{background:#4FC3F7}
.c-gold .stat-value{color:#FFB300!important}.c-gold .stat-accent{background:#FFB300}
.c-purple .stat-value{color:#CE93D8!important}.c-purple .stat-accent{background:#CE93D8}
.c-teal .stat-value{color:#4DD0E1!important}.c-teal .stat-accent{background:#4DD0E1}
.mission-banner{background:linear-gradient(135deg,#0D1E3A,#111E2E);border:1px solid #1A3A5C;border-radius:12px;padding:28px;margin-bottom:28px;display:grid;grid-template-columns:1fr 1fr 1fr 1fr;gap:20px}
@media(max-width:800px){.mission-banner{grid-template-columns:1fr 1fr}}
@media(max-width:500px){.mission-banner{grid-template-columns:1fr}}
.mission-item{}
.mission-icon{font-size:1.5rem;margin-bottom:8px}
.mission-title{font-size:0.72rem;color:var(--gold);font-weight:bold;margin-bottom:6px;letter-spacing:1px}
.mission-text{font-size:0.68rem;color:var(--muted);line-height:1.7}
.main-grid{display:grid;grid-template-columns:1fr 340px;gap:16px;padding:0 32px 32px}
@media(max-width:900px){.main-grid{grid-template-columns:1fr}}
.section{background:var(--card);border:1px solid var(--border);border-radius:12px;overflow:hidden}
.section-head{padding:16px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;background:#080F1E}
.section-title{font-size:0.68rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;display:flex;align-items:center;gap:8px}
.section-dot{width:6px;height:6px;border-radius:50%;background:var(--green)}
.section-count{font-size:0.65rem;color:var(--muted);background:var(--card2);padding:3px 10px;border-radius:10px;border:1px solid var(--border)}
.section-desc{padding:10px 20px;font-size:0.65rem;color:var(--muted);background:#080F1E;border-bottom:1px solid var(--border);line-height:1.7}
.block-item{padding:14px 20px;border-bottom:1px solid #0D1421;display:grid;grid-template-columns:72px 1fr auto;gap:12px;align-items:center;cursor:pointer;transition:background 0.15s}
.block-item:hover{background:#0D1421}
.block-item:last-child{border-bottom:none}
.block-num{font-size:0.85rem;font-weight:bold;color:var(--blue)}
.block-info{}
.block-hash{font-size:0.7rem;color:var(--muted);margin-bottom:3px;display:flex;align-items:center;gap:6px}
.block-parents{font-size:0.62rem;color:#3A5570}
.block-right{text-align:right}
.block-humans{font-size:0.72rem;color:var(--gold);margin-bottom:2px}
.block-time{font-size:0.62rem;color:var(--green)}
.badge-merge{background:#2D1B4E;color:var(--purple);font-size:0.58rem;padding:2px 6px;border-radius:4px;border:1px solid #4A2D7A}
.badge-tx{background:#0D2A1A;color:var(--green);font-size:0.58rem;padding:2px 6px;border-radius:4px;border:1px solid #1A4A2A}
.empty{padding:48px;text-align:center;color:var(--muted);font-size:0.75rem;line-height:2.5}
.right-col{display:flex;flex-direction:column;gap:14px}
.info-card{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:20px}
.info-card-title{font-size:0.65rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:16px;display:flex;align-items:center;gap:8px}
.info-row{display:flex;justify-content:space-between;align-items:center;padding:9px 0;border-bottom:1px solid #0D1421}
.info-row:last-child{border-bottom:none}
.info-key{font-size:0.68rem;color:var(--muted)}
.info-val{font-size:0.68rem;color:var(--text);text-align:right;max-width:55%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.info-val.green{color:var(--green)}
.info-val.blue{color:var(--blue)}
.info-val.gold{color:var(--gold)}
.info-val.purple{color:var(--purple)}
.mm-card{background:linear-gradient(135deg,#0D1E3A,#0A1628);border:1px solid #1A3A5C;border-radius:12px;padding:18px}
.mm-title{font-size:0.62rem;color:var(--blue);letter-spacing:2px;text-transform:uppercase;margin-bottom:14px}
.mm-row{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid #1A2D45;align-items:center}
.mm-row:last-child{border-bottom:none}
.mm-key{font-size:0.62rem;color:var(--muted)}
.mm-val{font-size:0.62rem;color:var(--purple)}
.mm-btn{width:100%;margin-top:12px;padding:10px;background:var(--blue);color:#050A14;border:none;border-radius:8px;cursor:pointer;font-family:monospace;font-size:0.72rem;font-weight:bold;letter-spacing:1px;transition:opacity 0.2s}
.mm-btn:hover{opacity:0.85}
.philosophy-card{background:linear-gradient(135deg,#1A1200,#0D1421);border:1px solid #3A2800;border-radius:12px;padding:22px;text-align:center}
.philosophy-quote{font-size:0.85rem;color:var(--gold);font-style:italic;line-height:1.9;margin-bottom:6px}
.philosophy-sub{font-size:0.62rem;color:var(--muted);letter-spacing:2px}
.humans-grid{padding:20px 32px 32px;display:grid;grid-template-columns:1fr 360px;gap:16px}
@media(max-width:900px){.humans-grid{grid-template-columns:1fr}}
.human-item{padding:14px 20px;border-bottom:1px solid #0D1421;display:flex;align-items:center;gap:14px;transition:background 0.15s}
.human-item:hover{background:#0D1421}
.human-item:last-child{border-bottom:none}
.human-avatar{width:40px;height:40px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.75rem;font-weight:bold;flex-shrink:0;border:2px solid}
.human-info{flex:1;min-width:0}
.human-balance{font-size:0.85rem;color:var(--gold);font-weight:bold;margin-bottom:2px}
.human-addr{font-size:0.68rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.human-badge{font-size:0.6rem;padding:3px 10px;border-radius:10px;flex-shrink:0;border:1px solid}
.index-wrap{padding:24px 32px 32px;display:grid;grid-template-columns:1fr 1fr;gap:16px}
@media(max-width:700px){.index-wrap{grid-template-columns:1fr}}
.index-card{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:24px}
.index-title{font-size:0.65rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:8px}
.index-desc{font-size:0.68rem;color:var(--muted);line-height:1.7;margin-bottom:20px}
.index-big{font-size:3rem;font-weight:900;color:var(--gold);margin-bottom:4px}
.index-label{font-size:0.65rem;color:var(--muted)}
.bar-bg{height:8px;background:#0D1421;border-radius:4px;overflow:hidden;margin:16px 0 8px}
.bar-fill{height:100%;border-radius:4px;background:linear-gradient(90deg,var(--green) 0%,var(--gold) 50%,var(--red) 100%);transition:width 1.5s ease}
.bar-labels{display:flex;justify-content:space-between;font-size:0.58rem;color:var(--muted)}
.metrics-row{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin-top:16px}
.metric-box{background:#080F1E;border-radius:8px;padding:12px;text-align:center}
.metric-val{font-size:1.3rem;font-weight:bold;color:var(--gold)}
.metric-label{font-size:0.58rem;color:var(--muted);margin-top:3px}
.net-wrap{padding:24px 32px 32px}
.net-card{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:24px;margin-bottom:16px}
.net-title{font-size:0.65rem;color:var(--muted);letter-spacing:2px;text-transform:uppercase;margin-bottom:16px}
.net-nodes{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:16px}
.node-box{background:#080F1E;border-radius:8px;padding:16px;border:1px solid var(--border)}
.node-status{display:flex;align-items:center;gap:6px;font-size:0.72rem;color:var(--green);margin-bottom:6px;font-weight:bold}
.node-url{font-size:0.6rem;color:var(--muted);word-break:break-all;line-height:1.5}
.node-dot{width:8px;height:8px;border-radius:50%;background:var(--green);box-shadow:0 0 8px var(--green)}
.spec-table{width:100%;border-collapse:collapse}
.spec-table td{padding:9px 0;border-bottom:1px solid #0D1421;font-size:0.68rem}
.spec-table tr:last-child td{border-bottom:none}
.spec-table td:first-child{color:var(--muted);width:45%}
.spec-table td:last-child{color:var(--text);text-align:right}
.bootstrap-box{background:#080F1E;border-radius:8px;padding:14px;font-size:0.65rem;color:var(--purple);word-break:break-all;line-height:1.8;border:1px solid #1A2D45}
.reg-wrap{padding:24px 32px;max-width:680px;margin:0 auto}
.reg-hero{background:linear-gradient(135deg,#0D1E3A,#111E2E);border:1px solid #1A3A5C;border-radius:12px;padding:28px;margin-bottom:20px;text-align:center}
.reg-hero-title{font-size:1.1rem;font-weight:bold;color:var(--text);margin-bottom:8px}
.reg-hero-sub{font-size:0.75rem;color:var(--muted);line-height:1.8}
.app-only{background:#0D1220;border:1px solid #1A2040;border-radius:10px;padding:24px;text-align:center;margin-bottom:20px}
.app-only-icon{font-size:2.5rem;margin-bottom:10px}
.app-only-title{font-size:0.78rem;color:var(--purple);font-weight:bold;letter-spacing:2px;margin-bottom:10px}
.app-only-text{font-size:0.72rem;color:var(--muted);line-height:1.9}
.reg-steps{display:grid;grid-template-columns:repeat(4,1fr);gap:10px;margin-bottom:20px}
@media(max-width:600px){.reg-steps{grid-template-columns:repeat(2,1fr)}}
.reg-step{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:16px;text-align:center}
.step-num{width:30px;height:30px;background:var(--blue);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 10px;font-weight:bold;font-size:0.8rem;color:#050A14}
.step-title{font-size:0.65rem;color:var(--text);font-weight:bold;margin-bottom:4px}
.step-desc{font-size:0.6rem;color:var(--muted);line-height:1.6}
.privacy-bar{background:#0D1A0D;border:1px solid #1A3020;border-radius:8px;padding:10px 16px;margin-bottom:16px;font-size:0.7rem;color:var(--green);text-align:center}
.reg-card{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:24px;margin-bottom:16px}
.wallet-box{background:#0D1A0D;border:1px solid #1A3020;border-radius:8px;padding:12px;margin-bottom:12px;display:none}
.wallet-label{font-size:0.6rem;color:var(--muted);margin-bottom:4px;letter-spacing:1px}
.wallet-addr{font-size:0.78rem;color:var(--green);font-weight:bold}
.proof-box{background:var(--card2);border:1px solid #3A2800;border-radius:8px;padding:12px;margin-bottom:12px;display:none}
.proof-label{font-size:0.6rem;color:var(--gold);margin-bottom:4px;letter-spacing:1px}
.proof-val{font-size:0.68rem;color:var(--muted)}
.reg-btn{width:100%;padding:16px;border-radius:8px;border:none;cursor:pointer;font-family:monospace;font-size:0.82rem;font-weight:bold;letter-spacing:1px;transition:all 0.2s;margin-bottom:10px}
.btn-connect{background:var(--blue);color:#050A14}
.btn-connect:hover{opacity:0.88}
.btn-register{background:var(--gold);color:#050A14}
.btn-register:hover{opacity:0.88}
.reg-btn:disabled{opacity:0.3;cursor:not-allowed}
.reg-log{background:#080F1E;border-radius:8px;padding:14px;font-size:0.7rem;line-height:2;min-height:60px;border:1px solid var(--border)}
.reg-log .ok{color:var(--green)}
.reg-log .err{color:var(--red)}
.reg-log .info{color:var(--gold)}
</style>
</head>
<body>

<header>
  <div class="logo-wrap">
    <div class="logo-icon">⚖</div>
    <div>
      <div class="logo-text">AEQUITAS</div>
      <div class="logo-sub">CHAIN EXPLORER</div>
    </div>
  </div>
  <div class="header-right">
    <div class="badge badge-live"><span class="pulse"></span>LIVE</div>
    <div class="badge badge-dag">● BLOCKDAG</div>
  </div>
</header>

<div class="tabs">
  <div class="tab active" onclick="showTab('explorer',this)">🔍 Explorer</div>
  <div class="tab" onclick="showTab('humans',this)">👥 Humans</div>
  <div class="tab" onclick="showTab('index',this)">📊 Index</div>
  <div class="tab" onclick="showTab('network',this)">🌐 Network</div>
  <div class="tab" onclick="showTab('register',this)">🔐 Register</div>
</div>

<!-- EXPLORER TAB -->
<div id="tab-explorer" class="tab-content active">
  <div class="hero">
    <div class="hero-label">Live Chain Statistics</div>
    <div class="stats-grid">
      <div class="stat c-blue">
        <div class="stat-accent"></div>
        <div class="stat-icon">🔗</div>
        <div class="stat-label">Block Height</div>
        <div class="stat-value" id="s-height">—</div>
        <div class="stat-sub">New block every 6 seconds<br>BlockDAG consensus</div>
      </div>
      <div class="stat c-green">
        <div class="stat-accent"></div>
        <div class="stat-icon">🧬</div>
        <div class="stat-label">Verified Humans</div>
        <div class="stat-value" id="s-humans">—</div>
        <div class="stat-sub">Biometric Proof of Humanity<br>One person · one wallet</div>
      </div>
      <div class="stat c-gold">
        <div class="stat-accent"></div>
        <div class="stat-icon">🪙</div>
        <div class="stat-label">Total Supply</div>
        <div class="stat-value" id="s-supply">—</div>
        <div class="stat-sub">Humans × 1,000 AEQ<br>Supply follows humanity</div>
      </div>
      <div class="stat c-purple">
        <div class="stat-accent"></div>
        <div class="stat-icon">⚖</div>
        <div class="stat-label">Aequitas Index</div>
        <div class="stat-value" id="s-index">—</div>
        <div class="stat-sub">0 = perfect equality<br>100 = max inequality</div>
      </div>
      <div class="stat c-teal">
        <div class="stat-accent"></div>
        <div class="stat-icon">⚡</div>
        <div class="stat-label">Uptime</div>
        <div class="stat-value" id="s-uptime" style="font-size:1.2rem">—</div>
        <div class="stat-sub">Node v0.3.0<br>2 nodes active</div>
      </div>
    </div>

    <div class="mission-banner">
      <div class="mission-item">
        <div class="mission-icon">🧬</div>
        <div class="mission-title">Proof of Humanity</div>
        <div class="mission-text">Every AEQ holder proves they are a unique human via biometric verification. No bots, no duplicates, no fake accounts. Ever.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">⚖</div>
        <div class="mission-title">Fair Distribution</div>
        <div class="mission-text">Every verified human receives exactly 1,000 AEQ. No pre-mine, no investor allocation. Total supply = verified humans × 1,000.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">🔗</div>
        <div class="mission-title">BlockDAG Chain</div>
        <div class="mission-text">A Directed Acyclic Graph allows parallel block production, higher throughput, and faster finality than traditional blockchains.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">⛽</div>
        <div class="mission-title">Gasless</div>
        <div class="mission-text">Registration is completely free. No ETH needed. If you are human, you can register. Period. No exceptions.</div>
      </div>
    </div>
  </div>

  <div class="main-grid">
    <div>
      <div class="section">
        <div class="section-head">
          <div class="section-title"><span class="section-dot"></span>Recent Blocks</div>
          <div class="section-count" id="block-count">Loading...</div>
        </div>
        <div class="section-desc">🔀 = BlockDAG merge (multiple parents) · ✅ TX = contains registration transactions · Block time: 6 seconds</div>
        <div id="blocks-list"><div class="empty">Loading blocks...</div></div>
      </div>
    </div>

    <div class="right-col">
      <div class="info-card">
        <div class="info-card-title">🌐 Network Info</div>
        <div class="info-row"><span class="info-key">Chain Name</span><span class="info-val gold">Aequitas Chain</span></div>
        <div class="info-row"><span class="info-key">Chain ID</span><span class="info-val blue">9001</span></div>
        <div class="info-row"><span class="info-key">Symbol</span><span class="info-val gold">AEQ</span></div>
        <div class="info-row"><span class="info-key">Block Time</span><span class="info-val">6 seconds</span></div>
        <div class="info-row"><span class="info-key">Consensus</span><span class="info-val purple">BlockDAG + PoH</span></div>
        <div class="info-row"><span class="info-key">Active Nodes</span><span class="info-val green">2 Online</span></div>
        <div class="info-row"><span class="info-key">ZKP System</span><span class="info-val">Groth16</span></div>
      </div>

      <div class="mm-card">
        <div class="mm-title">🦊 Add to MetaMask</div>
        <div class="mm-row"><span class="mm-key">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
        <div class="mm-row"><span class="mm-key">Chain ID</span><span class="mm-val">9001</span></div>
        <div class="mm-row"><span class="mm-key">Symbol</span><span class="mm-val">AEQ</span></div>
        <div class="mm-row"><span class="mm-key">Decimals</span><span class="mm-val">18</span></div>
        <button class="mm-btn" onclick="addToMetaMask()">+ ADD AEQUITAS NETWORK</button>
      </div>

      <div class="philosophy-card">
        <div class="philosophy-quote">"Money exists because people exist.<br>Nothing more, nothing less."</div>
        <div class="philosophy-sub">— THE AEQUITAS PRINCIPLE —</div>
      </div>
    </div>
  </div>
</div>

<!-- HUMANS TAB -->
<div id="tab-humans" class="tab-content">
  <div style="padding:32px 32px 0">
    <div class="hero-label">Verified Humans on Aequitas Chain</div>
    <div class="mission-banner" style="margin-bottom:20px">
      <div class="mission-item">
        <div class="mission-icon">🔒</div>
        <div class="mission-title">What is it?</div>
        <div class="mission-text">Each address listed here has been verified as a unique human using biometric data. The actual data never leaves the device — only a cryptographic proof is used.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">🧮</div>
        <div class="mission-title">How it works</div>
        <div class="mission-text">A Groth16 Zero-Knowledge Proof is generated from your biometric hash. This proves you are human without revealing any personal data whatsoever.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">🛡</div>
        <div class="mission-title">Sybil protection</div>
        <div class="mission-text">Each biometric hash is stored permanently. One fingerprint = one registration = one wallet = 1,000 AEQ. This can never be circumvented.</div>
      </div>
      <div class="mission-item">
        <div class="mission-icon">🌍</div>
        <div class="mission-title">Global access</div>
        <div class="mission-text">Anyone with a smartphone and a fingerprint can register. No bank account, no credit card, no ETH required. Truly inclusive by design.</div>
      </div>
    </div>
  </div>

  <div class="humans-grid">
    <div class="section">
      <div class="section-head">
        <div class="section-title"><span class="section-dot"></span>Registered Humans</div>
        <div class="section-count" id="human-count-badge">0 humans</div>
      </div>
      <div class="section-desc">All verified humans on the Aequitas Chain. Each human received 1,000 AEQ upon registration. Registration is permanent and non-transferable.</div>
      <div id="humans-list"><div class="empty">No humans registered yet.<br><br>Download the Aequitas App to register.<br>Be the first human on the chain!</div></div>
    </div>

    <div class="right-col">
      <div class="info-card">
        <div class="info-card-title">📊 Registry Stats</div>
        <div class="info-row"><span class="info-key">Total Humans</span><span class="info-val green" id="stat-humans">0</span></div>
        <div class="info-row"><span class="info-key">Total Supply</span><span class="info-val gold" id="stat-supply">0 AEQ</span></div>
        <div class="info-row"><span class="info-key">Avg Balance</span><span class="info-val gold">1,000 AEQ</span></div>
        <div class="info-row"><span class="info-key">Registration Fee</span><span class="info-val green">FREE</span></div>
        <div class="info-row"><span class="info-key">Grant per Human</span><span class="info-val gold">1,000 AEQ</span></div>
      </div>
    </div>
  </div>
</div>

<!-- INDEX TAB -->
<div id="tab-index" class="tab-content">
  <div class="index-wrap">
    <div class="index-card" style="grid-column:1/-1">
      <div class="index-title">Aequitas Index — Economic Equality Score</div>
      <div class="index-desc">The Aequitas Index measures economic equality on the chain on a scale from 0 (perfect equality) to 100 (maximum inequality). It combines the Gini coefficient with network growth metrics. The protocol automatically activates redistribution mechanisms when inequality exceeds certain thresholds — through wealth caps, inflation adjustments, and pool transfers.</div>
      <div style="display:grid;grid-template-columns:auto 1fr;gap:32px;align-items:center">
        <div>
          <div class="index-big" id="idx-score">—</div>
          <div class="index-label">Current Index</div>
        </div>
        <div>
          <div class="bar-bg"><div class="bar-fill" id="idx-bar" style="width:0%"></div></div>
          <div class="bar-labels"><span>0 — Perfect Equality</span><span>50 — Moderate</span><span>100 — Max Inequality</span></div>
          <div style="margin-top:12px;font-size:0.7rem;color:var(--muted);background:#080F1E;padding:10px;border-radius:6px" id="idx-phase-desc">Loading...</div>
        </div>
      </div>
      <div class="metrics-row">
        <div class="metric-box"><div class="metric-val" id="idx-gini">—</div><div class="metric-label">Gini Coefficient</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-supply2">—</div><div class="metric-label">Total Supply</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-phase">—</div><div class="metric-label">Protocol Phase</div></div>
        <div class="metric-box"><div class="metric-val" id="idx-humans2">—</div><div class="metric-label">Verified Humans</div></div>
      </div>
    </div>

    <div class="index-card">
      <div class="index-title">Redistribution Pools</div>
      <div class="index-desc">When inequality thresholds are exceeded, AEQ flows into these pools automatically.</div>
      <div class="metrics-row">
        <div class="metric-box"><div class="metric-val" id="pool-v">—</div><div class="metric-label">Velocity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-l">—</div><div class="metric-label">Liquidity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-u">—</div><div class="metric-label">Unity Pool</div></div>
        <div class="metric-box"><div class="metric-val" id="pool-t">—</div><div class="metric-label">Treasury</div></div>
      </div>
    </div>

    <div class="index-card">
      <div class="index-title">Protocol Phases</div>
      <div class="index-desc">The Aequitas protocol evolves through phases as the network grows.</div>
      <table class="spec-table">
        <tr><td>Phase 0</td><td style="color:var(--green)">Bootstrap — Building the network</td></tr>
        <tr><td>Phase 1</td><td style="color:var(--blue)">Growth — Expanding human registry</td></tr>
        <tr><td>Phase 2</td><td style="color:var(--gold)">Stability — Redistribution active</td></tr>
        <tr><td>Phase 3</td><td style="color:var(--purple)">Maturity — Full decentralization</td></tr>
      </table>
    </div>
  </div>
</div>

<!-- NETWORK TAB -->
<div id="tab-network" class="tab-content">
  <div class="net-wrap">
    <div class="net-card">
      <div class="net-title">Active Nodes</div>
      <div class="net-nodes">
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span>Node 1 — Railway (Primary)</div>
          <div class="node-url">aequitas-production-9fba.up.railway.app</div>
          <div style="margin-top:8px;font-size:0.62rem;color:var(--muted)">API + P2P + RPC · Block Producer</div>
        </div>
        <div class="node-box">
          <div class="node-status"><span class="node-dot"></span>Node 2 — Render (Secondary)</div>
          <div class="node-url">aequitas-node-2.onrender.com</div>
          <div style="margin-top:8px;font-size:0.62rem;color:var(--muted)">API + P2P · Block Producer + Sync</div>
        </div>
      </div>
    </div>

    <div class="net-card">
      <div class="net-title">Bootstrap Node Address</div>
      <div style="margin-bottom:12px;font-size:0.68rem;color:var(--muted)">Connect to this address to join the Aequitas P2P network using libp2p:</div>
      <div class="bootstrap-box">/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R</div>
    </div>

    <div class="net-card">
      <div class="net-title">Technical Specifications</div>
      <table class="spec-table">
        <tr><td>Chain ID</td><td style="color:var(--blue)">9001</td></tr>
        <tr><td>EVM Compatible</td><td style="color:var(--green)">Yes (JSON-RPC at /rpc)</td></tr>
        <tr><td>Block Time</td><td>6 seconds</td></tr>
        <tr><td>Consensus</td><td style="color:var(--purple)">BlockDAG + Proof of Humanity</td></tr>
        <tr><td>P2P Protocol</td><td>libp2p (Go)</td></tr>
        <tr><td>ZKP System</td><td>Groth16 (snarkjs)</td></tr>
        <tr><td>State Storage</td><td style="color:var(--green)">PostgreSQL (persistent)</td></tr>
        <tr><td>Language</td><td>Go 1.21</td></tr>
        <tr><td>Source Code</td><td style="color:var(--blue)"><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub ↗</a></td></tr>
      </table>
    </div>

    <div class="net-card">
      <div class="net-title">MetaMask RPC Configuration</div>
      <table class="spec-table">
        <tr><td>Network Name</td><td style="color:var(--gold)">Aequitas Chain</td></tr>
        <tr><td>RPC URL</td><td style="color:var(--blue);font-size:0.6rem">https://aequitas-production-9fba.up.railway.app/rpc</td></tr>
        <tr><td>Chain ID</td><td style="color:var(--blue)">9001</td></tr>
        <tr><td>Currency Symbol</td><td style="color:var(--gold)">AEQ</td></tr>
        <tr><td>Decimals</td><td>18</td></tr>
      </table>
      <button class="mm-btn" onclick="addToMetaMask()" style="margin-top:14px">+ ADD TO METAMASK</button>
    </div>
  </div>
</div>

<!-- REGISTER TAB -->
<div id="tab-register" class="tab-content">
  <div class="reg-wrap">
    <div class="reg-hero">
      <div class="reg-hero-title">🔐 Register as a Verified Human</div>
      <div class="reg-hero-sub">Join the Aequitas network and receive your 1,000 AEQ.<br>Registration requires biometric verification via the Android app.<br>No gas fees. No waiting. Permanent.</div>
    </div>

    <div class="app-only">
      <div class="app-only-icon">📱</div>
      <div class="app-only-title">REGISTRATION VIA ANDROID APP ONLY</div>
      <div class="app-only-text">Proof of Humanity requires biometric verification on your device.<br>Download the Aequitas App, scan your fingerprint,<br>and your 1,000 AEQ will be granted automatically.<br><br>Your biometric data <strong style="color:var(--gold)">never leaves your device</strong>.<br>Only a cryptographic zero-knowledge proof is transmitted.</div>
    </div>

    <div class="reg-steps">
      <div class="reg-step">
        <div class="step-num">1</div>
        <div class="step-title">Biometric Scan</div>
        <div class="step-desc">Fingerprint via Hardware Secure Element — data stays on device</div>
      </div>
      <div class="reg-step">
        <div class="step-num">2</div>
        <div class="step-title">ZKP Generated</div>
        <div class="step-desc">Groth16 Zero-Knowledge Proof — proves humanity without revealing data</div>
      </div>
      <div class="reg-step">
        <div class="step-num">3</div>
        <div class="step-title">Connect Wallet</div>
        <div class="step-desc">Connect MetaMask or any Web3 wallet to receive your AEQ</div>
      </div>
      <div class="reg-step">
        <div class="step-num">4</div>
        <div class="step-title">1,000 AEQ</div>
        <div class="step-desc">Instantly credited. No gas fees. No waiting. Permanent.</div>
      </div>
    </div>

    <div class="privacy-bar">🔒 Hardware Secure Element · Real Groth16 ZKP · No gas fees · Permanent Sybil protection · Data never leaves device</div>

    <div class="reg-card">
      <div class="wallet-box" id="wallet-box">
        <div class="wallet-label">CONNECTED WALLET</div>
        <div class="wallet-addr" id="wallet-addr">—</div>
      </div>
      <div class="proof-box" id="proof-box">
        <div class="proof-label">⚡ PROOF PARAMETERS DETECTED FROM APP</div>
        <div class="proof-val" id="proof-val">—</div>
      </div>
      <button class="reg-btn btn-connect" id="btn-connect" onclick="connectWallet()">🦊 CONNECT METAMASK</button>
      <button class="reg-btn btn-register" id="btn-register" onclick="register()" disabled>🔐 REGISTER ON-CHAIN</button>
      <div class="reg-log" id="reg-status"><span class="info">// Open Aequitas Android App to generate your proof...</span></div>
    </div>

    <div class="info-card">
      <div class="info-card-title">ℹ Registration Details</div>
      <div class="info-row"><span class="info-key">Network</span><span class="info-val purple">Aequitas Chain (BlockDAG)</span></div>
      <div class="info-row"><span class="info-key">Chain ID</span><span class="info-val gold">9001</span></div>
      <div class="info-row"><span class="info-key">Grant Amount</span><span class="info-val gold">1,000 AEQ</span></div>
      <div class="info-row"><span class="info-key">Gas Fee</span><span class="info-val green">FREE (gasless)</span></div>
      <div class="info-row"><span class="info-key">Registrations</span><span class="info-val">Once per human · permanent</span></div>
      <div class="info-row"><span class="info-key">Sybil Protection</span><span class="info-val green">Biometric ZKP</span></div>
    </div>
  </div>
</div>

<script>
const PROOF_SERVER='https://aequitas-proof-server-production.up.railway.app';
let walletAddr='',proofParams=null;

function showTab(name,el){
  document.querySelectorAll('.tab-content').forEach(t=>t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t=>t.classList.remove('active'));
  document.getElementById('tab-'+name).classList.add('active');
  el.classList.add('active');
}

function fmt(n){if(n===undefined||n===null||n==='—')return '—';if(typeof n==='number')return n.toLocaleString();return n}
function timeAgo(ts){const d=Math.floor(Date.now()/1000)-ts;if(d<60)return d+'s ago';if(d<3600)return Math.floor(d/60)+'m ago';return Math.floor(d/3600)+'h ago'}
function short(h,s=8,e=6){return h?h.slice(0,s)+'...'+h.slice(-e):'—'}

function avatarColor(addr){
  const c=['#4FC3F7','#00E676','#FFB300','#CE93D8','#EF5350','#4DD0E1'];
  return c[parseInt((addr||'0x00').slice(2,4),16)%c.length];
}

async function addToMetaMask(){
  if(!window.ethereum){alert('MetaMask not found');return}
  try{await window.ethereum.request({method:'wallet_addEthereumChain',params:[{chainId:'0x2329',chainName:'Aequitas Chain',nativeCurrency:{name:'AEQ',symbol:'AEQ',decimals:18},rpcUrls:['https://aequitas-production-9fba.up.railway.app/rpc'],blockExplorerUrls:['https://aequitas-production-9fba.up.railway.app']}]})}catch(e){console.error(e)}
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
      const phases=['Phase 0: Bootstrap — Building the network and onboarding early humans','Phase 1: Growth — Expanding the human registry globally','Phase 2: Stability — Wealth redistribution mechanisms active','Phase 3: Maturity — Full decentralization achieved'];
      document.getElementById('idx-phase-desc').textContent=phases[d.phase||0]||'Phase '+(d.phase||0);
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
      return ` + "`" + `<div class="block-item">
        <div class="block-num">#${b.height}</div>
        <div class="block-info">
          <div class="block-hash">${short(b.hash)}${merge?'<span class="badge-merge">🔀 MERGE</span>':''}${hasTx?'<span class="badge-tx">✅ TX</span>':''}</div>
          <div class="block-parents">${b.parent_hashes?b.parent_hashes.length+' parent(s) · '+short(b.proposer,8,4):''}</div>
        </div>
        <div class="block-right">
          <div class="block-humans">${b.humans||0} humans</div>
          <div class="block-time">${timeAgo(b.timestamp)}</div>
        </div>
      </div>` + "`" + `;
    }).join('');
  }catch(e){}
}

async function loadHumans(){
  try{
    const d=await(await fetch('/api/humans')).json();
    document.getElementById('human-count-badge').textContent=fmt(d.total)+' humans';
    const list=document.getElementById('humans-list');
    if(!d.humans||!d.humans.length){
      list.innerHTML='<div class="empty">No humans registered yet.<br><br>Be the first! Download the Aequitas App<br>and scan your fingerprint to register.</div>';
      return;
    }
    list.innerHTML=d.humans.map((h,i)=>{
      const color=avatarColor(h.address||'0x00');
      const init=(h.address||'??').slice(2,4).toUpperCase();
      return ` + "`" + `<div class="human-item">
        <div class="human-avatar" style="background:${color}20;color:${color};border-color:${color}40">${init}</div>
        <div class="human-info">
          <div class="human-balance">${fmt(h.balance)} AEQ</div>
          <div class="human-addr">${h.address}</div>
        </div>
        <div class="human-badge" style="background:#0D2A1A;color:#00E676;border-color:#1A4A2A">✓ HUMAN</div>
      </div>` + "`" + `;
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
      log('⚡ Proof parameters detected from app','info');
      document.querySelectorAll('.tab')[4].click();
    }catch(e){}
  }
}

async function connectWallet(){
  if(!window.ethereum){log('✗ MetaMask not found','err');return}
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
    log('✓ Wallet connected: '+walletAddr.slice(0,12)+'...','ok');
  }catch(e){log('✗ '+e.message,'err')}
}

function log(msg,type){
  const el=document.getElementById('reg-status');
  el.innerHTML+=` + "`" + `<div><span class="${type}">${msg}</span></div>` + "`" + `;
}

async function register(){
  if(!walletAddr){log('✗ Connect wallet first','err');return}
  if(!proofParams){log('✗ No proof parameters. Use the Android app first.','err');return}
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

window.ethereum?.on('accountsChanged',accounts=>{
  walletAddr=accounts[0]||'';
  if(walletAddr){
    document.getElementById('wallet-box').style.display='block';
    document.getElementById('wallet-addr').textContent=walletAddr;
    document.getElementById('btn-register').disabled=!proofParams;
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
