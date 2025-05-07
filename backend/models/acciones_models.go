package models

import (
	"time"
)

type Movimiento struct {
	ID        int       `json:"id"`
	Tipo      string    `json:"tipo"`
	Producto  string    `json:"producto"`
	Cantidad  int       `json:"cantidad"`
	Timestamp time.Time `json:"timestamp"`
}

// DB save/query function here
func SaveMovimiento(m Movimiento) error {
	// Insert into DB (PostgreSQL or Firestore logic)
	return nil
}
