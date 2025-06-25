package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse representa una respuesta exitosa estándar
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SendSuccess envía una respuesta de éxito estándar
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// SendCreated envía una respuesta de creación exitosa
func SendCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// SendNoContent envía una respuesta sin contenido (para eliminaciones)
func SendNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
