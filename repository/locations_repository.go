package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type LocationRepositories struct {
}

func NewCountriesRepository() *LocationRepositories {
	return &LocationRepositories{}
}

func (r *LocationRepositories) GetAllCountries() ([]models.Country, error) {
	var countries []models.Country
	if err := config.DB.Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *LocationRepositories) GetCountryByID(id uint) (*models.Country, error) {
	var country models.Country
	if err := config.DB.First(&country, id).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *LocationRepositories) CreateCountry(country *models.Country) error {
	if err := config.DB.Create(country).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocationRepositories) UpdateCountry(country *models.Country) error {
	if err := config.DB.Save(country).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocationRepositories) DeleteCountry(id uint) error {
	if err := config.DB.Delete(&models.Country{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocationRepositories) GetProvincesByCountryID(countryID uint) ([]models.Province, error) {
	var provinces []models.Province
	if err := config.DB.Where("country_id = ?", countryID).Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

// Create province
func (r *LocationRepositories) CreateProvince(province *models.Province) error {
	if err := config.DB.Create(province).Error; err != nil {
		return err
	}
	return nil
}

// Update province
func (r *LocationRepositories) UpdateProvince(province *models.Province) error {
	if err := config.DB.Save(province).Error; err != nil {
		return err
	}
	return nil
}

// Delete province
func (r *LocationRepositories) DeleteProvince(id uint) error {
	if err := config.DB.Delete(&models.Province{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocationRepositories) GetCantonsByProvinceID(provinceID uint) ([]models.Canton, error) {
	var cantons []models.Canton
	if err := config.DB.Where("province_id = ?", provinceID).Find(&cantons).Error; err != nil {
		return nil, err
	}
	return cantons, nil
}

// Get canton by id
func (r *LocationRepositories) GetCantonByID(id uint) (*models.Canton, error) {
	var canton models.Canton
	if err := config.DB.First(&canton, id).Error; err != nil {
		return nil, err
	}
	return &canton, nil
}

// Create canton
func (r *LocationRepositories) CreateCanton(canton *models.Canton) error {
	if err := config.DB.Create(canton).Error; err != nil {
		return err
	}
	return nil
}

// Update canton
func (r *LocationRepositories) UpdateCanton(canton *models.Canton) error {
	if err := config.DB.Save(canton).Error; err != nil {
		return err
	}
	return nil
}

// Delete canton
func (r *LocationRepositories) DeleteCanton(id uint) error {
	if err := config.DB.Delete(&models.Canton{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocationRepositories) GetDistrictsByCantonID(cantonID uint) ([]models.District, error) {
	var districts []models.District
	if err := config.DB.Where("canton_id = ?", cantonID).Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

// Get district by id
func (r *LocationRepositories) GetDistrictByID(id uint) (*models.District, error) {
	var district models.District
	if err := config.DB.First(&district, id).Error; err != nil {
		return nil, err
	}
	return &district, nil
}

// Create district
func (r *LocationRepositories) CreateDistrict(district *models.District) error {
	if err := config.DB.Create(district).Error; err != nil {
		return err
	}
	return nil
}

// Update district
func (r *LocationRepositories) UpdateDistrict(district *models.District) error {
	if err := config.DB.Save(district).Error; err != nil {
		return err
	}
	return nil
}

// Delete district
func (r *LocationRepositories) DeleteDistrict(id uint) error {
	if err := config.DB.Delete(&models.District{}, id).Error; err != nil {
		return err
	}
	return nil
}
