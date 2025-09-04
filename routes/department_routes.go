package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterDepartmentRoutes(r *gin.RouterGroup) {
	controller := controllers.NewDepartmentController()

	r.GET("/departments/company/:company_id", controller.GetByCompany)
	r.GET("/departments/:id", controller.GetByID)
	r.POST("/departments", controller.Create)
	r.PUT("/departments", controller.Update)
	r.DELETE("/departments/:id", controller.Delete)
}
