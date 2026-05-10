package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClosedCasesRoutes(api *gin.RouterGroup, closedCasesController *controllers.ClosedCasesController) {
	// Rutas bajo /api/closed-cases
	closedCasesGroup := api.Group("/closed-cases")
	{
		// /api/closed-cases/:sender_id
		closedCasesGroup.GET("/:sender_id", closedCasesController.GetClosedCases)

		// /api/closed-cases/messages/:case_id
		closedCasesGroup.GET("/messages/:case_id", closedCasesController.GetCaseMessages)
	}
}
