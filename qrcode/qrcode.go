package qrcodeTerminal

import (
	"os"

	"github.com/mdp/qrterminal"
)

// NewQRCodeTerminal initializes a new QR code terminal with default configuration
func NewQRCodeTerminal() *QRCodeTerminal {
	return &QRCodeTerminal{}
}

// QRCodeTerminal is a struct that handles QR code generation and printing
type QRCodeTerminal struct{}

// Print generates and prints the QR code for the given data to the provided writer
func (qt *QRCodeTerminal) Print(data string) {
	config := qrterminal.Config{
		Level:     qrterminal.L, // Set error correction level (Low)
		Writer:    os.Stdout,    // Output to terminal (stdout)
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1, // Optional quiet zone around the QR code
	}
	qrterminal.GenerateWithConfig(data, config)
}

// Save generates a QR code for the given data and saves it to the specified file
func (qt *QRCodeTerminal) Save(data string, filePath string) error {
	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate QR code and write to the file
	config := qrterminal.Config{
		Level:     qrterminal.L, // Set error correction level (Low)
		Writer:    file,         // Output to file
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1, // Optional quiet zone around the QR code
	}
	qrterminal.GenerateWithConfig(data, config)
	return nil
}
