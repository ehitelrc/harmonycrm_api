package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type PaymentValidationRepository struct {
}

func NewPaymentValidationRepository() *PaymentValidationRepository {
	return &PaymentValidationRepository{}
}

// GetPaymentValidations returns records with optional filters
func (r *PaymentValidationRepository) GetPaymentValidations(startDate, endDate string, caseId int, contractId int) ([]models.VwPaymentValidations, error) {
	var validations []models.VwPaymentValidations

	query := config.DB.Model(&models.VwPaymentValidations{})

	if startDate != "" && endDate != "" {
		query = query.Where("ocr_date >= ? AND ocr_date <= ?", startDate+" 00:00:00", endDate+" 23:59:59.999999")
	} else if startDate != "" {
		query = query.Where("ocr_date >= ?", startDate+" 00:00:00")
	}

	if caseId > 0 {
		query = query.Where("case_id = ?", caseId)
	}

	if contractId > 0 {
		query = query.Where("contract_id = ?", contractId)
	}

	err := query.Debug().Find(&validations).Error
	return validations, err
}

func (r *PaymentValidationRepository) GetReceiptBase64(erpId int64) (string, error) {
	var result struct {
		ReceiptBase64 string `gorm:"column:receipt_base64"`
	}

	err := config.DB.Table("erp_payment_confirmation").
		Select("receipt_base64").
		Where("id = ?", erpId).
		Take(&result).Error

	if err != nil {
		return "", err
	}
	return result.ReceiptBase64, nil
}
