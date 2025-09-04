package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewPermissionController()

	r.GET("/permissions", ctrl.GetAll)
	r.GET("/permissions/:id", ctrl.GetByID)

	// By users_id
	r.GET("/permissions/user/:user_id/company/:company_id", ctrl.GetByUserID)

	r.POST("/permissions", ctrl.Create)
	r.PUT("/permissions", ctrl.Update) // objeto completo con id
	r.DELETE("/permissions/:id", ctrl.Delete)
}
