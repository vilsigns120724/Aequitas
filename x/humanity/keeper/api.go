package keeper

import (
"encoding/json"
"fmt"
"io"
"math/big"
"net/http"
"os"
"strings"
"time"

"github.com/ethereum/go-ethereum/common"
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
mux.HandleFunc("/api/check-registration", a.handleCheckRegistration)
mux.HandleFunc("/api/check-registration-by-biohash", a.handleCheckRegistrationByBioHash)
mux.HandleFunc("/api/check-nullifier", a.handleCheckNullifier)
mux.HandleFunc("/api/swap", a.handleSwap)
mux.HandleFunc("/api/add-liquidity", a.handleAddLiquidity)
mux.HandleFunc("/api/remove-liquidity", a.handleRemoveLiquidity)
mux.HandleFunc("/api/lp-position", a.handleLPPosition)
mux.HandleFunc("/api/faucet", a.handleFaucet)
mux.HandleFunc("/api/pool", a.handlePoolStatus)
mux.HandleFunc("/api/snapshot", a.handleSnapshot)
mux.HandleFunc("/api/peers", a.handlePeers)
mux.HandleFunc("/api/peers/register", a.handlePeerRegister)
mux.HandleFunc("/registered", a.handleRegistered)
mux.HandleFunc("/download/app.apk", a.handleAppDownload)
fmt.Println("── Starting EVM RPC ─────────────────────")
evmRPC := NewEVMRPCServer(a.blockchain, a.state)
mux.HandleFunc("/rpc", evmRPC.handleRPC)
if evmRPC.evm != nil {
fmt.Println("✓ EVM Engine ready")
// Ensure V7 contract is deployed — redeploys from hardcoded bytecode
// if missing (e.g. after a DB reset). Without this the node fails with
// "no code at address" on every registration attempt.
deployerAddr := os.Getenv("RELAYER_ADDRESS")
if deployerAddr == "" {
deployerAddr = "0x0BE8b961CBf6564bd1931B0803D35C0659E0D016"
}
EnsureContractsDeployed(evmRPC.evm, a.state, deployerAddr)
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
// Use a.state (PostgreSQL-backed ChainState) as the single source of
// truth for human count, not a.keeper — the in-memory Keeper map is never
// persisted and resets to 0 on every restart, which previously made this
// "growth" figure silently diverge from total_humans below (which already
// correctly used a.state).
humans := a.state.TotalHumans()
growth := humans * 10
if growth > 100 {
growth = 100
}
// Calculate time until next UBI distribution (24h after server start)
uptimeSecs := time.Since(a.startTime).Seconds()
nextUBISecs := int64(86400 - int(uptimeSecs)%86400)
if nextUBISecs < 0 {
nextUBISecs = 0
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
"contract_v5":  V5_SEPOLIA_LEGACY_ADDR,
"contract_v6":  V6_CONTRACT_ADDR,
"contract_v7":  V7_CONTRACT_ADDR,
"bio_verifier": BIO_VERIFIER_ADDR,
"chain_evm_id": 1926,
"index":        a.state.CalcAequitasIndex(),
"gini":         a.state.CalcGini(),
"growth":       growth,
"velocity":     50,
"phase":        a.state.CalcPhase(),
"fee_bps":      10,
"pool_validators": fmt.Sprintf("%.4f", a.state.GetBalance("0x78c1c143e395b181f13bcb6868ff53aa86c3d2ba")),
"pool_lp":         fmt.Sprintf("%.4f", a.state.GetBalance("0xc181c3a4d09444b99089ae0f56c1e7f4c20d01eb")),
"pool_ubi":        fmt.Sprintf("%.4f", a.state.GetBalance("0x4a9b8f99f0d8cff0e510fef502100571203b054a")),
"pool_treasury":   fmt.Sprintf("%.4f", a.state.GetBalance("0x2273894fb781978d54e767f9fba2dcb33d93eb15")),
"ubi_next_payout_secs": nextUBISecs,
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
json.NewEncoder(w).Encode(map[string]interface{}{"balance": 0, "tusd_balance": 0, "is_human": false})
return
}

// Use ChainState (native balance) as the single source of truth.
// This used to query the V7 contract's balanceOf()/isHuman() directly,
// which was the right call back when registrations only wrote to EVM
// storage and ChainState was never updated. Since registration now also
// grants the native balance via state.RegisterHuman() (and transfers
// move the native balance via state.Transfer()), ChainState reflects
// the real, current state — while the contract's own balanceOf() can
// lag behind it (it's no longer touched by ordinary native transfers
// at all, and read-only contract calls are intentionally not persisted
// per-call). Querying the contract here would show a wallet's balance
// from whenever it last interacted with the contract directly, not its
// real current native balance.
balance := a.state.GetBalance(wallet)
tusdBalance := a.state.GetTUsdBalance(wallet)
isHuman := a.state.IsHuman(wallet)
demurrage := a.state.GetDemurrageStatus(wallet)

json.NewEncoder(w).Encode(map[string]interface{}{
"wallet":                   wallet,
"balance":                  balance,
"tusd_balance":              tusdBalance,
"is_human":                  isHuman,
"demurrage_active":          demurrage.Active,
"demurrage_days_until_start": demurrage.DaysUntilStart,
"show_14_day_notice":        demurrage.ShowFourteenDayNotice,
"show_7_day_notice":         demurrage.ShowSevenDayNotice,
})
}

// handleCheckRegistration lets the app ask "did MY specific proof commitment
// get registered, and to which wallet?" — instead of reading the last entry
// in a global, unfiltered /api/humans list (which showed every user the
// most recently registered wallet, regardless of who they actually were).
func (a *APIServer) handleCheckRegistration(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")

commitment := r.URL.Query().Get("commitment")
if commitment == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"registered": false})
return
}

