package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCaseDashboardRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewCaseDashboardController()

	r.GET("/case_dashboard/company/:company_id", ctrl.GetCompanyDashboard)
}
