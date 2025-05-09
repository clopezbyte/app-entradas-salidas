package models

type Salidas struct {
	BodegaSalida           string `json:"bodega_salida"`
	ProveedorSalida        string `json:"proveedor_salida"`
	NumeroOrdenConsecutivo int64  `json:"numero_orden_consecutivo"`
	PersonaEntrega         string `json:"persona_entrega"`
	FechaSalida            string `json:"fecha_salida"`
	EvidenciaRecepcion     string `json:"evidencia_recepcion"` // GCS URL or object path
	Comentarios            string `json:"comentarios"`
}
