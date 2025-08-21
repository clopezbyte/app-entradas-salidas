package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/clopezbyte/app-entradas-salidas/models"
	"google.golang.org/api/iterator"
)

func HandleRmaQuery(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	apiKey := os.Getenv("API_KEY")
	token := r.Header.Get("Authorization")
	if token != "Bearer "+apiKey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse rma number string from request body
	var rmaRequest models.RmaRequest
	if err := json.NewDecoder(r.Body).Decode(&rmaRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Initialize Firestore client
	ctx := context.Background()
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		http.Error(w, "Firestore error", http.StatusInternalServerError)
		log.Printf("Firestore error: %v", err)
		return
	}
	defer fsClient.Close()

	// Query Firestore
	iter := fsClient.Collection("entradas").
		Where("ASN", "==", rmaRequest.Rma).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		http.Error(w, "RMA or ASN not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error querying Firestore: %v", err)
		http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
		return
	}

	// Map Firestore doc to response struct
	var rmaResponse models.RmaResponse
	if err := doc.DataTo(&rmaResponse); err != nil {
		log.Printf("Error parsing Firestore document: %v", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rmaResponse); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

}
