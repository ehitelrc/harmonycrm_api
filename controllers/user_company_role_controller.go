package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserCompanyRoleController struct {
	repo *repository.UserCompanyRoleRepository
}

func NewUserCompanyRoleController() *UserCompanyRoleController {
	return &UserCompanyRoleController{repo: repository.NewUserCompanyRoleRepository()}
}

// GET /user-company-roles
func (uc *UserCompanyRoleController) GetAll(c *gin.Context) {
	rows, err := uc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener asignaciones", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones obtenidas", rows, nil)
}

// GET /user-company-roles/:id
func (uc *UserCompanyRoleController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := uc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Asignación no encontrada", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación encontrada", row, nil)
}

// GET /user-company-roles/user/:user_id
func (uc *UserCompanyRoleController) GetByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}
	rows, err := uc.repo.GetByUser(uint(userID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener por usuario", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por usuario", rows, nil)
}

// GET /user-company-roles/company/:company_id
func (uc *UserCompanyRoleController) GetByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}
	rows, err := uc.repo.GetByCompany(uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener por compañía", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por compañía", rows, nil)
}

// GET /user-company-roles/role/:role_id
func (uc *UserCompanyRoleController) GetByRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}
	rows, err := uc.repo.GetByRole(uint(roleID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener por rol", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por rol", rows, nil)
}

// POST /user-company-roles
func (uc *UserCompanyRoleController) Create(c *gin.Context) {
	var body models.UserCompanyRole
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := uc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "No se pudo crear la asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Asignación creada correctamente", body, nil)
}

// PUT /user-company-roles  (objeto completo con id)
func (uc *UserCompanyRoleController) Update(c *gin.Context) {
	var body models.UserCompanyRole
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := uc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "No se pudo actualizar la asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación actualizada correctamente", body, nil)
}

// DELETE /user-company-roles/:id
func (uc *UserCompanyRoleController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := uc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar la asignación", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignación eliminada correctamente", nil, nil)
}

// GET /user-company-roles/company/:company_id/role/:role_id/users-permissions
func (uc *UserCompanyRoleController) GetUsersAndPermissionsByCompanyRole(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "role_id inválido", nil, err)
		return
	}

	rows, err := uc.repo.GetUsersAndPermissionsByCompanyRole(uint(companyID), uint(roleID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener usuarios y permisos", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Usuarios y permisos por compañía y rol", rows, nil)
}

// GET /user-company-roles/company/:company_id/user/:user_id/permissions
func (uc *UserCompanyRoleController) GetPermissionsByCompanyUser(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}

	perms, err := uc.repo.GetPermissionsByCompanyUser(uint(companyID), uint(userID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener permisos por compañía y usuario", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Permisos efectivos por compañía y usuario", perms, nil)
}

// GET /user-company-roles/user/:user_id/company/1
func (uc *UserCompanyRoleController) GetByUserAndCompanyMixed(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "user_id inválido", nil, err)
		return
	}
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	rows, err := uc.repo.GetByUserAndCompanyMixed(uint(userID), uint(companyID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener por usuario y compañía", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Asignaciones por usuario y compañía", rows, nil)
}

// PUT /user-company-roles/batch
func (uc *UserCompanyRoleController) BatchUpdate(c *gin.Context) {
	var bodies []models.UserRoleCompanyManage
	if err := c.ShouldBindJSON(&bodies); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := uc.repo.BatchUpdate(bodies); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "No se pudo realizar la actualización masiva", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Actualización masiva realizada correctamente", bodies, nil)
}
