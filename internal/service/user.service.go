package service

import (
	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/repository"
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
	return s.userRepo.Create(user)
}

func (s *UserService) Update(user models.User) (models.User, error) {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
