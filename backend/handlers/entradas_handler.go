package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/clopezbyte/app-entradas-salidas/models"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

func HandleEntradasSubmit(w http.ResponseWriter, r *http.Request) {
	// Limit file size (e.g., 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract file
	file, handler, err := r.FormFile("evidencia")
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

	// Copy file data to GCS writer
	if _, err = io.Copy(wc, file); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Error finalizing image", http.StatusInternalServerError)
		return
	}

	// Construct public URL (if bucket is public)
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

	entrada := models.Entradas{
		TipoDelivery:          r.FormValue("tipo_delivery"),
		BodegaRecepcion:       r.FormValue("bodega_recepcion"),
		ProveedorRecepcion:    r.FormValue("proveedor_recepcion"),
		NumeroRemisionFactura: numRem,
		PersonaRecepcion:      r.FormValue("persona_recepcion"),
		FechaRecepcion:        r.FormValue("fecha_recepcion"),
		EvidenciaRecepcion:    imageURL,
		Cantidad:              cant,
		Comentarios:           r.FormValue("comentarios"),
	}

	// Save to Firestore
	fsClient, err := firestore.NewClient(ctx, "app-entradas-salidas-merc")
	if err != nil {
		http.Error(w, "Firestore error", http.StatusInternalServerError)
		return
	}
	defer fsClient.Close()

	_, _, err = fsClient.Collection("entradas").Add(ctx, entrada)
	if err != nil {
		http.Error(w, "Error saving to Firestore", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Entrada submitted successfully."}`))
}
