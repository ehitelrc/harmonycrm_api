package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewRoleController()

	r.GET("/roles", ctrl.GetAll)
	r.GET("/roles/:id", ctrl.GetByID)
	r.POST("/roles", ctrl.Create)
	r.PUT("/roles", ctrl.Update)
	r.DELETE("/roles/:id", ctrl.Delete)
}
