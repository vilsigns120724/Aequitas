package keeper

import (
"encoding/json"
"fmt"
"net/http"
"time"
)

type APIServer struct {
blockchain *Blockchain
p2pNode    *P2PNode
keeper     *Keeper
startTime  time.Time
}

func NewAPIServer(bc *Blockchain, p2p *P2PNode, k *Keeper) *APIServer {
return &APIServer{
blockchain: bc,
p2pNode:    p2p,
keeper:     k,
startTime:  time.Now(),
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
json.NewEncoder(w).Encode(map[string]interface{}{
"chain_id":      "aequitas-1",
"version":       "v0.3.0",
"height":        latest.Height,
"latest_hash":   latest.Hash,
"total_humans":  a.keeper.TotalHumans(),
"total_supply":  a.keeper.TotalHumans() * 1000,
"node_id":       a.p2pNode.GetNodeID(),
"uptime":        uptime,
"block_time":    6,
"sepolia_sync":  "https://aequitas-proof-server-production.up.railway.app",
"contract_v5":   "0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5",
"explorer":      "https://sepolia.etherscan.io/address/0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5",
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
:root{--green:#3fb950;--blue:#58a6ff;--gold:#e3b341;--red:#f85149;--border:#21262d;--card:#0d1117;--bg:#060810}

/* HEADER */
header{background:#0d1117;border-bottom:1px solid var(--border);padding:18px 32px;display:flex;align-items:center;justify-content:space-between;position:sticky;top:0;z-index:100}
.logo{display:flex;align-items:center;gap:12px}
.logo-icon{font-size:1.5rem}
.logo-text{font-size:1.2rem;font-weight:bold;color:#f0f6fc;letter-spacing:2px}
.logo-sub{font-size:0.7rem;color:#8b949e;letter-spacing:1px}
.header-right{display:flex;align-items:center;gap:16px}
.live-badge{display:flex;align-items:center;gap:6px;background:#1a3a2a;border:1px solid #2ea04326;padding:5px 12px;border-radius:20px;font-size:0.72rem;color:var(--green)}
.pulse{width:7px;height:7px;background:var(--green);border-radius:50%;animation:pulse 2s infinite}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.4;transform:scale(0.8)}}
.chain-id{font-size:0.72rem;color:#8b949e;background:#161b22;border:1px solid var(--border);padding:5px 12px;border-radius:20px}

/* STATS GRID */
.stats{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:12px;padding:24px 32px 0}
.stat{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;position:relative;overflow:hidden;transition:border-color 0.2s}
.stat:hover{border-color:#388bfd40}
.stat::before{content:'';position:absolute;top:0;left:0;right:0;height:2px}
.stat.green::before{background:var(--green)}
.stat.blue::before{background:var(--blue)}
.stat.gold::before{background:var(--gold)}
.stat.red::before{background:var(--red)}
.stat-label{font-size:0.68rem;color:#8b949e;text-transform:uppercase;letter-spacing:1.5px;margin-bottom:10px}
.stat-value{font-size:1.9rem;font-weight:bold;line-height:1}
.stat.green .stat-value{color:var(--green)}
.stat.blue .stat-value{color:var(--blue)}
.stat.gold .stat-value{color:var(--gold)}
.stat.red .stat-value{color:var(--red)}
.stat-sub{font-size:0.7rem;color:#8b949e;margin-top:6px}

/* MAIN LAYOUT */
.main{display:grid;grid-template-columns:1fr 380px;gap:16px;padding:24px 32px;align-items:start}
@media(max-width:900px){.main{grid-template-columns:1fr}}

/* SECTIONS */
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;overflow:hidden}
.section-header{padding:16px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between}
.section-title{font-size:0.75rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px}
.section-count{font-size:0.7rem;color:#8b949e;background:#161b22;padding:2px 8px;border-radius:10px}

/* BLOCKS */
.block-item{padding:12px 20px;border-bottom:1px solid #161b22;display:grid;grid-template-columns:80px 1fr 70px 90px;gap:12px;align-items:center;transition:background 0.15s}
.block-item:hover{background:#161b2280}
.block-item:last-child{border-bottom:none}
.block-height{color:var(--blue);font-weight:bold;font-size:0.85rem}
.block-hash{color:#8b949e;font-size:0.75rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.block-humans{color:var(--gold);font-size:0.75rem;text-align:right}
.block-time{color:var(--green);font-size:0.75rem;text-align:right}

/* HUMANS */
.human-item{padding:12px 20px;border-bottom:1px solid #161b22;display:flex;align-items:center;gap:12px}
.human-item:last-child{border-bottom:none}
.human-avatar{width:32px;height:32px;background:linear-gradient(135deg,#1f6feb,#388bfd);border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:0.75rem;color:white;font-weight:bold;flex-shrink:0}
.human-addr{font-size:0.75rem;color:#8b949e;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.human-badge{margin-left:auto;font-size:0.65rem;color:var(--green);background:#1a3a2a;padding:2px 8px;border-radius:10px;flex-shrink:0}

/* INFO PANEL */
.info-panel{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;margin-bottom:16px}
.info-title{font-size:0.7rem;color:#8b949e;text-transform:uppercase;letter-spacing:2px;margin-bottom:14px}
.info-row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid #161b22}
.info-row:last-child{border-bottom:none}
.info-key{font-size:0.72rem;color:#8b949e}
.info-val{font-size:0.72rem;color:#f0f6fc;text-align:right;max-width:60%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.info-val.green{color:var(--green)}
.info-val.blue{color:var(--blue)}
.info-val.gold{color:var(--gold)}

/* LINKS */
.link-row{display:flex;gap:8px;flex-wrap:wrap;margin-top:12px}
.ext-link{font-size:0.7rem;color:var(--blue);background:#1c2d40;border:1px solid #1f6feb40;padding:5px 12px;border-radius:6px;text-decoration:none;transition:background 0.2s}
.ext-link:hover{background:#1f6feb20}

/* UPTIME */
#uptime-display{font-variant-numeric:tabular-nums}
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
    <div class="chain-id">aequitas-1</div>
  </div>
</header>

<div class="stats">
  <div class="stat blue">
    <div class="stat-label">Block Height</div>
    <div class="stat-value" id="s-height">—</div>
    <div class="stat-sub">every 6 seconds</div>
  </div>
  <div class="stat green">
    <div class="stat-label">Verified Humans</div>
    <div class="stat-value" id="s-humans">—</div>
    <div class="stat-sub">Proof of Humanity</div>
  </div>
  <div class="stat gold">
    <div class="stat-label">Total Supply</div>
    <div class="stat-value" id="s-supply">—</div>
    <div class="stat-sub">AEQ · dynamic cap</div>
  </div>
  <div class="stat green">
    <div class="stat-label">Uptime</div>
    <div class="stat-value" id="s-uptime" style="font-size:1.2rem">—</div>
    <div class="stat-sub" id="s-version">—</div>
  </div>
</div>

<div class="main">
  <!-- LEFT: BLOCKS -->
  <div>
    <div class="section">
      <div class="section-header">
        <span class="section-title">Recent Blocks</span>
        <span class="section-count" id="block-count">—</span>
      </div>
      <div id="blocks-list"></div>
    </div>
  </div>

  <!-- RIGHT: INFO + HUMANS -->
  <div>
    <div class="info-panel">
      <div class="info-title">Node Info</div>
      <div class="info-row"><span class="info-key">Node ID</span><span class="info-val blue" id="i-nodeid">—</span></div>
      <div class="info-row"><span class="info-key">Latest Hash</span><span class="info-val" id="i-hash">—</span></div>
      <div class="info-row"><span class="info-key">Block Time</span><span class="info-val green">6 seconds</span></div>
      <div class="info-row"><span class="info-key">Contract V5</span><span class="info-val blue">0x4f147d...f0B8b5</span></div>
      <div class="info-row"><span class="info-key">Sepolia Sync</span><span class="info-val green">✓ Active</span></div>
      <div class="info-row"><span class="info-key">Fee</span><span class="info-val gold">0.1%</span></div>
      <div class="info-row"><span class="info-key">Initial Grant</span><span class="info-val gold">1,000 AEQ</span></div>
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
let uptimeBase = 0;

function fmt(n){return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g,',')}
function fmtUptime(s){
  const h=Math.floor(s/3600),m=Math.floor((s%3600)/60),sec=s%60;
  return (h?h+'h ':'')+m+'m '+sec+'s';
}

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
    document.getElementById('s-version').textContent=status.version;
    uptimeBase=status.uptime;

    // Node info
    const nid=status.node_id||'';
    document.getElementById('i-nodeid').textContent=nid.slice(0,16)+'...';
    document.getElementById('i-hash').textContent=status.latest_hash.slice(0,12)+'...';

    // Blocks
    const bl=document.getElementById('blocks-list');
    document.getElementById('block-count').textContent=status.height+' blocks';
    bl.innerHTML=[...blocks].reverse().map(b=>{
      const t=new Date(b.timestamp*1000).toLocaleTimeString();
      return '<div class="block-item">'+
        '<span class="block-height">#'+b.height+'</span>'+
        '<span class="block-hash">'+b.hash+'</span>'+
        '<span class="block-humans">👤 '+b.humans+'</span>'+
        '<span class="block-time">'+t+'</span>'+
      '</div>';
    }).join('');

    // Humans
    const hl=document.getElementById('humans-list');
    document.getElementById('human-count').textContent=(humans.total||0)+' registered';
    const list=humans.humans||[];
    hl.innerHTML=list.map((h,i)=>{
      const addr=h.address||'';
      const short=addr.length>12?addr.slice(0,10)+'...':addr;
      const init=addr.slice(0,2).toUpperCase();
      return '<div class="human-item">'+
        '<div class="human-avatar">'+init+'</div>'+
        '<span class="human-addr">'+addr+'</span>'+
        '<span class="human-badge">✓ verified</span>'+
      '</div>';
    }).join('');

  }catch(e){console.error(e)}
}

// Uptime ticker
setInterval(()=>{
  if(uptimeBase){
    uptimeBase++;
    document.getElementById('s-uptime').textContent=fmtUptime(uptimeBase);
  }
},1000);

update();
setInterval(update,6000);
</script>
</body>
</html>`)
}
