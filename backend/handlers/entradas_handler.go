package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/clopezbyte/app-entradas-salidas/utils"

	"github.com/clopezbyte/app-entradas-salidas/models"
)

func HandleEntradasSubmit(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization token from the header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
		return
	}

	// Strip the "Bearer " prefix from the token
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}
	idToken := authHeader[7:] // Extract the actual token

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	fmt.Println("Verified user ID:", token.UID)

	// Limit file size (5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract file
	file, handler, err := r.FormFile("evidencia_recepcion")
	if err != nil {
		http.Error(w, "Error reading image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload to GCS
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "GCS client error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	bucket := "app-entradas-salidas-merc"
	object := fmt.Sprintf("evidencias/%d_%s", time.Now().UnixNano(), handler.Filename)
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	if _, err = io.Copy(wc, file); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Error finalizing image", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)

	// Parse form values
	numRem, err := strconv.ParseInt(r.FormValue("numero_remision_factura"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_remision_factura", http.StatusBadRequest)
		return
	}

	cant, err := strconv.ParseInt(r.FormValue("cantidad"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	fechaRecepcion, err := time.Parse("2006-01-02", r.FormValue("fecha_recepcion"))
	if err != nil {
		http.Error(w, "Invalid fecha_recepcion format", http.StatusBadRequest)
		return
	}

	// Construct Entradas struct
	entrada := models.Entradas{
		TipoDelivery:          r.FormValue("tipo_delivery"),
		BodegaRecepcion:       r.FormValue("bodega_recepcion"),
		ProveedorRecepcion:    r.FormValue("proveedor_recepcion"),
		NumeroRemisionFactura: numRem,
		PersonaRecepcion:      r.FormValue("persona_recepcion"),
		FechaRecepcion:        fechaRecepcion,
		EvidenciaRecepcion:    imageURL,
		Cantidad:              cant,
		Comentarios:           r.FormValue("comentarios"),
	}

	// Save to Firestore
	// Specify project and target db, otherwise it will try to target default
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		http.Error(w, "Firestore error", http.StatusInternalServerError)
		return
	}
	defer fsClient.Close()

	// Add entrada form as new document to "entradas" collection
	_, _, err = fsClient.Collection("entradas").Add(ctx, entrada)
	if err != nil {
		log.Printf("Error saving to Firestore: %v", err)
		http.Error(w, fmt.Sprintf("Error saving to Firestore: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Entrada submitted successfully."}`))
}
