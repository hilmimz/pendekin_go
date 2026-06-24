package usecase

import (
	"fmt"
	"math/rand"
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/errors"
	"time"
)

type ShortLinkUseCase struct {
	shortLinkRepo domain.ShortLinkRepository
	cfg           *config.AppConfig
}

func NewShortLinkUseCase(shortLinkRepo domain.ShortLinkRepository, cfg *config.AppConfig) *ShortLinkUseCase {
	return &ShortLinkUseCase{
		shortLinkRepo: shortLinkRepo,
		cfg:           cfg,
	}
}

func generateRandomAlias(length int) *string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	alias := string(b)
	return &alias
}

func (s *ShortLinkUseCase) CreateShortLink(req *domain.CreateShortLinkRequest) (*domain.CreateShortLinkResponse, *errors.Error) {
	var alias *string
	var expiresIn int

	if req.Alias == nil {
		alias = generateRandomAlias(int(s.cfg.AliasLength))
	} else {
		alias = req.Alias
		isAliasExists, err := s.shortLinkRepo.FindByAlias(alias)
		if err != nil {
			err := errors.Internal("failed to find short link by alias", err)
			return nil, err
		}
		if isAliasExists {
			err := errors.Conflict("alias already exist", err)
			return nil, err
		}
	}

	if req.ExpiresIn == nil {
		expiresIn = s.cfg.ExpiresIn
	} else {
		expiresIn = *req.ExpiresIn
	}

	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Hour)
	shortLink := &domain.ShortLink{
		OriginalURL: req.OriginalURL,
		ClickCount:  0,
		UserID:      1,
		ExpiresAt:   expiresAt,
		Alias:       alias,
	}
	if err := s.shortLinkRepo.Create(shortLink); err != nil {
		err := errors.Internal("failed to create short link: ", err)
		return nil, err
	}

	res := &domain.CreateShortLinkResponse{
		ShortLink: fmt.Sprintf(s.cfg.AppName+"/%s", *shortLink.Alias),
		ExpiresAt: shortLink.ExpiresAt,
	}
	return res, nil
}
