package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgentController struct {
	repo *repository.AgentRepository
}

func NewAgentController() *AgentController {
	return &AgentController{repo: repository.NewAgentRepository()}
}

// GET /agents
func (ac *AgentController) GetAll(c *gin.Context) {
	rows, err := ac.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener agentes", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Agentes obtenidos correctamente", rows, nil)
}

// GET /agents/:user_id
func (ac *AgentController) GetByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}
	row, err := ac.repo.GetByUserID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Agente no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Agente encontrado", row, nil)
}

// POST /agents
func (ac *AgentController) Create(c *gin.Context) {
	var body models.Agent
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	// body.UserID debe venir en el JSON
	if err := ac.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear agente", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Agente creado correctamente", body, nil)
}

// DELETE /agents/:user_id
func (ac *AgentController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}
	if err := ac.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar agente", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Agente eliminado correctamente", nil, nil)
}
