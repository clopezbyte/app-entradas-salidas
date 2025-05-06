package main

import (
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/backend/routes"
)

func main() {
	// Set up the router
	router := routes.SetupRouter()

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
