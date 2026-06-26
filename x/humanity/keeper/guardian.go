package keeper

// Guardian system + inactive-wallet escrow implementation.
//
// A Guardian is a trusted registered human who can confirm another human is
// still alive, preventing their AEQ from moving to escrow due to inactivity.
//
// Timeline (from Whitepaper):
//   Year 0–2    Normal use
//   Year 2      Warning 1 (guardian can respond)
//   Year 2+60d  Warning 2
//   Year 2+120d Warning 3
//   Year 2+180d Balance → ESCROW (recoverable by owner)
//   Year 4      If still inactive → escrow released to UBI pool
//
// The "2 years + 180 days" threshold before escrow = 2.5 years = 913 days.
// The "still inactive at 4 years" threshold from registration means from the
// moment funds entered escrow (2.5 yr mark) there is another 1.5 years before
// the UBI release. We track moved_at in escrow_accounts for this.

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// inactivityEscrowSeconds is the inactivity threshold before funds move to escrow.
// 2 years + 180 days = 2.5 years ≈ 913 days.
const inactivityEscrowSeconds = int64((2*365 + 180) * 24 * 60 * 60)

// escrowToUBISeconds is how long funds sit in escrow before moving to the UBI
// pool if the owner stays inactive. 1.5 years ≈ 548 days.
const escrowToUBISeconds = int64((365 + 183) * 24 * 60 * 60)

// guardianTimelockSeconds is the 7-day timelock preventing guardian changes.
const guardianTimelockSeconds = int64(7 * 24 * 60 * 60)

// maxWardsPerGuardian limits how many wards a single guardian may have.
const maxWardsPerGuardian = 3

// ─── DB SCHEMA ────────────────────────────────────────────────────────────────

// InitGuardianTables creates the guardians and escrow_accounts tables if they
// don't already exist. Called from initDB so they're always present.
func (cs *ChainState) InitGuardianTables() error {
	if cs.db == nil {
		return nil
	}
	if _, err := cs.db.Exec(`CREATE TABLE IF NOT EXISTS guardians (
		wallet_address  TEXT PRIMARY KEY,
		guardian_address TEXT NOT NULL,
		set_at          BIGINT NOT NULL
	)`); err != nil {
		return fmt.Errorf("failed to create guardians table: %w", err)
	}
	if _, err := cs.db.Exec(`CREATE TABLE IF NOT EXISTS escrow_accounts (
		wallet_address TEXT PRIMARY KEY,
		amount         NUMERIC NOT NULL,
		moved_at       BIGINT NOT NULL
	)`); err != nil {
		return fmt.Errorf("failed to create escrow_accounts table: %w", err)
	}
	return nil
}

// ─── GUARDIAN STATE METHODS ───────────────────────────────────────────────────