wallet := a.state.GetWalletByCommitment(commitment)
if wallet == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"registered": false})
return
}

balance := a.state.GetBalance(wallet)
isHuman := a.state.IsHuman(wallet)
json.NewEncoder(w).Encode(map[string]interface{}{
"registered": true,
"wallet":     wallet,
"balance":    balance,
"is_human":   isHuman,
})
}

// handleCheckRegistrationByBioHash mirrors handleCheckRegistration, but
// keyed by the device's biometric identity hash rather than a proof
// commitment. Needed because, under the new website-side proof flow, the
// app only ever knows its own bioHash — it never computes a commitment
// itself anymore (that now happens on the website, after MetaMask
// supplies the real wallet) — so it can't poll by commitment the way the
// old flow did.
func (a *APIServer) handleCheckRegistrationByBioHash(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
if r.Method == "OPTIONS" { w.WriteHeader(200); return }

// Accept bioHash via POST body (preferred — keeps biometric identifier
// out of server logs) or GET query string (legacy, for old app builds).
var bioHash string
if r.Method == "POST" {
r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
var body struct{ BioHash string `json:"bioHash"` }
json.NewDecoder(r.Body).Decode(&body)
bioHash = body.BioHash
} else {
bioHash = r.URL.Query().Get("bioHash")
}
if bioHash == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"registered": false})
return
}

wallet := a.state.GetWalletByBioHash(bioHash)
if wallet == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"registered": false})
return
}

balance := a.state.GetBalance(wallet)
isHuman := a.state.IsHuman(wallet)

// If the bioHash exists in bio_registrations but the wallet is NOT yet
// marked as human on-chain, it means someone else used this biometric
// hash to generate a proof but hasn't completed registration yet —
// OR a different wallet tried to reuse this bioHash. Either way, the
// current user should NOT see "success". Return a distinct status so
// the app can show an appropriate message.
if !isHuman {
json.NewEncoder(w).Encode(map[string]interface{}{
"registered":       false,
"biometric_in_use": true,
"wallet":           wallet,
})
return
}

json.NewEncoder(w).Encode(map[string]interface{}{
"registered": true,
"wallet":     wallet,
"balance":    balance,
"is_human":   isHuman,
})
}

// handleCheckNullifier lets the client ask "has this nullifier been used?"
// before submitting a registration. GET /api/check-nullifier?n=<hex>
func (a *APIServer) handleCheckNullifier(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
nullifier := r.URL.Query().Get("n")
if nullifier == "" {
json.NewEncoder(w).Encode(map[string]interface{}{"used": false})
return
}
wallet := a.state.GetWalletByNullifier(nullifier)
json.NewEncoder(w).Encode(map[string]interface{}{
"used":   wallet != "",
"wallet": wallet,
})
}

// queryV7Status reads isHuman(address) and balanceOf(address) directly from
// the V7 contract via eth_call. Kept available for debugging/comparison
// against the contract's own bookkeeping, but no longer used by the
// balance-facing endpoints above — see handleBalance for why.
func (a *APIServer) queryV7Status(wallet string) (float64, bool) {
evmRPC := NewEVMRPCServer(a.blockchain, a.state)
if evmRPC.evm == nil {
return 0, false
}

to := common.HexToAddress(V7_CONTRACT_ADDR)
from := common.HexToAddress(wallet)

// isHuman(address) — selector 0xf72c436f
// persist=false: this is a read-only status query (used by the explorer
// frontend's balance/status display), not a real registration. Previously
// every poll of this endpoint silently wrote isHuman=true to evm_storage
// as a side effect, which is part of why "already registered" kept
// reappearing even right after a full database reset.
isHumanData := append(common.Hex2Bytes("f72c436f"), common.LeftPadBytes(from.Bytes(), 32)...)
isHumanRet, err := evmRPC.evm.CallContract(from, to, isHumanData, big.NewInt(0), false)
isHuman := false
if err == nil && len(isHumanRet) >= 32 {
isHuman = isHumanRet[31] == 1
}

if !isHuman {
return 0, false
}

// balanceOf(address) — selector 0x70a08231
balanceData := append(common.Hex2Bytes("70a08231"), common.LeftPadBytes(from.Bytes(), 32)...)
balanceRet, err := evmRPC.evm.CallContract(from, to, balanceData, big.NewInt(0), false)
balance := 0.0
if err == nil && len(balanceRet) >= 32 {
weiInt := new(big.Int).SetBytes(balanceRet)
decimals := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
balanceFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(weiInt), decimals).Float64()
balance = balanceFloat
}

