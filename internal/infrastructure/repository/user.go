package repository

import (
	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/UliVargas/blog-go/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}
	return users, nil
}

func (r *UserRepository) GetByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return model.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return model.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) Create(user model.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return errors.WrapDatabaseError(err)
	}
	return nil
}

func (r *UserRepository) Update(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return model.User{}, errors.WrapDatabaseError(err)
	}
	return user, nil
}

func (r *UserRepository) Delete(id uint) error {
	err := r.db.Delete(&model.User{}, id).Error
	if err != nil {
		return errors.WrapDatabaseError(err)
	}
	return nil
}
