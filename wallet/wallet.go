package wallet

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"nftfetch/wallet/qrcode"

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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	walletDir := filepath.Join(homeDir, ".nftfetch", "wallet")
	mnemonicsPath := filepath.Join(walletDir, walletFile)

	walletSeed := loadOrCreateWallet(walletDir, mnemonicsPath)

	w, err := wallet.FromSeed(api, walletSeed, wallet.V4R2)
	if err != nil {
		log.Fatalf("Failed to create wallet from seed: %v", err)
	}

	displayWalletInfo(ctx, api, w)
	generateQRCode(w, walletDir)
}

func loadOrCreateWallet(walletDir, mnemonicsPath string) []string {
	if _, err := os.Stat(mnemonicsPath); os.IsNotExist(err) {
		return createNewWallet(walletDir, mnemonicsPath)
	}
	return loadExistingWallet(mnemonicsPath)
}

func createNewWallet(walletDir, mnemonicsPath string) []string {
	fmt.Println("Wallet does not exist. Creating a new wallet...")

	walletSeed := wallet.NewSeed()
	mnemonics := strings.Join(walletSeed, " ")

	if err := os.MkdirAll(walletDir, 0755); err != nil {
		log.Fatalf("Failed to create wallet directory: %v", err)
	}

	if err := os.WriteFile(mnemonicsPath, []byte(mnemonics), 0644); err != nil {
		log.Fatalf("Failed to save mnemonics: %v", err)
	}

	fmt.Printf("Wallet created and mnemonics saved to: %s\n", mnemonicsPath)
	return walletSeed
}

func loadExistingWallet(mnemonicsPath string) []string {
	fmt.Println("Wallet already exists. Loading wallet...")

	mnemonics, err := os.ReadFile(mnemonicsPath)
	if err != nil {
		log.Fatalf("Failed to read mnemonics: %v", err)
	}

	return strings.Split(string(mnemonics), " ")
}

func displayWalletInfo(ctx context.Context, api *ton.APIClient, w *wallet.Wallet) {
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get current masterchain info: %v", err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalf("Failed to get wallet balance: %v", err)
	}

	fmt.Printf("Wallet balance: %v\n", balance)
	fmt.Println("Wallet address:", w.Address().String())
}

func generateQRCode(w *wallet.Wallet, walletDir string) {
	address := w.Address().String()
	walletImagePath := filepath.Join(walletDir, "address.png")

	qr := qrcode.NewQRCode()
	qr.GenerateAndSave(address, walletImagePath)

	fmt.Println("QR code for wallet address:")
	qr.Print(address)
	fmt.Printf("QR code image saved to: %s\n", walletImagePath)
}
