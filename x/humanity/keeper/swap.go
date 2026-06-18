package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// verifyPersonalSign checks that signature is a valid personal_sign
// signature of message, produced by the private key behind claimedWallet.
// This mirrors what MetaMask's personal_sign does client-side: it hashes
// "\x19Ethereum Signed Message:\n" + len(message) + message, then
// recovers the signer's address from the signature and compares it
// against claimedWallet.
//
// This exists because, unlike registerWithSig (which runs through the V7
// contract and lets the contract itself verify the signature on-chain),
// the swap/liquidity/faucet actions below are pure ChainState operations
// with no contract involved — so this Go code has to do the signature
// check itself. Without it, anyone could POST {"wallet": "someone else's
// address", ...} and act on a wallet they don't actually control.
func verifyPersonalSign(message, signatureHex, claimedWallet string) error {
	sigHex := strings.TrimPrefix(signatureHex, "0x")
	sigBytes := common.Hex2Bytes(sigHex)
	if len(sigBytes) != 65 {
		return fmt.Errorf("signature must be 65 bytes, got %d", len(sigBytes))
	}

	// go-ethereum's recovery expects the V value (last byte) to be 0 or 1,
	// but personal_sign / eth_sign in wallets produce 27 or 28 (EIP-191
	// convention) — normalize before recovery.
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	prefixed := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message))
	hash := crypto.Keccak256(prefixed)

	pubKey, err := crypto.SigToPub(hash, sigBytes)
	if err != nil {
		return fmt.Errorf("could not recover public key: %w", err)
	}

	recovered := strings.ToLower(crypto.PubkeyToAddress(*pubKey).Hex())
	claimed := strings.ToLower(claimedWallet)
	if recovered != claimed {
		return fmt.Errorf("signature does not match claimed wallet (recovered %s)", recovered)
	}
	return nil
}

// ── SWAP ──────────────────────────────────────────────────────────────────

type SwapRequest struct {
	Wallet    string  `json:"wallet"`
	Direction string  `json:"direction"` // "aeq_to_tusd" or "tusd_to_aeq"
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"` // personal_sign over a fixed message, see below
}

type SwapResponse struct {
	Success    bool    `json:"success"`
	Message    string  `json:"message"`
	AmountOut  float64 `json:"amount_out"`
	NewAEQ     float64 `json:"new_aeq_balance"`
	NewTUSD    float64 `json:"new_tusd_balance"`
}

func (a *APIServer) handleSwap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: "POST required"})
		return
	}

	body, _ := io.ReadAll(r.Body)
	var req SwapRequest
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: "invalid request"})
		return
	}

	wallet := strings.ToLower(req.Wallet)
	if wallet == "" || req.Amount <= 0 {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: "wallet and positive amount required"})
		return
	}
	if req.Direction != "aeq_to_tusd" && req.Direction != "tusd_to_aeq" {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: "direction must be aeq_to_tusd or tusd_to_aeq"})
		return
	}

	// The signed message is fixed and predictable from the request fields
	// themselves, so the wallet owner is explicitly confirming THIS exact
	// swap (amount + direction) — not just proving generic wallet
	// ownership, which could otherwise be replayed against a different
	// amount.
	message := fmt.Sprintf("Aequitas Swap: %s %.8f", req.Direction, req.Amount)
	if err := verifyPersonalSign(message, req.Signature, wallet); err != nil {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: "signature invalid: " + err.Error()})
		return
	}

	var amountOut float64
	var err error
	if req.Direction == "aeq_to_tusd" {
		amountOut, err = a.state.SwapAEQForTUSD(wallet, req.Amount)
	} else {
		amountOut, err = a.state.SwapTUSDForAEQ(wallet, req.Amount)
	}
	if err != nil {
		json.NewEncoder(w).Encode(SwapResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(SwapResponse{
		Success:   true,
		Message:   "swap successful",
		AmountOut: amountOut,
		NewAEQ:    a.state.GetBalance(wallet),
		NewTUSD:   a.state.GetTUsdBalance(wallet),
	})
}

// ── ADD LIQUIDITY ────────────────────────────────────────────────────────

type AddLiquidityRequest struct {
	Wallet     string  `json:"wallet"`
	AmountAEQ  float64 `json:"amount_aeq"`
	AmountTUSD float64 `json:"amount_tusd"`
	Signature  string  `json:"signature"`
}

type AddLiquidityResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (a *APIServer) handleAddLiquidity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(AddLiquidityResponse{Success: false, Message: "POST required"})
		return
	}

	body, _ := io.ReadAll(r.Body)
	var req AddLiquidityRequest
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(AddLiquidityResponse{Success: false, Message: "invalid request"})
		return
	}

	wallet := strings.ToLower(req.Wallet)
	if wallet == "" || req.AmountAEQ <= 0 || req.AmountTUSD <= 0 {
		json.NewEncoder(w).Encode(AddLiquidityResponse{Success: false, Message: "wallet and positive amounts required"})
		return
	}

	message := fmt.Sprintf("Aequitas Add Liquidity: %.8f AEQ + %.8f tUSD", req.AmountAEQ, req.AmountTUSD)
	if err := verifyPersonalSign(message, req.Signature, wallet); err != nil {
		json.NewEncoder(w).Encode(AddLiquidityResponse{Success: false, Message: "signature invalid: " + err.Error()})
		return
	}

	if err := a.state.AddLiquidity(wallet, req.AmountAEQ, req.AmountTUSD); err != nil {
		json.NewEncoder(w).Encode(AddLiquidityResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(AddLiquidityResponse{Success: true, Message: "liquidity added"})
}

// ── FAUCET ───────────────────────────────────────────────────────────────

type FaucetRequest struct {
	Wallet    string `json:"wallet"`
	Signature string `json:"signature"`
}

type FaucetResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Granted float64 `json:"granted"`
}

func (a *APIServer) handleFaucet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(FaucetResponse{Success: false, Message: "POST required"})
		return
	}

	body, _ := io.ReadAll(r.Body)
	var req FaucetRequest
	if err := json.Unmarshal(body, &req); err != nil {
		json.NewEncoder(w).Encode(FaucetResponse{Success: false, Message: "invalid request"})
		return
	}

	wallet := strings.ToLower(req.Wallet)
	if wallet == "" {
		json.NewEncoder(w).Encode(FaucetResponse{Success: false, Message: "wallet required"})
		return
	}

	message := fmt.Sprintf("Aequitas tUSD Faucet Claim: %s", wallet)
	if err := verifyPersonalSign(message, req.Signature, wallet); err != nil {
		json.NewEncoder(w).Encode(FaucetResponse{Success: false, Message: "signature invalid: " + err.Error()})
		return
	}

	if err := a.state.ClaimTUsdFaucet(wallet); err != nil {
		json.NewEncoder(w).Encode(FaucetResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(FaucetResponse{Success: true, Message: "faucet claimed", Granted: tusdFaucetAmount})
}

// ── POOL STATUS ──────────────────────────────────────────────────────────

func (a *APIServer) handlePoolStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reserveAEQ, reserveTUSD := a.state.GetPoolReserves()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reserve_aeq":  reserveAEQ,
		"reserve_tusd": reserveTUSD,
		"price_aeq_in_tusd": func() float64 {
			if reserveAEQ == 0 {
				return 0
			}
			return reserveTUSD / reserveAEQ
		}(),
	})
}
