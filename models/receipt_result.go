package models

import "time"

type ReceiptResult struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	CaseID uint   `json:"case_id" gorm:"index"`
	Status string `json:"status" gorm:"type:varchar(30);default:'new'"`

	BankName        string `json:"bank_name" gorm:"type:varchar(100)"`
	TransactionType string `json:"transaction_type" gorm:"type:varchar(100)"`
	ReferenceNumber string `json:"reference_number" gorm:"type:varchar(100)"`
	Date            string `json:"date" gorm:"type:varchar(20)"` // yyyy/MM/dd normalizado
	Time            string `json:"time" gorm:"type:varchar(10)"` // HH:mm normalizado

	Amount     float64 `json:"amount"`
	AmountSent float64 `json:"amount_sent"`

	SenderName    string `json:"sender_name" gorm:"type:varchar(200)"`
	SenderPhone   string `json:"sender_phone" gorm:"type:varchar(50)"`
	ReceiverName  string `json:"receiver_name" gorm:"type:varchar(200)"`
	ReceiverPhone string `json:"receiver_phone" gorm:"type:varchar(50)"`

	OriginAccount      string `json:"origin_account" gorm:"type:varchar(100)"`
	DestinationAccount string `json:"destination_account" gorm:"type:varchar(100)"`

	Description string `json:"description" gorm:"type:text"`
	RawText     string `json:"raw_text" gorm:"type:text"`

	Warnings string `json:"warnings" gorm:"type:text"` // se guardan como JSON serializado

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ReceiptResult) TableName() string {
	return "receipt_results"
}
