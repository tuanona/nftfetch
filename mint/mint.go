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

const MintOpcode uint64 = 4215445508

type Mint struct {
	Metadata *cell.Cell
}

// StoreMint stores the Mint message in a cell to be sent to the smart contract
func StoreMint(src Mint) *cell.Cell {
	b := cell.BeginCell()
	b.MustStoreUInt(MintOpcode, 32)
	b.MustStoreRef(src.Metadata)
	return b.EndCell()
}

func MintNft() error {
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
	collectionAddr := address.MustParseAddr("kQC14v2x_3Pc2hE1Z_z58Ks6FUSksrahn8qcerfTBKnXaJjC")

	// Creating the Cell with metadata
	metadataCell := cell.BeginCell().MustStoreStringSnake(metadataHash).EndCell()

	// Prepare the Mint message and the body for the transaction
	mintMsg := Mint{Metadata: metadataCell}
	body := StoreMint(mintMsg)

	// Send the transaction with the mint message to the contract
	err = w.Send(context.Background(), &wallet.Message{
		Mode: 1,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     collectionAddr,
			Amount:      tlb.MustFromTON("0.05"),
			Body:        body,
		},
	}, true)

	if err != nil {
		return fmt.Errorf("failed to mint NFT: %v", err)
	}

	fmt.Println("NFT minted successfully!")
	return nil
}

// Get wallet using the mnemonic seed from a file
func getWallet(api *ton.APIClient) (*wallet.Wallet, error) {
	words, err := readSeedFromFile()
	if err != nil {
		return nil, err
	}
	return wallet.FromSeed(api, words, wallet.V4R2)
}

// Read seed from mnemonics.txt file
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

// Read metadata hash from file
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
