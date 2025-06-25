package service

import (
	"fmt"

	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/repository"
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
	existingUser, _ := s.userRepo.GetByEmail(user.Email)
	if existingUser.ID != 0 {
		return fmt.Errorf("usuario con email %s ya existe", user.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepo.Create(user)
}

func (s *UserService) Update(user models.User) (models.User, error) {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
