package repository

import (
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
)

type FunnelRepository struct{}

func NewFunnelRepository() *FunnelRepository { return &FunnelRepository{} }

// GetAll: lista funnels con stages ordenados por position
func (r *FunnelRepository) GetAll() ([]models.Funnel, error) {
	var list []models.Funnel

	db := config.DB

	err := db.Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *FunnelRepository) GetByID(id uint) (*models.Funnel, error) {
	var f models.Funnel
	db := config.DB
	err := db.First(&f, id).Error

	if err != nil {
		return nil, err
	}

	return &f, nil
}

// Create: crea funnel y (opcional) sus stages (si vienen en f.Stages)
func (r *FunnelRepository) Create(f *models.Funnel) error {
	db := config.DB
	return db.Transaction(func(tx *gorm.DB) error {
		// se crean stages si están en f.Stages (gorm los inserta si FKs se setean)
		if err := tx.Create(f).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *FunnelRepository) Update(f *models.Funnel) error {
	return config.DB.Save(f).Error
}

func (r *FunnelRepository) Delete(id uint) error {
	db := config.DB
	return db.Transaction(func(tx *gorm.DB) error {
		// si FK no tiene cascade, borrar stages primero
		if err := tx.Where("funnel_id = ?", id).Delete(&models.FunnelStage{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Funnel{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// Stages
func (r *FunnelRepository) GetStages(funnelID uint) ([]models.FunnelStage, error) {
	var stages []models.FunnelStage
	db := config.DB
	err := db.Where("funnel_id = ?", funnelID).Find(&stages).Error
	if err != nil {
		return nil, err
	}
	return stages, nil
}

func (r *FunnelRepository) GetStageByID(funnelID, stageID uint) (*models.FunnelStage, error) {
	var stage models.FunnelStage
	db := config.DB
	err := db.Where("funnel_id = ? AND id = ?", funnelID, stageID).First(&stage).Error
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *FunnelRepository) CreateStage(stage *models.FunnelStage) error {
	db := config.DB
	return db.Create(stage).Error
}

func (r *FunnelRepository) UpdateStage(stage *models.FunnelStage) error {
	db := config.DB
	return db.Save(stage).Error
}

func (r *FunnelRepository) DeleteStage(id uint) error {
	db := config.DB
	return db.Delete(&models.FunnelStage{}, id).Error
}
