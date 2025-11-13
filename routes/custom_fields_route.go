package routes

import (
	"harmony_api/controllers"
	"harmony_api/repository"

	"github.com/gin-gonic/gin"
)

func RegristerCustomFieldRoutes(r *gin.RouterGroup) {
	controller := controllers.NewCustomFieldController(
		repository.NewCustomFieldRepository(),
	)

	group := r.Group("/custom-fields")

	// Save custom field value
	group.POST("/entity/value", controller.SaveEntityCustomFieldValue)

}
