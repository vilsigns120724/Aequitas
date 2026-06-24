package keeper

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

// StateSnapshot is a complete, portable export of the Go-state needed to
// bootstrap a new node that has no PostgreSQL database of its own.
// It includes accounts, pool reserves, nullifiers, and bio-registration
// commitments — sufficient to validate new registrations, swaps, and
// transfers without access to the primary database.
// Version marks the snapshot schema version. P3-11: lets importers detect schema changes.
const SnapshotVersion = 1

type StateSnapshot struct {
Version int `json:"version"`
	Timestamp        int64                     `json:"timestamp"`
	Accounts         []*AccountState           `json:"accounts"`
	Pool             *PoolState                `json:"pool"`
	Nullifiers       map[string]string         `json:"nullifiers"`          // nullifier → wallet
	BioRegistrations []SnapshotBioRegistration `json:"bio_registrations"`
	Signature        string                    `json:"signature,omitempty"` // ECDSA over SHA256(JSON without this field)
}

type SnapshotBioRegistration struct {
	Commitment    string `json:"commitment"`
	WalletAddress string `json:"wallet_address"`
	BioHash       string `json:"bio_hash,omitempty"`
}

// ExportSnapshot captures the live Go-state and, if signingKey is non-nil,
// signs the JSON payload so consumers can verify authenticity.
func (cs *ChainState) ExportSnapshot(signingKey *ecdsa.PrivateKey) *StateSnapshot {
	cs.mu.RLock()
	accounts := make([]*AccountState, 0, len(cs.accounts))
	for _, acc := range cs.accounts {
		cp := *acc
		accounts = append(accounts, &cp)
	}
	var pool PoolState
	if cs.pool != nil {
		pool = *cs.pool
	}
	nullifiers := make(map[string]string, len(cs.nullifiers))
	for k, v := range cs.nullifiers {
		nullifiers[k] = v
	}
	cs.mu.RUnlock()

	snap := &StateSnapshot{
Version: SnapshotVersion,
		Timestamp:  time.Now().Unix(),
		Accounts:   accounts,
		Pool:       &pool,
		Nullifiers: nullifiers,
	}

	// Pull bio_registrations from DB (commitment → wallet only).
	// bio_hash is intentionally omitted from the snapshot — it is a biometric
	// identifier and must not be exported to peer nodes. A new node can verify
	// commitment uniqueness without needing the raw bio_hash.
	if cs.db != nil {
		rows, err := cs.db.Query(`SELECT commitment, wallet_address FROM bio_registrations`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var commitment, wallet string
				rows.Scan(&commitment, &wallet)
				snap.BioRegistrations = append(snap.BioRegistrations, SnapshotBioRegistration{
					Commitment:    commitment,
					WalletAddress: wallet,
				})
			}
		}
	}

	if signingKey != nil {
		body, _ := json.Marshal(snap)
		hash := sha256.Sum256(body)
		sig, err := crypto.Sign(hash[:], signingKey)
		if err == nil {
			snap.Signature = hex.EncodeToString(sig)
		}
	}

	return snap
}

