package keeper

import (
"crypto/sha256"
"database/sql"
"encoding/hex"
"encoding/json"
"fmt"
"math"
"os"
"sync"
	"strings"

_ "github.com/lib/pq"
)

type AccountState struct {
Address     string  `json:"address"`
Balance     float64 `json:"balance"`
IsHuman     bool    `json:"is_human"`
// TUsdBalance is the account's holding of tUSD — a simulated, chain-native
// test-dollar token used to exercise the swap/liquidity-pool mechanism
// without touching any real external currency or bridge. See PoolState
// below for the actual AEQ<->tUSD liquidity pool this balance interacts
// with.
TUsdBalance float64 `json:"tusd_balance"`
// LPShares is this account's claim on the liquidity pool, in the same
// units as PoolState.TotalLPShares. An account's withdrawable amount at
// any moment is (LPShares / TotalLPShares) * each reserve — see
// RemoveLiquidity. This is the standard Uniswap v2 share-accounting
// model: shares are minted on deposit and burned on withdrawal, so each
// LP's claim automatically reflects fees/price-impact accumulated by the
// pool since they joined, without needing per-LP bookkeeping of "their"
// specific tokens.
LPShares float64 `json:"lp_shares"`
}

// PoolState holds the two reserves of the single AEQ<->tUSD liquidity pool.
// Pricing follows the constant-product formula (reserveAEQ * reserveTUSD =
// k), the same model Uniswap v2 popularized: the more of one side someone
// swaps in, the worse the price gets for the next unit, which is what
// makes the pool self-balancing without needing an oracle or admin to set
// a price. A 0.1% fee is taken from every swap's input amount before the
// constant-product math runs, and is distributed across the four pools
// from the original tokenomics design (validators/LPs/UBI/treasury) —
// see DistributeSwapFee. Ordinary AEQ-to-AEQ transfers (state.Transfer)
// are NOT touched by this fee; it only applies to swaps through this pool.
type PoolState struct {
ReserveAEQ    float64 `json:"reserve_aeq"`
ReserveTUSD   float64 `json:"reserve_tusd"`
// TotalLPShares is the sum of every account's LPShares. Starts at 0; the
// very first deposit mints sqrt(amountAEQ * amountTUSD) shares (the
// standard Uniswap v2 formula — using the geometric mean means the
// first depositor's chosen ratio doesn't let them mint an arbitrarily
// large or small initial share count by gaming the two amounts).
TotalLPShares float64 `json:"total_lp_shares"`
}

type ChainState struct {
mu       sync.RWMutex
accounts map[string]*AccountState
pool     *PoolState
db       *sql.DB
useDB    bool
}

func NewChainState(dataFile string) *ChainState {
cs := &ChainState{
accounts: make(map[string]*AccountState),
}

// Try PostgreSQL first
dbURL := os.Getenv("DATABASE_URL")
if dbURL != "" {
// Add sslmode if not present
if !strings.Contains(dbURL, "sslmode") {
if strings.Contains(dbURL, "?") {
dbURL += "&sslmode=disable"
} else {
dbURL += "?sslmode=disable"
}
}
db, err := sql.Open("postgres", dbURL)
if err == nil {
err = db.Ping()
if err == nil {
cs.db = db
cs.useDB = true
cs.initDB()
cs.loadFromDB()
fmt.Println("✓ ChainState using PostgreSQL")
return cs
}
}
fmt.Printf("⚠ PostgreSQL failed: %v - using file\n", err)
}

// Fallback to file
cs.useDB = false
if os.Getenv("RESET_STATE") == "true" {
fmt.Println("✓ RESET_STATE=true — starting fresh")
os.Remove(dataFile)
} else {
cs.loadFromFile(dataFile)
}
return cs
}

