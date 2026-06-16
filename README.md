# AEQUITAS — Decentralized Human Currency

> *"Money exists because people exist. Nothing more, nothing less."*

[![Chain](https://img.shields.io/badge/Chain-Aequitas%20Chain%20V6-gold)](https://aequitas-production-9fba.up.railway.app)
[![Chain ID](https://img.shields.io/badge/Chain%20ID-1926-blue)](https://aequitas-production-9fba.up.railway.app/rpc)
[![EVM](https://img.shields.io/badge/EVM-Compatible-green)](https://aequitas-production-9fba.up.railway.app/rpc)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

---

## What is Aequitas?

Aequitas is the first monetary system where money supply is directly and mathematically tied to verified human existence.

**Every verified human receives exactly 1,000 AEQ upon registration — unconditionally, equally, permanently.**

```
Total Supply = Verified Active Humans × 1,000 AEQ
```

No pre-mine. No founder allocation. No investor round. No early adopter advantage. The first person to register and the billionth person to register receive identical amounts. This is not a policy choice — it is a mathematical law encoded in immutable smart contract code.

---

## Live Infrastructure

| Component | URL |
|-----------|-----|
| 🌐 Block Explorer | https://aequitas-production-9fba.up.railway.app |
| ⛓ RPC Endpoint | https://aequitas-production-9fba.up.railway.app/rpc |
| 🔒 Proof Server | https://aequitas-proof-server-production.up.railway.app |
| 🔗 Node 2 (Render) | https://aequitas-node-2.onrender.com |
| ⛓ Bootstrap Node | `/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R` |

## MetaMask Configuration

| Parameter | Value |
|-----------|-------|
| Network Name | Aequitas Chain |
| RPC URL | https://aequitas-production-9fba.up.railway.app/rpc |
| Chain ID | **1926** |
| Symbol | AEQ |
| Block Explorer | https://aequitas-production-9fba.up.railway.app |

---

## V6 Smart Contracts (Aequitas Chain — Chain ID 1926)

| Contract | Address |
|----------|---------|
| AequitasV6 (Main) | `0xA76cA3bf34F2Ae5dFA0608696627e42b81180488` |
| BioVerifier (Groth16) | `0xc369D27b49DE017d113Bbcb9A1884a9e745B6BE2` |
| AequitasV5 (Sepolia Legacy) | `0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5` |

---

## Core Principles

### 1. Proof of Humanity
Every AEQ holder must prove they are a unique living human through biometric verification. This is not optional — it is the foundation of the entire system.

- Fingerprint processed by Hardware Secure Element (HSE)
- Raw biometric data **never** leaves the device
- Groth16 Zero-Knowledge Proof generated on Proof Server
- Commitment hash stored permanently — no double registration possible
- **One human, one wallet, forever**

### 2. BlockDAG Architecture
Aequitas runs on a custom Layer 1 blockchain built from scratch in Go, with BlockDAG consensus:

- Multiple blocks can be produced simultaneously by different nodes
- Blocks are later merged into "merge blocks" with multiple parents
- Higher throughput, lower latency, better fault tolerance
- Block time: ~6 seconds average

### 3. V6 Economic Model

#### Proof of Alive
Inactive wallets are handled gracefully — not punished. After 2 years of inactivity, warnings are sent. After 2.5 years, AEQ moves to personal escrow (not lost). After 4 years of total inactivity, AEQ enters the UBI Pool.

#### Guardian System
Every verified human can appoint one Guardian (another verified human) who can confirm they are still alive — with **zero transaction rights**. Max 3 wards per Guardian. 7-day timelock on assignment prevents forced assignment under duress.

#### Demurrage
1% annual fee on any balance **above** fairShare. Money flows to UBI Pool — never deleted. Historical precedent: Wörgl, Austria (1932) reduced unemployment 25% in one year through demurrage currency.

#### Wealth Cap
Hard ceiling enforced from human #1:
| Phase | Humans | Cap |
|-------|--------|-----|
| 0 | 1–100 | 50× fairShare |
| 1 | 101–1,000 | 20× fairShare |
| 2 | 1,001–10,000 | 10× fairShare |
| 3 | 10,001–100,000 | 5× fairShare |
| 4 | 100,000+ | 3× fairShare |

Excess is instantly redistributed equally to ALL active humans.

#### Universal Basic Income
UBI from protocol economics — no taxation, no government, no political decision required.

Sources:
- Transaction fees (0.1%) → 20% to UBI Pool
- Wealth cap overflow → immediate equal redistribution
- Demurrage on excess balances
- Inactive wallet escrow after 4 years

#### No Algorithmic Inflation
The **only** event that creates new AEQ: a new verified human registers → 1,000 AEQ created. No mining rewards, no staking rewards, no protocol emissions. Supply is mathematically guaranteed to equal humans × 1,000.

---

## System Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Android App                        │
│    Hardware Secure Element · Fingerprint → ZKP      │
│    Biometric data never leaves device               │
└──────────────────────┬──────────────────────────────┘
                       │ biometric hash
┌──────────────────────▼──────────────────────────────┐
│              Proof Server (Node.js)                 │
│    Groth16 ZKP Generation · Sybil Protection        │
│    Proof Storage · Short ID System                  │
└──────────────────────┬──────────────────────────────┘
                       │ pA, pB, pC, pubSignals
┌──────────────────────▼──────────────────────────────┐
│           Aequitas Layer 1 (Go)                     │
│    Node 1 (Railway) ←── libp2p ──→ Node 2 (Render)  │
│    BlockDAG · EVM Engine · JSON-RPC                 │
│    PostgreSQL (shared persistent state)             │
└─────────────────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│           AequitasV6 Smart Contract                 │
│    BioVerifier (Groth16) · Wealth Cap               │
│    Guardian System · Demurrage · UBI Pool           │
└─────────────────────────────────────────────────────┘
```

---

## Technical Specifications

| Parameter | Value |
|-----------|-------|
| Language | Go 1.24 |
| Consensus | BlockDAG + Proof of Humanity |
| Block Time | ~6 seconds |
| Chain ID | 1926 (0x786) |
| EVM | Yes — go-ethereum engine |
| ZKP System | Groth16 / snarkjs / circom |
| ZKP Curve | BN128 (alt-bn128) |
| Proof Size | ~200 bytes |
| Verification Time | ~10ms |
| State Storage | PostgreSQL (persistent) |
| P2P Protocol | libp2p (Go) |
| Initial Grant | 1,000 AEQ per human |
| Transaction Fee | 0.1% |

---

## Repository Structure

```
aequitas-chain/
├── cmd/aequitasd/          — Node binary entry point
├── x/humanity/keeper/
│   ├── api.go              — HTTP API server
│   ├── api_html.go         — Explorer web UI (multilingual)
│   ├── evm_engine.go       — EVM execution (go-ethereum)
│   ├── evm_rpc.go          — JSON-RPC handler
│   ├── evm_storage.go      — Contract storage (PostgreSQL)
│   ├── evm_v6mirror.go     — V6 state mirror
│   ├── blockdag.go         — BlockDAG consensus
│   ├── p2p.go              — libp2p networking
│   ├── state.go            — Chain state + PostgreSQL
│   └── register.go         — ZKP registration handler
├── AequitasV6.sol          — V6 smart contract
├── AequitasV5.sol          — V5 legacy contract
├── BioVerifier.sol         — Groth16 ZKP verifier
├── Aequitas_Whitepaper_v1.0_EN.pdf
├── Aequitas_Whitepaper_v1.0_DE.pdf
└── README.md
```

---

## Registration Flow

1. **Android App** → fingerprint scan via Hardware Secure Element
2. **App** → derives biometric hash (never leaves device)
3. **App** → sends hash to Proof Server → Groth16 ZKP generated
4. **App** → stores proof on server → gets short Proof ID
5. **App** → opens MetaMask with `?proofId=xxx`
6. **Explorer** → loads proof via ID → auto-connects wallet
7. **Explorer** → submits `/api/register` with ZKP proof
8. **Node** → verifies ZKP → calls AequitasV6 contract
9. **Node** → mirrors state to PostgreSQL
10. **App** → polls until confirmed → shows 1,000 AEQ balance

---

## Roadmap

| Phase | Status | Description |
|-------|--------|-------------|
| 0 | ✅ | Smart contracts, ZKP system, Android app, Proof Server |
| 0+ | ✅ | Aequitas Layer 1 Chain (Go), BlockDAG, P2P, Block Explorer |
| V6 | ✅ | Sovereign chain, EVM engine, Proof of Alive, Guardian, Demurrage, UBI, Wealth Cap |
| 1 | 🔄 | APK release, grant applications, growing human registry |
| 2 | ⬜ | iOS app, Proof of Alive activation, Guardian system live |
| 3 | ⬜ | DEX, lending protocol, cross-chain bridges |
| 4 | ⬜ | Full decentralization, community governance |

---

## Why Aequitas?

Bitcoin's estimated Gini coefficient exceeds 0.85 — higher than any country on Earth. The top 1% of Bitcoin addresses control over 90% of all Bitcoin. The cryptocurrency that was supposed to democratize finance created the most extreme wealth concentration in human history.

Aequitas was created to answer: *"What would a cryptocurrency look like if designed from first principles to be fair to every human being?"*

The answer is surprisingly simple: **Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.**

---

## Links

- 🌐 [Block Explorer](https://aequitas-production-9fba.up.railway.app)
- 📄 [Whitepaper EN](Aequitas_Whitepaper_v1.0_EN.pdf)
- 📄 [Whitepaper DE](Aequitas_Whitepaper_v1.0_DE.pdf)
- 💻 [GitHub](https://github.com/hanoi96international-gif/Aequitas)
- 🔍 [V5 Etherscan](https://sepolia.etherscan.io/address/0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5)

---

*Aequitas launched June 2026 · Phase 0 · Chain ID 1926*
