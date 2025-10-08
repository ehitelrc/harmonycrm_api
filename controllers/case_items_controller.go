package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CaseItemsController struct {
	repo *repository.CaseItemsRepositories
}

func NewCaseItemController() *CaseItemsController {
	return &CaseItemsController{repo: repository.NewCaseItemsRepository()}
}

func (c *CaseItemsController) GetAllItemsByCaseID(ctx *gin.Context) {
	caseID, err := strconv.Atoi(ctx.Param("case_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de caso inválido", nil, err)
		return
	}

	items, err := c.repo.GetAllItemsByCaseID(uint(caseID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener los ítems del caso", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Lista de ítems del caso", items, nil)
}

func (c *CaseItemsController) GetItemByCaseItemID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	item, err := c.repo.GetItemByCaseItemID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Ítem no encontrado", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Ítem encontrado", item, nil)
}

func (c *CaseItemsController) CreateCaseItem(ctx *gin.Context) {
	var item models.CaseItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.CreateCaseItem(&item); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear ítem del caso", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Ítem del caso creado exitosamente", item, nil)
}

func (c *CaseItemsController) UpdateCaseItem(ctx *gin.Context) {
	var item models.CaseItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	if err := c.repo.UpdateCaseItem(&item); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar ítem del caso", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Ítem del caso actualizado exitosamente", item, nil)
}

func (c *CaseItemsController) DeleteCaseItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.DeleteCaseItem(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar ítem del caso", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Ítem del caso eliminado exitosamente", nil, nil)
}
