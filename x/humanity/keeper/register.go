package keeper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var registerRateLimit sync.Map

func init() {
// P3-4: periodically clean up expired rate-limit entries to prevent unbounded growth.
go func() {
for {
time.Sleep(60 * time.Second)
now := time.Now()
registerRateLimit.Range(func(k, v interface{}) bool {
if now.Sub(v.(time.Time)) > 11*time.Second {
registerRateLimit.Delete(k)
}
return true
})
}
}()
}

// isPrivateOrLoopback returns true for RFC-1918 private ranges and loopback
// addresses — used to decide whether to trust X-Forwarded-For (only safe when
// the direct connection comes from a known reverse-proxy, not the open internet).
func isPrivateOrLoopback(ipStr string) bool {
	parsed := net.ParseIP(ipStr)
	if parsed == nil {
		return false
	}
	for _, cidr := range []string{
		"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"::1/128", "fc00::/7",
	} {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err == nil && ipnet.Contains(parsed) {
			return true
		}
	}
	return false
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	// Only trust X-Forwarded-For when the TCP connection itself comes from a
	// private/loopback address — i.e. through Railway's or Render's proxy.
	// A direct internet client must not be able to spoof their IP via this header.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" && isPrivateOrLoopback(host) {
		first := strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
		if ip, _, err := net.SplitHostPort(first); err == nil {
			return ip
		}
		if first != "" {
			return first
		}
	}
	return host
}

// V7 ABI fragment for registerWithSig — used only to encode calldata correctly,
// including the dynamic `bytes signature` parameter (offset/length tail encoding).
const registerWithSigABI = `[{
"name": "registerWithSig",
"type": "function",
"inputs": [
{"name": "pA", "type": "uint256[2]"},
{"name": "pB", "type": "uint256[2][2]"},
{"name": "pC", "type": "uint256[2]"},
{"name": "pubSignals", "type": "uint256[2]"},
{"name": "claimedHuman", "type": "address"},
{"name": "signature", "type": "bytes"},
{"name": "nullifier", "type": "bytes32"}
]
}]`

const bioVerifierABI = `[{
"name": "verifyProof",
"type": "function",
"inputs": [
{"name": "pA", "type": "uint256[2]"},
{"name": "pB", "type": "uint256[2][2]"},
{"name": "pC", "type": "uint256[2]"},
{"name": "pubSignals", "type": "uint256[2]"}
],
"outputs": [{"name": "", "type": "bool"}],
"stateMutability": "view"
}]`

type RegisterRequest struct {
	Wallet     string     `json:"wallet"`
	PubSignals []string   `json:"pubSignals"`
	PA         []string   `json:"pA"`
	PB         [][]string `json:"pB"`
	PC         []string   `json:"pC"`
	Signature  string     `json:"signature"` // hex-encoded, 65 bytes, from personal_sign
	// BioHash is the raw bigint biometric hash from the device, used to
	// record the registration in bio_registrations so the app can poll for
	// its own registration by bioHash on startup.
	BioHash string `json:"bioHash"`
	// BioHashKey is the keccak256-hashed version of BioHash, as computed
	// by the proof server. Stored in bio_hashes so the proof server's
	// duplicate check (which also uses keccak256) can detect re-registrations
	// even when the user generates a new wallet on the same device.
	BioHashKey string `json:"bioHashKey"`
	// Nullifier is SHA256(bioHash + ":aequitas-ubi-v1") for v1 circuit, or
	// the hex representation of pubSignals[1] for v2 circuit (ZK-bound).
	Nullifier string `json:"nullifier"`
	// ZKNullifier is pubSignals[1] from the v2 circuit — the nullifier
	// derived INSIDE the ZK proof, making it cryptographically binding.
	// When present, it overrides the client-SHA256 nullifier.
	ZKNullifier   string `json:"zkNullifier"`
	CircuitVersion int   `json:"circuitVersion"`
}

type RegisterResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Balance float64 `json:"balance"`
	TxHash  string  `json:"tx_hash"`
}

func (a *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "POST required"})
		return
	}

	ip := clientIP(r)
	now := time.Now()
	if last, ok := registerRateLimit.Load(ip); ok && now.Sub(last.(time.Time)) < 10*time.Second {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "too many requests — please wait 10 seconds"})
		return
	}
	registerRateLimit.Store(ip, now)

	r.Body = http.MaxBytesReader(w, r.Body, 256<<10) // 256 KB — ZK proofs are large
	body, _ := io.ReadAll(r.Body)
	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "invalid request"})
		return
	}

	wallet := strings.ToLower(req.Wallet)
	if wallet == "" {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "wallet required"})
		return
	}
	if len(req.PA) < 2 || len(req.PB) < 2 || len(req.PC) < 2 || len(req.PubSignals) < 2 {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "incomplete ZK proof"})
		return
	}
	if req.Signature == "" {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "signature required"})
		return
	}

	// Early reject: wallet already registered — saves the expensive EVM call.
	if a.state.IsHuman(wallet) {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "wallet already registered"})
		return
	}

	fmt.Printf("[REGISTER] Relaying registerWithSig for: %s\n", wallet)

	// Use the shared EVMRPCServer so all parallel registrations share
	// the same nonce map + mutex — prevents two concurrent registrations
	// from reading the same DB-Nonce and writing the same follower value.
	evmRPC := a.evmRPC
	if evmRPC == nil || evmRPC.evm == nil {
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "EVM engine unavailable"})
		return
	}

	txHash, err := a.registerOnV7(evmRPC, wallet, req)
	if err != nil {
		fmt.Printf("[REGISTER] V7 registerWithSig failed: %v\n", err)
		json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: err.Error()})
		return
	}

	fmt.Printf("[REGISTER] ✓ Relayed registerWithSig for %s, tx=%s\n", wallet, txHash)
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: true,
		Message: "Registered as human on Aequitas V7! 1,000 AEQ granted.",
		Balance: 1000,
		TxHash:  txHash,
	})
}

