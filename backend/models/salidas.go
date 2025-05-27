package models

import (
	"time"
)

type Salidas struct {
	BodegaSalida           string    `json:"bodega_salida"`
	ProveedorSalida        string    `json:"proveedor_salida"`
	Cliente                string    `json:"cliente"`
	NumeroOrdenConsecutivo string    `json:"numero_orden_consecutivo"`
	PersonaEntrega         string    `json:"persona_entrega"`
	PersonaRecoge          string    `json:"persona_recoge"`
	FirmaPersonaRecoge     string    `json:"firma_persona_recoge"` // GCS URL or object path
	FechaSalida            time.Time `json:"fecha_salida"`
	EvidenciaSalida        string    `json:"evidencia_salida"` // GCS URL or object path
	Comentarios            string    `json:"comentarios"`
}

type SalidasData struct {
	BodegaSalida           string    `firestore:"BodegaSalida"`
	ProveedorSalida        string    `firestore:"ProveedorSalida"`
	Cliente                string    `firestore:"Cliente"`
	NumeroOrdenConsecutivo string    `firestore:"NumeroOrdenConsecutivo"`
	PersonaEntrega         string    `firestore:"PersonaEntrega"`
	PersonaRecoge          string    `firestore:"PersonaRecoge"`
	FirmaPersonaRecoge     string    `firestore:"FirmaPersonaRecoge"` // GCS URL or object path
	FechaSalida            time.Time `firestore:"FechaSalida"`
	EvidenciaSalida        string    `firestore:"EvidenciaSalida"` // GCS URL or object path
	Comentarios            string    `firestore:"Comentarios"`
}
