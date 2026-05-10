package repository

import (
	"encoding/json"
	"harmony_api/config"
	"harmony_api/dto"
	"harmony_api/models"
	"time"
)

type ReceiptRepository struct{}

func NewReceiptRepository() *ReceiptRepository {
	return &ReceiptRepository{}
}

func (r *ReceiptRepository) SaveReceiptResult(result *dto.ReceiptExtractionResult, caseID uint, messageID *uint64, createdAt *time.Time) (*models.ReceiptResult, error) {

	warningsJSON, _ := json.Marshal(result.Warnings)

	record := models.ReceiptResult{
		CaseID:    caseID,
		MessageID: messageID,
		Status:    "new",

		BankName:        result.BankName,
		TransactionType: result.TransactionType,
		ReferenceNumber: result.ReferenceNumber,
		Date:            result.Date,
		Time:            result.Time,
		Amount:          result.Amount,
		AmountSent:      result.AmountSent,

		SenderName:    result.SenderName,
		SenderPhone:   result.SenderPhone,
		ReceiverName:  result.ReceiverName,
		ReceiverPhone: result.ReceiverPhone,

		OriginAccount:      result.OriginAccount,
		DestinationAccount: result.DestinationAccount,

		Description: result.Description,
		RawText:     result.RawText,
		Warnings:    string(warningsJSON),
	}

	if createdAt != nil {
		record.CreatedAt = *createdAt
	}

	if err := config.DB.Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

// Obtener todos los que están "new"
func (r *ReceiptRepository) GetNewReceipts() ([]models.ReceiptResult, error) {
	var receipts []models.ReceiptResult

	err := config.DB.
		Where("status = ?", "new").
		Order("created_at ASC").
		Find(&receipts).Error

	return receipts, err
}

// Cambiar estado a "read"
func (r *ReceiptRepository) MarkAsRead(id uint) error {
	return config.DB.Model(&models.ReceiptResult{}).
		Where("id = ?", id).
		Update("status", "read").Error
}

// Cambiar estado a "processed"
func (r *ReceiptRepository) MarkAsProcessed(id uint) error {
	return config.DB.Model(&models.ReceiptResult{}).
		Where("id = ?", id).
		Update("status", "processed").Error
}
