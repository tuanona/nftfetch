package compare

import "fmt"

func Compare() error {
	fmt.Println("Error broh")
	return nil
}

/*import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

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

	metadataHash, err := readMetadataHash()
	if err != nil {
		return fmt.Errorf("failed to read metadata hash: %v", err)
	}

	// Assuming the NFT collection contract address is known
	collectionAddr := address.MustParseAddr("EQBQdiwDar5r7oMczFegkBw5Upa-ZyayQ_UMQyKtcrc4d85h")

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get current block: %v", err)
	}

	// Create a Cell for the metadata hash
	metadataCell := cell.BeginCell().
		MustStoreStringSnake(metadataHash).
		EndCell()

	// Call the getNftData method on the collection contract using the Cell
	res, err := api.RunGetMethod(context.Background(), block, collectionAddr, "getNftData", metadataCell)
	if err != nil {
		return fmt.Errorf("failed to run getNftData method: %v", err)
	}

	// Correct way to retrieve the stack values
	stack := res.GetStack()                   // Use GetStack() to access the execution result stack
	ownerAddress, err := stack.FetchAddress() // Fetch the first address from the stack
	if err != nil {
		return fmt.Errorf("failed to get owner address: %v", err)
	}

	if ownerAddress.String() == w.Address().String() {
		fmt.Printf("The NFT with metadata hash %s is owned by your wallet.\n", metadataHash)
	} else {
		fmt.Printf("The NFT with metadata hash %s is not owned by your wallet.\n", metadataHash)
	}

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
*/
