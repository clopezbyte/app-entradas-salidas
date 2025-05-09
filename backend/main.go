package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/routes"
	_ "github.com/lib/pq" //PostgreSQL driver
)

func main() {

	// password := "admin"
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	fmt.Println("Error hashing password:", err)
	// 	return
	// }
	// fmt.Println("Hashed password:", string(hashedPassword))

	//DB conn (pg or firestore)
	db, err := sql.Open("postgres", "host=localhost port=5434 user=admin password=admin dbname=test sslmode=disable")
	if err != nil {
		log.Fatal("Failed connecting to database", err)
	}
	defer db.Close()

	router := routes.SetupRouter(db)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
