package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(r *gin.RouterGroup) {
	controller := controllers.NewDashboardController()

	group := r.Group("/dashboard")

	group.GET("/campaign/summary/company/:company_id", controller.GetCampaignSummary)
	group.GET("/summary/company/:company_id", controller.GetGeneralDashboardByCompany)

	// Funnel por campaña
	group.GET("/campaign/funnel_summary/:campaign_id", controller.GetCampaignFunnelSummary)

}
