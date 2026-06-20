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
