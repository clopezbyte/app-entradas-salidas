package models

type DatosEntrada struct {
	BodegaEntrada string `json:"bodega"`
	MesEntrada    string `json:"mes_entrada"`
}

type DatosSalida struct {
	BodegaSalida string `json:"bodega"`
	MesSalida    string `json:"mes_salida"`
}

type BodegaLocation struct {
	Location string `json:"location"`
}

type HistorialBodega struct {
	Mes      string `json:"mes"`
	Bodega   string `json:"bodega"`
	Quantity int16  `json:"cantidad"`
}
