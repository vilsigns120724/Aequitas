// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Aequitas {

// ─── STATE ───────────────────────────────────────
mapping(bytes32 => bool) private registeredHashes;
mapping(address => uint256) public balances;
uint256 public totalSupply;
uint256 public totalHumans;
uint256 public constant INITIAL_GRANT = 1000 * 10**18;

// ─── EVENTS ──────────────────────────────────────
event HumanRegistered(address wallet, uint256 newCap);
event Transfer(address from, address to, uint256 amount);
event Inflation(uint256 amountPerWallet);
event Deflation(uint256 amountPerWallet);

// ─── REGISTRIERUNG ───────────────────────────────
function registerHuman(bytes32 biometricHash) external {
require(!registeredHashes[biometricHash], "Bereits registriert");
require(balances[msg.sender] == 0, "Wallet bereits aktiv");

registeredHashes[biometricHash] = true;
totalHumans += 1;
balances[msg.sender] = INITIAL_GRANT;
totalSupply += INITIAL_GRANT;

emit HumanRegistered(msg.sender, maxCap());
}

// ─── DYNAMISCHER CAP ─────────────────────────────
function maxCap() public view returns (uint256) {
return totalHumans * INITIAL_GRANT;
}

// ─── TRANSAKTION ─────────────────────────────────
function transfer(address to, uint256 amount) external {
require(balances[msg.sender] >= amount, "Nicht genug Guthaben");
balances[msg.sender] -= amount;
balances[to] += amount;
emit Transfer(msg.sender, to, amount);
}

// ─── GELDMECHANISMUS ─────────────────────────────
function runCycle(bool isNegative, address[] calldata users) external {
uint256 rate = totalSupply / 200; // 0.5%
uint256 perUser = rate / users.length;

if (isNegative) {
for (uint i = 0; i < users.length; i++) {
balances[users[i]] += perUser;
}
totalSupply += rate;
emit Inflation(perUser);
} else {
for (uint i = 0; i < users.length; i++) {
if (balances[users[i]] >= perUser) {
balances[users[i]] -= perUser;
}
}
totalSupply -= rate;
emit Deflation(perUser);
}
}

// ─── STATUS ──────────────────────────────────────
function getStatus() external view returns (
uint256 humans,
uint256 supply,
uint256 cap
) {
return (totalHumans, totalSupply, maxCap());
}
}