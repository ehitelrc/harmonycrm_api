package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	repo   *repository.ClientRepository
	cfrepo *repository.CustomFieldRepository
}

func NewClientController() *ClientController {
	return &ClientController{
		repo:   repository.NewClientRepository(),
		cfrepo: repository.NewCustomFieldRepository(),
	}
}

// GET /clients
func (cc *ClientController) GetAll(c *gin.Context) {
	rows, err := cc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener clientes", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Clientes obtenidos correctamente", rows, nil)
}

// GET /clients/:id
func (cc *ClientController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := cc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Cliente no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cliente encontrado", row, nil)
}

// POST /clients
func (cc *ClientController) Create(c *gin.Context) {
	var body models.Client
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear cliente", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Cliente creado correctamente", body, nil)
}

// PUT /clients  (recibe objeto completo con id)
func (cc *ClientController) Update(c *gin.Context) {
	var body models.Client
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar cliente", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cliente actualizado correctamente", body, nil)
}

// DELETE /clients/:id
func (cc *ClientController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := cc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar cliente", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Cliente eliminado correctamente", nil, nil)
}

// POST /clients/leads
func (cc *ClientController) CreateLead(c *gin.Context) {
	var body models.LeadRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := cc.repo.CreateLead(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear lead", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Lead creado correctamente", body, nil)

}

func (cc *ClientController) GetCustomFields(c *gin.Context) {

	entityIDParam := c.Param("entity_id")
	entityID, err := strconv.Atoi(entityIDParam)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID de entidad inválido", nil, err)
		return
	}

	fields, err := cc.cfrepo.GetFields("clients", uint(entityID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener campos personalizados", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Campos personalizados obtenidos correctamente", fields, nil)
}

// controllers/client_controller.go
func (cc *ClientController) GetDuplicatePhones(c *gin.Context) {
	rows, err := cc.repo.GetClientsWithDuplicatePhones()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener clientes con teléfonos duplicados", nil, err)
		return
	}

	utils.Respond(
		c,
		http.StatusOK,
		true,
		"Clientes con teléfonos duplicados obtenidos correctamente",
		rows,
		nil,
	)
}

// controllers/client_controller.go
func (cc *ClientController) GetDuplicatePhonesDTO(c *gin.Context) {
	rows, err := cc.repo.GetDuplicatePhonesDTO()
	if err != nil {
		utils.Respond(
			c,
			http.StatusInternalServerError,
			false,
			"Error al obtener clientes con teléfonos duplicados",
			nil,
			err,
		)
		return
	}

	utils.Respond(
		c,
		http.StatusOK,
		true,
		"Clientes con teléfonos duplicados obtenidos correctamente",
		rows,
		nil,
	)
}

// controllers/client_controller.go

func (cc *ClientController) GetByExternalID(c *gin.Context) {
	externalID := c.Param("external_id")
	if externalID == "" {
		utils.Respond(c, http.StatusBadRequest, false, "external_id inválido", nil, nil)
		return
	}

	rows, err := cc.repo.GetByExternalID(externalID)
	if err != nil {
		utils.Respond(
			c,
			http.StatusInternalServerError,
			false,
			"Error al consultar clientes por external_id",
			nil,
			err,
		)
		return
	}

	utils.Respond(
		c,
		http.StatusOK,
		true,
		"Clientes obtenidos correctamente",
		rows,
		nil,
	)
}
