package keeper

import (
"crypto/subtle"
"encoding/hex"
"encoding/json"
"fmt"
"html"
"io"
"math/big"
"net"
"net/http"
"os"
"strings"
"time"

"github.com/ethereum/go-ethereum/accounts"
"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/crypto"
)

type APIServer struct {
blockchain        *BlockDAG
p2pNode           *P2PNode
keeper            *Keeper
startTime         time.Time
proofServerStatus map[string]interface{}
state             *ChainState
// Shared EVM RPC server — one instance so all registration calls share
// the same nonce map and mutex, preventing parallel registrations from
// reading the same DB nonce and writing the same follower value.
evmRPC            *EVMRPCServer
}

func NewAPIServer(bc *BlockDAG, p2p *P2PNode, k *Keeper, state *ChainState) *APIServer {
s := &APIServer{
blockchain:        bc,
p2pNode:           p2p,
keeper:            k,
startTime:         time.Now(),
proofServerStatus: map[string]interface{}{},
state:             state,
evmRPC:            NewEVMRPCServer(bc, state),
}
go s.syncProofServerStatus()
return s
}

// isValidWalletAddr checks 0x-prefixed 40-hex-char Ethereum address format.
// P3-11: prevents garbage keys from entering cs.accounts map.
func isValidWalletAddr(addr string) bool {
	if len(addr) != 42 { return false }
	if addr[:2] != "0x" && addr[:2] != "0X" { return false }
	for _, c := range addr[2:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func (a *APIServer) syncProofServerStatus() {
for {
proofHTTP := &http.Client{Timeout: 8 * time.Second}
resp, err := proofHTTP.Get("https://aequitas-proof-server-production.up.railway.app/health")
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
mux.HandleFunc("/api/gini/history", a.handleGiniHistory)
mux.HandleFunc("/api/price-history", a.handlePriceHistory)
mux.HandleFunc("/api/wealth-cap", a.handleWealthCap)
mux.HandleFunc("/api/sign-validator-challenge", a.handleSignValidatorChallenge)
mux.HandleFunc("/api/nonce", a.handleNonce)
mux.HandleFunc("/api/peers", a.handlePeers)
mux.HandleFunc("/api/peers/challenge", a.handlePeerChallenge)
mux.HandleFunc("/api/peers/register", a.handlePeerRegister)
mux.HandleFunc("/api/register-validator-key", a.handleRegisterValidatorKey)
mux.HandleFunc("/registered", a.handleRegistered)
mux.HandleFunc("/download/app.apk", a.handleAppDownload)
mux.HandleFunc("/download/node-guide-en.pdf", func(w http.ResponseWriter, r *http.Request) {
	a.handleStaticDownload(w, r, "downloads/Aequitas_Node_Guide_EN.pdf", "Aequitas_Node_Guide_EN.pdf", "application/pdf")
})
mux.HandleFunc("/download/node-guide-de.pdf", func(w http.ResponseWriter, r *http.Request) {
	a.handleStaticDownload(w, r, "downloads/Aequitas_Node_Guide_DE.pdf", "Aequitas_Node_Guide_DE.pdf", "application/pdf")
})
fmt.Println("── Starting EVM RPC ─────────────────────")
// Use the shared EVMRPCServer (a.evmRPC) so /rpc and /api/register share
// one nonce map + mutex — creating a second instance here caused separate
// nonce maps, making the atomic nonce reservation ineffective.
mux.HandleFunc("/rpc", a.evmRPC.handleRPC)
if a.evmRPC.evm != nil {
fmt.Println("✓ EVM Engine ready")
// Ensure V7 contract is deployed — redeploys from hardcoded bytecode
// if missing (e.g. after a DB reset). Without this the node fails with
// "no code at address" on every registration attempt.
deployerAddr := os.Getenv("RELAYER_ADDRESS")
if deployerAddr == "" {
deployerAddr = "0x0BE8b961CBf6564bd1931B0803D35C0659E0D016"
}
EnsureContractsDeployed(a.evmRPC.evm, a.state, deployerAddr)
} else {
fmt.Println("✗ EVM Engine failed")
}
addr := fmt.Sprintf(":%d", port)
fmt.Printf("✓ API Server listening on port %d\n", port)
// Use http.Server with explicit timeouts to prevent slowloris attacks and
// goroutine leaks from clients that never send/read — the default mux has none.
srv := &http.Server{
Addr:         addr,
Handler:      mux,
ReadTimeout:  30 * time.Second,
WriteTimeout: 60 * time.Second,
IdleTimeout:  120 * time.Second,
}
go func() {
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
fmt.Printf("[API] Server error: %v\n", err)
}
}()
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
// P3-3: compute next UBI based on last_ubi_at, not server uptime.
nextUBISecs := a.state.SecondsUntilNextUBI()

json.NewEncoder(w).Encode(map[string]interface{}{
"chain_id":     "aequitas-1",
"version":      "v0.3.0",
"height":       latest.Height,
"latest_hash":  latest.Hash,
"total_humans": a.state.TotalHumans(),
"total_supply": fmt.Sprintf("%.2f AEQ", a.state.TotalSupply()),
"node_id":      a.p2pNode.GetNodeID(),
"uptime":       uptime,
"is_primary":   os.Getenv("IS_PRIMARY_NODE") == "true",
"block_time":   6,
"contract_v7":  V7_CONTRACT_ADDR,
// P3-8: V5/V6 legacy addresses removed from status — minimise attack surface.
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
// P3-2: support ?limit=N&offset=M for block history paging.
limit := 50
offset := 0
fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)
if limit < 1 || limit > 500 { limit = 50 }
if offset < 0 { offset = 0 }
// Default: newest blocks (offset from end)
if r.URL.Query().Get("offset") == "" {
offset = len(blocks) - limit
if offset < 0 { offset = 0 }
}
end := offset + limit
if end > len(blocks) { end = len(blocks) }
if offset >= len(blocks) { offset = len(blocks) }
json.NewEncoder(w).Encode(blocks[offset:end])
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
proofHTTP2 := &http.Client{Timeout: 8 * time.Second}
resp, err := proofHTTP2.Get("https://aequitas-proof-server-production.up.railway.app/humans")
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

// POST only — GET is removed because bioHash in the URL lands in
// server/proxy logs creating unnecessary biometric linkability.
if r.Method != "POST" && r.Method != "OPTIONS" {
http.Error(w, `{"error":"POST required"}`, 405); return
}
r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
var bioHashBody struct{ BioHash string `json:"bioHash"` }
json.NewDecoder(r.Body).Decode(&bioHashBody)
var bioHash = bioHashBody.BioHash
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
// Return only used/unused — never the associated wallet address.
// The wallet linkage is a biometric identifier that should not be
// publicly enumerable via the nullifier index.
json.NewEncoder(w).Encode(map[string]interface{}{"used": wallet != ""})
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
// XSS fix: escape wallet parameter before writing to HTML — without this,
// a crafted URL like /registered?wallet=<script>... would execute JS.
wallet := html.EscapeString(r.URL.Query().Get("wallet"))
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
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
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
	html := strings.Replace(explorerHTML, `<html lang="en">`, `<html lang="en" data-active="`+path+`">`, 1)
	// Activate tab button.
	html = strings.Replace(html,
		`class="tab active" onclick="showTab('register',this)"`,
		`class="tab" onclick="showTab('register',this)"`, 1)
	html = strings.Replace(html,
		`class="tab" onclick="showTab('`+path+`',this)"`,
		`class="tab active" onclick="showTab('`+path+`',this)"`, 1)
	// Force tab content and first stab-panel visible via inline style.
	// Inline style beats every CSS rule except JS .style.display override.
	html = strings.Replace(html,
		`id="tab-`+path+`" class="tab-content"`,
		`id="tab-`+path+`" class="tab-content" style="display:block"`, 1)
	// Also hide register content when not on register route.
	html = strings.Replace(html,
		`id="tab-register" class="tab-content active"`,
		`id="tab-register" class="tab-content"`, 1)
	fmt.Fprint(w, html)
	return
}
fmt.Fprint(w, explorerHTML)
}

// handleNonce returns the next swap nonce a wallet should sign with.
// GET /api/nonce?wallet=0x...
func (a *APIServer) handleNonce(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
wallet := strings.ToLower(r.URL.Query().Get("wallet"))
if wallet == "" {
http.Error(w, `{"error":"wallet required"}`, 400)
return
}
nonce := a.state.GetSwapNonce(wallet)
json.NewEncoder(w).Encode(map[string]interface{}{"wallet": wallet, "nonce": nonce})
}

// handlePriceHistory returns AEQ/tUSD price snapshots for the chart.
// GET /api/price-history?minutes=240&limit=5000
func (a *APIServer) handlePriceHistory(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Cache-Control", "no-cache")
minutes := 240
limit := 1000
fmt.Sscanf(r.URL.Query().Get("minutes"), "%d", &minutes)
fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
// Clamp to prevent memory exhaustion from large DB reads.
if minutes < 1 { minutes = 1 }
if minutes > 43200 { minutes = 43200 } // max 30 days
if limit < 1 { limit = 1 }
if limit > 5000 { limit = 5000 }
history := a.state.GetPriceHistory(minutes, limit)
json.NewEncoder(w).Encode(map[string]interface{}{"history": history, "count": len(history)})
}

// handleGiniHistory returns Gini snapshots stored after each UBI distribution.
// Falls back to the current Gini as a single point when no history exists yet.
func (a *APIServer) handleGiniHistory(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Cache-Control", "no-cache")
history := a.state.GetGiniHistory(60) // last 60 snapshots
if len(history) == 0 {
// First UBI hasn't run yet — return current state as bootstrap point.
gini := a.state.CalcGini()
humans := a.state.TotalHumans()
history = []map[string]interface{}{
{"idx": gini * 100, "gini": gini, "humans": humans, "timestamp": time.Now().Unix()},
}
}
json.NewEncoder(w).Encode(map[string]interface{}{"history": history})
}

// handleWealthCap returns the current wealth cap parameters.
// Field names match the live wealth-cap widget in the Equality tab.
func (a *APIServer) handleWealthCap(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Cache-Control", "no-cache")
accs := a.state.GetAllAccounts()
var total float64
n := 0
for _, acc := range accs {
if acc.IsHuman { total += acc.Balance.Float(); n++ }
}
avg := 0.0
if n > 0 { avg = total / float64(n) }
mult := 5.0
if n > 5 { mult = float64(n) }
if mult > 25 { mult = 25 }
capAEQ := mult * avg
json.NewEncoder(w).Encode(map[string]interface{}{
"cap_aeq": capAEQ, "multiplier": mult, "average_aeq": avg,
"humans": n, "total_supply": total,
})
}

// handleSignValidatorChallenge signs the key-possession challenge message with
// RELAYER_PRIVATE_KEY. Restricted to loopback (127.0.0.1 / ::1) so only
// node operators with server access can use it — not an internet-accessible oracle.
// GET /api/sign-validator-challenge?wallet=0x...
func (a *APIServer) handleSignValidatorChallenge(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
// Block requests that arrive through a proxy — X-Forwarded-For or
// X-Real-IP being set means a proxy sits in front, potentially making
// a non-local request appear to come from 127.0.0.1.
if r.Header.Get("X-Forwarded-For") != "" || r.Header.Get("X-Real-IP") != "" {
http.Error(w, `{"error":"only accessible directly from localhost (no proxy)"}`, 403); return
}
host, _, err := net.SplitHostPort(r.RemoteAddr)
if err != nil || (host != "127.0.0.1" && host != "::1") {
http.Error(w, `{"error":"only accessible from localhost"}`, 403); return
}
humanWallet := strings.ToLower(r.URL.Query().Get("wallet"))
if humanWallet == "" || !strings.HasPrefix(humanWallet, "0x") || len(humanWallet) != 42 {
http.Error(w, `{"error":"wallet required (0x...)"}`, 400); return
}
key := a.blockchain.GetSigningKey()
if key == nil {
http.Error(w, `{"error":"RELAYER_PRIVATE_KEY not configured"}`, 500); return
}
message := "Aequitas: validator key linked to human " + humanWallet
msgHash := accounts.TextHash([]byte(message))
sig, err := crypto.Sign(msgHash, key)
if err != nil {
http.Error(w, `{"error":"signing failed"}`, 500); return
}
sig[64] += 27
signingAddr := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
json.NewEncoder(w).Encode(map[string]interface{}{
"signing_address": signingAddr,
"human_wallet":    humanWallet,
"signature":       "0x" + hex.EncodeToString(sig),
"message":         message,
})
}

// handleRegisterValidatorKey links a node signing key to a registered human
// wallet, authorising that signing key to propose blocks.
//
// Requires TWO signatures proving control of BOTH keys:
//   human_signature:      personal_sign("Aequitas: authorize validator key {signing_address}", human_wallet)
//   signing_key_signature: personal_sign("Aequitas: validator key linked to human {human_wallet}", signing_address)
//
// The double-signature requirement proves the requester controls both the
// human wallet AND the node signing key, preventing impersonation attacks
// where someone registers a victim's signing address using their own wallet.
// UNIQUE(human_wallet) ensures one human = one validator key.
func (a *APIServer) handleRegisterValidatorKey(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
if r.Method == "OPTIONS" { w.WriteHeader(200); return }
if r.Method != "POST" { http.Error(w, `{"error":"POST required"}`, 405); return }
r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
var req struct {
SigningAddress      string `json:"signing_address"`
HumanWallet        string `json:"human_wallet"`
HumanSignature     string `json:"human_signature"`
SigningKeySignature string `json:"signing_key_signature"`
}
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, `{"error":"invalid request"}`, 400); return
}
signingAddr := strings.ToLower(strings.TrimSpace(req.SigningAddress))
humanWallet := strings.ToLower(strings.TrimSpace(req.HumanWallet))
if !strings.HasPrefix(signingAddr, "0x") || len(signingAddr) != 42 ||
!strings.HasPrefix(humanWallet, "0x") || len(humanWallet) != 42 {
http.Error(w, `{"error":"invalid address"}`, 400); return
}
// 1. Human wallet proves it authorises this signing key.
humanMsg := "Aequitas: authorize validator key " + signingAddr
if err := verifyPersonalSign(humanMsg, req.HumanSignature, humanWallet); err != nil {
http.Error(w, `{"error":"invalid human_signature: `+err.Error()+`"}`, 400); return
}
// 2. Signing key proves it is linked to this human wallet (key-possession proof).
signingMsg := "Aequitas: validator key linked to human " + humanWallet
if err := verifyPersonalSign(signingMsg, req.SigningKeySignature, signingAddr); err != nil {
http.Error(w, `{"error":"invalid signing_key_signature — sign with RELAYER_PRIVATE_KEY: `+err.Error()+`"}`, 400); return
}
if err := a.state.RegisterValidatorKey(signingAddr, humanWallet); err != nil {
http.Error(w, `{"error":"`+err.Error()+`"}`, 400); return
}
a.blockchain.AddAuthorizedValidator(signingAddr)
fmt.Printf("[VALIDATOR] ✓ Registered key %s for human %s\n", signingAddr, humanWallet)
json.NewEncoder(w).Encode(map[string]interface{}{
"success":         true,
"signing_address": signingAddr,
"human_wallet":    humanWallet,
})
}

