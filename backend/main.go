package main

import (
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/routes"
)

func main() {

	// Load dotenv only needed in development
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	router := routes.SetupRouter()

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
