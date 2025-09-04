package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgentDepartmentAssignmentController struct {
	repo *repository.AgentDepartmentAssignmentRepository
}

func NewAgentDepartmentAssignmentController() *AgentDepartmentAssignmentController {
	return &AgentDepartmentAssignmentController{
		repo: repository.NewAgentDepartmentAssignmentRepository(),
	}
}

// GET /agent-department-assignments/:id
func (ac *AgentDepartmentAssignmentController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := ac.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Asignación no encontrada", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación encontrada", row, nil)
}

// GET /agent-department-assignments/agent/:agent_id
func (ac *AgentDepartmentAssignmentController) GetByAgent(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "agent_id inválido", nil, err)
		return
	}
	rows, err := ac.repo.GetByAgent(uint(agentID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener asignaciones por agente", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por agente", rows, nil)
}

// GET /agent-department-assignments/department/:department_id
func (ac *AgentDepartmentAssignmentController) GetByDepartment(c *gin.Context) {
	deptID, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}
	rows, err := ac.repo.GetByDepartment(uint(deptID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener asignaciones por departamento", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por departamento", rows, nil)
}

// POST /agent-department-assignments
func (ac *AgentDepartmentAssignmentController) Create(c *gin.Context) {
	var body models.AgentDepartmentAssignment
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := ac.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Asignación creada correctamente", body, nil)
}

// PUT /agent-department-assignments  (objeto completo con id)
func (ac *AgentDepartmentAssignmentController) Update(c *gin.Context) {
	var body models.AgentDepartmentAssignment
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := ac.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación actualizada correctamente", body, nil)
}

// DELETE /agent-department-assignments/:id
func (ac *AgentDepartmentAssignmentController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := ac.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación eliminada correctamente", nil, nil)
}
