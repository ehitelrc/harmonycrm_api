package repository

import (
	"harmony_api/models"

	"gorm.io/gorm"
)

type CustomListRepository struct {
	DB *gorm.DB
}

func NewCustomListRepository(db *gorm.DB) *CustomListRepository {
	return &CustomListRepository{DB: db}
}

func (r *CustomListRepository) GetDefinitions(entity string) ([]models.CustomListDefinition, error) {
	var defs []models.CustomListDefinition
	err := r.DB.Where("entity_name = ?", entity).Find(&defs).Error
	return defs, err
}

func (r *CustomListRepository) GetValues(listID uint) ([]models.CustomListValue, error) {
	var vals []models.CustomListValue
	err := r.DB.Where("list_id = ?", listID).Find(&vals).Error
	return vals, err
}

func (r *CustomListRepository) GetSelectedValue(entity string, entityID uint, listID uint) (*uint, error) {
	var rel models.CustomListEntityValue

	err := r.DB.Debug().
		Where("entity_name = ? AND entity_id = ? AND list_id = ?", entity, entityID, listID).
		First(&rel).Error

	if err != nil {
		return nil, nil // no selección previa
	}

	return &rel.ListValue, nil
}

func (r *CustomListRepository) SaveSelection(entity string, entityID uint, valueID uint, listID uint) error {
	// eliminar selecciones previas para ese entity+lista
	err := r.DB.Where("list_value = ?", valueID).
		Delete(&models.CustomListEntityValue{}).Error
	if err != nil {
		return err
	}

	// insertar nueva selección
	newSel := models.CustomListEntityValue{
		EntityName: entity,
		EntityID:   entityID,
		ListValue:  valueID,
		ListID:     listID,
	}

	return r.DB.Create(&newSel).Error
}
