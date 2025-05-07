package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clopezbyte/app-entradas-salidas/models"
)

func RegisterEntrada(w http.ResponseWriter, r *http.Request) {
	var entrada models.Movimiento
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entrada.Tipo = "entrada"
	// save to DB (pseudo code)
	if err := models.SaveMovimiento(entrada); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entrada)
}

func RegisterSalida(w http.ResponseWriter, r *http.Request) {
	var salida models.Movimiento
	if err := json.NewDecoder(r.Body).Decode(&salida); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	salida.Tipo = "salida"
	// save to DB (pseudo code)
	if err := models.SaveMovimiento(salida); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(salida)
}
