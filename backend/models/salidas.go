package models

import (
	"time"
)

type Salidas struct {
	BodegaSalida           string    `json:"bodega_salida"`
	ProveedorSalida        string    `json:"proveedor_salida"`
	NumeroOrdenConsecutivo int64     `json:"numero_orden_consecutivo"`
	PersonaEntrega         string    `json:"persona_entrega"`
	FechaSalida            time.Time `json:"fecha_salida"`
	EvidenciaSalida        string    `json:"evidencia_salida"` // GCS URL or object path
	Comentarios            string    `json:"comentarios"`
}

type SalidasData struct {
	BodegaSalida           string    `firestore:"BodegaSalida"`
	ProveedorSalida        string    `firestore:"ProveedorSalida"`
	NumeroOrdenConsecutivo int64     `firestore:"NumeroOrdenConsecutivo"`
	PersonaEntrega         string    `firestore:"PersonaEntrega"`
	FechaSalida            time.Time `firestore:"FechaSalida"`
	EvidenciaSalida        string    `firestore:"EvidenciaSalida"` // GCS URL or object path
	Comentarios            string    `firestore:"Comentarios"`
}
