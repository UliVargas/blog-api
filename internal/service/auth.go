package service

import (
	"errors"
	"time"

	"github.com/UliVargas/blog-go/internal/config"
	"github.com/UliVargas/blog-go/internal/repository"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Buscar usuario por email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return "", appErrors.ErrInvalidCredentials
		}
		return "", err
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", appErrors.ErrInvalidCredentials
	}

	// Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	cfg := config.Load()

	tokenString, err := token.SignedString([]byte(cfg.JWTSECRET))
	if err != nil {
		return "", appErrors.NewInternalServerError(err, "Error al generar token")
	}

	return tokenString, nil
}
