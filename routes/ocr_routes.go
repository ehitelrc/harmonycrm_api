package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func OCRRoutes(router *gin.RouterGroup, ctrl *controllers.OCRController) {
	ocr := router.Group("/ocr")
	{
		ocr.POST("/", ctrl.OCR)
	}
}