func (cs *ChainState) initDB() {
cs.db.Exec(`CREATE TABLE IF NOT EXISTS evm_contracts (
address TEXT PRIMARY KEY,
bytecode TEXT NOT NULL,
deployer TEXT,
deployed_at TIMESTAMP DEFAULT NOW()
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS evm_storage (
address TEXT NOT NULL,
slot TEXT NOT NULL,
value TEXT NOT NULL,
PRIMARY KEY (address, slot)
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS evm_nonces (
address TEXT PRIMARY KEY,
nonce BIGINT DEFAULT 0
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS chain_accounts (
address TEXT PRIMARY KEY,
balance FLOAT NOT NULL DEFAULT 0,
is_human BOOLEAN NOT NULL DEFAULT false
)`)
// tusd_balance added separately (ALTER instead of being in the original
// CREATE TABLE) so this upgrade doesn't require recreating the table on
// chains that already have chain_accounts from before this feature.
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS tusd_balance FLOAT NOT NULL DEFAULT 0`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS lp_shares FLOAT NOT NULL DEFAULT 0`)
// Links a ZK proof commitment to the wallet that successfully registered
// with it, so the app can ask "did MY proof get registered, and to which
// wallet?" instead of guessing from a global, unfiltered list.
cs.db.Exec(`CREATE TABLE IF NOT EXISTS bio_registrations (
commitment TEXT PRIMARY KEY,
wallet_address TEXT NOT NULL,
tx_hash TEXT,
registered_at TIMESTAMP DEFAULT NOW()
)`)
// Single-row table holding the AEQ<->tUSD pool reserves. A fixed id=1 row
// is used instead of a key-value table since there's only ever one pool
// right now — simpler queries, and trivial to extend to multiple pools
// later (id column is already there) if more pairs are ever added.
cs.db.Exec(`CREATE TABLE IF NOT EXISTS liquidity_pool (
id INTEGER PRIMARY KEY DEFAULT 1,
reserve_aeq FLOAT NOT NULL DEFAULT 0,
reserve_tusd FLOAT NOT NULL DEFAULT 0,
total_lp_shares FLOAT NOT NULL DEFAULT 0
)`)
cs.db.Exec(`ALTER TABLE liquidity_pool ADD COLUMN IF NOT EXISTS total_lp_shares FLOAT NOT NULL DEFAULT 0`)
}

func (cs *ChainState) loadFromDB() {
rows, err := cs.db.Query("SELECT address, balance, is_human, tusd_balance, lp_shares FROM chain_accounts")
if err != nil {
fmt.Printf("⚠ Could not load from DB: %v\n", err)
return
}
defer rows.Close()
count := 0
for rows.Next() {
acc := &AccountState{}
rows.Scan(&acc.Address, &acc.Balance, &acc.IsHuman, &acc.TUsdBalance, &acc.LPShares)
cs.accounts[acc.Address] = acc
count++
}
fmt.Printf("✓ Loaded %d accounts from PostgreSQL\n", count)

cs.loadOrInitPool()
}

// loadOrInitPool reads the single liquidity_pool row, creating it (at
// 0/0/0) if it doesn't exist yet. The pool intentionally does NOT get
// auto-filled with any starting reserves: every AEQ in this system only
// ever exists because a real human registered for it ("money exists
// because people exist"), so a pool can't be seeded out of thin air
// without breaking that principle. Real liquidity has to come from
// someone actually depositing AEQ they earned via AddLiquidity below.
func (cs *ChainState) loadOrInitPool() {
var reserveAEQ, reserveTUSD, totalShares float64
err := cs.db.QueryRow("SELECT reserve_aeq, reserve_tusd, total_lp_shares FROM liquidity_pool WHERE id = 1").Scan(&reserveAEQ, &reserveTUSD, &totalShares)
if err != nil {
_, insertErr := cs.db.Exec(`INSERT INTO liquidity_pool (id, reserve_aeq, reserve_tusd, total_lp_shares) VALUES (1, 0, 0, 0)
ON CONFLICT (id) DO NOTHING`)
if insertErr != nil {
fmt.Printf("⚠ Could not create liquidity pool row: %v\n", insertErr)
}
cs.pool = &PoolState{ReserveAEQ: 0, ReserveTUSD: 0, TotalLPShares: 0}
fmt.Printf("✓ Liquidity pool created (empty — awaiting first deposit via AddLiquidity)\n")
return
}
cs.pool = &PoolState{ReserveAEQ: reserveAEQ, ReserveTUSD: reserveTUSD, TotalLPShares: totalShares}
fmt.Printf("✓ Liquidity pool loaded: %.2f AEQ / %.2f tUSD / %.6f shares\n", reserveAEQ, reserveTUSD, totalShares)
}

