package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (r *CompanyRepository) GetAll() ([]models.Company, error) {
	var companies []models.Company
	err := config.DB.Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) GetByID(id uint) (*models.Company, error) {
	var company models.Company
	err := config.DB.First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) Create(company *models.Company) error {
	return config.DB.Create(company).Error
}

func (r *CompanyRepository) Update(company *models.Company) error {
	return config.DB.Save(company).Error
}

func (r *CompanyRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Company{}, id).Error
}

func (r *CompanyRepository) GetByUserID(userID uint) ([]models.CompanyUserView, error) {
	var companies []models.CompanyUserView
	err := config.DB.Where("user_id = ?", userID).Find(&companies).Error
	return companies, err
}
