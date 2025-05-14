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

	// Build query and query firestore
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

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
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

	// Build query and query firestore
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
		var entrada models.SalidasData
		if err := doc.DataTo(&entrada); err != nil {
			log.Printf("Error parsing Firestore document: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		results = append(results, entrada)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}
