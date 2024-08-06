package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"nftfetch/wallet"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdxyw/generativeart"
	"github.com/jdxyw/generativeart/arts"
	"github.com/nfnt/resize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

// Metadata struct untuk menyimpan data metadata
type Metadata struct {
	CPUId      string `json:"cpu_id"`
	CPUModel   string `json:"cpu_model"`
	MBSerial   string `json:"mb_serial"`
	DiskSerial string `json:"disk_serial"`
	MACAddress string `json:"mac_address"`
	BIOSUUID   string `json:"bios_uuid"`
}

// Karakter ASCII yang digunakan untuk memetakan nilai skala abu-abu
const asciiChars = " .:-=+*#%@"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Command: nftfetch <generate|compare|wallet>")
		return
	}

	command := os.Args[1]

	switch command {
	case "generate":
		if err := generate(); err != nil {
			fmt.Printf("Error saat menghasilkan NFT: %v\n", err)
		}
	case "compare":
		if err := compare(); err != nil {
			fmt.Printf("Error saat membandingkan NFT: %v\n", err)
		}
	case "wallet":
		wallet.CreateWallet()
	default:
		fmt.Println("Perintah tidak valid. Penggunaan: nftfetch <generate|compare|wallet>")
	}
}

// Mendapatkan metadata perangkat
func getDeviceMetadata() (Metadata, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return Metadata{}, err
	}

	diskInfo, err := disk.Partitions(true)
	if err != nil || len(diskInfo) == 0 {
		return Metadata{}, err
	}

	macAddrs, err := net.Interfaces()
	if err != nil || len(macAddrs) == 0 {
		return Metadata{}, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return Metadata{}, err
	}

	metadata := Metadata{
		CPUId:      cpuInfo[0].VendorID,
		CPUModel:   cpuInfo[0].ModelName,
		MBSerial:   hostInfo.HostID,
		DiskSerial: diskInfo[0].Device,
		MACAddress: macAddrs[0].HardwareAddr,
		BIOSUUID:   hostInfo.HostID,
	}

	return metadata, nil
}

// Menyimpan metadata sebagai JSON
func saveMetadataAsJSON(metadata Metadata, filePath string) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// Menghitung hash SHA512
func calculateSHA512(data []byte) string {
	hash := sha512.Sum512(data)
	return hex.EncodeToString(hash[:])
}

// Menyimpan data ke file
func saveToFile(data, filePath string) error {
	return os.WriteFile(filePath, []byte(data), 0644)
}

// Menghasilkan seni piksel berdasarkan hash
func generatePixelArt(hash string, outputPath string) error {
	colors := generateColorsFromHash(hash)
	c := generativeart.NewCanva(512, 512)
	c.SetBackground(color.RGBA{0x1a, 0x06, 0x33, 0xFF})
	c.FillBackground()
	c.SetColorSchema(colors)
	c.Draw(arts.NewContourLine(512))
	return c.ToPNG(outputPath)
}

// Menghasilkan warna dari hash
func generateColorsFromHash(hash string) []color.RGBA {
	var colors []color.RGBA
	for i := 0; i < len(hash); i += 6 {
		if i+6 <= len(hash) {
			r, _ := hex.DecodeString(hash[i : i+2])
			g, _ := hex.DecodeString(hash[i+2 : i+4])
			b, _ := hex.DecodeString(hash[i+4 : i+6])
			colors = append(colors, color.RGBA{r[0], g[0], b[0], 0xFF})
		}
	}
	return colors
}

// Menghasilkan NFT
func generate() error {
	nftfetchDir, err := getNFTFetchDir()
	if err != nil {
		return err
	}

	metadata, err := getDeviceMetadata()
	if err != nil {
		return err
	}

	metadataPath := filepath.Join(nftfetchDir, "metadata", "metadata.json")
	if err := saveMetadataAsJSON(metadata, metadataPath); err != nil {
		return err
	}

	metadataJSON, err := os.ReadFile(metadataPath)
	if err != nil {
		return err
	}

	metahash := calculateSHA512(metadataJSON)
	metahashPath := filepath.Join(nftfetchDir, "metadata", "metahash.txt")
	if err := saveToFile(metahash, metahashPath); err != nil {
		return err
	}

	nftPath := filepath.Join(nftfetchDir, "nft", "nft.png")
	if err := generatePixelArt(metahash, nftPath); err != nil {
		return err
	}

	return displayNFT(nftPath, metadataPath)
}

