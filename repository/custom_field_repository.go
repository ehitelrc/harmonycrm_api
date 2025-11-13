package repository

import (
	"encoding/json"
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
	"strings"
	"time"
)

type CustomFieldRepository struct{}

func NewCustomFieldRepository() *CustomFieldRepository {
	return &CustomFieldRepository{}
}

// GetFields retorna todos los campos definidos + valores (si entityID != nil)
func (r *CustomFieldRepository) GetFields(entityName string, entityID uint) ([]map[string]interface{}, error) {
	var defs []models.CustomFieldDefinition
	if err := config.DB.
		Where("entity_name = ? AND is_active = TRUE", entityName).
		Order("sort_order").
		Find(&defs).Error; err != nil {
		return nil, err
	}

	var values []models.CustomFieldValue

	if err := config.DB.
		Where("entity_name = ? AND entity_id = ?", entityName, entityID).
		Find(&values).Error; err != nil {
		return nil, err
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
					field["value"] = v.ValueInteger
				case "decimal":
					field["value"] = v.ValueDecimal
				case "boolean":
					field["value"] = v.ValueBoolean
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

// SaveEntityCustomFieldValue guarda el valor de un campo personalizado para una entidad
func (r *CustomFieldRepository) SaveEntityCustomFieldValue(data []models.CustomFieldValuePayload) error {

	if len(data) == 0 {
		return nil // Nada que guardar
	}

	entityName := data[0].EntityName

	// 1️⃣ Cargar definiciones (tipos) de los campos para esta entidad
	var definitions []models.CustomFieldDefinition
	if err := config.DB.
		Where("entity_name = ?", entityName).
		Find(&definitions).Error; err != nil {
		return err
	}

	// 2️⃣ Crear un mapa rápido: field_key → definición
	defMap := make(map[string]models.CustomFieldDefinition)
	for _, d := range definitions {
		defMap[d.FieldKey] = d
	}

	// 3️⃣ Empezar transacción
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, item := range data {

		def, exists := defMap[item.FieldKey]
		if !exists {
			tx.Rollback()
			return fmt.Errorf("field_key '%s' no está definido para '%s'", item.FieldKey, entityName)
		}

		// 4️⃣ Convertir dinámicamente el `FieldValue` según FieldType
		cfValue := models.CustomFieldValue{
			FieldID:    def.ID,
			EntityName: item.EntityName,
			EntityID:   item.EntityID,
		}

		// 5️⃣ Asignar el valor correcto según el tipo
		switch def.FieldType {

		case "text":
			if item.FieldValue != nil {
				v := fmt.Sprintf("%v", item.FieldValue)
				cfValue.ValueText = &v
			}

		case "integer":
			if item.FieldValue != nil {
				switch v := item.FieldValue.(type) {
				case float64:
					i := int(v)
					cfValue.ValueInteger = &i
				case int:
					cfValue.ValueInteger = &v
				case json.Number:
					i, _ := v.Int64()
					ii := int(i)
					cfValue.ValueInteger = &ii
				default:
					tx.Rollback()
					return fmt.Errorf("integer inválido en field_key '%s'", item.FieldKey)
				}
			}

		case "decimal":
			if item.FieldValue != nil {
				switch v := item.FieldValue.(type) {
				case float64:
					cfValue.ValueDecimal = &v
				case json.Number:
					f, _ := v.Float64()
					cfValue.ValueDecimal = &f
				default:
					tx.Rollback()
					return fmt.Errorf("decimal inválido en field_key '%s'", item.FieldKey)
				}
			}

		case "boolean":
			if item.FieldValue != nil {
				switch v := item.FieldValue.(type) {
				case bool:
					cfValue.ValueBoolean = &v
				case string:
					b := strings.ToLower(v) == "true"
					cfValue.ValueBoolean = &b
				default:
					tx.Rollback()
					return fmt.Errorf("boolean inválido en field_key '%s'", item.FieldKey)
				}
			}

		case "date":
			if item.FieldValue != nil {
				s := fmt.Sprintf("%v", item.FieldValue)
				parsed, err := time.Parse("2006-01-02", s)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("fecha inválida en field_key '%s': %v", item.FieldKey, err)
				}
				cfValue.ValueDate = &parsed
			}

		default:
			tx.Rollback()
			return fmt.Errorf("tipo de campo no soportado: %s", def.FieldType)
		}

		// 6️⃣ Borrar anterior registro del mismo campo & entidad
		if err := tx.
			Where("field_id = ? AND entity_id = ?", def.ID, item.EntityID).
			Delete(&models.CustomFieldValue{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 7️⃣ Insert nuevo valor
		if err := tx.Create(&cfValue).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 8️⃣ Confirmar
	return tx.Commit().Error
}
