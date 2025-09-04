package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type DepartmentRepository struct{}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{}
}

func (r *DepartmentRepository) GetAll() ([]models.Department, error) {
	var departments []models.Department
	err := config.DB.Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) GetByCompanyID(companyID uint) ([]models.Department, error) {
	var departments []models.Department
	err := config.DB.Where("company_id = ?", companyID).Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) GetByID(id uint) (*models.Department, error) {
	var dept models.Department
	err := config.DB.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepository) Create(dept *models.Department) error {
	return config.DB.Create(dept).Error
}

func (r *DepartmentRepository) Update(dept *models.Department) error {
	return config.DB.Save(dept).Error
}

func (r *DepartmentRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Department{}, id).Error
}
