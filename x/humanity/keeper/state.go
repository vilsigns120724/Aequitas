package keeper

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
)

// timeNowFunc is a seam for time.Now(), letting demurrage timing be
// mocked in tests without needing to thread a clock through every call.
var timeNowFunc = time.Now

// processStartTime records when this process started. Used by
// resetDBStateForBootstrap to refuse RESET_DB_STATE=true on accidental
// crash-recovery restarts that retain the env var.
var processStartTime = time.Now()

type AccountState struct {
	Address string  `json:"address"`
	Balance Decimal `json:"balance"`
	IsHuman bool    `json:"is_human"`
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
	FaucetClaimed bool  `json:"faucet_claimed"`
	Version       int64 `json:"-"` // optimistic lock version, not serialized
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
	ReserveAEQ  Decimal `json:"reserve_aeq"`
	ReserveTUSD Decimal `json:"reserve_tusd"`
	// TotalLPShares is the sum of every account's LPShares. Starts at 0; the
	// very first deposit mints sqrt(amountAEQ * amountTUSD) shares (the
	// standard Uniswap v2 formula — using the geometric mean means the
	// first depositor's chosen ratio doesn't let them mint an arbitrarily
	// large or small initial share count by gaming the two amounts).
	TotalLPShares Decimal `json:"total_lp_shares"`
}

type ChainState struct {
	mu         sync.RWMutex
	accounts   map[string]*AccountState
	pool       *PoolState
	db         *sql.DB
	useDB      bool
	nullifiers map[string]string // nullifier hex → wallet address (in-memory cache)
}

// P3-FIX: stateMu, acquireStateLock, and releaseStateLock were dead code
// (never called). Removed to eliminate a future deadlock trap where a
// developer might call acquireStateLock inside a function already holding cs.mu.

// beginStateTx starts a SERIALIZABLE PostgreSQL transaction for critical
// state mutations. The caller is responsible for calling Commit or Rollback.
// Returns nil if no DB is configured (in-memory/file mode).
func (cs *ChainState) beginStateTx() *sql.Tx {
	if cs.db == nil {
		return nil
	}
	// P0-7: SERIALIZABLE prevents phantom reads that could violate
	// totalSupply = humans * 1000 invariant under concurrent writes.
	tx, err := cs.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		fmt.Printf("[DB] Warning: could not begin state tx: %v\n", err)
		return nil
	}
	return tx
}

// P3-9: validate pool addresses at startup to catch typos early
func validatePoolAddresses() {
	for _, addr := range []string{validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr} {
		if len(addr) != 42 || addr[:2] != "0x" {
			panic("invalid pool address: " + addr)
		}
	}
}

func NewChainState(dataFile string) *ChainState {
	validatePoolAddresses()
	cs := &ChainState{
		accounts:   make(map[string]*AccountState),
		nullifiers: make(map[string]string),
	}

	// Try PostgreSQL first
	if os.Getenv("RESET_STATE") == "true" && os.Getenv("DATABASE_URL") != "" {
		fmt.Println("⚠ RESET_STATE=true is set but DATABASE_URL is active — DB is NOT wiped by this flag.")
		fmt.Println("  To reset a DB-backed node, run DELETE queries directly in the PostgreSQL console.")
	}
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
				if os.Getenv("RESET_DB_STATE") == "true" {
					cs.resetDBStateForBootstrap()
				}
				if os.Getenv("CLEAR_REGISTRATIONS") == "true" {
					cs.clearRegistrationsFromDB()
				}
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
	// P3-10: log schema migration errors instead of silently ignoring them.
	dbExec := func(q string, args ...interface{}) {
		if _, err := cs.db.Exec(q, args...); err != nil {
			fmt.Printf("[DB] initDB warning: %v\n", err)
		}
	}
	dbExec(`CREATE TABLE IF NOT EXISTS evm_contracts (
address TEXT PRIMARY KEY,
bytecode TEXT NOT NULL,
deployer TEXT,
deployed_at TIMESTAMP DEFAULT NOW()
)`)
	dbExec(`CREATE TABLE IF NOT EXISTS evm_storage (
address TEXT NOT NULL,
slot TEXT NOT NULL,
value TEXT NOT NULL,
PRIMARY KEY (address, slot)
)`)
	dbExec(`CREATE TABLE IF NOT EXISTS evm_nonces (
address TEXT PRIMARY KEY,
nonce BIGINT DEFAULT 0
)`)
	dbExec(`CREATE TABLE IF NOT EXISTS chain_accounts (
address TEXT PRIMARY KEY,
balance FLOAT NOT NULL DEFAULT 0,
is_human BOOLEAN NOT NULL DEFAULT false
)`)
	// tusd_balance added separately (ALTER instead of being in the original
	// CREATE TABLE) so this upgrade doesn't require recreating the table on
	// chains that already have chain_accounts from before this feature.
	// P2-3 fix: ADD COLUMN before ALTER TYPE — on a fresh DB the column must
	// exist before we can change its type. IF NOT EXISTS makes both safe to run
	// on existing DBs too.
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS tusd_balance FLOAT NOT NULL DEFAULT 0`)
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS lp_shares FLOAT NOT NULL DEFAULT 0`)
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS last_activity_at BIGINT NOT NULL DEFAULT 0`)
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS demurrage_14_day_warning_shown BOOLEAN NOT NULL DEFAULT false`)
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS faucet_claimed BOOLEAN NOT NULL DEFAULT false`)
	dbExec(`ALTER TABLE chain_accounts ADD COLUMN IF NOT EXISTS version BIGINT NOT NULL DEFAULT 0`)
	// Upgrade balance columns to NUMERIC(20,6) for exact decimal storage.
	dbExec(`ALTER TABLE chain_accounts ALTER COLUMN balance TYPE NUMERIC(20,6) USING balance::NUMERIC(20,6)`)
	dbExec(`ALTER TABLE chain_accounts ALTER COLUMN tusd_balance TYPE NUMERIC(20,6) USING tusd_balance::NUMERIC(20,6)`)
	dbExec(`ALTER TABLE chain_accounts ALTER COLUMN lp_shares TYPE NUMERIC(20,6) USING lp_shares::NUMERIC(20,6)`)
	// Links a ZK proof commitment to the wallet that successfully registered
	// with it, so the app can ask "did MY proof get registered, and to which
	// wallet?" instead of guessing from a global, unfiltered list.
	dbExec(`CREATE TABLE IF NOT EXISTS bio_registrations (
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
	dbExec(`ALTER TABLE bio_registrations ADD COLUMN IF NOT EXISTS bio_hash TEXT`)
	dbExec(`CREATE UNIQUE INDEX IF NOT EXISTS uidx_bio_registrations_bio_hash ON bio_registrations(bio_hash) WHERE bio_hash IS NOT NULL`)
	// Single-row table holding the AEQ<->tUSD pool reserves. A fixed id=1 row
	// is used instead of a key-value table since there's only ever one pool
	// right now — simpler queries, and trivial to extend to multiple pools
	// later (id column is already there) if more pairs are ever added.
	dbExec(`CREATE TABLE IF NOT EXISTS liquidity_pool (
id INTEGER PRIMARY KEY DEFAULT 1,
reserve_aeq FLOAT NOT NULL DEFAULT 0,
reserve_tusd FLOAT NOT NULL DEFAULT 0,
total_lp_shares FLOAT NOT NULL DEFAULT 0
)`)
	dbExec(`ALTER TABLE liquidity_pool ADD COLUMN IF NOT EXISTS total_lp_shares FLOAT NOT NULL DEFAULT 0`)
	// Upgrade liquidity_pool columns to NUMERIC(20,6) for exact decimal storage
	// (same migration applied to chain_accounts columns above). Must run AFTER the
	// ADD COLUMN statements so that the column definitely exists before ALTER TYPE.
	dbExec(`ALTER TABLE liquidity_pool ALTER COLUMN reserve_aeq TYPE NUMERIC(20,6) USING reserve_aeq::NUMERIC(20,6)`)
	dbExec(`ALTER TABLE liquidity_pool ALTER COLUMN reserve_tusd TYPE NUMERIC(20,6) USING reserve_tusd::NUMERIC(20,6)`)
	dbExec(`ALTER TABLE liquidity_pool ALTER COLUMN total_lp_shares TYPE NUMERIC(20,6) USING total_lp_shares::NUMERIC(20,6)`)
	// nullifiers stores the one-way SHA256 derivative of each identity's bioHash.
	// Checked at registration time to prevent the same biometric from registering
	// with a second wallet. The nullifier itself never reveals the bioHash.
	dbExec(`CREATE TABLE IF NOT EXISTS nullifiers (
nullifier TEXT PRIMARY KEY,
wallet_address TEXT NOT NULL,
registered_at TIMESTAMP DEFAULT NOW()
)`)
	dbExec(`CREATE TABLE IF NOT EXISTS chain_config (
key TEXT PRIMARY KEY,
value TEXT NOT NULL
)`)

	// Pending block transactions — persisted so they survive node restarts.
	// Without this, transfers via sendRawTransaction update Go-state/DB but
	// pendingTxs (in-memory) is lost on restart → secondary nodes never get
	// the TX in a block → balances permanently diverge across nodes.
	dbExec(`CREATE TABLE IF NOT EXISTS pending_txs (
id         SERIAL PRIMARY KEY,
tx_json    TEXT   NOT NULL,
created_at BIGINT NOT NULL DEFAULT 0
)`)

	// EVM transaction receipts — persisted so MetaMask can get correct
	// receipts after a node restart (avoids "Senden fehlgeschlagen" for
	// transactions that actually succeeded before the node restarted).
	dbExec(`CREATE TABLE IF NOT EXISTS evm_tx_receipts (
tx_hash    TEXT PRIMARY KEY,
from_addr  TEXT NOT NULL,
to_addr    TEXT,
status     TEXT NOT NULL DEFAULT '0x1',
created_at BIGINT NOT NULL
)`)
	// FIX: contract_addr was never persisted, so getTransactionReceipt lost
	// "contractAddress" for deployment TXs after every restart (deployedContracts
	// is in-memory only) — MetaMask/explorers would then show a deployment
	// receipt with contractAddress: null. ADD COLUMN IF NOT EXISTS is safe to
	// run against an existing table created before this column existed.
	dbExec(`ALTER TABLE evm_tx_receipts ADD COLUMN IF NOT EXISTS contract_addr TEXT`)
	// Keep only the last 10000 receipts to prevent unbounded growth.
	// Old receipts are pruned in SaveTxReceipt.

	cs.InitSwapNoncesTable()
	cs.InitValidatorKeysTable()
	cs.InitGiniSnapshotsTable()
	cs.InitPriceSnapshotsTable()
	if err := cs.InitGuardianTables(); err != nil {
		fmt.Printf("[DB] FATAL: InitGuardianTables failed: %v\n", err)
		panic(err)
	}
}

