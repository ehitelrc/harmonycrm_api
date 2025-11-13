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

	// Guardar selección de una lista para una entidad
	group.POST("/select", controller.SaveSelection)

}
