package models

type VwPaymentValidations struct {
	CaseID                 int64   `json:"case_id" gorm:"column:case_id"`
	SenderID               string  `json:"sender_id" gorm:"column:sender_id"`
	BankName               string  `json:"bank_name" gorm:"column:bank_name"`
	TransactionType        string  `json:"transaction_type" gorm:"column:transaction_type"`
	ReceiptReferenceNumber string  `json:"receipt_reference_number" gorm:"column:receipt_reference_number"`
	ReceiptDate            string  `json:"receipt_date" gorm:"column:receipt_date"`
	ReceiptTime            string  `json:"receipt_time" gorm:"column:receipt_time"`
	ReceiptAmount          float64 `json:"receipt_amount" gorm:"column:receipt_amount"`
	AmountSent             float64 `json:"amount_sent" gorm:"column:amount_sent"`
	SenderName             string  `json:"sender_name" gorm:"column:sender_name"`
	RawText                string  `json:"raw_text" gorm:"column:raw_text"`
	ErpStatus              string  `json:"erp_status" gorm:"column:erp_status"`
	HarmonyState           string  `json:"harmony_state" gorm:"column:harmony_state"`
	ErpReferenceNumber     string  `json:"erp_reference_number" gorm:"column:erp_reference_number"`
	PaymentDate            string  `json:"payment_date" gorm:"column:payment_date"`
	PaymentTime            string  `json:"payment_time" gorm:"column:payment_time"`
	ErpAmount              float64 `json:"erp_amount" gorm:"column:erp_amount"`
	ClientDocument         string  `json:"client_document" gorm:"column:client_document"`
	ErpID                  int64   `json:"erp_id" gorm:"column:erp_id"`
	ClientName             string  `json:"client_name" gorm:"column:client_name"`
	ContractID             int64   `json:"contract_id" gorm:"column:contract_id"`
}

func (VwPaymentValidations) TableName() string {
	return "vw_payment_validations"
}
