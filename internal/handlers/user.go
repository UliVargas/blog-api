package handlers

import (
	"net/http"
	"strconv"

	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/service"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/UliVargas/blog-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.HandleError(c, appErrors.ErrInvalidID)
		return
	}
	user, err := h.userService.GetByID(uint(idUint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user models.CreateUserRequest

	// 1. Bind JSON a la estructura de usuario
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleBadRequest(c, "Datos inválidos proporcionados")
		return
	}

	// 2. Validar los datos con estructura
	if err := utils.GetValidator().Struct(&user); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 3. Crear el usuario
	if err := h.userService.Create(user.ToUser()); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}
