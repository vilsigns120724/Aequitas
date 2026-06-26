package keeper

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

// round6 rounds a float64 to 6 decimal places, eliminating floating-point
// accumulation errors that build up in ledger operations over many transactions.
// This is the application-level fix for float64 imprecision; a full integer
// refactor (microAEQ int64) remains a future architecture task.
func round6(v float64) float64 {
	return math.Round(v*1_000_000) / 1_000_000
}

// ─── CONTRACT STORAGE ─────────────────────────────────────────────────────────

func (cs *ChainState) SaveContract(address string, bytecode []byte, deployer string) error {
	if cs.db == nil {
		return nil
	}
	address = strings.ToLower(address)
	_, err := cs.db.Exec(
		`INSERT INTO evm_contracts (address, bytecode, deployer) VALUES ($1, $2, $3)
 ON CONFLICT (address) DO UPDATE SET bytecode = $2`,
		address, hex.EncodeToString(bytecode), deployer,
	)
	if err != nil {
		fmt.Printf("[EVM] Error saving contract: %v\n", err)
	}
	return err
}

func (cs *ChainState) LoadContract(address string) ([]byte, error) {
	if cs.db == nil {
		return nil, nil
	}
	address = strings.ToLower(address)
	var bytecodeHex string
	err := cs.db.QueryRow(
		`SELECT bytecode FROM evm_contracts WHERE lower(address) = $1`, address,
	).Scan(&bytecodeHex)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	bytecodeHex = strings.TrimPrefix(strings.TrimPrefix(bytecodeHex, `\x`), "0x")
	return hex.DecodeString(bytecodeHex)
}

func (cs *ChainState) GetAllContracts() []string {
	if cs.db == nil {
		return nil
	}
	rows, err := cs.db.Query(`SELECT address FROM evm_contracts`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var addrs []string
	for rows.Next() {
		var addr string
		rows.Scan(&addr)
		addrs = append(addrs, addr)
	}
	return addrs
}

// ─── NONCE STORAGE ────────────────────────────────────────────────────────────

func (cs *ChainState) SaveNonce(address string, nonce uint64) error {
	if cs.db == nil {
		return nil
	}
	address = strings.ToLower(address)
	// Compare-and-swap: only advance the nonce, never decrease it.
	// Two nodes racing to reserve the same nonce would both issue
	// INSERT … nonce=$2; the second node's UPDATE fires but the
	// WHERE nonce < $2 clause rejects it, so the DB always holds
	// the highest reserved nonce.
	_, err := cs.db.Exec(
		`INSERT INTO evm_nonces (address, nonce) VALUES ($1, $2)
 ON CONFLICT (address) DO UPDATE SET nonce = $2 WHERE evm_nonces.nonce < $2`,
		address, nonce,
	)
	return err
}

// ReserveNonce atomically advances address from expected to next.
// It returns false when another process/node already reserved the same nonce.
func (cs *ChainState) ReserveNonce(address string, expected, next uint64) (bool, error) {
	if cs.db == nil {
		return true, nil
	}
	address = strings.ToLower(address)
	if expected == 0 {
		res, err := cs.db.Exec(
			`INSERT INTO evm_nonces (address, nonce) VALUES ($1, $2)
 ON CONFLICT (address) DO NOTHING`,
			address, next,
		)
		if err != nil {
			return false, err
		}
		if rows, _ := res.RowsAffected(); rows == 1 {
			return true, nil
		}
	}
	res, err := cs.db.Exec(
		`UPDATE evm_nonces SET nonce = $3 WHERE lower(address) = $1 AND nonce = $2`,
		address, expected, next,
	)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows == 1, nil
}

func (cs *ChainState) LoadNonce(address string) uint64 {
	if cs.db == nil {
		return 0
	}
	address = strings.ToLower(address)
	var nonce uint64
	cs.db.QueryRow(`SELECT nonce FROM evm_nonces WHERE lower(address) = $1`, address).Scan(&nonce)
	return nonce
}

// ─── CONTRACT STORAGE SLOTS ───────────────────────────────────────────────────

func (cs *ChainState) SaveStorageSlot(address, slot, value string) error {
	if cs.db == nil {
		return nil
	}
	address = strings.ToLower(address)
	_, err := cs.db.Exec(
		`INSERT INTO evm_storage (address, slot, value) VALUES ($1, $2, $3)
 ON CONFLICT (address, slot) DO UPDATE SET value = $3`,
		address, slot, value,
	)
	return err
}

// LoadAllStorageSlots returns every slot stored for address, used to back up
// contract state before a destructive upgrade so it can be restored on failure.
func (cs *ChainState) LoadAllStorageSlots(address string) (map[string]string, error) {
	if cs.db == nil {
		return nil, nil
	}
	address = strings.ToLower(address)
	rows, err := cs.db.Query(`SELECT slot, value FROM evm_storage WHERE lower(address) = $1`, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]string)
	for rows.Next() {
		var slot, value string
		if err := rows.Scan(&slot, &value); err == nil {
			out[slot] = value
		}
	}
	return out, nil
}