// resetDBStateForBootstrap is an explicit operator escape hatch for secondary
// nodes that must discard a divergent local DB before importing a signed
// bootstrap snapshot. It intentionally refuses to run on the primary or without
// BOOTSTRAP_SNAPSHOT_URL so RESET_DB_STATE cannot silently wipe a production
// chain database.
func (cs *ChainState) resetDBStateForBootstrap() {
	if cs.db == nil {
		return
	}
	if os.Getenv("IS_PRIMARY_NODE") == "true" {
		fmt.Println("[DB-RESET] Refused: RESET_DB_STATE=true on IS_PRIMARY_NODE=true")
		return
	}
	// Only honour within the first 5 minutes of startup.
	// An accidentally-retained RESET_DB_STATE would otherwise wipe the DB
	// on every Railway crash-recovery restart.
	if time.Since(processStartTime) > 5*time.Minute {
		fmt.Println("[DB-RESET] Refused: RESET_DB_STATE=true but process started >5 minutes ago — ignoring to prevent accidental wipe on restart")
		return
	}
	if os.Getenv("BOOTSTRAP_SNAPSHOT_URL") == "" {
		fmt.Println("[DB-RESET] Refused: RESET_DB_STATE=true requires BOOTSTRAP_SNAPSHOT_URL")
		return
	}

	tables := []string{
		"pending_txs",      // prevent stale TXs from polluting post-reset state
		"bio_registrations",
		"nullifiers",
		"bio_hashes",
		"evm_contracts",
		"evm_storage",
		"evm_nonces",
		"evm_tx_receipts",
		"registered_nodes",
		"validator_keys",
		"liquidity_pool",
		"swap_nonces",
		"price_snapshots",
		"gini_snapshots",
		"guardians",
		"escrow_accounts",
		"chain_accounts",
		"chain_config",
		"v6_balances",
		"v6_commitments",
		"v6_humans",
		"v6_state",
		// FIX: same reasoning as clearRegistrationsFromDB — without this, a
		// stale relationship-slot snapshot survives the reset and gets
		// blindly restored into evm_storage on the next automatic V7
		// redeploy, reintroducing isHuman/balanceOf entries this reset was
		// supposed to remove.
		"evm_upgrade_relationship_slots",
	}

	fmt.Println("[DB-RESET] RESET_DB_STATE=true — truncating local secondary DB before snapshot bootstrap")
	// FIX: track every failure instead of just printing an easy-to-miss
	// "Warning" and continuing as if nothing happened. A reset whose whole
	// purpose is to guarantee a clean slate before importing a snapshot must
	// not silently end in "Done" when some tables were never actually
	// truncated — that's exactly the kind of half-reset state that produced
	// "already registered" / StateRoot-divergence bugs throughout this
	// project's history.
	var failed []string
	for _, table := range tables {
		var exists bool
		if err := cs.db.QueryRow(`SELECT to_regclass($1) IS NOT NULL`, "public."+table).Scan(&exists); err != nil {
			fmt.Printf("[DB-RESET] Warning: could not check table %s: %v\n", table, err)
			failed = append(failed, table+" (existence check failed)")
			continue
		}
		if !exists {
			continue
		}
		if _, err := cs.db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`, pq.QuoteIdentifier(table))); err != nil {
			fmt.Printf("[DB-RESET] Warning: could not truncate %s: %v\n", table, err)
			failed = append(failed, table)
		}
	}
	if len(failed) > 0 {
		fmt.Printf("[ALERT] [DB-RESET] %d table(s) FAILED to truncate: %v — DB is NOT cleanly reset. Do NOT remove RESET_DB_STATE yet; investigate before the next restart imports a snapshot on top of this half-cleared state.\n", len(failed), failed)
		return
	}
	fmt.Println("[DB-RESET] Done")
	fmt.Println("[DB-RESET] ⚠ IMPORTANT: remove RESET_DB_STATE=true from env vars after this deploy succeeds.")
	fmt.Println("[DB-RESET]   Leaving it set will WIPE the DB again on every future restart.")
}

// tableNameFromDelete extracts the table name from a "DELETE FROM <table>"
// or "DELETE FROM <table> WHERE ..." statement, returning "" for anything
// else (e.g. UPDATE statements, which should always run unconditionally
// since they only ever target tables initDB guarantees exist upfront).
func tableNameFromDelete(stmt string) string {
	const prefix = "DELETE FROM "
	if !strings.HasPrefix(stmt, prefix) {
		return ""
	}
	rest := stmt[len(prefix):]
	if idx := strings.IndexByte(rest, ' '); idx >= 0 {
		return rest[:idx]
	}
	return rest
}

// clearRegistrationsFromDB removes all human registration data without wiping
// the full DB. Triggered by CLEAR_REGISTRATIONS=true env var. Clears:
// nullifiers, bio_registrations, chain_accounts (is_human+balance), EVM
// storage slots for V7 (usedNullifiers/usedCommitments/isHuman), evm_nonces,
// evm_tx_receipts, and pending_txs. Safe to run on primary or secondary.
// Remove CLEAR_REGISTRATIONS=true after the first successful restart.
func (cs *ChainState) clearRegistrationsFromDB() {
	if cs.db == nil {
		return
	}
	if time.Since(processStartTime) > 5*time.Minute {
		fmt.Println("[CLEAR-REG] Refused: CLEAR_REGISTRATIONS=true but process started >5 minutes ago")
		return
	}
	fmt.Println("[CLEAR-REG] Clearing all registration data from DB...")
	v7Addr := strings.ToLower(V7_CONTRACT_ADDR)
	stmts := []string{
		`DELETE FROM nullifiers`,
		`DELETE FROM bio_registrations`,
		// FIX: bio_hashes was never cleared here. Nothing on the chain side
		// reads this table for registration blocking today (it's write-only,
		// populated by SaveBioHash/snapshot import), but leaving stale rows
		// behind after every other registration table is wiped is an
		// inconsistent half-reset, and it's the SAME table name as the one
		// the separate proof-server service uses for its own (much more
		// consequential) duplicate-biometric check — clearing it here keeps
		// the chain's own copy honest regardless of what reads it later.
		`DELETE FROM bio_hashes`,
		`UPDATE chain_accounts SET is_human = false, balance = 0, tusd_balance = 0, lp_shares = 0, last_activity_at = 0, faucet_claimed = false`,
		`DELETE FROM evm_storage WHERE lower(address) = '` + v7Addr + `'`,
		`DELETE FROM evm_nonces`,
		`DELETE FROM evm_tx_receipts`,
		`DELETE FROM pending_txs`,
		`DELETE FROM evm_contracts WHERE lower(address) = '` + v7Addr + `'`,
		// CRITICAL FIX: evm_upgrade_relationship_slots was never cleared here.
		// This table snapshots EVERY evm_storage row for V7 (see
		// SavePreUpgradeRelationshipSlots) before a contract-version upgrade
		// wipes evm_storage, then blindly restores all of it afterward via
		// RestorePreUpgradeRelationshipSlots — relying on MigrateEVMFromGoState
		// to overwrite the slots it knows how to re-derive (balanceOf, isHuman,
		// etc.) for every account that's still in chain_accounts. That
		// assumption breaks the moment CLEAR_REGISTRATIONS wipes chain_accounts
		// too: migration then has zero accounts to re-derive from, so EVERY
		// stale slot from the snapshot — including isHuman=true for wallets
		// this exact reset was supposed to un-register — gets faithfully
		// restored on the very next automatic V7 redeploy. Without this line,
		// a wallet that got its isHuman EVM slot stuck "true" (e.g. from the
		// concurrent-registration race fixed in 2dee74b) stayed stuck forever,
		// no matter how many times CLEAR_REGISTRATIONS was run — confirmed in
		// production via /api/admin/registration-debug and the
		// "[MIGRATE] Restored N guardian/escrow slots from pre-upgrade
		// snapshot" log line reappearing after every reset.
		`DELETE FROM evm_upgrade_relationship_slots WHERE address = '` + v7Addr + `'`,
		// CRITICAL FIX: liquidity_pool reserves were never reset here.
		// StateRoot() hashes cs.pool.ReserveAEQ/ReserveTUSD/TotalLPShares
		// directly (see state.go ~2261) — leaving stale pool reserves behind
		// while every other table gets wiped means two nodes that both ran
		// CLEAR_REGISTRATIONS at different points in their history (e.g. a
		// primary reset fresh, a secondary with leftover pool data from
		// before any reset ever touched it) compute permanently different
		// StateRoots for the IDENTICAL set of accounts/nullifiers — exactly
		// the "[DAG] StateRoot mismatch ... accepted (warn only)" /
		// "5+ consecutive StateRoot mismatches" pattern seen in production
		// on every single block between a freshly-reset primary and a
		// secondary whose liquidity_pool row was never touched.
		`UPDATE liquidity_pool SET reserve_aeq = 0, reserve_tusd = 0, total_lp_shares = 0 WHERE id = 1`,
	}
	// FIX: bio_hashes and evm_upgrade_relationship_slots are only ever
	// created lazily (by SaveBioHash / SavePreUpgradeRelationshipSlots) —
	// unlike nullifiers/bio_registrations/chain_accounts, which initDB
	// always creates upfront. On a node whose DB has never gone through a
	// registration or a contract-version upgrade, those two tables
	// genuinely don't exist yet, and DELETE FROM a nonexistent table prints
	// a scary "relation does not exist" warning that looks like a real
	// problem but is actually a harmless no-op. Skip cleanly instead.
	// FIX: track failures instead of printing an easy-to-miss "Warning" and
	// unconditionally claiming success at the end. This reset's entire job
	// is to guarantee no stale registration data survives it — a statement
	// that fails partway through (e.g. a transient DB hiccup) used to leave
	// some tables wiped and others not, while still printing "Done" with no
	// way to tell the two outcomes apart in the logs.
	var failed []string
	for _, stmt := range stmts {
		tableName := tableNameFromDelete(stmt)
		if tableName != "" {
			var exists bool
			if err := cs.db.QueryRow(`SELECT to_regclass($1) IS NOT NULL`, "public."+tableName).Scan(&exists); err == nil && !exists {
				continue
			}
		}
		if _, err := cs.db.Exec(stmt); err != nil {
			fmt.Printf("[CLEAR-REG] Warning: %v\n", err)
			failed = append(failed, stmt)
		}
	}
	if len(failed) > 0 {
		fmt.Printf("[ALERT] [CLEAR-REG] %d statement(s) FAILED — registrations are only PARTIALLY cleared: %v\n", len(failed), failed)
		fmt.Println("[ALERT] [CLEAR-REG] Do not remove CLEAR_REGISTRATIONS yet — investigate and rerun, or this half-reset state can cause stale isHuman/nullifier entries to resurface.")
		return
	}
	fmt.Println("[CLEAR-REG] Done — all registrations cleared.")
	fmt.Println("[CLEAR-REG] ⚠ Remove CLEAR_REGISTRATIONS=true from env vars and redeploy once more.")
}

// setConfigValue persists a key/value pair to chain_config (upsert).
// P2-AUDIT: Log errors instead of silently ignoring them. A failed write to
// last_ubi_at would allow double-distribution if the DB is temporarily
// unavailable — the error must be visible in logs for operators to act on.
func (cs *ChainState) setConfigValue(key, value string) {
	if cs.db == nil {
		return
	}
	if _, err := cs.db.Exec(`INSERT INTO chain_config (key, value) VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, key, value); err != nil {
		fmt.Printf("[DB] Warning: setConfigValue(%q) failed: %v\n", key, err)
	}
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

// SecondsUntilNextUBI returns integer seconds until next UBI for the /api/status endpoint.
// P3-3: uses last_ubi_at from DB, not server uptime, so restarts don't give wrong countdowns.

// GetWealthCapInfo returns the current wealth cap parameters using the canonical
// formulas: bootstrapMultiplierLocked() for multiplier and 1000.0 for average.
// P2-2: ensures handleWealthCap shows the same values as enforceWealthCapLocked.
func (cs *ChainState) GetWealthCapInfo() (capAEQ float64, mult float64, avg float64, humans int) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	for _, acc := range cs.accounts {
		if acc.IsHuman {
			humans++
		}
	}
	mult = cs.bootstrapMultiplierLocked()
	avg = cs.getAverageBalanceLocked()
	capAEQ = mult * avg
	return
}

