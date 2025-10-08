package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationsController struct {
	repo *repository.LocationRepositories
}

func NewLocationsController() *LocationsController {
	return &LocationsController{repo: repository.NewCountriesRepository()}
}

func (c *LocationsController) GetAllCountries(ctx *gin.Context) {
	locations, err := c.repo.GetAllCountries()

	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener países", nil, err)
		return
	}
	utils.Respond(ctx, http.StatusOK, true, "Lista de países", locations, nil)
}

func (c *LocationsController) GetCountryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	country, err := c.repo.GetCountryByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "País no encontrado", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "País encontrado", country, nil)
}

func (c *LocationsController) CreateCountry(ctx *gin.Context) {
	var country models.Country
	if err := ctx.ShouldBindJSON(&country); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.CreateCountry(&country); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear país", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "País creado exitosamente", country, nil)
}

func (c *LocationsController) UpdateCountry(ctx *gin.Context) {

	var country models.Country
	if err := ctx.ShouldBindJSON(&country); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.UpdateCountry(&country); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar país", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "País actualizado exitosamente", country, nil)
}

func (c *LocationsController) DeleteCountry(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.DeleteCountry(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar país", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "País eliminado exitosamente", nil, nil)
}

func (c *LocationsController) GetProvincesByCountryID(ctx *gin.Context) {
	countryID, err := strconv.Atoi(ctx.Param("country_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de país inválido", nil, err)
		return
	}

	provinces, err := c.repo.GetProvincesByCountryID(uint(countryID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener provincias", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Lista de provincias", provinces, nil)
}

func (c *LocationsController) GetCantonsByProvinceID(ctx *gin.Context) {
	provinceID, err := strconv.Atoi(ctx.Param("province_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de provincia inválido", nil, err)
		return
	}

	cantons, err := c.repo.GetCantonsByProvinceID(uint(provinceID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener cantones", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Lista de cantones", cantons, nil)
}

// Create province
func (c *LocationsController) CreateProvince(ctx *gin.Context) {
	var province models.Province
	if err := ctx.ShouldBindJSON(&province); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.CreateProvince(&province); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear provincia", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Provincia creada exitosamente", province, nil)
}

// Update province
func (c *LocationsController) UpdateProvince(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	var province models.Province
	if err := ctx.ShouldBindJSON(&province); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	province.ID = uint(id)

	if err := c.repo.UpdateProvince(&province); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar provincia", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Provincia actualizada exitosamente", province, nil)
}

// Delete province
func (c *LocationsController) DeleteProvince(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.DeleteProvince(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar provincia", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Provincia eliminada exitosamente", nil, nil)
}

// Add canton
func (c *LocationsController) CreateCanton(ctx *gin.Context) {
	var canton models.Canton
	if err := ctx.ShouldBindJSON(&canton); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.CreateCanton(&canton); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear cantón", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Cantón creado exitosamente", canton, nil)
}

// Update canton
func (c *LocationsController) UpdateCanton(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	var canton models.Canton
	if err := ctx.ShouldBindJSON(&canton); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	canton.ID = uint(id)

	if err := c.repo.UpdateCanton(&canton); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar cantón", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Cantón actualizado exitosamente", canton, nil)
}

// Delete canton
func (c *LocationsController) DeleteCanton(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.DeleteCanton(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar cantón", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Cantón eliminado exitosamente", nil, nil)
}

// Districts

func (c *LocationsController) GetDistrictsByCantonID(ctx *gin.Context) {
	cantonID, err := strconv.Atoi(ctx.Param("canton_id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de cantón inválido", nil, err)
		return
	}

	districts, err := c.repo.GetDistrictsByCantonID(uint(cantonID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener distritos", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Lista de distritos", districts, nil)
}

// Get district by id
func (c *LocationsController) GetDistrictByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	district, err := c.repo.GetDistrictByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Distrito no encontrado", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Distrito encontrado", district, nil)
}

func (c *LocationsController) CreateDistrict(ctx *gin.Context) {
	var district models.District
	if err := ctx.ShouldBindJSON(&district); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.CreateDistrict(&district); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear distrito", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Distrito creado exitosamente", district, nil)
}

func (c *LocationsController) UpdateDistrict(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	var district models.District
	if err := ctx.ShouldBindJSON(&district); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}
	district.ID = uint(id)

	if err := c.repo.UpdateDistrict(&district); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar distrito", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Distrito actualizado exitosamente", district, nil)
}

func (c *LocationsController) DeleteDistrict(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.DeleteDistrict(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar distrito", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Distrito eliminado exitosamente", nil, nil)
}
