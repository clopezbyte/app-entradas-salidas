package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/clopezbyte/app-entradas-salidas/models"
	"github.com/clopezbyte/app-entradas-salidas/utils"
)

func HandleEntradasSubmit(w http.ResponseWriter, r *http.Request) {
	// Get token from header using utils
	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	// Extract the base64Data from the "evidencia_recepcion" object
	evidenciaObject := r.FormValue("evidencia_recepcion")           // Contain object as string
	log.Println("Raw evidencia_recepcion string:", evidenciaObject) // Debug how does it look like

	// Assuming the evidence object is passed as a JSON string, you can unmarshal it
	var evidencia map[string]interface{}
	if err := json.Unmarshal([]byte(evidenciaObject), &evidencia); err != nil {
		http.Error(w, "Failed to parse evidencia_recepcion", http.StatusBadRequest)
		return
	}

	// Extract the base64 data from the map
	b64, ok := evidencia["base64Data"].(string)
	if !ok || b64 == "" {
		http.Error(w, "Missing base64 image", http.StatusBadRequest)
		return
	}

	// Decode the base64 string
	decoded, err := utils.DecodeB64(b64)
	if err != nil {
		http.Error(w, "B64 decoding error", http.StatusBadRequest)
		return
	}

	contentType := http.DetectContentType(decoded)
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Invalid image content", http.StatusBadRequest)
		return
	}

	// Upload to GCS
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "GCS client error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	bucket := "app-entradas-salidas-merc"
	object := fmt.Sprintf("evidencias_entradas/%s.jpeg", uuid.New().String())
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	// Set proper content type and metadata
	wc.ContentType = contentType
	wc.Metadata = map[string]string{
		"upload-source":         "retool-app-entradas",
		"original-content-type": contentType,
	}

	// Write the file data
	if _, err := wc.Write(decoded); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	if err := wc.Close(); err != nil {
		http.Error(w, "Error finalizing image", http.StatusInternalServerError)
		return
	}

	// Create public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	log.Printf("File uploaded successfully to: %s with content type: %s", imageURL, contentType)

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

	cliente := r.FormValue("cliente")
	if cliente == "null" { //Not a devolución rma case
		cliente = "N/A"
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
		Cliente:               cliente,
		NumeroRemisionFactura: numRem,
		PersonaRecepcion:      r.FormValue("persona_recepcion"),
		FechaRecepcion:        fechaRecepcion,
		EvidenciaRecepcion:    imageURL,
		Cantidad:              cant,
		Comentarios:           r.FormValue("comentarios"),
	}

	// Save to Firestore
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

func HandleSalidasSubmit(w http.ResponseWriter, r *http.Request) {
	// Get token from header using utils
	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	// Extract and parse evidencia_recepcion as JSON string
	evidenciaObject := r.FormValue("evidencia_salida")
	log.Println("Raw evidencia_salida string:", evidenciaObject)

	var evidencia map[string]interface{}
	if err := json.Unmarshal([]byte(evidenciaObject), &evidencia); err != nil {
		http.Error(w, "Failed to parse evidencia_salida", http.StatusBadRequest)
		return
	}

	b64, ok := evidencia["base64Data"].(string)
	if !ok || b64 == "" {
		http.Error(w, "Missing base64 image", http.StatusBadRequest)
		return
	}

	// Decode and validate image
	decoded, err := utils.DecodeB64(b64)
	if err != nil {
		http.Error(w, "B64 decoding error", http.StatusBadRequest)
		return
	}

	contentType := http.DetectContentType(decoded)
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Invalid image content", http.StatusBadRequest)
		return
	}

	// Upload to GCS
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "GCS client error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	bucket := "app-entradas-salidas-merc"
	object := fmt.Sprintf("evidencias_salidas/%s.jpeg", uuid.New().String())
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)

	// Set proper content type and metadata
	wc.ContentType = contentType
	wc.Metadata = map[string]string{
		"upload-source":         "retool-app-salidas",
		"original-content-type": contentType,
	}

	if _, err := wc.Write(decoded); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	if err := wc.Close(); err != nil {
		http.Error(w, "Error finalizing image", http.StatusInternalServerError)
		return
	}

	// Create public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	log.Printf("File uploaded successfully to: %s with content type: %s", imageURL, contentType)

	// Parse form values
	numOrdenCons, err := strconv.ParseInt(r.FormValue("numero_orden_consecutivo"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_orden_consecutivo", http.StatusBadRequest)
		return
	}

	fechaSalida, err := time.Parse("2006-01-02", r.FormValue("fecha_salida"))
	if err != nil {
		http.Error(w, "Invalid fecha_salida format", http.StatusBadRequest)
		return
	}

	// Construct Salidas struct
	salida := models.Salidas{
		BodegaSalida:           r.FormValue("bodega_salida"),
		ProveedorSalida:        r.FormValue("proveedor_salida"),
		NumeroOrdenConsecutivo: numOrdenCons,
		PersonaEntrega:         r.FormValue("persona_entrega"),
		FechaSalida:            fechaSalida,
		EvidenciaSalida:        imageURL,
		Comentarios:            r.FormValue("comentarios"),
	}

	// Save to Firestore
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		http.Error(w, "Firestore error", http.StatusInternalServerError)
		return
	}
	defer fsClient.Close()

	// Add entrada form as new document to "salidas" collection
	_, _, err = fsClient.Collection("salidas").Add(ctx, salida)
	if err != nil {
		log.Printf("Error saving to Firestore: %v", err)
		http.Error(w, fmt.Sprintf("Error saving to Firestore: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Salida submitted successfully."}`))
}