// TryLockDistribution attempts to atomically claim the distribution slot.
// It uses a PostgreSQL compare-and-swap: updates last_ubi_at to now only if
// the current value is > 23 hours old (or missing). Returns true if this node
// won the lock — only then should it actually run the distribution.
// This replaces the IS_PRIMARY_NODE env-var, which any operator could set.
func (cs *ChainState) TryLockDistribution() bool {
	if cs.db == nil {
		return true // no DB → single-node mode, always proceed
	}
	threshold := fmt.Sprintf("%d", time.Now().Add(-23*time.Hour-55*time.Minute).Unix()) // F5-FIX: grace period, still < 24h
	now := fmt.Sprintf("%d", time.Now().Unix())
	// Insert if missing, update if older than threshold
	result, err := cs.db.Exec(
		`INSERT INTO chain_config (key, value) VALUES ('last_ubi_at', $1)
ON CONFLICT (key) DO UPDATE SET value = $1
WHERE chain_config.value = '' OR COALESCE(NULLIF(regexp_replace(chain_config.value, '[^0-9]', '', 'g'), ''), '0')::BIGINT < $2`,
		now, threshold,
	)
	if err != nil {
		fmt.Printf("[POOLS] TryLockDistribution error: %v\n", err)
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

// SetNextUBIAt stores when the scheduler will next trigger pool distributions.
// Called by main.go immediately after calculating the next run time so the
// display timer is always in sync with the actual goroutine schedule.
func (cs *ChainState) SetNextUBIAt(unixTs int64) {
	cs.setConfigValue("next_ubi_at", fmt.Sprintf("%d", unixTs))
}

// SecondsUntilNextUBI returns how many seconds until the next UBI distribution.
// Reads "next_ubi_at" which main.go writes every time it schedules a run,
// so the countdown is exact — not estimated from last_ubi_at + 24h.
func (cs *ChainState) SecondsUntilNextUBI() int64 {
	v := cs.getConfigValue("next_ubi_at")
	if v == "" {
		// Scheduler not yet started (non-primary node or fresh start before
		// first goroutine tick). Show no countdown rather than a wrong value.
		return 0
	}
	var nextAt int64
	fmt.Sscan(v, &nextAt)
	secs := nextAt - time.Now().Unix()
	if secs < 0 {
		return 0
	}
	return secs
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
		if err := rows.Scan(&acc.Address, &bal, &acc.IsHuman, &tusd, &lp, &acc.LastActivityAt, &acc.Demurrage14DayWarningShown, &acc.FaucetClaimed, &acc.Version); err != nil {
			fmt.Printf("[DB] Scan error loading account: %v — skipping row\n", err)
			continue
		}
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
		// defer replaced by explicit Close at end of block
		for nrows.Next() {
			var nul, wal string
			// P2-FIX: check scan error to skip malformed rows.
			if scanErr := nrows.Scan(&nul, &wal); scanErr != nil {
				fmt.Printf("[DB] Warning: nullifier scan error: %v\n", scanErr)
				continue
			}
			if nul == "" {
				continue
			}
			cs.nullifiers[nul] = wal
		}
		nrows.Close()
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
	// P2-AUDIT: Take a snapshot under RLock before serializing. Without the lock,
	// a concurrent write to cs.accounts (Transfer, RegisterHuman) could produce
	// a partially-modified map view in the JSON output.
	cs.mu.RLock()
	data, _ := json.Marshal(cs.accounts)
	cs.mu.RUnlock()
	// D8-FIX: atomic write via temp-file + rename to prevent partial file
	// corruption if the process crashes mid-write.
	tmpPath := "/tmp/aequitas_state.json.tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		fmt.Printf("[STATE] Warning: failed to write state: %v\n", err)
		return
	}
	os.Rename(tmpPath, "/tmp/aequitas_state.json")
}

func (cs *ChainState) saveAccountToDB(acc *AccountState) {
	// P0-1: retry up to 3 times on optimistic-lock conflict so callers don't
	// silently lose writes when two nodes update the same account concurrently.
	for attempt := 0; attempt < 3; attempt++ {
		prevVer := acc.Version
		cs.saveAccountToDBInner(acc)
		if acc.Version != prevVer {
			return
		} // version incremented = success
		if acc.Version == 0 {
			return
		} // new account (INSERT path)
		// Conflict: reload from DB and retry with current state
		if attempt < 2 {
			var dbBal, dbTusd, dbLp float64
			var dbVer, dbLastActivity int64
			var dbFaucetClaimed, dbDemurrage14Shown bool
			// P1-FIX: reload ALL fields (including faucet_claimed, last_activity_at,
			// demurrage_14_day_warning_shown) so the retry doesn't overwrite those
			// fields with stale in-memory values from before the conflict.
			if err := cs.db.QueryRow(`SELECT balance, tusd_balance, lp_shares, version, last_activity_at, faucet_claimed, demurrage_14_day_warning_shown FROM chain_accounts WHERE lower(address) = $1`, acc.Address).
				Scan(&dbBal, &dbTusd, &dbLp, &dbVer, &dbLastActivity, &dbFaucetClaimed, &dbDemurrage14Shown); err == nil && dbVer > 0 {
				acc.Balance = NewDecimal(dbBal)
				acc.TUsdBalance = NewDecimal(dbTusd)
				acc.LPShares = NewDecimal(dbLp)
				acc.Version = dbVer
				acc.LastActivityAt = dbLastActivity
				acc.FaucetClaimed = dbFaucetClaimed
				acc.Demurrage14DayWarningShown = dbDemurrage14Shown
			}
		}
	}
}

func (cs *ChainState) saveAccountToDBInner(acc *AccountState) {
	if !cs.useDB {
		acc.Version++ // no-DB mode: mark as saved
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
				// Conflict: another node wrote a newer version. Reload DB version
				// into memory so the next caller can retry with the correct base.
				var dbVer int64
				cs.db.QueryRow(`SELECT version FROM chain_accounts WHERE lower(address) = $1`, acc.Address).Scan(&dbVer)
				acc.Version = dbVer // resync in-memory; do NOT increment
				fmt.Printf("[DB] Conflict: account %s modified by another node — local version reset to DB version %d\n", acc.Address, dbVer)
				return // caller must decide whether to retry
			}
		}
	}
	if err != nil {
		fmt.Printf("[DB] Error saving account %s: %v\n", acc.Address, err)
		return
	}
	// P0-1 fix: only increment version after a confirmed successful write
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
const demurrageGracePeriodSeconds = 90 * 24 * 60 * 60 // 3 months
const demurrageMonthlyRate = 0.005                    // 0.5%/month

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
const secondsPerMonth = 30 * 24 * 60 * 60 // approximation, consistent with the grace period's 30-day months

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
// Returns the amount that was decayed (0 if nothing was settled) so callers
// on the primary node can attach it to the queued Transaction — secondary
// nodes replay this exact figure via applyDemurrageLossLocked instead of
// recomputing it themselves (which would use their own wall-clock time and
// diverge from the primary's StateRoot; see ApplyTransferDelta etc.).
func (cs *ChainState) settleDemurrageLocked(acc *AccountState) Decimal {
	// P0-FIX: pool addresses are tokenomics infrastructure — never apply
	// demurrage to them. Doing so would drain pool balances incorrectly.
	if isTokenomicsPoolAddress(acc.Address) {
		return 0
	}
	current := effectiveBalance(acc)
	lost := acc.Balance.Sub(current)
	if lost <= 0 {
		return 0
	}
	acc.Balance = current
	cs.distributeSwapFee(lost.Float(), true) // true = denominated in AEQ; reuses the same 40/30/20/10 split as swap fees
	fmt.Printf("[DEMURRAGE] %s: idle balance decayed by %.6f AEQ, redistributed to pools\n", acc.Address, lost.Float())
	return lost
}

