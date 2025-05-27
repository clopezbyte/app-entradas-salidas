package models

import "time"

type EmailData struct {
	Email              string    `firestore:"email"`
	RepName            string    `firestore:"rep_name"`
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
}