func (cs *ChainState) LoadStorageSlot(address, slot string) (string, error) {
	if cs.db == nil {
		return "", nil
	}
	address = strings.ToLower(address)
	var value string
	err := cs.db.QueryRow(
		`SELECT value FROM evm_storage WHERE lower(address) = $1 AND slot = $2`,
		address, slot,
	).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// ─── DUAL-LEDGER SYNC ────────────────────────────────────────────────────────

// MigrateEVMFromGoState rebuilds all V7 contract storage slots from the
// authoritative Go-state and database after an evm_storage wipe (e.g. on
// contract upgrade). Writes: totalSupply (slot 0), totalHumans (slot 1),
// balanceOf (slot 4), isHuman (slot 6), usedCommitments (slot 7),
// usedNullifiers (slot 8). Safe to call without holding cs.mu.
func (cs *ChainState) MigrateEVMFromGoState(contractAddr string) error {
	if cs.db == nil {
		return nil
	}
	contractAddr = strings.ToLower(contractAddr)
	fmt.Printf("[MIGRATE] Rebuilding EVM storage from Go-state for %s...\n", contractAddr)

	weiPerAEQ := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	var totalSupply float64
	var totalHumans int64

	cs.mu.RLock()
	for addr, acc := range cs.accounts {
		balBig, _ := new(big.Float).SetPrec(256).Mul(
			new(big.Float).SetFloat64(acc.Balance.Float()),
			new(big.Float).SetInt(weiPerAEQ),
		).Int(nil)
		if balBig == nil {
			balBig = new(big.Int)
		}
		addrBytes := common.HexToAddress(addr).Bytes()
		cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 4).Hex(), common.BigToHash(balBig).Hex())
		totalSupply += acc.Balance.Float()
		if acc.IsHuman {
			cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 6).Hex(), common.HexToHash("0x01").Hex())
			totalHumans++
			// Preserve lastActivity (slot 10) and lastDemurrage (slot 11).
			if acc.LastActivityAt > 0 {
				ts := big.NewInt(acc.LastActivityAt)
				cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 10).Hex(), common.BigToHash(ts).Hex())
				cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 11).Hex(), common.BigToHash(ts).Hex())
			}
			// Set ubiClaimed (slot 12) to the CURRENT ubiPerHumanAccumulated (slot 3).
			// This prevents double-claiming: after an upgrade, each human's "already claimed"
			// marker is set to the current accumulator so they can't re-claim historical UBI.
			// They can still earn new UBI from future distributions.
			// ubiPerHumanAccumulated (slot 3) will be read from EVM storage below.
			// We store a marker here; the actual slot-3 value is written after the loop.
		}
	}
	cs.mu.RUnlock()

	// totalSupply (slot 0) and totalHumans (slot 1)
	supplyWei, _ := new(big.Float).SetPrec(256).Mul(
		new(big.Float).SetFloat64(totalSupply),
		new(big.Float).SetInt(weiPerAEQ),
	).Int(nil)
	if supplyWei == nil {
		supplyWei = new(big.Int)
	}
	cs.SaveStorageSlot(contractAddr, common.BigToHash(big.NewInt(0)).Hex(), common.BigToHash(supplyWei).Hex())
	cs.SaveStorageSlot(contractAddr, common.BigToHash(big.NewInt(1)).Hex(), common.BigToHash(big.NewInt(totalHumans)).Hex())

	// Read current ubiPerHumanAccumulated (slot 3) from DB so we can set
	// ubiClaimed = that value for every human, preventing double-claim on upgrade.
	ubiAccumSlot := common.BigToHash(big.NewInt(3)).Hex()
	ubiAccumVal, _ := cs.LoadStorageSlot(contractAddr, ubiAccumSlot)
	if ubiAccumVal == "" {
		ubiAccumVal = common.Hash{}.Hex()
	}
	// Also write slot 2 (ubiPool) and slot 3 (ubiPerHumanAccumulated) — preserve existing
	// slot 3 value; it is NOT part of Go-state so we keep what was last in EVM.

	// Set ubiClaimed (slot 12) = ubiPerHumanAccumulated for every human to prevent double-claiming.
	cs.mu.RLock()
	for addr, acc := range cs.accounts {
		if acc.IsHuman {
			addrB := common.HexToAddress(addr).Bytes()
			cs.SaveStorageSlot(contractAddr, mappingSlot(addrB, 12).Hex(), ubiAccumVal)
		}
	}
	cs.mu.RUnlock()

	// usedNullifiers (slot 8): nullifier → wallet
	rows, err := cs.db.Query(`SELECT nullifier, wallet_address FROM nullifiers`)
	if err == nil {
		// P2-FIX: explicit Close instead of defer — defer fires at function
		// return, keeping both DB cursors open simultaneously during migration.
		for rows.Next() {
			var nullifier, wallet string
			if scanErr := rows.Scan(&nullifier, &wallet); scanErr != nil {
				fmt.Printf("[EVM] Warning: nullifier scan error in MigrateEVM: %v\n", scanErr)
				continue
			}
			nullKey := common.HexToHash(strings.TrimPrefix(nullifier, "0x"))
			nullSlot := mappingSlotBytes32(nullKey, 8)
			walletHash := common.BigToHash(common.HexToAddress(wallet).Big())
			cs.SaveStorageSlot(contractAddr, nullSlot.Hex(), walletHash.Hex())
		}
		rows.Close()
	}

	// usedCommitments (slot 7) + commitmentOf (slot 9): from bio_registrations
	rows2, err2 := cs.db.Query(`SELECT commitment, wallet_address FROM bio_registrations`)
	if err2 == nil {
		for rows2.Next() {
			var commitment, wallet string
			if err := rows2.Scan(&commitment, &wallet); err != nil {
				fmt.Printf("[MIGRATE] WARNING: bio_registrations scan error: %v — skipping row\n", err)
				continue
			}
			commitBig, ok := new(big.Int).SetString(strings.TrimPrefix(commitment, "0x"), 10)
			if !ok {
				commitBig, ok = new(big.Int).SetString(strings.TrimPrefix(commitment, "0x"), 16)
			}
			if !ok || commitBig == nil {
				continue
			}
			// usedCommitments[commitment] = true (slot 7)
			commitSlot7 := mappingSlot(common.LeftPadBytes(commitBig.Bytes(), 32), 7)
			cs.SaveStorageSlot(contractAddr, commitSlot7.Hex(), common.HexToHash("0x01").Hex())
			// commitmentOf[wallet] = commitment (slot 9)
			if wallet != "" {
				commitSlot9 := mappingSlot(common.HexToAddress(wallet).Bytes(), 9)
				cs.SaveStorageSlot(contractAddr, commitSlot9.Hex(), common.BigToHash(commitBig).Hex())
			}
		}
		rows2.Close()
	}

	// lastActivity (slot 10) + lastDemurrage (slot 11): from chain_accounts
	cs.mu.RLock()
	for addr, acc := range cs.accounts {
		if acc.LastActivityAt == 0 {
			continue
		}
		ts := big.NewInt(acc.LastActivityAt)
		addrBytes := common.HexToAddress(addr).Bytes()
		cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 10).Hex(), common.BigToHash(ts).Hex())
		cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 11).Hex(), common.BigToHash(ts).Hex())
	}
	cs.mu.RUnlock()

	// Restore guardian/escrow relationship slots (5, 13-16) that were saved
	// before the storage wipe. These are not tracked in any Go-state table so
	// they can only be preserved by snapshot + restore across the upgrade.
	cs.RestorePreUpgradeRelationshipSlots(contractAddr)

	fmt.Printf("[MIGRATE] ✓ EVM storage rebuilt: %d humans, %.2f AEQ total supply\n", totalHumans, totalSupply)
	return nil
}