func (cs *ChainState) loadFromFile(dataFile string) {
data, err := os.ReadFile(dataFile)
if err != nil {
fmt.Println("✓ Starting with fresh chain state")
return
}
var accounts map[string]*AccountState
if err := json.Unmarshal(data, &accounts); err != nil {
fmt.Println("⚠ Could not load state, starting fresh")
return
}
cs.accounts = accounts
fmt.Printf("✓ Loaded chain state: %d accounts\n", len(accounts))
}

func (cs *ChainState) save() {
if cs.useDB {
return // DB saves immediately in RegisterHuman/Transfer
}
data, _ := json.Marshal(cs.accounts)
os.WriteFile("/tmp/aequitas_state.json", data, 0644)
}

func (cs *ChainState) saveAccountToDB(acc *AccountState) {
if !cs.useDB {
return
}
_, err := cs.db.Exec(`INSERT INTO chain_accounts (address, balance, is_human, tusd_balance, lp_shares) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (address) DO UPDATE SET balance = $2, is_human = $3, tusd_balance = $4, lp_shares = $5`,
acc.Address, acc.Balance, acc.IsHuman, acc.TUsdBalance, acc.LPShares)
if err != nil {
fmt.Printf("[DB] Error saving account %s: %v\n", acc.Address, err)
} else {
fmt.Printf("[DB] Saved account %s | balance=%.2f | tusd=%.2f | lp=%.6f | is_human=%v\n", acc.Address, acc.Balance, acc.TUsdBalance, acc.LPShares, acc.IsHuman)
}
}

func (cs *ChainState) GetBalance(address string) float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
if acc, ok := cs.accounts[address]; ok {
return acc.Balance
}
return 0
}

func (cs *ChainState) GetTUsdBalance(address string) float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
if acc, ok := cs.accounts[address]; ok {
return acc.TUsdBalance
}
return 0
}

func (cs *ChainState) GetPoolReserves() (float64, float64) {
cs.mu.RLock()
defer cs.mu.RUnlock()
if cs.pool == nil {
return 0, 0
}
return cs.pool.ReserveAEQ, cs.pool.ReserveTUSD
}

func (cs *ChainState) IsHuman(address string) bool {
cs.mu.RLock()
defer cs.mu.RUnlock()
if acc, ok := cs.accounts[address]; ok {
return acc.IsHuman
}
return false
}

func (cs *ChainState) RegisterHuman(address string) error {
cs.mu.Lock()
defer cs.mu.Unlock()

if acc, ok := cs.accounts[address]; ok && acc.IsHuman {
return fmt.Errorf("already registered")
}

if _, ok := cs.accounts[address]; !ok {
cs.accounts[address] = &AccountState{Address: address}
}

cs.accounts[address].IsHuman = true
cs.accounts[address].Balance += 1000
cs.saveAccountToDB(cs.accounts[address])
cs.save()

fmt.Printf("[STATE] ✓ Human registered: %s | Balance: %.2f AEQ\n",
address, cs.accounts[address].Balance)
return nil
}

func (cs *ChainState) Transfer(from, to string, amount float64) error {
cs.mu.Lock()
defer cs.mu.Unlock()

fromAcc, ok := cs.accounts[from]
if !ok || fromAcc.Balance < amount {
return fmt.Errorf("insufficient balance")
}

fromAcc.Balance -= amount
cs.saveAccountToDB(fromAcc)

if _, ok := cs.accounts[to]; !ok {
cs.accounts[to] = &AccountState{Address: to}
}
cs.accounts[to].Balance += amount
cs.saveAccountToDB(cs.accounts[to])
cs.save()

fmt.Printf("[STATE] ✓ Transfer %.2f AEQ: %s → %s\n", amount, from, to)
return nil
}

// swapFeeBps is the fee taken from every swap's input amount, in basis
// points (10 = 0.1%). This ONLY applies to swaps through the AEQ<->tUSD
// pool — ordinary AEQ-to-AEQ transfers via Transfer() above remain
// completely free, per the project's design decision that moving AEQ
// between people should never cost anything; only exchanging it for a
// different currency does.
// Fee recipient addresses for the four tokenomics pools, per the original
// design (40% validators / 30% LPs / 20% UBI / 10% treasury). These are
// real wallet addresses Daniel controls — provided explicitly so swap
// fees are credited somewhere actually accessible, rather than to
// addresses with no corresponding private key.
const (
	validatorsPoolAddr = "0x78c1c143e395b181f13bcb6868ff53aa86c3d2ba"
	lpPoolAddr         = "0xc181c3a4d09444b99089ae0f56c1e7f4c20d01eb"
	ubiPoolAddr        = "0x4a9b8f99f0d8cff0e510fef502100571203b054a"
	treasuryPoolAddr   = "0x2273894fb781978d54e767f9fba2dcb33d93eb15"
)

