package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterMessageStatusReportRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewMessageStatusReportController()
	r.GET("/reports/message-statuses/company/:company_id", ctrl.GetSummary)
}
