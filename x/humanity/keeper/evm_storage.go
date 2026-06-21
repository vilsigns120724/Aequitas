package keeper

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

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
	_, err := cs.db.Exec(
		`INSERT INTO evm_nonces (address, nonce) VALUES ($1, $2)
 ON CONFLICT (address) DO UPDATE SET nonce = $2`,
		address, nonce,
	)
	return err
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
func (cs *ChainState) MigrateEVMFromGoState(contractAddr string) {
	if cs.db == nil {
		return
	}
	contractAddr = strings.ToLower(contractAddr)
	fmt.Printf("[MIGRATE] Rebuilding EVM storage from Go-state for %s...\n", contractAddr)

	weiPerAEQ := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	var totalSupply float64
	var totalHumans int64

	cs.mu.RLock()
	for addr, acc := range cs.accounts {
		balBig, _ := new(big.Float).SetPrec(256).Mul(
			new(big.Float).SetFloat64(acc.Balance),
			new(big.Float).SetInt(weiPerAEQ),
		).Int(nil)
		if balBig == nil {
			balBig = new(big.Int)
		}
		addrBytes := common.HexToAddress(addr).Bytes()
		cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 4).Hex(), common.BigToHash(balBig).Hex())
		totalSupply += acc.Balance
		if acc.IsHuman {
			cs.SaveStorageSlot(contractAddr, mappingSlot(addrBytes, 6).Hex(), common.HexToHash("0x01").Hex())
			totalHumans++
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

	// usedNullifiers (slot 8): nullifier → wallet
	rows, err := cs.db.Query(`SELECT nullifier, wallet_address FROM nullifiers`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nullifier, wallet string
			rows.Scan(&nullifier, &wallet)
			nullKey := common.HexToHash(strings.TrimPrefix(nullifier, "0x"))
			nullSlot := mappingSlotBytes32(nullKey, 8)
			walletHash := common.BigToHash(common.HexToAddress(wallet).Big())
			cs.SaveStorageSlot(contractAddr, nullSlot.Hex(), walletHash.Hex())
		}
	}

	// usedCommitments (slot 7) + commitmentOf (slot 9): from bio_registrations
	rows2, err2 := cs.db.Query(`SELECT commitment, wallet_address FROM bio_registrations`)
	if err2 == nil {
		defer rows2.Close()
		for rows2.Next() {
			var commitment, wallet string
			rows2.Scan(&commitment, &wallet)
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

	fmt.Printf("[MIGRATE] ✓ EVM storage rebuilt: %d humans, %.2f AEQ total supply\n", totalHumans, totalSupply)
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
			bal = acc.Balance
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

// syncHumanRegistrationLocked writes both the balanceOf (slot 4) and
// isHuman (slot 6) EVM slots for a newly registered human. Must be called
// only while the caller already holds cs.mu (write lock).
func (cs *ChainState) syncHumanRegistrationLocked(contractAddr string, addr string) {
	cs.syncBalanceLocked(contractAddr, addr)
	isHumanSlot := mappingSlot(common.HexToAddress(addr).Bytes(), 6)
	_ = cs.SaveStorageSlot(strings.ToLower(contractAddr), isHumanSlot.Hex(), common.HexToHash("0x01").Hex())
}

// syncBalanceLocked is like SyncBalancesToEVM but reads cs.accounts directly
// without acquiring cs.mu. Must be called only while the caller already holds
// cs.mu (read or write lock) — calling SyncBalancesToEVM from inside a locked
// function would deadlock on the inner RLock().
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
			bal = acc.Balance
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

	for _, acc := range cs.GetAllAccounts() {
		addr := common.HexToAddress(acc.Address)
		decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		wei := new(big.Int).Mul(big.NewInt(int64(acc.Balance)), decimals)
		sdb.SetBalance(addr, wei)
		sdb.SetNonce(addr, cs.LoadNonce(acc.Address))
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
	_, err := cs.db.Exec(
		`INSERT INTO bio_registrations (commitment, wallet_address, tx_hash, bio_hash) VALUES ($1, $2, $3, $4)
 ON CONFLICT (commitment) DO UPDATE SET wallet_address = $2, tx_hash = $3, bio_hash = $4`,
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

func (cs *ChainState) SaveNullifier(nullifier, walletAddress string) {
	if nullifier == "" {
		return
	}
	walletAddress = strings.ToLower(walletAddress)
	cs.mu.Lock()
	cs.nullifiers[nullifier] = walletAddress
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

// ─── SWAP NONCES ─────────────────────────────────────────────────────────────
//
// Each wallet has a monotonically increasing nonce for swap/liquidity actions.
// The nonce is included in the signed message, so a captured signature cannot
// be replayed — the nonce check atomically rejects any second use.

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
