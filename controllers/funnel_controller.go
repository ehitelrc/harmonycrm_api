package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FunnelController struct {
	repo *repository.FunnelRepository
}

func NewFunnelController() *FunnelController {
	return &FunnelController{
		repo: repository.NewFunnelRepository(),
	}
}

// GET /funnels
func (ctl *FunnelController) GetAll(c *gin.Context) {
	list, err := ctl.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al listar funnels", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Lista de funnels", list, nil)
}

// GET /funnels/:id
func (ctl *FunnelController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	f, err := ctl.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener funnel", nil, err)
		return
	}
	if f == nil {
		utils.Respond(c, http.StatusNotFound, false, "Funnel no encontrado", nil, nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Funnel encontrado", f, nil)
}

// POST /funnels  (crear funnel + stages opcionales)
func (ctl *FunnelController) Create(c *gin.Context) {
	var in models.Funnel
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	// si vienen stages, GORM los crea ligado por foreignKey
	if err := ctl.repo.Create(&in); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear funnel", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Funnel creado correctamente", in, nil)
}

func (ctl *FunnelController) Update(c *gin.Context) {
	var in models.Funnel
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	// si vienen stages, GORM los crea ligado por foreignKey
	if err := ctl.repo.Update(&in); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar funnel", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Funnel actualizado correctamente", in, nil)
}

// DELETE /funnels/:id
func (ctl *FunnelController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := ctl.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar funnel", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Funnel eliminado correctamente", nil, nil)
}

// Stages
func (ctl *FunnelController) GetStages(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	stages, err := ctl.repo.GetStages(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener etapas", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Etapas del funnel", stages, nil)
}

func (ctl *FunnelController) GetStageByID(c *gin.Context) {
	funnelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de funnel inválido", nil, err)
		return
	}
	stageID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de etapa inválido", nil, err)
		return
	}
	stage, err := ctl.repo.GetStageByID(uint(funnelID), uint(stageID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener etapa", nil, err)
		return
	}
	if stage == nil {
		utils.Respond(c, http.StatusNotFound, false, "Etapa no encontrada", nil, nil)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Etapa encontrada", stage, nil)
}

func (ctl *FunnelController) CreateStage(c *gin.Context) {
	var in models.FunnelStage
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := ctl.repo.CreateStage(&in); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear etapa", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Etapa creada correctamente", in, nil)
}

func (ctl *FunnelController) UpdateStage(c *gin.Context) {
	var in models.FunnelStage
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := ctl.repo.UpdateStage(&in); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar etapa", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Etapa actualizada correctamente", in, nil)
}

func (ctl *FunnelController) DeleteStage(c *gin.Context) {
	stageID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de etapa inválido", nil, err)
		return
	}
	if err := ctl.repo.DeleteStage(uint(stageID)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar etapa", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Etapa eliminada correctamente", nil, nil)
}
