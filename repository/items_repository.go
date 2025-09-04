package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type ItemRepository struct{}

func NewItemRepository() *ItemRepository { return &ItemRepository{} }

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	var rows []models.Item
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *ItemRepository) GetByID(id uint) (*models.Item, error) {
	var row models.Item
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ItemRepository) GetByCompany(companyID uint) ([]models.Item, error) {
	var rows []models.Item
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}

func (r *ItemRepository) Create(m *models.Item) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *ItemRepository) Update(m *models.Item) error {
	return config.DB.Save(m).Error
}

func (r *ItemRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Item{}, id).Error
}
