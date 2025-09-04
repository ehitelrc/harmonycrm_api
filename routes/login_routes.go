package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(r *gin.RouterGroup) {
	loginCtrl := controllers.NewLoginController()
	r.POST("/auth/login", loginCtrl.Login)
}
