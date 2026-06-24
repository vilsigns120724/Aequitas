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
			// P2-FIX: use explicit Close() not defer — defer fires at function return
			// (after the signing step), keeping the DB connection occupied unnecessarily.
			for rows.Next() {
				var commitment, wallet string
				rows.Scan(&commitment, &wallet)
				snap.BioRegistrations = append(snap.BioRegistrations, SnapshotBioRegistration{
					Commitment:    commitment,
					WalletAddress: wallet,
				})
			}
			rows.Close()
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
	if cs.TotalHumans() > 0 {
		fmt.Printf("[SNAPSHOT] DB already populated (%d humans) — skipping bootstrap import\n", cs.TotalHumans())
		return nil
	}
	return cs.downloadAndApplySnapshot(peerURL, expectedSignerHex)
}

func (cs *ChainState) downloadAndApplySnapshot(peerURL, expectedSignerHex string) error {
	client := &http.Client{Timeout: 60 * time.Second}

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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20))
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	var snap StateSnapshot
	if err := json.Unmarshal(body, &snap); err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	now := time.Now().Unix()
	if snap.Timestamp > now+60 {
		return fmt.Errorf("snapshot timestamp is in the future (%d seconds ahead)", snap.Timestamp-now)
	}
	if now-snap.Timestamp > 86400 {
		return fmt.Errorf("snapshot is too old (%d seconds)", now-snap.Timestamp)
	}

	// Signature verification is mandatory when BOOTSTRAP_SIGNER is configured.
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

	// Apply in-memory under lock, then persist outside lock to avoid deadlock.
	var accountsToPersist []*AccountState
	cs.mu.Lock()
	for _, acc := range snap.Accounts {
		acc.Address = strings.ToLower(acc.Address)
		cs.accounts[acc.Address] = acc
		accountsToPersist = append(accountsToPersist, acc)
	}
	if snap.Pool != nil {
		cs.pool = snap.Pool
	}
	for nullifier, wallet := range snap.Nullifiers {
		cs.nullifiers[nullifier] = wallet
	}
	cs.mu.Unlock()

	if cs.db != nil {
		for _, acc := range accountsToPersist {
			cs.saveAccountToDB(acc)
		}
		cs.savePoolToDB()
		for nullifier, wallet := range snap.Nullifiers {
			cs.db.Exec(
				`INSERT INTO nullifiers (nullifier, wallet_address) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
				nullifier, strings.ToLower(wallet),
			)
		}
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

	cs.MigrateEVMFromGoState(V7_CONTRACT_ADDR)

	fmt.Printf("[SNAPSHOT] ✓ Applied %d accounts, %d nullifiers, %d bio-registrations\n",
		len(snap.Accounts), len(snap.Nullifiers), len(snap.BioRegistrations))
	return nil
}

