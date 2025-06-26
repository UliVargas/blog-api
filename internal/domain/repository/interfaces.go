package repository

import "github.com/UliVargas/blog-go/internal/domain/model"

// UserRepositoryInterface define el contrato para las operaciones del repositorio de usuarios
// Esta interfaz pertenece a la capa de dominio ya que define el contrato
// que el dominio espera de la capa de infraestructura
type UserRepositoryInterface interface {
	GetAll() ([]model.User, error)
	GetByID(id uint) (model.User, error)
	GetByEmail(email string) (model.User, error)
	Create(user model.User) error
	Update(user model.User) (model.User, error)
	Delete(id uint) error
}