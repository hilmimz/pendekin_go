package domain

import (
	"errors"
	"pendekin_go/pkg/errs"
	"time"
)

var (
	ErrEmailNotFound = errors.New("email not found")
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRegisterResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

type UserUsecase interface {
	Register(req *UserRegisterRequest) (*UserRegisterResponse, *errs.Error)
}
