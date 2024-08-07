package compare

func Compare() error {
	return nil
}

/*Membandingkan NFT
func Compare() error {
	nftfetchDir, err := getNFTFetchDir()
	if err != nil {
		return err
	}

	metahashPath := filepath.Join(nftfetchDir, "metadata", "metahash.txt")
	metahash, err := os.ReadFile(metahashPath)
	if err != nil {
		return fmt.Errorf("error membaca metahash: %v", err)
	}

	nfthashPath := filepath.Join(nftfetchDir, "nfthash", "nfthash.txt")
	nfthash, err := os.ReadFile(nfthashPath)
	if err != nil {
		return fmt.Errorf("file belum dihasilkan: %v", err)
	}

	if string(metahash) == string(nfthash) {
		fmt.Println("NFT terverifikasi.")
	} else {
		fmt.Println("NFT tidak terverifikasi.")
	}
	return nil
}
*/