// SetGuardian persists a guardian relationship. Validates:
//   - both wallet and guardian are registered humans
//   - no circular relationship
//   - guardian does not already have maxWardsPerGuardian wards
//   - 7-day timelock since last guardian change for this wallet
//
// The caller (API handler) is responsible for signature verification.
func (cs *ChainState) SetGuardian(wallet, guardian string) error {
	// Always use server time for the timelock — ignoring the caller-supplied
	// timestamp prevents a caller from passing a future value to bypass the lock.
	setAt := time.Now().Unix()
	wallet = strings.ToLower(wallet)
	guardian = strings.ToLower(guardian)

	if wallet == guardian {
		return fmt.Errorf("wallet cannot be its own guardian")
	}

	// Both must be registered humans.
	if !cs.IsHuman(wallet) {
		return fmt.Errorf("wallet %s is not a registered human", wallet)
	}
	if !cs.IsHuman(guardian) {
		return fmt.Errorf("guardian %s is not a registered human", guardian)
	}

	if cs.db == nil {
		return fmt.Errorf("guardian operations require a database — run with DATABASE_URL set")
	}

	// Anti-circular: A cannot be guardian of B if B is guardian of A.
	// Check outside the transaction (read-only, no TOCTOU risk for this check).
	var guardianOfGuardian string
	scanErr := cs.db.QueryRow(
		`SELECT guardian_address FROM guardians WHERE wallet_address = $1`, guardian,
	).Scan(&guardianOfGuardian)
	if scanErr == nil && strings.ToLower(guardianOfGuardian) == wallet {
		return fmt.Errorf("circular guardian relationship: %s is already guardian of %s", wallet, guardian)
	}

	// Max 3 wards per guardian (pre-check; re-checked under lock in confirmGuardian).
	var wardCount int
	cs.db.QueryRow(
		`SELECT COUNT(*) FROM guardians WHERE lower(guardian_address) = lower($1) AND lower(wallet_address) != lower($2)`,
		guardian, wallet,
	).Scan(&wardCount)
	if wardCount >= maxWardsPerGuardian {
		return fmt.Errorf("guardian %s already has %d wards (maximum %d)", guardian, wardCount, maxWardsPerGuardian)
	}

	// FIX 4: Wrap the 7-day timelock check + update in a single transaction with
	// SELECT FOR UPDATE to prevent concurrent requests from racing past the timelock.
	tx, err := cs.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin guardian transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var existingSetAt int64
	rowErr := tx.QueryRow(
		`SELECT COALESCE(set_at, 0) FROM guardians WHERE wallet_address = $1 FOR UPDATE`,
		wallet,
	).Scan(&existingSetAt)
	if rowErr != nil && rowErr != sql.ErrNoRows {
		return fmt.Errorf("guardian lookup failed: %w", rowErr)
	}

	if existingSetAt > 0 && setAt-existingSetAt < guardianTimelockSeconds {
		daysLeft := (guardianTimelockSeconds - (setAt - existingSetAt)) / 86400
		return fmt.Errorf("guardian was set %d days ago — must wait 7 days before changing (%.0f days remaining)",
			(setAt-existingSetAt)/86400, float64(daysLeft))
	}

	// FIX 5: Re-check ward count inside the transaction (under the lock) to
	// prevent a TOCTOU race where two concurrent requests both pass the pre-check
	// and then both increment the ward count past the maximum.
	var currentWards int
	tx.QueryRow( //nolint:errcheck
		`SELECT COUNT(*) FROM guardians WHERE lower(guardian_address) = lower($1) AND lower(wallet_address) != lower($2)`,
		guardian, wallet,
	).Scan(&currentWards)
	if currentWards >= maxWardsPerGuardian {
		return fmt.Errorf("guardian %s already has %d wards (maximum %d) — re-check under lock", guardian, currentWards, maxWardsPerGuardian)
	}

	// FIX 1: Anti-cycle depth-limited walk. The simple A→B/B→A check above only
	// catches 2-node cycles. Walk up to 5 levels of the guardian chain starting
	// from `guardian` to detect 3+ node cycles (A→B→C→A etc.) before inserting.
	current := guardian
	for depth := 0; depth < 5; depth++ {
		var nextGuardian string
		err := tx.QueryRow(
			`SELECT guardian_address FROM guardians WHERE lower(wallet_address) = lower($1)`, current,
		).Scan(&nextGuardian)
		if err != nil {
			break // no guardian for `current` — end of chain, safe to insert
		}
		if strings.EqualFold(nextGuardian, wallet) {
			tx.Rollback() //nolint:errcheck
			return fmt.Errorf("circular guardian relationship detected (cycle depth %d)", depth+2)
		}
		current = nextGuardian
	}

	if _, execErr := tx.Exec(
		`INSERT INTO guardians (wallet_address, guardian_address, set_at)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (wallet_address) DO UPDATE
		   SET guardian_address = $2, set_at = $3`,
		wallet, guardian, setAt,
	); execErr != nil {
		return fmt.Errorf("db error: %w", execErr)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("guardian commit failed: %w", err)
	}

	fmt.Printf("[GUARDIAN] ✓ %s set guardian to %s\n", wallet, guardian)
	return nil
}

