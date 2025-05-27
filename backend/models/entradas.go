package models

import (
	"time"
)

type Entradas struct {
	TipoDelivery          string    `json:"tipo_delivery"`
	BodegaRecepcion       string    `json:"bodega_recepcion"`
	ProveedorRecepcion    string    `json:"proveedor_recepcion"`
	Cliente               string    `json:"cliente"`
	NumeroRemisionFactura string    `json:"numero_remision_factura"`
	PersonaRecepcion      string    `json:"persona_recepcion"`
	FechaRecepcion        time.Time `json:"fecha_recepcion"`
	EvidenciaRecepcion    string    `json:"evidencia_recepcion"` // GCS URL or object path
	Cantidad              int64     `json:"cantidad"`
	Comentarios           string    `json:"comentarios"`
}

type EntradasData struct {
	BodegaRecepcion    string    `firestore:"BodegaRecepcion"`
	Cantidad           int       `firestore:"Cantidad"`
	Comentarios        string    `firestore:"Comentarios"`
	EvidenciaRecepcion string    `firestore:"EvidenciaRecepcion"`
	FechaRecepcion     time.Time `firestore:"FechaRecepcion"`
	NumeroRemision     string    `firestore:"NumeroRemisionFactura"`
	PersonaRecepcion   string    `firestore:"PersonaRecepcion"`
	ProveedorRecepcion string    `firestore:"ProveedorRecepcion"`
	Cliente            string    `firestore:"Cliente"`
	TipoDelivery       string    `firestore:"TipoDelivery"`
	ASN                string    `firestore:"ASN"`
	FechaAjusteASN     time.Time `firestore:"FechaAjusteASN"`
}
