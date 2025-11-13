package routes

import (
	"harmony_api/config"
	"harmony_api/controllers"
	"harmony_api/repository"

	"github.com/gin-gonic/gin"
)

func RegisterCustomListRoutes(r *gin.RouterGroup) {
	controller := controllers.NewCustomListController(
		repository.NewCustomListRepository(config.DB),
	)

	group := r.Group("/custom-lists")

	// Obtener todas las listas dinámicas para una entidad específica
	// Ejemplo: /api/custom-lists/clients/25
	group.GET("/:entity/:entity_id", controller.GetLists)

	// Get lists
	group.GET("/list", controller.GetAllLists)

	// Guardar selección de una lista para una entidad
	group.POST("/select", controller.SaveSelection)

	// Create list value
	group.POST("/value", controller.CreateListValue)

	// Update list value
	group.PUT("/value", controller.UpdateListValue)

	// Delete list value
	group.DELETE("/value/:id", controller.DeleteListValue)

}
