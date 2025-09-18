package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyController struct {
	repo *repository.CompanyRepository
}

func NewCompanyController() *CompanyController {
	return &CompanyController{
		repo: repository.NewCompanyRepository(),
	}
}

func (c *CompanyController) GetAll(ctx *gin.Context) {
	companies, err := c.repo.GetAll()
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener compañías", nil, err)
		return
	}
	utils.Respond(ctx, http.StatusOK, true, "Lista de compañías", companies, nil)
}

func (c *CompanyController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	company, err := c.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Compañía no encontrada", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Compañía encontrada", company, nil)
}

func (c *CompanyController) Create(ctx *gin.Context) {
	var company models.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.Create(&company); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear compañía", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Compañía creada correctamente", company, nil)
}

func (c *CompanyController) Update(ctx *gin.Context) {
	var company models.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.Update(&company); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar compañía", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Compañía actualizada correctamente", company, nil)
}

func (c *CompanyController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar compañía", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Compañía eliminada correctamente", nil, nil)
}

func (c *CompanyController) GetByUserID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de usuario inválido", nil, err)
		return
	}

	companies, err := c.repo.GetByUserID(uint(userID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener compañías del usuario", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Compañías del usuario obtenidas correctamente", companies, nil)
}
