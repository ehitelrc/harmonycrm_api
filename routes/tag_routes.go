package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(router *gin.RouterGroup) {
	tagController := controllers.NewTagController()

	tags := router.Group("/tags")
	{
		tags.POST("", tagController.Create)
		tags.PUT("/:id", tagController.Update)
		tags.DELETE("/:id", tagController.Delete)
		tags.GET("", tagController.GetAll)
	}

	caseTags := router.Group("/cases/:caseId/tags")
	{
		caseTags.GET("", tagController.GetTagsByCase)
		caseTags.POST("", tagController.AssignToCase)
		caseTags.DELETE("/:tagId", tagController.RemoveFromCase)
	}
}
