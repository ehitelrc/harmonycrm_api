package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	repo *repository.ItemRepository
}

func NewItemController() *ItemController {
	return &ItemController{repo: repository.NewItemRepository()}
}

// GET /items
func (ic *ItemController) GetAll(c *gin.Context) {
	rows, err := ic.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener ítems", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Ítems obtenidos correctamente", rows, nil)
}

// GET /items/:id
func (ic *ItemController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := ic.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Ítem no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Ítem encontrado", row, nil)
}

// GET /items/company/:company_id
func (ic *ItemController) GetByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Company ID inválido", nil, err)
		return
	}
	rows, err := ic.repo.GetByCompany(uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener ítems por compañía", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Ítems por compañía obtenidos correctamente", rows, nil)
}

// POST /items
func (ic *ItemController) Create(c *gin.Context) {
	var body models.Item
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := ic.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear ítem", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Ítem creado correctamente", body, nil)
}

// PUT /items (recibe objeto completo con id)
func (ic *ItemController) Update(c *gin.Context) {
	var body models.Item
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := ic.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar ítem", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Ítem actualizado correctamente", body, nil)
}

// DELETE /items/:id
func (ic *ItemController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := ic.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar ítem", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Ítem eliminado correctamente", nil, nil)
}