// Membandingkan NFT
func compare() error {
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

// Mendapatkan direktori NFTFetch
func getNFTFetchDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	nftfetchDir := filepath.Join(homeDir, ".nftfetch")
	if _, err := os.Stat(nftfetchDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(nftfetchDir, "metadata"), 0755); err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Join(nftfetchDir, "nft"), 0755); err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Join(nftfetchDir, "nfthash"), 0755); err != nil {
			return "", err
		}
	}
	return nftfetchDir, nil
}

// Menampilkan NFT
func displayNFT(nftPath, metadataPath string) error {
	// Baca metadata
	metadataJSON, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("error membaca metadata: %v", err)
	}

	var metadata Metadata
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		return fmt.Errorf("error mengurai metadata: %v", err)
	}

	frameWidth := 60

	// Tampilkan NFT Generated
	fmt.Println("┌" + strings.Repeat("─", frameWidth) + "┐")
	fmt.Printf("│%s│\n", centerText("NFT Generated", frameWidth))
	fmt.Println("└" + strings.Repeat("─", frameWidth) + "┘")

	// Tampilkan ASCII
	fmt.Println("┌" + strings.Repeat("─", frameWidth) + "┐")
	fmt.Printf("│%s│\n", centerText("ASCII", frameWidth))
	fmt.Println("├" + strings.Repeat("─", frameWidth) + "┤")
	if err := generateASCII(nftPath); err != nil {
		return fmt.Errorf("error menghasilkan ASCII art: %v", err)
	}
	fmt.Println("└" + strings.Repeat("─", frameWidth) + "┘")

	// Tampilkan Metadata
	fmt.Println("┌" + strings.Repeat("─", frameWidth) + "┐")
	fmt.Printf("│%s│\n", centerText("METADATA", frameWidth))
	fmt.Println("├" + strings.Repeat("─", frameWidth) + "┤")
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("CPU ID: %s", metadata.CPUId))
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("CPU Model: %s", metadata.CPUModel))
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("Motherboard Serial: %s", metadata.MBSerial))
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("Disk Serial: %s", metadata.DiskSerial))
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("MAC Address: %s", metadata.MACAddress))
	fmt.Printf("│ %-58s │\n", fmt.Sprintf("BIOS UUID: %s", metadata.BIOSUUID))
	fmt.Println("└" + strings.Repeat("─", frameWidth) + "┘")

	fmt.Println("This machine is belong to you!")

	return nil
}

// Menghasilkan ASCII art dari gambar
func generateASCII(imagePath string) error {
	// Buka file gambar
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("error membuka file gambar: %v", err)
	}
	defer file.Close()

	// Decode gambar
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error mendecode gambar: %v", err)
	}

	// Resize gambar menjadi 58x29 karakter (untuk menyesuaikan dengan frame)
	img = resize.Resize(58, 29, img, resize.Lanczos3)

	// Konversi gambar ke skala abu-abu dan petakan ke karakter ASCII
	asciiArt := convertToASCII(img)

	// Cetak ASCII art
	printASCII(asciiArt)

	return nil
}

// Mengkonversi gambar ke ASCII
func convertToASCII(img image.Image) string {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var asciiArt string
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			char := mapGrayToASCII(gray.Y)
			asciiArt += string(char)
		}
		asciiArt += "\n"
	}
	return asciiArt
}

// Memetakan nilai skala abu-abu ke karakter ASCII
func mapGrayToASCII(grayValue uint8) byte {
	index := int((float64(grayValue) / 255.0) * float64(len(asciiChars)-1))
	return asciiChars[index]
}

// Mencetak ASCII art
func printASCII(asciiArt string) {
	lines := strings.Split(asciiArt, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("│ %-58s │\n", line)
		}
	}
}

// Memusatkan teks dalam lebar tertentu
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-len(text))
}
