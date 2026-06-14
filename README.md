# AEQUITAS — Decentralized Human Currency

> "Money exists because people exist. Nothing more, nothing less."

[![Sepolia](https://img.shields.io/badge/Network-Sepolia%20Testnet-green)](https://sepolia.etherscan.io)
[![Chain](https://img.shields.io/badge/Layer1-Aequitas%20Chain-purple)](https://aequitas-production-9fba.up.railway.app)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## Overview

Aequitas is the first decentralized monetary system that ties money supply directly to verified human existence. Every verified human receives 1,000 AEQ upon registration — unconditionally, equally, and permanently.

**Total Supply = Verified Humans × 1,000 AEQ**

## Live Infrastructure

| Component | URL / Address |
|-----------|--------------|
| 🌐 Block Explorer | https://aequitas-production-9fba.up.railway.app |
| ⛓ Bootstrap Node | `/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R` |
| 🔗 Node 2 | https://aequitas-node-2.onrender.com |
| 🔒 Proof Server | https://aequitas-proof-server-production.up.railway.app |
| 📱 Android App | APK (private, coming to Play Store) |

## Live Contracts (Ethereum Sepolia)

| Contract | Address |
|----------|---------|
| AequitasV5 (ERC-20 + Full Economic Model) | `0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5` |
| Bio Verifier (Groth16) | `0x39Ac1431C94F6391B92d39615aB56B888Bbf2389` |

## What Makes Aequitas Unique

### 1. BlockDAG Consensus
Aequitas runs on a custom Layer 1 blockchain with **BlockDAG architecture** — multiple blocks can be produced in parallel and merged into a single block with multiple parents. This enables high throughput without sacrificing decentralization. Block #154 (MERGE)

├── Parent: a3f8c2... (Node 1 - Railway)

└── Parent: 9d2e1f... (Node 2 - Render)
### 2. Proof of Humanity
Registration requires a **Groth16 Zero-Knowledge Proof** of biometric uniqueness:
commitment = ZKP( biometricHash × walletAddress + deviceSalt )

Biometric data never leaves the device. Only the mathematical proof is submitted on-chain.

**Coming soon:** PPG cardiac biometrics via MAX30102 sensor — believed to be the world's first application of PPG-based biometrics for decentralized identity.

### 3. Algorithmic Monetary Policy
The **Aequitas Index** automatically measures economic health and adjusts monetary policy without human intervention:
Index = Velocity × 40% + Growth × 35% + (100 - Gini) × 25%
Index < 40  → Inflation (0–1.5% annual, equal distribution)

Index 40-60 → Neutral

Index > 60  → Wealth cap active, overflow redistributed

### 4. Fair Economics
- **0.1% transaction fee** → 40% validators / 30% LPs / 20% UBI / 10% treasury
- **Dynamic wealth cap** — prevents accumulation (activates at Phase 1+)
- **No deflation** — overflow redistributed equally to all humans

## Economic Model

| Parameter | Value |
|-----------|-------|
| Initial Grant | 1,000 AEQ per verified human |
| Supply Cap | Humans × 1,000 AEQ |
| Transaction Fee | 0.1% |
| Max Inflation | 1.5%/year (algorithmic) |
| Wealth Cap (Phase 3) | 5× fairShare |
| Block Time | 6 seconds |
| Consensus | Proof of Humanity + BlockDAG |

## Architecture
┌─────────────────────────────────────────────┐

│              Aequitas Layer 1               │

│         BlockDAG · Go · libp2p              │

│                                             │

│  Node 1 (Railway) ←──P2P──→ Node 2 (Render)│

│       ↕ HTTP Block Sync ↕                   │

└─────────────────────────────────────────────┘

↕ Sync

┌─────────────────────────────────────────────┐

│           Ethereum Sepolia                  │

│    AequitasV5 Smart Contract (ERC-20)       │

│    Groth16 Bio Verifier                     │

└─────────────────────────────────────────────┘

↕

┌─────────────────────────────────────────────┐

│           Proof Server (Railway)            │

│    Groth16 ZKP Generation                   │

│    Sybil Protection                         │

│    Keeper Bot (hourly cycles)               │

└─────────────────────────────────────────────┘

↕

┌─────────────────────────────────────────────┐

│           Android App                       │

│    Hardware Secure Element                  │

│    Fingerprint → ZKP → MetaMask             │

└─────────────────────────────────────────────┘

## Wealth Cap Phases

| Phase | Registrations | Max Balance |
|-------|--------------|-------------|
| 0 | 0–100 | No cap |
| 1 | 100–1,000 | 20× fairShare |
| 2 | 1,000–10,000 | 10× fairShare |
| 3 | 10,000–100,000 | 5× fairShare |
| 4 | 100,000+ | Gini-dynamic 3–5× |

## Repository Structure
AequitasV5.sol               — Main contract (V5, full economic model)

bio_verifier.sol             — Groth16 ZKP verifier

biometric.circom             — ZKP circuit

biometric.wasm               — Compiled circuit

bio_0001.zkey                — Proving key

bio_verification_key.json    — Verification key

cmd/aequitasd/               — Layer 1 node binary

x/humanity/keeper/           — BlockDAG, P2P, API, Sync

Aequitas_Whitepaper_v0.9_DE.pdf

Aequitas_Whitepaper_v0.9_EN.pdf

## Roadmap

- ✅ Phase 0: Smart contracts, ZKP, Android app, Proof Server, Keeper bot
- ✅ Phase 0+: Layer 1 Chain (Go), BlockDAG, P2P Network, Block Explorer
- 🔄 Phase 1: MAX30102 PPG sensor integration, iOS app, Grant application
- ⬜ Phase 2: Public testnet, lending protocol, DEX
- ⬜ Phase 3: Mainnet, cross-chain bridges

## Links

- 🌐 [Block Explorer](https://aequitas-production-9fba.up.railway.app)
- 📄 [Whitepaper DE](https://github.com/hanoi96international-gif/Aequitas/raw/main/Aequitas_Whitepaper_v0.9_DE.pdf)
- 📄 [Whitepaper EN](https://github.com/hanoi96international-gif/Aequitas/raw/main/Aequitas_Whitepaper_v0.9_EN.pdf)
- 🔍 [Etherscan V5](https://sepolia.etherscan.io/address/0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5)
- 💻 [GitHub](https://github.com/hanoi96international-gif/Aequitas)
