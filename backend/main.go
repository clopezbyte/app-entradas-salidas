package main

import (
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/routes"
)

func main() {

	router := routes.SetupRouter()

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
