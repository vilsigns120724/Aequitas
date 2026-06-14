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
:root {
  --bg: #0A0E1A;
  --card: #111827;
  --border: #1E2D45;
  --green: #22C55E;
  --blue: #3B82F6;
  --gold: #C9A84C;
  --purple: #BC8CFF;
  --red: #EF4444;
  --text: #E8EDF5;
  --muted: #6B7A99;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
body { background: var(--bg); color: var(--text); font-family: 'Courier New', monospace; min-height: 100vh; }

/* HEADER */
header { display: flex; align-items: center; justify-content: space-between; padding: 16px 32px; border-bottom: 1px solid var(--border); background: #0d1117; position: sticky; top: 0; z-index: 100; }
.logo { display: flex; align-items: center; gap: 12px; }
.logo-text { font-size: 1.4rem; font-weight: bold; color: var(--gold); letter-spacing: 6px; }
.logo-sub { font-size: 0.65rem; color: var(--muted); letter-spacing: 3px; margin-top: 2px; }
.header-right { display: flex; gap: 8px; align-items: center; }
.live-badge { display: flex; align-items: center; gap: 6px; background: #0d1a0d; border: 1px solid #1a3020; padding: 5px 12px; border-radius: 20px; font-size: 0.72rem; color: var(--green); }
.pulse { width: 7px; height: 7px; background: var(--green); border-radius: 50%; animation: pulse 2s infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.3} }
.chain-badge { font-size: 0.72rem; color: var(--purple); background: #0D1220; border: 1px solid #1A2040; padding: 5px 12px; border-radius: 20px; }

/* TABS */
.tabs { display: flex; gap: 0; border-bottom: 1px solid var(--border); background: #0d1117; padding: 0 32px; overflow-x: auto; }
.tab { padding: 12px 20px; font-size: 0.75rem; color: var(--muted); cursor: pointer; border-bottom: 2px solid transparent; letter-spacing: 1px; text-transform: uppercase; transition: all 0.2s; white-space: nowrap; }
.tab:hover { color: var(--text); }
.tab.active { color: var(--blue); border-bottom-color: var(--blue); }
.tab-content { display: none; }
.tab-content.active { display: block; }

/* HERO STATS */
.hero { padding: 24px 32px 0; }
.hero-title { font-size: 0.65rem; color: var(--muted); letter-spacing: 3px; text-transform: uppercase; margin-bottom: 16px; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 12px; margin-bottom: 24px; }
.stat { background: var(--card); border: 1px solid var(--border); border-radius: 10px; padding: 18px; position: relative; overflow: hidden; }
.stat::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; }
.stat.green::before { background: var(--green); }
.stat.blue::before { background: var(--blue); }
.stat.gold::before { background: var(--gold); }
.stat.purple::before { background: var(--purple); }
.stat-label { font-size: 0.62rem; color: var(--muted); text-transform: uppercase; letter-spacing: 1.5px; margin-bottom: 8px; }
.stat-value { font-size: 1.7rem; font-weight: bold; line-height: 1; }
.stat.green .stat-value { color: var(--green); }
.stat.blue .stat-value { color: var(--blue); }
.stat.gold .stat-value { color: var(--gold); }
.stat.purple .stat-value { color: var(--purple); }
.stat-sub { font-size: 0.65rem; color: var(--muted); margin-top: 6px; line-height: 1.5; }

/* EXPLAINER BOX */
.explainer { background: #0D1220; border: 1px solid #1A2040; border-radius: 10px; padding: 20px; margin-bottom: 24px; }
.explainer-title { font-size: 0.7rem; color: var(--purple); letter-spacing: 2px; text-transform: uppercase; margin-bottom: 12px; }
.explainer-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; }
.explainer-item { }
.explainer-item-title { font-size: 0.7rem; color: var(--gold); margin-bottom: 4px; }
.explainer-item-text { font-size: 0.68rem; color: var(--muted); line-height: 1.6; }

/* MAIN GRID */
.main { display: grid; grid-template-columns: 1fr 360px; gap: 16px; padding: 0 32px 32px; }
@media(max-width: 900px) { .main { grid-template-columns: 1fr; } }

/* SECTIONS */
.section { background: var(--card); border: 1px solid var(--border); border-radius: 10px; overflow: hidden; }
.section-header { padding: 14px 20px; border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; }
.section-title { font-size: 0.7rem; color: var(--muted); text-transform: uppercase; letter-spacing: 2px; }
.section-count { font-size: 0.68rem; color: var(--muted); background: #161b22; padding: 2px 8px; border-radius: 10px; }
.section-desc { font-size: 0.65rem; color: var(--muted); padding: 10px 20px; border-bottom: 1px solid var(--border); line-height: 1.6; background: #0d1117; }

/* BLOCKS */
.block-item { padding: 12px 20px; border-bottom: 1px solid #161b22; display: grid; grid-template-columns: 70px 1fr auto; gap: 10px; align-items: center; transition: background 0.2s; }
.block-item:hover { background: #161b2250; }
.block-item:last-child { border-bottom: none; }
.block-height { color: var(--blue); font-weight: bold; font-size: 0.82rem; }
.block-info { }
.block-hash { color: var(--muted); font-size: 0.7rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; margin-bottom: 2px; }
.block-meta { font-size: 0.65rem; color: #4a5568; }
.block-right { text-align: right; }
.block-humans { color: var(--gold); font-size: 0.7rem; }
.block-time { color: var(--green); font-size: 0.65rem; margin-top: 2px; }
.merge-badge { display: inline-block; background: #2d1b4e; color: var(--purple); font-size: 0.6rem; padding: 1px 6px; border-radius: 4px; margin-left: 6px; }
.tx-badge { display: inline-block; background: #1a3a2a; color: var(--green); font-size: 0.6rem; padding: 1px 6px; border-radius: 4px; margin-left: 6px; }

/* HUMANS */
.human-item { padding: 12px 18px; border-bottom: 1px solid #161b22; display: flex; align-items: center; gap: 12px; }
.human-item:last-child { border-bottom: none; }
.human-avatar { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 0.75rem; color: white; font-weight: bold; flex-shrink: 0; border: 2px solid var(--green); }
.human-info { flex: 1; min-width: 0; }
.human-addr { font-size: 0.72rem; color: var(--muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.human-balance { font-size: 0.78rem; color: var(--gold); font-weight: bold; }
.human-badge { font-size: 0.6rem; color: var(--green); background: #1a3a2a; padding: 2px 8px; border-radius: 8px; flex-shrink: 0; }
.empty-state { padding: 40px; text-align: center; color: var(--muted); font-size: 0.75rem; line-height: 2; }

/* RIGHT PANEL */
.right-panel { display: flex; flex-direction: column; gap: 14px; }
.info-panel { background: var(--card); border: 1px solid var(--border); border-radius: 10px; padding: 18px; }
.info-title { font-size: 0.65rem; color: var(--muted); text-transform: uppercase; letter-spacing: 2px; margin-bottom: 14px; }
.info-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #161b22; }
.info-row:last-child { border-bottom: none; }
.info-key { font-size: 0.68rem; color: var(--muted); }
.info-val { font-size: 0.68rem; color: var(--text); text-align: right; max-width: 55%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.info-val.green { color: var(--green); }
.info-val.blue { color: var(--blue); }
.info-val.gold { color: var(--gold); }
.info-val.purple { color: var(--purple); }

/* METAMASK CONFIG */
.mm-config { background: #0D1220; border: 1px solid #1A2040; border-radius: 8px; padding: 14px; margin-top: 2px; }
.mm-title { font-size: 0.65rem; color: var(--purple); letter-spacing: 2px; margin-bottom: 10px; }
.mm-row { display: flex; justify-content: space-between; padding: 5px 0; border-bottom: 1px solid #1A2040; }
.mm-row:last-child { border-bottom: none; }
.mm-key { font-size: 0.65rem; color: var(--muted); }
.mm-val { font-size: 0.65rem; color: var(--purple); }
.copy-btn { font-size: 0.6rem; color: var(--blue); cursor: pointer; margin-left: 6px; background: none; border: none; font-family: monospace; }

/* PHILOSOPHY */
.philosophy { background: linear-gradient(135deg, #0D1220, #111827); border: 1px solid var(--border); border-radius: 10px; padding: 20px; text-align: center; }
.philosophy-quote { font-size: 0.85rem; color: var(--gold); font-style: italic; line-height: 1.8; margin-bottom: 8px; }
.philosophy-sub { font-size: 0.65rem; color: var(--muted); }

/* REGISTER TAB */
.reg-wrap { padding: 24px 32px; max-width: 640px; margin: 0 auto; }
.reg-card { background: var(--card); border: 1px solid var(--border); border-radius: 12px; padding: 28px; margin-bottom: 16px; }
.reg-title { font-size: 1rem; font-weight: bold; color: var(--text); margin-bottom: 8px; }
.reg-sub { font-size: 0.75rem; color: var(--muted); margin-bottom: 24px; line-height: 1.7; }
.reg-steps { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 24px; }
@media(max-width: 600px) { .reg-steps { grid-template-columns: repeat(2, 1fr); } }
.reg-step { text-align: center; padding: 14px 8px; background: #161b22; border-radius: 8px; border: 1px solid var(--border); }
.reg-step-num { width: 28px; height: 28px; background: var(--blue); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 8px; font-weight: bold; font-size: 0.8rem; color: white; }
.reg-step-title { font-size: 0.65rem; color: var(--text); margin-bottom: 4px; font-weight: bold; }
.reg-step-desc { font-size: 0.6rem; color: var(--muted); line-height: 1.5; }
.app-only-box { background: #0D1220; border: 1px solid #1A2040; border-radius: 10px; padding: 20px; text-align: center; margin-bottom: 16px; }
.app-only-icon { font-size: 2rem; margin-bottom: 8px; }
.app-only-title { font-size: 0.8rem; color: var(--purple); font-weight: bold; letter-spacing: 2px; margin-bottom: 8px; }
.app-only-text { font-size: 0.72rem; color: var(--muted); line-height: 1.7; }
.reg-btn { width: 100%; padding: 16px; border-radius: 8px; border: none; cursor: pointer; font-family: monospace; font-size: 0.85rem; font-weight: bold; letter-spacing: 1px; transition: all 0.2s; margin-bottom: 10px; }
.reg-btn-connect { background: var(--blue); color: white; }
.reg-btn-connect:hover { opacity: 0.85; }
.reg-btn-register { background: var(--gold); color: #0A0E1A; }
.reg-btn-register:hover { opacity: 0.85; }
.reg-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.reg-status { background: #161b22; border-radius: 8px; padding: 14px; font-size: 0.72rem; line-height: 1.9; min-height: 60px; }
.reg-status .ok { color: var(--green); }
.reg-status .err { color: var(--red); }
.reg-status .info { color: var(--gold); }
.reg-wallet-box { background: #0d1a0d; border: 1px solid #1a3020; border-radius: 8px; padding: 12px; margin-bottom: 12px; display: none; }
.reg-wallet-label { font-size: 0.62rem; color: var(--muted); margin-bottom: 4px; }
.reg-wallet-addr { font-size: 0.78rem; color: var(--green); }
.reg-proof-box { background: #161b22; border: 1px solid #e3b34130; border-radius: 8px; padding: 12px; margin-bottom: 12px; display: none; }
.reg-proof-label { font-size: 0.62rem; color: var(--gold); margin-bottom: 4px; }
.reg-proof-val { font-size: 0.7rem; color: var(--muted); }
.privacy-badge { background: #0d1a0d; border: 1px solid #1a3020; border-radius: 6px; padding: 10px; margin-bottom: 16px; font-size: 0.7rem; color: var(--green); text-align: center; }

/* INDEX TAB */
.index-wrap { padding: 24px 32px; }
.index-card { background: var(--card); border: 1px solid var(--border); border-radius: 10px; padding: 24px; margin-bottom: 16px; }
.index-title { font-size: 0.7rem; color: var(--muted); text-transform: uppercase; letter-spacing: 2px; margin-bottom: 8px; }
.index-desc { font-size: 0.72rem; color: var(--muted); line-height: 1.7; margin-bottom: 20px; }
.index-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 20px; }
@media(max-width: 600px) { .index-grid { grid-template-columns: repeat(2, 1fr); } }
.index-metric { text-align: center; background: #161b22; border-radius: 8px; padding: 16px; }
.index-metric-label { font-size: 0.62rem; color: var(--muted); margin-bottom: 6px; }
.index-metric-val { font-size: 1.4rem; font-weight: bold; color: var(--gold); }
.index-metric-sub { font-size: 0.6rem; color: var(--muted); margin-top: 4px; }
.index-bar-bg { height: 12px; background: #161b22; border-radius: 6px; overflow: hidden; margin-bottom: 8px; }
.index-bar-fill { height: 100%; border-radius: 6px; transition: width 1.5s ease; background: linear-gradient(90deg, var(--red) 0%, var(--gold) 50%, var(--green) 100%); }
.index-labels { display: flex; justify-content: space-between; font-size: 0.6rem; color: var(--muted); }
.index-phase { margin-top: 12px; font-size: 0.72rem; color: var(--muted); text-align: center; padding: 10px; background: #161b22; border-radius: 8px; }
</style>
</head>
<body>

<header>
  <div class="logo">
    <span style="font-size:1.5rem">⚖</span>
    <div style="margin-left:10px">
      <div class="logo-text">AEQUITAS</div>
      <div class="logo-sub">CHAIN EXPLORER</div>
    </div>
  </div>
  <div class="header-right">
    <div class="live-badge"><span class="pulse"></span>LIVE</div>
    <div class="chain-badge">● BLOCKDAG</div>
  </div>
</header>

<div class="tabs">
  <div class="tab active" onclick="showTab('explorer')">🔍 Explorer</div>
  <div class="tab" onclick="showTab('humans')">👥 Humans</div>
  <div class="tab" onclick="showTab('index')">📊 Aequitas Index</div>
  <div class="tab" onclick="showTab('network')">🌐 Network</div>
  <div class="tab" onclick="showTab('register')">🔐 Register</div>
</div>

<!-- TAB: EXPLORER -->
<div id="tab-explorer" class="tab-content active">
  <div class="hero">
    <div class="hero-title">⚡ Live Chain Statistics</div>
    <div class="stats-grid">
      <div class="stat blue">
        <div class="stat-label">Block Height</div>
        <div class="stat-value" id="s-height">—</div>
        <div class="stat-sub">New block every 6 seconds<br>BlockDAG consensus</div>
      </div>
      <div class="stat green">
        <div class="stat-label">Verified Humans</div>
        <div class="stat-value" id="s-humans">—</div>
        <div class="stat-sub">Biometric Proof of Humanity<br>One person, one wallet</div>
      </div>
      <div class="stat gold">
        <div class="stat-label">Total Supply</div>
        <div class="stat-value" id="s-supply">—</div>
        <div class="stat-sub">= Humans × 1,000 AEQ<br>Supply follows humanity</div>
      </div>
      <div class="stat purple">
        <div class="stat-label">Aequitas Index</div>
        <div class="stat-value" id="s-index">—</div>
        <div class="stat-sub">0 = perfect equality<br>100 = maximum inequality</div>
      </div>
      <div class="stat green">
        <div class="stat-label">Node Uptime</div>
        <div class="stat-value" id="s-uptime" style="font-size:1.1rem">—</div>
        <div class="stat-sub" id="s-version">Aequitas Chain v0.3</div>
      </div>
    </div>

    <div class="explainer">
      <div class="explainer-title">💡 How Aequitas Works</div>
      <div class="explainer-grid">
        <div class="explainer-item">
          <div class="explainer-item-title">🧬 Proof of Humanity</div>
          <div class="explainer-item-text">Every AEQ holder must prove they are a unique human via biometric verification. No bots, no duplicates, no fake accounts.</div>
        </div>
        <div class="explainer-item">
          <div class="explainer-item-title">⚖ Fair Distribution</div>
          <div class="explainer-item-text">Every verified human receives exactly 1,000 AEQ. No pre-mine, no investor allocation. Total supply = verified humans × 1,000.</div>
        </div>
        <div class="explainer-item">
          <div class="explainer-item-title">🔗 BlockDAG Chain</div>
          <div class="explainer-item-text">Aequitas uses a Directed Acyclic Graph (DAG) for blocks, allowing parallel block production and higher throughput than traditional blockchains.</div>
        </div>
        <div class="explainer-item">
          <div class="explainer-item-title">⛽ Gasless</div>
          <div class="explainer-item-text">Registration is completely free. No gas fees, no ETH needed. If you're human, you can register. Period.</div>
        </div>
      </div>
    </div>
  </div>

  <div class="main">
    <div>
      <div class="section">
        <div class="section-header">
          <div class="section-title">Recent Blocks</div>
          <div class="section-count" id="block-count">Loading...</div>
        </div>
        <div class="section-desc">Each block contains transactions and is linked to previous blocks via cryptographic hashes. 🔀 = merged multiple tips (BlockDAG feature). ✅ = contains registration transactions.</div>
        <div id="blocks-list"><div class="empty-state">Loading blocks...</div></div>
      </div>
    </div>

    <div class="right-panel">
      <div class="info-panel">
        <div class="info-title">🌐 Network Info</div>
        <div class="info-row"><span class="info-key">Chain ID</span><span class="info-val blue">9001</span></div>
        <div class="info-row"><span class="info-key">Symbol</span><span class="info-val gold">AEQ</span></div>
        <div class="info-row"><span class="info-key">Block Time</span><span class="info-val">6 seconds</span></div>
        <div class="info-row"><span class="info-key">Consensus</span><span class="info-val purple">BlockDAG PoH</span></div>
        <div class="info-row"><span class="info-key">Nodes</span><span class="info-val green">2 Active</span></div>
        <div class="info-row"><span class="info-key">RPC</span><span class="info-val blue">/rpc</span></div>
      </div>

      <div class="mm-config">
        <div class="mm-title">🦊 ADD TO METAMASK</div>
        <div class="mm-row"><span class="mm-key">Network Name</span><span class="mm-val">Aequitas Chain</span></div>
        <div class="mm-row"><span class="mm-key">RPC URL</span><span class="mm-val" style="font-size:0.55rem">aequitas-production-9fba.up.railway.app/rpc</span></div>
        <div class="mm-row"><span class="mm-key">Chain ID</span><span class="mm-val">9001</span></div>
        <div class="mm-row"><span class="mm-key">Symbol</span><span class="mm-val">AEQ</span></div>
        <button onclick="addToMetaMask()" style="width:100%;margin-top:10px;padding:8px;background:var(--blue);color:white;border:none;border-radius:6px;cursor:pointer;font-family:monospace;font-size:0.72rem;letter-spacing:1px">+ ADD NETWORK</button>
      </div>

      <div class="philosophy">
        <div class="philosophy-quote">"Money exists because people exist.<br>Nothing more, nothing less."</div>
        <div class="philosophy-sub">The Aequitas Principle</div>
      </div>
    </div>
  </div>
</div>

<!-- TAB: HUMANS -->
<div id="tab-humans" class="tab-content">
  <div class="hero">
    <div class="hero-title">👥 Verified Humans on Aequitas Chain</div>
    <div class="explainer" style="margin-bottom:16px">
      <div class="explainer-title">🔒 Proof of Humanity</div>
      <div class="explainer-grid">
        <div class="explainer-item">
          <div class="explainer-item-title">What is it?</div>
          <div class="explainer-item-text">Each address listed here has been verified as a unique human using biometric data (fingerprint). The actual biometric data never leaves the device — only a cryptographic proof is used.</div>
        </div>
        <div class="explainer-item">
          <div class="explainer-item-title">How does it work?</div>
          <div class="explainer-item-text">A Groth16 Zero-Knowledge Proof (ZKP) is generated from your biometric hash. This proves you are human without revealing any personal data. One human = one wallet = 1,000 AEQ.</div>
        </div>
      </div>
    </div>
    <div class="section" style="margin:0 32px 32px">
      <div class="section-header">
        <div class="section-title">Registered Humans</div>
        <div class="section-count" id="human-count-tab">0 humans</div>
      </div>
      <div class="section-desc">All verified humans on the Aequitas Chain. Each human received 1,000 AEQ upon registration. Registration is permanent and cannot be transferred.</div>
      <div id="humans-list"><div class="empty-state">No humans registered yet.<br>Download the Aequitas App to register.</div></div>
    </div>
  </div>
</div>

<!-- TAB: INDEX -->
<div id="tab-index" class="tab-content">
  <div class="index-wrap">
    <div class="index-card">
      <div class="index-title">📊 Aequitas Index</div>
      <div class="index-desc">The Aequitas Index measures economic equality on the chain. It combines the Gini coefficient (wealth distribution) with network growth metrics to produce a single score from 0 (perfect equality) to 100 (maximum inequality). The goal is to keep this index as low as possible through the automatic wealth redistribution mechanisms built into the protocol.</div>
      <div class="index-grid">
        <div class="index-metric">
          <div class="index-metric-label">Index Score</div>
          <div class="index-metric-val" id="idx-score">—</div>
          <div class="index-metric-sub">0=equal, 100=unequal</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Gini Coefficient</div>
          <div class="index-metric-val" id="idx-gini">—</div>
          <div class="index-metric-sub">wealth distribution</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Total Supply</div>
          <div class="index-metric-val" id="idx-supply">—</div>
          <div class="index-metric-sub">AEQ in circulation</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Phase</div>
          <div class="index-metric-val" id="idx-phase">—</div>
          <div class="index-metric-sub">protocol phase</div>
        </div>
      </div>
      <div class="index-bar-bg">
        <div class="index-bar-fill" id="idx-bar" style="width:0%"></div>
      </div>
      <div class="index-labels">
        <span>0 — Perfect Equality</span>
        <span>50 — Moderate</span>
        <span>100 — Max Inequality</span>
      </div>
      <div class="index-phase" id="idx-phase-desc">Loading...</div>
    </div>

    <div class="index-card">
      <div class="index-title">🏦 Automatic Redistribution</div>
      <div class="index-desc">When economic inequality exceeds certain thresholds, the protocol automatically activates redistribution mechanisms. These include wealth caps, inflation adjustments, and pool transfers — all governed by smart contract logic, not human decisions.</div>
      <div class="index-grid">
        <div class="index-metric">
          <div class="index-metric-label">Velocity Pool</div>
          <div class="index-metric-val" id="pool-v">—</div>
          <div class="index-metric-sub">transaction incentives</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Liquidity Pool</div>
          <div class="index-metric-val" id="pool-l">—</div>
          <div class="index-metric-sub">market stability</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Unity Pool</div>
          <div class="index-metric-val" id="pool-u">—</div>
          <div class="index-metric-sub">new registrations</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Treasury</div>
          <div class="index-metric-val" id="pool-t">—</div>
          <div class="index-metric-sub">protocol development</div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- TAB: NETWORK -->
<div id="tab-network" class="tab-content">
  <div class="index-wrap">
    <div class="index-card">
      <div class="index-title">🌐 Network Topology</div>
      <div class="index-desc">The Aequitas network currently runs on 2 nodes — one on Railway (primary) and one on Render (secondary). Both nodes participate in block production and sync with each other via HTTP block sync and libp2p P2P networking.</div>
      <div class="index-grid" style="grid-template-columns: repeat(2, 1fr)">
        <div class="index-metric">
          <div class="index-metric-label">Node 1 — Railway</div>
          <div class="index-metric-val" style="font-size:0.9rem;color:var(--green)">● Online</div>
          <div class="index-metric-sub">aequitas-production-9fba.up.railway.app</div>
        </div>
        <div class="index-metric">
          <div class="index-metric-label">Node 2 — Render</div>
          <div class="index-metric-val" style="font-size:0.9rem;color:var(--green)">● Online</div>
          <div class="index-metric-sub">aequitas-node-2.onrender.com</div>
        </div>
      </div>
    </div>

    <div class="index-card">
      <div class="index-title">🔗 Bootstrap Node</div>
      <div class="index-desc">To join the Aequitas network, connect to the bootstrap node using libp2p. New nodes can sync the full block history and participate in consensus.</div>
      <div style="background:#161b22;border-radius:8px;padding:14px;font-size:0.68rem;color:var(--purple);word-break:break-all;line-height:1.8">
        /dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R
      </div>
    </div>

    <div class="index-card">
      <div class="index-title">⚙ Technical Specifications</div>
      <div class="index-desc">Aequitas Chain technical details for developers and node operators.</div>
      <div style="display:grid;gap:8px">
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">Chain ID</span><span class="info-val blue">9001</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">EVM Compatible</span><span class="info-val green">Yes (JSON-RPC)</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">Block Time</span><span class="info-val">6 seconds</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">Consensus</span><span class="info-val purple">BlockDAG + Proof of Humanity</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">P2P Protocol</span><span class="info-val">libp2p</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">ZKP System</span><span class="info-val">Groth16 (snarkjs)</span></div>
        <div class="info-row" style="padding:8px 0;border-bottom:1px solid #161b22"><span class="info-key">State Storage</span><span class="info-val green">PostgreSQL (persistent)</span></div>
        <div class="info-row" style="padding:8px 0"><span class="info-key">Source Code</span><span class="info-val blue"><a href="https://github.com/hanoi96international-gif/Aequitas" target="_blank" style="color:var(--blue)">GitHub ↗</a></span></div>
      </div>
    </div>
  </div>
</div>

<!-- TAB: REGISTER -->
<div id="tab-register" class="tab-content">
  <div class="reg-wrap">
    <div class="reg-card">
      <div class="reg-title">🔐 Register as Human</div>
      <div class="reg-sub">Join the Aequitas network and receive your 1,000 AEQ. Registration requires biometric verification via the Android app to ensure you are a unique human. No gas fees required.</div>

      <div class="app-only-box">
        <div class="app-only-icon">📱</div>
        <div class="app-only-title">REGISTRATION VIA ANDROID APP ONLY</div>
        <div class="app-only-text">Proof of Humanity requires biometric verification on your device.<br>Download the Aequitas App, scan your fingerprint,<br>and your 1,000 AEQ will be granted automatically.<br><br>Your biometric data <strong style="color:var(--gold)">never leaves your device</strong>.<br>Only a cryptographic proof is transmitted.</div>
      </div>

      <div class="reg-steps">
        <div class="reg-step">
          <div class="reg-step-num">1</div>
          <div class="reg-step-title">Biometric Scan</div>
          <div class="reg-step-desc">Fingerprint scanned via Hardware Secure Element on your device</div>
        </div>
        <div class="reg-step">
          <div class="reg-step-num">2</div>
          <div class="reg-step-title">ZKP Generation</div>
          <div class="reg-step-desc">Groth16 Zero-Knowledge Proof generated — proves humanity without revealing data</div>
        </div>
        <div class="reg-step">
          <div class="reg-step-num">3</div>
          <div class="reg-step-title">Connect Wallet</div>
          <div class="reg-step-desc">Connect MetaMask or any Web3 wallet to receive your AEQ</div>
        </div>
        <div class="reg-step">
          <div class="reg-step-num">4</div>
          <div class="reg-step-title">1,000 AEQ</div>
          <div class="reg-step-desc">Instantly credited to your wallet. No gas fees. No waiting.</div>
        </div>
      </div>

      <div class="privacy-badge">🔒 Hardware Secure Element · Real Groth16 ZKP · No gas fees · Permanent Sybil protection</div>

      <div class="reg-wallet-box" id="wallet-box">
        <div class="reg-wallet-label">CONNECTED WALLET</div>
        <div class="reg-wallet-addr" id="wallet-addr">—</div>
      </div>

      <div class="reg-proof-box" id="proof-box">
        <div class="reg-proof-label">⚡ Proof parameters detected from app</div>
        <div class="reg-proof-val" id="proof-val">—</div>
      </div>

      <button class="reg-btn reg-btn-connect" id="btn-connect" onclick="connectWallet()">🦊 CONNECT METAMASK</button>
      <button class="reg-btn reg-btn-register" id="btn-register" onclick="register()" disabled>🔐 REGISTER ON-CHAIN</button>

      <div class="reg-status" id="reg-status">
        <span class="info">// Open Aequitas Android App to generate your proof...</span>
      </div>
    </div>

    <div class="info-panel">
      <div class="info-title">ℹ Registration Details</div>
      <div class="info-row"><span class="info-key">Network</span><span class="info-val purple">Aequitas Chain (BlockDAG)</span></div>
      <div class="info-row"><span class="info-key">RPC URL</span><span class="info-val blue" style="font-size:0.6rem">aequitas-production-9fba.up.railway.app/rpc</span></div>
      <div class="info-row"><span class="info-key">Chain ID</span><span class="info-val gold">9001</span></div>
      <div class="info-row"><span class="info-key">Grant Amount</span><span class="info-val gold">1,000 AEQ</span></div>
      <div class="info-row"><span class="info-key">Gas Fee</span><span class="info-val green">FREE (gasless)</span></div>
      <div class="info-row"><span class="info-key">Registrations</span><span class="info-val">Once per human, permanent</span></div>
    </div>
  </div>
</div>

<script>
const PROOF_SERVER = 'https://aequitas-proof-server-production.up.railway.app';
let walletAddr = '';
let proofParams = null;

function showTab(name) {
  document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-' + name).classList.add('active');
  event.target.classList.add('active');
}

function fmt(n) {
  if (n === undefined || n === null || n === '—') return '—';
  if (typeof n === 'number') return n.toLocaleString();
  return n;
}

function timeAgo(ts) {
  const diff = Math.floor(Date.now()/1000) - ts;
  if (diff < 60) return diff + 's ago';
  if (diff < 3600) return Math.floor(diff/60) + 'm ago';
  return Math.floor(diff/3600) + 'h ago';
}

function shortHash(h) {
  return h ? h.slice(0,8) + '...' + h.slice(-6) : '—';
}

function shortAddr(a) {
  return a ? a.slice(0,8) + '...' + a.slice(-4) : '—';
}

function avatarColor(addr) {
  const colors = ['#3B82F6','#22C55E','#C9A84C','#BC8CFF','#EF4444','#F97316'];
  const idx = parseInt(addr.slice(2,4), 16) % colors.length;
  return colors[idx];
}

async function addToMetaMask() {
  if (!window.ethereum) { alert('MetaMask not found'); return; }
  try {
    await window.ethereum.request({
      method: 'wallet_addEthereumChain',
      params: [{
        chainId: '0x2329',
        chainName: 'Aequitas Chain',
        nativeCurrency: { name: 'AEQ', symbol: 'AEQ', decimals: 18 },
        rpcUrls: ['https://aequitas-production-9fba.up.railway.app/rpc'],
        blockExplorerUrls: ['https://aequitas-production-9fba.up.railway.app']
      }]
    });
  } catch(e) { console.error(e); }
}

async function loadStatus() {
  try {
    const r = await fetch('/api/status');
    const d = await r.json();
    document.getElementById('s-height').textContent = fmt(d.block_height);
    document.getElementById('s-humans').textContent = fmt(d.total_humans);
    document.getElementById('s-supply').textContent = d.total_supply || '—';
    document.getElementById('s-index').textContent = fmt(d.index);
    document.getElementById('s-uptime').textContent = d.uptime || '—';
    document.getElementById('idx-score').textContent = fmt(d.index);
    document.getElementById('idx-gini').textContent = fmt(d.gini);
    document.getElementById('idx-supply').textContent = d.total_supply || '—';
    document.getElementById('idx-phase').textContent = fmt(d.phase);
    document.getElementById('human-count-tab').textContent = fmt(d.total_humans) + ' humans';
    if (d.index !== undefined) {
      document.getElementById('idx-bar').style.width = Math.min(d.index, 100) + '%';
      const phase = d.phase || 0;
      const phases = ['Phase 0: Bootstrap — Building the network','Phase 1: Growth — Expanding human registry','Phase 2: Stability — Wealth redistribution active','Phase 3: Maturity — Full decentralization'];
      document.getElementById('idx-phase-desc').textContent = phases[phase] || 'Phase ' + phase;
    }
  } catch(e) {}
}

async function loadBlocks() {
  try {
    const r = await fetch('/api/blocks?limit=20');
    const blocks = await r.json();
    const list = document.getElementById('blocks-list');
    if (!blocks || !blocks.length) { list.innerHTML = '<div class="empty-state">No blocks yet</div>'; return; }
    document.getElementById('block-count').textContent = blocks.length + ' blocks';
    list.innerHTML = blocks.map(b => {
      const isMerge = b.parent_hashes && b.parent_hashes.length > 1;
      const hasTx = b.transactions && b.transactions.length > 0;
      return ` + "`" + `<div class="block-item">
        <div class="block-height">#${b.height}</div>
        <div class="block-info">
          <div class="block-hash">${shortHash(b.hash)}${isMerge ? '<span class="merge-badge">🔀 MERGE</span>' : ''}${hasTx ? '<span class="tx-badge">✅ TX</span>' : ''}</div>
          <div class="block-meta">${b.parent_hashes ? b.parent_hashes.length + ' parent(s)' : ''}</div>
        </div>
        <div class="block-right">
          <div class="block-humans">${b.humans || 0} humans</div>
          <div class="block-time">${timeAgo(b.timestamp)}</div>
        </div>
      </div>` + "`" + `;
    }).join('');
  } catch(e) {}
}

async function loadHumans() {
  try {
    const r = await fetch('/api/humans');
    const d = await r.json();
    const list = document.getElementById('humans-list');
    if (!d.humans || !d.humans.length) {
      list.innerHTML = '<div class="empty-state">No humans registered yet.<br><br>Be the first! Download the Aequitas App<br>and scan your fingerprint to register.</div>';
      return;
    }
    list.innerHTML = d.humans.map(h => {
      const color = avatarColor(h.address || '0x00');
      const initials = (h.address || '?').slice(2,4).toUpperCase();
      return ` + "`" + `<div class="human-item">
        <div class="human-avatar" style="background:${color}">${initials}</div>
        <div class="human-info">
          <div class="human-balance">${fmt(h.balance)} AEQ</div>
          <div class="human-addr">${h.address}</div>
        </div>
        <div class="human-badge">✓ HUMAN</div>
      </div>` + "`" + `;
    }).join('');
  } catch(e) {}
}

// Check for proof params from app
function checkProofParams() {
  const params = new URLSearchParams(window.location.search);
  const proof = params.get('proof');
  if (proof) {
    try {
      proofParams = JSON.parse(decodeURIComponent(proof));
      const box = document.getElementById('proof-box');
      box.style.display = 'block';
      document.getElementById('proof-val').textContent =
        'bio: ' + proofParams.bio.slice(0,15) + '... | salt: ' + proofParams.salt.slice(0,10) + '...';
      log('⚡ Proof parameters detected from app', 'info');
      showTab('register');
      document.querySelectorAll('.tab')[4].classList.add('active');
    } catch(e) {}
  }
}

async function connectWallet() {
  if (!window.ethereum) {
    log('✗ MetaMask not found', 'err');
    return;
  }
  try {
    await addToMetaMask();
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    walletAddr = accounts[0];
    const box = document.getElementById('wallet-box');
    box.style.display = 'block';
    document.getElementById('wallet-addr').textContent = walletAddr;
    document.getElementById('btn-register').disabled = !proofParams;
    document.getElementById('btn-connect').textContent = '✓ ' + walletAddr.slice(0,10) + '...' + walletAddr.slice(-4);
    document.getElementById('btn-connect').style.background = 'var(--green)';
    document.getElementById('btn-connect').style.color = '#0A0E1A';
    log('✓ Wallet connected: ' + walletAddr.slice(0,12) + '...', 'ok');
  } catch(e) {
    log('✗ ' + e.message, 'err');
  }
}

function log(msg, type) {
  const el = document.getElementById('reg-status');
  el.innerHTML += ` + "`" + `<div><span class="${type}">${msg}</span></div>` + "`" + `;
}

async function register() {
  if (!walletAddr) { log('✗ Connect wallet first', 'err'); return; }
  if (!proofParams) { log('✗ No proof parameters. Use the Android app first.', 'err'); return; }
  try {
    log('⏳ Step 1/2: Generating ZK proof...', 'info');
    document.getElementById('btn-register').disabled = true;
    const proveResp = await fetch(PROOF_SERVER + '/prove', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ bio: proofParams.bio, salt: proofParams.salt, wallet: walletAddr }) });
    const proveData = await proveResp.json();
    if (!proveResp.ok) {
      log('✗ ' + (proveData.error || 'Proof failed'), 'err');
      document.getElementById('btn-register').disabled = false;
      return;
    }
    log('✓ ZK Proof generated! Step 2/2: Registering on chain...', 'ok');
    const resp = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bio: proofParams.bio, salt: proofParams.salt, wallet: walletAddr })
    });
    const data = await resp.json();
    if (!data.success) {
      log('✗ ' + data.message, 'err');
      document.getElementById('btn-register').disabled = false;
      return;
    }
    log('🎉 ' + data.message + ' | TX: ' + data.tx_hash, 'ok');
    // Redirect to registered page
    setTimeout(() => {
      window.location.href = '/registered?success=true&wallet=' + walletAddr;
    }, 1500);
  } catch(e) {
    log('✗ ' + e.message, 'err');
    document.getElementById('btn-register').disabled = false;
  }
}

// Init
checkProofParams();
loadStatus();
loadBlocks();
loadHumans();
setInterval(loadStatus, 6000);
setInterval(loadBlocks, 6000);
setInterval(loadHumans, 10000);

window.ethereum?.on('accountsChanged', (accounts) => {
  walletAddr = accounts[0] || '';
  if (walletAddr) {
    document.getElementById('wallet-box').style.display = 'block';
    document.getElementById('wallet-addr').textContent = walletAddr;
    document.getElementById('btn-register').disabled = !proofParams;
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