// upgradeRelationshipSlotsTable is the name of the temporary table used to
// preserve guardian/escrow EVM storage slots across a contract upgrade wipe.
const upgradeRelationshipSlotsTable = "evm_upgrade_relationship_slots"

// SavePreUpgradeRelationshipSlots reads EVM storage slots 5 (escrowOf) and
// 13-16 (guardianOf / pendingGuardian / guardianRequestedAt / wardCount) from
// the live evm_storage table and copies them to a temporary snapshot table.
// Call this BEFORE wiping evm_storage on a contract upgrade; then call
// RestorePreUpgradeRelationshipSlots AFTER MigrateEVMFromGoState has rebuilt
// the rest of the storage to re-inject these slots back in.
func (cs *ChainState) SavePreUpgradeRelationshipSlots(contractAddr string) {
	if cs.db == nil {
		return
	}
	contractAddr = strings.ToLower(contractAddr)
	cs.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		address TEXT NOT NULL,
		slot    TEXT NOT NULL,
		value   TEXT NOT NULL,
		PRIMARY KEY (address, slot)
	)`, upgradeRelationshipSlotsTable))
	// Clear any stale snapshot from a previous upgrade cycle.
	cs.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE address = $1`, upgradeRelationshipSlotsTable), contractAddr)

	// We can't filter by "slot prefix for base slot N" efficiently in SQL because
	// the slot hash is opaque. Instead, snapshot ALL slots for this address that
	// we cannot reconstruct from Go-state (i.e., everything EXCEPT the slots
	// MigrateEVMFromGoState already writes: 0,1,4,6,7,8,9,10,11,12). We do
	// this by saving all slots and then letting MigrateEVMFromGoState overwrite
	// the ones it knows about, so only the truly-opaque slots (5,13-16) survive.
	rows, err := cs.db.Query(`SELECT slot, value FROM evm_storage WHERE lower(address) = $1`, contractAddr)
	if err != nil {
		fmt.Printf("[DEPLOY] Warning: could not snapshot relationship slots: %v\n", err)
		return
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var slot, value string
		rows.Scan(&slot, &value)
		cs.db.Exec(fmt.Sprintf(`INSERT INTO %s (address, slot, value) VALUES ($1, $2, $3)
			ON CONFLICT (address, slot) DO UPDATE SET value = $3`, upgradeRelationshipSlotsTable),
			contractAddr, slot, value)
		count++
	}
	fmt.Printf("[DEPLOY] Saved %d EVM storage slots for guardian/escrow preservation\n", count)
}

// RestorePreUpgradeRelationshipSlots writes the slots that were saved by
// SavePreUpgradeRelationshipSlots back into evm_storage. Called at the end of
// MigrateEVMFromGoState so these survive the upgrade storage wipe.
func (cs *ChainState) RestorePreUpgradeRelationshipSlots(contractAddr string) {
	if cs.db == nil {
		return
	}
	contractAddr = strings.ToLower(contractAddr)
	rows, err := cs.db.Query(fmt.Sprintf(`SELECT slot, value FROM %s WHERE address = $1`, upgradeRelationshipSlotsTable), contractAddr)
	if err != nil {
		return // table doesn't exist yet (first-ever deploy, no prior snapshot)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var slot, value string
		rows.Scan(&slot, &value)
		// Use INSERT … ON CONFLICT DO NOTHING so that slots MigrateEVMFromGoState
		// already wrote (balanceOf, isHuman, etc.) are not overwritten by stale
		// pre-upgrade values — only truly-missing slots get restored.
		cs.db.Exec(`INSERT INTO evm_storage (address, slot, value) VALUES ($1, $2, $3)
			ON CONFLICT (address, slot) DO NOTHING`,
			contractAddr, slot, value)
		count++
	}
	if count > 0 {
		fmt.Printf("[MIGRATE] Restored %d guardian/escrow slots from pre-upgrade snapshot\n", count)
	}
}

// SyncBalancesToEVM writes the current Go-state AEQ balance for each addr into
// the AequitasV7 contract's balanceOf storage slot (mapping at position 4),
// keeping both ledgers consistent after every Go-state change.
func (cs *ChainState) SyncBalancesToEVM(contractAddr string, addrs ...string) {
	if cs.db == nil {
		return
	}
	contractAddr = strings.ToLower(contractAddr)
	weiPerAEQ := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	for _, addr := range addrs {
		addr = strings.ToLower(addr)
		cs.mu.RLock()
		acc, ok := cs.accounts[addr]
		cs.mu.RUnlock()
		var bal float64
		if ok {
			// P2-12: use effectiveBalance (with demurrage decay) so EVM
			// storage matches what the user actually holds right now,
			// not the raw stored value which may be higher than actual.
			bal = effectiveBalance(acc).Float()
		}
		balBig, _ := new(big.Float).SetPrec(256).Mul(
			new(big.Float).SetFloat64(bal),
			new(big.Float).SetInt(weiPerAEQ),
		).Int(nil)
		if balBig == nil {
			balBig = new(big.Int)
		}
		slot := mappingSlot(common.HexToAddress(addr).Bytes(), 4).Hex()
		val := common.BigToHash(balBig).Hex()
		if err := cs.SaveStorageSlot(contractAddr, slot, val); err != nil {
			fmt.Printf("[EVM] Warning: could not sync balance for %s: %v\n", addr, err)
		}
	}
}

