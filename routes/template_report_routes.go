package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterTemplateReportRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewTemplateReportController()
	r.GET("/reports/templates/company/:company_id", ctrl.GetTemplateReport)
}
