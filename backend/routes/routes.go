package routes

import (
	"database/sql"

	"github.com/clopezbyte/app-entradas-salidas/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/locations", handlers.GetBodegaLocation(db)).Methods("GET")
	r.HandleFunc("/export", handlers.ExportAndEmail).Methods("POST")
	r.HandleFunc("/bodega-history", handlers.GetBodegaData(db)).Methods("GET")
	r.HandleFunc("/creds", handlers.GetCreds(db)).Methods("POST")
	r.HandleFunc("/entradas", handlers.HandleEntradasSubmit).Methods("POST")
	return r
}
