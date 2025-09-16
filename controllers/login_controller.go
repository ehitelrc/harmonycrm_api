package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	repo *repository.LoginRepository
}

func NewLoginController() *LoginController {
	return &LoginController{
		repo: repository.NewLoginRepository(),
	}
}

type LoginRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CompanyID int    `json:"company_id" binding:"required"`
}

func (ctrl *LoginController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	// Validar usuario en la función SQL
	data, err := ctrl.repo.Login(req.Email, req.Password, req.CompanyID)
	if err != nil {
		utils.Respond(c, http.StatusOK, false, "Error al validar usuario", nil, err)
		return
	}

	// if data has information
	if data == nil {
		utils.Respond(c, http.StatusOK, false, "Usuario o contraseña incorrectos", nil, nil)
		return
	}

	data.Token = "NOHAY_TOKEN"

	utils.Respond(c, http.StatusOK, true, "Login exitoso", data, nil)
}

func (ctrl *LoginController) Logout(c *gin.Context) {
	// Aquí puedes agregar la lógica para invalidar el token si es necesario
	utils.Respond(c, http.StatusOK, true, "Logout exitoso", nil, nil)
}

func (c *LoginController) GetPermissions(ctx *gin.Context) {

	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		utils.Respond(ctx, 400, false, "ID de usuario inválido", nil, err)
		return
	}

	companyID, err := strconv.Atoi(ctx.Param("company_id"))
	if err != nil {
		utils.Respond(ctx, 400, false, "ID de compañía inválido", nil, err)
		return
	}

	perms, err := c.repo.GetUserPermissions(userID, companyID)
	if err != nil {
		utils.Respond(ctx, 500, false, "Error obteniendo permisos", nil, err)
		return
	}

	utils.Respond(ctx, 200, true, "Permisos obtenidos", perms, nil)
}
