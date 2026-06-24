package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/errs"
	"time"
)

type ShortUrlUsecase struct {
	shortUrlRepo domain.ShortUrlRepository
	clickLogRepo domain.ClickLogRepository
	cfg          *config.AppConfig
}

func NewShortUrlUsecase(shortUrlRepo domain.ShortUrlRepository, clickLogRepo domain.ClickLogRepository, cfg *config.AppConfig) *ShortUrlUsecase {
	return &ShortUrlUsecase{
		shortUrlRepo: shortUrlRepo,
		clickLogRepo: clickLogRepo,
		cfg:          cfg,
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

func (s *ShortUrlUsecase) CreateShortUrl(req *domain.CreateShortUrlRequest) (*domain.CreateShortUrlResponse, *errs.Error) {
	var alias *string
	var expiresIn int

	if req.Alias == nil {
		alias = generateRandomAlias(int(s.cfg.AliasLength))
	} else {
		alias = req.Alias
		shortUrl, err := s.shortUrlRepo.FindByAlias(alias)

		if shortUrl != nil {
			return nil, errs.Conflict("alias already exist", err)
		}

		if err != nil && !errors.Is(err, domain.ErrAliasNotFound) {
			return nil, errs.Internal("failed to find short link by alias", err)
		}

	}

	if req.ExpiresIn == nil {
		expiresIn = s.cfg.ExpiresIn
	} else {
		expiresIn = *req.ExpiresIn
	}

	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Hour)
	shortLink := &domain.ShortUrl{
		OriginalURL: req.OriginalURL,
		ClickCount:  0,
		UserID:      1,
		ExpiresAt:   expiresAt,
		Alias:       alias,
	}
	if err := s.shortUrlRepo.Create(shortLink); err != nil {
		err := errs.Internal("failed to create short url", err)
		return nil, err
	}

	res := &domain.CreateShortUrlResponse{
		ID:          shortLink.ID,
		OriginalURL: shortLink.OriginalURL,
		Alias:       shortLink.Alias,
		ShortUrl:    fmt.Sprint("http://" + s.cfg.AppName + "/" + *shortLink.Alias),
		ExpiresAt:   shortLink.ExpiresAt,
		CreatedAt:   shortLink.CreatedAt,
	}
	return res, nil
}

func (s *ShortUrlUsecase) RedirectShortUrl(req *domain.RedirectShortUrlRequest) (*domain.RedirectShortUrlResponse, *errs.Error) {
	shortUrl, err := s.shortUrlRepo.FindByAlias(req.Alias)
	if err != nil && !errors.Is(err, domain.ErrAliasNotFound) {
		err := errs.Internal("failed to find short url by alias", err)
		return nil, err
	}

	if shortUrl == nil {
		err := errs.NotFound("short url not found", err)
		return nil, err
	}

	res := &domain.RedirectShortUrlResponse{
		OriginalURL: shortUrl.OriginalURL,
	}

	if err := s.shortUrlRepo.UpdateClickCount(shortUrl); err != nil {
		err := errs.Internal("failed to update click count", err)
		return nil, err
	}

	if err := s.clickLogRepo.Create(&domain.ClickLog{
		ClickedAt:  time.Now(),
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		Referer:    req.Referer,
		ShortURLID: shortUrl.ID,
	}); err != nil {
		err := errs.Internal("failed to create click log", err)
		return nil, err
	}

	return res, nil
}
