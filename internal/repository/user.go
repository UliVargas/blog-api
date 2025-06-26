package repository

import (
	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/pkg/errors"
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
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}
	return users, nil
}

func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return models.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) Create(user models.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return errors.WrapDatabaseError(err)
	}
	return nil
}

func (r *UserRepository) Update(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return models.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) Delete(id uint) error {
	err := r.db.Delete(&models.User{}, id).Error
	if err != nil {
		return errors.WrapDatabaseError(err)
	}
	return nil
}
