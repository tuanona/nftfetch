package qrcodeTerminal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/mdp/qrterminal/v3"
)

// QRCodeTerminal is a struct that handles QR code generation and printing
type QRCodeTerminal struct{}

// NewQRCodeTerminal initializes a new QR code terminal with default configuration
func NewQRCodeTerminal() *QRCodeTerminal {
	return &QRCodeTerminal{}
}

// Print generates and prints the QR code for the given data to the provided writer
func (qt *QRCodeTerminal) Print(data string) {
	config := qrterminal.Config{
		Level:      qrterminal.L, // Set error correction level (Low)
		Writer:     os.Stdout,    // Output to terminal (stdout)
		BlackChar:  qrterminal.BLACK,
		WhiteChar:  qrterminal.WHITE,
		QuietZone:  1,     // Optional quiet zone around the QR code
		HalfBlocks: false, // Use full size blocks
	}
	qrterminal.GenerateWithConfig(data, config)
}

// Save generates a QR code for the given data and saves it to the specified file as a PNG
func (qt *QRCodeTerminal) Save(data string, filePath string) error {
	// Buffer to hold the QR code data
	var buf bytes.Buffer

	// Generate QR code and write to the buffer
	config := qrterminal.Config{
		Level:      qrterminal.L, // Set error correction level (Low)
		Writer:     &buf,         // Output to buffer
		BlackChar:  qrterminal.BLACK,
		WhiteChar:  qrterminal.WHITE,
		QuietZone:  1,     // Optional quiet zone around the QR code
		HalfBlocks: false, // Use full size blocks
	}
	qrterminal.GenerateWithConfig(data, config)

	// Convert buffer data to an image
	img, err := createImageFromBuffer(buf.Bytes())
	if err != nil {
		return err
	}

	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the image as PNG and write to the file
	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// createImageFromBuffer creates an image.Image from the QR code buffer data
func createImageFromBuffer(buf []byte) (image.Image, error) {
	// Determine the size of the QR code from the buffer length
	lines := bytes.Split(buf, []byte{'\n'})
	height := len(lines)
	width := len(lines[0])

	img := image.NewGray(image.Rect(0, 0, width, height))

	for y, line := range lines {
		for x, char := range line {
			if char == qrterminal.BLACK[0] {
				img.SetGray(x, y, color.Gray{Y: 0})
			} else {
				img.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}

	return img, nil

}
