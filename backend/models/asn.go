package models

import (
	"time"
)

type ASN struct {
	NumeroRemisionFactura int64     `json:"numero_remision_factura"`
	ASN                   string    `json:"asn"`
	FechaAjusteASN        time.Time `json:"fecha_ajuste_asn"`
}
