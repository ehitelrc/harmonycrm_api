package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	repo *repository.UserRepository
}

func NewUserController() *UserController {
	return &UserController{repo: repository.NewUserRepository()}
}

// GET /users
func (uc *UserController) GetAll(c *gin.Context) {
	rows, err := uc.repo.GetAll()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener usuarios", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Usuarios obtenidos correctamente", rows, nil)
}

// GET /users/:id
func (uc *UserController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	row, err := uc.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Usuario no encontrado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Usuario encontrado", row, nil)
}

// POST /users
func (uc *UserController) Create(c *gin.Context) {
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := uc.repo.Create(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al crear usuario", nil, err)
		return
	}
	utils.Respond(c, http.StatusCreated, true, "Usuario creado correctamente", body, nil)
}

// PUT /users  (recibe objeto completo con id)
func (uc *UserController) Update(c *gin.Context) {
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}
	if err := uc.repo.Update(&body); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar usuario", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Usuario actualizado correctamente", body, nil)
}

// DELETE /users/:id
func (uc *UserController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	if err := uc.repo.Delete(uint(id)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al eliminar usuario", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Usuario eliminado correctamente", nil, nil)
}
