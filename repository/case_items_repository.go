package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type CaseItemsRepositories struct {
}

func NewCaseItemsRepository() *CaseItemsRepositories {
	return &CaseItemsRepositories{}
}

func (r *CaseItemsRepositories) GetAllItemsByCaseID(caseID uint) ([]models.VwCaseItemsDetail, error) {
	var items []models.VwCaseItemsDetail
	if err := config.DB.Where("case_id = ?", caseID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CaseItemsRepositories) GetItemByCaseItemID(id uint) (*models.CaseItem, error) {
	var item models.CaseItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CaseItemsRepositories) CreateCaseItem(item *models.CaseItem) error {
	if err := config.DB.Create(item).Error; err != nil {
		return err
	}
	return nil
}

func (r *CaseItemsRepositories) UpdateCaseItem(item *models.CaseItem) error {
	if err := config.DB.Save(item).Error; err != nil {
		return err
	}
	return nil
}

func (r *CaseItemsRepositories) DeleteCaseItem(id uint) error {
	if err := config.DB.Delete(&models.CaseItem{}, id).Error; err != nil {
		return err
	}
	return nil
}
