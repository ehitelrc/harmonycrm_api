package models

type WhatsAppMessageControl struct {
	ID          int64  `gorm:"primary_key;autoIncrement"`
	WSMessageID string `gorm:"column:ws_message__id"`
}

func (WhatsAppMessageControl) TableName() string {
	return "whatsapp_message_control"
}
