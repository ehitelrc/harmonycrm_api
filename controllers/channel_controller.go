package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChannelController struct {
	Repo repository.ChannelRepository
}

func (ctrl *ChannelController) GetAll(c *gin.Context) {
	channels, err := ctrl.Repo.GetAllChannels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al obtener los canales", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": channels})
}

func (ctrl *ChannelController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID inválido"})
		return
	}

	channel, err := ctrl.Repo.GetChannelByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Canal no encontrado", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": channel})
}

func (ctrl *ChannelController) Create(c *gin.Context) {
	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos", "error": err.Error()})
		return
	}

	if err := ctrl.Repo.CreateChannel(&channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al crear canal", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": channel})
}

func (ctrl *ChannelController) Update(c *gin.Context) {
	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos", "error": err.Error()})
		return
	}

	if err := ctrl.Repo.UpdateChannel(&channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al actualizar canal", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": channel})
}

func (ctrl *ChannelController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID inválido"})
		return
	}

	if err := ctrl.Repo.DeleteChannel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al eliminar canal", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Canal eliminado"})
}

func (ctrl *ChannelController) CreateWhatsappTemplate(c *gin.Context) {
	var template models.ChannelWhatsAppTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos", "error": err.Error()})
		return
	}

	if err := ctrl.Repo.CreateWhatsappTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al crear plantilla de Whatsapp", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": template})
}

func (ctrl *ChannelController) GetWhatsappTemplatesByCompanyID(c *gin.Context) {
	companyIDParam := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de empresa inválido"})
		return
	}

	templates, err := ctrl.Repo.GetWhatsappTemplatesByCompanyID(uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al obtener plantillas de Whatsapp", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": templates})
}
