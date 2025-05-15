package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/clopezbyte/app-entradas-salidas/models"
	"github.com/clopezbyte/app-entradas-salidas/utils"
)

func HandleProvideEntradasData(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // Error message from the utility function
		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	fmt.Println("Verified user ID:", token.UID)

	// Use the request's context for proper cancellation
	ctx := r.Context()

	// Initialize Firestore client
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		log.Printf("Error initializing Firestore client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer fsClient.Close()

	// Use a channel to collect results
	resultsChan := make(chan []models.EntradasData)
	errorChan := make(chan error)

	// Build query and query firestore
	go func() {
		query := fsClient.Collection("entradas").OrderBy("FechaRecepcion", firestore.Desc).Limit(10)
		docs, err := query.Documents(ctx).GetAll()
		if err != nil {
			log.Printf("Error querying Firestore: %v", err)
			http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
			return
		}

		// Parse Firestore documents into a slice of EntradasData
		var results []models.EntradasData
		for _, doc := range docs {
			var entrada models.EntradasData
			if err := doc.DataTo(&entrada); err != nil {
				log.Printf("Error parsing Firestore document: %v", err)
				http.Error(w, "Error processing data", http.StatusInternalServerError)
				return
			}
			results = append(results, entrada)
		}
		resultsChan <- results
	}()

	// Wait for the goroutine to finish
	select {
	case results := <-resultsChan:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	case err := <-errorChan:
		log.Printf("Error querying Firestore: %v", err)
		http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
	case <-ctx.Done():
		log.Printf("Request context canceled")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	}

}

func HandleProvideSalidasData(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // Error message from the utility function
		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	fmt.Println("Verified user ID:", token.UID)

	// Use the request's context for proper cancellation
	ctx := r.Context()

	// Initialize Firestore client
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		log.Printf("Error initializing Firestore client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer fsClient.Close()

	// Use a channel to collect results
	resultsChan := make(chan []models.SalidasData)
	errorChan := make(chan error)

	// Build query and query firestore
	go func() {
		query := fsClient.Collection("salidas").OrderBy("FechaSalida", firestore.Desc).Limit(10)
		docs, err := query.Documents(ctx).GetAll()
		if err != nil {
			log.Printf("Error querying Firestore: %v", err)
			http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
			return
		}

		// Parse Firestore documents into a slice of EntradasData
		var results []models.SalidasData
		for _, doc := range docs {
			var salida models.SalidasData
			if err := doc.DataTo(&salida); err != nil {
				log.Printf("Error parsing Firestore document: %v", err)
				http.Error(w, "Error processing data", http.StatusInternalServerError)
				return
			}
			results = append(results, salida)
		}
		resultsChan <- results
	}()

	// Return JSON response
	// Wait for the goroutine to finish
	select {
	case results := <-resultsChan:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	case err := <-errorChan:
		log.Printf("Error querying Firestore: %v", err)
		http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
	case <-ctx.Done():
		log.Printf("Request context canceled")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	}

}
