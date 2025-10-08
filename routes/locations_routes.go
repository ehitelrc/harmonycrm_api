package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func LocationsRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewLocationsController()

	locations := r.Group("/locations")

	// Countries routes
	locations.GET("/countries", ctrl.GetAllCountries)
	locations.GET("/countries/:id", ctrl.GetCountryByID)
	locations.POST("/countries", ctrl.CreateCountry)
	locations.PUT("/countries/:id", ctrl.UpdateCountry)
	locations.DELETE("/countries/:id", ctrl.DeleteCountry)

	// Provinces routes

	locations.GET("/provinces/country/:country_id", ctrl.GetProvincesByCountryID)
	locations.GET("/provinces/:province_id", ctrl.GetCantonsByProvinceID)
	locations.POST("/provinces", ctrl.CreateProvince)
	locations.PUT("/provinces/", ctrl.UpdateProvince)
	locations.DELETE("/provinces/:province_id", ctrl.DeleteProvince)

	// Cantons routes
	locations.GET("/cantons/province/:province_id", ctrl.GetCantonsByProvinceID)
	locations.GET("/cantons/:canton_id", ctrl.GetDistrictsByCantonID)
	locations.POST("/cantons", ctrl.CreateCanton)
	locations.PUT("/cantons/", ctrl.UpdateCanton)
	locations.DELETE("/cantons/:canton_id", ctrl.DeleteCanton)

	// Districts routes
	locations.GET("/districts/canton/:canton_id", ctrl.GetDistrictsByCantonID)
	locations.GET("/district/:district_id", ctrl.GetDistrictByID)
	locations.POST("/districts", ctrl.CreateDistrict)
	locations.PUT("/districts/", ctrl.UpdateDistrict)
	locations.DELETE("/districts/:district_id", ctrl.DeleteDistrict)

}