// syncHumanRegistrationLocked writes balanceOf (slot 4), isHuman (slot 6),
// lastActivity (slot 10), and lastDemurrage (slot 11) EVM slots for a newly
// registered human. Must be called only while the caller already holds cs.mu (write lock).
// syncBalanceLocked now handles all four slots, so this is a simple delegation.
func (cs *ChainState) syncHumanRegistrationLocked(contractAddr string, addr string) {
	cs.syncBalanceLocked(contractAddr, addr)
}

// syncBalanceLocked is like SyncBalancesToEVM but reads cs.accounts directly
// without acquiring cs.mu. Must be called only while the caller already holds
// cs.mu (read or write lock) — calling SyncBalancesToEVM from inside a locked
// function would deadlock on the inner RLock().
// Syncs slots: 4 (balanceOf), 6 (isHuman), 10 (lastActivity), 11 (lastDemurrage).
func (cs *ChainState) syncBalanceLocked(contractAddr string, addrs ...string) {
	if cs.db == nil {
		return
	}
	contractAddr = strings.ToLower(contractAddr)
	weiPerAEQ := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	for _, addr := range addrs {
		addr = strings.ToLower(addr)
		acc, ok := cs.accounts[addr]
		var bal float64
		if ok {
			// P1-4: use effectiveBalance (demurrage-adjusted) so the EVM slot
			// matches the user's real spendable amount, not the stored pre-decay value.
			bal = effectiveBalance(acc).Float()
		}
		balBig, _ := new(big.Float).SetPrec(256).Mul(
			new(big.Float).SetFloat64(bal),
			new(big.Float).SetInt(weiPerAEQ),
		).Int(nil)
		if balBig == nil {
			balBig = new(big.Int)
		}
		addrBytes := common.HexToAddress(addr).Bytes()
		// slot 4: balanceOf
		if err := cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 4).Hex(), common.BigToHash(balBig).Hex()); err != nil {
			fmt.Printf("[EVM] Warning: could not sync balance for %s: %v\n", addr, err)
		}
		if !ok {
			continue
		}
		// slot 6: isHuman
		isHumanVal := common.HexToHash("0x00")
		if acc.IsHuman {
			isHumanVal = common.HexToHash("0x01")
		}
		cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 6).Hex(), isHumanVal.Hex())
		// slots 10 + 11: lastActivity / lastDemurrage
		if acc.LastActivityAt > 0 {
			ts := common.BigToHash(big.NewInt(acc.LastActivityAt))
			cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 10).Hex(), ts.Hex())
			cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 11).Hex(), ts.Hex())
		}
	}
}

// ─── EVM ENGINE HELPERS ───────────────────────────────────────────────────────

// PersistContractStorage reads storage slots from a stateDB and saves to PostgreSQL.
// Since we no longer have a persistent stateDB, this is a no-op log.
func (e *EVMEngine) PersistContractStorage(contractAddr common.Address) {
	fmt.Printf("[EVM] Contract %s active in session\n", strings.ToLower(contractAddr.Hex()))
}

// NewPersistentStateDB creates a StateDB loaded from PostgreSQL.
// Used by tests and legacy code. For production use EVMEngine.newStateDB().
func NewPersistentStateDB(cs *ChainState) (*state.StateDB, error) {
	memDB := rawdb.NewMemoryDatabase()
	sdb, err := state.New(common.Hash{}, state.NewDatabase(memDB), nil)
	if err != nil {
		return nil, err
	}

	// P2-AUDIT: Do NOT call LoadNonce per account — that issues N PostgreSQL
	// queries (one per account) and creates a DoS vector. EVM nonces for
	// sends are managed by the RPC layer; the legacy StateDB doesn't need
	// per-account nonces for call execution. Matches the fix in newStateDB.
	for _, acc := range cs.GetAllAccounts() {
		addr := common.HexToAddress(acc.Address)
		// P1-FIX: acc.Balance is a Decimal (int64 micro-units). Use .Float()
		// to get the AEQ float value before converting to wei. Using
		// int64(acc.Balance) directly would re-interpret micro-AEQ as whole-AEQ
		// and multiply by 1e18 a second time, overstating balances by 1e6×.
		decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		balWei, _ := new(big.Float).SetPrec(256).Mul(
			new(big.Float).SetFloat64(acc.Balance.Float()),
			new(big.Float).SetInt(decimals),
		).Int(nil)
		if balWei == nil {
			balWei = new(big.Int)
		}
		sdb.SetBalance(addr, balWei)
	}

	for _, addrStr := range cs.GetAllContracts() {
		addr := common.HexToAddress(addrStr)
		code, err := cs.LoadContract(addrStr)
		if err != nil || len(code) == 0 {
			continue
		}
		sdb.SetCode(addr, code)
		fmt.Printf("[EVM] Loaded contract: %s (%d bytes)\n", addrStr, len(code))

		if cs.db != nil {
			rows, err := cs.db.Query(
				`SELECT slot, value FROM evm_storage WHERE address = $1`, addrStr)
			if err == nil {
				for rows.Next() {
					var slot, value string
					rows.Scan(&slot, &value)
					sdb.SetState(addr, common.HexToHash(slot), common.HexToHash(value))
				}
				rows.Close()
			}
		}
	}

	sdb.Commit(0, false)
	return sdb, nil
}

