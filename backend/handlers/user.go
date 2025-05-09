package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func GetCreds(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse JSON request body for username and password
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Println("Invalid request body:", err)
			return
		}

		if creds.Username == "" || creds.Password == "" {
			http.Error(w, "Missing username or password", http.StatusBadRequest)
			log.Println("Missing username or password")
			return
		}

		// Query the database for the hashed password
		query := "SELECT password FROM users WHERE username = ?"
		var hashedPassword string
		err := db.QueryRow(query, creds.Username).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				log.Println("Authentication failed: invalid username or password")
			} else {
				http.Error(w, "Error querying database", http.StatusInternalServerError)
				log.Printf("Error querying database: %v", err)
			}
			return
		}

		// Compare the provided password with the hashed password
		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			log.Println("Authentication failed: invalid username or password")
			return
		}

		// Authentication successful
		log.Printf("Successfully authenticated user: %v", creds.Username)

		// Return a safe response (e.g., username or user role)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Username string `json:"username"`
		}{Username: creds.Username})
	}
}
