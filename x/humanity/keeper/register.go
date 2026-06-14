package keeper

import (
"encoding/json"
"fmt"
"io"
"net/http"
"strings"
)

type RegisterRequest struct {
Bio    string `json:"bio"`
Salt   string `json:"salt"`
Wallet string `json:"wallet"`
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

// Check if already registered on this chain
if a.balances.GetBalance(wallet) > 0 {
json.NewEncoder(w).Encode(RegisterResponse{
Success: false,
Message: "already registered on Aequitas Chain",
Balance: a.balances.GetBalance(wallet),
})
return
}

// Verify proof via Proof Server
proofResp, err := http.Post(
"https://aequitas-proof-server-production.up.railway.app/prove",
"application/json",
strings.NewReader(string(body)),
)
if err != nil {
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: "proof server error"})
return
}
defer proofResp.Body.Close()

proofBody, _ := io.ReadAll(proofResp.Body)
var proofData map[string]interface{}
json.Unmarshal(proofBody, &proofData)

if proofResp.StatusCode != 200 {
errMsg := "proof verification failed"
if msg, ok := proofData["error"].(string); ok {
errMsg = msg
}
json.NewEncoder(w).Encode(RegisterResponse{Success: false, Message: errMsg})
return
}

// Grant 1000 AEQ - GASLESS
a.balances.Grant(wallet, InitialGrant)
a.keeper.RegisterHuman(wallet, "zkp-verified", 0)

txHash := fmt.Sprintf("0x%x", a.blockchain.LatestBlock().Height)

fmt.Printf("[REGISTER] ✓ Human registered: %s | 1000 AEQ granted (gasless)\n", wallet)

json.NewEncoder(w).Encode(RegisterResponse{
Success: true,
Message: "✓ Registered as human! 1,000 AEQ granted.",
Balance: InitialGrant,
TxHash:  txHash,
})
}