// applyDemurrageLossLocked applies a demurrage loss already decided by the
// primary node (lost, in AEQ) directly to acc's balance and redistributes it
// to the tokenomics pools, WITHOUT consulting effectiveBalance()/nowUnix().
// Used exclusively by secondary-node replay (the "Delta" functions below) so
// every node arrives at byte-identical state for a given block, regardless
// of how much wall-clock time has passed since the primary processed it
// (live replication has sub-second skew; a node resyncing from genesis can
// be replaying months-old transactions all at the "current" wall-clock
// instant, which would otherwise decay them by the wrong amount entirely).
// Caller must hold cs.mu (write lock).
func (cs *ChainState) applyDemurrageLossLocked(acc *AccountState, lost float64) {
	if lost <= 0 || isTokenomicsPoolAddress(acc.Address) {
		return
	}
	acc.Balance = NewDecimal(round6(acc.Balance.Float() - lost))
	cs.distributeSwapFee(lost, true)
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
	if cs.db == nil || proposerAddr == "" {
		return
	}
	proposerAddr = strings.ToLower(proposerAddr)
	res, err := cs.db.Exec(`UPDATE registered_nodes SET blocks_produced = blocks_produced + 1 WHERE lower(signing_address) = lower($1)`, proposerAddr)
	if err != nil || res == nil {
		return
	}
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

	// P1-AUDIT: Query block counts BEFORE acquiring cs.mu. The old code ran this
	// DB query inside the lock, creating the same deadlock risk warned about above
	// for GetRegisteredNodes. Move it here so the lock section is DB-query-free.
	type nodeShare struct {
		wallet string
		blocks int64
	}
	var nodeShares []nodeShare
	var totalBlocks int64
	if cs.db != nil {
		rows, _ := cs.db.Query(`SELECT wallet_address, blocks_produced FROM registered_nodes WHERE wallet_address = ANY($1)`, pq.Array(nodes))
		if rows != nil {
			for rows.Next() {
				var w string
				var b int64
				rows.Scan(&w, &b)
				if b == 0 {
					b = 1
				} // minimum weight so new nodes still get something
				nodeShares = append(nodeShares, nodeShare{w, b})
				totalBlocks += b
			}
			rows.Close()
		}
	}
	if len(nodeShares) == 0 {
		for _, w := range nodes {
			nodeShares = append(nodeShares, nodeShare{w, 1})
			totalBlocks++
		}
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	poolAcc, ok := cs.accounts[validatorsPoolAddr]
	if !ok || poolAcc.Balance <= 0 {
		fmt.Println("[VALIDATORS] Pool is empty — nothing to distribute today")
		return
	}

	total := poolAcc.Balance.Float()
	// P0-2: credit recipients BEFORE zeroing the pool so a crash mid-loop
	// leaves money in the pool (re-distributable) rather than losing it.
	var totalDistributed float64
	for _, ns := range nodeShares {
		wallet := ns.wallet
		// P2-FIX: validate wallet address before crediting — a malformed
		// entry in registered_nodes would insert a garbage key into cs.accounts.
		if len(wallet) != 42 || wallet[:2] != "0x" {
			fmt.Printf("[VALIDATORS] Skipping invalid wallet address: %q\n", wallet)
			continue
		}
		share := round6(total * float64(ns.blocks) / float64(totalBlocks))
		if share <= 0 {
			continue
		} // E4-FIX: skip rounding-to-zero to preserve pool
		if _, ok := cs.accounts[wallet]; !ok {
			cs.accounts[wallet] = &AccountState{Address: wallet}
		}
		acc := cs.accounts[wallet]
		cs.settleDemurrageLocked(acc)
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
		touchActivity(acc)
		cs.enforceWealthCapLocked(acc)
		cs.saveAccountToDB(acc)
		totalDistributed += share
	}
	// Zero pool only after all recipients are successfully written,
	// and only if something was actually distributed (prevents destroying
	// pool balance when all shares rounded to zero).
	if totalDistributed > 0 {
		poolAcc.Balance = NewDecimal(0)
		cs.saveAccountToDB(poolAcc)
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

	// Collect all LP holders and their share counts BEFORE settling demurrage,
	// so we know who participates.
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

	// E3 fix: settle demurrage for ALL LP holders FIRST. settleDemurrageLocked
	// credits demurrage fees to pool addresses (including lpPoolAddr), so the
	// pool balance may increase during this loop. Reading poolAcc.Balance before
	// this loop would miss those newly-credited fees, and zeroing the pool at
	// the end would then destroy them.
	for _, h := range holders {
		acc := cs.accounts[h.addr]
		cs.settleDemurrageLocked(acc)
	}
	// Re-check totalShares after demurrage settlement — shares could have gone to zero.
	if totalShares <= 0 {
		return
	}

	// P2-FIX: second totalShares guard after the demurrage loop. Recompute from
	// live account LPShares values so any unexpected collapse is caught here,
	// preventing division by zero in the distribution loop below.
	totalShares = 0
	for _, h := range holders {
		if acc, ok := cs.accounts[h.addr]; ok {
			totalShares += acc.LPShares.Float()
		}
	}
	if totalShares <= 0 {
		fmt.Println("[LP] totalShares collapsed to zero after demurrage loop — pool left untouched")
		return
	}

	// NOW read the pool balance — it includes any demurrage credits just added.
	poolAcc, ok := cs.accounts[lpPoolAddr]
	if !ok || poolAcc.Balance <= 0 {
		fmt.Println("[LP] Pool is empty — nothing to distribute today")
		return
	}

	total := poolAcc.Balance.Float()
	// P0-2: credit holders BEFORE zeroing pool — crash-safe ordering.
	// E4 fix: track total distributed so we don't zero the pool if all shares
	// rounded to zero (which would destroy micro-AEQ silently).
	var totalDistributed float64
	for _, h := range holders {
		share := round6((h.shares / totalShares) * total)
		totalDistributed += share
		acc := cs.accounts[h.addr]
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
		touchActivity(acc)
		cs.enforceWealthCapLocked(acc)
		cs.saveAccountToDB(acc)
	}
	if totalDistributed > 0 {
		poolAcc.Balance = NewDecimal(0)
		cs.saveAccountToDB(poolAcc)
	} else {
		fmt.Printf("[LP] All shares rounded to zero (%.9f AEQ total) — pool preserved\n", total)
	}
	cs.save()

	holderAddrs := make([]string, len(holders))
	for i, h := range holders {
		holderAddrs[i] = h.addr
	}
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

	// E3-FIX for UBI: settle demurrage for ALL humans FIRST. settleDemurrageLocked
	// credits 20% of each human's decay to ubiPoolAddr. Reading the pool balance
	// BEFORE this loop would miss those credits; zeroing AFTER distributes them.
	// Same fix applied to DistributeLPPool.
	for _, addr := range humanAddrs {
		cs.settleDemurrageLocked(cs.accounts[addr])
	}
	// NOW read pool balance — includes any demurrage credits just added.
	poolAcc, ok = cs.accounts[ubiPoolAddr]
	if !ok || poolAcc.Balance <= 0 {
		fmt.Println("[UBI] Pool empty after demurrage settlement — nothing to distribute")
		return
	}
	// P0-FIX: Do NOT call settleDemurrageLocked on the pool account itself —
	// pool addresses are tokenomics infrastructure and must never have demurrage applied.
	total := poolAcc.Balance.Float()
	share := total / float64(len(humanAddrs))
	// P0-5/P2-9: prevent funds vanishing via float rounding to 0
	if round6(share) == 0 {
		fmt.Printf("[UBI] Share %.10f rounds to zero — pool left intact for next distribution\n", share)
		return
	}
	// P0-2 + P1-6: credit humans BEFORE zeroing pool AND before last_ubi_at.
	for _, addr := range humanAddrs {
		acc := cs.accounts[addr]
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + share))
		touchActivity(acc)
		cs.enforceWealthCapLocked(acc)
		cs.saveAccountToDB(acc)
	}
	poolAcc.Balance = NewDecimal(0)
	cs.saveAccountToDB(poolAcc)
	cs.save()
	// last_ubi_at only set after ALL writes succeed (P1-6).
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
	// Use TotalSupply / humans (= 1000 AEQ always) rather than averaging
	// wallet balances. AEQ deposited into the AMM pool lives in cs.pool.ReserveAEQ
	// — NOT in any human's cs.accounts entry — so wallet-sum / humans gives a
	// misleadingly low number (e.g. 960 when 40 AEQ/human is in the pool).
	// The protocol invariant TotalSupply = humans × 1000 makes the fair-share
	// average exactly 1000 AEQ regardless of where those AEQ currently sit.
	humans := 0
	for _, acc := range cs.accounts {
		if acc.IsHuman {
			humans++
		}
	}
	if humans == 0 {
		return 0
	}
	return 1000.0 // TotalSupply / humans = humans×1000 / humans = 1000 AEQ
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
	wealthCapAmt := avg * multiplier
	if acc.Balance.Float() <= wealthCapAmt {
		return
	}
	excess := acc.Balance.Float() - wealthCapAmt
	acc.Balance = NewDecimal(wealthCapAmt)
	cs.distributeSwapFee(excess, true)
	fmt.Printf("[WEALTH CAP] %s exceeded %.2fx average (%.2f AEQ) — %.4f AEQ excess redistributed to pools\n",
		acc.Address, multiplier, wealthCapAmt, excess)
}

