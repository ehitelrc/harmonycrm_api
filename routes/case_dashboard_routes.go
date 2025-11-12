package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCaseDashboardRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewCaseDashboardController()

	r.GET("/case_dashboard/company/:company_id", ctrl.GetCompanyDashboard)
	r.GET("/case_dashboard/company/:company_id/department/:department_id", ctrl.GetDepartmentDashboard)

	// By company and user
	r.GET("/case_dashboard/company/:company_id/user/:user_id", ctrl.GetUserDashboard)
}
