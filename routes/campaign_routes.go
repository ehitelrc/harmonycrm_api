package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCampaignRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewCampaignController()

	r.GET("/campaigns/company/:company_id", ctrl.GetByCompany)
	r.GET("/campaigns/:id", ctrl.GetByID)
	r.POST("/campaigns", ctrl.Create)
	r.PUT("/campaigns", ctrl.Update) // objeto completo con id
	r.DELETE("/campaigns/:id", ctrl.Delete)
}
