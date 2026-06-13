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
}

func NewAPIServer(bc *Blockchain, p2p *P2PNode, k *Keeper) *APIServer {
return &APIServer{
blockchain: bc,
p2pNode:    p2p,
keeper:     k,
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
json.NewEncoder(w).Encode(map[string]interface{}{
"chain_id":     "aequitas-1",
"version":      "v0.2.0",
"height":       latest.Height,
"latest_hash":  latest.Hash,
"total_humans": a.keeper.TotalHumans(),
"total_supply": a.keeper.TotalHumans() * 1000,
"node_id":      a.p2pNode.GetNodeID(),
"uptime":       time.Now().Unix(),
})
}

func (a *APIServer) handleBlocks(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
blocks := a.blockchain.GetBlocks()
// Return last 20 blocks
start := 0
if len(blocks) > 20 {
start = len(blocks) - 20
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
<title>Aequitas Chain</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { background: #0a0a0f; color: #e0e0e0; font-family: 'Courier New', monospace; }
  header { background: #0d1117; border-bottom: 1px solid #21262d; padding: 20px 40px; display: flex; align-items: center; gap: 16px; }
  header h1 { font-size: 1.4rem; color: #f0f0f0; }
  header .tag { background: #1a3a2a; color: #3fb950; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; padding: 30px 40px 0; }
  .card { background: #0d1117; border: 1px solid #21262d; border-radius: 10px; padding: 20px; }
  .card .label { color: #8b949e; font-size: 0.75rem; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 1px; }
  .card .value { font-size: 1.8rem; font-weight: bold; color: #3fb950; }
  .card .value.blue { color: #58a6ff; }
  .card .value.yellow { color: #e3b341; }
  .section { padding: 30px 40px; }
  .section h2 { color: #8b949e; font-size: 0.85rem; text-transform: uppercase; letter-spacing: 2px; margin-bottom: 16px; }
  .block { background: #0d1117; border: 1px solid #21262d; border-radius: 8px; padding: 14px 18px; margin-bottom: 8px; display: flex; align-items: center; gap: 20px; }
  .block .height { color: #58a6ff; font-weight: bold; min-width: 80px; }
  .block .hash { color: #8b949e; font-size: 0.85rem; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .block .time { color: #3fb950; font-size: 0.8rem; min-width: 80px; text-align: right; }
  .block .humans { color: #e3b341; font-size: 0.8rem; min-width: 80px; text-align: right; }
  .pulse { display: inline-block; width: 8px; height: 8px; background: #3fb950; border-radius: 50%; animation: pulse 2s infinite; margin-right: 8px; }
  @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.3} }
</style>
</head>
<body>
<header>
  <span class="pulse"></span>
  <h1>⚖ Aequitas Chain</h1>
  <span class="tag">LIVE</span>
</header>

<div class="grid">
  <div class="card"><div class="label">Block Height</div><div class="value" id="height">...</div></div>
  <div class="card"><div class="label">Total Humans</div><div class="value yellow" id="humans">...</div></div>
  <div class="card"><div class="label">Total Supply</div><div class="value blue" id="supply">...</div></div>
  <div class="card"><div class="label">Latest Hash</div><div class="value" style="font-size:0.8rem;word-break:break-all" id="hash">...</div></div>
</div>

<div class="section">
  <h2>Recent Blocks</h2>
  <div id="blocks"></div>
</div>

<script>
async function update() {
  try {
    const [status, blocks] = await Promise.all([
      fetch('/api/status').then(r => r.json()),
      fetch('/api/blocks').then(r => r.json())
    ]);
    document.getElementById('height').textContent = '#' + status.height;
    document.getElementById('humans').textContent = status.total_humans;
    document.getElementById('supply').textContent = status.total_supply + ' AEQ';
    document.getElementById('hash').textContent = status.latest_hash.slice(0,16) + '...';
    const blocksEl = document.getElementById('blocks');
    blocksEl.innerHTML = [...blocks].reverse().map(b => {
      const t = new Date(b.timestamp * 1000).toLocaleTimeString();
      return '<div class="block"><span class="height">#' + b.height + '</span><span class="hash">' + b.hash + '</span><span class="humans">👤 ' + b.humans + '</span><span class="time">' + t + '</span></div>';
    }).join('');
  } catch(e) {}
}
update();
setInterval(update, 6000);
</script>
</body>
</html>`)
}
