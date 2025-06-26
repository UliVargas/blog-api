package dto

import "github.com/UliVargas/blog-go/internal/domain/model"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) ToUser() model.User {
	return model.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