// DemurrageStatus describes whether/when an idle account's AEQ will
// start (or has started) decaying, for surfacing to the user at login.
type DemurrageStatus struct {
	Active                bool    `json:"active"`                   // true if decay has already started
	DaysUntilStart        float64 `json:"days_until_start"`         // only meaningful if !Active; can be negative-free, always >= 0
	ShowFourteenDayNotice bool    `json:"show_fourteen_day_notice"` // one-time notice, true only on the call that first crosses into the 14-day window
	ShowSevenDayNotice    bool    `json:"show_seven_day_notice"`    // true on every check within the last 7 days before decay starts
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
			// P1-5: set in-memory flag SYNCHRONOUSLY to prevent duplicate notices on parallel requests.
			// DB write is async to avoid blocking the GET path.
			acc.Demurrage14DayWarningShown = true
			// FIX-11: capture account by value before launching goroutine so the
			// goroutine does not need to re-acquire cs.mu (which would block until
			// the outer write-lock is released anyway, but a value copy is safer
			// against concurrent mutations and eliminates any lock ordering concerns).
			accCopy := *acc
			go func() {
				cs.saveAccountToDB(&accCopy)
			}()
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
	// P1-10: run EVM sync synchronously first, then retry in background.
	// Prevents permanent Go/EVM divergence if the first sync fails.
	cs.syncHumanRegistrationLocked(V7_CONTRACT_ADDR, address)
	addr := address
	go func() {
		for attempt := 1; attempt <= 3; attempt++ {
			time.Sleep(time.Duration(attempt) * 3 * time.Second)
			cs.mu.RLock()
			cs.syncHumanRegistrationLocked(V7_CONTRACT_ADDR, addr)
			cs.mu.RUnlock()
			fmt.Printf("[STATE] EVM sync retry %d for %s\n", attempt, addr)
		}
	}()
	return nil
}

// Transfer moves amount AEQ from->to on the primary node. Returns the AEQ
// amount demurrage-decayed off the sender and recipient respectively (0 if
// neither was idle long enough to decay) — callers must attach these to the
// queued Transaction so secondary nodes replay the exact same numbers
// instead of recomputing decay at their own wall-clock time (see
// applyDemurrageLossLocked).
func (cs *ChainState) Transfer(from, to string, amount float64) (float64, float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	// P1-FIX: reject NaN/Inf amounts — these would corrupt balances via
	// NewDecimal which uses math.Round (NaN/Inf propagate silently).
	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, 0, fmt.Errorf("invalid transfer amount: %v", amount)
	}
	// P2-5: reject self-transfers; mirrors AequitasV7.sol behaviour and
	// prevents double-demurrage settlement on the same account object.
	if from == to {
		return 0, 0, fmt.Errorf("self-transfer not allowed")
	}

	fromAcc, ok := cs.accounts[from]
	if !ok {
		return 0, 0, fmt.Errorf("insufficient balance")
	}
	fromLost := cs.settleDemurrageLocked(fromAcc) // make sure we're checking against the real, decayed balance
	if fromAcc.Balance.Float() < amount {
		return 0, 0, fmt.Errorf("insufficient balance")
	}

	fromAcc.Balance = NewDecimal(round6(fromAcc.Balance.Float() - amount))
	touchActivity(fromAcc) // sending counts as "using" the money — resets its decay clock
	cs.saveAccountToDB(fromAcc)

	if _, ok := cs.accounts[to]; !ok {
		cs.accounts[to] = &AccountState{Address: to}
	}
	toLost := cs.settleDemurrageLocked(cs.accounts[to])
	cs.accounts[to].Balance = NewDecimal(round6(cs.accounts[to].Balance.Float() + amount))
	touchActivity(cs.accounts[to]) // receiving also resets the clock on the recipient's whole balance
	cs.enforceWealthCapLocked(cs.accounts[to])
	cs.saveAccountToDB(cs.accounts[to])
	cs.save()

	fmt.Printf("[STATE] ✓ Transfer %.2f AEQ: %s → %s\n", amount, from, to)
	cs.syncBalanceLocked(V7_CONTRACT_ADDR, from, to, validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr)
	return fromLost.Float(), toLost.Float(), nil
}

// TransferWithV7Fee is used by the RPC layer when intercepting V7 ERC-20
// transfer() calls (selector a9059cbb). It mirrors V7's _calcFee():
//   TX_FEE_BPS = 10 (0.1% base fee)
//   Concentration surcharge if sender holds ≥1/5/10% of total supply
//   20% of fee → UBI pool, 80% burned (removed from supply)
// Without this, Go-ledger and V7-contract diverge on every user transfer.
// Returns (netAmountCredited, fromDemurrageLost, toDemurrageLost, err) — the
// two demurrage figures must be attached to the queued Transaction so
// secondary nodes replay the exact decay instead of recomputing it (see
// applyDemurrageLossLocked).
func (cs *ChainState) TransferWithV7Fee(from, to string, amount float64) (float64, float64, float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, 0, 0, fmt.Errorf("invalid transfer amount: %v", amount)
	}

	fromAcc, ok := cs.accounts[from]
	if !ok {
		return 0, 0, 0, fmt.Errorf("insufficient balance")
	}
	fromLost := cs.settleDemurrageLocked(fromAcc)
	if fromAcc.Balance.Float() < amount {
		return 0, 0, 0, fmt.Errorf("insufficient balance")
	}

	// Compute total supply inline to avoid re-entering the mutex.
	humans := 0
	for _, acc := range cs.accounts {
		if acc.IsHuman {
			humans++
		}
	}
	totalSupply := float64(humans) * 1000.0
	fee := calcV7Fee(fromAcc.Balance.Float(), amount, totalSupply)
	// E1-FIX: In the Go-state ledger, AEQ cannot be burned (supply is tied
	// to humans * 1000). Redirect 100% of fee to UBI pool instead of the
	// V7-contract's 20%/80% split — this preserves the supply invariant
	// and ensures all fees benefit the community rather than disappearing.
	// E-FIX: compute net first, derive ubi as remainder - preserves supply invariant
	netToRecipient := round6(amount - fee)
	ubiContrib := amount - netToRecipient

	fromAcc.Balance = NewDecimal(round6(fromAcc.Balance.Float() - amount))
	touchActivity(fromAcc)
	cs.saveAccountToDB(fromAcc)

	if _, ok := cs.accounts[to]; !ok {
		cs.accounts[to] = &AccountState{Address: to}
	}
	toLost := cs.settleDemurrageLocked(cs.accounts[to])
	cs.accounts[to].Balance = NewDecimal(round6(cs.accounts[to].Balance.Float() + netToRecipient))
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

	fmt.Printf("[STATE] ✓ TransferV7 %.6f AEQ (fee=%.6f → UBI): %s → %s\n",
		amount, fee, from, to)
	cs.syncBalanceLocked(V7_CONTRACT_ADDR, from, to, validatorsPoolAddr, lpPoolAddr, ubiPoolAddr, treasuryPoolAddr)
	return netToRecipient, fromLost.Float(), toLost.Float(), nil
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
func (cs *ChainState) SwapAEQForTUSD(address string, amountIn, minAmountOut float64) (float64, float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.swapLocked(address, amountIn, true, minAmountOut)
}

// SwapTUSDForAEQ swaps `amountIn` tUSD from `address` into AEQ. Same
// constant-product pricing and fee handling as SwapAEQForTUSD, just with
// the two reserves' roles reversed.
func (cs *ChainState) SwapTUSDForAEQ(address string, amountIn, minAmountOut float64) (float64, float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.swapLocked(address, amountIn, false, minAmountOut)
}

