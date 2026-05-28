package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCasesBulkRoutes(router *gin.RouterGroup) {
	controller := controllers.NewCasesBulkController()

	group := router.Group("/v1/cases/bulk-close")
	{
		group.GET("/search", controller.Search)
		group.POST("/execute", controller.Execute)
	}
}