return balance, isHuman
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
.box{background:#111827;border:1px solid #1E2D45;border-radius:12px;padding:32px;max-width:440px;width:100%%}
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
w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
path := strings.Trim(r.URL.Path, "/")
if idx := strings.Index(path, "/"); idx >= 0 {
	path = path[:idx]
}
// Backwards-compat: /swap redirects to /exchange.
if path == "swap" {
	http.Redirect(w, r, "/exchange", http.StatusMovedPermanently)
	return
}
// For non-register tabs, swap the active class server-side so the correct
// tab is visible on direct URL load — no JS timing dependency.
validTabs := map[string]bool{"explorer": true, "index": true, "network": true, "exchange": true}
if validTabs[path] {
	html := strings.Replace(explorerHTML,
		`class="tab active" onclick="showTab('register',this)"`,
		`class="tab" onclick="showTab('register',this)"`, 1)
	html = strings.Replace(html,
		`class="tab" onclick="showTab('`+path+`',this)"`,
		`class="tab active" onclick="showTab('`+path+`',this)"`, 1)
	html = strings.Replace(html,
		`id="tab-register" class="tab-content active"`,
		`id="tab-register" class="tab-content"`, 1)
	html = strings.Replace(html,
		`id="tab-`+path+`" class="tab-content"`,
		`id="tab-`+path+`" class="tab-content active"`, 1)
	fmt.Fprint(w, html)
	return
}
fmt.Fprint(w, explorerHTML)
}

// handlePeerRegister accepts a node registration and returns the current peer list.
// POST /api/peers/register  body: {"url":"https://..."}
func (a *APIServer) handlePeerRegister(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
if r.Method == "OPTIONS" { w.WriteHeader(200); return }
var req struct{ URL string `json:"url"` }
r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
json.NewDecoder(r.Body).Decode(&req)
if req.URL != "" {
GlobalPeerRegistry.Register(req.URL)
fmt.Printf("[PEERS] Registered: %s\n", req.URL)
// Also start syncing this new peer from our side
a.blockchain.startSyncForPeer(req.URL)
}
json.NewEncoder(w).Encode(map[string]interface{}{"peers": GlobalPeerRegistry.AllPeers()})
}

// handlePeers returns the list of all known peer nodes.
// GET /api/peers
func (a *APIServer) handlePeers(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
json.NewEncoder(w).Encode(map[string]interface{}{"peers": GlobalPeerRegistry.AllPeers()})
}

// handleSnapshot exports the full Go-state as a signed JSON snapshot.
// Protected by SNAPSHOT_TOKEN env var if set. A new node can bootstrap
// itself by setting BOOTSTRAP_SNAPSHOT_URL to this endpoint's URL.
func (a *APIServer) handleSnapshot(w http.ResponseWriter, r *http.Request) {
if token := os.Getenv("SNAPSHOT_TOKEN"); token != "" {
if r.URL.Query().Get("token") != token {
http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
return
}
}
snap := a.state.ExportSnapshot(a.blockchain.GetSigningKey())
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
json.NewEncoder(w).Encode(snap)
}

func (a *APIServer) handleAppDownload(w http.ResponseWriter, r *http.Request) {
const apkPath = "downloads/aequitas-app.apk"
const fallbackURL = "https://github.com/hanoi96international-gif/Aequitas/raw/main/downloads/aequitas-app.apk"
f, err := os.Open(apkPath)
if err != nil {
	// File not found in container — redirect to GitHub raw URL.
	http.Redirect(w, r, fallbackURL, http.StatusFound)
	return
}
defer f.Close()
fi, err := f.Stat()
if err != nil {
	http.Redirect(w, r, fallbackURL, http.StatusFound)
	return
}
w.Header().Set("Content-Disposition", "attachment; filename=aequitas-app.apk")
w.Header().Set("Content-Type", "application/vnd.android.package-archive")
w.Header().Set("Access-Control-Allow-Origin", "*")
http.ServeContent(w, r, "aequitas-app.apk", fi.ModTime(), f)
}
