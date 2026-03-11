package models

import "time"

type VwUnreconciledErpPayments struct {
	IdRecord        int64     `json:"id_record" gorm:"column:id_record"`
	ErpID           int64     `json:"erp_id" gorm:"column:erp_id"`
	ErpStatus       string    `json:"erp_status" gorm:"column:erp_status"`
	HarmonyState    string    `json:"harmony_state" gorm:"column:harmony_state"`
	BankName        string    `json:"bank_name" gorm:"column:bank_name"`
	ReferenceNumber string    `json:"reference_number" gorm:"column:reference_number"`
	PaymentDate     string    `json:"payment_date" gorm:"column:payment_date"`
	PaymentTime     string    `json:"payment_time" gorm:"column:payment_time"`
	ErpAmount       float64   `json:"erp_amount" gorm:"column:erp_amount"`
	ClientDocument  string    `json:"client_document" gorm:"column:client_document"`
	ClientName      string    `json:"client_name" gorm:"column:client_name"`
	ContractID      int64     `json:"contract_id" gorm:"column:contract_id"`
	ReceiptType     string    `json:"receipt_type" gorm:"column:receipt_type"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
}

func (VwUnreconciledErpPayments) TableName() string {
	return "vw_unreconciled_erp_payments"
}