// swapFeeBps is the fee taken from every swap's input amount, in basis
// points (10 = 0.1%). This ONLY applies to swaps through the AEQ<->tUSD
// pool — ordinary AEQ-to-AEQ transfers via Transfer() above remain
// completely free, per the project's design decision that moving AEQ
// between people should never cost anything; only exchanging it for a
// different currency does.
const swapFeeBps = 10

// SwapAEQForTUSD swaps `amountIn` AEQ from `address` into tUSD, using the
// constant-product formula (reserveAEQ * reserveTUSD = k) for pricing. A
// 0.1% fee is deducted from amountIn before the swap math runs, and is
// distributed across the four tokenomics pools (see DistributeSwapFee)
// rather than added to the liquidity pool's reserves — so the pool's own
// k grows only from genesis seeding, not from accumulated fees, keeping
// the fee-distribution logic in one place instead of split between the
// pool and the four-way split.
func (cs *ChainState) SwapAEQForTUSD(address string, amountIn float64) (float64, error) {
cs.mu.Lock()
defer cs.mu.Unlock()
return cs.swapLocked(address, amountIn, true)
}

// SwapTUSDForAEQ swaps `amountIn` tUSD from `address` into AEQ. Same
// constant-product pricing and fee handling as SwapAEQForTUSD, just with
// the two reserves' roles reversed.
func (cs *ChainState) SwapTUSDForAEQ(address string, amountIn float64) (float64, error) {
cs.mu.Lock()
defer cs.mu.Unlock()
return cs.swapLocked(address, amountIn, false)
}

// swapLocked implements both swap directions. aeqToTusd=true means AEQ is
// the input side and tUSD is the output side; false is the reverse.
// Caller must hold cs.mu.
func (cs *ChainState) swapLocked(address string, amountIn float64, aeqToTusd bool) (float64, error) {
if amountIn <= 0 {
return 0, fmt.Errorf("amount must be positive")
}
if cs.pool == nil {
return 0, fmt.Errorf("liquidity pool not initialized")
}

acc, ok := cs.accounts[address]
if !ok {
return 0, fmt.Errorf("account not found")
}

if aeqToTusd {
if acc.Balance < amountIn {
return 0, fmt.Errorf("insufficient AEQ balance")
}
} else {
if acc.TUsdBalance < amountIn {
return 0, fmt.Errorf("insufficient tUSD balance")
}
}

// Fee is taken off the top of the input amount; only the remainder
// participates in the constant-product swap.
fee := amountIn * float64(swapFeeBps) / 10000.0
amountInAfterFee := amountIn - fee

var amountOut float64
if aeqToTusd {
// x*y=k: reserveAEQ * reserveTUSD = (reserveAEQ + amountInAfterFee) * (reserveTUSD - amountOut)
amountOut = (cs.pool.ReserveTUSD * amountInAfterFee) / (cs.pool.ReserveAEQ + amountInAfterFee)
if amountOut >= cs.pool.ReserveTUSD {
return 0, fmt.Errorf("swap too large for pool liquidity")
}
cs.pool.ReserveAEQ += amountInAfterFee
cs.pool.ReserveTUSD -= amountOut
acc.Balance -= amountIn
acc.TUsdBalance += amountOut
} else {
amountOut = (cs.pool.ReserveAEQ * amountInAfterFee) / (cs.pool.ReserveTUSD + amountInAfterFee)
if amountOut >= cs.pool.ReserveAEQ {
return 0, fmt.Errorf("swap too large for pool liquidity")
}
cs.pool.ReserveTUSD += amountInAfterFee
cs.pool.ReserveAEQ -= amountOut
acc.TUsdBalance -= amountIn
acc.Balance += amountOut
}

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.distributeSwapFee(fee, aeqToTusd)
cs.save()

fmt.Printf("[SWAP] %s: %.4f %s → %.4f %s (fee %.4f)\n",
address, amountIn, sideLabel(aeqToTusd, true), amountOut, sideLabel(aeqToTusd, false), fee)

return amountOut, nil
}

