package routes

import (
	"database/sql"

	"github.com/clopezbyte/app-entradas-salidas/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/export", handlers.ExportAndEmail).Methods("POST")
	r.HandleFunc("/entradas", handlers.HandleEntradasSubmit).Methods("POST")
	return r
}