// GetGuardian returns the guardian address and the timestamp it was set, or
// ("", 0, sql.ErrNoRows) if no guardian is configured for wallet.
func (cs *ChainState) GetGuardian(wallet string) (guardian string, setAt int64, err error) {
	wallet = strings.ToLower(wallet)
	if cs.db == nil {
		return "", 0, fmt.Errorf("no database")
	}
	row := cs.db.QueryRow(
		`SELECT guardian_address, set_at FROM guardians WHERE wallet_address = $1`, wallet,
	)
	err = row.Scan(&guardian, &setAt)
	return guardian, setAt, err
}

// ConfirmAlive resets wallet's last_activity_at to now. The caller must be the
// guardian of wallet; the API handler verifies this before calling here.
// The guardian has ZERO financial access — this only resets the inactivity timer.
//
// expectedGuardian is the guardian address the API handler looked up before
// acquiring this lock. We re-fetch and compare inside the lock to close the
// TOCTOU window where a concurrent SetGuardian could swap the guardian between
// the API handler's DB lookup and this update.
func (cs *ChainState) ConfirmAlive(wallet, expectedGuardian string) error {
	wallet = strings.ToLower(wallet)
	expectedGuardian = strings.ToLower(expectedGuardian)
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Re-verify the guardian inside the lock to prevent TOCTOU: a concurrent
	// SetGuardian could change the guardian between the API handler's DB lookup
	// and here. Re-fetching under the write lock ensures atomicity.
	if cs.db != nil && expectedGuardian != "" {
		var currentGuardian string
		dbErr := cs.db.QueryRow(
			`SELECT guardian_address FROM guardians WHERE wallet_address = $1`, wallet,
		).Scan(&currentGuardian)
		if dbErr != nil || strings.ToLower(currentGuardian) != expectedGuardian {
			return fmt.Errorf("guardian mismatch — concurrent guardian change detected, retry")
		}
	}

	acc, ok := cs.accounts[wallet]
	if !ok {
		return fmt.Errorf("account %s not found", wallet)
	}
	// FIX 3: Only registered humans participate in the guardian / inactivity system.
	// Non-human accounts have no inactivity clock, so resetting it is a no-op at
	// best and a logic error at worst (could mask a misconfigured guardian entry).
	if !acc.IsHuman {
		return fmt.Errorf("account %s is not a registered human — guardian confirm-alive not applicable", wallet)
	}
	touchActivity(acc)
	cs.saveAccountToDB(acc)
	fmt.Printf("[GUARDIAN] ✓ Guardian confirmed %s is alive — activity timer reset\n", wallet)
	return nil
}

// GetEscrow returns the escrow amount and moved_at timestamp for wallet, or
// (0, 0, nil) if no escrow entry exists.
func (cs *ChainState) GetEscrow(wallet string) (amount float64, movedAt int64, err error) {
	wallet = strings.ToLower(wallet)
	if cs.db == nil {
		return 0, 0, nil
	}
	row := cs.db.QueryRow(
		`SELECT amount, moved_at FROM escrow_accounts WHERE wallet_address = $1`, wallet,
	)
	err = row.Scan(&amount, &movedAt)
	if err == sql.ErrNoRows {
		return 0, 0, nil // no escrow entry — not an error
	} else if err != nil {
		return 0, 0, fmt.Errorf("escrow lookup failed: %w", err)
	}
	return amount, movedAt, nil
}