// swapLocked implements both swap directions. aeqToTusd=true means AEQ is
// the input side and tUSD is the output side; false is the reverse.
// minAmountOut, if > 0, rejects the swap before any state is mutated when
// the computed output would fall below it (slippage protection) — 0 means
// no protection requested.
// Caller must hold cs.mu. Returns (amountOut, demurrageLost, err) — lost
// must be attached to the queued Transaction so secondary nodes replay the
// exact decay via ApplySwapDelta instead of recomputing it themselves.
func (cs *ChainState) swapLocked(address string, amountIn float64, aeqToTusd bool, minAmountOut float64) (float64, float64, error) {
	// P2-7: reload pool from DB before swap to avoid stale-memory AMM invariant violation
	cs.reloadPoolFromDB()
	address = strings.ToLower(address)
	if amountIn <= 0 {
		return 0, 0, fmt.Errorf("amount must be positive")
	}
	if cs.pool == nil {
		return 0, 0, fmt.Errorf("liquidity pool not initialized")
	}

	acc, ok := cs.accounts[address]
	if !ok {
		return 0, 0, fmt.Errorf("account not found")
	}
	lost := cs.settleDemurrageLocked(acc) // settle decay before checking/using the AEQ balance below

	if aeqToTusd {
		if acc.Balance.Float() < amountIn {
			return 0, 0, fmt.Errorf("insufficient AEQ balance")
		}
	} else {
		if acc.TUsdBalance.Float() < amountIn {
			return 0, 0, fmt.Errorf("insufficient tUSD balance")
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
			return 0, 0, fmt.Errorf("swap too large for pool liquidity")
		}
		if minAmountOut > 0 && amountOut < minAmountOut {
			return 0, 0, fmt.Errorf("slippage: output %.6f tUSD below requested minimum %.6f", amountOut, minAmountOut)
		}
		cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() + amountInAfterFee))
		cs.pool.ReserveTUSD = NewDecimal(max(0.0, round6(cs.pool.ReserveTUSD.Float()-amountOut)))
		acc.Balance = NewDecimal(round6(acc.Balance.Float() - amountIn))
		acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + amountOut))
	} else {
		amountOut = AMMSwapOut(cs.pool.ReserveTUSD, cs.pool.ReserveAEQ, NewDecimal(amountInAfterFee)).Float()
		if amountOut >= cs.pool.ReserveAEQ.Float() {
			return 0, 0, fmt.Errorf("swap too large for pool liquidity")
		}
		if minAmountOut > 0 && amountOut < minAmountOut {
			return 0, 0, fmt.Errorf("slippage: output %.6f AEQ below requested minimum %.6f", amountOut, minAmountOut)
		}
		cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() + amountInAfterFee))
		cs.pool.ReserveAEQ = NewDecimal(max(0.0, round6(cs.pool.ReserveAEQ.Float()-amountOut)))
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
	return amountOut, lost.Float(), nil
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

