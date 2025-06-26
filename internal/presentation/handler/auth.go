package handler

import (
	"net/http"

	services "github.com/UliVargas/blog-go/internal/application/service"
	"github.com/UliVargas/blog-go/internal/domain/dto"
	domainService "github.com/UliVargas/blog-go/internal/domain/service"
	"github.com/UliVargas/blog-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domainService.AuthServiceInterface
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "Inicio de sesión exitoso",
		Token:   token,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user dto.RegisterRequest

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
	if err := h.authService.Register(user.ToUser()); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}
