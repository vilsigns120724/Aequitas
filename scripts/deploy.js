import { network } from "hardhat";
import { parseEther } from "viem";

const { viem } = await import("hardhat");

const aequitas = await viem.deployContract("Aequitas");
console.log("Aequitas deployed to:", aequitas.address);

const [humans, supply, cap] = await aequitas.read.getStatus();
console.log("Humans: ", humans.toString());
console.log("Total Supply:", supply.toString());
console.log("Max Cap: ", cap.toString());
console.log("\nAequitas laeuft auf der Chain.");