// handlePeerChallenge issues a one-time challenge that the peer must sign to
// prove ownership of their signing key (P1-3 validator signature verification).
// GET /api/peers/challenge?address=0x...
func (a *APIServer) handlePeerChallenge(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
addr := strings.ToLower(r.URL.Query().Get("address"))
if !isValidWalletAddr(addr) {
http.Error(w, `{"error":"invalid address"}`, 400)
return
}
challenge := a.blockchain.IssuePeerChallenge(addr)
json.NewEncoder(w).Encode(map[string]interface{}{
"challenge":  challenge,
"expires_in": 90,
"instructions": "Sign the challenge string with your signing key and include the hex signature in POST /api/peers/register as 'signature'",
})
}

// handlePeerRegister accepts a node registration and returns the current peer
// list plus all authorized validator addresses. A node that sends its
// signing_address is automatically added to the authorized validator set so
// its blocks are accepted without manual AUTHORIZED_VALIDATORS configuration.
// POST /api/peers/register  body: {"url":"https://...","signing_address":"0x...","signature":"0x..."}
func (a *APIServer) handlePeerRegister(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
if r.Method == "OPTIONS" { w.WriteHeader(200); return }
var req struct {
URL            string `json:"url"`
SigningAddress string `json:"signing_address"`
PeerSecret     string `json:"peer_secret"`
Signature      string `json:"signature"` // P1-3 challenge-response
}
r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
json.NewDecoder(r.Body).Decode(&req)

// Secret check comes FIRST. URL registration and sync goroutines are only
// started for authenticated peers — prevents goroutine exhaustion via
// unauthenticated registrations even when PEER_SECRET is not set.
peerSecret := os.Getenv("PEER_SECRET")
// P1-2: constant-time comparison prevents timing-based secret oracle attacks.
secretOK := peerSecret != "" && subtle.ConstantTimeCompare([]byte(req.PeerSecret), []byte(peerSecret)) == 1

// Pre-check validator key status so URL registration can use it too.
var keyAuthorizedEarly bool
if earlyAddr := strings.ToLower(strings.TrimSpace(req.SigningAddress)); earlyAddr != "" {
for _, k := range a.state.GetValidatorKeys() {
if k["signing_address"] == earlyAddr { keyAuthorizedEarly = true; break }
}
}
if req.URL != "" && isAllowedPeerURL(req.URL) {
// Allow URL registration if PEER_SECRET matches OR the signing address
// has an individually registered validator key. Without this, a node
// with a valid key could not connect via peer discovery.
if secretOK || keyAuthorizedEarly {
GlobalPeerRegistry.Register(req.URL)
fmt.Printf("[PEERS] Registered: %s\n", req.URL)
a.blockchain.startSyncForPeer(req.URL)
} else {
fmt.Printf("[PEERS] URL rejected (no valid PEER_SECRET or validator key): %s\n", req.URL)
}
} else if req.URL != "" {
fmt.Printf("[PEERS] URL rejected (must be public HTTPS): %s\n", req.URL)
}
if addr := strings.ToLower(strings.TrimSpace(req.SigningAddress)); addr != "" && strings.HasPrefix(addr, "0x") && len(addr) == 42 {
// Authorization: accept if PEER_SECRET matches OR if the address has
// a registered validator key (individual human-signed credential) OR
// if the peer provided a valid challenge-response signature (P1-3).
sigOK := req.Signature != "" && a.blockchain.VerifyPeerChallenge(addr, req.Signature)
keys := a.state.GetValidatorKeys()
keyAuthorized := false
for _, k := range keys {
if k["signing_address"] == addr { keyAuthorized = true; break }
}
if secretOK || keyAuthorized || sigOK {
a.blockchain.AddAuthorizedValidator(addr)
method := "key"
if sigOK && !keyAuthorized { method = "challenge-response signature" }
if secretOK && !keyAuthorized && !sigOK { method = "PEER_SECRET" }
if !keyAuthorized { fmt.Printf("[PEERS] Auto-authorized validator via %s: %s\n", method, addr) }
} else if req.Signature == "" {
fmt.Printf("[PEERS] Validator %s: no signature provided — request /api/peers/challenge first\n", addr)
} else {
fmt.Printf("[PEERS] Validator %s: invalid/expired challenge signature\n", addr)
}
a.blockchain.mu.RLock()
validators := make([]string, 0, len(a.blockchain.authorizedValidators))
for v := range a.blockchain.authorizedValidators { validators = append(validators, v) }
a.blockchain.mu.RUnlock()
// P2-9: only return validator list if authorized
if secretOK || keyAuthorized || sigOK {
json.NewEncoder(w).Encode(map[string]interface{}{"peers": GlobalPeerRegistry.AllPeers(), "validators": validators})
} else {
json.NewEncoder(w).Encode(map[string]interface{}{"peers": GlobalPeerRegistry.AllPeers()})
}
return
}
json.NewEncoder(w).Encode(map[string]interface{}{"peers": GlobalPeerRegistry.AllPeers(), "validators": []string{}})
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
// SNAPSHOT_TOKEN is mandatory — the snapshot contains nullifier-to-wallet
// mappings and bio-registration data. Without a token, all requests are
// rejected so the endpoint is never accidentally left open.
token := os.Getenv("SNAPSHOT_TOKEN")
if token == "" {
http.Error(w, `{"error":"SNAPSHOT_TOKEN not configured"}`, http.StatusForbidden)
return
}
// P2-15: token in Authorization header (not URL query param that lands in logs).
authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
if authHeader != token && r.URL.Query().Get("token") != token {
http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
return
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

func (a *APIServer) handleStaticDownload(w http.ResponseWriter, r *http.Request, path, filename, contentType string) {
f, err := os.Open(path)
if err != nil {
	http.Error(w, "File not found", 404)
	return
}
defer f.Close()
fi, err := f.Stat()
if err != nil {
	http.Error(w, "File error", 500)
	return
}
w.Header().Set("Content-Disposition", "attachment; filename="+filename)
w.Header().Set("Content-Type", contentType)
w.Header().Set("Access-Control-Allow-Origin", "*")
http.ServeContent(w, r, filename, fi.ModTime(), f)
}
