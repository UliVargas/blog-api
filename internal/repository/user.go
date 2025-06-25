package repository

import (
	"github.com/UliVargas/blog-go/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email =?", email).First(&user).Error
	return user, err
}

func (r *UserRepository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *UserRepository) Update(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *UserRepository) Delete(id uint) error {
	err := r.db.Delete(&models.User{}, id).Error
	return err
}
