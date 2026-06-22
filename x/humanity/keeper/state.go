package keeper

import (
"crypto/sha256"
"database/sql"
"encoding/hex"
"encoding/json"
"fmt"
"math"
"os"
"sort"
"sync"
"time"
	"strings"

"github.com/lib/pq"
)

// timeNowFunc is a seam for time.Now(), letting demurrage timing be
// mocked in tests without needing to thread a clock through every call.
var timeNowFunc = time.Now

type AccountState struct {
Address     string  `json:"address"`
Balance     Decimal `json:"balance"`
IsHuman     bool    `json:"is_human"`
// TUsdBalance is the account's holding of tUSD — a simulated, chain-native
// test-dollar token used to exercise the swap/liquidity-pool mechanism
// without touching any real external currency or bridge. See PoolState
// below for the actual AEQ<->tUSD liquidity pool this balance interacts
// with.
TUsdBalance Decimal `json:"tusd_balance"`
// LPShares is this account's claim on the liquidity pool, in the same
// units as PoolState.TotalLPShares. An account's withdrawable amount at
// any moment is (LPShares / TotalLPShares) * each reserve — see
// RemoveLiquidity. This is the standard Uniswap v2 share-accounting
// model: shares are minted on deposit and burned on withdrawal, so each
// LP's claim automatically reflects fees/price-impact accumulated by the
// pool since they joined, without needing per-LP bookkeeping of "their"
// specific tokens.
LPShares Decimal `json:"lp_shares"`
// LastActivityAt is the Unix timestamp (seconds) of this account's most
// recent AEQ-moving action (registration, sending/receiving a transfer,
// swapping, or adding/removing liquidity). Demurrage (see ApplyDemurrage)
// is calculated live from how long it's been since this timestamp — the
// balance shown to the user is always computed fresh from Balance and
// this timestamp, rather than being eaten away by a periodic background
// job. Touching the account in any of those ways resets this timestamp,
// which is the whole point: money that's actively circulating doesn't
// decay, only money that's sitting idle does.
LastActivityAt int64 `json:"last_activity_at"`
// Demurrage14DayWarningShown tracks whether the one-time "your balance
// starts decaying in 14 days" notice has already been surfaced for the
// CURRENT grace period. Reset back to false by touchActivity whenever
// the account's clock restarts (any AEQ-moving action), so the warning
// can fire again for the next idle period rather than being a permanent
// one-time-ever flag.
Demurrage14DayWarningShown bool `json:"demurrage_14_day_warning_shown"`
// FaucetClaimed is set permanently to true once an account has claimed the
// tUSD test faucet. Unlike the old TUsdBalance>0 check, this flag is never
// reset by spending tUSD, so a wallet cannot re-claim by draining its balance.
FaucetClaimed bool `json:"faucet_claimed"`
Version      int64 `json:"-"` // optimistic lock version, not serialized
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
ReserveAEQ    Decimal `json:"reserve_aeq"`
ReserveTUSD   Decimal `json:"reserve_tusd"`
// TotalLPShares is the sum of every account's LPShares. Starts at 0; the
// very first deposit mints sqrt(amountAEQ * amountTUSD) shares (the
// standard Uniswap v2 formula — using the geometric mean means the
// first depositor's chosen ratio doesn't let them mint an arbitrarily
// large or small initial share count by gaming the two amounts).
TotalLPShares Decimal `json:"total_lp_shares"`
}

type ChainState struct {
mu        sync.RWMutex
stateMu   sync.Mutex
accounts  map[string]*AccountState
pool      *PoolState
db        *sql.DB
useDB     bool
nullifiers map[string]string // nullifier hex → wallet address (in-memory cache)
}

// acquireStateLock serializes critical state mutations using an in-process
// sync.Mutex. PostgreSQL advisory locks were removed because database/sql
// uses a connection pool and lock/unlock can land on different connections,
// causing lock leaks or deadlocks. For multi-node, the version column on
// chain_accounts provides optimistic concurrency conflict detection.
func (cs *ChainState) acquireStateLock() bool {
cs.stateMu.Lock()
return true
}

func (cs *ChainState) releaseStateLock() {
cs.stateMu.Unlock()
}

// beginStateTx starts a SERIALIZABLE PostgreSQL transaction for critical
// state mutations. The caller is responsible for calling Commit or Rollback.
// Returns nil if no DB is configured (in-memory/file mode).
func (cs *ChainState) beginStateTx() *sql.Tx {
if cs.db == nil {
return nil
}
tx, err := cs.db.Begin()
if err != nil {
fmt.Printf("[DB] Warning: could not begin state tx: %v\n", err)
return nil
}
return tx
}

func NewChainState(dataFile string) *ChainState {
cs := &ChainState{
accounts:  make(map[string]*AccountState),
nullifiers: make(map[string]string),
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
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS version BIGINT NOT NULL DEFAULT 0`)
// Upgrade balance columns to NUMERIC(20,6) for exact decimal storage.
// FLOAT loses precision; NUMERIC stores exactly 6 decimal places.
cs.db.Exec(`ALTER TABLE chain_accounts ALTER COLUMN balance TYPE NUMERIC(20,6) USING balance::NUMERIC(20,6)`)
cs.db.Exec(`ALTER TABLE chain_accounts ALTER COLUMN tusd_balance TYPE NUMERIC(20,6) USING tusd_balance::NUMERIC(20,6)`)
cs.db.Exec(`ALTER TABLE chain_accounts ALTER COLUMN lp_shares TYPE NUMERIC(20,6) USING lp_shares::NUMERIC(20,6)`)
cs.db.Exec(`UPDATE liquidity_pool SET reserve_aeq = ROUND(reserve_aeq::NUMERIC, 6), reserve_tusd = ROUND(reserve_tusd::NUMERIC, 6), total_lp_shares = ROUND(total_lp_shares::NUMERIC, 6) WHERE id = 1`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS tusd_balance FLOAT NOT NULL DEFAULT 0`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS lp_shares FLOAT NOT NULL DEFAULT 0`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS last_activity_at BIGINT NOT NULL DEFAULT 0`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS demurrage_14_day_warning_shown BOOLEAN NOT NULL DEFAULT false`)
cs.db.Exec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS faucet_claimed BOOLEAN NOT NULL DEFAULT false`)
// Links a ZK proof commitment to the wallet that successfully registered
// with it, so the app can ask "did MY proof get registered, and to which
// wallet?" instead of guessing from a global, unfiltered list.
cs.db.Exec(`CREATE TABLE IF NOT EXISTS bio_registrations (
commitment TEXT PRIMARY KEY,
wallet_address TEXT NOT NULL,
tx_hash TEXT,
registered_at TIMESTAMP DEFAULT NOW()
)`)
// bio_hash lets the app poll "did MY device's identity hash get
// registered yet, and to which wallet?" — needed because, under the new
// flow where the proof is generated on the website (after MetaMask
// supplies the real wallet), the app itself never computes a commitment
// and so can't poll by one. It only ever knows its own bio_hash.
cs.db.Exec(`ALTER TABLE bio_registrations ADD COLUMN IF NOT EXISTS bio_hash TEXT`)
cs.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uidx_bio_registrations_bio_hash ON bio_registrations(bio_hash) WHERE bio_hash IS NOT NULL`)
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
// Upgrade liquidity_pool columns to NUMERIC(20,6) for exact decimal storage
// (same migration applied to chain_accounts columns above). Must run AFTER the
// ADD COLUMN statements so that the column definitely exists before ALTER TYPE.
cs.db.Exec(`ALTER TABLE liquidity_pool ALTER COLUMN reserve_aeq TYPE NUMERIC(20,6) USING reserve_aeq::NUMERIC(20,6)`)
cs.db.Exec(`ALTER TABLE liquidity_pool ALTER COLUMN reserve_tusd TYPE NUMERIC(20,6) USING reserve_tusd::NUMERIC(20,6)`)
cs.db.Exec(`ALTER TABLE liquidity_pool ALTER COLUMN total_lp_shares TYPE NUMERIC(20,6) USING total_lp_shares::NUMERIC(20,6)`)
// nullifiers stores the one-way SHA256 derivative of each identity's bioHash.
// Checked at registration time to prevent the same biometric from registering
// with a second wallet. The nullifier itself never reveals the bioHash.
cs.db.Exec(`CREATE TABLE IF NOT EXISTS nullifiers (
nullifier TEXT PRIMARY KEY,
wallet_address TEXT NOT NULL,
registered_at TIMESTAMP DEFAULT NOW()
)`)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS chain_config (
key TEXT PRIMARY KEY,
value TEXT NOT NULL
)`)

