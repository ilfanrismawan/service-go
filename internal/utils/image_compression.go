package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"path/filepath"
)

// CompressImage compresses an image to reduce file size
// maxWidth and maxHeight specify maximum dimensions (0 = no limit)
// quality is for JPEG compression (1-100, default 85)
func CompressImage(input io.Reader, output io.Writer, maxWidth, maxHeight int, quality int) error {
	if quality <= 0 || quality > 100 {
		quality = 85
	}

	// Decode image
	img, format, err := image.Decode(input)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize if needed
	if maxWidth > 0 || maxHeight > 0 {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Calculate new dimensions while maintaining aspect ratio
		if maxWidth > 0 && width > maxWidth {
			ratio := float64(maxWidth) / float64(width)
			width = maxWidth
			height = int(float64(height) * ratio)
		}
		if maxHeight > 0 && height > maxHeight {
			ratio := float64(maxHeight) / float64(height)
			height = maxHeight
			width = int(float64(width) * ratio)
		}

		// Simple resize (for production, use a proper image library like github.com/disintegration/imaging)
		// For now, we'll just encode with compression
	}

	// Encode with compression
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(output, img, &jpeg.Options{Quality: quality})
	case "png":
		return png.Encode(output, img)
	default:
		// Default to JPEG
		return jpeg.Encode(output, img, &jpeg.Options{Quality: quality})
	}
}

// CompressImageBytes compresses image bytes and returns compressed bytes
func CompressImageBytes(input []byte, maxWidth, maxHeight int, quality int) ([]byte, error) {
	inputReader := bytes.NewReader(input)
	var outputBuffer bytes.Buffer

	if err := CompressImage(inputReader, &outputBuffer, maxWidth, maxHeight, quality); err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}

// GetImageFormat returns image format from filename or mime type
func GetImageFormat(filename string, contentType string) string {
	if contentType != "" {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			return filepath.Ext(exts[0])
		}
	}
	return filepath.Ext(filename)
}

// IsValidImageFormat checks if the file format is supported
func IsValidImageFormat(filename string) bool {
	ext := filepath.Ext(filename)
	validFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	return validFormats[ext]
}

