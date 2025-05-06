package routes

import (
	"your-app/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/entrada", handlers.RegisterEntrada).Methods("POST")
	r.HandleFunc("/salida", handlers.RegisterSalida).Methods("POST")
	return r
}
