package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewUserController()

	r.GET("/users", ctrl.GetAll)
	r.GET("/users/:id", ctrl.GetByID)
	r.POST("/users", ctrl.Create)
	r.PUT("/users", ctrl.Update) // objeto completo con id

	// Set New password
	r.PUT("/users/:id/password", ctrl.UpdatePassword)

	r.DELETE("/users/:id", ctrl.Delete)

}