// SaveBioRegistration links a ZK proof commitment to the wallet that
// successfully registered with it. Called once, right after a
// registerWithSig transaction is confirmed successful — never speculatively.
// bioHash is also stored alongside the commitment so the app (which only
// ever knows its own bioHash, not the commitment computed on the website
// under the new flow) can poll for its registration — see
// GetWalletByBioHash below.
func (cs *ChainState) SaveBioRegistration(commitment, walletAddress, txHash, bioHash string) error {
	if commitment == "" {
		return fmt.Errorf("empty commitment rejected")
	}
	if cs.db == nil {
		return nil
	}
	walletAddress = strings.ToLower(walletAddress)
	if bioHash != "" {
		existing := cs.GetWalletByBioHash(bioHash)
		if existing != "" && strings.ToLower(existing) != walletAddress {
			return fmt.Errorf("biometric already registered to %s", existing)
		}
	}
	// P2-AUDIT: Use ON CONFLICT DO NOTHING to protect the first successful
	// registration from being overwritten by a concurrent/replay registration
	// with the same commitment. The contract itself enforces commitment uniqueness
	// on-chain; the DB row is just a mirror for polling — never the authority.
	_, err := cs.db.Exec(
		`INSERT INTO bio_registrations (commitment, wallet_address, tx_hash, bio_hash) VALUES ($1, $2, $3, $4)
 ON CONFLICT (commitment) DO NOTHING`,
		commitment, walletAddress, txHash, bioHash,
	)
	return err
}

// GetWalletByCommitment looks up which wallet (if any) successfully
// registered with a given proof commitment. Returns "" if none found —
// this lets the app ask "did MY specific proof get registered?" instead of
// reading the last entry in a global, unfiltered humans list.
func (cs *ChainState) GetWalletByCommitment(commitment string) string {
	if cs.db == nil {
		return ""
	}
	var wallet string
	err := cs.db.QueryRow(`SELECT wallet_address FROM bio_registrations WHERE commitment = $1`, commitment).Scan(&wallet)
	if err != nil {
		return ""
	}
	return wallet
}

// GetWalletByBioHash looks up which wallet (if any) most recently
// completed registration for a given device biometric identity hash.
// Used by the app's post-bioHash-flow polling (startPollingByBioHash) —
// the app never computes a commitment itself under that flow, only its
// own bioHash, so this is the only key it can reliably poll by.
func (cs *ChainState) GetWalletByBioHash(bioHash string) string {
	if cs.db == nil {
		return ""
	}
	var wallet string
	err := cs.db.QueryRow(`SELECT wallet_address FROM bio_registrations WHERE bio_hash = $1 ORDER BY registered_at DESC LIMIT 1`, bioHash).Scan(&wallet)
	if err != nil {
		return ""
	}
	return wallet
}

// SaveBioHash writes the biometric hash into the bio_hashes table after a
// confirmed registration. The proof server's /check and /prove endpoints
// read from this table to block duplicate biometric registrations — keeping
// it in sync with bio_registrations ensures both layers see the same state.
func (cs *ChainState) SaveBioHash(bioHash, walletAddress string) {
	if cs.db == nil || bioHash == "" {
		return
	}
	walletAddress = strings.ToLower(walletAddress)
	_, err := cs.db.Exec(
		`INSERT INTO bio_hashes (hash, wallet_address) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING`,
		bioHash, walletAddress,
	)
	if err != nil {
		fmt.Printf("[REGISTER] Warning: could not sync bio_hashes: %v\n", err)
	}
}

// ─── NULLIFIERS ───────────────────────────────────────────────────────────────
//
// A nullifier is a one-way derivation of the biometric secret:
//   nullifier = SHA256(bioHash + ":aequitas-ubi-v1")
//
// It is computed by the client and stored on-chain after a successful
// registration. Because the same biometric always produces the same bioHash
// (on the same device), it always produces the same nullifier — so a second
// registration attempt reveals an already-used nullifier and is rejected,
// even if the user switches wallets. The server never sees the raw bioHash
// in this step, only its SHA256 derivative. In a future ZK upgrade the
// nullifier will be generated inside the Groth16 circuit itself (Semaphore
// style), removing even the SHA256 link.

func (cs *ChainState) IsNullifierUsed(nullifier string) bool {
	cs.mu.RLock()
	_, inMem := cs.nullifiers[nullifier]
	cs.mu.RUnlock()
	if inMem {
		return true
	}
	if cs.db == nil {
		return false
	}
	var wallet string
	err := cs.db.QueryRow(`SELECT wallet_address FROM nullifiers WHERE nullifier = $1`, nullifier).Scan(&wallet)
	return err == nil && wallet != ""
}

// maxInMemNullifiers caps the in-memory nullifier cache to ~50 MB at 1M entries.
// P3-7: above this threshold new nullifiers are only written to DB; lookups
// fall through to the DB automatically via IsNullifierUsed.
const maxInMemNullifiers = 500_000

// TryClaimNullifier atomically inserts the nullifier and returns true if it
// was newly inserted (this caller owns the registration), false if it already
// existed (another goroutine or a previous replay already claimed it).
// Using a DB-level INSERT … ON CONFLICT eliminates the TOCTOU window between
// IsNullifierUsed and SaveNullifier — no separate mutex required.
func (cs *ChainState) TryClaimNullifier(nullifier, walletAddress string) bool {
	if nullifier == "" {
		return false
	}
	walletAddress = strings.ToLower(walletAddress)
	if cs.db == nil {
		cs.mu.Lock()
		if _, exists := cs.nullifiers[nullifier]; exists {
			cs.mu.Unlock()
			return false
		}
		cs.nullifiers[nullifier] = walletAddress
		cs.mu.Unlock()
		return true
	}
	res, err := cs.db.Exec(
		`INSERT INTO nullifiers (nullifier, wallet_address) VALUES ($1, $2) ON CONFLICT (nullifier) DO NOTHING`,
		nullifier, walletAddress,
	)
	if err != nil {
		fmt.Printf("[NULLIFIER] TryClaimNullifier DB error: %v\n", err)
		return false
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return false // already existed
	}
	// Insert succeeded — update in-memory cache.
	cs.mu.Lock()
	if len(cs.nullifiers) < maxInMemNullifiers {
		cs.nullifiers[nullifier] = walletAddress
	}
	cs.mu.Unlock()
	return true
}

