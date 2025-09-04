package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := config.DB.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(role *models.Role) error {
	return config.DB.Create(role).Error
}

func (r *RoleRepository) Update(role *models.Role) error {
	return config.DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Role{}, id).Error
}
