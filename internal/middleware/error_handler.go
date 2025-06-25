package middleware

import (
	"net/http"
	"strings"

	"github.com/UliVargas/blog-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse representa la estructura básica de respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

// HandleValidationError maneja errores de validación de forma centralizada
func HandleValidationError(c *gin.Context, err error) {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		response := utils.CreateValidationErrorResponse(validationErr)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Si no es un error de validación, devolver error genérico
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Datos inválidos",
		Details: err.Error(),
	})
}

// HandleInternalError maneja errores internos del servidor
func HandleInternalError(c *gin.Context, err error, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   message,
		Details: err.Error(),
	})
}

// HandleNotFoundError maneja errores de recurso no encontrado
func HandleNotFoundError(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: message,
	})
}

// HandleBadRequestError maneja errores de solicitud incorrecta
func HandleBadRequestError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: message,
	})
}

// HandleDatabaseError maneja errores de base de datos y los convierte en mensajes amigables
func HandleDatabaseError(c *gin.Context, err error, defaultMessage string) {
	errorMsg := err.Error()
	userFriendlyMsg := defaultMessage

	// Detectar errores comunes de base de datos
	if strings.Contains(errorMsg, "duplicate key value violates unique constraint") {
		if strings.Contains(errorMsg, "uni_users_email") {
			userFriendlyMsg = "El email ya está registrado"
		} else if strings.Contains(errorMsg, "uni_users_username") {
			userFriendlyMsg = "El nombre de usuario ya está en uso"
		} else {
			userFriendlyMsg = "Ya existe un registro con estos datos"
		}
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: userFriendlyMsg,
		})
		return
	}

	// Detectar errores de clave foránea
	if strings.Contains(errorMsg, "foreign key constraint") {
		userFriendlyMsg = "No se puede completar la operación debido a dependencias"
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: userFriendlyMsg,
		})
		return
	}

	// Detectar errores de conexión
	if strings.Contains(errorMsg, "connection refused") || strings.Contains(errorMsg, "no connection") {
		userFriendlyMsg = "Error de conexión con la base de datos"
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: userFriendlyMsg,
		})
		return
	}

	// Para otros errores de base de datos, usar mensaje genérico
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: userFriendlyMsg,
	})
}
