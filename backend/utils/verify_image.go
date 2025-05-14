package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
)

// Takes raw image bytes and returns standardized JPEG bytes
func DecodeAndConvertToJPEG(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to decode image: %v", err)
		return nil, errors.New("invalid image format")
	}
	log.Printf("Decoded image format: %s", format)

	// Re-encode to JPEG
	var jpegBuffer bytes.Buffer
	if err := jpeg.Encode(&jpegBuffer, img, &jpeg.Options{Quality: 85}); err != nil {
		log.Printf("Failed to encode image as JPEG: %v", err)
		return nil, errors.New("failed to encode image")
	}

	return jpegBuffer.Bytes(), nil
}

// DecodeB64 decodes a base64 string and returns the byte array
func DecodeB64(data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Printf("Failed to decode base64 data: %v", err)
		return nil, errors.New("failed to decode base64 data")
	}
	return decodedData, nil
}
