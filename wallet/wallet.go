package wallet

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	qrcodeTerminal "nftfetch/wallet/qrcode" // Import your updated qrcodeTerminal package

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	walletFile = "mnemonics.txt"
	configUrl  = "https://ton-blockchain.github.io/testnet-global.config.json"
)

// CreateOrLoadWallet creates a new wallet if it doesn't exist, or loads the existing wallet
func CreateOrLoadWallet() {
	client := liteclient.NewConnectionPool()
	ctx := client.StickyContext(context.Background())

	err := client.AddConnectionsFromConfigUrl(ctx, configUrl)
	if err != nil {
		log.Fatalf("Failed to connect to TON: %v", err)
	}

	api := ton.NewAPIClient(client)

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	// Define wallet directory in the home directory
	walletDir := filepath.Join(homeDir, ".nftfetch", "wallet")

	// Path to the mnemonics file
	mnemonicsPath := filepath.Join(walletDir, walletFile)

	var walletSeed []string

	// Check if wallet already exists
	if _, err := os.Stat(mnemonicsPath); os.IsNotExist(err) {
		// Wallet does not exist, create a new one
		fmt.Println("Wallet does not exist. Creating a new wallet...")

		walletSeed = wallet.NewSeed()
		mnemonics := strings.Join(walletSeed, " ")

		// Create wallet directory if it does not exist
		if err := os.MkdirAll(walletDir, 0755); err != nil {
			log.Fatalf("Failed to create wallet directory: %v", err)
		}

		// Save mnemonics to file
		err = os.WriteFile(mnemonicsPath, []byte(mnemonics), 0644)
		if err != nil {
			log.Fatalf("Failed to save mnemonics: %v", err)
		}

		fmt.Printf("Wallet created and mnemonics saved to: %s\n", mnemonicsPath)
	} else {
		// Wallet exists, load the mnemonics
		fmt.Println("Wallet already exists. Loading wallet...")

		mnemonics, err := os.ReadFile(mnemonicsPath)
		if err != nil {
			log.Fatalf("Failed to read mnemonics: %v", err)
		}

		walletSeed = strings.Split(string(mnemonics), " ")
	}

	// Load the wallet from the mnemonics
	w, err := wallet.FromSeed(api, walletSeed, wallet.V4R2)
	if err != nil {
		log.Fatalf("Failed to create wallet from seed: %v", err)
	}

	// Get the current masterchain info
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get current masterchain info: %v", err)
	}

	// Get the wallet balance
	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalf("Failed to get wallet balance: %v", err)
	}

	fmt.Printf("Wallet balance: %v\n", balance)

	// Generate QR code for wallet address
	address := w.Address().String() // Ensure address is in string format

	// Create an instance of QRCodeTerminal
	qrTerminal := qrcodeTerminal.NewQRCodeTerminal()

	// Define the path for the QR code image
	walletImagePath := filepath.Join(walletDir, "wallet.png")

	// Print the QR code for the wallet address to terminal
	fmt.Println("Printing QR code to terminal...")
	qrTerminal.Print(address)
	fmt.Println("QR code for wallet address printed in terminal.")
	fmt.Println("Address: ")
	fmt.Println(address)
	// Generate and save QR code to PNG file
	fmt.Println("Generating QR code for wallet address...")
	qrTerminal.Save(address, walletImagePath)
	fmt.Printf("QR CODE in PNG format for wallet address saved to %s\n", walletImagePath)

}
