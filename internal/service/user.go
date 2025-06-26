package service

import (
	"errors"

	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/repository"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetByID(id uint) (models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Create(user models.User) error {
	// Verificar si el usuario ya existe
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, appErrors.ErrUserNotFound) {
		// Error inesperado al consultar la base de datos
		return err
	}
	if existingUser.ID != 0 {
		// El usuario ya existe
		return appErrors.ErrEmailExists
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.NewInternalServerError(err, "Error al procesar la contraseña")
	}
	user.Password = string(hashedPassword)
	
	// Crear el usuario
	return s.userRepo.Create(user)
}

func (s *UserService) Update(user models.User) (models.User, error) {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
