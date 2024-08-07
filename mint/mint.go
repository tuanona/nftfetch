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
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func MintNFT() error {
	// Initialize connection to TON network
	client := liteclient.NewConnectionPool()
	err := client.AddConnection(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		return fmt.Errorf("failed to connect to TON network: %w", err)
	}

	// Initialize API client
	api := ton.NewAPIClient(client)

	// Load wallet seed
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	seedPath := filepath.Join(homeDir, ".nftfetch", "wallet", "mnemonics.txt")
	seedContent, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read wallet seed: %w", err)
	}
	walletSeed := strings.Split(string(seedContent), " ")

	// Load wallet
	w, err := wallet.FromSeed(api, walletSeed, wallet.V4R2)
	if err != nil {
		return fmt.Errorf("failed to load wallet: %w", err)
	}

	// Parse collection address
	collectionAddr := address.MustParseAddr("EQCSrRIKVEBaRd8aQfsOaNq3C4FVZGY5Oka55A5oFMVEs0lY")
	collection := nft.NewCollectionClient(api, collectionAddr)

	// Get collection data
	collectionData, err := collection.GetCollectionData(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get collection data: %w", err)
	}

	// Read metahash
	metahashPath := filepath.Join(homeDir, ".nftfetch", "metadata", "metahash.txt")
	metahash, err := os.ReadFile(metahashPath)
	if err != nil {
		return fmt.Errorf("failed to read metahash: %w", err)
	}

	// Prepare mint data
	mintData, err := collection.BuildMintPayload(
		collectionData.NextItemIndex,
		w.WalletAddress(),
		tlb.MustFromTON("0.05"),
		&nft.ContentOffchain{URI: string(metahash)},
	)
	if err != nil {
		return fmt.Errorf("failed to build mint payload: %w", err)
	}

	// Send mint transaction
	msg := wallet.SimpleMessage(collectionAddr, tlb.MustFromTON("0.05"), mintData)
	err = w.Send(context.Background(), msg, true)
	if err != nil {
		return fmt.Errorf("failed to send mint transaction: %w", err)
	}

	fmt.Println("NFT minted successfully!")
	return nil
}
