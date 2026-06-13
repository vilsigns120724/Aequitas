package keeper

import (
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

type APIServer struct {
blockchain    *Blockchain
p2pNode       *P2PNode
keeper        *Keeper
startTime     time.Time
sepoliaStatus map[string]interface{}
}

func NewAPIServer(bc *Blockchain, p2p *P2PNode, k *Keeper) *APIServer {
s := &APIServer{
blockchain:    bc,
p2pNode:       p2p,
keeper:        k,
startTime:     time.Now(),
sepoliaStatus: map[string]interface{}{},
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

// Aequitas Index (simplified on-chain calculation)
gini := 0 // grows as humans accumulate different balances
velocity := 50
growth := humans * 10
if growth > 100 {
growth = 100
}
index := (velocity*40 + growth*35 + (100-gini)*25) / 100

json.NewEncoder(w).Encode(map[string]interface{}{
"chain_id":      "aequitas-1",
"version":       "v0.3.0",
"height":        latest.Height,
"latest_hash":   latest.Hash,
"total_humans":  humans,
"total_supply":  humans * 1000,
"node_id":       a.p2pNode.GetNodeID(),
"uptime":        uptime,
"block_time":    6,
"contract_v5":   "0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5",
// Aequitas Index
"index":    index,
"velocity": velocity,
"growth":   growth,
"gini":     gini,
// Fee pools (accumulate from 0.1% tx fees on Sepolia)
"fee_bps":         10,
"validator_pool":  a.sepoliaStatus["validatorPool"],
"lp_pool":         a.sepoliaStatus["lpPool"],
"ubi_pool":        a.sepoliaStatus["ubiPool"],
"treasury":        a.sepoliaStatus["treasury"],
// Sepolia live
"sepolia": a.sepoliaStatus,
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
json.NewEncoder(w).Encode(map[string]interface{}{
"total":  a.keeper.TotalHumans(),
"humans": a.keeper.GetAllHumans(),
})
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
*{margin:0;padding:0;box-sizing:border-box}
body{background:#060810;color:#c9d1d9;font-family:'Courier New',monospace;min-height:100vh}
:root{--green:#3fb950;--blue:#58a6ff;--gold:#e3b341;--red:#f85149;--border:#21262d;--card:#0d1117;--bg:#060810;--purple:#bc8cff}
header{background:#0d1117;border-bottom:1px solid var(--border);padding:18px 32px;display:flex;align-items:center;justify-content:space-between;position:sticky;top:0;z-index:100}
.logo{display:flex;align-items:center;gap:12px}
.logo-icon{font-size:1.5rem}
.logo-text{font-size:1.2rem;font-weight:bold;color:#f0f6fc;letter-spacing:2px}
.logo-sub{font-size:0.7rem;color:#8b949e;letter-spacing:1px}
.header-right{display:flex;align-items:center;gap:12px}
.live-badge{display:flex;align-items:center;gap:6px;background:#1a3a2a;border:1px solid #2ea04326;padding:5px 12px;border-radius:20px;font-size:0.72rem;color:var(--green)}
.pulse{width:7px;height:7px;background:var(--green);border-radius:50%;animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:0.3}}
.chain-badge{font-size:0.72rem;color:#8b949e;background:#161b22;border:1px solid var(--border);padding:5px 12px;border-radius:20px}

/* STATS */
.stats{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:12px;padding:24px 32px 0}
.stat{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px;position:relative;overflow:hidden}
.stat::before{content:'';position:absolute;top:0;left:0;right:0;height:2px}
.stat.green::before{background:var(--green)}
.stat.blue::before{background:var(--blue)}
.stat.gold::before{background:var(--gold)}
.stat.purple::before{background:var(--purple)}
.stat-label{font-size:0.65rem;color:#8b949e;text-transform:uppercase;letter-spacing:1.5px;margin-bottom:8px}
.stat-value{font-size:1.7rem;font-weight:bold;line-height:1}
.stat.green .stat-value{color:var(--green)}
.stat.blue .stat-value{color:var(--blue)}
.stat.gold .stat-value{color:var(--gold)}
.stat.purple .stat-value{color:var(--purple)}
.stat-sub{font-size:0.68rem;color:#8b949e;margin-top:5px}

/* INDEX SECTION */
.index-bar-wrap{padding:24px 32px 0}
.index-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.index-title{font-size:0.7rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px;margin-bottom:16px}
.index-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:12px;margin-bottom:16px}
.index-metric{text-align:center}
.index-metric-label{font-size:0.65rem;color:#8b949e;margin-bottom:4px}
.index-metric-val{font-size:1.3rem;font-weight:bold;color:var(--gold)}
.index-bar-bg{height:8px;background:#161b22;border-radius:4px;overflow:hidden}
.index-bar-fill{height:100%;border-radius:4px;transition:width 1s ease;background:linear-gradient(90deg,var(--red),var(--gold),var(--green))}
.index-labels{display:flex;justify-content:space-between;font-size:0.6rem;color:#8b949e;margin-top:4px}

/* FEE POOLS */
.pools-wrap{padding:16px 32px 0}
.pools-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.pools-title{font-size:0.7rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px;margin-bottom:16px}
.pools-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:12px}
@media(max-width:700px){.pools-grid{grid-template-columns:repeat(2,1fr)}}
.pool-item{background:#161b22;border-radius:8px;padding:14px;text-align:center}
.pool-icon{font-size:1.2rem;margin-bottom:6px}
.pool-label{font-size:0.62rem;color:#8b949e;margin-bottom:4px}
.pool-val{font-size:1rem;font-weight:bold;color:var(--gold)}
.pool-pct{font-size:0.6rem;color:#8b949e;margin-top:2px}

/* SEPOLIA STATUS */
.sepolia-wrap{padding:16px 32px 0}
.sepolia-card{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px}
.sepolia-title{font-size:0.7rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px;margin-bottom:14px;display:flex;align-items:center;gap:8px}
.sepolia-dot{width:6px;height:6px;background:var(--green);border-radius:50%;animation:pulse 2s infinite}
.sepolia-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(120px,1fr));gap:10px}
.sepolia-item{background:#161b22;border-radius:6px;padding:10px}
.sepolia-item-label{font-size:0.6rem;color:#8b949e;margin-bottom:3px}
.sepolia-item-val{font-size:0.9rem;font-weight:bold;color:var(--blue)}

/* MAIN */
.main{display:grid;grid-template-columns:1fr 360px;gap:16px;padding:16px 32px 32px}
@media(max-width:900px){.main{grid-template-columns:1fr}}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.section-header{padding:14px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between}
.section-title{font-size:0.72rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px}
.section-count{font-size:0.68rem;color:#8b949e;background:#161b22;padding:2px 8px;border-radius:10px}
.block-item{padding:11px 20px;border-bottom:1px solid #161b22;display:grid;grid-template-columns:80px 1fr 60px 85px;gap:10px;align-items:center}
.block-item:hover{background:#161b2270}
.block-item:last-child{border-bottom:none}
.block-height{color:var(--blue);font-weight:bold;font-size:0.82rem}
.block-hash{color:#8b949e;font-size:0.72rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.block-humans{color:var(--gold);font-size:0.72rem;text-align:right}
.block-time{color:var(--green);font-size:0.72rem;text-align:right}
.human-item{padding:11px 18px;border-bottom:1px solid #161b22;display:flex;align-items:center;gap:10px}
.human-item:last-child{border-bottom:none}
.human-avatar{width:30px;height:30px;background:linear-gradient(135deg,#1f6feb,#388bfd);border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.7rem;color:white;font-weight:bold;flex-shrink:0}
.human-addr{font-size:0.72rem;color:#8b949e;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1}
.human-balance{font-size:0.68rem;color:var(--gold);flex-shrink:0}
.human-badge{font-size:0.62rem;color:var(--green);background:#1a3a2a;padding:2px 7px;border-radius:8px;flex-shrink:0}
.info-panel{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:18px;margin-bottom:14px}
.info-title{font-size:0.68rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px;margin-bottom:12px}
.info-row{display:flex;justify-content:space-between;padding:7px 0;border-bottom:1px solid #161b22}
.info-row:last-child{border-bottom:none}
.info-key{font-size:0.7rem;color:#8b949e}
.info-val{font-size:0.7rem;color:#f0f6fc;text-align:right;max-width:55%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.info-val.green{color:var(--green)}
.info-val.blue{color:var(--blue)}
.info-val.gold{color:var(--gold)}
.link-row{display:flex;gap:8px;flex-wrap:wrap;margin-top:12px}
.ext-link{font-size:0.68rem;color:var(--blue);background:#1c2d40;border:1px solid #1f6feb30;padding:5px 11px;border-radius:6px;text-decoration:none}
.ext-link:hover{background:#1f6feb20}
</style>
</head>
<body>

<header>
  <div class="logo">
    <span class="logo-icon">⚖</span>
    <div>
      <div class="logo-text">AEQUITAS</div>
      <div class="logo-sub">CHAIN EXPLORER</div>
    </div>
  </div>
  <div class="header-right">
    <div class="live-badge"><span class="pulse"></span>LIVE</div>
    <div class="chain-badge">aequitas-1</div>
  </div>
</header>

<!-- STATS -->
<div class="stats">
  <div class="stat blue"><div class="stat-label">Block Height</div><div class="stat-value" id="s-height">—</div><div class="stat-sub">every 6 seconds</div></div>
  <div class="stat green"><div class="stat-label">Verified Humans</div><div class="stat-value" id="s-humans">—</div><div class="stat-sub">Proof of Humanity</div></div>
  <div class="stat gold"><div class="stat-label">Total Supply</div><div class="stat-value" id="s-supply">—</div><div class="stat-sub">AEQ · dynamic cap</div></div>
  <div class="stat purple"><div class="stat-label">Aequitas Index</div><div class="stat-value" id="s-index">—</div><div class="stat-sub">economic health</div></div>
  <div class="stat green"><div class="stat-label">Uptime</div><div class="stat-value" id="s-uptime" style="font-size:1.1rem">—</div><div class="stat-sub" id="s-version">—</div></div>
</div>

<!-- AEQUITAS INDEX -->
<div class="index-bar-wrap">
  <div class="index-card">
    <div class="index-title">⚖ Aequitas Index</div>
    <div class="index-grid">
      <div class="index-metric"><div class="index-metric-label">Velocity</div><div class="index-metric-val" id="idx-velocity">—</div></div>
      <div class="index-metric"><div class="index-metric-label">Growth</div><div class="index-metric-val" id="idx-growth">—</div></div>
      <div class="index-metric"><div class="index-metric-label">Gini</div><div class="index-metric-val" id="idx-gini">—</div></div>
      <div class="index-metric"><div class="index-metric-label">Index Score</div><div class="index-metric-val" id="idx-score">—</div></div>
    </div>
    <div class="index-bar-bg"><div class="index-bar-fill" id="idx-bar" style="width:0%"></div></div>
    <div class="index-labels"><span>Recession</span><span>Neutral</span><span>Boom</span></div>
  </div>
</div>

<!-- FEE POOLS -->
<div class="pools-wrap">
  <div class="pools-card">
    <div class="pools-title">💰 Fee Pools (0.1% per transaction)</div>
    <div class="pools-grid">
      <div class="pool-item"><div class="pool-icon">⛏</div><div class="pool-label">Validators</div><div class="pool-val" id="pool-v">—</div><div class="pool-pct">40% of fees</div></div>
      <div class="pool-item"><div class="pool-icon">💧</div><div class="pool-label">Liquidity</div><div class="pool-val" id="pool-l">—</div><div class="pool-pct">30% of fees</div></div>
      <div class="pool-item"><div class="pool-icon">🌍</div><div class="pool-label">UBI</div><div class="pool-val" id="pool-u">—</div><div class="pool-pct">20% of fees</div></div>
      <div class="pool-item"><div class="pool-icon">🏛</div><div class="pool-label">Treasury</div><div class="pool-val" id="pool-t">—</div><div class="pool-pct">10% of fees</div></div>
    </div>
  </div>
</div>

<!-- SEPOLIA STATUS -->
<div class="sepolia-wrap">
  <div class="sepolia-card">
    <div class="sepolia-title"><span class="sepolia-dot"></span>Sepolia Contract V5 — Live Status</div>
    <div class="sepolia-grid">
      <div class="sepolia-item"><div class="sepolia-item-label">Registrations</div><div class="sepolia-item-val" id="sep-humans">—</div></div>
      <div class="sepolia-item"><div class="sepolia-item-label">Supply</div><div class="sepolia-item-val" id="sep-supply">—</div></div>
      <div class="sepolia-item"><div class="sepolia-item-label">Gini</div><div class="sepolia-item-val" id="sep-gini">—</div></div>
      <div class="sepolia-item"><div class="sepolia-item-label">Index</div><div class="sepolia-item-val" id="sep-index">—</div></div>
      <div class="sepolia-item"><div class="sepolia-item-label">Phase</div><div class="sepolia-item-val" id="sep-phase">—</div></div>
      <div class="sepolia-item"><div class="sepolia-item-label">Status</div><div class="sepolia-item-val" style="color:#3fb950" id="sep-status">—</div></div>
    </div>
  </div>
</div>

<!-- MAIN: BLOCKS + SIDEBAR -->
<div class="main">
  <div class="section">
    <div class="section-header">
      <span class="section-title">Recent Blocks</span>
      <span class="section-count" id="block-count">—</span>
    </div>
    <div id="blocks-list"></div>
  </div>

  <div>
    <div class="info-panel">
      <div class="info-title">Node Info</div>
      <div class="info-row"><span class="info-key">Node ID</span><span class="info-val blue" id="i-nodeid">—</span></div>
      <div class="info-row"><span class="info-key">Latest Hash</span><span class="info-val" id="i-hash">—</span></div>
      <div class="info-row"><span class="info-key">Block Time</span><span class="info-val green">6 seconds</span></div>
      <div class="info-row"><span class="info-key">Fee</span><span class="info-val gold">0.1% per tx</span></div>
      <div class="info-row"><span class="info-key">Initial Grant</span><span class="info-val gold">1,000 AEQ</span></div>
      <div class="info-row"><span class="info-key">Sepolia Sync</span><span class="info-val green">✓ Active</span></div>
      <div class="info-row"><span class="info-key">Contract V5</span><span class="info-val blue">0x4f147d...f0B8b5</span></div>
      <div class="link-row">
        <a class="ext-link" href="https://sepolia.etherscan.io/address/0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5" target="_blank">Etherscan ↗</a>
        <a class="ext-link" href="https://github.com/hanoi96international-gif/Aequitas" target="_blank">GitHub ↗</a>
        <a class="ext-link" href="/api/status" target="_blank">API ↗</a>
      </div>
    </div>

    <div class="section">
      <div class="section-header">
        <span class="section-title">Verified Humans</span>
        <span class="section-count" id="human-count">—</span>
      </div>
      <div id="humans-list"></div>
    </div>
  </div>
</div>

<script>
let uptimeBase=0;
function fmt(n){return Number(n||0).toLocaleString()}
function fmtUptime(s){const h=Math.floor(s/3600),m=Math.floor((s%3600)/60),sec=s%60;return(h?h+'h ':'')+m+'m '+sec+'s'}
function fmtPool(v){if(!v)return'0';const n=Number(v)/1e18;return n<0.001?'<0.001':n.toFixed(3)}

async function update(){
  try{
    const [status,blocks,humans]=await Promise.all([
      fetch('/api/status').then(r=>r.json()),
      fetch('/api/blocks').then(r=>r.json()),
      fetch('/api/humans').then(r=>r.json())
    ]);

    // Stats
    document.getElementById('s-height').textContent='#'+fmt(status.height);
    document.getElementById('s-humans').textContent=fmt(status.total_humans);
    document.getElementById('s-supply').textContent=fmt(status.total_supply);
    document.getElementById('s-index').textContent=status.index||'—';
    document.getElementById('s-version').textContent=status.version||'';
    uptimeBase=status.uptime||0;

    // Node info
    const nid=status.node_id||'';
    document.getElementById('i-nodeid').textContent=nid.slice(0,16)+'...';
    document.getElementById('i-hash').textContent=(status.latest_hash||'').slice(0,12)+'...';

    // Aequitas Index
    document.getElementById('idx-velocity').textContent=status.velocity??'—';
    document.getElementById('idx-growth').textContent=status.growth??'—';
    document.getElementById('idx-gini').textContent=status.gini??'—';
    document.getElementById('idx-score').textContent=status.index??'—';
    document.getElementById('idx-bar').style.width=(status.index||0)+'%';

    // Fee pools
    document.getElementById('pool-v').textContent=fmtPool(status.validator_pool)+' AEQ';
    document.getElementById('pool-l').textContent=fmtPool(status.lp_pool)+' AEQ';
    document.getElementById('pool-u').textContent=fmtPool(status.ubi_pool)+' AEQ';
    document.getElementById('pool-t').textContent=fmtPool(status.treasury)+' AEQ';

    // Sepolia
    const sep=status.sepolia||{};
    document.getElementById('sep-humans').textContent=sep.registrations??'—';
    document.getElementById('sep-supply').textContent=sep.supply??'—';
    document.getElementById('sep-gini').textContent=sep.gini??'—';
    document.getElementById('sep-index').textContent=sep.index??'—';
    document.getElementById('sep-phase').textContent='Phase '+(sep.phase??'—');
    document.getElementById('sep-status').textContent=sep.status==='ok'?'✓ Online':'—';

    // Blocks
    document.getElementById('block-count').textContent=fmt(status.height)+' blocks';
    document.getElementById('blocks-list').innerHTML=[...blocks].reverse().map(b=>{
      const t=new Date(b.timestamp*1000).toLocaleTimeString();
      return '<div class="block-item">'+
        '<span class="block-height">#'+b.height+'</span>'+
        '<span class="block-hash">'+b.hash+'</span>'+
        '<span class="block-humans">👤 '+b.humans+'</span>'+
        '<span class="block-time">'+t+'</span>'+
      '</div>';
    }).join('');

    // Humans
    document.getElementById('human-count').textContent=(humans.total||0)+' registered';
    const list=humans.humans||[];
    document.getElementById('humans-list').innerHTML=list.map(h=>{
      const addr=h.address||'';
      const init=addr.slice(0,2).toUpperCase();
      const short=addr.length>18?addr.slice(0,14)+'...':addr;
      return '<div class="human-item">'+
        '<div class="human-avatar">'+init+'</div>'+
        '<span class="human-addr">'+short+'</span>'+
        '<span class="human-balance">1,000 AEQ</span>'+
        '<span class="human-badge">✓</span>'+
      '</div>';
    }).join('');

  }catch(e){console.error(e)}
}

setInterval(()=>{if(uptimeBase){uptimeBase++;document.getElementById('s-uptime').textContent=fmtUptime(uptimeBase);}},1000);
update();
setInterval(update,6000);
</script>
</body>
</html>`)
}
