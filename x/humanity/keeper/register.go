package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

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
{"name": "signature", "type": "bytes"}
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

	evmRPC := NewEVMRPCServer(a.blockchain, a.state)
	if evmRPC.evm == nil {
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
	if req.BioHash != "" {
		if existingWallet := a.state.GetWalletByBioHash(req.BioHash); existingWallet != "" {
			return "", fmt.Errorf("biometric already registered to %s", existingWallet)
		}
	}

	calldata, err := parsedABI.Pack("registerWithSig", pA, pB, pC, pubSignals, claimedHuman, sigBytes)
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
		fmt.Printf("[REGISTER] EVM registerWithSig dry-run reverted (%v); validating through native V7 mirror\n", dryRunErr)
		txHash, mirrorErr := a.persistRegisterWithSigMirror(evmRPC, to, claimedHuman, pA, pB, pC, pubSignals, sigBytes, calldata)
		if mirrorErr != nil {
			return "", fmt.Errorf("registration would fail on-chain: %w; mirror validation failed: %v", dryRunErr, mirrorErr)
		}
		if regErr := a.state.RegisterHuman(wallet); regErr != nil {
			fmt.Printf("[REGISTER] Warning: native balance grant failed (mirror registration still succeeded): %v\n", regErr)
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
		fmt.Printf("[REGISTER] Warning: native balance grant failed (contract registration still succeeded): %v\n", regErr)
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[REGISTER] Warning: proof-server /store-bio call failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("[REGISTER] Warning: proof-server /store-bio returned %d\n", resp.StatusCode)
	}
}

func (a *APIServer) persistRegisterWithSigMirror(evmRPC *EVMRPCServer, contractAddr, claimedHuman common.Address, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, pubSignals [2]*big.Int, sigBytes []byte, calldata []byte) (string, error) {
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

	if !validRegisterSignature(contractAddr, claimedHuman, commitment, sigBytes) {
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
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 8).Hex(), common.BigToHash(commitment).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 9).Hex(), common.BigToHash(now).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 10).Hex(), common.BigToHash(now).Hex())
	a.state.SaveStorageSlot(addrStr, mappingSlot(claimedHuman.Bytes(), 11).Hex(), common.BigToHash(ubiAccumulated).Hex())

	txHash := crypto.Keccak256Hash(append(calldata, claimedHuman.Bytes()...)).Hex()
	fmt.Printf("[REGISTER] Native V7 mirror persisted for %s, tx=%s\n", wallet, txHash)
	return txHash, nil
}

func validRegisterSignature(contractAddr, claimedHuman common.Address, commitment *big.Int, sigBytes []byte) bool {
	packed := make([]byte, 0, 32+20+8+32)
	packed = append(packed, common.LeftPadBytes(big.NewInt(1926).Bytes(), 32)...)
	packed = append(packed, contractAddr.Bytes()...)
	packed = append(packed, []byte("register")...)
	packed = append(packed, common.LeftPadBytes(commitment.Bytes(), 32)...)
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

func parseUint2(values []string) ([2]*big.Int, error) {
	var out [2]*big.Int
	for i := 0; i < 2; i++ {
		n := new(big.Int)
		_, ok := n.SetString(values[i], 10)
		if !ok {
			return out, fmt.Errorf("invalid number at index %d: %s", i, values[i])
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
			out[i][j] = n
		}
	}
	return out, nil
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
