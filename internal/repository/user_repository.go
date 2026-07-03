package repository

import (
	"errors"
	"pendekin_go/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := ur.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEmailNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Create(user *domain.User) error {
	err := ur.db.Create(user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (ur *UserRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEmailNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
