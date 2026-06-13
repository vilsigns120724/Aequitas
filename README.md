# AEQUITAS — Decentralized Human Currency

> "Money exists because people exist. Nothing more, nothing less."

[![Sepolia](https://img.shields.io/badge/Network-Sepolia%20Testnet-green)](https://sepolia.etherscan.io)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## Overview

Aequitas is the first decentralized monetary system that ties money supply directly to verified human existence. Every verified human receives 1,000 AEQ upon registration — unconditionally, equally, and permanently.

**Total Supply = Verified Humans × 1,000 AEQ**

## Live Contracts (Ethereum Sepolia)

| Contract | Address |
|----------|---------|
| AequitasV5 (ERC-20 + Full Economic Model) | `0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5` |
| Bio Verifier (Groth16) | `0x39Ac1431C94F6391B92d39615aB56B888Bbf2389` |
| AequitasV4 (ERC-20 + Bio ZKP) | `0x2B5ACedF2c41c70d51A2cbAd927b8940EE725DA7` |
| ZK Verifier V4 | `0x6502A5745Ca13d14cDe3E77EDa8b279fF3b72E0A` |

## Live Demo

**DApp:** https://hanoi96international-gif.github.io/Aequitas/aequitas-dapp.html

## What's New in V5

- 0.1% transaction fee (40% validators, 30% LPs, 20% UBI, 10% treasury)
- Algorithmic inflation 0–1.5% based on on-chain data only
- Dynamic wealth cap with waterfall redistribution
- Phased activation (Phase 0–4)
- No deflation — overflow redistributed instead
- Full economic model on-chain

## Economic Model

| Parameter | Value |
|-----------|-------|
| Initial Grant | 1,000 AEQ per verified human |
| Supply Cap | Humans × 1,000 AEQ |
| Transaction Fee | 0.1% |
| Max Inflation | 1.5%/year (algorithmic) |
| Wealth Cap (Phase 3) | 5× fairShare |
| Consensus (Phase 2) | Proof of Humanity (BlockDAG) |

## Wealth Cap Phases

| Phase | Registrations | Max Balance |
|-------|--------------|-------------|
| 0 | 0–100 | No cap |
| 1 | 100–1,000 | 20× fairShare |
| 2 | 1,000–10,000 | 10× fairShare |
| 3 | 10,000–100,000 | 5× fairShare |
| 4 | 100,000+ | Gini-dynamic 3–5× |

## Repository Structure

```
AequitasV5.sol          — Main contract (V5, full economic model)
AequitasV4.sol          — Previous version
bio_verifier.sol        — Groth16 ZKP verifier
biometric.circom        — ZKP circuit
biometric.wasm          — Compiled circuit
bio_0001.zkey           — Proving key
bio_verification_key.json — Verification key
aequitas-dapp.html      — Live Web DApp
Aequitas_Whitepaper_v0.9_DE.pdf
Aequitas_Whitepaper_v0.9_EN.pdf
```

## Biometric Identity

Registration requires a Zero-Knowledge Proof of biometric uniqueness:

```
commitment = ZKP( biometricHash × walletAddress + deviceSalt )
```

Biometric data never leaves the device. Only the mathematical proof is submitted on-chain.

**Phase 1 (upcoming):** PPG cardiac biometrics via MAX30102 sensor — believed to be the world's first application of PPG-based biometrics for decentralized identity.

## Roadmap

- ✅ Phase 0: Smart contracts, ZKP, Android app, Web DApp, Keeper bot
- 🔄 Phase 1: MAX30102 PPG sensor, AequitasV5, Cosmos SDK chain
- ⬜ Phase 2: Public testnet, BlockDAG, lending protocol, iOS app
- ⬜ Phase 3: Mainnet, DEX, cross-chain bridges

## Links

- Whitepaper: See PDF files in this repository
- Etherscan V5: https://sepolia.etherscan.io/address/0x4f147d5B3388AF07993CC4fC548502A78Af0B8b5
- Etherscan V4: https://sepolia.etherscan.io/address/0x2B5ACedF2c41c70d51A2cbAd927b8940EE725DA7
