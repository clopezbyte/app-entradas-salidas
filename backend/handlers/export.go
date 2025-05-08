package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/utils"
)

type ExportRequest struct {
	Email   string            `json:"email"`
	Filters map[string]string `json:"filters"` //filtering
}

func ExportAndEmail(w http.ResponseWriter, r *http.Request) {
	var req ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Use req.Filters to query your DB, for now, mock data
	headers := []string{"ID", "Name", "Status"}
	rows := [][]string{
		{"1", "Item A", "In"},
		{"2", "Item B", "Out"},
	}

	csvData, err := utils.GenerateCSV(headers, rows)
	if err != nil {
		http.Error(w, "Could not generate CSV", http.StatusInternalServerError)
		return
	}

	err = utils.SendEmailWithCSV(req.Email, "Your Exported Data", "See attached CSV.", csvData)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
