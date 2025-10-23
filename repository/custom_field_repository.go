package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type CustomFieldRepository struct{}

func NewCustomFieldRepository() *CustomFieldRepository {
	return &CustomFieldRepository{}
}

// GetFields retorna todos los campos definidos + valores (si entityID != nil)
func (r *CustomFieldRepository) GetFields(entityName string, entityID *uint) ([]map[string]interface{}, error) {
	var defs []models.CustomFieldDefinition
	if err := config.DB.
		Where("entity_name = ? AND is_active = TRUE", entityName).
		Order("sort_order").
		Find(&defs).Error; err != nil {
		return nil, err
	}

	var values []models.CustomFieldValue
	if entityID != nil {
		if err := config.DB.
			Where("entity_name = ? AND entity_id = ?", entityName, *entityID).
			Find(&values).Error; err != nil {
			return nil, err
		}
	}

	// Crear el mix
	result := make([]map[string]interface{}, 0, len(defs))
	for _, def := range defs {
		field := map[string]interface{}{
			"field_id":    def.ID,
			"field_key":   def.FieldKey,
			"label":       def.Label,
			"field_type":  def.FieldType,
			"is_required": def.IsRequired,
			"value":       nil,
		}

		// Buscar valor si aplica
		for _, v := range values {
			if v.FieldID == def.ID {
				switch def.FieldType {
				case "text":
					field["value"] = v.ValueText
				case "integer":
					field["value"] = v.ValueInt
				case "decimal":
					field["value"] = v.ValueDec
				case "boolean":
					field["value"] = v.ValueBool
				case "date":
					field["value"] = v.ValueDate
				}
				break
			}
		}

		result = append(result, field)
	}

	return result, nil
}
