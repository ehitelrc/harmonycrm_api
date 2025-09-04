package controllers

import (
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignController struct {
	repo *repository.CampaignRepository
}

func NewCampaignController() *CampaignController {
	return &CampaignController{repo: repository.NewCampaignRepository()}
}

// GET /campaigns/company/:company_id
func (cc *CampaignController) GetByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Parámetro company_id inválido", nil, err)
		return
	}
	rows, err := cc.repo.GetByCompany(uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener campañas", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Campañas obtenidas correctamente", rows, nil)
}

// GET /campaigns/:id
func (cc *CampaignController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := cc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Campaña no encontrada", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Campaña encontrada", row, nil)
}

// POST /campaigns
func (cc *CampaignController) Create(c *gin.Context) {
	var body models.Campaign
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear campaña", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Campaña creada correctamente", body, nil)
}

// PUT /campaigns  (recibe el objeto completo con id)
func (cc *CampaignController) Update(c *gin.Context) {
	var body models.Campaign
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar campaña", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Campaña actualizada correctamente", body, nil)
}

// DELETE /campaigns/:id
func (cc *CampaignController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := cc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar campaña", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Campaña eliminada correctamente", nil, nil)
}
