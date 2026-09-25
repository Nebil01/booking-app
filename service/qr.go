package service

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenerateQRCode(data string, filePath string) error {
	// Ensure destination directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Encode data into a QR code
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Scale QR code
	scaledQR, err := barcode.Scale(qrCode, 300, 300)
	if err != nil {
		return fmt.Errorf("failed to scale QR code: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode as PNG
	if err := png.Encode(file, scaledQR); err != nil {
		return fmt.Errorf("failed to encode QR code image: %w", err)
	}

	fmt.Println("✅ QR Code saved at:", filePath)
	return nil
}
