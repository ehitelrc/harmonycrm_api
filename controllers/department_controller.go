package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{
		repo: repository.NewDepartmentRepository(),
	}
}

func (c *DepartmentController) GetAll(ctx *gin.Context) {
	departments, err := c.repo.GetAll()
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener departamentos", nil, err)
		return
	}
	utils.Respond(ctx, http.StatusOK, true, "Lista de departamentos", departments, nil)
}

func (c *DepartmentController) GetByCompany(ctx *gin.Context) {
	companyID, err := strconv.Atoi(ctx.Param("company_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Parámetro company_id inválido", nil, err)
		return
	}

	departments, err := c.repo.GetByCompanyID(uint(companyID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener departamentos", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Departamentos encontrados", departments, nil)
}

func (c *DepartmentController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	dept, err := c.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Departamento no encontrado", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Departamento encontrado", dept, nil)
}

func (c *DepartmentController) Create(ctx *gin.Context) {
	var dept models.Department
	if err := ctx.ShouldBindJSON(&dept); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.Create(&dept); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear departamento", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Departamento creado correctamente", dept, nil)
}

func (c *DepartmentController) Update(ctx *gin.Context) {
	var dept models.Department
	if err := ctx.ShouldBindJSON(&dept); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.Update(&dept); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar departamento", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Departamento actualizado correctamente", dept, nil)
}

func (c *DepartmentController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar departamento", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Departamento eliminado correctamente", nil, nil)
}