func sideLabel(aeqToTusd, isInput bool) string {
if aeqToTusd == isInput {
return "AEQ"
}
return "tUSD"
}

func (cs *ChainState) savePoolToDB() {
if !cs.useDB || cs.pool == nil {
return
}
_, err := cs.db.Exec(`UPDATE liquidity_pool SET reserve_aeq = $1, reserve_tusd = $2, total_lp_shares = $3 WHERE id = 1`,
cs.pool.ReserveAEQ, cs.pool.ReserveTUSD, cs.pool.TotalLPShares)
if err != nil {
fmt.Printf("[DB] Error saving pool: %v\n", err)
}
}

// distributeSwapFee splits the fee collected from a swap across the four
// tokenomics pools from the original design: 40% validators, 30%
// liquidity providers, 20% UBI, 10% treasury — crediting each of the four
// real addresses above. feeInAEQ is true when the fee was collected in
// AEQ (an AEQ->tUSD swap); false means it was collected in tUSD (a
// tUSD->AEQ swap) — the split percentages are the same either way, only
// the currency the fee is credited in differs. Caller must hold cs.mu.
func (cs *ChainState) distributeSwapFee(fee float64, feeInAEQ bool) {
if fee <= 0 {
return
}
shares := map[string]float64{
validatorsPoolAddr: fee * 0.40,
lpPoolAddr:         fee * 0.30,
ubiPoolAddr:        fee * 0.20,
treasuryPoolAddr:   fee * 0.10,
}
for addr, amount := range shares {
if _, ok := cs.accounts[addr]; !ok {
cs.accounts[addr] = &AccountState{Address: addr}
}
if feeInAEQ {
cs.accounts[addr].Balance += amount
} else {
cs.accounts[addr].TUsdBalance += amount
}
cs.saveAccountToDB(cs.accounts[addr])
}
currency := "tUSD"
if feeInAEQ {
currency = "AEQ"
}
fmt.Printf("[FEE] Swap fee %.6f %s distributed across validators/lps/ubi/treasury\n", fee, currency)
}

