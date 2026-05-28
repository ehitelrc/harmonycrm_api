package controllers

import (
	"harmony_api/dto"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CasesBulkController struct {
	repo     *repository.CasesBulkRepository
	userRepo *repository.UserRepository
}

func NewCasesBulkController() *CasesBulkController {
	return &CasesBulkController{
		repo:     repository.NewCasesBulkRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

func (c *CasesBulkController) Search(ctx *gin.Context) {
	companyIDStr := ctx.Query("company_id")
	departmentIDStr := ctx.Query("department_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		// fallback to local date format if needed, but standardizing to RFC3339 is better
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.Respond(ctx, http.StatusBadRequest, false, "start_date inválido", nil, err)
			return
		}
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.Respond(ctx, http.StatusBadRequest, false, "end_date inválido", nil, err)
			return
		}
	}

	// Fetch cases using repository
	results, err := c.repo.SearchCasesToBulkClose(uint(companyID), uint(departmentID), startDate, endDate)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error buscando casos", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Casos obtenidos correctamente", results, nil)
}

func (c *CasesBulkController) Execute(ctx *gin.Context) {
	var req dto.BulkCloseExecuteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	userID := req.UserID

	if userID == 0 {
		utils.Respond(ctx, http.StatusUnauthorized, false, "No se proporcionó user_id", nil, nil)
		return
	}

	user, err := c.userRepo.GetByID(userID)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error obteniendo datos del usuario", nil, err)
		return
	}

	// Check password
	if !utils.ComparePasswordHash(req.Password, user.PasswordHash) {
		utils.Respond(ctx, http.StatusUnauthorized, false, "Contraseña incorrecta", nil, nil)
		return
	}

	if len(req.CaseIDs) == 0 {
		utils.Respond(ctx, http.StatusBadRequest, false, "No se seleccionaron casos", nil, nil)
		return
	}

	if err := c.repo.ExecuteBulkClose(req.CaseIDs); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al cerrar casos", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Casos cerrados correctamente", nil, nil)
}
