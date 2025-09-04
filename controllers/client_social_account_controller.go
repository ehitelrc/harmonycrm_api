package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientSocialAccountController struct {
	repo *repository.ClientSocialAccountRepository
}

func NewClientSocialAccountController() *ClientSocialAccountController {
	return &ClientSocialAccountController{
		repo: repository.NewClientSocialAccountRepository(),
	}
}

// GET /client-social-accounts
func (cc *ClientSocialAccountController) GetAll(c *gin.Context) {
	rows, err := cc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener cuentas sociales", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuentas sociales obtenidas", rows, nil)
}

// GET /client-social-accounts/:id
func (cc *ClientSocialAccountController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := cc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Cuenta social no encontrada", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta social encontrada", row, nil)
}

// GET /client-social-accounts/client/:client_id
func (cc *ClientSocialAccountController) GetByClient(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "client_id inválido", nil, err)
		return
	}
	rows, err := cc.repo.GetByClient(uint(clientID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener cuentas por cliente", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuentas por cliente", rows, nil)
}

// GET /client-social-accounts/channel/:channel_id
func (cc *ClientSocialAccountController) GetByChannel(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "channel_id inválido", nil, err)
		return
	}
	rows, err := cc.repo.GetByChannel(uint(channelID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener cuentas por canal", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuentas por canal", rows, nil)
}

// GET /client-social-accounts/channel/:channel_id/external/:external_id
func (cc *ClientSocialAccountController) GetByChannelAndExternal(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "channel_id inválido", nil, err)
		return
	}
	externalID := c.Param("external_id")
	row, err := cc.repo.GetByChannelAndExternal(uint(channelID), externalID)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Cuenta social no encontrada", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta social encontrada", row, nil)
}

// POST /client-social-accounts
func (cc *ClientSocialAccountController) Create(c *gin.Context) {
	var body models.ClientSocialAccount
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Error al crear cuenta social (verifique unicidad channel_id+external_id)", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Cuenta social creada correctamente", body, nil)
}

// PUT /client-social-accounts  (objeto completo con id)
func (cc *ClientSocialAccountController) Update(c *gin.Context) {
	var body models.ClientSocialAccount
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Error al actualizar cuenta social", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta social actualizada", body, nil)
}

// DELETE /client-social-accounts/:id
func (cc *ClientSocialAccountController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := cc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar cuenta social", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta social eliminada", nil, nil)
}

// PATCH /client-social-accounts/:id/activate
func (cc *ClientSocialAccountController) Activate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := cc.repo.SetActive(uint(id), true); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "No se pudo activar la cuenta", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta activada", nil, nil)
}

// PATCH /client-social-accounts/:id/deactivate
func (cc *ClientSocialAccountController) Deactivate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := cc.repo.SetActive(uint(id), false); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "No se pudo desactivar la cuenta", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cuenta desactivada", nil, nil)
}
