package keeper

// Fix 8: syncHumans and syncFromSepolia have been removed.
// These functions bypassed biometric verification by creating fake "sepolia_human_N"
// addresses without ZK proof validation. Human registration must go through the
// proper RegisterHuman flow with nullifier + ZK proof verification.
// The StartSync goroutine that called syncFromSepolia on a 60-second timer
// has also been removed.
