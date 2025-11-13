package controllers

import (
	"net/http"
	"strconv"

	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"

	"github.com/gin-gonic/gin"
)

type CustomListController struct {
	Repo *repository.CustomListRepository
}

func NewCustomListController(repo *repository.CustomListRepository) *CustomListController {
	return &CustomListController{Repo: repo}
}

// GET /api/custom-lists/:entity/:entityId
func (ctl *CustomListController) GetLists(c *gin.Context) {
	entity := c.Param("entity")
	entityIDStr := c.Param("entity_id")

	entityID, _ := strconv.ParseUint(entityIDStr, 10, 64)

	defs, err := ctl.Repo.GetDefinitions(entity)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error obteniendo definiciones", nil, err)
		return
	}

	var result []map[string]interface{}

	for _, d := range defs {
		values, _ := ctl.Repo.GetValues(d.ID)
		selected, _ := ctl.Repo.GetSelectedValue(entity, uint(entityID), d.ID)

		listDTO := map[string]interface{}{
			"list_id":           d.ID,
			"list_name":         d.ListName,
			"code_label":        d.CodeLabel,
			"description_label": d.DescriptionLabel,
			"list_label":        d.ListLabel,
			"selected_value":    selected,
			"values": func() []map[string]interface{} {
				arr := []map[string]interface{}{}
				for _, v := range values {
					arr = append(arr, map[string]interface{}{
						"id":          v.ID,
						"code":        v.CodeValue,
						"description": v.DescriptionValue,
					})
				}
				return arr
			}(),
		}

		result = append(result, listDTO)
	}

	utils.Respond(c, http.StatusOK, true, "Listas obtenidas", result, nil)
}

// GET /api/custom-lists/list
func (ctl *CustomListController) GetAllLists(c *gin.Context) {
	defs, err := ctl.Repo.GetAllLists()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error obteniendo definiciones", nil, err)
		return
	}

	var result []map[string]interface{}

	for _, d := range defs {
		values, _ := ctl.Repo.GetValues(d.ID)

		listDTO := map[string]interface{}{
			"list_id":           d.ID,
			"list_name":         d.ListName,
			"code_label":        d.CodeLabel,
			"description_label": d.DescriptionLabel,
			"list_label":        d.ListLabel,
			"values": func() []map[string]interface{} {
				arr := []map[string]interface{}{}
				for _, v := range values {
					arr = append(arr, map[string]interface{}{
						"id":          v.ID,
						"code":        v.CodeValue,
						"description": v.DescriptionValue,
					})
				}
				return arr
			}(),
		}

		result = append(result, listDTO)
	}

	utils.Respond(c, http.StatusOK, true, "Listas obtenidas", result, nil)
}

// POST /api/custom-lists/select
func (ctl *CustomListController) SaveSelection(c *gin.Context) {
	var payload struct {
		EntityName string `json:"entity_name"`
		EntityID   uint   `json:"entity_id"`
		ListValue  uint   `json:"list_value"`
		ListID     uint   `json:"list_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Payload inválido", nil, err)
		return
	}

	err := ctl.Repo.SaveSelection(payload.EntityName, payload.EntityID, payload.ListValue, payload.ListID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error guardando selección", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Selección guardada", nil, nil)
}

func (ctl *CustomListController) CreateListValue(c *gin.Context) {
	var payload struct {
		ListID           uint   `json:"list_id"`
		CodeValue        string `json:"code_value"`
		DescriptionValue string `json:"description_value"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Payload inválido", nil, err)
		return
	}

	newValue := models.CustomListValue{
		ListID:           payload.ListID,
		CodeValue:        payload.CodeValue,
		DescriptionValue: payload.DescriptionValue,
	}

	err := ctl.Repo.CreateListValue(&newValue)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error creando valor de lista", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Valor de lista creado", newValue, nil)
}

func (ctl *CustomListController) UpdateListValue(c *gin.Context) {

	var payload struct {
		ID               uint   `json:"id"`
		ListID           uint   `json:"list_id"`
		CodeValue        string `json:"code_value"`
		DescriptionValue string `json:"description_value"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Payload inválido", nil, err)
		return
	}

	updateValue := models.CustomListValue{
		ID:               payload.ID,
		ListID:           payload.ListID,
		CodeValue:        payload.CodeValue,
		DescriptionValue: payload.DescriptionValue,
	}

	err := ctl.Repo.UpdateListValue(&updateValue)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error actualizando valor de lista", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Valor de lista actualizado", updateValue, nil)
}

func (ctl *CustomListController) DeleteListValue(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	hasData, err := ctl.Repo.ListHasData(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error verificando datos asociados", nil, err)
		return
	}

	if hasData {
		utils.Respond(c, http.StatusBadRequest, false, "No se puede eliminar el valor de lista porque tiene datos asociados", nil, nil)
		return
	}

	err = ctl.Repo.DeleteListValue(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error eliminando valor de lista", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Valor de lista eliminado", nil, nil)

}
