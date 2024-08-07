package main

import (
	"fmt"
	"nftfetch/compare"
	"nftfetch/generate"
	"nftfetch/mint"
	"nftfetch/wallet"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "generate":
		if err := generate.Generate(); err != nil {
			fmt.Printf("Error generating NFT: %v\n", err)
		}
	case "compare":
		if err := compare.Compare(); err != nil {
			fmt.Printf("Error comparing NFTs: %v\n", err)
		}
	case "wallet":
		wallet.CreateOrLoadWallet()
	case "mint":
		if err := mint.MintNFT(); err != nil {
			fmt.Printf("Error minting NFT: %v\n", err)
		}
	case "help":
		printHelp()
	default:
		fmt.Println("Invalid command. Use 'nftfetch help' for usage information.")
	}
}

func printHelp() {
	fmt.Println("Usage: nftfetch <command>")
	fmt.Println("Commands:")
	fmt.Println("  generate     Generate a new NFT based on device metadata")
	fmt.Println("  compare      Compare the validity of the owner")
	fmt.Println("  wallet       Create or load a wallet")
	fmt.Println("  mint         Mint a new NFT")
	fmt.Println("  help         Show this help message")
}
