package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository { return &UserRepository{} }

func (r *UserRepository) GetAll() ([]models.User, error) {
	var rows []models.User
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var row models.User
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) Create(m *models.User) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *UserRepository) Update(m *models.User) error {
	return config.DB.Save(m).Error
}

func (r *UserRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
