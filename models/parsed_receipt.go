package models

type ParsedReceipt struct {
	BankName           string  `json:"bank_name"`
	TransactionType    string  `json:"transaction_type"`
	ReferenceNumber    string  `json:"reference_number"`
	Amount             float64 `json:"amount"`
	AmountSent         float64 `json:"amount_sent"`
	Date               string  `json:"date"`
	Time               string  `json:"time"`
	SenderName         string  `json:"sender_name"`
	SenderPhone        string  `json:"sender_phone"`
	ReceiverName       string  `json:"receiver_name"`
	ReceiverPhone      string  `json:"receiver_phone"`
	OriginAccount      string  `json:"origin_account"`
	DestinationAccount string  `json:"destination_account"`
	Description        string  `json:"description"`
}
