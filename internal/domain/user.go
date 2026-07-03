package domain

import (
	"errors"
	"pendekin_go/pkg/errs"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type JWTClaims struct {
	UserID int    `json:"sub"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
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

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

type UserUsecase interface {
	Register(req *UserRegisterRequest) (*UserRegisterResponse, *errs.Error)
	Login(req *UserLoginRequest) (*UserLoginResponse, *errs.Error)
}
