package qrcode

import (
	"fmt"
	"image/png"
	"os"

	"github.com/skip2/go-qrcode"
)

type QRCode struct{}

func NewQRCode() *QRCode {
	return &QRCode{}
}

func (q *QRCode) GenerateAndSave(data, filePath string) error {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, qr.Image(256)); err != nil {
		return fmt.Errorf("failed to encode QR code as PNG: %v", err)
	}

	return nil
}

func (q *QRCode) Print(data string) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		fmt.Printf("Failed to generate QR code: %v\n", err)
		return
	}

	fmt.Println(qr.ToSmallString(false))
}
