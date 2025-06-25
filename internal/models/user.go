package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name" validate:"required,min=3,max=50"`
	Email     string    `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=6"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (req *CreateUserRequest) ToUser() User {
	return User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
