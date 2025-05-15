package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strings"

	"fmt"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"golang.org/x/net/context"
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

// UploadImageToGCS uploads a base64 encoded image to Google Cloud Storage
func UploadImageToGCS(ctx context.Context, client *storage.Client, bucket, folder, b64Data, uploadSource string) (string, error) {
	decoded, err := DecodeB64(b64Data)
	if err != nil {
		return "", fmt.Errorf("base64 decoding failed: %w", err)
	}
	contentType := http.DetectContentType(decoded)
	if !strings.HasPrefix(contentType, "image/") {
		return "", errors.New("invalid image content")
	}
	object := fmt.Sprintf("%s/%s.jpeg", folder, uuid.New().String())
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ContentType = contentType
	wc.Metadata = map[string]string{
		"upload-source":         uploadSource,
		"original-content-type": contentType,
	}
	if _, err := wc.Write(decoded); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object), nil
}
