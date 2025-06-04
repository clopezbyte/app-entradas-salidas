package models

import (
	"time"
)

type ASN struct {
	ID             string    `json:"id"`
	ASN            string    `json:"asn"`
	FechaAjusteASN time.Time `json:"fecha_ajuste_asn"`
}