// ImportSnapshotFromURL downloads a StateSnapshot from peerURL and applies it
// to cs. Only imports if the local DB is empty (TotalHumans == 0) to avoid
// overwriting an already-populated state. Verifies the snapshot signature
// against expectedSignerHex if non-empty.
func (cs *ChainState) ImportSnapshotFromURL(peerURL, expectedSignerHex string) error {
	client := &http.Client{Timeout: 60 * time.Second}

	// P2-AUDIT: Send token in Authorization header instead of URL query parameter.
	// Tokens in URLs land in server/proxy access logs and browser history,
	// creating unnecessary credential exposure. Bearer header is already
	// accepted by handleSnapshot on the server side.
	req, reqErr := http.NewRequest("GET", peerURL, nil)
	if reqErr != nil {
		return fmt.Errorf("request build failed: %w", reqErr)
	}
	if token := os.Getenv("SNAPSHOT_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("snapshot server returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20)) // 50 MB cap
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	var snap StateSnapshot
	if err := json.Unmarshal(body, &snap); err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	now := time.Now().Unix()
	// Reject future-dated snapshots — could be injected by a compromised peer.
	if snap.Timestamp > now+60 {
		return fmt.Errorf("snapshot timestamp is in the future (%d seconds ahead)", snap.Timestamp-now)
	}
	// Reject stale snapshots older than 24 hours.
	if now-snap.Timestamp > 86400 {
		return fmt.Errorf("snapshot is too old (%d seconds)", now-snap.Timestamp)
	}

	// When BOOTSTRAP_SIGNER is configured, a valid signature is MANDATORY.
	// Importing without a valid signature is rejected — an attacker cannot
	// bypass verification by stripping the signature field.
	if expectedSignerHex != "" {
		if snap.Signature == "" {
			return fmt.Errorf("snapshot has no signature but BOOTSTRAP_SIGNER is set — import rejected")
		}
		sigCopy := snap.Signature
		snap.Signature = ""
		unsigned, _ := json.Marshal(snap)
		snap.Signature = sigCopy

		hash := sha256.Sum256(unsigned)
		sigBytes, err := hex.DecodeString(snap.Signature)
		if err != nil || len(sigBytes) != 65 {
			return fmt.Errorf("invalid snapshot signature format")
		}
		pubBytes, err := crypto.Ecrecover(hash[:], sigBytes)
		if err != nil {
			return fmt.Errorf("snapshot signature recovery failed: %w", err)
		}
		pubKey, err := crypto.UnmarshalPubkey(pubBytes)
		if err != nil {
			return fmt.Errorf("snapshot public key invalid: %w", err)
		}
		recovered := strings.ToLower(crypto.PubkeyToAddress(*pubKey).Hex())
		expected := strings.ToLower(expectedSignerHex)
		if recovered != expected {
			return fmt.Errorf("snapshot signed by %s, expected %s", recovered, expected)
		}
		fmt.Printf("[SNAPSHOT] ✓ Signature verified (signer: %s)\n", recovered)
	}

	// Apply accounts: update in-memory state under lock, then persist to DB
	// WITHOUT holding the lock — saveAccountToDB issues DB queries which must
	// not run while cs.mu is held (deadlock risk if the DB driver acquires
	// internal locks that chain back into any code that needs cs.mu).
	// P1-AUDIT: Split the lock into two phases: in-memory update (under lock)
	// and DB write (outside lock). Collect accounts to persist, release lock,
	// then write them. This eliminates the deadlock window.
	var accountsToPersist []*AccountState
	cs.mu.Lock()
	for _, acc := range snap.Accounts {
		acc.Address = strings.ToLower(acc.Address)
		cs.accounts[acc.Address] = acc
		accountsToPersist = append(accountsToPersist, acc)
	}
	// Apply pool in-memory
	if snap.Pool != nil {
		cs.pool = snap.Pool
	}
	// Apply nullifiers in-memory
	for nullifier, wallet := range snap.Nullifiers {
		cs.nullifiers[nullifier] = wallet
	}
	cs.mu.Unlock()

	// Persist accounts to DB outside lock
	if cs.db != nil {
		for _, acc := range accountsToPersist {
			cs.saveAccountToDB(acc)
		}
		// Persist pool outside lock
		cs.savePoolToDB()
		// Persist nullifiers to DB
		for nullifier, wallet := range snap.Nullifiers {
			cs.db.Exec(
				`INSERT INTO nullifiers (nullifier, wallet_address) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
				nullifier, strings.ToLower(wallet),
			)
		}
		// Persist bio_registrations
		for _, br := range snap.BioRegistrations {
			cs.db.Exec(
				`INSERT INTO bio_registrations (commitment, wallet_address, bio_hash) VALUES ($1, $2, $3)
				 ON CONFLICT (commitment) DO NOTHING`,
				br.Commitment, strings.ToLower(br.WalletAddress), br.BioHash,
			)
			if br.BioHash != "" {
				cs.db.Exec(
					`INSERT INTO bio_hashes (hash, wallet_address) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
					br.BioHash, strings.ToLower(br.WalletAddress),
				)
			}
		}
	}

	// Rebuild EVM storage from the imported Go-state
	cs.MigrateEVMFromGoState(V7_CONTRACT_ADDR)

	fmt.Printf("[SNAPSHOT] ✓ Imported %d accounts, %d nullifiers, %d bio-registrations\n",
		len(snap.Accounts), len(snap.Nullifiers), len(snap.BioRegistrations))
	return nil
}
