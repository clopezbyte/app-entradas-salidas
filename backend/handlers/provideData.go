package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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

	// Get and validate query params
	monthStr := r.FormValue("month")
	yearStr := r.FormValue("year")
	if monthStr == "" || yearStr == "" {
		http.Error(w, "Missing 'month' or 'year' query parameter", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > time.Now().Year()+1 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // first day of the next month

	// Query Firestore for EntradasData within the specified date range
	query := fsClient.Collection("entradas").
		Where("FechaRecepcion", ">=", startDate).
		Where("FechaRecepcion", "<", endDate).
		OrderBy("FechaRecepcion", firestore.Asc)
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

	// Get and validate query params
	monthStr := r.FormValue("month")
	yearStr := r.FormValue("year")
	if monthStr == "" || yearStr == "" {
		http.Error(w, "Missing 'month' or 'year' query parameter", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > time.Now().Year()+1 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // first day of the next month

	// Query Firestore for SalidasData within the specified date range
	query := fsClient.Collection("salidas").
		Where("FechaSalida", ">=", startDate).
		Where("FechaSalida", "<", endDate).
		OrderBy("FechaSalida", firestore.Desc)
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

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

var (
	customerIDsCache      []string
	customerIDsCacheTime  time.Time
	customerIDsCacheMutex sync.Mutex
	cacheDuration         = 15 * time.Minute
)

func HandleProvideCustomers(w http.ResponseWriter, r *http.Request) {
	customerIDsCacheMutex.Lock()
	defer customerIDsCacheMutex.Unlock()

	// Serve from cache if not expired
	if time.Since(customerIDsCacheTime) < cacheDuration && customerIDsCache != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customerIDsCache)
		return
	}

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

	// Query Firestore for SalidasData within the specified date range
	//cached
	query := fsClient.Collection("customers")
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error querying Firestore: %v", err)
		http.Error(w, "Error querying Firestore", http.StatusInternalServerError)
		return
	}

	// Parse Firestore documents into a slice of Customers
	var ids []string
	for _, doc := range docs {
		ids = append(ids, doc.Ref.ID)
	}

	// Update cache
	customerIDsCache = ids
	customerIDsCacheTime = time.Now()

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}
