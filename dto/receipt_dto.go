package dto

type ReceiptExtractionResult struct {
	BankName           string   `json:"bank_name"`
	TransactionType    string   `json:"transaction_type"`
	ReferenceNumber    string   `json:"reference_number"`
	Date               string   `json:"date"`
	Time               string   `json:"time"`
	Amount             float64  `json:"amount"`      // Monto debitado
	AmountSent         float64  `json:"amount_sent"` // Monto enviado (si aplica)
	SenderName         string   `json:"sender_name"`
	SenderPhone        string   `json:"sender_phone"`
	ReceiverName       string   `json:"receiver_name"`
	ReceiverPhone      string   `json:"receiver_phone"`
	OriginAccount      string   `json:"origin_account"`
	DestinationAccount string   `json:"destination_account"`
	Description        string   `json:"description"`
	RawText            string   `json:"raw_text"` // Texto OCR completo por cualquier cosa
	Warnings           []string `json:"warnings"` // Mensajes tipo “no se encontró X”
}
