package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCaseItemRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewCaseItemController()

	// CRUD básico
	r.GET("/case-items/case/:case_id", ctrl.GetAllItemsByCaseID)
	r.GET("/case-items/:id", ctrl.GetItemByCaseItemID)

	r.POST("/case-items", ctrl.CreateCaseItem)
	r.PUT("/case-items", ctrl.UpdateCaseItem)

	// Eliminar
	r.DELETE("/case-items/:id", ctrl.DeleteCaseItem)

}