func (cs *ChainState) SaveNullifier(nullifier, walletAddress string) {
	if nullifier == "" {
		return
	}
	walletAddress = strings.ToLower(walletAddress)
	cs.mu.Lock()
	if len(cs.nullifiers) < maxInMemNullifiers {
		cs.nullifiers[nullifier] = walletAddress
	}
	cs.mu.Unlock()
	if cs.db == nil {
		return
	}
	if _, err := cs.db.Exec(
		`INSERT INTO nullifiers (nullifier, wallet_address) VALUES ($1, $2) ON CONFLICT (nullifier) DO NOTHING`,
		nullifier, walletAddress,
	); err != nil {
		fmt.Printf("[NULLIFIER] Warning: could not persist nullifier: %v\n", err)
	}
}

func (cs *ChainState) GetWalletByNullifier(nullifier string) string {
	cs.mu.RLock()
	w, ok := cs.nullifiers[nullifier]
	cs.mu.RUnlock()
	if ok {
		return w
	}
	if cs.db == nil {
		return ""
	}
	var wallet string
	cs.db.QueryRow(`SELECT wallet_address FROM nullifiers WHERE nullifier = $1`, nullifier).Scan(&wallet)
	return wallet
}

// ─── PRICE HISTORY ───────────────────────────────────────────────────────────

func (cs *ChainState) InitPriceSnapshotsTable() {
	if cs.db == nil {
		return
	}
	cs.db.Exec(`CREATE TABLE IF NOT EXISTS price_snapshots (
		id           SERIAL PRIMARY KEY,
		price        DOUBLE PRECISION NOT NULL,
		reserve_aeq  DOUBLE PRECISION NOT NULL,
		reserve_tusd DOUBLE PRECISION NOT NULL,
		captured_at  TIMESTAMP DEFAULT NOW()
	)`)
	// Keep only last 30 days (~324000 rows at 8s intervals) — purge older rows
	cs.db.Exec(`DELETE FROM price_snapshots WHERE captured_at < NOW() - INTERVAL '30 days'`)
}

// SavePriceSnapshot records the current AEQ/tUSD price. Must be safe to call
// concurrently — copies pool values under RLock before the DB write so a
// concurrent swap cannot modify cs.pool while we're reading it.
func (cs *ChainState) SavePriceSnapshot() {
	if cs.db == nil {
		return
	}
	cs.mu.RLock()
	if cs.pool == nil || cs.pool.ReserveAEQ <= 0 || cs.pool.ReserveTUSD <= 0 {
		cs.mu.RUnlock()
		return
	}
	price := cs.pool.ReserveTUSD.Float() / cs.pool.ReserveAEQ.Float()
	aeq := cs.pool.ReserveAEQ.Float()
	tusd := cs.pool.ReserveTUSD.Float()
	cs.mu.RUnlock()
	cs.db.Exec(`INSERT INTO price_snapshots (price, reserve_aeq, reserve_tusd) VALUES ($1, $2, $3)`,
		price, aeq, tusd)
}

