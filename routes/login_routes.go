package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(r *gin.RouterGroup) {
	loginCtrl := controllers.NewLoginController()
	r.POST("/auth/login", loginCtrl.Login)

	r.POST("/auth/logout", loginCtrl.Logout)

	r.GET("/auth/permissions/:user_id/:company_id", loginCtrl.GetPermissions)
}
