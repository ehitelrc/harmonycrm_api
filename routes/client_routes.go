package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewClientController()

	r.GET("/clients", ctrl.GetAll)
	r.GET("/clients/:id", ctrl.GetByID)
	r.POST("/clients", ctrl.Create)
	r.PUT("/clients", ctrl.Update) // objeto completo con id
	r.DELETE("/clients/:id", ctrl.Delete)
}
