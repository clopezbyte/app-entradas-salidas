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
		log.Printf("Invalid header token: %v", err)

		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		log.Printf("Invalid token: %v", err)

		return
	}
	fmt.Println("Verified user ID:", token.UID)

	// Limit file size (5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v", err)
		return
	}

	// Extract the base64Data from the "evidencia_recepcion" object
	evidenciaObject := r.FormValue("evidencia_recepcion")           // Contain object as string
	log.Println("Raw evidencia_recepcion string:", evidenciaObject) // Debug how does it look like

	// Assuming the evidence object is passed as a JSON string, you can unmarshal it
	var evidencia map[string]interface{}
	if err := json.Unmarshal([]byte(evidenciaObject), &evidencia); err != nil {
		http.Error(w, "Failed to parse evidencia_recepcion", http.StatusBadRequest)
		log.Printf("Error parsing evidencia_recepcion: %v", err)
		return
	}

	// Extract the base64 data from the map
	b64, ok := evidencia["base64Data"].(string)
	if !ok || b64 == "" {
		http.Error(w, "Missing base64 image", http.StatusBadRequest)
		log.Println("Missing base64 image", err)
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
		log.Printf("GCS client error: %v", err)
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
		log.Printf("Error uploading image: %v", err)
		return
	}

	if err := wc.Close(); err != nil {
		http.Error(w, "Error finalizing image", http.StatusInternalServerError)
		log.Printf("Error finalizing image: %v", err)
		return
	}

	// Create public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	log.Printf("File uploaded successfully to: %s with content type: %s", imageURL, contentType)

	// Parse form values
	numRem, err := strconv.ParseInt(r.FormValue("numero_remision_factura"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_remision_factura", http.StatusBadRequest)
		log.Printf("Error parsing numero_remision_factura: %v", err)
		return
	}

	//Validations Block
	cant, err := strconv.ParseInt(r.FormValue("cantidad"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		log.Printf("Error parsing quantity: %v", err)
		return
	}

	cliente := r.FormValue("cliente")
	if cliente == "null" { //Not a devolución rma case
		cliente = "N/A"
	}

	fechaRecepcion, err := time.Parse("2006-01-02T15:04:05.000-0700", r.FormValue("fecha_recepcion"))
	if err != nil {
		http.Error(w, "Invalid fecha_recepcion format", http.StatusBadRequest)
		log.Printf("Error parsing fecha_recepcion: %v", err)
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
		log.Printf("Firestore error: %v", err)
		return
	}
	defer fsClient.Close()

	//Email block
	if entrada.TipoDelivery == "Devolución (RMA)" && entrada.Cliente != "" {
		utils.HandleClientEmailNotification(ctx, fsClient, entrada)
		//replaced goroutine for testing
	}

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

func QueryNumRem(w http.ResponseWriter, r *http.Request) {
	// Get token from header using utils
	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("Invalid header token: %v", err)

		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		log.Printf("Invalid token: %v", err)

		return
	}
	fmt.Println("Verified user ID:", token.UID)

	// Parse numero remision factura
	numRem, err := strconv.ParseInt(r.FormValue("numero_remision_factura"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_remision_factura format", http.StatusBadRequest)
		log.Printf("Error parsing numero_remision_factura: %v", err)
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

	// Build query and query firestore
	query := fsClient.Collection("entradas").Where("NumeroRemisionFactura", "==", numRem).Limit(1)
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

func HandleASNSubmit(w http.ResponseWriter, r *http.Request) {
	// Get token from header using utils
	authHeader := r.Header.Get("Authorization")
	idToken, err := utils.GetTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("Invalid header token: %v", err)

		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		log.Printf("Invalid token: %v", err)

		return
	}
	fmt.Println("Verified user ID:", token.UID)

	//Parse ASN update date
	FechaAjusteASN, err := time.Parse("2006-01-02", r.FormValue("fecha_ajuste_asn"))
	if err != nil {
		http.Error(w, "Invalid fecha_recepcion format", http.StatusBadRequest)
		log.Printf("Error parsing fecha_recepcion: %v", err)
		return
	}

	// Parse numero remision factura
	numRem, err := strconv.ParseInt(r.FormValue("numero_remision_factura"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_remision_factura format", http.StatusBadRequest)
		log.Printf("Error parsing numero_remision_factura: %v", err)
		return
	}

	// Construct ASN struct
	asn := models.ASN{
		NumeroRemisionFactura: numRem,
		ASN:                   r.FormValue("asn"),
		FechaAjusteASN:        FechaAjusteASN,
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

	// Update entrada with ASN
	// Find the document in "entradas" with the given numero_remision_factura
	iter := fsClient.Collection("entradas").Where("NumeroRemisionFactura", "==", numRem).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		http.Error(w, "No matching entrada found", http.StatusNotFound)
		log.Printf("No matching entrada found for numero_remision_factura: %d", numRem)
		return
	}

	// Update the ASN and FechaAjusteASN fields
	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "ASN", Value: asn.ASN},
		{Path: "FechaAjusteASN", Value: asn.FechaAjusteASN},
	})
	if err != nil {
		http.Error(w, "Failed to update ASN", http.StatusInternalServerError)
		log.Printf("Failed to update ASN: %v", err)
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
		log.Printf("Invalid header token: %v", err)
		return
	}

	// Verify the token using the firebase package
	token, err := utils.VerifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		log.Printf("Invalid token: %v", err)
		return
	}
	fmt.Println("Verified user ID:", token.UID)
	////////////////////////////////////////////////////////////////////////////

	// Limit file size (5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v", err)
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Extract and parse evidencia_recepcion as JSON string
	evidenciaObject := r.FormValue("evidencia_salida")
	log.Println("Raw evidencia_salida string:", evidenciaObject)

	var evidencia map[string]interface{}
	if err := json.Unmarshal([]byte(evidenciaObject), &evidencia); err != nil {
		http.Error(w, "Failed to parse evidencia_salida", http.StatusBadRequest)
		log.Printf("Error parsing evidencia_salida: %v", err)
		return
	}

	b64, ok := evidencia["base64Data"].(string)
	if !ok || b64 == "" {
		http.Error(w, "Missing base64 image", http.StatusBadRequest)
		log.Printf("Missing base64 image: %v", err)
		return
	}

	// Decode and upload evidencia_salida
	bucket := "app-entradas-salidas-merc"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "GCS client error", http.StatusInternalServerError)
		log.Printf("GCS client error: %v", err)
		return
	}
	defer client.Close()

	imageURL, err := utils.UploadImageToGCS(ctx, client, bucket, "evidencias_salidas", b64, "retool-app-salidas")
	if err != nil {
		http.Error(w, "Image upload failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Image upload failed: %v", err)
		return
	}
	log.Printf("File uploaded successfully to: %s", imageURL)

	////////////////////////////////////////////////////////////////////////////

	//Extract and parse signature as JSON string
	signatureObject := r.FormValue("firma_persona_recoge")
	log.Println("Raw firma_persona_recoge string:", signatureObject)
	var firma map[string]interface{}
	if err := json.Unmarshal([]byte(signatureObject), &firma); err != nil {
		http.Error(w, "Failed to parse firma_persona_recoge", http.StatusBadRequest)
		log.Printf("Error parsing firma_persona_recoge: %v", err)
		return
	}
	b64Firma, ok := firma["base64Data"].(string)
	if !ok || b64Firma == "" {
		http.Error(w, "Missing base64 firma", http.StatusBadRequest)
		log.Printf("Missing base64 firma: %v", err)
		return
	}

	// Decode and upload firma_persona_recoge
	signatureImageURL, err := utils.UploadImageToGCS(ctx, client, bucket, "evidencias_salidas/salidas_firmas", b64Firma, "retool-app-salidas")
	if err != nil {
		http.Error(w, "Signature upload failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Signature upload failed: %v", err)
		return
	}
	log.Printf("Signature uploaded successfully to: %s", signatureImageURL)

	////////////////////////////////////////////////////////////////////////////

	// Parse form values
	numOrdenCons, err := strconv.ParseInt(r.FormValue("numero_orden_consecutivo"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid numero_orden_consecutivo", http.StatusBadRequest)
		log.Printf("Error parsing numero_orden_consecutivo: %v", err)
		return
	}

	fechaSalida, err := time.Parse("2006-01-02T15:04:05.000-0700", r.FormValue("fecha_salida"))
	if err != nil {
		http.Error(w, "Invalid fecha_salida format", http.StatusBadRequest)
		log.Printf("Error parsing fecha_salida: %v", err)
		return
	}

	// Construct Salidas struct
	salida := models.Salidas{
		BodegaSalida:           r.FormValue("bodega_salida"),
		ProveedorSalida:        r.FormValue("proveedor_salida"),
		NumeroOrdenConsecutivo: numOrdenCons,
		PersonaEntrega:         r.FormValue("persona_entrega"),
		PersonaRecoge:          r.FormValue("persona_recoge"),
		FirmaPersonaRecoge:     signatureImageURL,
		FechaSalida:            fechaSalida,
		EvidenciaSalida:        imageURL,
		Comentarios:            r.FormValue("comentarios"),
	}

	// Save to Firestore
	fsClient, err := firestore.NewClientWithDatabase(ctx, "b-materials", "app-in-out-good")
	if err != nil {
		http.Error(w, "Firestore error", http.StatusInternalServerError)
		log.Printf("Firestore error: %v", err)
		return
	}
	defer fsClient.Close()

	// Add entrada form as new document to "salidas" collection
	_, _, err = fsClient.Collection("salidas").Add(ctx, salida)
	if err != nil {
		log.Printf("Error saving to Firestore: %v", err)
		http.Error(w, fmt.Sprintf("Error saving to Firestore: %v", err), http.StatusInternalServerError)
		log.Printf("Error saving to Firestore: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Salida submitted successfully."}`))
}