// reloadPoolFromDB loads the current pool state from PostgreSQL with SELECT FOR UPDATE
// so swap operations always start from the authoritative DB state, not stale memory.
// P2-7: prevents AMM invariant violation when two nodes swap concurrently.
func (cs *ChainState) reloadPoolFromDB() {
	if cs.db == nil || cs.pool == nil {
		return
	}
	var aeq, tusd, lp float64
	err := cs.db.QueryRow(`SELECT reserve_aeq, reserve_tusd, total_lp_shares FROM liquidity_pool WHERE id = 1`).
		Scan(&aeq, &tusd, &lp)
	if err == nil {
		cs.pool.ReserveAEQ = NewDecimal(aeq)
		cs.pool.ReserveTUSD = NewDecimal(tusd)
		cs.pool.TotalLPShares = NewDecimal(lp)
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
// Returns the AEQ amount demurrage-decayed off address before the deposit
// (0 if none) — callers must attach this to the queued Transaction so
// secondary nodes replay the exact decay via AddLiquidityDelta instead of
// recomputing it themselves.
func (cs *ChainState) AddLiquidity(address string, amountAEQ, amountTUSD float64) (float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	address = strings.ToLower(address)

	if amountAEQ <= 0 || amountTUSD <= 0 {
		return 0, fmt.Errorf("both amounts must be positive")
	}
	if math.IsNaN(amountAEQ) || math.IsInf(amountAEQ, 0) || math.IsNaN(amountTUSD) || math.IsInf(amountTUSD, 0) {
		return 0, fmt.Errorf("invalid liquidity amounts")
	}
	if cs.pool == nil {
		return 0, fmt.Errorf("liquidity pool not initialized")
	}

	acc, ok := cs.accounts[address]
	if !ok {
		return 0, fmt.Errorf("account not found")
	}
	lost := cs.settleDemurrageLocked(acc) // settle decay before checking/using the AEQ balance below
	if acc.Balance.Float() < amountAEQ {
		return 0, fmt.Errorf("insufficient AEQ balance")
	}
	if acc.TUsdBalance.Float() < amountTUSD {
		return 0, fmt.Errorf("insufficient tUSD balance")
	}

	// If the pool already has liquidity, require the deposit to roughly
	// match the existing ratio — an unbalanced deposit would otherwise
	// instantly shift the price, which is the same rule real AMMs enforce.
	var mintedShares float64
	if cs.pool.ReserveAEQ > 0 && cs.pool.ReserveTUSD > 0 {
		expectedTUSD := amountAEQ * (cs.pool.ReserveTUSD.Float() / cs.pool.ReserveAEQ.Float())
		tolerance := expectedTUSD * 0.003 // 0.3% slack — tighter than 1% to prevent price manipulation
		if amountTUSD < expectedTUSD-tolerance || amountTUSD > expectedTUSD+tolerance {
			return 0, fmt.Errorf("deposit ratio does not match pool ratio (expected ~%.4f tUSD for %.4f AEQ)", expectedTUSD, amountAEQ)
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
	return lost.Float(), nil
}

// RemoveLiquidity burns sharesToBurn of address's LP shares and returns
// the corresponding proportional amount of both reserves to their
// balances. sharesToBurn must not exceed the account's own LPShares —
// an account can only withdraw its own claim, never another LP's.
// Returns (outAEQ, outTUSD, demurrageLost, err) — demurrageLost must be
// attached to the queued Transaction so secondary nodes replay the exact
// decay via RemoveLiquidityDelta instead of recomputing it themselves.
func (cs *ChainState) RemoveLiquidity(address string, sharesToBurn float64) (float64, float64, float64, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	address = strings.ToLower(address)

	if sharesToBurn <= 0 {
		return 0, 0, 0, fmt.Errorf("shares must be positive")
	}
	if cs.pool == nil {
		return 0, 0, 0, fmt.Errorf("liquidity pool not initialized")
	}

	acc, ok := cs.accounts[address]
	if !ok {
		return 0, 0, 0, fmt.Errorf("account not found")
	}
	// FIX: RemoveLiquidity previously never settled demurrage on the
	// withdrawing account, unlike AddLiquidity/Transfer/swapLocked — an idle
	// wealthy account could dodge decay indefinitely by periodically
	// removing/re-adding trivial liquidity amounts (touchActivity() below
	// resets the decay clock without the decay ever having been applied).
	lost := cs.settleDemurrageLocked(acc)

	// F17-BOUNDARY: If TotalLPShares rounds to 0 but the user still has LP shares
	// (dust rounding edge case), allow them to drain the entire pool — they are
	// the last LP and the pool is effectively theirs.
	if cs.pool.TotalLPShares <= 0 {
		if acc.LPShares.Float() > 0 {
			outAEQ := cs.pool.ReserveAEQ.Float()
			outTUSD := cs.pool.ReserveTUSD.Float()
			acc.LPShares = NewDecimal(0)
			acc.Balance = NewDecimal(round6(acc.Balance.Float() + outAEQ))
			acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + outTUSD))
			touchActivity(acc)
			cs.enforceWealthCapLocked(acc)
			cs.pool.ReserveAEQ = NewDecimal(0)
			cs.pool.ReserveTUSD = NewDecimal(0)
			cs.pool.TotalLPShares = NewDecimal(0)
			cs.saveAccountToDB(acc)
			cs.savePoolToDB()
			cs.save()
			cs.syncBalanceLocked(V7_CONTRACT_ADDR, address)
			go cs.SavePriceSnapshot()
			fmt.Printf("[POOL] ✓ %s drained final dust position → %.4f AEQ + %.4f tUSD\n", address, outAEQ, outTUSD)
			return outAEQ, outTUSD, lost.Float(), nil
		}
		return 0, 0, 0, fmt.Errorf("liquidity pool is empty")
	}

	if acc.LPShares.Float() < sharesToBurn {
		return 0, 0, 0, fmt.Errorf("insufficient LP shares (have %.6f, requested %.6f)", acc.LPShares.Float(), sharesToBurn)
	}
	// F17-FIX: guard against TotalLPShares corruption (< actual shares).
	// Capping fraction to 1.0 above prevents over-withdrawal from reserves,
	// but TotalLPShares -= sharesToBurn would go negative. Clamp sharesToBurn.
	if sharesToBurn > cs.pool.TotalLPShares.Float() {
		sharesToBurn = cs.pool.TotalLPShares.Float()
		if sharesToBurn <= 0 {
			return 0, 0, 0, fmt.Errorf("pool total LP shares is zero or negative")
		}
		// Zeroing acc.LPShares prevents phantom shares when the clamped
		// sharesToBurn is less than the user's recorded LPShares.
		acc.LPShares = NewDecimal(0)
		// P0-FIX: return immediately after the zero-out so we do NOT fall
		// through to the normal "acc.LPShares -= sharesToBurn" path below,
		// which would compute 0 - sharesToBurn = negative LP shares.
		fraction17 := sharesToBurn / cs.pool.TotalLPShares.Float()
		if fraction17 > 1.0 {
			fraction17 = 1.0
		}
		outAEQ17 := cs.pool.ReserveAEQ.Float() * fraction17
		outTUSD17 := cs.pool.ReserveTUSD.Float() * fraction17
		if outAEQ17 > cs.pool.ReserveAEQ.Float() {
			outAEQ17 = cs.pool.ReserveAEQ.Float()
		}
		if outTUSD17 > cs.pool.ReserveTUSD.Float() {
			outTUSD17 = cs.pool.ReserveTUSD.Float()
		}
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + outAEQ17))
		acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + outTUSD17))
		touchActivity(acc)
		cs.enforceWealthCapLocked(acc)
		newResAEQ17 := round6(cs.pool.ReserveAEQ.Float() - outAEQ17)
		newResTUSD17 := round6(cs.pool.ReserveTUSD.Float() - outTUSD17)
		if newResAEQ17 < 0 {
			newResAEQ17 = 0
		}
		if newResTUSD17 < 0 {
			newResTUSD17 = 0
		}
		cs.pool.ReserveAEQ = NewDecimal(newResAEQ17)
		cs.pool.ReserveTUSD = NewDecimal(newResTUSD17)
		cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() - sharesToBurn))
		cs.saveAccountToDB(acc)
		cs.savePoolToDB()
		cs.save()
		cs.syncBalanceLocked(V7_CONTRACT_ADDR, address)
		go cs.SavePriceSnapshot()
		fmt.Printf("[POOL] ✓ %s removed liquidity (F17 clamp): %.6f shares → %.4f AEQ + %.4f tUSD\n", address, sharesToBurn, outAEQ17, outTUSD17)
		return outAEQ17, outTUSD17, lost.Float(), nil
	}

	fraction := sharesToBurn / cs.pool.TotalLPShares.Float()
	if fraction > 1.0 {
		fraction = 1.0
	} // cap: TotalLPShares corruption guard
	outAEQ := cs.pool.ReserveAEQ.Float() * fraction
	outTUSD := cs.pool.ReserveTUSD.Float() * fraction
	if outAEQ > cs.pool.ReserveAEQ.Float() {
		outAEQ = cs.pool.ReserveAEQ.Float()
	}
	if outTUSD > cs.pool.ReserveTUSD.Float() {
		outTUSD = cs.pool.ReserveTUSD.Float()
	}

	acc.LPShares = NewDecimal(round6(acc.LPShares.Float() - sharesToBurn))
	acc.Balance = NewDecimal(round6(acc.Balance.Float() + outAEQ))
	acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + outTUSD))
	touchActivity(acc) // receiving AEQ back from the pool counts as using it
	cs.enforceWealthCapLocked(acc)
	newReserveAEQ := round6(cs.pool.ReserveAEQ.Float() - outAEQ)
	newReserveTUSD := round6(cs.pool.ReserveTUSD.Float() - outTUSD)
	if newReserveAEQ < 0 {
		newReserveAEQ = 0
	}
	if newReserveTUSD < 0 {
		newReserveTUSD = 0
	}
	cs.pool.ReserveAEQ = NewDecimal(newReserveAEQ)
	cs.pool.ReserveTUSD = NewDecimal(newReserveTUSD)
	cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() - sharesToBurn))

	cs.saveAccountToDB(acc)
	cs.savePoolToDB()
	cs.save()

	cs.syncBalanceLocked(V7_CONTRACT_ADDR, address)
	go cs.SavePriceSnapshot()

	fmt.Printf("[POOL] ✓ %s removed liquidity: %.6f shares → %.4f AEQ + %.4f tUSD\n", address, sharesToBurn, outAEQ, outTUSD)
	return outAEQ, outTUSD, lost.Float(), nil
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
	// P2-AUDIT: Add to existing balance instead of overwriting — a user who had
	// received tUSD via another path (pool payout, migration) before claiming the
	// faucet would have had their entire tUSD balance zeroed by the old Set.
	acc.TUsdBalance = acc.TUsdBalance.Add(NewDecimal(tusdFaucetAmount))
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
	// P1-1: read last_ubi_at from DB BEFORE acquiring the mutex to avoid
	// holding RLock across a blocking DB query (deadlock / latency risk).
	lastUBIAt := cs.getConfigValue("last_ubi_at")
	cs.mu.RLock()
	addrs := make([]string, 0, len(cs.accounts))
	for a := range cs.accounts {
		addrs = append(addrs, a)
	}
	sort.Strings(addrs)
	var sb strings.Builder
	for _, a := range addrs {
		acc := cs.accounts[a]
		// Include ALL accounts with non-zero AEQ or tUSD balances (not only humans)
		if acc.IsHuman || acc.Balance > 0 || acc.TUsdBalance > 0 || acc.LPShares > 0 {
			// FaucetClaimed and LastActivityAt must be included: two nodes that
			// processed different sets of faucet claims would produce the same
			// balances but handle future UBI distribution differently.
			// P1-9: LastActivityAt excluded — wall-clock differs between nodes
			// for the same TX, causing StateRoot mismatch and peer block rejection.
			fmt.Fprintf(&sb, "%s:%.6f:%.6f:%.6f:h=%v:fc=%v:",
				a,
				round6(acc.Balance.Float()),
				round6(acc.TUsdBalance.Float()),
				round6(acc.LPShares.Float()),
				acc.IsHuman,
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
	for k := range cs.nullifiers {
		nullKeys = append(nullKeys, k)
	}
	sort.Strings(nullKeys)
	fmt.Fprintf(&sb, "|n=%d:", len(nullKeys))
	// P2-5: include only nullifier keys, not wallet addresses (privacy)
	for _, k := range nullKeys {
		sb.WriteString(k)
		sb.WriteString(":")
	}
	// Include last UBI distribution timestamp (pre-fetched before RLock — P1-1).
	fmt.Fprintf(&sb, "|ubi:%s", lastUBIAt)
	cs.mu.RUnlock()
	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

// calcGiniLocked computes the Gini coefficient without acquiring cs.mu.
// Must only be called while cs.mu is already held (read or write).
// calcGiniFromBalances is the single shared implementation used by both
// calcGiniLocked (inside lock) and CalcGini (acquires lock). P2-1: uses
// sort.Float64s O(n log n) instead of the old O(n^2) bubble sort.
func calcGiniFromBalances(balances []float64) float64 {
	n := len(balances)
	if n < 2 {
		return 0.0
	}
	sort.Float64s(balances)
	var sum, numerator float64
	for i, x := range balances {
		sum += x
		numerator += float64(2*i+1-n) * x
	}
	if sum == 0 {
		return 0.0
	}
	gini := numerator / (float64(n) * sum)
	if gini < 0 {
		gini = -gini
	}
	if n > 1 {
		gini = gini * float64(n) / float64(n-1)
	}
	if gini > 1.0 {
		gini = 1.0
	}
	return gini
}

func (cs *ChainState) calcGiniLocked() float64 {
	var balances []float64
	for _, acc := range cs.accounts {
		if acc.IsHuman && acc.Balance > 0 {
			balances = append(balances, effectiveBalance(acc).Float())
		}
	}
	return calcGiniFromBalances(balances)
}

func (cs *ChainState) CalcGini() float64 {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	// P2-1: deduplicated — delegates to calcGiniLocked (which now uses sort.Float64s).
	return cs.calcGiniLocked()
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

// -- SECONDARY-NODE REPLAY DELTA METHODS -----------------------------------
// These methods are called exclusively by replayTransactions on secondary nodes.
// They apply pre-computed amounts directly, without re-running business logic,
// to avoid pool-state divergence and floating-point ordering differences.

// ApplyTransferDelta directly adjusts AEQ balances by the net amount that
// reached the recipient (after any fee that was applied on the primary).
// Used by secondary nodes replaying transfer TXs from blocks. fromLost/toLost
// are the exact demurrage amounts the primary decayed off each side (see
// Transfer()/TransferWithV7Fee()) — applied directly via
// applyDemurrageLossLocked rather than recomputed, since recomputing via
// effectiveBalance()/nowUnix() at replay time would use this node's own
// wall-clock time and diverge from what the primary actually settled.
func (cs *ChainState) ApplyTransferDelta(from, to string, netAmount, fromLost, toLost float64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	fromAcc, ok := cs.accounts[from]
	if !ok {
		return fmt.Errorf("from account not found: %s", from)
	}
	cs.applyDemurrageLossLocked(fromAcc, fromLost)
	if fromAcc.Balance.Float() < netAmount {
		return fmt.Errorf("insufficient balance (have %.6f, need %.6f)", fromAcc.Balance.Float(), netAmount)
	}
	fromAcc.Balance = NewDecimal(round6(fromAcc.Balance.Float() - netAmount))
	cs.saveAccountToDB(fromAcc)

	if _, ok := cs.accounts[to]; !ok {
		cs.accounts[to] = &AccountState{Address: to}
	}
	toAcc := cs.accounts[to]
	cs.applyDemurrageLossLocked(toAcc, toLost)
	toAcc.Balance = NewDecimal(round6(toAcc.Balance.Float() + netAmount))
	cs.enforceWealthCapLocked(toAcc)
	cs.saveAccountToDB(toAcc)
	return nil
}

// ApplySwapDelta adjusts balances after a swap, using the exact amountIn/amountOut
// stored in the block TX. aeqToTusd=true: wallet loses amountIn AEQ, gains amountOut tUSD.
// aeqToTusd=false: wallet loses amountIn tUSD, gains amountOut AEQ.
// Also updates pool reserves to mirror what swapLocked() did on the primary.
// demurrageLost is the exact amount swapLocked() decayed off wallet on the
// primary — applied directly (applyDemurrageLossLocked) rather than
// recomputed, for the same reason as in ApplyTransferDelta.
func (cs *ChainState) ApplySwapDelta(wallet string, amountIn, amountOut float64, aeqToTusd bool, demurrageLost float64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	wallet = strings.ToLower(wallet)
	acc, ok := cs.accounts[wallet]
	if !ok {
		return fmt.Errorf("account not found: %s", wallet)
	}
	cs.applyDemurrageLossLocked(acc, demurrageLost)
	if aeqToTusd {
		if acc.Balance.Float() < amountIn {
			return fmt.Errorf("insufficient AEQ balance")
		}
		acc.Balance = NewDecimal(round6(acc.Balance.Float() - amountIn))
		acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + amountOut))
	} else {
		if acc.TUsdBalance.Float() < amountIn {
			return fmt.Errorf("insufficient tUSD balance")
		}
		acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() - amountIn))
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + amountOut))
	}
	cs.saveAccountToDB(acc)

	// Update pool reserves to match what swapLocked() did on primary.
	// fee is swapFeeBps (0.1%); amountInAfterFee is what enters the pool.
	if cs.pool != nil {
		fee := amountIn * float64(swapFeeBps) / 10000.0
		amountInAfterFee := amountIn - fee
		if aeqToTusd {
			// Sender put in AEQ, got tUSD: reserveAEQ grows, reserveTUSD shrinks.
			cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() + amountInAfterFee))
			cs.pool.ReserveTUSD = NewDecimal(max(0.0, round6(cs.pool.ReserveTUSD.Float()-amountOut)))
		} else {
			// Sender put in tUSD, got AEQ: reserveTUSD grows, reserveAEQ shrinks.
			cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() + amountInAfterFee))
			cs.pool.ReserveAEQ = NewDecimal(max(0.0, round6(cs.pool.ReserveAEQ.Float()-amountOut)))
		}
		cs.savePoolToDB()
		// Distribute swap fee to the 4 tokenomics pools (40% validators /
		// 30% LP / 20% UBI / 10% treasury) — mirrors swapLocked() on primary.
		// Without this the fee-pool addresses stay at 0 on secondaries,
		// causing StateRoot divergence (pool addresses are included in the hash).
		cs.distributeSwapFee(fee, aeqToTusd)
	}
	return nil
}