// RecoverFromEscrow lets a wallet owner reclaim their escrowed balance.
// The caller (API handler) is responsible for signature verification.
// On success the escrow row is deleted and the balance restored to the wallet.
func (cs *ChainState) RecoverFromEscrow(wallet string) error {
	wallet = strings.ToLower(wallet)
	if cs.db == nil {
		return fmt.Errorf("no database")
	}

	// FIX 8: Use DELETE...RETURNING for an atomic claim. A plain SELECT then
	// DELETE has a TOCTOU race: two concurrent recovery calls can both pass the
	// SELECT and then both credit the balance. The atomic DELETE returns the row
	// only to the goroutine that deletes it — the other gets sql.ErrNoRows.
	var amount float64
	var movedAt int64
	err := cs.db.QueryRow(
		`DELETE FROM escrow_accounts WHERE wallet_address = $1 RETURNING amount, moved_at`,
		wallet,
	).Scan(&amount, &movedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no escrow found for wallet %s", wallet)
	}
	if err != nil {
		return fmt.Errorf("escrow retrieval failed: %w", err)
	}
	if amount <= 0 {
		return fmt.Errorf("escrow amount is zero")
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Credit balance back. If the account was lost from memory (e.g. after a
	// state reset), recreate it as a human — escrow only exists for registered
	// humans, so restoring it as non-human would break the supply invariant.
	if _, ok := cs.accounts[wallet]; !ok {
		cs.accounts[wallet] = &AccountState{Address: wallet, IsHuman: true}
	}
	acc := cs.accounts[wallet]
	cs.settleDemurrageLocked(acc)
	acc.Balance = NewDecimal(round6(acc.Balance.Float() + amount))
	touchActivity(acc) // recovering from escrow resets the inactivity clock
	cs.enforceWealthCapLocked(acc)
	cs.saveAccountToDB(acc)

	fmt.Printf("[ESCROW] ✓ %s recovered %.6f AEQ from escrow\n", wallet, amount)
	return nil
}

// ─── DAILY SCHEDULER METHODS ─────────────────────────────────────────────────

// CheckAndMoveToEscrow is called once per day. It finds every wallet whose
// last_activity_at is older than inactivityEscrowSeconds and moves their AEQ
// balance to the escrow_accounts table. The balance is removed from the
// wallet immediately but remains recoverable by the owner for up to 1.5 more
// years (see ReleaseEscrowToUBI).
func (cs *ChainState) CheckAndMoveToEscrow() {
	if cs.db == nil {
		return
	}
	threshold := time.Now().Unix() - inactivityEscrowSeconds
	now := time.Now().Unix()

	// Phase 1: collect candidates from in-memory state under RLock (no DB calls).
	type candidate struct {
		addr    string
		balance float64
	}
	var preCandiates []string
	cs.mu.RLock()
	for addr, acc := range cs.accounts {
		if !acc.IsHuman { continue }
		if acc.LastActivityAt == 0 || acc.LastActivityAt > threshold { continue }
		bal := effectiveBalance(acc).Float()
		if bal <= 0 { continue }
		preCandiates = append(preCandiates, addr)
	}
	cs.mu.RUnlock()

	if len(preCandiates) == 0 {
		return
	}

	// Phase 2: filter out already-escrowed wallets via DB (outside RLock).
	type candidateBalance struct{ addr string; balance float64 }
	var candidates []candidateBalance
	for _, addr := range preCandiates {
		var existing float64
		scanErr := cs.db.QueryRow(
			`SELECT amount FROM escrow_accounts WHERE wallet_address = $1`, addr,
		).Scan(&existing)
		if scanErr == nil && existing > 0 {
			continue // already escrowed
		}
		// Re-read balance under no lock (best-effort; confirmed under write lock below).
		cs.mu.RLock()
		acc, ok := cs.accounts[addr]
		var bal float64
		if ok { bal = effectiveBalance(acc).Float() }
		cs.mu.RUnlock()
		if bal <= 0 { continue }
		candidates = append(candidates, candidateBalance{addr, bal})
	}

	if len(candidates) == 0 {
		return
	}

	// FIX 1: Collect decisions under write lock; do DB writes OUTSIDE the lock
	// to avoid blocking all chain operations during potentially slow DB I/O.
	type escrowEntry struct {
		acc            AccountState // value copy — safe to use after lock released
		balance        float64
		now            int64
		inactiveSince  string
	}
	var toEscrow []escrowEntry

	cs.mu.Lock()
	for _, c := range candidates {
		acc, ok := cs.accounts[c.addr]
		if !ok {
			continue
		}
		// FIX 5: Re-check IsHuman under write lock. Between the candidate scan
		// and here an account could theoretically have been de-registered or was
		// never human (e.g. a pool address that slipped through the RLock scan).
		// Non-human accounts must never be escrowed.
		if !acc.IsHuman {
			continue
		}
		// Re-check LastActivityAt under write lock — between the candidate scan
		// and here, the account may have been touched by a Transfer/Swap/ConfirmAlive
		// call, which would update LastActivityAt.
		if acc.LastActivityAt > threshold {
			continue // no longer inactive
		}
		bal := round6(effectiveBalance(acc).Float())
		if bal <= 0 {
			continue
		}

		// Settle demurrage first so Balance reflects reality.
		cs.settleDemurrageLocked(acc)
		bal = round6(acc.Balance.Float())
		if bal <= 0 {
			continue
		}

		inactiveSince := time.Unix(acc.LastActivityAt, 0).Format("2006-01-02")

		// Update in-memory state NOW (under the lock).
		acc.Balance = NewDecimal(0)

		// Collect a value copy for DB writes after the lock is released.
		toEscrow = append(toEscrow, escrowEntry{
			acc:           *acc,
			balance:       bal,
			now:           now,
			inactiveSince: inactiveSince,
		})
	}
	cs.mu.Unlock()

	// FIX 1 (cont.) + FIX 2: DB writes outside the lock.
	// FIX 2: Use ON CONFLICT DO NOTHING so an existing escrow row's moved_at
	// clock is never reset — resetting it would restart the 1.5-year countdown.
	for _, entry := range toEscrow {
		if _, err := cs.db.Exec(
			`INSERT INTO escrow_accounts (wallet_address, amount, moved_at)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (wallet_address) DO NOTHING`,
			entry.acc.Address, entry.balance, entry.now,
		); err != nil {
			fmt.Printf("[ESCROW] Error writing escrow for %s: %v\n", entry.acc.Address, err)
			continue
		}
		cs.saveAccountToDB(&entry.acc)
		fmt.Printf("[ESCROW] ✓ Moved %.6f AEQ from %s to escrow (inactive since %s)\n",
			entry.balance, entry.acc.Address, entry.inactiveSince)
	}
}

// ReleaseEscrowToUBI is called once per day. It finds escrow entries older
// than escrowToUBISeconds (1.5 years from when the funds were escrowed) and
// moves them into the UBI pool for distribution.
func (cs *ChainState) ReleaseEscrowToUBI() {
	if cs.db == nil {
		return
	}
	threshold := time.Now().Unix() - escrowToUBISeconds

	// FIX 3: Use DELETE...RETURNING for an atomic claim — no other goroutine
	// can claim the same row between a SELECT and a subsequent DELETE.
	rows, err := cs.db.Query(
		`DELETE FROM escrow_accounts
		 WHERE moved_at < $1
		 RETURNING wallet_address, amount`,
		threshold,
	)
	if err != nil {
		fmt.Printf("[ESCROW] ReleaseEscrowToUBI query error: %v\n", err)
		return
	}
	type escrowEntry struct {
		addr   string
		amount float64
	}
	var entries []escrowEntry
	for rows.Next() {
		var e escrowEntry
		if scanErr := rows.Scan(&e.addr, &e.amount); scanErr == nil && e.amount > 0 {
			entries = append(entries, e)
		}
	}
	rows.Close() // close before acquiring cs.mu to avoid holding DB connection under lock

	if len(entries) == 0 {
		return
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	for _, e := range entries {
		// Credit UBI pool.
		if _, ok := cs.accounts[ubiPoolAddr]; !ok {
			cs.accounts[ubiPoolAddr] = &AccountState{Address: ubiPoolAddr}
		}
		cs.accounts[ubiPoolAddr].Balance = cs.accounts[ubiPoolAddr].Balance.Add(NewDecimal(round6(e.amount)))
		cs.saveAccountToDB(cs.accounts[ubiPoolAddr])
		fmt.Printf("[ESCROW] ✓ Released %.6f AEQ from %s escrow → UBI pool\n", e.amount, e.addr)
	}

	cs.save()
}
