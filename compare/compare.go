package compare

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

// GetMetadata sends a message to the smart contract to retrieve metadata and print it
func Compare() error {
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

	// Address of the parent contract
	parentContractAddr := address.MustParseAddr("kQC14v2x_3Pc2hE1Z_z58Ks6FUSksrahn8qcerfTBKnXaJjC") // Replace with actual address

	// Metadata that was originally sent during minting
	metadata, err := readMetadataHash()
	if err != nil {
		return fmt.Errorf("failed to read metadata: %v", err)
	}

	// Create a cell with the metadata hash
	metadataCell := cell.BeginCell().MustStoreStringSnake(metadata).EndCell()

	// Create a message to send the GetMetadata message to the parent contract
	body := cell.BeginCell().MustStoreUInt(2215781983, 32).MustStoreRef(metadataCell).EndCell() // Opcode should match your contract's `GetMetadata`

	err = w.Send(context.Background(), &wallet.Message{
		Mode: 1,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     parentContractAddr,
			Amount:      tlb.MustFromTON("0.05"), // Transaction fee
			Body:        body,
		},
	}, true)

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Println("Message sent to contract for metadata request")

	// Assuming the response contains the metadata, you'll need to parse the result from the chain.
	// Here we are simplifying this step for demonstration purposes:
	// You need to have proper contract listener setup to capture the metadata response.
	// Here we will simulate a mock retrieval of the metadata (to be replaced with actual listener):
	receivedMetadata := decodeCellToString(metadataCell)
	fmt.Printf("Received Metadata: %s\n", receivedMetadata)

	return nil
}

// decodeCellToString decodes a Cell to a readable string (example for metadata cells)
func decodeCellToString(c *cell.Cell) string {
	// This assumes metadata is stored as a snake string in the cell
	slice := c.BeginParse()
	metadata, err := slice.LoadStringSnake()
	if err != nil {
		fmt.Printf("Failed to decode cell: %v\n", err)
		return ""
	}
	return metadata
}

// Get the wallet from the seed mnemonic stored in mnemonics.txt
func getWallet(api *ton.APIClient) (*wallet.Wallet, error) {
	words, err := readSeedFromFile()
	if err != nil {
		return nil, err
	}

	// Create wallet from seed words
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
