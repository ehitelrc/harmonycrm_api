package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	repo *repository.RoleRepository
}

func NewRoleController() *RoleController {
	return &RoleController{
		repo: repository.NewRoleRepository(),
	}
}

// GET /roles
func (rc *RoleController) GetAll(c *gin.Context) {
	roles, err := rc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener roles", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Lista de roles", roles, nil)
}

// GET /roles/:id
func (rc *RoleController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	role, err := rc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Rol no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Rol encontrado", role, nil)
}

// POST /roles
func (rc *RoleController) Create(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := rc.repo.Create(&role); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Rol creado correctamente", role, nil)
}

// PUT /roles
func (rc *RoleController) Update(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := rc.repo.Update(&role); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Rol actualizado correctamente", role, nil)
}

// DELETE /roles/:id
func (rc *RoleController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := rc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Rol eliminado correctamente", nil, nil)
}
