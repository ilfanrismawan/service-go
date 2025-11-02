package utils

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"service/internal/config"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code image for the given text
// size: QR code size in pixels (default: 256)
// Returns PNG image bytes
func GenerateQRCode(text string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}

	// Generate QR code with error correction level M (about 15% error correction)
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// Get PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, qr.Image(size)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateQRCodeForOrder generates a QR code for a service order
// Contains order number and URL for tracking
func GenerateQRCodeForOrder(orderNumber string) ([]byte, error) {
	baseURL := config.Config.BaseURL
	if baseURL == "" {
		baseURL = "https://service.example.com"
	}

	qrText := baseURL + "/orders/" + orderNumber
	return GenerateQRCode(qrText, 256)
}

// GenerateQRCodeWriter generates QR code and writes to io.Writer
func GenerateQRCodeWriter(text string, size int, writer io.Writer) error {
	qrBytes, err := GenerateQRCode(text, size)
	if err != nil {
		return err
	}
	_, err = writer.Write(qrBytes)
	return err
}

// DecodeQRCode decodes a QR code from an image
// Note: This requires a QR code decoder library like github.com/makiuchi-d/gozxing
// For now, this is a placeholder
func DecodeQRCode(img image.Image) (string, error) {
	// TODO: Implement QR code decoding
	// This would require adding a QR decoder library
	return "", nil
}