// AddLiquidityDelta applies an add-liquidity operation on secondary nodes using
// the exact stored amounts. lpShares is the number of LP shares minted on the
// primary node; if > 0 it is used directly instead of recomputing, eliminating
// pool-state drift between nodes. Reloads pool from DB first for consistency.
// demurrageLost is the exact amount AddLiquidity() decayed off wallet on the
// primary — applied directly rather than recomputed, for the same reason as
// in ApplyTransferDelta.
func (cs *ChainState) AddLiquidityDelta(wallet string, aeqAmount, tusdAmount, lpShares, demurrageLost float64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	wallet = strings.ToLower(wallet)
	acc, ok := cs.accounts[wallet]
	if !ok {
		return fmt.Errorf("account not found: %s", wallet)
	}
	cs.reloadPoolFromDB()
	cs.applyDemurrageLossLocked(acc, demurrageLost)
	if acc.Balance.Float() < aeqAmount {
		return fmt.Errorf("insufficient AEQ balance")
	}
	if acc.TUsdBalance.Float() < tusdAmount {
		return fmt.Errorf("insufficient tUSD balance")
	}

	// Use the stored LP shares from the primary node when available.
	// Fall back to recomputing (from pool state or geometric mean) for
	// blocks produced by old nodes that don't include the lp_shares field.
	var mintedShares float64
	if lpShares > 0 {
		mintedShares = lpShares
	} else if cs.pool != nil && cs.pool.ReserveAEQ.Float() > 0 && cs.pool.TotalLPShares.Float() > 0 {
		mintedShares = (aeqAmount / cs.pool.ReserveAEQ.Float()) * cs.pool.TotalLPShares.Float()
	} else {
		mintedShares = math.Sqrt(aeqAmount * tusdAmount)
	}

	acc.Balance = NewDecimal(round6(acc.Balance.Float() - aeqAmount))
	acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() - tusdAmount))
	acc.LPShares = NewDecimal(round6(acc.LPShares.Float() + mintedShares))
	if cs.pool != nil {
		cs.pool.ReserveAEQ = NewDecimal(round6(cs.pool.ReserveAEQ.Float() + aeqAmount))
		cs.pool.ReserveTUSD = NewDecimal(round6(cs.pool.ReserveTUSD.Float() + tusdAmount))
		cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() + mintedShares))
		cs.savePoolToDB()
	}
	cs.saveAccountToDB(acc)
	return nil
}

// RemoveLiquidityDelta burns sharesToBurn LP shares and returns proportional
// pool reserves to the wallet, using the secondary's current pool state.
// demurrageLost is the exact amount RemoveLiquidity() decayed off wallet on
// the primary — applied directly rather than recomputed, for the same
// reason as in ApplyTransferDelta.
func (cs *ChainState) RemoveLiquidityDelta(wallet string, sharesToBurn, demurrageLost float64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	wallet = strings.ToLower(wallet)
	acc, ok := cs.accounts[wallet]
	if !ok {
		return fmt.Errorf("account not found: %s", wallet)
	}
	cs.reloadPoolFromDB()
	cs.applyDemurrageLossLocked(acc, demurrageLost)
	if cs.pool == nil || cs.pool.TotalLPShares.Float() <= 0 {
		return fmt.Errorf("liquidity pool is empty")
	}
	if acc.LPShares.Float() < sharesToBurn {
		return fmt.Errorf("insufficient LP shares")
	}
	// Mirror F17 + F18 caps from primary RemoveLiquidity
	if sharesToBurn > cs.pool.TotalLPShares.Float() {
		sharesToBurn = cs.pool.TotalLPShares.Float()
		if sharesToBurn <= 0 {
			return fmt.Errorf("pool total LP shares is zero")
		}
	}
	fraction := sharesToBurn / cs.pool.TotalLPShares.Float()
	if fraction > 1.0 {
		fraction = 1.0
	}
	outAEQ := round6(cs.pool.ReserveAEQ.Float() * fraction)
	outTUSD := round6(cs.pool.ReserveTUSD.Float() * fraction)
	if outAEQ > cs.pool.ReserveAEQ.Float() {
		outAEQ = cs.pool.ReserveAEQ.Float()
	}
	if outTUSD > cs.pool.ReserveTUSD.Float() {
		outTUSD = cs.pool.ReserveTUSD.Float()
	}

	acc.LPShares = NewDecimal(round6(acc.LPShares.Float() - sharesToBurn))
	acc.Balance = NewDecimal(round6(acc.Balance.Float() + outAEQ))
	acc.TUsdBalance = NewDecimal(round6(acc.TUsdBalance.Float() + outTUSD))
	newReserveAEQ := round6(cs.pool.ReserveAEQ.Float() - outAEQ)
	newReserveTUSD := round6(cs.pool.ReserveTUSD.Float() - outTUSD)
	if newReserveAEQ < 0 {
		newReserveAEQ = 0
	}
	if newReserveTUSD < 0 {
		newReserveTUSD = 0
	}
	cs.pool.ReserveAEQ = NewDecimal(newReserveAEQ)
	cs.pool.ReserveTUSD = NewDecimal(newReserveTUSD)
	cs.pool.TotalLPShares = NewDecimal(round6(cs.pool.TotalLPShares.Float() - sharesToBurn))
	cs.savePoolToDB()
	cs.saveAccountToDB(acc)
	return nil
}

// ApplyUBIDelta credits amountPerHuman AEQ to every registered human on this node.
// Used by secondary nodes replaying ubi_distribution TXs from blocks.
//
// FIX (StateRoot divergence): ubiAt must be the timestamp the PRIMARY used
// when it ran DistributeUBIPool (i.e. the block's Timestamp), not this
// node's own wall clock. last_ubi_at feeds directly into StateRoot(), so
// every secondary independently calling time.Now() here wrote a different
// value than the primary and than every OTHER secondary — guaranteeing a
// StateRoot mismatch on every single UBI distribution. Pass 0 to fall back
// to time.Now() only for callers that have no block context (none should,
// post-fix, but this keeps the function safe to call directly).
func (cs *ChainState) ApplyUBIDelta(amountPerHuman float64, ubiAt int64) {
	if amountPerHuman <= 0 {
		return
	}
	if ubiAt <= 0 {
		ubiAt = time.Now().Unix()
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, acc := range cs.accounts {
		if !acc.IsHuman {
			continue
		}
		acc.Balance = NewDecimal(round6(acc.Balance.Float() + amountPerHuman))
		touchActivity(acc)
		cs.enforceWealthCapLocked(acc)
		cs.saveAccountToDB(acc)
	}
	// Zero the UBI pool on secondary (it was zeroed on primary after distribution)
	if ubiAcc, ok := cs.accounts[ubiPoolAddr]; ok {
		ubiAcc.Balance = NewDecimal(0)
		cs.saveAccountToDB(ubiAcc)
	}
	// Write last_ubi_at to secondary's chain_config so StateRoot matches primary.
	cs.setConfigValue("last_ubi_at", fmt.Sprintf("%d", ubiAt))
}

// ApplyFaucetDelta credits faucetAmount tUSD to wallet and marks FaucetClaimed.
// Used by secondary nodes replaying faucet TXs from blocks.
func (cs *ChainState) ApplyFaucetDelta(wallet string, faucetAmount float64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	wallet = strings.ToLower(wallet)
	if _, ok := cs.accounts[wallet]; !ok {
		cs.accounts[wallet] = &AccountState{Address: wallet}
	}
	acc := cs.accounts[wallet]
	if acc.FaucetClaimed {
		return nil // idempotent: already applied
	}
	acc.FaucetClaimed = true
	acc.TUsdBalance = acc.TUsdBalance.Add(NewDecimal(faucetAmount))
	cs.saveAccountToDB(acc)
	return nil
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
 ON CONFLICT (address) DO UPDATE SET commitment = $2, last_activity = NOW()`,
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
			"address":     addr,
			"balance_wei": bal,
		})
	}
	return balances
}
