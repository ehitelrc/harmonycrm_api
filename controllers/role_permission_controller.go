package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RolePermissionController struct {
	repo *repository.RolePermissionRepository
}

func NewRolePermissionController() *RolePermissionController {
	return &RolePermissionController{repo: repository.NewRolePermissionRepository()}
}

// DTOs

type replaceRequest struct {
	PermissionIDs []uint `json:"permission_ids"` // puede venir vacío para dejar el rol sin permisos
}

// GET /role-permissions/role/:role_id
func (rc *RolePermissionController) GetByRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}
	rows, err := rc.repo.GetByRole(uint(roleID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener permisos por rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos por rol", rows, nil)
}

// GET /role-permissions/permission/:permission_id
func (rc *RolePermissionController) GetByPermission(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "permission_id inválido", nil, err)
		return
	}
	rows, err := rc.repo.GetByPermission(uint(pid))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener roles por permiso", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Roles por permiso", rows, nil)
}

// POST /role-permissions  (asignar uno)
func (rc *RolePermissionController) Assign(c *gin.Context) {
	var body []models.AssignRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	if err := rc.repo.AssignBatch(body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar permisos al rol", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Permisos asignados al rol", nil, nil)

}

// DELETE /role-permissions/role/:role_id/permission/:permission_id (desasignar uno)
func (rc *RolePermissionController) Unassign(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}
	pid, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "permission_id inválido", nil, err)
		return
	}
	if err := rc.repo.Unassign(uint(roleID), uint(pid)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al desasignar permiso del rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permiso desasignado del rol", nil, nil)
}

// PUT /role-permissions/role/:role_id (reemplazo total de permisos del rol)
func (rc *RolePermissionController) ReplaceForRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}
	var body replaceRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := rc.repo.ReplaceRolePermissions(uint(roleID), body.PermissionIDs); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar permisos del rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos del rol actualizados", gin.H{
		"role_id":        roleID,
		"permission_ids": body.PermissionIDs,
	}, nil)
}

// GET /role-permissions/view/role/:role_id
func (rc *RolePermissionController) GetViewByRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}
	rows, err := rc.repo.GetViewByRole(uint(roleID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener permisos por rol (vista)", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos por rol (vista)", rows, nil)
}
