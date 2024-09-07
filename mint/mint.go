package mint

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func Mint() error {
	client := liteclient.NewConnectionPool()
	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return fmt.Errorf("failed to add connections: %v", err)
	}

	api := ton.NewAPIClient(client)
	w, err := getWallet(api)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	fmt.Printf("Wallet address: %s\n", w.Address())

	metadataHash, err := readMetadataHash()
	if err != nil {
		return fmt.Errorf("failed to read metadata hash: %v", err)
	}

	// Assuming the NFT collection contract address is known
	collectionAddr := address.MustParseAddr("kQDEyWHyHbw7RuUiVOHfTgQhAYwjJjmtDkbZbPMEnL9V0nA_")

	// Creating the Cell with metadata
	metadataCell := cell.BeginCell().
		MustStoreStringSnake(metadataHash).EndCell() // Storing metadata in a Cell

	metadata := cell.BeginCell().
		MustStoreUInt(0x1, 32). // op code for mint
		MustStoreRef(metadataCell).
		EndCell()

	err = w.Send(context.Background(), &wallet.Message{
		Mode: 1,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     collectionAddr,
			Amount:      tlb.MustFromTON("0.05"),
			Body:        metadata,
		},
	}, true)

	if err != nil {
		return fmt.Errorf("failed to mint NFT: %v", err)
	}

	fmt.Println("NFT minted successfully!")
	return nil
}

func getWallet(api *ton.APIClient) (*wallet.Wallet, error) {
	words, err := readSeedFromFile()
	if err != nil {
		return nil, err
	}
	return wallet.FromSeed(api, words, wallet.V4R2)
}

func readSeedFromFile() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	seedFile := filepath.Join(home, ".nftfetch", "wallet", "mnemonics.txt")
	content, err := os.ReadFile(seedFile)
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(content)), nil
}

func readMetadataHash() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	metaFile := filepath.Join(home, ".nftfetch", "metadata", "metahash.txt")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
