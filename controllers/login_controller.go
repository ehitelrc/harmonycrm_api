package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"

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
