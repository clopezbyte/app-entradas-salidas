package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func GetSignedURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bucket := "app-entradas-salidas-merc"
	filename := r.URL.Query().Get("filename")
	contentType := r.URL.Query().Get("contentType")

	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "storage client error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		Method:         "PUT",
		Expires:        time.Now().Add(15 * time.Minute),
		ContentType:    contentType,
		GoogleAccessID: os.Getenv("GCS_SIGNER_EMAIL"),
		PrivateKey:     []byte(os.Getenv("GCS_PRIVATE_KEY")),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign URL: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url":       url,
		"publicUrl": fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filename),
	})
}
