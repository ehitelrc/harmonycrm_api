package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type UnreconciledPaymentsRepository struct {
}

func NewUnreconciledPaymentsRepository() *UnreconciledPaymentsRepository {
	return &UnreconciledPaymentsRepository{}
}

func (r *UnreconciledPaymentsRepository) GetUnreconciledPayments(startDate, endDate string) ([]models.VwUnreconciledErpPayments, error) {
	var payments []models.VwUnreconciledErpPayments

	query := config.DB.Model(&models.VwUnreconciledErpPayments{})

	if startDate != "" && endDate != "" {
		query = query.Where("created_at >= ? AND created_at <= ?", startDate+" 00:00:00", endDate+" 23:59:59.999999")
	} else if startDate != "" {
		query = query.Where("created_at >= ?", startDate+" 00:00:00")
	}

	err := query.Find(&payments).Error
	return payments, err
}
