package domain

import (
	"pendekin_go/pkg/errors"
	"time"
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

type ShortLinkRepository interface {
	FindByAlias(alias *string) (bool, error)
	Create(shortLink *ShortLink) error
}

type ShortLinkUsecase interface {
	CreateShortLink(req *CreateShortLinkRequest) (*CreateShortLinkResponse, *errors.Error)
	// DeleteShortLink(shortLinkID int, userID int) *errors.Error
	// GetShortLinkStats(shortLinkID int) (*ShortLink, *errors.Error)
}
