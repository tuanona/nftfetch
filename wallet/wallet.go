package wallet

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	qrcodeTerminal "nftfetch/qrcode" // Import your qrcodeTerminal package

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	walletDir  = ".nftfetch/wallet"
	walletFile = "mnemonics.txt"
	configUrl  = "https://ton-blockchain.github.io/testnet-global.config.json"
)

// CreateWallet creates a new wallet, saves mnemonics to a file, and prints a QR code of the wallet address
func CreateWallet() {
	client := liteclient.NewConnectionPool()
	ctx := client.StickyContext(context.Background())

	err := client.AddConnectionsFromConfigUrl(ctx, configUrl)
	if err != nil {
		log.Fatalf("Failed to connect to TON: %v", err)
	}

	api := ton.NewAPIClient(client)

	// Check if wallet already exists
	if _, err := os.Stat(filepath.Join(walletDir, walletFile)); !os.IsNotExist(err) {
		fmt.Println("Wallet already exists. Skipping creation.")
		return
	}

	// Generate a new wallet
	walletSeed := wallet.NewSeed()
	w, err := wallet.FromSeed(api, walletSeed, wallet.V4R2)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
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

	// Create wallet directory if not exists
	if err := os.MkdirAll(walletDir, 0755); err != nil {
		log.Fatalf("Failed to create wallet directory: %v", err)
	}

	// Save mnemonics to file
	mnemonicsPath := filepath.Join(walletDir, walletFile)
	mnemonics := strings.Join(walletSeed, " ")
	err = os.WriteFile(mnemonicsPath, []byte(mnemonics), 0644)
	if err != nil {
		log.Fatalf("Failed to save mnemonics: %v", err)
	}

	fmt.Printf("Wallet created and mnemonics saved to: %s\n", mnemonicsPath)

	// Generate QR code for wallet address
	address := w.Address().String() // Ensure address is in string format

	// Create an instance of QRCodeTerminal
	qrTerminal := qrcodeTerminal.NewQRCodeTerminal()

	// Print the QR code for the wallet address
	fmt.Println("Generating QR code for wallet address...")
	qrTerminal.Print(address)

	fmt.Println("QR code for wallet address printed in terminal.")
}
