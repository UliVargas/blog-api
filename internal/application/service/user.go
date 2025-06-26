package service

import (
	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/UliVargas/blog-go/internal/infrastructure/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetByID(id uint) (model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(user model.User) (model.User, error) {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