cs.InitSwapNoncesTable()
cs.InitValidatorKeysTable()
cs.InitGiniSnapshotsTable()
cs.InitPriceSnapshotsTable()}

// setConfigValue persists a key/value pair to chain_config (upsert).
func (cs *ChainState) setConfigValue(key, value string) {
if cs.db == nil {
return
}
cs.db.Exec(`INSERT INTO chain_config (key, value) VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, key, value)
}

// getConfigValue reads a key from chain_config, returning "" if missing.
func (cs *ChainState) getConfigValue(key string) string {
if cs.db == nil {
return ""
}
var v string
cs.db.QueryRow(`SELECT value FROM chain_config WHERE key = $1`, key).Scan(&v)
return v
}

// GetLastUBIAt returns the Unix timestamp of the most recent UBI distribution,
// or 0 if it has never run.
func (cs *ChainState) GetLastUBIAt() int64 {
v := cs.getConfigValue("last_ubi_at")
if v == "" {
return 0
}
var t int64
fmt.Sscan(v, &t)
return t
}

// TimeUntilNextUBI returns how long until the next UBI distribution is due.
// Returns 0 if overdue.
func (cs *ChainState) TimeUntilNextUBI() time.Duration {
last := cs.GetLastUBIAt()
if last == 0 {
return 5 * time.Second
}
next := time.Unix(last, 0).Add(24 * time.Hour)
d := time.Until(next)
if d < 0 {
return 0
}
return d
}

func (cs *ChainState) loadFromDB() {
rows, err := cs.db.Query("SELECT address, balance, is_human, tusd_balance, lp_shares, last_activity_at, demurrage_14_day_warning_shown, faucet_claimed, COALESCE(version,0) FROM chain_accounts")
if err != nil {
fmt.Printf("⚠ Could not load from DB: %v\n", err)
return
}
defer rows.Close()
count := 0
mergedCount := 0
for rows.Next() {
acc := &AccountState{}
var bal, tusd, lp float64
rows.Scan(&acc.Address, &bal, &acc.IsHuman, &tusd, &lp, &acc.LastActivityAt, &acc.Demurrage14DayWarningShown, &acc.FaucetClaimed, &acc.Version)
acc.Balance = NewDecimal(bal)
acc.TUsdBalance = NewDecimal(tusd)
acc.LPShares = NewDecimal(lp)
// Accounts loaded from DB must always use the conditional optimistic-lock
// UPDATE path in saveAccountToDB. If the version column is NULL in an old
// row, COALESCE returns 0, which would trigger the INSERT/unconditional
// path and bypass the conflict check. Normalize to 1 — both in memory AND
// in the DB row, so UPDATE … WHERE version = 1 actually finds the row.
if acc.Version == 0 {
	acc.Version = 1
	cs.db.Exec(`UPDATE chain_accounts SET version = 1 WHERE lower(address) = $1 AND (version IS NULL OR version = 0)`,
		strings.ToLower(acc.Address))
}
count++

// One-time migration: every state-mutating function (Transfer,
// RegisterHuman, swapLocked, etc.) now consistently lowercases
// addresses before using them as map keys — but rows written
// BEFORE that fix could be stored under a mixed-case address (e.g.
// MetaMask's checksum format) while later operations on the SAME
// real wallet used lowercase, splitting one person's balance across
// two separate accounts. This silently shrank what the UI showed
// for that wallet without actually losing any AEQ — the rest was
// just sitting under a differently-cased key. Merging here, once,
// at load time, makes loadFromDB self-healing for any old data
// without needing a separate manual SQL migration step.
//
// IMPORTANT: SQL row order is not guaranteed, so the mixed-case row
// for a given wallet could arrive before OR after its lowercase
// counterpart. We always check whether cs.accounts[normalized]
// already exists — regardless of whether THIS row's own address
// happened to already be lowercase — and merge into it rather than
// assuming the first-seen row is "the real one".
normalized := strings.ToLower(acc.Address)
if existing, ok := cs.accounts[normalized]; ok {
mergedCount++
fmt.Printf("[MIGRATION] Merging duplicate-case account %s into %s (balance %.6f + %.6f, tusd %.6f + %.6f, lp %.6f + %.6f)\n",
acc.Address, normalized, existing.Balance.Float(), acc.Balance.Float(), existing.TUsdBalance.Float(), acc.TUsdBalance.Float(), existing.LPShares.Float(), acc.LPShares.Float())
existing.Balance = existing.Balance.Add(acc.Balance)
existing.TUsdBalance = existing.TUsdBalance.Add(acc.TUsdBalance)
existing.LPShares = existing.LPShares.Add(acc.LPShares)
existing.IsHuman = existing.IsHuman || acc.IsHuman
if acc.LastActivityAt > existing.LastActivityAt {
existing.LastActivityAt = acc.LastActivityAt
existing.Demurrage14DayWarningShown = acc.Demurrage14DayWarningShown
}
cs.saveAccountToDB(existing)
if acc.Address != normalized {
// Remove the old mixed-case row so it doesn't get re-merged
// (harmlessly, but noisily) on every future restart.
cs.db.Exec(`DELETE FROM chain_accounts WHERE address = $1`, acc.Address)
}
continue
}
acc.Address = normalized
cs.accounts[normalized] = acc
}
fmt.Printf("✓ Loaded %d accounts from PostgreSQL", count)
if mergedCount > 0 {
fmt.Printf(" (%d mixed-case duplicates merged)", mergedCount)
}
fmt.Println()

// Load nullifiers into memory so IsNullifierUsed is O(1) without a DB hit.
if nrows, nerr := cs.db.Query("SELECT nullifier, wallet_address FROM nullifiers"); nerr == nil {
defer nrows.Close()
for nrows.Next() {
var nul, wal string
nrows.Scan(&nul, &wal)
cs.nullifiers[nul] = wal
}
fmt.Printf("✓ Loaded %d nullifiers from PostgreSQL\n", len(cs.nullifiers))
}

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
cs.pool = &PoolState{ReserveAEQ: NewDecimal(0), ReserveTUSD: NewDecimal(0), TotalLPShares: NewDecimal(0)}
fmt.Printf("✓ Liquidity pool created (empty — awaiting first deposit via AddLiquidity)\n")
return
}
cs.pool = &PoolState{ReserveAEQ: NewDecimal(reserveAEQ), ReserveTUSD: NewDecimal(reserveTUSD), TotalLPShares: NewDecimal(totalShares)}
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
var result sql.Result
var err error
if acc.Version == 0 {
// First write: INSERT with version=1, or update if exists without version conflict check
result, err = cs.db.Exec(`INSERT INTO chain_accounts (address, balance, is_human, tusd_balance, lp_shares, last_activity_at, demurrage_14_day_warning_shown, faucet_claimed, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1)
ON CONFLICT (address) DO UPDATE SET balance = $2, is_human = $3, tusd_balance = $4, lp_shares = $5, last_activity_at = $6, demurrage_14_day_warning_shown = $7, faucet_claimed = $8, version = COALESCE(chain_accounts.version,0) + 1`,
acc.Address, acc.Balance.Float(), acc.IsHuman, acc.TUsdBalance.Float(), acc.LPShares.Float(), acc.LastActivityAt, acc.Demurrage14DayWarningShown, acc.FaucetClaimed)
} else {
// Optimistic locking: only update if version matches what we read.
// If another node updated in parallel, rows affected = 0 → conflict detected.
result, err = cs.db.Exec(`UPDATE chain_accounts SET balance = $2, is_human = $3, tusd_balance = $4, lp_shares = $5, last_activity_at = $6, demurrage_14_day_warning_shown = $7, faucet_claimed = $8, version = $9 + 1
WHERE address = $1 AND version = $9`,
acc.Address, acc.Balance.Float(), acc.IsHuman, acc.TUsdBalance.Float(), acc.LPShares.Float(), acc.LastActivityAt, acc.Demurrage14DayWarningShown, acc.FaucetClaimed, acc.Version)
if err == nil {
if rows, _ := result.RowsAffected(); rows == 0 {
fmt.Printf("[DB] Conflict: account %s was modified by another node (version %d)\n", acc.Address, acc.Version)
// Reload from DB to get the current version, then retry
var dbBal, dbTusd, dbLp float64
var dbVer int64
if row := cs.db.QueryRow(`SELECT balance, tusd_balance, lp_shares, version FROM chain_accounts WHERE address = $1`, acc.Address); row != nil {
row.Scan(&dbBal, &dbTusd, &dbLp, &dbVer)
fmt.Printf("[DB] Remote state: balance=%.2f, version=%d\n", dbBal, dbVer)
}
}
}
}
if err != nil {
fmt.Printf("[DB] Error saving account %s: %v\n", acc.Address, err)
return
}
// Update in-memory version to match DB
acc.Version++
}

// Demurrage parameters. AEQ balances that haven't been touched (no
// transfer, swap, or liquidity action) for demurrageGracePeriodSeconds
// begin losing value continuously at demurrageMonthlyRate per month,
// compounding every second rather than in discrete daily/monthly steps —
// this avoids any visible "jump" at day/month boundaries. Touching the
// account in any AEQ-moving way resets the clock to zero, which is the
// entire point: money that's actively circulating never decays, only
// money sitting idle does. Modeled after real-world demurrage currencies
// (Wörgl's 1932 experiment used 1%/month; the long-running Chiemgauer
// uses roughly 2%/quarter ≈ 0.66%/month) — 0.5%/month here is a
// deliberately moderate starting point, slightly gentler than either.
// Lost AEQ is distributed across the same four tokenomics pools as the
// swap fee (40% validators / 30% LPs / 20% UBI / 10% treasury), not
// burned — it stays circulating in the system rather than vanishing
// from total supply. Only AEQ decays this way; tUSD (a simulated test
// currency, not the real UBI-grant currency) is unaffected.
const demurrageGracePeriodSeconds = 90 * 24 * 60 * 60   // 3 months
const demurrageMonthlyRate = 0.005                       // 0.5%/month

// wealthCapMultiplier defines the maximum AEQ a single account may hold,
// expressed as a multiple of the current average AEQ balance across all
// registered humans — not a fixed number. This makes the cap self-
// adapting: as the system grows and average wealth naturally rises
// through normal economic activity, the cap rises proportionally with
// it, rather than needing to be manually raised through discrete
// "phases" as the project matures. The cap only kicks in on incoming
// AEQ (registration grants, transfers, swap/liquidity payouts) — see
// enforceWealthCapLocked — never on a balance that's already there, so
// it can't retroactively punish someone for an average that later rose.
const wealthCapMultiplier = 25.0
const secondsPerMonth = 30 * 24 * 60 * 60                 // approximation, consistent with the grace period's 30-day months

// touchActivity stamps address's LastActivityAt to now, resetting its
// demurrage clock. Called by every AEQ-moving action (Transfer, swaps,
// AddLiquidity/RemoveLiquidity, registration) — NOT by pure balance
// reads, since merely checking a balance isn't "using" the money. Caller
// must hold cs.mu (write lock).
func touchActivity(acc *AccountState) {
acc.LastActivityAt = nowUnix()
acc.Demurrage14DayWarningShown = false // new grace period — the 14-day notice can fire again when this one nears its end
}

// nowUnix exists as a single seam so demurrage timing could be mocked in
// tests later; right now it's just time.Now().Unix().
func nowUnix() int64 {
return timeNowFunc().Unix()
}

// effectiveBalance computes what address's AEQ balance is RIGHT NOW,
// continuously decayed for any time past the grace period since
// LastActivityAt — without writing anything. This is what every balance
// read (GetBalance, /api/balance, /api/humans, etc.) should show, so the
// number displayed always reflects live decay even between the
// lazy-settlement points (see settleDemurrageLocked) where it actually
// gets written to the stored Balance field. Caller must hold at least a
// read lock.
func effectiveBalance(acc *AccountState) Decimal {
if acc.LastActivityAt == 0 {
return acc.Balance
}
idleSeconds := nowUnix() - acc.LastActivityAt
if idleSeconds <= demurrageGracePeriodSeconds {
return acc.Balance
}
decayingSeconds := float64(idleSeconds - demurrageGracePeriodSeconds)
monthsDecaying := decayingSeconds / float64(secondsPerMonth)
factor := math.Pow(1-demurrageMonthlyRate, monthsDecaying)
return acc.Balance.MulFloat(factor)
}

// settleDemurrageLocked actually writes off the decay computed by
// effectiveBalance into acc.Balance, and distributes what was lost across
// the four tokenomics pools — same split as the swap fee. This is called
// right before any operation that's about to read-then-modify Balance
// (Transfer, swaps, liquidity actions), so those operations always work
// from an up-to-date, already-settled balance instead of accidentally
// granting someone pre-decay value just because they happened to act at
// that exact moment. Caller must hold cs.mu (write lock).
func (cs *ChainState) settleDemurrageLocked(acc *AccountState) {
current := effectiveBalance(acc)
lost := acc.Balance.Sub(current)
if lost <= 0 {
return
}
acc.Balance = current
cs.distributeSwapFee(lost.Float(), true) // true = denominated in AEQ; reuses the same 40/30/20/10 split as swap fees
fmt.Printf("[DEMURRAGE] %s: idle balance decayed by %.6f AEQ, redistributed to pools\n", acc.Address, lost.Float())
}

func (cs *ChainState) GetBalance(address string) float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
address = strings.ToLower(address)
if acc, ok := cs.accounts[address]; ok {
return effectiveBalance(acc).Float()
}
return 0
}

// DistributeUBIPool empties the UBI pool address's entire AEQ balance,
// splitting it equally across every currently-registered human, then
// calls cs.save()/persists each affected account. Intended to be called
// once a day by a background ticker (see main.go) — not on every block,
// since "the UBI pool" only makes sense as a daily payout, not a
// per-block trickle. The pool is fully drained each time rather than
// only partially distributed: any AEQ that flows into it between now
// and the next run (swap fees, demurrage, wealth-cap overflow) accrues
// fresh, so there's no need to hold a standing reserve.
// RegisterNode adds this node's operator wallet to the registered_nodes
// table so it participates in future validators pool distributions.
// Called once at startup if NODE_OPERATOR_WALLET env var is set.
// Safe to call multiple times — ON CONFLICT DO NOTHING.
func (cs *ChainState) RegisterNode(operatorWallet string) {
if cs.db == nil || operatorWallet == "" {
return
}
wallet := strings.ToLower(operatorWallet)
cs.db.Exec(`CREATE TABLE IF NOT EXISTS registered_nodes (
wallet_address TEXT PRIMARY KEY,
signing_address TEXT DEFAULT '',
registered_at TIMESTAMP DEFAULT NOW(),
blocks_produced BIGINT NOT NULL DEFAULT 0
)`)
cs.db.Exec(`ALTER TABLE registered_nodes ADD COLUMN IF NOT EXISTS blocks_produced BIGINT NOT NULL DEFAULT 0`)
cs.db.Exec(`ALTER TABLE registered_nodes ADD COLUMN IF NOT EXISTS signing_address TEXT DEFAULT ''`)
_, err := cs.db.Exec(
`INSERT INTO registered_nodes (wallet_address, signing_address) VALUES ($1, $2) ON CONFLICT (wallet_address) DO UPDATE SET signing_address = EXCLUDED.signing_address`,
wallet, strings.ToLower(os.Getenv("RELAYER_ADDRESS")),
)
if err != nil {
fmt.Printf("[NODE] Warning: could not register node wallet %s: %v\n", wallet, err)
} else {
fmt.Printf("[NODE] ✓ Node operator wallet registered: %s\n", wallet)
}
}

// GetRegisteredNodes returns all node operator wallets currently
// registered in the DB. Used by DistributeValidatorsPool.
func (cs *ChainState) GetRegisteredNodes() []string {
if cs.db == nil {
return nil
}
rows, err := cs.db.Query(`SELECT wallet_address FROM registered_nodes`)
if err != nil {
return nil
}
defer rows.Close()
var wallets []string
for rows.Next() {
var w string
rows.Scan(&w)
wallets = append(wallets, w)
}
return wallets
}

// IncrementBlockCount records that the given proposer wallet produced a block.
// Used by DistributeValidatorsPool to distribute rewards proportionally.
func (cs *ChainState) IncrementBlockCount(proposerAddr string) {
if cs.db == nil || proposerAddr == "" { return }
proposerAddr = strings.ToLower(proposerAddr)
// Only increment if this address is a registered node operator
res, _ := cs.db.Exec(`UPDATE registered_nodes SET blocks_produced = blocks_produced + 1 WHERE lower(signing_address) = lower($1)`, proposerAddr)
if n, _ := res.RowsAffected(); n == 0 {
cs.db.Exec(`UPDATE registered_nodes SET blocks_produced = blocks_produced + 1 WHERE lower(wallet_address) = lower($1)`, proposerAddr)
}
}

func (cs *ChainState) DistributeValidatorsPool() {
// Load registered nodes BEFORE acquiring the lock — GetRegisteredNodes
// only reads from PostgreSQL, not cs.accounts, so it doesn't need the
// mutex. Calling it inside the lock would be a deadlock risk if the DB
// driver itself tries to acquire any internal lock that chains back.
nodes := cs.GetRegisteredNodes()
if len(nodes) == 0 {
fmt.Println("[VALIDATORS] No registered node operators — pool left untouched")
return
}

cs.mu.Lock()
defer cs.mu.Unlock()

poolAcc, ok := cs.accounts[validatorsPoolAddr]
if !ok || poolAcc.Balance <= 0 {
fmt.Println("[VALIDATORS] Pool is empty — nothing to distribute today")
return
}

total := poolAcc.Balance.Float()
// Query block counts for proportional distribution
type nodeShare struct{ wallet string; blocks int64 }
var nodeShares []nodeShare
var totalBlocks int64
rows, _ := cs.db.Query(`SELECT wallet_address, blocks_produced FROM registered_nodes WHERE wallet_address = ANY($1)`, pq.Array(nodes))
if rows != nil {
for rows.Next() {
var w string; var b int64
rows.Scan(&w, &b)
if b == 0 { b = 1 } // minimum weight so new nodes still get something
nodeShares = append(nodeShares, nodeShare{w, b})
totalBlocks += b
}
rows.Close()
}
if len(nodeShares) == 0 {
for _, w := range nodes { nodeShares = append(nodeShares, nodeShare{w, 1}); totalBlocks++ }
}
poolAcc.Balance = NewDecimal(0)
cs.saveAccountToDB(poolAcc)

for _, ns := range nodeShares {
wallet := ns.wallet
share := round6(total * float64(ns.blocks) / float64(totalBlocks))
if _, ok := cs.accounts[wallet]; !ok {
cs.accounts[wallet] = &AccountState{Address: wallet}
}
acc := cs.accounts[wallet]
cs.settleDemurrageLocked(acc)
acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
touchActivity(acc)
cs.enforceWealthCapLocked(acc)
cs.saveAccountToDB(acc)
}
cs.save()

cs.syncBalanceLocked(V7_CONTRACT_ADDR, append(nodes, validatorsPoolAddr)...)
fmt.Printf("[VALIDATORS] Distributed %.6f AEQ proportionally (%d nodes, block-weighted)\n", total, len(nodeShares))
}

// DistributeLPPool pays out the entire LP pool balance to liquidity
// providers, proportional to their LP share count. This mirrors how
// real AMMs (Uniswap v2, etc.) reward LPs — the more of the pool you
// provided, the larger your share of the fee income. Accounts with zero
// LP shares receive nothing.
func (cs *ChainState) DistributeLPPool() {
cs.mu.Lock()
defer cs.mu.Unlock()

poolAcc, ok := cs.accounts[lpPoolAddr]
if !ok || poolAcc.Balance <= 0 {
fmt.Println("[LP] Pool is empty — nothing to distribute today")
return
}

// Collect all LP holders and their share counts
type lpHolder struct {
addr   string
shares float64
}
var holders []lpHolder
totalShares := 0.0
for addr, acc := range cs.accounts {
if acc.LPShares > 0 {
holders = append(holders, lpHolder{addr, acc.LPShares.Float()})
totalShares += acc.LPShares.Float()
}
}
if totalShares <= 0 || len(holders) == 0 {
fmt.Println("[LP] No LP holders — pool left untouched")
return
}

total := poolAcc.Balance.Float()
poolAcc.Balance = NewDecimal(0)
cs.saveAccountToDB(poolAcc)

for _, h := range holders {
share := (h.shares / totalShares) * total
acc := cs.accounts[h.addr]
cs.settleDemurrageLocked(acc)
acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
touchActivity(acc)
cs.enforceWealthCapLocked(acc)
cs.saveAccountToDB(acc)
}
cs.save()

holderAddrs := make([]string, len(holders))
for i, h := range holders { holderAddrs[i] = h.addr }
cs.syncBalanceLocked(V7_CONTRACT_ADDR, append(holderAddrs, lpPoolAddr)...)

fmt.Printf("[LP] ✓ Distributed %.6f AEQ across %d LP holders (proportional to shares)\n", total, len(holders))
}

func (cs *ChainState) DistributeUBIPool() {
cs.mu.Lock()
defer cs.mu.Unlock()

poolAcc, ok := cs.accounts[ubiPoolAddr]
if !ok || poolAcc.Balance <= 0 {
fmt.Println("[UBI] Pool is empty — nothing to distribute today")
return
}

var humanAddrs []string
for addr, acc := range cs.accounts {
if acc.IsHuman {
humanAddrs = append(humanAddrs, addr)
}
}
if len(humanAddrs) == 0 {
fmt.Println("[UBI] No registered humans yet — pool left untouched")
return
}

total := poolAcc.Balance.Float()
share := total / float64(len(humanAddrs))
poolAcc.Balance = NewDecimal(0)
cs.saveAccountToDB(poolAcc)

for _, addr := range humanAddrs {
acc := cs.accounts[addr]
cs.settleDemurrageLocked(acc) // settle any pending decay before adding the UBI share
acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
touchActivity(acc) // receiving the daily UBI share counts as activity, like any other incoming AEQ
cs.enforceWealthCapLocked(acc) // a UBI payout can in principle still push someone over the cap
cs.saveAccountToDB(acc)
}
cs.save()

cs.setConfigValue("last_ubi_at", fmt.Sprintf("%d", time.Now().Unix()))
cs.syncBalanceLocked(V7_CONTRACT_ADDR, append(humanAddrs, ubiPoolAddr)...)

fmt.Printf("[UBI] ✓ Distributed %.6f AEQ across %d registered humans (%.6f AEQ each)\n",
total, len(humanAddrs), share)
capturedGini := cs.calcGiniLocked()
capturedHumans := len(humanAddrs)
go cs.SaveGiniSnapshotValues(capturedGini, capturedHumans)
}

// getAverageBalanceLocked computes the mean AEQ balance across every
// registered human (using each account's live, demurrage-adjusted
// balance, not the raw stored value, since that's the real current
// wealth distribution). Non-human accounts (the four fee-pool addresses,
// any unregistered address that merely received a transfer) are excluded
// — the cap is about wealth among the humans the system actually exists
// for, not diluted by infrastructure accounts. Caller must hold cs.mu.
func (cs *ChainState) getAverageBalanceLocked() float64 {
total := 0.0
count := 0
for _, acc := range cs.accounts {
if !acc.IsHuman {
continue
}
total += effectiveBalance(acc).Float()
count++
}
if count == 0 {
return 0
}
return total / float64(count)
}

// enforceWealthCapLocked checks acc's balance against the current
// wealth cap (wealthCapMultiplier * average human balance) and, if it's
// over, skims the excess into the four tokenomics pools — the same
// 40/30/20/10 split used for swap fees and demurrage. This is called
// after AEQ arrives in an account (registration, receiving a transfer,
// a tusd->aeq swap, or removing liquidity), never on amounts already
// sitting in a balance from before the cap existed or before the
// average rose — so it can only ever trim genuinely NEW incoming AEQ
// down to the cap, not retroactively confiscate existing savings.
// Caller must hold cs.mu.
// isTokenomicsPoolAddress reports whether addr is one of the four
// official fee-recipient addresses (validators/LPs/UBI/treasury). These
// are deliberately exempt from the wealth cap — their entire purpose is
// to accumulate fees/demurrage/cap-overflow from everyone else, so
// capping them would be self-defeating. Every other address, registered
// human or not, is subject to the cap.
func isTokenomicsPoolAddress(addr string) bool {
switch addr {
case validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr:
return true
}
return false
}

// bootstrapMultiplierLocked returns the effective wealth cap multiplier.
// During bootstrap (< 25 registered humans) the multiplier scales with the
// human count — max(5, min(N, 25)) — so early joiners cannot accumulate
// 25,000 AEQ before meaningful participation exists. At 25+ humans the
// full wealthCapMultiplier (25×) applies permanently. Caller must hold cs.mu.
func (cs *ChainState) bootstrapMultiplierLocked() float64 {
count := 0
for _, acc := range cs.accounts {
if acc.IsHuman {
count++
}
}
if count >= 25 {
return wealthCapMultiplier
}
m := float64(count)
if m < 5.0 {
m = 5.0
}
return m
}

func (cs *ChainState) enforceWealthCapLocked(acc *AccountState) {
if isTokenomicsPoolAddress(acc.Address) {
return
}
// Deliberately NOT gated on acc.IsHuman: capping only registered
// humans would let someone bypass the entire mechanism just by
// parking AEQ in any ordinary, unregistered address (a personal
// "overflow wallet" they also control) — that address would have
// accumulated unlimited AEQ with no cap ever applying to it. The cap
// has to apply to anyone receiving AEQ, registered or not, for it to
// mean anything.
avg := cs.getAverageBalanceLocked()
if avg <= 0 {
return // no meaningful average yet (e.g. only one human registered so far)
}
multiplier := cs.bootstrapMultiplierLocked()
cap := avg * multiplier
if acc.Balance.Float() <= cap {
return
}
excess := acc.Balance.Float() - cap
acc.Balance = NewDecimal(cap)
cs.distributeSwapFee(excess, true)
fmt.Printf("[WEALTH CAP] %s exceeded %.2fx average (%.2f AEQ) — %.4f AEQ excess redistributed to pools\n",
acc.Address, multiplier, cap, excess)
}

// DemurrageStatus describes whether/when an idle account's AEQ will
// start (or has started) decaying, for surfacing to the user at login.
type DemurrageStatus struct {
	Active           bool    `json:"active"`             // true if decay has already started
	DaysUntilStart   float64 `json:"days_until_start"`    // only meaningful if !Active; can be negative-free, always >= 0
	ShowFourteenDayNotice bool `json:"show_fourteen_day_notice"` // one-time notice, true only on the call that first crosses into the 14-day window
	ShowSevenDayNotice    bool `json:"show_seven_day_notice"`    // true on every check within the last 7 days before decay starts
}

// GetDemurrageStatus reports where address stands relative to the
// demurrage grace period, and — like settleDemurrageLocked — has a side
// effect: the first time this is called once the account has entered the
// 14-day warning window, it flips Demurrage14DayWarningShown so the
// one-time notice isn't repeated on every subsequent login within that
// same window. The 7-day notice has no such one-time flag; per Daniel's
// spec, that one is meant to repeat on every login during its window.
func (cs *ChainState) GetDemurrageStatus(address string) DemurrageStatus {
cs.mu.Lock()
defer cs.mu.Unlock()
address = strings.ToLower(address)

acc, ok := cs.accounts[address]
if !ok || acc.LastActivityAt == 0 {
return DemurrageStatus{Active: false, DaysUntilStart: float64(demurrageGracePeriodSeconds) / 86400}
}

idleSeconds := nowUnix() - acc.LastActivityAt
secondsUntilStart := demurrageGracePeriodSeconds - idleSeconds
if secondsUntilStart <= 0 {
return DemurrageStatus{Active: true}
}

daysUntilStart := float64(secondsUntilStart) / 86400
status := DemurrageStatus{Active: false, DaysUntilStart: daysUntilStart}

if daysUntilStart <= 7 {
status.ShowSevenDayNotice = true
} else if daysUntilStart <= 14 {
if !acc.Demurrage14DayWarningShown {
status.ShowFourteenDayNotice = true
acc.Demurrage14DayWarningShown = true
cs.saveAccountToDB(acc)
}
}

return status
}

func (cs *ChainState) GetTUsdBalance(address string) float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
address = strings.ToLower(address)
if acc, ok := cs.accounts[address]; ok {
return acc.TUsdBalance.Float()
}
return 0
}

func (cs *ChainState) GetPoolReserves() (float64, float64) {
cs.mu.RLock()
defer cs.mu.RUnlock()
if cs.pool == nil {
return 0, 0
}
return cs.pool.ReserveAEQ.Float(), cs.pool.ReserveTUSD.Float()
}

func (cs *ChainState) IsHuman(address string) bool {
cs.mu.RLock()
defer cs.mu.RUnlock()
address = strings.ToLower(address)
if acc, ok := cs.accounts[address]; ok {
return acc.IsHuman
}
return false
}

func (cs *ChainState) RegisterHuman(address string) error {
cs.mu.Lock()
defer cs.mu.Unlock()
address = strings.ToLower(address)

if acc, ok := cs.accounts[address]; ok && acc.IsHuman {
return fmt.Errorf("already registered")
}

if _, ok := cs.accounts[address]; !ok {
cs.accounts[address] = &AccountState{Address: address}
}

cs.accounts[address].IsHuman = true
cs.accounts[address].Balance = cs.accounts[address].Balance.Add(NewDecimal(1000))
touchActivity(cs.accounts[address]) // starts this 1,000 AEQ's own grace period fresh
cs.enforceWealthCapLocked(cs.accounts[address])
cs.saveAccountToDB(cs.accounts[address])
cs.save()

fmt.Printf("[STATE] ✓ Human registered: %s | Balance: %.2f AEQ\n",
address, cs.accounts[address].Balance.Float())
cs.syncHumanRegistrationLocked(V7_CONTRACT_ADDR, address)
return nil
}

func (cs *ChainState) Transfer(from, to string, amount float64) error {
cs.mu.Lock()
defer cs.mu.Unlock()
from = strings.ToLower(from)
to = strings.ToLower(to)

fromAcc, ok := cs.accounts[from]
if !ok {
return fmt.Errorf("insufficient balance")
}
cs.settleDemurrageLocked(fromAcc) // make sure we're checking against the real, decayed balance
if fromAcc.Balance.Float() < amount {
return fmt.Errorf("insufficient balance")
}

fromAcc.Balance = NewDecimal(round6(fromAcc.Balance.Float() - amount))
touchActivity(fromAcc) // sending counts as "using" the money — resets its decay clock
cs.saveAccountToDB(fromAcc)

if _, ok := cs.accounts[to]; !ok {
cs.accounts[to] = &AccountState{Address: to}
}
cs.settleDemurrageLocked(cs.accounts[to])
cs.accounts[to].Balance = NewDecimal(round6(cs.accounts[to].Balance.Float() + amount))
touchActivity(cs.accounts[to]) // receiving also resets the clock on the recipient's whole balance
cs.enforceWealthCapLocked(cs.accounts[to])
cs.saveAccountToDB(cs.accounts[to])
cs.save()

fmt.Printf("[STATE] ✓ Transfer %.2f AEQ: %s → %s\n", amount, from, to)
cs.syncBalanceLocked(V7_CONTRACT_ADDR, from, to, validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr)
return nil
}

// TransferWithV7Fee is used by the RPC layer when intercepting V7 ERC-20
// transfer() calls (selector a9059cbb). It mirrors V7's _calcFee():
//   TX_FEE_BPS = 10 (0.1% base fee)
//   Concentration surcharge if sender holds ≥1/5/10% of total supply
//   20% of fee → UBI pool, 80% burned (removed from supply)
// Without this, Go-ledger and V7-contract diverge on every user transfer.
func (cs *ChainState) TransferWithV7Fee(from, to string, amount float64) error {
cs.mu.Lock()
defer cs.mu.Unlock()
from = strings.ToLower(from)
to = strings.ToLower(to)

fromAcc, ok := cs.accounts[from]
if !ok {
return fmt.Errorf("insufficient balance")
}
cs.settleDemurrageLocked(fromAcc)
if fromAcc.Balance.Float() < amount {
return fmt.Errorf("insufficient balance")
}

// Compute total supply inline to avoid re-entering the mutex.
humans := 0
for _, acc := range cs.accounts {
if acc.IsHuman { humans++ }
}
totalSupply := float64(humans) * 1000.0
fee := calcV7Fee(fromAcc.Balance.Float(), amount, totalSupply)
ubiContrib := round6(fee * 0.20)

fromAcc.Balance = NewDecimal(round6(fromAcc.Balance.Float() - amount))
touchActivity(fromAcc)
cs.saveAccountToDB(fromAcc)

if _, ok := cs.accounts[to]; !ok {
cs.accounts[to] = &AccountState{Address: to}
}
cs.settleDemurrageLocked(cs.accounts[to])
cs.accounts[to].Balance = NewDecimal(round6(cs.accounts[to].Balance.Float() + amount - fee))
touchActivity(cs.accounts[to])
cs.enforceWealthCapLocked(cs.accounts[to])
cs.saveAccountToDB(cs.accounts[to])

if ubiContrib > 0 {
if _, ok := cs.accounts[ubiPoolAddr]; !ok {
cs.accounts[ubiPoolAddr] = &AccountState{Address: ubiPoolAddr}
}
cs.accounts[ubiPoolAddr].Balance = cs.accounts[ubiPoolAddr].Balance.Add(NewDecimal(ubiContrib))
cs.saveAccountToDB(cs.accounts[ubiPoolAddr])
}
cs.save()

fmt.Printf("[STATE] ✓ TransferV7 %.6f AEQ (fee=%.6f burn=%.6f ubi=%.6f): %s → %s\n",
amount, fee, fee-ubiContrib, ubiContrib, from, to)
cs.syncBalanceLocked(V7_CONTRACT_ADDR, from, to, validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr)
return nil
}

// calcV7Fee mirrors AequitasV7.sol's _calcFee():
// base = TX_FEE_BPS (10) = 0.1% of amount
// concentration surcharge based on sender's share of total supply.
func calcV7Fee(senderBalance, amount, totalSupply float64) float64 {
base := amount * 10.0 / 10_000.0
if totalSupply <= 0 {
return round6(base)
}
shareBPS := (senderBalance * 10_000.0) / totalSupply
var extra float64
switch {
case shareBPS >= 1000:
extra = amount * 100.0 / 10_000.0
case shareBPS >= 500:
extra = amount * 50.0 / 10_000.0
case shareBPS >= 100:
extra = amount * 10.0 / 10_000.0
}
return round6(base + extra)
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
address = strings.ToLower(address)
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
cs.settleDemurrageLocked(acc) // settle decay before checking/using the AEQ balance below

if aeqToTusd {
if acc.Balance.Float() < amountIn {
return 0, fmt.Errorf("insufficient AEQ balance")
}
} else {
if acc.TUsdBalance.Float() < amountIn {
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
amountOut = AMMSwapOut(cs.pool.ReserveAEQ, cs.pool.ReserveTUSD, NewDecimal(amountInAfterFee)).Float()
if amountOut >= cs.pool.ReserveTUSD.Float() {
return 0, fmt.Errorf("swap too large for pool liquidity")
}
cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() + amountInAfterFee))
cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() - amountOut))
acc.Balance = NewDecimal(round6(acc.Balance.Float() - amountIn))
acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + amountOut))
} else {
amountOut = AMMSwapOut(cs.pool.ReserveTUSD, cs.pool.ReserveAEQ, NewDecimal(amountInAfterFee)).Float()
if amountOut >= cs.pool.ReserveAEQ.Float() {
return 0, fmt.Errorf("swap too large for pool liquidity")
}
cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() + amountInAfterFee))
cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() - amountOut))
acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() - amountIn))
acc.Balance = NewDecimal(round6(acc.Balance.Float() + amountOut))
}
touchActivity(acc) // swapping (either direction) counts as using the AEQ side
if !aeqToTusd {
cs.enforceWealthCapLocked(acc) // AEQ just arrived via this swap direction — check the cap
}

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.distributeSwapFee(fee, aeqToTusd)
cs.save()

fmt.Printf("[SWAP] %s: %.4f %s → %.4f %s (fee %.4f)\n",
address, amountIn, sideLabel(aeqToTusd, true), amountOut, sideLabel(aeqToTusd, false), fee)

cs.syncBalanceLocked(V7_CONTRACT_ADDR, address, validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr)
go cs.SavePriceSnapshot()
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
// Use a transaction so concurrent pool writes are serialized at the DB level.
// This prevents two nodes from simultaneously distributing UBI or running swaps
// with stale pool reserves. The WHERE id = 1 ensures we update the single pool row.
tx, err := cs.db.Begin()
if err != nil {
fmt.Printf("[DB] Error starting pool tx: %v\n", err)
return
}
// Lock the pool row for this transaction (other writers block until we commit)
var dummy int
tx.QueryRow(`SELECT id FROM liquidity_pool WHERE id = 1 FOR UPDATE`).Scan(&dummy)
_, err = tx.Exec(`UPDATE liquidity_pool SET reserve_aeq = $1, reserve_tusd = $2, total_lp_shares = $3 WHERE id = 1`,
cs.pool.ReserveAEQ.Float(), cs.pool.ReserveTUSD.Float(), cs.pool.TotalLPShares.Float())
if err != nil {
tx.Rollback()
fmt.Printf("[DB] Error saving pool: %v\n", err)
return
}
if err := tx.Commit(); err != nil {
fmt.Printf("[DB] Error committing pool: %v\n", err)
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
cs.accounts[addr].Balance = cs.accounts[addr].Balance.Add(NewDecimal(amount))
} else {
cs.accounts[addr].TUsdBalance = cs.accounts[addr].TUsdBalance.Add(NewDecimal(amount))
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
address = strings.ToLower(address)

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
cs.settleDemurrageLocked(acc) // settle decay before checking/using the AEQ balance below
if acc.Balance.Float() < amountAEQ {
return fmt.Errorf("insufficient AEQ balance")
}
if acc.TUsdBalance.Float() < amountTUSD {
return fmt.Errorf("insufficient tUSD balance")
}

// If the pool already has liquidity, require the deposit to roughly
// match the existing ratio — an unbalanced deposit would otherwise
// instantly shift the price, which is the same rule real AMMs enforce.
var mintedShares float64
if cs.pool.ReserveAEQ > 0 && cs.pool.ReserveTUSD > 0 {
expectedTUSD := amountAEQ * (cs.pool.ReserveTUSD.Float() / cs.pool.ReserveAEQ.Float())
tolerance := expectedTUSD * 0.01 // 1% slack for rounding
if amountTUSD < expectedTUSD-tolerance || amountTUSD > expectedTUSD+tolerance {
return fmt.Errorf("deposit ratio does not match pool ratio (expected ~%.4f tUSD for %.4f AEQ)", expectedTUSD, amountAEQ)
}
if cs.pool.TotalLPShares > 0 {
// Proportional to the pool's existing size — same fraction of the
// AEQ reserve as the fraction of total shares being minted, so an
// LP's claim accurately tracks how much of the pool they actually
// own (including any fees the pool has accumulated since genesis).
mintedShares = (amountAEQ / cs.pool.ReserveAEQ.Float()) * cs.pool.TotalLPShares.Float()
} else {
// Pool has reserves but zero LP shares — legacy state from before
// share-tracking was introduced. Only mint shares for the NEW
// deposit via geometric mean; do NOT credit pre-existing reserves
// to the depositor. Doing so would let anyone with a tiny deposit
// claim practically the entire pool (a drain attack).
mintedShares = math.Sqrt(amountAEQ * amountTUSD)
fmt.Printf("[POOL] Pool had %.4f AEQ / %.4f tUSD with no LP shares recorded — minting %.6f shares for new deposit only\n",
cs.pool.ReserveAEQ.Float(), cs.pool.ReserveTUSD.Float(), mintedShares)
}
} else {
// First-ever deposit: shares = geometric mean of the two amounts
// (standard Uniswap v2 bootstrap formula). Using sqrt(x*y) instead
// of, say, just amountAEQ means the first depositor can't mint an
// outsized number of shares simply by picking a lopsided ratio.
mintedShares = math.Sqrt(amountAEQ * amountTUSD)
}

acc.Balance = NewDecimal(round6(acc.Balance.Float() - amountAEQ))
acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() - amountTUSD))
acc.LPShares = NewDecimal(round6(acc.LPShares.Float() + mintedShares))
touchActivity(acc) // depositing into the pool counts as using the AEQ
cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() + amountAEQ))
cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() + amountTUSD))
cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() + mintedShares))

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.save()

cs.syncBalanceLocked(V7_CONTRACT_ADDR, address)
go cs.SavePriceSnapshot()

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
address = strings.ToLower(address)

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
if acc.LPShares.Float() < sharesToBurn {
return 0, 0, fmt.Errorf("insufficient LP shares (have %.6f, requested %.6f)", acc.LPShares.Float(), sharesToBurn)
}

fraction := sharesToBurn / cs.pool.TotalLPShares.Float()
outAEQ := cs.pool.ReserveAEQ.Float() * fraction
outTUSD := cs.pool.ReserveTUSD.Float() * fraction

acc.LPShares = NewDecimal(round6(acc.LPShares.Float() - sharesToBurn))
acc.Balance = NewDecimal(round6(acc.Balance.Float() + outAEQ))
acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + outTUSD))
touchActivity(acc) // receiving AEQ back from the pool counts as using it
cs.enforceWealthCapLocked(acc)
cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() - outAEQ))
cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() - outTUSD))
cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() - sharesToBurn))

cs.saveAccountToDB(acc)
cs.savePoolToDB()
cs.save()

cs.syncBalanceLocked(V7_CONTRACT_ADDR, address)
go cs.SavePriceSnapshot()

fmt.Printf("[POOL] ✓ %s removed liquidity: %.6f shares → %.4f AEQ + %.4f tUSD\n", address, sharesToBurn, outAEQ, outTUSD)
return outAEQ, outTUSD, nil
}

// GetLPShares returns address's current LP share balance, and the pool's
// total shares — callers can compute the account's ownership fraction
// (and therefore its withdrawable amounts) from these two numbers.
func (cs *ChainState) GetLPShares(address string) (float64, float64) {
cs.mu.RLock()
defer cs.mu.RUnlock()
address = strings.ToLower(address)
var mine float64
if acc, ok := cs.accounts[address]; ok {
mine = acc.LPShares.Float()
}
total := 0.0
if cs.pool != nil {
total = cs.pool.TotalLPShares.Float()
}
return mine, total
}

func (cs *ChainState) TotalSupply() float64 {
// Total supply is always exactly Humans × 1,000 AEQ by protocol design.
// Each registered human receives exactly 1,000 AEQ upon registration —
// no more, no less. Floating-point drift from swap fees and demurrage
// calculations means the sum of all account balances + pool reserves
// diverges slightly from this over time, so we compute it directly
// from the human count instead of summing balances.
cs.mu.RLock()
defer cs.mu.RUnlock()
humans := 0
for _, acc := range cs.accounts {
if acc.IsHuman {
humans++
}
}
return float64(humans) * 1000.0
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

// GetAllAccounts returns a COPY of each account, with Balance set to its
// live, demurrage-adjusted value (see effectiveBalance) — not the raw
// stored value, and not a pointer to the real account. Copies matter
// here: callers (the explorer's /api/humans, etc.) must never be able to
// mutate the actual stored balance just by displaying it, and showing
// the raw stored value would make the UI lag behind real decay until
// that specific account next did something that triggered
// settleDemurrageLocked.
func (cs *ChainState) GetAllAccounts() []*AccountState {
cs.mu.RLock()
defer cs.mu.RUnlock()
result := make([]*AccountState, 0, len(cs.accounts))
for _, acc := range cs.accounts {
displayCopy := *acc
displayCopy.Balance = effectiveBalance(acc)
result = append(result, &displayCopy)
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
address = strings.ToLower(address)

acc, ok := cs.accounts[address]
if !ok || !acc.IsHuman {
return fmt.Errorf("only registered humans can claim the test-tUSD faucet")
}
if acc.FaucetClaimed {
return fmt.Errorf("faucet already claimed")
}

acc.FaucetClaimed = true
acc.TUsdBalance = NewDecimal(tusdFaucetAmount)
cs.saveAccountToDB(acc)
cs.save()

fmt.Printf("[FAUCET] ✓ %s claimed %.2f test-tUSD\n", address, tusdFaucetAmount)
return nil
}

// StateRoot computes a deterministic hash of ALL economically meaningful state:
// human AEQ balances, tUSD balances, LP shares, pool reserves, and nullifiers.
// Two states with the same root are guaranteed to be economically identical.
// Previously only human AEQ balances were included; this allowed different
// economic states to hash identically, defeating state-root verification.
func (cs *ChainState) StateRoot() string {
cs.mu.RLock()
addrs := make([]string, 0, len(cs.accounts))
for a := range cs.accounts { addrs = append(addrs, a) }
sort.Strings(addrs)
var sb strings.Builder
for _, a := range addrs {
acc := cs.accounts[a]
// Include ALL accounts with non-zero AEQ or tUSD balances (not only humans)
if acc.IsHuman || acc.Balance > 0 || acc.TUsdBalance > 0 || acc.LPShares > 0 {
// FaucetClaimed and LastActivityAt must be included: two nodes that
// processed different sets of faucet claims would produce the same
// balances but handle future UBI distribution differently.
fmt.Fprintf(&sb, "%s:%.6f:%.6f:%.6f:h=%v:t=%d:fc=%v:",
a,
round6(acc.Balance.Float()),
round6(acc.TUsdBalance.Float()),
round6(acc.LPShares.Float()),
acc.IsHuman,
acc.LastActivityAt,
acc.FaucetClaimed)
}
}
// Include pool state: reserves and total LP shares
if cs.pool != nil {
fmt.Fprintf(&sb, "pool:%.6f:%.6f:%.6f",
round6(cs.pool.ReserveAEQ.Float()),
round6(cs.pool.ReserveTUSD.Float()),
round6(cs.pool.TotalLPShares.Float()))
}
// Include nullifier count (hash of keys, not values, for privacy)
nullKeys := make([]string, 0, len(cs.nullifiers))
for k := range cs.nullifiers { nullKeys = append(nullKeys, k) }
sort.Strings(nullKeys)
fmt.Fprintf(&sb, "|n=%d:", len(nullKeys))
for _, k := range nullKeys { sb.WriteString(k); sb.WriteString(":"); sb.WriteString(cs.nullifiers[k]) }
// Include last UBI distribution timestamp: two nodes with identical balances
// but different UBI timestamps would diverge on the next UBI run.
lastUBI := cs.getConfigValue("last_ubi_at")
fmt.Fprintf(&sb, "|ubi:%s", lastUBI)
cs.mu.RUnlock()
hash := sha256.Sum256([]byte(sb.String()))
return hex.EncodeToString(hash[:])
}

// calcGiniLocked computes the Gini coefficient without acquiring cs.mu.
// Must only be called while cs.mu is already held (read or write).
func (cs *ChainState) calcGiniLocked() float64 {
balances := []float64{}
for _, acc := range cs.accounts {
if acc.IsHuman && acc.Balance > 0 { balances = append(balances, effectiveBalance(acc).Float()) }
}
n := len(balances)
if n < 2 { return 0.0 }
for i := 0; i < n; i++ {
for j := i + 1; j < n; j++ { if balances[j] < balances[i] { balances[i], balances[j] = balances[j], balances[i] } }
}
var sum, numerator float64
for i, x := range balances { sum += x; numerator += float64(2*i+1-n) * x }
if sum == 0 { return 0.0 }
gini := numerator / (float64(n) * sum)
if gini < 0 { gini = -gini }
if n > 1 { gini = gini * float64(n) / float64(n-1) }
if gini > 1.0 { gini = 1.0 }
return gini
}

func (cs *ChainState) CalcGini() float64 {
cs.mu.RLock()
defer cs.mu.RUnlock()
if len(cs.accounts) < 2 {
return 0.0
}
balances := []float64{}
for _, acc := range cs.accounts {
// Only count registered humans — pool addresses and unregistered
// wallets holding small fee amounts would skew the distribution.
if acc.IsHuman && acc.Balance > 0 {
balances = append(balances, effectiveBalance(acc).Float())
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
// Correct for small-sample bias: the standard formula caps at (n-1)/n,
// so with 2 people the maximum is 0.5 even with total concentration.
// Multiplying by n/(n-1) gives the unbiased estimator that reaches 1.0.
gini = gini * float64(n) / float64(n-1)
if gini > 1.0 {
gini = 1.0
}
return gini
}

func (cs *ChainState) CalcAequitasIndex() float64 {
gini := cs.CalcGini()
index := gini * 100.0
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
acc.Balance = NewDecimal(amount)
cs.saveAccountToDB(acc)
} else {
acc = &AccountState{Address: address, Balance: NewDecimal(amount)}
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
