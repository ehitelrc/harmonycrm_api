package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	repo *repository.PermissionRepository
}

func NewPermissionController() *PermissionController {
	return &PermissionController{repo: repository.NewPermissionRepository()}
}

// GET /permissions
func (pc *PermissionController) GetAll(c *gin.Context) {
	rows, err := pc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener permisos", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos obtenidos correctamente", rows, nil)
}

// GET /permissions/:id
func (pc *PermissionController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := pc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Permiso no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permiso encontrado", row, nil)
}

func (pc *PermissionController) GetByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de usuario inválido", nil, err)
		return
	}

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de compañía inválido", nil, err)
		return
	}

	rows, err := pc.repo.GetByUserID(uint(userID), uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Permisos no encontrados", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos encontrados", rows, nil)
}

// POST /permissions
func (pc *PermissionController) Create(c *gin.Context) {
	var body models.Permission
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := pc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear permiso", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Permiso creado correctamente", body, nil)
}

// PUT /permissions  (recibe objeto completo con id)
func (pc *PermissionController) Update(c *gin.Context) {
	var body models.Permission
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := pc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar permiso", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permiso actualizado correctamente", body, nil)
}

// DELETE /permissions/:id
func (pc *PermissionController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := pc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar permiso", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permiso eliminado correctamente", nil, nil)
}
