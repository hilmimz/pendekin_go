package domain

import (
	"errors"
	"pendekin_go/pkg/errs"
	"time"
)

var (
	ErrAliasNotFound = errors.New("alias not found")
)

type ShortUrl struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      *int      `json:"user_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Alias       *string   `json:"alias"`
}

type CreateShortUrlRequest struct {
	OriginalURL string  `json:"original_url" binding:"required,url"`
	UserID      *int    `json:"user_id"`
	ExpiresIn   *int    `json:"expires_in" binding:"omitempty,min=1"`
	Alias       *string `json:"alias" binding:"omitempty,min=4,max=20,alphanum"`
}

type CreateShortUrlResponse struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	Alias       *string   `json:"alias"`
	ShortUrl    string    `json:"short_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RedirectShortUrlRequest struct {
	Alias     *string `json:"alias" binding:"required"`
	IPAddress string  `json:"ip_address"`
	Referer   string  `json:"referer"`
	UserAgent string  `json:"user_agent"`
}

type RedirectShortUrlResponse struct {
	OriginalURL string `json:"original_url"`
}

type DeleteShortUrlRequest struct {
	ID     int  `json:"id"`
	UserID *int `json:"user_id"`
}

type DeleteShortUrlResponse struct {
	ShortUrlID int     `json:"short_url_id"`
	Alias      *string `json:"alias"`
}

type ShortUrlRepository interface {
	FindByAlias(alias *string) (*ShortUrl, error)
	FindById(id int) (*ShortUrl, error)
	Create(shortUrl *ShortUrl) error
	Delete(shortUrl *ShortUrl) error
	UpdateClickCount(shortUrl *ShortUrl) error
}

type ShortUrlUsecase interface {
	CreateShortUrl(req *CreateShortUrlRequest) (*CreateShortUrlResponse, *errs.Error)
	RedirectShortUrl(req *RedirectShortUrlRequest) (*RedirectShortUrlResponse, *errs.Error)
	DeleteShortUrl(req *DeleteShortUrlRequest) (*DeleteShortUrlResponse, *errs.Error)
	// GetShortUrlStats(shortUrlID int) (*ShortUrl, *errors.Error)
}
