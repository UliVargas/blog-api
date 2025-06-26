package utils

import (
	"errors"
	"net/http"

	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ErrorResponse representa la estructura de respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// HandleError maneja errores de aplicación de forma centralizada
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Verificar si es un AppError personalizado
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error: appErr.Error(),
		})
		return
	}

	// Manejar errores de dominio específicos
	switch {
	case errors.Is(err, appErrors.ErrUserNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Usuario no encontrado",
		})
	case errors.Is(err, appErrors.ErrEmailExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "El email ya está registrado",
		})
	case errors.Is(err, appErrors.ErrUsernameExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "El nombre de usuario ya está en uso",
		})
	case errors.Is(err, appErrors.ErrUserExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "El usuario ya existe",
		})
	case errors.Is(err, appErrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Credenciales inválidas",
		})
	case errors.Is(err, appErrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "No autorizado",
		})
	case errors.Is(err, appErrors.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Datos de entrada inválidos",
		})
	case errors.Is(err, appErrors.ErrInvalidID):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "ID inválido",
		})
	case errors.Is(err, appErrors.ErrDatabaseConnection):
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: "Error de conexión con la base de datos",
		})
	case errors.Is(err, appErrors.ErrForeignKeyViolation):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "No se puede completar la operación debido a dependencias",
		})
	case errors.Is(err, appErrors.ErrDatabaseOperation):
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Error en operación de base de datos",
		})
	default:
		// Error genérico no manejado
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Error interno del servidor",
			Details: err.Error(),
		})
	}
}

// HandleValidationError maneja errores de validación
func HandleValidationError(c *gin.Context, err error) {
	validationErrors := FormatValidationErrors(err)
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Error:  "Datos de validación incorrectos",
		Errors: validationErrors,
	})
}

// HandleBadRequest maneja errores de solicitud incorrecta
func HandleBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: message,
	})
}
