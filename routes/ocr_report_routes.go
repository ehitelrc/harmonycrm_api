package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOcrReportRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewOcrReportController()
	r.GET("/reports/ocr/company/:company_id", ctrl.GetOcrReport)
}