// AddLiquidity lets a real account deposit AEQ and tUSD into the pool in
// proportion to the pool's current ratio (or, if the pool is currently
// empty, at whatever ratio the depositor chooses — that first deposit
// sets the initial price). This is the ONLY way reserves enter the pool;
// there is no admin/genesis fill, since every AEQ here has to trace back
// to a real human's registration grant, consistent with "money exists
// because people exist."
//
// NOTE: this does not yet mint or track LP shares/tokens — it only moves
// balances into the pool. A depositor currently has no on-chain claim to
// withdraw their share back out. Tracking proportional LP ownership (so
// deposits are genuinely reversible) is a deliberate follow-up, not
// included in this first pass.
func (cs *ChainState) AddLiquidity(address string, amountAEQ, amountTUSD float64) error {
cs.mu.Lock()
defer cs.mu.Unlock()

if amountAEQ <= 0 || amountTUSD <= 0 {
return fmt.Errorf("both amounts must be positive")
}
if cs.pool == nil {
return fmt.Errorf("liquidity pool not initialized")
}

acc, ok := cs.accounts[address]
if !ok {
return fmt.Errorf("account not found")
}
if acc.Balance < amountAEQ {
return fmt.Errorf("insufficient AEQ balance")
}
if acc.TUsdBalance < amountTUSD {
return fmt.Errorf("insufficient tUSD balance")
}

// If the pool already has liquidity, require the deposit to roughly
// match the existing ratio — an unbalanced deposit would otherwise
// instantly shift the price, which is the same rule real AMMs enforce.
var mintedShares float64
if cs.pool.ReserveAEQ > 0 && cs.pool.ReserveTUSD > 0 {
expectedTUSD := amountAEQ * (cs.pool.ReserveTUSD / cs.pool.ReserveAEQ)
tolerance := expectedTUSD * 0.01 // 1% slack for rounding
if amountTUSD < expectedTUSD-tolerance || amountTUSD > expectedTUSD+tolerance {
return fmt.Errorf("deposit ratio does not match pool ratio (expected ~%.4f tUSD for %.4f AEQ)", expectedTUSD, amountAEQ)
}
if cs.pool.TotalLPShares > 0 {
// Proportional to the pool's existing size — same fraction of the
// AEQ reserve as the fraction of total shares being minted, so an
// LP's claim accurately tracks how much of the pool they actually
// own (including any fees the pool has accumulated since genesis).
mintedShares = (amountAEQ / cs.pool.ReserveAEQ) * cs.pool.TotalLPShares
} else {
// The pool has real reserves but zero recorded shares — this
// happens for reserves that were deposited before LP-share
// tracking existed (no account was ever credited shares for
// them). Treat this deposit as if it were bootstrapping a pool
// that already happens to have this much in it: mint shares for
// the NEW deposit using the geometric-mean formula, AND retroactively
// mint the depositor shares for the pool's pre-existing, currently
// unclaimed reserves too — otherwise those original reserves would
// permanently sit in the pool with no one able to ever withdraw
// them, since shares would stay stuck at zero forever (any future
// deposit would hit this same zero-total-shares branch again).
newShares := math.Sqrt(amountAEQ * amountTUSD)
preExistingShares := math.Sqrt(cs.pool.ReserveAEQ * cs.pool.ReserveTUSD)
mintedShares = newShares + preExistingShares
fmt.Printf("[POOL] ⚠ Pool had %.4f AEQ / %.4f tUSD in reserves with zero recorded LP shares (pre-dates share tracking) — crediting depositor %.6f shares for those alongside %.6f shares for this new deposit\n",
cs.pool.ReserveAEQ, cs.pool.ReserveTUSD, preExistingShares, newShares)
}
} else {
// First-ever deposit: shares = geometric mean of the two amounts
// (standard Uniswap v2 bootstrap formula). Using sqrt(x*y) instead
// of, say, just amountAEQ means the first depositor can't mint an
// outsized number of shares simply by picking a lopsided ratio.
mintedShares = math.Sqrt(amountAEQ * amountTUSD)
}

acc.Balance -= amountAEQ
acc.TUsdBalance -= amountTUSD
acc.LPShares += mintedShares
cs.pool.ReserveAEQ += amountAEQ
cs.pool.ReserveTUSD += amountTUSD
cs.pool.TotalLPShares += mintedShares

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.save()

fmt.Printf("[POOL] ✓ %s added liquidity: %.4f AEQ + %.4f tUSD → %.6f LP shares\n", address, amountAEQ, amountTUSD, mintedShares)
return nil
}

// RemoveLiquidity burns sharesToBurn of address's LP shares and returns
// the corresponding proportional amount of both reserves to their
// balances. sharesToBurn must not exceed the account's own LPShares —
// an account can only withdraw its own claim, never another LP's.
func (cs *ChainState) RemoveLiquidity(address string, sharesToBurn float64) (float64, float64, error) {
cs.mu.Lock()
defer cs.mu.Unlock()

if sharesToBurn <= 0 {
return 0, 0, fmt.Errorf("shares must be positive")
}
if cs.pool == nil || cs.pool.TotalLPShares <= 0 {
return 0, 0, fmt.Errorf("liquidity pool is empty")
}

acc, ok := cs.accounts[address]
if !ok {
return 0, 0, fmt.Errorf("account not found")
}
if acc.LPShares < sharesToBurn {
return 0, 0, fmt.Errorf("insufficient LP shares (have %.6f, requested %.6f)", acc.LPShares, sharesToBurn)
}

fraction := sharesToBurn / cs.pool.TotalLPShares
outAEQ := cs.pool.ReserveAEQ * fraction
outTUSD := cs.pool.ReserveTUSD * fraction

acc.LPShares -= sharesToBurn
acc.Balance += outAEQ
acc.TUsdBalance += outTUSD
cs.pool.ReserveAEQ -= outAEQ
cs.pool.ReserveTUSD -= outTUSD
cs.pool.TotalLPShares -= sharesToBurn

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.save()

fmt.Printf("[POOL] ✓ %s removed liquidity: %.6f shares → %.4f AEQ + %.4f tUSD\n", address, sharesToBurn, outAEQ, outTUSD)
return outAEQ, outTUSD, nil
}

// GetLPShares returns address's current LP share balance, and the pool's
// total shares — callers can compute the account's ownership fraction
// (and therefore its withdrawable amounts) from these two numbers.
func (cs *ChainState) GetLPShares(address string) (float64, float64) {
cs.mu.RLock()
defer cs.mu.RUnlock()
var mine float64
if acc, ok := cs.accounts[address]; ok {
mine = acc.LPShares
}
total := 0.0
if cs.pool != nil {
total = cs.pool.TotalLPShares
}
return mine, total
}

