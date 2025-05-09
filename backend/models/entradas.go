package models

type Entradas struct {
	TipoDelivery          string `json:"tipo_delivery"`
	BodegaRecepcion       string `json:"bodega_recepcion"`
	ProveedorRecepcion    string `json:"proveedor_recepcion"`
	NumeroRemisionFactura int64  `json:"numero_remision_factura"`
	PersonaRecepcion      string `json:"persona_recepcion"`
	FechaRecepcion        string `json:"fecha_recepcion"`
	EvidenciaRecepcion    string `json:"evidencia_recepcion"` // GCS URL or object path
	Cantidad              int64  `json:"cantidad"`
	Comentarios           string `json:"comentarios"`
}
