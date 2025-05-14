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
	EvidenciaRecepcion     string    `json:"evidencia_recepcion"` // GCS URL or object path
	Comentarios            string    `json:"comentarios"`
}

type SalidasData struct {
	BodegaSalida           string    `json:"bodega_salida"`
	ProveedorSalida        string    `json:"proveedor_salida"`
	NumeroOrdenConsecutivo int64     `json:"numero_orden_consecutivo"`
	PersonaEntrega         string    `json:"persona_entrega"`
	FechaSalida            time.Time `json:"fecha_salida"`
	EvidenciaRecepcion     string    `json:"evidencia_recepcion"` // GCS URL or object path
	Comentarios            string    `json:"comentarios"`
}
