package dto

type ProcessReceiptRequest struct {
	Nombre           string  `json:"nombre"`
	Monto            float64 `json:"monto"`
	Fecha            string  `json:"fecha"`
	NumeroReferencia string  `json:"numeroReferencia"`
	Observaciones    string  `json:"observaciones"`
	Mensaje          string  `json:"mensaje"`
}