// registerOnV7 builds, signs (with the relayer key), and submits a registerWithSig
// transaction on the V7 contract. claimedHuman in the contract call is the user's
// own wallet — verified via the signature, not via msg.sender (which is the relayer).
func (a *APIServer) registerOnV7(evmRPC *EVMRPCServer, wallet string, req RegisterRequest) (string, error) {
	relayerPK := os.Getenv("RELAYER_PRIVATE_KEY")
	if relayerPK == "" {
		return "", fmt.Errorf("server misconfigured: RELAYER_PRIVATE_KEY not set")
	}
	relayerPK = strings.TrimPrefix(relayerPK, "0x")

	relayerKey, err := crypto.HexToECDSA(relayerPK)
	if err != nil {
		return "", fmt.Errorf("invalid relayer key: %w", err)
	}
	relayerAddr := crypto.PubkeyToAddress(relayerKey.PublicKey)

	parsedABI, err := abi.JSON(strings.NewReader(registerWithSigABI))
	if err != nil {
		return "", fmt.Errorf("abi parse failed: %w", err)
	}

	pA, err := parseUint2(req.PA)
	if err != nil {
		return "", fmt.Errorf("invalid pA: %w", err)
	}
	pB, err := parseUint2x2(req.PB)
	if err != nil {
		return "", fmt.Errorf("invalid pB: %w", err)
	}
	pC, err := parseUint2(req.PC)
	if err != nil {
		return "", fmt.Errorf("invalid pC: %w", err)
	}
	pubSignals, err := parseUint2(req.PubSignals)
	if err != nil {
		return "", fmt.Errorf("invalid pubSignals: %w", err)
	}

	sigHex := strings.TrimPrefix(req.Signature, "0x")
	sigBytes := common.Hex2Bytes(sigHex)
	if len(sigBytes) != 65 {
		return "", fmt.Errorf("signature must be 65 bytes, got %d", len(sigBytes))
	}

	claimedHuman := common.HexToAddress(wallet)

	// Prefer ZK-circuit-derived nullifier (v2 circuit, pubSignals[1]) over
	// client-SHA256 nullifier. When the v2 circuit is used, the nullifier is
	// cryptographically attested by the Groth16 proof itself.
	effectiveNullifier := req.Nullifier
	if req.ZKNullifier != "" && req.CircuitVersion >= 2 {
		effectiveNullifier = req.ZKNullifier
		fmt.Printf("[REGISTER] Using ZK-bound nullifier (circuit v%d)\n", req.CircuitVersion)
	}
	// Nullifier is mandatory — reject registrations without one.
	if effectiveNullifier == "" {
		return "", fmt.Errorf("nullifier required")
	}
	if existingWallet := a.state.GetWalletByNullifier(effectiveNullifier); existingWallet != "" {
		return "", fmt.Errorf("identity already registered (nullifier used by %s)", existingWallet)
	}
	if req.BioHash != "" {
		if existingWallet := a.state.GetWalletByBioHash(req.BioHash); existingWallet != "" {
			return "", fmt.Errorf("biometric already registered to %s", existingWallet)
		}
		// P1-FIX: SHA256-derivation only applies to v1 circuit (nullifier=SHA256(bioHash)).
		// v2 circuit nullifiers are ZK-bound (pubSignals[1]) and NOT SHA256(bioHash).
		// Applying this check to a v2 ZKNullifier always fails, rejecting valid v2
		// registrations that arrive with a bioHash for duplicate-detection purposes.
		if req.CircuitVersion < 2 || req.ZKNullifier == "" {
			// Application-level nullifier binding: verify the nullifier is correctly
			// derived from the biometric hash via SHA256(bioHash+":aequitas-ubi-v1").
			// This prevents a client from submitting a valid ZK proof with an
			// arbitrary nullifier — the server independently recomputes and checks.
			h := sha256.Sum256([]byte(req.BioHash + ":aequitas-ubi-v1"))
			expectedNullifier := hex.EncodeToString(h[:])
			actualNullifier := strings.ToLower(strings.TrimPrefix(effectiveNullifier, "0x"))
			if expectedNullifier != actualNullifier {
				return "", fmt.Errorf("nullifier does not match biometric hash derivation")
			}
		}
	}

	// Encode nullifier as bytes32 (big-endian uint256).
	// v1 nullifiers are SHA256 hex strings ("0xabc..." or "abc...").
	// v2 nullifiers are pubSignals[1] — a decimal integer string like "17579322874185".
	// Both must be encoded as big-endian 32-byte integers for the contract.
	var nullifierBytes [32]byte
	if effectiveNullifier != "" {
		n := new(big.Int)
		s := strings.TrimPrefix(effectiveNullifier, "0x")
		if effectiveNullifier != s {
			// Had "0x" prefix → hex string (v1 circuit)
			n.SetString(s, 16)
		} else {
			// No prefix → try decimal first (v2 circuit), fall back to hex
			if _, ok := n.SetString(s, 10); !ok {
				n.SetString(s, 16)
			}
		}
		b := n.Bytes()
		if len(b) <= 32 {
			copy(nullifierBytes[32-len(b):], b) // right-align (big-endian)
		}
	}

	calldata, err := parsedABI.Pack("registerWithSig", pA, pB, pC, pubSignals, claimedHuman, sigBytes, nullifierBytes)
	if err != nil {
		return "", fmt.Errorf("encoding failed: %w", err)
	}

	to := common.HexToAddress(V7_CONTRACT_ADDR)

	// Dry-run first: simulate the call before spending a real transaction/nonce.
	// If the contract would revert (invalid signature, invalid proof, already
	// registered, etc.), CallContract returns an error and we abort here with
	// the real reason — instead of reporting false success to the caller.
	// persist=false is critical here: this is ONLY a simulation. Previously
	// CallContract always persisted its result regardless of intent, which
	// meant this dry-run alone — even for an attempt that was never actually
	// submitted, or whose real submission later failed — already wrote
	// isHuman=true/balance=1000 to evm_storage as a side effect of merely
	// checking whether registration would succeed.
	_, dryRunErr := evmRPC.evm.CallContract(relayerAddr, to, calldata, big.NewInt(0), false)
	if dryRunErr != nil {
		errStr := strings.ToLower(dryRunErr.Error())
		contractMissing := strings.Contains(errStr, "no code") ||
			strings.Contains(errStr, "empty code") ||
			strings.Contains(errStr, "contract not deployed")
		if !contractMissing {
			// Proof invalid, already registered, bad signature, etc. — surface
			// the real EVM revert reason instead of silently bypassing via mirror.
			return "", fmt.Errorf("registration rejected: %w", dryRunErr)
		}
		// V7 not yet deployed (startup race). Mirror validates the proof via
		// BioVerifier and writes storage slots directly. Only allowed here.
		fmt.Printf("[REGISTER] V7 not yet deployed — using mirror registration for %s\n", wallet)
		txHash, mirrorErr := a.persistRegisterWithSigMirror(evmRPC, to, claimedHuman, pA, pB, pC, pubSignals, sigBytes, calldata, effectiveNullifier)
		if mirrorErr != nil {
			return "", fmt.Errorf("registration failed (V7 missing + mirror failed): %w; mirror: %v", dryRunErr, mirrorErr)
		}
		if regErr := a.state.RegisterHuman(wallet); regErr != nil {
			fmt.Printf("[REGISTER] Warning: native balance grant failed after mirror: %v\n", regErr)
		}
		if len(req.PubSignals) > 0 {
			commitment := req.PubSignals[0]
			if saveErr := a.state.SaveBioRegistration(commitment, wallet, txHash, req.BioHash); saveErr != nil {
				fmt.Printf("[REGISTER] Warning: could not save bio registration link: %v\n", saveErr)
			}
		}
		if req.BioHashKey != "" {
			a.state.SaveBioHash(req.BioHashKey, wallet)
		} else if req.BioHash != "" {
			a.state.SaveBioHash(req.BioHash, wallet)
		}
		// Mirror path must also emit a block TX so secondary nodes learn about
		// this registration. Without this, secondaries never see mirror-path
		// registrations and their nullifier tables diverge, enabling double-spend.
		mirrorCommitment := ""
		if len(req.PubSignals) > 0 {
			mirrorCommitment = req.PubSignals[0]
		}
		if effectiveNullifier != "" {
			mirrorPAslice := []*big.Int{pA[0], pA[1]}
			mirrorPCslice := []*big.Int{pC[0], pC[1]}
			mirrorPSslice := []*big.Int{pubSignals[0], pubSignals[1]}
			a.blockchain.AddTransaction(Transaction{
				Type:       "register_human",
				Wallet:     wallet,
				TxHash:     txHash,
				Nullifier:  effectiveNullifier,
				Commitment: mirrorCommitment,
				ProofA:     bigIntsToHexStrings(mirrorPAslice),
				ProofB:     bigInt2x2ToHexStrings(pB),
				ProofC:     bigIntsToHexStrings(mirrorPCslice),
				PubSignals: bigIntsToHexStrings(mirrorPSslice),
			})
		}
		return txHash, nil
	}

	relayerAddrStr := strings.ToLower(relayerAddr.Hex())
	nonce := a.state.LoadNonce(relayerAddrStr)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(0),
		Gas:      6_000_000,
		GasPrice: big.NewInt(0),
		Data:     calldata,
	})

	signer := types.NewEIP155Signer(big.NewInt(1926))
	signedTx, err := types.SignTx(tx, signer, relayerKey)
	if err != nil {
		return "", fmt.Errorf("relayer signing failed: %w", err)
	}

	rawBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("tx encoding failed: %w", err)
	}

	rawHex := "0x" + common.Bytes2Hex(rawBytes)
	// Layer 2: actually submit. sendRawTransaction itself now returns a real
	// RPCError if the on-chain execution reverts (see evm_rpc.go), so a
	// non-nil rpcErr here already means the registration genuinely failed —
	// not just that submission failed.
	result, rpcErr := evmRPC.sendRawTransaction([]json.RawMessage{
		mustMarshal(rawHex),
	})
	if rpcErr != nil {
		return "", fmt.Errorf("registration failed on-chain: %s", rpcErr.Message)
	}

	txHash, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response from relay")
	}

	// ── NATIVE BALANCE (Phase 1 of native-coin migration) ────────────────
	// Previously AEQ only existed as an ERC20-style balanceOf() mapping
	// inside the V7 contract's EVM storage — eth_getBalance (the actual
	// native-coin query, what MetaMask would use for the chain's own gas
	// currency) read from a completely separate, never-updated ChainState/
	// chain_accounts table. That's why a wallet could show 1,000 AEQ in
	// MetaMask (via the custom ERC20 token display) while genuinely having
	// 0 native balance on the chain itself. This call makes ChainState the
	// real source of truth for the 1,000 AEQ grant — registration remains
	// fully gasless regardless of this; granting a native balance here has
	// no gas cost of its own, it's a direct database write, not a
	// transaction the user pays for.
	if regErr := a.state.RegisterHuman(wallet); regErr != nil {
		// P0-4: retry RegisterHuman — EVM succeeded but Go-State failed.
		// Without retry, the wallet has EVM balance but no native Go balance = permanent divergence.
		registered := false
		for retry := 1; retry <= 3; retry++ {
			time.Sleep(time.Duration(retry) * 500 * time.Millisecond)
			if err2 := a.state.RegisterHuman(wallet); err2 == nil {
				registered = true
				break
			}
		}
		if !registered {
			Log.Error("CRITICAL: RegisterHuman failed 3x after EVM success — Go/EVM diverged", "wallet", wallet, "error", regErr)
		}
		a.state.SyncBalancesToEVM(V7_CONTRACT_ADDR, wallet)
	} else {
		a.state.SyncBalancesToEVM(V7_CONTRACT_ADDR, wallet)
	}

	// Record which wallet this proof's commitment actually registered to,
	// so the app can later ask "did MY proof get registered, and where?"
	// instead of reading the last entry in a global, unfiltered list.
	if len(req.PubSignals) > 0 {
		commitment := req.PubSignals[0]
		if saveErr := a.state.SaveBioRegistration(commitment, wallet, txHash, req.BioHash); saveErr != nil {
			fmt.Printf("[REGISTER] Warning: could not save bio registration link: %v\n", saveErr)
		}
	}

	// Keep bio_hashes in sync so the proof server's /check and /prove
	// endpoints can block duplicate biometric registrations via that table.
	// Use the keccak256 bioHashKey when available — it matches the format
	// the proof server uses for its own duplicate checks.
	bioHashKey := req.BioHashKey
	if bioHashKey == "" {
		bioHashKey = req.BioHash
	}
	if bioHashKey != "" {
		a.state.SaveBioHash(bioHashKey, wallet)
		// Fire-and-forget: sync to proof server so its /prove duplicate check
		// is actually populated. Non-blocking — a slow proof server must not
		// delay the registration response.
		go notifyProofServer(bioHashKey, wallet)
	}
	// Always save effectiveNullifier (the ZK-derived one when using v2 circuit),
	// not just req.Nullifier — req.Nullifier may be empty or the old SHA256 value
	// while effectiveNullifier is the one actually stored on-chain.
	// TryClaimNullifier (called in RegisterHuman flow) already writes the nullifier
	// atomically. SaveNullifier here would overwrite the in-memory cache with a
	// potentially stale value from a concurrent race winner. Use effectiveNullifier
	// only for the block TX, not for a second write.
	nullifierToStore := effectiveNullifier
	if nullifierToStore == "" {
		nullifierToStore = req.Nullifier
	}

	// Add a register_human TX to the DAG so secondary nodes learn about this
	// registration via normal block sync and can apply it to their own state.
	// Nullifier + Commitment are included so the secondary can replay the full
	// state change (balance + nullifier + bio-registration) idempotently.
	commitment := ""
	if len(req.PubSignals) > 0 {
		commitment = req.PubSignals[0]
	}
	// Only emit block TX when we have a nullifier — secondary nodes use it as
	// the idempotency key. A TX with empty nullifier would be silently dropped
	// by replayRegistrations, hiding the registration from all secondary nodes.
	if nullifierToStore != "" {
		pAslice := []*big.Int{pA[0], pA[1]}
		pCslice := []*big.Int{pC[0], pC[1]}
		psSlice := []*big.Int{pubSignals[0], pubSignals[1]}
		a.blockchain.AddTransaction(Transaction{
			Type:       "register_human",
			Wallet:     wallet,
			TxHash:     txHash,
			Nullifier:  nullifierToStore,
			Commitment: commitment,
			ProofA:     bigIntsToHexStrings(pAslice),
			ProofB:     bigInt2x2ToHexStrings(pB),
			ProofC:     bigIntsToHexStrings(pCslice),
			PubSignals: bigIntsToHexStrings(psSlice),
		})
	}

	return txHash, nil
}