// GetPriceHistory returns price snapshots from the last `minutes` minutes,
// limited to `limit` points. Returns [{t, p, aeq, tusd}, ...].
// minutes is clamped to 1-43200, limit to 1-5000.
func (cs *ChainState) GetPriceHistory(minutes, limit int) []map[string]interface{} {
	if cs.db == nil {
		return nil
	}
	if minutes < 1 {
		minutes = 1
	}
	if minutes > 43200 {
		minutes = 43200
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 5000 {
		limit = 5000
	}
	// P1-11: use ($1 * INTERVAL '1 minute') instead of string concat to
	// prevent any future SQL-injection if $1 type changes to string.
	rows, err := cs.db.Query(`
		SELECT EXTRACT(EPOCH FROM captured_at)::BIGINT, price, reserve_aeq, reserve_tusd
		FROM price_snapshots
		WHERE captured_at >= NOW() - ($1 * INTERVAL '1 minute')
		ORDER BY captured_at ASC
		LIMIT $2`, minutes, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var ts int64
		var price, aeq, tusd float64
		rows.Scan(&ts, &price, &aeq, &tusd)
		result = append(result, map[string]interface{}{
			"t": ts * 1000, // milliseconds for JS Date
			"p": price,
			"a": aeq,
			"u": tusd,
		})
	}
	return result
}

// ─── GINI HISTORY ────────────────────────────────────────────────────────────

func (cs *ChainState) InitGiniSnapshotsTable() {
	if cs.db == nil {
		return
	}
	cs.db.Exec(`CREATE TABLE IF NOT EXISTS gini_snapshots (
		id          SERIAL PRIMARY KEY,
		gini        DOUBLE PRECISION NOT NULL,
		humans      INT NOT NULL,
		captured_at TIMESTAMP DEFAULT NOW()
	)`)
}

// SaveGiniSnapshot persists the current Gini coefficient. Called after each
// UBI distribution so the history chart has real data points over time.
// Must NOT be called while cs.mu is held — CalcGini acquires RLock internally.
func (cs *ChainState) SaveGiniSnapshot() {
	if cs.db == nil {
		return
	}
	gini := cs.CalcGini()      // acquires RLock
	humans := cs.TotalHumans() // acquires RLock
	cs.db.Exec(`INSERT INTO gini_snapshots (gini, humans) VALUES ($1, $2)`, gini, humans)
}

// SaveGiniSnapshotValues saves a pre-computed Gini/humans pair without
// acquiring any lock. Call this from inside a locked function by passing
// values already read under the lock, to avoid lock-reentrancy deadlocks.
func (cs *ChainState) SaveGiniSnapshotValues(gini float64, humans int) {
	if cs.db == nil {
		return
	}
	cs.db.Exec(`INSERT INTO gini_snapshots (gini, humans) VALUES ($1, $2)`, gini, humans)
}

// GetGiniHistory returns the last n Gini snapshots in chronological order.
// Returns a slice of maps with keys: idx (0-100), gini (0-1), humans, timestamp.
func (cs *ChainState) GetGiniHistory(n int) []map[string]interface{} {
	if cs.db == nil {
		return nil
	}
	rows, err := cs.db.Query(
		`SELECT gini, humans, EXTRACT(EPOCH FROM captured_at)::BIGINT
		 FROM gini_snapshots ORDER BY captured_at DESC LIMIT $1`, n)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var gini float64
		var humans int
		var ts int64
		rows.Scan(&gini, &humans, &ts)
		result = append(result, map[string]interface{}{
			"idx":       gini * 100,
			"gini":      gini,
			"humans":    humans,
			"timestamp": ts,
		})
	}
	// Reverse to get chronological order (we queried DESC).
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// ─── SWAP NONCES ─────────────────────────────────────────────────────────────
//
// Each wallet has a monotonically increasing nonce for swap/liquidity actions.
// The nonce is included in the signed message, so a captured signature cannot
// be replayed — the nonce check atomically rejects any second use.

// ─── VALIDATOR KEY REGISTRY ──────────────────────────────────────────────────
//
// Replaces the shared PEER_SECRET model with individual, human-authorized
// validator keys. Each node operator signs their signing key with their
// registered human wallet, creating a 1:1 link: "this human authorizes
// this signing key to produce blocks on their behalf."
//
// A compromised node key can be revoked individually without affecting any
// other validator. Authorization is tied to on-chain human identity.

func (cs *ChainState) InitValidatorKeysTable() {
	if cs.db == nil {
		return
	}
	cs.db.Exec(`CREATE TABLE IF NOT EXISTS validator_keys (
		signing_address TEXT PRIMARY KEY,
		human_wallet    TEXT NOT NULL UNIQUE,
		registered_at   TIMESTAMP DEFAULT NOW()
	)`)
	// Add UNIQUE on human_wallet if the table already existed without it.
	// Remove any existing duplicates first so the index creation succeeds.
	cs.db.Exec(`DELETE FROM validator_keys vk1
		USING validator_keys vk2
		WHERE vk1.registered_at < vk2.registered_at
		  AND vk1.human_wallet = vk2.human_wallet`)
	if _, err := cs.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_validator_keys_human_wallet
		ON validator_keys (human_wallet)`); err != nil {
		Log.Error("Could not enforce UNIQUE(human_wallet) on validator_keys", "error", err)
		// Node continues but duplicate keys may affect reward distribution
	}
}

// RegisterValidatorKey links a node signing address to a registered human
// wallet, authorizing that signing key to propose blocks. The human_wallet
// must be a registered human; the signature must be a valid personal_sign
// of "Aequitas: authorize validator key {signing_address}".
func (cs *ChainState) RegisterValidatorKey(signingAddress, humanWallet string) error {
	if cs.db == nil {
		return fmt.Errorf("no database")
	}
	signingAddress = strings.ToLower(strings.TrimSpace(signingAddress))
	humanWallet = strings.ToLower(strings.TrimSpace(humanWallet))
	if !cs.IsHuman(humanWallet) {
		return fmt.Errorf("human_wallet %s is not a registered human", humanWallet)
	}
	_, err := cs.db.Exec(
		`INSERT INTO validator_keys (signing_address, human_wallet) VALUES ($1, $2)
		 ON CONFLICT (signing_address) DO UPDATE SET human_wallet = $2, registered_at = NOW()`,
		signingAddress, humanWallet)
	return err
}

// LoadValidatorKeysIntoDAG reads all registered validator signing addresses
// from the DB and adds them to the DAG's authorized validators set.
// Called at startup so keys registered before the node restarted are effective.
func (cs *ChainState) LoadValidatorKeysIntoDAG(dag interface{ AddAuthorizedValidator(string) }) {
	if cs.db == nil {
		return
	}
	rows, err := cs.db.Query(`SELECT signing_address FROM validator_keys`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var addr string
		rows.Scan(&addr)
		dag.AddAuthorizedValidator(strings.ToLower(strings.TrimSpace(addr)))
	}
}

func (cs *ChainState) GetValidatorKeys() []map[string]string {
	if cs.db == nil {
		return nil
	}
	rows, err := cs.db.Query(`SELECT signing_address, human_wallet FROM validator_keys ORDER BY registered_at`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []map[string]string
	for rows.Next() {
		var addr, wallet string
		rows.Scan(&addr, &wallet)
		result = append(result, map[string]string{"signing_address": addr, "human_wallet": wallet})
	}
	return result
}

// ValidateNodeOperatorWallet returns an error string if the wallet is not a
// registered human. The calling code must STOP registration if this returns
// non-empty — rewards go only to verified humans, no exceptions.
func (cs *ChainState) ValidateNodeOperatorWallet(wallet string) string {
	if !cs.IsHuman(strings.ToLower(wallet)) {
		return "wallet " + wallet + " is NOT a registered human — register via the Android app first before running a node"
	}
	return ""
}

// InitSwapNoncesTable creates the swap_nonces table if it doesn't exist.
func (cs *ChainState) InitSwapNoncesTable() {
	if cs.db == nil {
		return
	}
	cs.db.Exec(`CREATE TABLE IF NOT EXISTS swap_nonces (
		wallet_address TEXT PRIMARY KEY,
		next_nonce     BIGINT NOT NULL DEFAULT 0
	)`)
}

// GetSwapNonce returns the next nonce a wallet should sign with.
// Returns 0 for wallets that have never performed a swap.
func (cs *ChainState) GetSwapNonce(wallet string) int64 {
	if cs.db == nil {
		return 0
	}
	wallet = strings.ToLower(wallet)
	var nonce int64
	cs.db.QueryRow(`SELECT next_nonce FROM swap_nonces WHERE wallet_address = $1`, wallet).Scan(&nonce)
	return nonce
}

// RestoreSwapNonce decrements the nonce back to its pre-swap value when a
// swap fails after the nonce was already consumed. Safe to call: if the nonce
// has already advanced past nonce+1 (extremely unlikely concurrent case) the
// UPDATE finds no rows and the decrement is skipped — user must re-sign.
func (cs *ChainState) RestoreSwapNonce(wallet string, nonce int64) {
	if cs.db == nil {
		return
	}
	wallet = strings.ToLower(wallet)
	cs.db.Exec(`UPDATE swap_nonces SET next_nonce = $2 WHERE wallet_address = $1 AND next_nonce = $2 + 1`,
		wallet, nonce)
}

// ConsumeSwapNonce atomically verifies that nonce matches the expected value
// and increments it. Returns an error if the nonce doesn't match (replay or
// wrong value). Must be called only after the signature has been verified.
func (cs *ChainState) ConsumeSwapNonce(wallet string, nonce int64) error {
	if cs.db == nil {
		return nil // no DB — skip in development
	}
	wallet = strings.ToLower(wallet)
	var result interface{ RowsAffected() (int64, error) }
	var err error
	if nonce == 0 {
		// First ever swap for this wallet — insert with next_nonce=1.
		result, err = cs.db.Exec(
			`INSERT INTO swap_nonces (wallet_address, next_nonce) VALUES ($1, 1)
			 ON CONFLICT (wallet_address) DO NOTHING`, wallet)
	} else {
		// Subsequent swap — increment only if current value matches.
		result, err = cs.db.Exec(
			`UPDATE swap_nonces SET next_nonce = next_nonce + 1
			 WHERE wallet_address = $1 AND next_nonce = $2`, wallet, nonce)
	}
	if err != nil {
		return fmt.Errorf("nonce db error: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("nonce %d already used or invalid — replay rejected", nonce)
	}
	return nil
}

// ─── EVM TX RECEIPTS (persistent — survives node restart) ────────────────────

// SaveTxReceipt persists an EVM transaction receipt to the database so MetaMask
// can retrieve it after a node restart. Without this, restarts cleared the
// in-memory txStatus map and MetaMask would show successful transactions as
// "Senden fehlgeschlagen" (failed) because receipts returned null.
func (cs *ChainState) SaveTxReceipt(txHash, fromAddr, toAddr, status string) {
	if cs.db == nil {
		return
	}
	_, err := cs.db.Exec(
		`INSERT INTO evm_tx_receipts (tx_hash, from_addr, to_addr, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (tx_hash) DO UPDATE SET status = $4`,
		strings.ToLower(txHash), strings.ToLower(fromAddr),
		strings.ToLower(toAddr), status, time.Now().Unix(),
	)
	if err != nil {
		fmt.Printf("[EVM] SaveTxReceipt error for %s: %v\n", txHash, err)
		return
	}
	// Prune old receipts — keep only the latest 10,000 to prevent unbounded growth.
	cs.db.Exec(`DELETE FROM evm_tx_receipts WHERE tx_hash NOT IN (
		SELECT tx_hash FROM evm_tx_receipts ORDER BY created_at DESC LIMIT 10000
	)`)
}

// GetTxReceipt looks up a persisted receipt. Returns (fromAddr, toAddr, status, found).
// Called by getTransactionReceipt when the txHash is not in the in-memory cache.
func (cs *ChainState) GetTxReceipt(txHash string) (fromAddr, toAddr, status string, found bool) {
	if cs.db == nil {
		return "", "", "", false
	}
	err := cs.db.QueryRow(
		`SELECT from_addr, COALESCE(to_addr, ''), status FROM evm_tx_receipts WHERE tx_hash = $1`,
		strings.ToLower(txHash),
	).Scan(&fromAddr, &toAddr, &status)
	if err == sql.ErrNoRows || err != nil {
		return "", "", "", false
	}
	return fromAddr, toAddr, status, true
}

// ─── PENDING TXs (persistent — survive node restart) ─────────────────────────

// SavePendingTx writes a Transaction to the DB so it survives node restarts.
// ProduceBlock calls LoadAndClearPendingTxs to drain these and include them
// in the next block, ensuring secondary nodes receive every state change.
func (cs *ChainState) SavePendingTx(tx Transaction) {
	if cs.db == nil {
		return
	}
	data, err := json.Marshal(tx)
	if err != nil {
		fmt.Printf("[TX] SavePendingTx marshal error: %v\n", err)
		return
	}
	if _, err := cs.db.Exec(
		`INSERT INTO pending_txs (tx_json, created_at) VALUES ($1, $2)`,
		string(data), time.Now().Unix(),
	); err != nil {
		fmt.Printf("[TX] SavePendingTx db error: %v\n", err)
	}
}

// LoadAndClearPendingTxs atomically reads all DB-pending TXs and deletes them.
// Called by ProduceBlock so that restart-surviving TXs are included once.
func (cs *ChainState) LoadAndClearPendingTxs() []Transaction {
	if cs.db == nil {
		return nil
	}
	rows, err := cs.db.Query(`SELECT id, tx_json FROM pending_txs ORDER BY id`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var txs []Transaction
	var ids []int64
	for rows.Next() {
		var id int64
		var raw string
		if err := rows.Scan(&id, &raw); err != nil {
			continue
		}
		var tx Transaction
		if err := json.Unmarshal([]byte(raw), &tx); err != nil {
			fmt.Printf("[TX] LoadAndClearPendingTxs unmarshal error: %v\n", err)
			continue
		}
		txs = append(txs, tx)
		ids = append(ids, id)
	}
	rows.Close()
	if len(ids) > 0 {
		// Delete in a separate query after closing rows
		for _, id := range ids {
			cs.db.Exec(`DELETE FROM pending_txs WHERE id = $1`, id)
		}
	}
	return txs
}
