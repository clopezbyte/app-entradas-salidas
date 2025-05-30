package routes

import (
	"github.com/clopezbyte/app-entradas-salidas/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/entradas", handlers.HandleEntradasSubmit).Methods("POST")
	r.HandleFunc("/entradas-data", handlers.HandleProvideEntradasData).Methods("POST")
	r.HandleFunc("/salidas", handlers.HandleSalidasSubmit).Methods("POST")
	r.HandleFunc("/salidas-data", handlers.HandleProvideSalidasData).Methods("POST")
	r.HandleFunc("/query-numrem", handlers.QueryNumRem).Methods("POST")
	r.HandleFunc("/update-asn", handlers.HandleASNSubmit).Methods("POST")
	return r
}
