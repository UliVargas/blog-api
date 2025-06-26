package service

import "github.com/UliVargas/blog-go/internal/domain/model"

// UserServiceInterface define el contrato para las operaciones del servicio de usuarios
// Esta interfaz pertenece a la capa de dominio ya que define el contrato
// que el dominio espera de la capa de aplicación
type UserServiceInterface interface {
	GetAll() ([]model.User, error)
	GetByID(id uint) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id uint) error
}

// AuthServiceInterface define el contrato para las operaciones del servicio de autenticación
// Esta interfaz pertenece a la capa de dominio ya que define el contrato
// que el dominio espera de la capa de aplicación
type AuthServiceInterface interface {
	Login(email, password string) (string, error)
	Register(user model.User) error
}