func (cs *ChainState) TotalSupply() float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
total := 0.0
for _, acc := range cs.accounts {
total += acc.Balance
}
return total
}

func (cs *ChainState) TotalHumans() int {
cs.mu.RLock()
defer cs.mu.RUnlock()
count := 0
for _, acc := range cs.accounts {
if acc.IsHuman {
count++
}
}
return count
}

func (cs *ChainState) GetAllAccounts() []*AccountState {
cs.mu.RLock()
defer cs.mu.RUnlock()
result := make([]*AccountState, 0, len(cs.accounts))
for _, acc := range cs.accounts {
result = append(result, acc)
}
return result
}

// tusdFaucetAmount is how much test-tUSD ClaimTUsdFaucet grants per
// account, once. tUSD is a simulated currency with no real-world value —
// unlike AEQ (which only ever exists because a real human registered for
// it), there's no "money exists because people exist" principle being
// violated by handing test-tUSD out directly. This exists purely so a
// registered human has something to pair with their real AEQ the first
// time they call AddLiquidity, since otherwise nobody could ever provide
// the tUSD side of the very first deposit.
const tusdFaucetAmount = 1000.0

// ClaimTUsdFaucet grants tusdFaucetAmount of test-tUSD to address, once.
// Returns an error if the account isn't registered, or already claimed.
func (cs *ChainState) ClaimTUsdFaucet(address string) error {
cs.mu.Lock()
defer cs.mu.Unlock()

acc, ok := cs.accounts[address]
if !ok || !acc.IsHuman {
return fmt.Errorf("only registered humans can claim the test-tUSD faucet")
}
if acc.TUsdBalance > 0 {
return fmt.Errorf("faucet already claimed")
}
// NOTE: this check only blocks re-claiming while a balance > 0 remains.
// Spending it all via a swap or AddLiquidity would make TUsdBalance hit
// 0 again, after which this same account could claim once more. Fine
// for a first-pass test faucet; a real one-time flag (e.g. a separate
// "claimed" column) would be needed before this matters in practice.

acc.TUsdBalance = tusdFaucetAmount
cs.saveAccountToDB(acc)
cs.save()

fmt.Printf("[FAUCET] ✓ %s claimed %.2f test-tUSD\n", address, tusdFaucetAmount)
return nil
}

func (cs *ChainState) StateRoot() string {
cs.mu.RLock()
data, _ := json.Marshal(cs.accounts)
cs.mu.RUnlock()
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:])
}

func (cs *ChainState) CalcGini() float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
if len(cs.accounts) < 2 {
return 0.0
}
balances := []float64{}
for _, acc := range cs.accounts {
if acc.Balance > 0 {
balances = append(balances, acc.Balance)
}
}
n := len(balances)
if n < 2 {
return 0.0
}
// Sort ascending
for i := 0; i < n; i++ {
for j := i + 1; j < n; j++ {
if balances[j] < balances[i] {
balances[i], balances[j] = balances[j], balances[i]
}
}
}
// Gini formula
sum := 0.0
for _, b := range balances {
sum += b
}
if sum == 0 {
return 0.0
}
numerator := 0.0
for i, b := range balances {
numerator += float64(2*i-n+1) * b
}
gini := numerator / (float64(n) * sum)
if gini < 0 {
gini = -gini
}
return gini
}

func (cs *ChainState) CalcAequitasIndex() float64 {
gini := cs.CalcGini()
humans := float64(cs.TotalHumans())
// Base index from Gini (0-100)
index := gini * 100.0
// Adjust for network size (small networks have inherently low inequality)
if humans < 10 {
// Bootstrap phase - index reflects growth potential
index = index * (humans / 10.0)
}
// Round to 1 decimal
return float64(int(index*10)) / 10.0
}

func (cs *ChainState) CalcPhase() int {
humans := cs.TotalHumans()
supply := cs.TotalSupply()
gini := cs.CalcGini()
switch {
case humans >= 1000000 && gini < 0.3:
return 3
case humans >= 10000 || supply >= 10000000:
return 2
case humans >= 100:
return 1
default:
return 0
}
}

