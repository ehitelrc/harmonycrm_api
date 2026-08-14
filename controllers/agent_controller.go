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

	if cached, found := utils.GetAgentsAllFromCache(); found {
		utils.Respond(c, http.StatusOK, true, "Agentes obtenidos correctamente (cache)", cached, nil)
		return
	}

	rows, err := ac.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener agentes", nil, err)
		return
	}

	utils.CacheAgentsAll(rows)

	utils.Respond(c, http.StatusOK, true, "Agentes obtenidos correctamente", rows, nil)
}

// GET /agents/:user_id
func (ac *AgentController) GetByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetAgentByUserIDFromCache(uint(id)); found {
		utils.Respond(c, http.StatusOK, true, "Agente encontrado (cache)", cached, nil)
		return
	}

	row, err := ac.repo.GetByUserID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Agente no encontrado", nil, err)
		return
	}

	utils.CacheAgentByUserID(uint(id), row)

	utils.Respond(c, http.StatusOK, true, "Agente encontrado", row, nil)
}

// POST /agents
func (ac *AgentController) Create(c *gin.Context) {
	var body models.Agent
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	if err := ac.repo.Delete(uint(body.UserID)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar agente", nil, err)
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

// GET /agents/company/:company_id/agents-with-user-info
func (ac *AgentController) GetAllByCompanyIDWithUserInfo(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetAgentsByCompanyWithUserInfoFromCache(uint(companyID)); found {
		utils.Respond(c, http.StatusOK, true, "Agentes obtenidos correctamente (cache)", cached, nil)
		return
	}

	rows, err := ac.repo.GetAllByCompanyIDWithUserInfo(uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener agentes con info de usuario por company_id", nil, err)
		return
	}

	utils.CacheAgentsByCompanyWithUserInfo(uint(companyID), rows)

	utils.Respond(c, http.StatusOK, true, "Agentes con info de usuario por company_id obtenidos correctamente", rows, nil)
}

// GET /agents/company/:company_id/department/:department_id/agents-with-user-info
func (ac *AgentController) GetAllByCompanyIDAndDepartmentIDWithUserInfo(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	departmentID, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetAgentsByCompanyAndDepartmentWithUserInfoFromCache(
		uint(companyID),
		uint(departmentID),
	); found {
		utils.Respond(c, http.StatusOK, true, "Agentes obtenidos correctamente (cache)", cached, nil)
		return
	}

	rows, err := ac.repo.GetAllByCompanyIDAndDepartmentIDWithUserInfo(uint(companyID), uint(departmentID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener agentes con info de usuario por company_id y department_id", nil, err)
		return
	}

	utils.CacheAgentsByCompanyAndDepartmentWithUserInfo(uint(companyID), uint(departmentID), rows)

	utils.Respond(c, http.StatusOK, true, "Agentes con info de usuario por company_id y department_id obtenidos correctamente", rows, nil)
}

// GET /agents-with-user-info
func (ac *AgentController) GetAllWithUserInfo(c *gin.Context) {

	if cached, found := utils.GetAgentsWithUserInfoAllFromCache(); found {
		utils.Respond(c, http.StatusOK, true, "Agentes con info de usuario obtenidos correctamente (cache)", cached, nil)
		return
	}

	rows, err := ac.repo.GetAllWithUserInfo()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener agentes con info de usuario", nil, err)
		return
	}

	utils.CacheAgentsWithUserInfoAll(rows)

	utils.Respond(c, http.StatusOK, true, "Agentes con info de usuario obtenidos correctamente", rows, nil)
}

// GET /agents-with-user-info/:user_id
func (ac *AgentController) GetByUserIDWithUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetAgentWithUserInfoByUserIDFromCache(uint(id)); found {
		utils.Respond(c, http.StatusOK, true, "Agente con info de usuario encontrado (cache)", cached, nil)
		return
	}

	row, err := ac.repo.GetByUserIDWithUserInfo(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Agente con info de usuario no encontrado", nil, err)
		return
	}

	utils.CacheAgentWithUserInfoByUserID(uint(id), row)

	utils.Respond(c, http.StatusOK, true, "Agente con info de usuario encontrado", row, nil)
}

// GET /agents/non-agents
func (ac *AgentController) GetAllNonAgents(c *gin.Context) {

	if cached, found := utils.GetNonAgentsAllFromCache(); found {
		utils.Respond(c, http.StatusOK, true, "Usuarios no agentes obtenidos correctamente (cache)", cached, nil)
		return
	}

	rows, err := ac.repo.GetAllNonAgents()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener usuarios no agentes", nil, err)
		return
	}

	utils.CacheNonAgentsAll(rows)

	utils.Respond(c, http.StatusOK, true, "Usuarios no agentes obtenidos correctamente", rows, nil)
}

type UnifiedAgentRequest struct {
	Email         string `json:"email" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Phone         string `json:"phone"`
	Password      string `json:"password" binding:"required"`
	CompanyID     uint   `json:"company_id" binding:"required"`
	RoleID        uint   `json:"role_id" binding:"required"`
	DepartmentIDs []uint `json:"department_ids"`
}

func (ac *AgentController) CreateUnifiedAgent(c *gin.Context) {
	var body UnifiedAgentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// 1. Hashear password
	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar contraseña", nil, err)
		return
	}

	// 2. Ejecutar creación transaccional
	user, err := ac.repo.CreateUnifiedAgent(body.Email, body.FullName, body.Phone, hash, body.CompanyID, body.RoleID, body.DepartmentIDs)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, err.Error(), nil, err)
		return
	}

	// 3. Invalidar cachés
	utils.InvalidateAgentsCache()
	utils.DeleteCache(utils.AgentsByCompanyWithUserInfoKey(body.CompanyID))

	utils.Respond(c, http.StatusCreated, true, "Agente creado y configurado correctamente!", user, nil)
}