// notifyProofServer POSTs the registered bioHashKey to the proof server's
// /store-bio endpoint so its duplicate check stays in sync with the chain.
// Requires PROOF_SERVER_URL and CHAIN_SERVICE_TOKEN env vars on the chain node;
// if either is missing the call is skipped silently (registration already succeeded).
func notifyProofServer(bioHashKey, wallet string) {
	proofServerURL := os.Getenv("PROOF_SERVER_URL")
	if proofServerURL == "" {
		proofServerURL = "https://aequitas-proof-server-production.up.railway.app"
	}
	token := os.Getenv("CHAIN_SERVICE_TOKEN")
	if token == "" {
		return
	}
	body, _ := json.Marshal(map[string]string{"bioHashKey": bioHashKey, "wallet": wallet})
	req, err := http.NewRequest("POST", proofServerURL+"/store-bio", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("[REGISTER] Warning: could not build proof-server notify request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-chain-token", token)
	// C3-FIX: use redirect-blocking client to prevent SSRF via PROOF_SERVER_URL.
	proofClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := proofClient.Do(req)
	if err != nil {
		fmt.Printf("[REGISTER] Warning: proof-server /store-bio call failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("[REGISTER] Warning: proof-server /store-bio returned %d\n", resp.StatusCode)
	}
}

func (a *APIServer) persistRegisterWithSigMirror(evmRPC *EVMRPCServer, contractAddr, claimedHuman common.Address, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, pubSignals [2]*big.Int, sigBytes []byte, calldata []byte, nullifierHex string) (string, error) {
	if nullifierHex == "" {
		return "", fmt.Errorf("nullifier required")
	}
	// Mirror the V7 contract's binding check: for circuit v2, pubSignals[1]
	// IS the nullifier output. If the caller's nullifierHex doesn't match
	// pubSignals[1] the proof and the nullifier are from different sessions,
	// which would break the "one human, one nullifier" invariant.
	if pubSignals[1] != nil && pubSignals[1].Sign() > 0 {
		// Compare numerically — effectiveNullifier for v2 circuit is a decimal
		// string ("17579322874185"), not hex. String comparison would always fail.
		providedBig := new(big.Int)
		s := strings.TrimPrefix(nullifierHex, "0x")
		if nullifierHex != s {
			providedBig.SetString(s, 16) // had 0x prefix → hex
		} else if _, ok := providedBig.SetString(s, 10); !ok {
			providedBig.SetString(s, 16) // no prefix → try decimal, then hex
		}
		if providedBig.Cmp(pubSignals[1]) != 0 {
			return "", fmt.Errorf("nullifier mismatch: provided %s != circuit output %s", nullifierHex, pubSignals[1].Text(10))
		}
	}
	wallet := strings.ToLower(claimedHuman.Hex())
	if a.state.IsHuman(wallet) {
		return "", fmt.Errorf("already registered")
	}

	commitment := pubSignals[0]
	addrStr := strings.ToLower(contractAddr.Hex())
	usedSlot := mappingSlotBytes32(common.BigToHash(commitment), 7)
	used, err := a.state.LoadStorageSlot(addrStr, usedSlot.Hex())
	if err != nil {
		return "", err
	}
	if common.HexToHash(used) != (common.Hash{}) {
		return "", fmt.Errorf("commitment used")
	}

	// Check nullifier slot (slot 8) for replay protection.
	nullKey := common.HexToHash(strings.TrimPrefix(nullifierHex, "0x"))
	nullSlotCheck := mappingSlotBytes32(nullKey, 8)
	nullVal, _ := a.state.LoadStorageSlot(addrStr, nullSlotCheck.Hex())
	if common.HexToHash(nullVal) != (common.Hash{}) {
		return "", fmt.Errorf("nullifier already used")
	}

	if !validRegisterSignature(contractAddr, claimedHuman, commitment, nullifierHex, sigBytes) {
		return "", fmt.Errorf("invalid signature")
	}

	verifierABI, err := abi.JSON(strings.NewReader(bioVerifierABI))
	if err != nil {
		return "", err
	}
	verifyData, err := verifierABI.Pack("verifyProof", pA, pB, pC, pubSignals)
	if err != nil {
		return "", err
	}
	ret, err := evmRPC.evm.CallContract(claimedHuman, common.HexToAddress(BIO_VERIFIER_ADDR), verifyData, big.NewInt(0), false)
	if err != nil {
		return "", fmt.Errorf("invalid proof: %w", err)
	}
	if len(ret) != 32 || ret[31] != 1 {
		return "", fmt.Errorf("invalid proof")
	}

	initialGrant := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	totalSupply := loadSlotBig(a.state, addrStr, 0)
	totalHumans := loadSlotBig(a.state, addrStr, 1)
	ubiAccumulated := loadSlotBig(a.state, addrStr, 3)
	now := big.NewInt(int64(blockContext().Time))

	a.state.SaveStorageSlot(addrStr, common.BigToHash(big.NewInt(0)).Hex(), common.BigToHash(new(big.Int).Add(totalSupply, initialGrant)).Hex())
	a.state.SaveStorageSlot(addrStr, common.BigToHash(big.NewInt(1)).Hex(), common.BigToHash(new(big.Int).Add(totalHumans, big.NewInt(1))).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 4).Hex(), common.BigToHash(initialGrant).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 6).Hex(), common.HexToHash("0x01").Hex())
	a.state.SaveStorageSlot(addrStr, usedSlot.Hex(), common.HexToHash("0x01").Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 9).Hex(), common.BigToHash(commitment).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 10).Hex(), common.BigToHash(now).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 11).Hex(), common.BigToHash(now).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 12).Hex(), common.BigToHash(ubiAccumulated).Hex())

	// usedNullifiers[nullifier] = claimedHuman — slot 8 (bytes32 → address mapping)
	nullSlot := mappingSlotBytes32(nullKey, 8)
	addrVal := common.BigToHash(claimedHuman.Big())
	a.state.SaveStorageSlot(addrStr, nullSlot.Hex(), addrVal.Hex())

	txHash := crypto.Keccak256Hash(append(calldata, claimedHuman.Bytes()...)).Hex()
	fmt.Printf("[REGISTER] Native V7 mirror persisted for %s, tx=%s\n", wallet, txHash)
	return txHash, nil
}

