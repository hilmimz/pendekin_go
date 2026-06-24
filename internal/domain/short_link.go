package domain

import (
	"errors"
	"pendekin_go/pkg/errs"
	"time"
)

var (
	ErrAliasNotFound = errors.New("alias not found")
)

type ShortLink struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Alias       *string   `json:"alias"`
}

type CreateShortLinkRequest struct {
	OriginalURL string  `json:"original_url" binding:"required,url"`
	UserID      int     `json:"user_id"`
	ExpiresIn   *int    `json:"expires_in" binding:"omitempty,min=1"`
	Alias       *string `json:"alias" binding:"omitempty,min=4,max=20,alphanum"`
}

type CreateShortLinkResponse struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	Alias       *string   `json:"alias"`
	ShortUrl    string    `json:"short_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RedirectShortLinkRequest struct {
	Alias *string `json:"alias" binding:"required"`
}

type RedirectShortLinkResponse struct {
	OriginalURL string `json:"original_url"`
}

type ShortLinkRepository interface {
	FindByAlias(alias *string) (*ShortLink, error)
	Create(shortLink *ShortLink) error
	UpdateClickCount(shortLink *ShortLink) error
}

type ShortLinkUsecase interface {
	CreateShortLink(req *CreateShortLinkRequest) (*CreateShortLinkResponse, *errs.Error)
	RedirectShortLink(req *RedirectShortLinkRequest) (*RedirectShortLinkResponse, *errs.Error)
	// DeleteShortLink(shortLinkID int, userID int) *errors.Error
	// GetShortLinkStats(shortLinkID int) (*ShortLink, *errors.Error)
}
