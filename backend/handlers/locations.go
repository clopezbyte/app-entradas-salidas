package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/models"
)

func GetBodegaLocation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var locations []models.BodegaLocation

		rows, err := db.Query("SELECT DISTINCT bodega FROM warehouse_movements")
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			log.Printf("Error querying database: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var location models.BodegaLocation
			if err := rows.Scan(&location.Location); err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				log.Printf("Error scanning row: %v", err)
				return
			}
			locations = append(locations, location)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows", http.StatusInternalServerError)
			log.Printf("Error iterating rows: %v", err)
			return
		}

		log.Printf("Successfully queried bodega locations: %v", locations)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(locations)
	}
}