func validRegisterSignature(contractAddr, claimedHuman common.Address, commitment *big.Int, nullifierHex string, sigBytes []byte) bool {
	nullKey := common.HexToHash(strings.TrimPrefix(nullifierHex, "0x"))
	packed := make([]byte, 0, 32+20+8+32+32)
	packed = append(packed, common.LeftPadBytes(big.NewInt(1926).Bytes(), 32)...)
	packed = append(packed, contractAddr.Bytes()...)
	packed = append(packed, []byte("register")...)
	packed = append(packed, common.LeftPadBytes(commitment.Bytes(), 32)...)
	packed = append(packed, nullKey.Bytes()...)
	messageHash := crypto.Keccak256Hash(packed)
	ethHash := accounts.TextHash(messageHash.Bytes())
	sig := append([]byte(nil), sigBytes...)
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	pub, err := crypto.SigToPub(ethHash, sig)
	if err != nil {
		return false
	}
	return bytes.Equal(crypto.PubkeyToAddress(*pub).Bytes(), claimedHuman.Bytes())
}

func loadSlotBig(state *ChainState, addr string, slot int64) *big.Int {
	value, _ := state.LoadStorageSlot(addr, common.BigToHash(big.NewInt(slot)).Hex())
	out := new(big.Int)
	if value != "" {
		out.SetBytes(common.HexToHash(value).Bytes())
	}
	return out
}

