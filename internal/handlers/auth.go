package handlers

import (
	"net/http"

	"github.com/UliVargas/blog-go/internal/service"
	"github.com/UliVargas/blog-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBadRequest(c, "Datos inválidos")
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inicio de sesión exitoso",
		"token":   token,
	})
}
