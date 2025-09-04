package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCompanyRoutes(r *gin.RouterGroup) {
	controller := controllers.NewCompanyController()

	r.GET("/companies", controller.GetAll)
	r.GET("/companies/:id", controller.GetByID)
	r.POST("/companies", controller.Create)
	r.PUT("/companies", controller.Update)
	r.DELETE("/companies/:id", controller.Delete)
}