// uint256Max is the maximum value for a Solidity uint256 parameter.
var uint256Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

func parseUint2(values []string) ([2]*big.Int, error) {
	var out [2]*big.Int
	for i := 0; i < 2; i++ {
		n := new(big.Int)
		_, ok := n.SetString(values[i], 10)
		if !ok {
			return out, fmt.Errorf("invalid number at index %d: %s", i, values[i])
		}
		// P2-17: reject values that exceed uint256 — ABI-encode would silently truncate them.
		if n.Sign() < 0 || n.Cmp(uint256Max) > 0 {
			return out, fmt.Errorf("value at index %d exceeds uint256 range", i)
		}
		out[i] = n
	}
	return out, nil
}

func parseUint2x2(values [][]string) ([2][2]*big.Int, error) {
	var out [2][2]*big.Int
	for i := 0; i < 2; i++ {
		if len(values[i]) < 2 {
			return out, fmt.Errorf("row %d has fewer than 2 elements", i)
		}
		for j := 0; j < 2; j++ {
			n := new(big.Int)
			_, ok := n.SetString(values[i][j], 10)
			if !ok {
				return out, fmt.Errorf("invalid number at [%d][%d]: %s", i, j, values[i][j])
			}
			// P2-AUDIT: Reject values exceeding uint256 — ABI-encode would
			// silently truncate them, letting attackers pass oversize values
			// that appear to be valid ZK proof components but are not.
			if n.Sign() < 0 || n.Cmp(uint256Max) > 0 {
				return out, fmt.Errorf("value at [%d][%d] exceeds uint256 range", i, j)
			}
			out[i][j] = n
		}
	}
	return out, nil
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

// bigIntsToHexStrings converts a slice of *big.Int to their decimal string
// representations for JSON serialization in block transactions.
// Decimal (base-10) is used instead of hex to match the ZK proof wire format
// used by the client and the Groth16 verifier ABI.
func bigIntsToHexStrings(vals []*big.Int) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		if v == nil {
			out[i] = "0"
		} else {
			out[i] = v.Text(10)
		}
	}
	return out
}

// bigInt2x2ToHexStrings converts a [2][2]*big.Int matrix to [][]string
// for JSON serialization in block transactions.
func bigInt2x2ToHexStrings(vals [2][2]*big.Int) [][]string {
	out := make([][]string, 2)
	for i := 0; i < 2; i++ {
		out[i] = make([]string, 2)
		for j := 0; j < 2; j++ {
			if vals[i][j] == nil {
				out[i][j] = "0"
			} else {
				out[i][j] = vals[i][j].Text(10)
			}
		}
	}
	return out
}
