package jwt

import (
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SecretKey []byte
	ExpiresIn time.Duration
}

func NewJWT(cfg *config.AppConfig) *JWT {
	return &JWT{
		SecretKey: []byte(cfg.JWTSecret),
		ExpiresIn: time.Duration(cfg.JWTExpiresIn) * time.Second,
	}
}

func (j *JWT) GenerateToken(id int, email string, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.JWTClaims{
		UserID: id,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) VerifyToken(tokenString string) (*jwt.Token, error) {
	// Use ParseWithClaims for custom claims, for default jwt map claims use Parse instead
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (j *JWT) GetClaims(tokenString string) (*domain.JWTClaims, error) {
	token, err := j.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
