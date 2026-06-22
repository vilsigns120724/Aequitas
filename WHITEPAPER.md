# AEQUITAS WHITEPAPER v2.0

**Proof of Humanity Chain — Eine faire Währung für alle Menschen**
**Proof of Humanity Chain — A Fair Currency for All of Humanity**

*Version 2.0 · Juni / June 2026*
*Chain ID 1926 · aequitas.digital*

---

## Inhalt / Table of Contents

1. [Das Problem / The Problem](#1-das-problem--the-problem)
2. [Die Vision / The Vision](#2-die-vision--the-vision)
3. [Proof of Humanity](#3-proof-of-humanity)
4. [Tokenomics & Wirtschaftsmodell / Economic Model](#4-tokenomics--wirtschaftsmodell--economic-model)
5. [Technische Architektur / Technical Architecture](#5-technische-architektur--technical-architecture)
6. [Smart Contract V7](#6-smart-contract-v7)
7. [Zero-Knowledge-Proofs & Privatsphäre / Privacy](#7-zero-knowledge-proofs--privatsphäre--privacy)
8. [Gleichheitsindex / Equality Index](#8-gleichheitsindex--equality-index)
9. [Exchange & Liquiditätspool / Liquidity Pool](#9-exchange--liquiditätspool--liquidity-pool)
10. [Sicherheit / Security](#10-sicherheit--security)
11. [Roadmap](#11-roadmap)
12. [Fazit / Conclusion](#12-fazit--conclusion)

---

## 1. Das Problem / The Problem

### DE
Bitcoin hat einen Gini-Koeffizienten von über 0,85 — höher als jedes Land der Erde. Die Top 1% der Adressen kontrollieren mehr als 90% aller Bitcoin. Was als dezentrales, demokratisches Geld begann, hat die extremste Vermögenskonzentration der Finanzgeschichte erschaffen.

Das ist kein Versagen der Blockchain-Technologie. Es ist ein Versagen des Designs: Bitcoin wurde ohne Rücksicht auf die Frage entworfen, wer Zugang zu initialem Kapital hat. Wer früh dabei war oder über Rechenleistung verfügte, gewann. Wer später kam oder arm war, verlor.

Das gleiche Muster wiederholt sich bei allen PoW- und PoS-Kryptowährungen: Die Reichen werden reicher, weil sie mehr Kapital einsetzen können. Das ist kein Bug — es ist die Systemarchitektur.

### EN
Bitcoin has a Gini coefficient exceeding 0.85 — higher than any country on Earth. The top 1% of addresses control more than 90% of all Bitcoin. What started as decentralized, democratic money created the most extreme wealth concentration in financial history.

This is not a failure of blockchain technology. It is a failure of design: Bitcoin was built without considering who has access to initial capital. Those who were early or had computing power won. Those who came later or were poor lost.

The same pattern repeats across all PoW and PoS cryptocurrencies: the rich get richer because they can deploy more capital. This is not a bug — it is the system architecture.

---

## 2. Die Vision / The Vision

### DE
Aequitas stellt eine radikale Frage: **Was wäre eine Kryptowährung, wenn sie von Grund auf fair für jeden Menschen auf der Erde konzipiert worden wäre?**

Die Antwort ist überraschend einfach:

> *Geld existiert, weil Menschen existieren. Daher sollte jede Person einen gleichen Anteil am Geld haben — allein weil sie ein Mensch ist.*

Aequitas setzt dieses Prinzip mathematisch um:

```
Gesamtangebot = Verifizierte Menschen × 1.000 AEQ
```

Kein Pre-Mine. Keine Gründer-Zuteilung. Kein früher Vorteil. Wer sich heute registriert und wer sich in zehn Jahren registriert, erhält exakt dasselbe. Es ist kein politisches Versprechen — es ist Code.

### EN
Aequitas asks a radical question: **What would a cryptocurrency look like if designed from first principles to be fair to every human being on Earth?**

The answer is surprisingly simple:

> *Money exists because people exist. Therefore, every person should have an equal share of money simply by virtue of being human.*

Aequitas implements this principle mathematically:

```
Total Supply = Verified Humans × 1,000 AEQ
```

No pre-mine. No founder allocation. No early-adopter advantage. Someone registering today and someone registering in ten years receive exactly the same. This is not a political promise — it is code.

---

## 3. Proof of Humanity

### DE
Das zentrale Problem eines auf menschlicher Existenz basierenden Währungssystems ist die Verifikation: Wie beweist man, dass eine Adresse einem echten, einzigartigen Menschen gehört — ohne persönliche Daten zu speichern?

Aequitas löst dies mit biometrischer Verifikation und Zero-Knowledge-Proofs:

**Registrierungsablauf:**
1. Die Android-App scannt den Fingerabdruck via **Hardware Secure Element (HSE)**
2. Das HSE leitet einen deterministischen biometrischen Hash ab — die Rohdaten verlassen das Gerät **niemals**
3. Der Hash wird an den Proof-Server gesendet
4. Der Proof-Server generiert einen **Groth16 Zero-Knowledge-Proof** (Groth16/BN128-Kurve)
5. Der Proof enthält `commitment` (Einmaligkeit-Nachweis) und `nullifier` (Replay-Schutz)
6. Die Blockchain verifiziert den Proof on-chain via BioVerifier-Contract
7. Bei Erfolg: 1.000 AEQ werden der Wallet gutgeschrieben

**Garantien:**
- Ein Mensch kann sich nur einmal registrieren (Nullifier-Bindung)
- Kein persönliches Datum wird gespeichert
- Verifikation ist dauerhaft und unveränderbar
- Kein Dritter (auch nicht Aequitas) kann eine Registrierung rückgängig machen

### EN
The central problem of a monetary system based on human existence is verification: how do you prove that an address belongs to a real, unique human — without storing personal data?

Aequitas solves this with biometric verification and Zero-Knowledge Proofs:

**Registration Flow:**
1. The Android app scans the fingerprint via **Hardware Secure Element (HSE)**
2. The HSE derives a deterministic biometric hash — raw data **never** leaves the device
3. The hash is sent to the Proof Server
4. The Proof Server generates a **Groth16 Zero-Knowledge Proof** (Groth16/BN128 curve)
5. The proof contains `commitment` (uniqueness proof) and `nullifier` (replay protection)
6. The blockchain verifies the proof on-chain via BioVerifier contract
7. On success: 1,000 AEQ credited to the wallet

**Guarantees:**
- A human can only register once (nullifier binding)
- No personal data is stored
- Verification is permanent and immutable
- No third party (not even Aequitas) can reverse a registration

---

## 4. Tokenomics & Wirtschaftsmodell / Economic Model

### 4.1 Geldangebot / Money Supply

Das Angebot ist eine mathematische Funktion menschlicher Existenz. Es gibt genau **eine** Möglichkeit neues AEQ zu erschaffen: ein neuer verifizierter Mensch registriert sich.

The supply is a mathematical function of human existence. There is exactly **one** way to create new AEQ: a new verified human registers.

| Ereignis / Event | AEQ-Änderung / AEQ Change |
|-----------------|--------------------------|
| Neue Registrierung / New Registration | +1.000 AEQ |
| Transfer | ±0 (nur Umverteilung / redistribution only) |
| Swap | ±0 |
| Demurrage | −x von Überschuss / from excess → UBI Pool |
| Wealth Cap Overflow | −x → sofortige Gleichverteilung / instant redistribution |

### 4.2 Universal Basic Income (UBI)

UBI wird täglich aus den Protokoll-Einnahmen verteilt — kein Staat, keine Steuer, kein Beschluss erforderlich.

UBI is distributed daily from protocol revenue — no state, no tax, no vote required.

**UBI-Pool-Quellen / UBI Pool Sources:**
- 20% aller Transaktionsgebühren (0,1% × 20% = 0,02% pro Transfer)
- Wealth-Cap-Überläufe (sofortige Gleichverteilung)
- Demurrage auf Überschussguthaben (1%/Jahr über fairShare)
- Inaktive Wallets: nach 4 Jahren Inaktivität → Escrow → UBI Pool

### 4.3 Demurrage — Haltegebühr

**Philosophie:** Geld ist ein Werkzeug, kein Selbstzweck. Horten von Geld über dem fairen Anteil kostet etwas — genau wie das Mieten eines Parkplatzes.

**Philosophy:** Money is a tool, not an end in itself. Hoarding money above the fair share costs something — just like renting a parking space.

```
Haltegebühr = (Guthaben − fairShare) × 1% × Haltezeit/Jahr
Demurrage   = (Balance − fairShare) × 1% × HoldingTime/Year
```

Die Gebühr fließt vollständig in den UBI-Pool. Kein AEQ wird vernichtet.
The fee flows entirely into the UBI Pool. No AEQ is destroyed.

### 4.4 Wealth Cap — Vermögensobergrenze

Eine dynamische Obergrenze verhindert extreme Konzentration und sinkt automatisch mit wachsender Gemeinschaft:

A dynamic ceiling prevents extreme concentration and automatically decreases as the community grows:

| Phase | Größe / Size | Obergrenze / Cap |
|-------|-------------|-----------------|
| 0 | 1–100 Menschen / humans | 50× fairShare = 50.000 AEQ |
| 1 | 101–1.000 | 20× fairShare = 20.000 AEQ |
| 2 | 1.001–10.000 | 10× fairShare = 10.000 AEQ |
| 3 | 10.001–100.000 | 5× fairShare = 5.000 AEQ |
| 4 | 100.000+ | 3× fairShare = 3.000 AEQ |

Überschuss wird sofort gleichmäßig unter **allen** aktiven Menschen verteilt.
Excess is instantly distributed equally among **all** active humans.

### 4.5 Transaktionsgebühren / Transaction Fees

| Empfänger / Recipient | Anteil / Share |
|----------------------|---------------|
| Validators | 40% |
| Liquidity Providers | 30% |
| UBI Pool | 20% |
| Treasury | 10% |

---

## 5. Technische Architektur / Technical Architecture

### 5.1 Layer 1 — Aequitas Chain

Aequitas läuft auf einer eigens entwickelten Layer-1-Blockchain, geschrieben in **Go 1.24**, mit einem hybriden BlockDAG-Konsensus.

Aequitas runs on a custom-built Layer 1 blockchain written in **Go 1.24**, with a hybrid BlockDAG consensus.

**BlockDAG:**
- Mehrere Blöcke können gleichzeitig von verschiedenen Nodes produziert werden
- Blöcke werden später in Merge-Blöcke zusammengeführt (mehrere Eltern)
- Höherer Durchsatz, niedrigere Latenz, bessere Fehlertoleranz
- Multiple blocks can be produced simultaneously by different nodes
- Blocks are merged into merge blocks with multiple parents
- Higher throughput, lower latency, better fault tolerance

**Dual-Ledger:**
Aequitas führt zwei synchronisierte Ledger parallel:
- **Go-Ledger**: PostgreSQL-gesichert, primäre Wahrheit für Salden und Menschen
- **EVM-Ledger**: go-ethereum Engine, kompatibel mit MetaMask und Web3

Aequitas maintains two synchronized ledgers in parallel:
- **Go-Ledger**: PostgreSQL-backed, primary truth for balances and humans
- **EVM-Ledger**: go-ethereum engine, compatible with MetaMask and Web3

### 5.2 Netzwerk-Topologie / Network Topology

```
Node 1 (Railway, Berlin)          Node 2 (Render, Frankfurt)
├── Primärer API-Server           ├── Sekundärer API-Server
├── Block-Produzent               ├── Block-Produzent
├── UBI-Verteilung (täglich)      ├── P2P-Peer
├── P2P Bootstrap-Node            └── HTTP Block-Sync
└── Geteilter PostgreSQL State ───────────────────────────┘
```

### 5.3 Technische Kenndaten / Technical Specifications

| Parameter | Wert / Value |
|-----------|-------------|
| Programmiersprache / Language | Go 1.24 |
| Konsens / Consensus | BlockDAG + Proof of Humanity |
| Blockzeit / Block Time | ~6 Sekunden / seconds |
| Chain ID | 1926 |
| EVM-Kompatibilität / EVM Compat. | Vollständig / Full (go-ethereum) |
| P2P-Protokoll / P2P Protocol | libp2p |
| State-Storage / State Storage | PostgreSQL (persistent) |
| ZKP-System / ZKP System | Groth16 / snarkjs / circom |
| Elliptische Kurve / Elliptic Curve | BN128 (alt-bn128) |
| Bio-Hash | keccak256 |
| Dezimalgenauigkeit / Precision | 6 Stellen / decimal places (1 AEQ = 1.000.000 Micro-AEQ) |

---

## 6. Smart Contract V7

### DE
Der AequitasV7-Contract ist das Herzstück des Protokolls. Er ist in Solidity geschrieben, auf der Aequitas Chain deployed und enthält die gesamte Wirtschaftslogik.

**Kernfunktionen:**
- `register()` — Registrierung mit ZKP-Beweis (direkt)
- `registerWithSig()` — Registrierung via Relayer (gaslos für den Nutzer)
- `transfer()` — Token-Transfer mit automatischer Demurrage und Gebühren
- `claimUBI()` — Tägliches UBI einfordern
- `addGuardian()` — Guardian-System für Proof of Alive
- `applyWealthCap()` — Vermögensobergrenze durchsetzen

**Sicherheitsfeatures:**
- Nullifier-Bindung verhindert Doppel-Registrierung
- `registerWithSig` nur von authorisierter Relayer-Adresse aufrufbar
- Optimistic Locking für Multi-Node-Schreibvorgänge
- Vollständiger Storage-Backup vor Contract-Upgrades

### EN
The AequitasV7 contract is the core of the protocol. Written in Solidity, deployed on Aequitas Chain, it contains all economic logic.

**Core Functions:**
- `register()` — Registration with ZKP proof (direct)
- `registerWithSig()` — Registration via relayer (gasless for user)
- `transfer()` — Token transfer with automatic demurrage and fees
- `claimUBI()` — Claim daily UBI
- `addGuardian()` — Guardian system for Proof of Alive
- `applyWealthCap()` — Enforce wealth ceiling

**Security Features:**
- Nullifier binding prevents double registration
- `registerWithSig` callable only from authorized relayer address
- Optimistic locking for multi-node writes
- Full storage backup before contract upgrades

---

## 7. Zero-Knowledge-Proofs & Privatsphäre / Privacy

### DE
Aequitas nutzt Groth16-Proofs auf der BN128-Kurve — eines der effizientesten ZKP-Systeme mit kleinen Proofs (~200 Bytes) und schneller On-Chain-Verifikation (~10ms).

**Nullifier-Bindung:** Der ZKP enthält einen eindeutigen Nullifier (`pubSignals[1]`), der kryptographisch an den biometrischen Hash gebunden ist. Derselbe Mensch kann denselben Nullifier nie zweimal verwenden — Sybil-Attacken sind mathematisch ausgeschlossen.

**Was gespeichert wird / What is stored:**
- ✅ `commitment` — kryptographischer Hash (nicht rückführbar auf Biometrie)
- ✅ `nullifier` — eindeutiger Einmal-Nachweis
- ✅ Wallet-Adresse
- ❌ Fingerabdruck-Daten — niemals
- ❌ Name, Adresse, ID — niemals
- ❌ IP-Adresse — nicht gespeichert

### EN
Aequitas uses Groth16 proofs on the BN128 curve — one of the most efficient ZKP systems with small proofs (~200 bytes) and fast on-chain verification (~10ms).

**Nullifier Binding:** The ZKP contains a unique nullifier (`pubSignals[1]`), cryptographically bound to the biometric hash. The same human can never use the same nullifier twice — Sybil attacks are mathematically impossible.

**What is stored:**
- ✅ `commitment` — cryptographic hash (not traceable to biometrics)
- ✅ `nullifier` — unique one-time proof
- ✅ Wallet address
- ❌ Fingerprint data — never
- ❌ Name, address, ID — never
- ❌ IP address — not stored

---

## 8. Gleichheitsindex / Equality Index

### DE
Aequitas ist das erste Währungssystem, das seinen eigenen Gleichheitsgrad live und transparent misst und veröffentlicht.

**Gini-Koeffizient:** Misst die Ungleichverteilung des AEQ-Vermögens. 0 = perfekte Gleichheit, 1 = totale Konzentration.

**Lorenz-Kurve:** Zeigt grafisch, wie viel Prozent des Reichtums die ärmsten X% der Menschen besitzen.

**Aequitas-Index:** Kombinierter Score aus Gini, Verteilung, Aktivität und Wachstum.

**Ziel / Target:** Gini < 0,35 (Skandinavien-Niveau)

| Währung / Currency | Gini |
|-------------------|------|
| **Aequitas AEQ** | **~0,08** |
| Skandinavien / Scandinavia | ~0,27 |
| Deutschland / Germany | ~0,31 |
| USA | ~0,41 |
| Brasilien / Brazil | ~0,53 |
| Bitcoin | ~0,85 |

### EN
Aequitas is the first monetary system that measures and publishes its own equality level live and transparently.

**Gini Coefficient:** Measures inequality of AEQ wealth distribution. 0 = perfect equality, 1 = total concentration.

**Lorenz Curve:** Graphically shows what percentage of wealth the poorest X% of humans own.

**Aequitas Index:** Combined score from Gini, distribution, activity, and growth.

**Target:** Gini < 0.35 (Scandinavia level)

---

## 9. Exchange & Liquiditätspool / Liquidity Pool

### DE
Aequitas enthält einen integrierten automatischen Market Maker (AMM) für den Handel zwischen AEQ und tUSD (einem simulierten Test-Dollar auf der Aequitas Chain).

**Mechanismus:** Das klassische `x·y=k`-Modell — der Pool hält automatisch einen Gleichgewichtspreis aufrecht.

**Gebührenverteilung / Fee Distribution:**
- 0,1% Swap-Gebühr wird automatisch aufgeteilt:
  - 40% → Validator-Pool (Netzwerkanreiz)
  - 30% → Liquidity Provider (LP-Rendite)
  - 20% → UBI-Pool (Grundeinkommen)
  - 10% → Treasury (Protokoll-Entwicklung)

**Liquiditäts-Shares:** LPs erhalten proportionale Shares und können jederzeit ihre Anteile plus akkumulierte Gebühren abheben.

### EN
Aequitas contains a built-in Automated Market Maker (AMM) for trading between AEQ and tUSD (a simulated test-dollar on Aequitas Chain).

**Mechanism:** The classic `x·y=k` model — the pool automatically maintains an equilibrium price.

**Fee Distribution:**
- 0.1% swap fee automatically split:
  - 40% → Validator Pool (network incentive)
  - 30% → Liquidity Providers (LP yield)
  - 20% → UBI Pool (basic income)
  - 10% → Treasury (protocol development)

**Liquidity Shares:** LPs receive proportional shares and can withdraw their stakes plus accumulated fees at any time.

---

## 10. Sicherheit / Security

### Auditierte Sicherheitsmaßnahmen / Audited Security Measures

| Bedrohung / Threat | Schutz / Protection |
|-------------------|---------------------|
| Doppel-Registrierung / Double registration | Nullifier-Bindung on-chain / Nullifier binding on-chain |
| Replay-Attacke / Replay attack | Nonce-System mit CAS / Nonce system with compare-and-swap |
| Sybil-Attacke / Sybil attack | Biometrie + ZKP + Hardware Secure Element |
| Pool-Drain | Wealth Cap + Demurrage + optimistic locking |
| Contract-Upgrade-Risiko | Vollständiger Storage-Backup vor Wipe / Full storage backup before wipe |
| Multi-Node-Konflikte / Multi-node conflicts | PostgreSQL optimistic locking + SELECT FOR UPDATE |
| Öffentliche Contract-Deployments / Public deployments | Deployment auf Relayer-Adresse beschränkt / Restricted to relayer |
| Private Keys in Logs | Ausgabe nur auf stderr (nicht in Log-Aggregatoren) / stderr only |
| XSS | HTML-Escaping aller User-Eingaben / HTML escaping of all user inputs |
| DNS-Rebinding | Peer-URL-Validierung + Öffentliche-IP-Prüfung / Public IP check |

### Dezentralisierung / Decentralization

Aequitas befindet sich in Phase 0 mit zwei betriebenen Nodes. Das Protokoll ist für beliebig viele Nodes ausgelegt — jeder Node-Betreiber kann mit `PEER_NODES` beitreten.

Aequitas is in Phase 0 with two operated nodes. The protocol is designed for any number of nodes — any operator can join with `PEER_NODES`.

---

## 11. Roadmap

| Phase | Status | DE | EN |
|-------|--------|----|----|
| 0 | ✅ | Smart Contracts · ZKP · Android App · Proof Server | Smart Contracts · ZKP · Android App · Proof Server |
| 0+ | ✅ | Aequitas Layer 1 · BlockDAG · P2P · Explorer | Aequitas Layer 1 · BlockDAG · P2P · Explorer |
| V7 | ✅ | EVM · Dual-Ledger · Exchange/AMM · UBI · Demurrage · Wealth Cap · Lorenz/Gini | EVM · Dual-Ledger · Exchange/AMM · UBI · Demurrage · Wealth Cap · Lorenz/Gini |
| 1 | 🔄 | APK-Veröffentlichung · Community-Aufbau · Grant-Anträge | APK Release · Community Growth · Grant Applications |
| 2 | ⬜ | iOS App · Proof of Alive · Guardian-System live | iOS App · Proof of Alive · Guardian System live |
| 3 | ⬜ | Cross-Chain Bridges · Externe DEX-Integration | Cross-Chain Bridges · External DEX Integration |
| 4 | ⬜ | Vollständige Dezentralisierung · Community Governance | Full Decentralization · Community Governance |

---

## 12. Fazit / Conclusion

### DE
Aequitas ist kein weiteres Experiment in Kryptospekulation. Es ist ein ernsthafter Versuch, Geld neu zu denken — von Grund auf, für alle Menschen, fair.

Die mathematische Garantie ist simpel und radikal zugleich: Solange Menschen existieren, existiert AEQ. Kein Zentralstaat, keine Bank, kein Algorithmus kann das Grundeinkommen entziehen oder die Gleichheit untergraben — es ist Code.

Der Gini-Koeffizient von Aequitas liegt heute bei ~0,08. Bitcoin liegt bei ~0,85. Der Unterschied ist nicht zufällig — er ist das Ergebnis des Designs.

### EN
Aequitas is not another experiment in crypto speculation. It is a serious attempt to rethink money — from first principles, for all people, fairly.

The mathematical guarantee is simple and radical at once: as long as humans exist, AEQ exists. No central state, no bank, no algorithm can remove the basic income or undermine the equality — it is code.

Aequitas's Gini coefficient today is ~0.08. Bitcoin's is ~0.85. The difference is not coincidence — it is the result of design.

---

*Aequitas · Chain ID 1926 · aequitas.digital*
*Version 2.0 · Juni / June 2026*
*Lizenz / License: MIT · Open Source: github.com/hanoi96international-gif/Aequitas*
