package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository { return &ClientRepository{} }

func (r *ClientRepository) GetAll() ([]models.Client, error) {
	var rows []models.Client
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *ClientRepository) GetByID(id uint) (*models.Client, error) {
	var row models.Client
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClientRepository) Create(m *models.Client) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *ClientRepository) Update(m *models.Client) error {
	return config.DB.Save(m).Error
}

func (r *ClientRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Client{}, id).Error
}
