package routes

import (
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/entrada", handlers.RegisterEntrada).Methods("POST")
	r.HandleFunc("/salida", handlers.RegisterSalida).Methods("POST")
	http.Handle("/", r)
	return r
}
