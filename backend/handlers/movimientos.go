package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/models"
)

func GetBodegaData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodegaData []models.HistorialBodega

		rows, err := db.Query("select date_part('MONTH', movement_date) as mes, bodega, sum(quantity) as quantity from warehouse_movements group by mes, bodega order by mes asc;")
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			log.Printf("Error querying database: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var datos models.HistorialBodega
			if err := rows.Scan(&datos.Mes, &datos.Bodega, &datos.Quantity); err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				log.Printf("Error scanning row: %v", err)
			}
			bodegaData = append(bodegaData, datos)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows", http.StatusInternalServerError)
			log.Printf("Error iterating rows: %v", err)
			return
		}
		log.Printf("Successfully queried bodega history")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bodegaData)

	}

}
