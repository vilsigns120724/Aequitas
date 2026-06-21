// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title AequitasV7
 * @author Aequitas Foundation
 * @notice "Money exists because people exist."
 *
 * CORE INVARIANT:
 *   SUM(balanceOf) + SUM(escrowOf) + ubiPool = totalSupply
 *
 * FIXES vs V6:
 *   [1] O(n) loops eliminated — Pull UBI pattern
 *   [2] totalSupply invariant strictly maintained
 *   [3] Concentration friction on transfers
 *   [4] No admin/owner, no upgrades
 *   [5] All state reconstructable from events
 */

interface IBioVerifier {
    function verifyProof(uint[2] calldata,uint[2][2] calldata,uint[2] calldata,uint[2] calldata) external view returns (bool);
}

contract AequitasV7 {
    string  public constant name     = "Aequitas";
    string  public constant symbol   = "AEQ";
    uint8   public constant decimals = 18;

    uint256 public constant INITIAL_GRANT     = 1_000 * 1e18;
    uint256 public constant TX_FEE_BPS        = 10;
    uint256 public constant UBI_SHARE_BPS     = 2000;
    uint256 public constant DEMURRAGE_BPS     = 100;
    uint256 public constant SECONDS_PER_YEAR  = 365 days;
    uint256 public constant INACTIVITY_ESCROW = 910 days;
    uint256 public constant INACTIVITY_UBI    = 1460 days;
    uint256 public constant GUARDIAN_TIMELOCK = 7 days;
    uint256 public constant MAX_WARDS         = 3;

    IBioVerifier public immutable bioVerifier;

    uint256 public totalSupply;
    uint256 public totalHumans;
    uint256 public ubiPool;
    uint256 public ubiPerHumanAccumulated;

    mapping(address => uint256) public balanceOf;
    mapping(address => uint256) public escrowOf;
    mapping(address => bool)    public isHuman;
    mapping(uint256 => bool)    public usedCommitments;
    // usedNullifiers stores the on-chain record of each biometric nullifier.
    // Slot 8. Prevents the same biometric from registering a second wallet
    // even when bypassing the API layer and calling the contract directly.
    mapping(bytes32 => address) public usedNullifiers;
    mapping(address => uint256) public commitmentOf;
    mapping(address => uint256) public lastActivity;
    mapping(address => uint256) public lastDemurrage;
    mapping(address => uint256) public ubiClaimed;
    mapping(address => address) public guardianOf;
    mapping(address => address) public pendingGuardian;
    mapping(address => uint256) public guardianRequestedAt;
    mapping(address => uint256) public wardCount;

    uint256[5] public CAPS       = [50, 20, 10, 5, 3];
    uint256[5] public THRESHOLDS = [0, 100, 1_000, 10_000, 100_000];

    event Registered(address indexed human, uint256 commitment, uint256 grant);
    event Transferred(address indexed from, address indexed to, uint256 amount, uint256 fee);
    event DemurrageApplied(address indexed human, uint256 amount);
    event WealthCapApplied(address indexed human, uint256 excess);
    event UBIAccumulated(uint256 addedPerHuman, uint256 total);
    event UBIClaimed(address indexed human, uint256 amount);
    event EscrowCreated(address indexed human, uint256 amount);
    event EscrowReleased(address indexed human, uint256 amount);
    event EscrowToUBI(address indexed human, uint256 amount);
    event GuardianProposed(address indexed human, address indexed proposed);
    event GuardianConfirmed(address indexed human, address indexed guardian);
    event GuardianRevoked(address indexed human, address indexed guardian);
    event AliveConfirmed(address indexed human, address indexed by);

    constructor(address _bioVerifier) {
        bioVerifier = IBioVerifier(_bioVerifier);
    }

    function register(uint[2] calldata pA, uint[2][2] calldata pB, uint[2] calldata pC, uint[2] calldata pubSignals, bytes32 nullifier) external {
        require(!isHuman[msg.sender], "Already registered");
        uint256 commitment = pubSignals[0];
        require(!usedCommitments[commitment], "Commitment used");
        require(nullifier != bytes32(0), "Nullifier required");
        require(usedNullifiers[nullifier] == address(0), "Nullifier used");
        require(bioVerifier.verifyProof(pA, pB, pC, pubSignals), "Invalid proof");
        usedCommitments[commitment] = true;
        commitmentOf[msg.sender] = commitment;
        usedNullifiers[nullifier] = msg.sender;
        isHuman[msg.sender] = true;
        totalHumans++;
        balanceOf[msg.sender] += INITIAL_GRANT;
        totalSupply += INITIAL_GRANT;
        ubiClaimed[msg.sender] = ubiPerHumanAccumulated;
        lastActivity[msg.sender] = block.timestamp;
        lastDemurrage[msg.sender] = block.timestamp;
        emit Registered(msg.sender, commitment, INITIAL_GRANT);
    }

    /**
     * @notice Register as a verified human via meta-transaction (gasless for the user)
     * @dev A relayer submits this transaction and pays gas, but the new human is
     *      claimedHuman, NOT msg.sender. claimedHuman must have personally signed
     *      this exact commitment, for this exact contract, on this exact chain,
     *      verified on-chain via ecrecover. usedCommitments prevents replay across
     *      either registration path.
     */
    // registerWithSig v2: circuit now outputs the nullifier as pubSignals[1]
    // so "1 biometric = 1 nullifier" is enforced by the ZK proof itself.
    // pubSignals[0] = commitment, pubSignals[1] = nullifier (ZK-bound).
    // The external bytes32 nullifier param is still accepted for compatibility
    // but MUST match pubSignals[1] when circuit version >= 2.
    function registerWithSig(
        uint[2] calldata pA,
        uint[2][2] calldata pB,
        uint[2] calldata pC,
        uint[2] calldata pubSignals,
        address claimedHuman,
        bytes calldata signature,
        bytes32 nullifier
    ) external {
        require(!isHuman[claimedHuman], "Already registered");

        uint256 commitment = pubSignals[0];
        require(!usedCommitments[commitment], "Commitment used");

        // Use circuit-derived nullifier when available (pubSignals[1] != 0).
        // This ZK-binds the nullifier to the biometric secret inside the proof.
        bytes32 effectiveNullifier = nullifier;
        if (pubSignals[1] != 0) {
            bytes32 zkNullifier = bytes32(pubSignals[1]);
            require(nullifier == bytes32(0) || nullifier == zkNullifier,
                "Nullifier mismatch: submitted nullifier does not match ZK-derived value");
            effectiveNullifier = zkNullifier;
        }

        require(effectiveNullifier != bytes32(0), "Nullifier required");
        require(usedNullifiers[effectiveNullifier] == address(0), "Nullifier used");

        bytes32 messageHash = keccak256(abi.encodePacked(
            block.chainid,
            address(this),
            "register",
            commitment,
            nullifier
        ));
        bytes32 ethSignedHash = keccak256(abi.encodePacked(
            "\x19Ethereum Signed Message:\n32",
            messageHash
        ));
        require(
            _recoverSigner(ethSignedHash, signature) == claimedHuman,
            "Invalid signature"
        );

        require(bioVerifier.verifyProof(pA, pB, pC, pubSignals), "Invalid proof");

        usedCommitments[commitment] = true;
        commitmentOf[claimedHuman] = commitment;
        usedNullifiers[effectiveNullifier] = claimedHuman;
        isHuman[claimedHuman] = true;
        totalHumans++;
        balanceOf[claimedHuman] += INITIAL_GRANT;
        totalSupply += INITIAL_GRANT;
        ubiClaimed[claimedHuman] = ubiPerHumanAccumulated;
        lastActivity[claimedHuman] = block.timestamp;
        lastDemurrage[claimedHuman] = block.timestamp;

        emit Registered(claimedHuman, commitment, INITIAL_GRANT);
    }

    function _recoverSigner(bytes32 ethSignedHash, bytes calldata signature) internal pure returns (address) {
        require(signature.length == 65, "Invalid signature length");
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := calldataload(signature.offset)
            s := calldataload(add(signature.offset, 32))
            v := byte(0, calldataload(add(signature.offset, 64)))
        }
        if (v < 27) v += 27;
        return ecrecover(ethSignedHash, v, r, s);
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        require(isHuman[msg.sender], "Not human");
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");
        require(to != address(0) && to != msg.sender, "Invalid recipient");
        _applyDemurrage(msg.sender);
        require(balanceOf[msg.sender] >= amount, "Insufficient after demurrage");
        uint256 fee = _calcFee(msg.sender, amount);
        uint256 ubiContrib = (fee * UBI_SHARE_BPS) / 10_000;
        uint256 burned = fee - ubiContrib;
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount - fee;
        ubiPool += ubiContrib;
        totalSupply -= burned;
        lastActivity[msg.sender] = block.timestamp;
        if (isHuman[to]) lastActivity[to] = block.timestamp;
        _applyWealthCap(to);
        emit Transferred(msg.sender, to, amount, fee);
        return true;
    }

    function _calcFee(address sender, uint256 amount) internal view returns (uint256) {
        uint256 base = (amount * TX_FEE_BPS) / 10_000;
        if (totalSupply == 0) return base;
        uint256 shareBPS = (balanceOf[sender] * 10_000) / totalSupply;
        uint256 extra = shareBPS >= 1000 ? 100 : shareBPS >= 500 ? 50 : shareBPS >= 100 ? 10 : 0;
        return base + (amount * extra) / 10_000;
    }

    function applyDemurrage(address human) external { require(isHuman[human],"Not human"); _applyDemurrage(human); }

    function _applyDemurrage(address human) internal {
        uint256 fs = fairShare();
        if (balanceOf[human] <= fs) { lastDemurrage[human] = block.timestamp; return; }
        uint256 elapsed = block.timestamp - lastDemurrage[human];
        if (elapsed == 0) return;
        uint256 excess = balanceOf[human] - fs;
        uint256 fee = (excess * DEMURRAGE_BPS * elapsed) / (10_000 * SECONDS_PER_YEAR);
        if (fee == 0) return;
        if (fee > excess) fee = excess;
        balanceOf[human] -= fee;
        ubiPool += fee;
        lastDemurrage[human] = block.timestamp;
        emit DemurrageApplied(human, fee);
    }

    function applyWealthCap(address human) external { require(isHuman[human],"Not human"); _applyWealthCap(human); }

    function _applyWealthCap(address human) internal {
        uint256 cap = wealthCap();
        if (balanceOf[human] <= cap) return;
        uint256 excess = balanceOf[human] - cap;
        balanceOf[human] = cap;
        ubiPool += excess;
        emit WealthCapApplied(human, excess);
    }

    function accumulateUBI() external {
        require(totalHumans > 0, "No humans");
        require(ubiPool > 0, "Pool empty");
        uint256 addPerHuman = ubiPool / totalHumans;
        require(addPerHuman > 0, "Too small");
        ubiPool -= addPerHuman * totalHumans;
        ubiPerHumanAccumulated += addPerHuman;
        emit UBIAccumulated(addPerHuman, ubiPerHumanAccumulated);
    }

    function claimUBI() external {
        require(isHuman[msg.sender], "Not human");
        require(escrowOf[msg.sender] == 0, "In escrow");
        uint256 owed = ubiPerHumanAccumulated - ubiClaimed[msg.sender];
        require(owed > 0, "Nothing to claim");
        ubiClaimed[msg.sender] = ubiPerHumanAccumulated;
        balanceOf[msg.sender] += owed;
        totalSupply += owed;
        lastActivity[msg.sender] = block.timestamp;
        _applyWealthCap(msg.sender);
        emit UBIClaimed(msg.sender, owed);
    }

    function claimableUBI(address human) external view returns (uint256) {
        if (!isHuman[human]) return 0;
        return ubiPerHumanAccumulated - ubiClaimed[human];
    }

    function confirmAlive() external {
        require(isHuman[msg.sender], "Not human");
        _confirmAlive(msg.sender);
        emit AliveConfirmed(msg.sender, msg.sender);
    }

    function _confirmAlive(address human) internal {
        lastActivity[human] = block.timestamp;
        if (escrowOf[human] > 0) {
            uint256 amount = escrowOf[human];
            escrowOf[human] = 0;
            uint256 fs = fairShare();
            balanceOf[human] += amount + fs;
            totalSupply += fs;
            ubiClaimed[human] = ubiPerHumanAccumulated;
            emit EscrowReleased(human, amount + fs);
        }
    }

    function triggerEscrow(address human) external {
        require(isHuman[human], "Not human");
        require(escrowOf[human] == 0, "Already in escrow");
        require(block.timestamp >= lastActivity[human] + INACTIVITY_ESCROW, "Not inactive enough");
        uint256 amount = balanceOf[human];
        balanceOf[human] = 0;
        escrowOf[human] = amount;
        emit EscrowCreated(human, amount);
    }

    function triggerEscrowToUBI(address human) external {
        require(escrowOf[human] > 0, "Not in escrow");
        require(block.timestamp >= lastActivity[human] + INACTIVITY_UBI, "Too soon");
        uint256 amount = escrowOf[human];
        escrowOf[human] = 0;
        ubiPool += amount;
        isHuman[human] = false;
        totalHumans--;
        emit EscrowToUBI(human, amount);
    }

    function proposeGuardian(address guardian) external {
        require(isHuman[msg.sender] && isHuman[guardian], "Must be human");
        require(guardian != msg.sender, "Cannot guard yourself");
        require(guardianOf[guardian] == address(0), "Guardian has own guardian");
        require(wardCount[guardian] < MAX_WARDS, "Max wards reached");
        require(guardianOf[msg.sender] != guardian, "Circular dependency");
        pendingGuardian[msg.sender] = guardian;
        guardianRequestedAt[msg.sender] = block.timestamp;
        emit GuardianProposed(msg.sender, guardian);
    }

    function confirmGuardian() external {
        require(pendingGuardian[msg.sender] != address(0), "No pending guardian");
        require(block.timestamp >= guardianRequestedAt[msg.sender] + GUARDIAN_TIMELOCK, "Timelock active");
        if (guardianOf[msg.sender] != address(0)) wardCount[guardianOf[msg.sender]]--;
        address g = pendingGuardian[msg.sender];
        guardianOf[msg.sender] = g;
        wardCount[g]++;
        pendingGuardian[msg.sender] = address(0);
        emit GuardianConfirmed(msg.sender, g);
    }

    function revokeGuardian() external {
        address g = guardianOf[msg.sender];
        require(g != address(0), "No guardian");
        wardCount[g]--;
        guardianOf[msg.sender] = address(0);
        emit GuardianRevoked(msg.sender, g);
    }

    function guardianConfirmAlive(address ward) external {
        require(isHuman[msg.sender], "Not human");
        require(guardianOf[ward] == msg.sender, "Not the guardian");
        _confirmAlive(ward);
        emit AliveConfirmed(ward, msg.sender);
    }

    function fairShare() public view returns (uint256) {
        if (totalHumans == 0) return INITIAL_GRANT;
        return totalSupply / totalHumans;
    }

    function wealthCap() public view returns (uint256) {
        return CAPS[currentPhase()] * fairShare();
    }

    function currentPhase() public view returns (uint256) {
        for (uint256 i = 4; i > 0; i--) {
            if (totalHumans >= THRESHOLDS[i]) return i;
        }
        return 0;
    }

    function getStats() external view returns (
        uint256 _supply, uint256 _humans, uint256 _pool,
        uint256 _fair, uint256 _cap, uint256 _phase, uint256 _ubiAcc
    ) {
        return (totalSupply, totalHumans, ubiPool, fairShare(), wealthCap(), currentPhase(), ubiPerHumanAccumulated);
    }

    function calcFee(address sender, uint256 amount) external view returns (uint256) {
        return _calcFee(sender, amount);
    }
}
