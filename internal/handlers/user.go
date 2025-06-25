package handlers

import (
	"net/http"
	"strconv"

	"github.com/UliVargas/blog-go/internal/middleware"
	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService}
}

var validate = validator.New()

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		middleware.HandleDatabaseError(c, err, "Error al obtener usuarios")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.HandleBadRequestError(c, "Formato de ID inválido")
		return
	}
	user, err := h.userService.GetByID(uint(idUint))
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "usuario no encontrado" {
			middleware.HandleNotFoundError(c, "Usuario no encontrado")
		} else {
			middleware.HandleDatabaseError(c, err, "Error al obtener el usuario")
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user models.CreateUserRequest

	// 1. Bind JSON a la estructura de usuario
	if err := c.ShouldBindJSON(&user); err != nil {
		middleware.HandleBadRequestError(c, "Datos inválidos proporcionados")
		return
	}

	// 2. Validar los datos con estructura
	if err := validate.Struct(&user); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	// 3. Crear el usuario
	if err := h.userService.Create(user.ToUser()); err != nil {
		middleware.HandleDatabaseError(c, err, "Error al crear el usuario")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}
