package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type TagRepository struct{}

func NewTagRepository() *TagRepository {
	return &TagRepository{}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	return config.DB.Create(tag).Error
}

func (r *TagRepository) Update(tag *models.Tag) error {
	return config.DB.Save(tag).Error
}

func (r *TagRepository) Delete(tagID uint) error {
	return config.DB.Delete(&models.Tag{}, tagID).Error
}

func (r *TagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := config.DB.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := config.DB.Find(&tags).Error
	return tags, err
}

func (r *TagRepository) AssignToCase(caseID uint, tagID uint) error {
	// Usamos SQL nativo para insertar en case_tags
	return config.DB.Exec("INSERT INTO case_tags (case_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING", caseID, tagID).Error
}

func (r *TagRepository) RemoveFromCase(caseID uint, tagID uint) error {
	return config.DB.Exec("DELETE FROM case_tags WHERE case_id = ? AND tag_id = ?", caseID, tagID).Error
}

func (r *TagRepository) GetTagsByCase(caseID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := config.DB.Joins("JOIN case_tags ct ON ct.tag_id = tags.id").
		Where("ct.case_id = ?", caseID).
		Find(&tags).Error
	return tags, err
}
