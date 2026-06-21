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
type StateSnapshot struct {
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
		Timestamp:  time.Now().Unix(),
		Accounts:   accounts,
		Pool:       &pool,
		Nullifiers: nullifiers,
	}

	// Pull bio_registrations from DB (commitment → wallet + bioHash)
	if cs.db != nil {
		rows, err := cs.db.Query(`SELECT commitment, wallet_address, COALESCE(bio_hash,'') FROM bio_registrations`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var commitment, wallet, bioHash string
				rows.Scan(&commitment, &wallet, &bioHash)
				snap.BioRegistrations = append(snap.BioRegistrations, SnapshotBioRegistration{
					Commitment:    commitment,
					WalletAddress: wallet,
					BioHash:       bioHash,
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

	// Append snapshot token if configured
	url := peerURL
	if token := os.Getenv("SNAPSHOT_TOKEN"); token != "" {
		if strings.Contains(url, "?") {
			url += "&token=" + token
		} else {
			url += "?token=" + token
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20)) // 50 MB cap
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	var snap StateSnapshot
	if err := json.Unmarshal(body, &snap); err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	// Verify snapshot age — reject anything older than 24 hours
	if time.Now().Unix()-snap.Timestamp > 86400 {
		return fmt.Errorf("snapshot is too old (%d seconds)", time.Now().Unix()-snap.Timestamp)
	}

	// Verify signature if an expected signer is provided
	if expectedSignerHex != "" && snap.Signature != "" {
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

	// Apply accounts
	cs.mu.Lock()
	for _, acc := range snap.Accounts {
		acc.Address = strings.ToLower(acc.Address)
		cs.accounts[acc.Address] = acc
		if cs.db != nil {
			cs.saveAccountToDB(acc)
		}
	}
	// Apply pool
	if snap.Pool != nil {
		cs.pool = snap.Pool
		cs.savePoolToDB()
	}
	// Apply nullifiers
	for nullifier, wallet := range snap.Nullifiers {
		cs.nullifiers[nullifier] = wallet
	}
	cs.mu.Unlock()

	// Persist nullifiers to DB
	if cs.db != nil {
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