func (cs *ChainState) SetBalance(address string, amount float64) {
cs.mu.Lock()
defer cs.mu.Unlock()
address = strings.ToLower(address)
if acc, ok := cs.accounts[address]; ok {
acc.Balance = amount
cs.saveAccountToDB(acc)
} else {
acc = &AccountState{Address: address, Balance: amount}
cs.accounts[address] = acc
cs.saveAccountToDB(acc)
}
}


// V6 Contract State Mirror - persists EVM contract state to PostgreSQL
func (cs *ChainState) InitV6StateTables() {
if cs.db == nil {
return
}
cs.db.Exec(`CREATE TABLE IF NOT EXISTS v6_state (
key TEXT PRIMARY KEY,
value TEXT NOT NULL,
updated_at TIMESTAMP DEFAULT NOW()
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS v6_humans (
address TEXT PRIMARY KEY,
commitment TEXT,
is_human BOOLEAN DEFAULT true,
is_inactive BOOLEAN DEFAULT false,
registered_at TIMESTAMP DEFAULT NOW(),
last_activity TIMESTAMP DEFAULT NOW()
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS v6_balances (
address TEXT PRIMARY KEY,
balance_wei TEXT NOT NULL,
updated_at TIMESTAMP DEFAULT NOW()
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS v6_commitments (
commitment TEXT PRIMARY KEY,
wallet TEXT NOT NULL,
used_at TIMESTAMP DEFAULT NOW()
)`)
fmt.Println("[V6] State tables initialized")
}

func (cs *ChainState) SaveV6State(key, value string) {
if cs.db == nil {
return
}
cs.db.Exec(
`INSERT INTO v6_state (key, value) VALUES ($1, $2)
 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`,
key, value,
)
}

func (cs *ChainState) LoadV6State(key string) string {
if cs.db == nil {
return ""
}
var value string
cs.db.QueryRow(`SELECT value FROM v6_state WHERE key = $1`, key).Scan(&value)
return value
}

func (cs *ChainState) SaveV6Balance(address, balanceWei string) {
if cs.db == nil {
return
}
cs.db.Exec(
`INSERT INTO v6_balances (address, balance_wei) VALUES ($1, $2)
 ON CONFLICT (address) DO UPDATE SET balance_wei = $2, updated_at = NOW()`,
address, balanceWei,
)
}

func (cs *ChainState) LoadV6Balance(address string) string {
if cs.db == nil {
return "0"
}
var balanceWei string
cs.db.QueryRow(`SELECT balance_wei FROM v6_balances WHERE address = $1`, address).Scan(&balanceWei)
if balanceWei == "" {
return "0"
}
return balanceWei
}

func (cs *ChainState) SaveV6Human(address, commitment string) {
if cs.db == nil {
return
}
cs.db.Exec(
`INSERT INTO v6_humans (address, commitment) VALUES ($1, $2)
 ON CONFLICT (address) DO UPDATE SET commitment = $2, updated_at = NOW()`,
address, commitment,
)
}

func (cs *ChainState) SaveV6Commitment(commitment, wallet string) {
if cs.db == nil {
return
}
cs.db.Exec(
`INSERT INTO v6_commitments (commitment, wallet) VALUES ($1, $2)
 ON CONFLICT (commitment) DO NOTHING`,
commitment, wallet,
)
}

func (cs *ChainState) GetAllV6Humans() []map[string]string {
if cs.db == nil {
return nil
}
rows, err := cs.db.Query(
`SELECT address, commitment FROM v6_humans WHERE is_human = true AND is_inactive = false`,
)
if err != nil {
return nil
}
defer rows.Close()
var humans []map[string]string
for rows.Next() {
var addr, commitment string
rows.Scan(&addr, &commitment)
humans = append(humans, map[string]string{
"address":    addr,
"commitment": commitment,
})
}
return humans
}

func (cs *ChainState) GetAllV6Balances() []map[string]string {
if cs.db == nil {
return nil
}
rows, err := cs.db.Query(`SELECT address, balance_wei FROM v6_balances`)
if err != nil {
return nil
}
defer rows.Close()
var balances []map[string]string
for rows.Next() {
var addr, bal string
rows.Scan(&addr, &bal)
balances = append(balances, map[string]string{
"address":    addr,
"balance_wei": bal,
})
}
return balances
}
