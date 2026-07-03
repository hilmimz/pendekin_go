package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/errs"
	"pendekin_go/pkg/logger"
	"time"

	"go.uber.org/zap"
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
			logger.Log.Warn("alias already exists",
				zap.String("alias", *alias),
			)
			return nil, errs.Conflict("alias already exist", err)
		}

		if err != nil && !errors.Is(err, domain.ErrAliasNotFound) {
			logger.Log.Error("failed to find short url by alias",
				zap.String("alias", *alias),
				zap.Error(err),
			)
			return nil, errs.Internal("failed to find short url by alias", err)
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
		UserID:      req.UserID,
		ExpiresAt:   expiresAt,
		Alias:       alias,
	}
	if err := s.shortUrlRepo.Create(shortLink); err != nil {
		logger.Log.Error("failed to create short url",
			zap.String("alias", *alias),
			zap.Error(err),
		)
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
	logger.Log.Info("short url created successfully",
		zap.String("alias", *alias),
		zap.String("original_url", shortLink.OriginalURL),
		zap.String("user_id", fmt.Sprint(shortLink.UserID)),
		zap.String("expires_at", shortLink.ExpiresAt.String()),
		zap.String("created_at", shortLink.CreatedAt.String()),
	)
	return res, nil
}

func (s *ShortUrlUsecase) RedirectShortUrl(req *domain.RedirectShortUrlRequest) (*domain.RedirectShortUrlResponse, *errs.Error) {
	shortUrl, err := s.shortUrlRepo.FindByAlias(req.Alias)
	if err != nil && !errors.Is(err, domain.ErrAliasNotFound) {
		logger.Log.Error("failed to find short url by alias",
			zap.String("alias", *req.Alias),
			zap.Error(err),
		)
		err := errs.Internal("failed to find short url by alias", err)
		return nil, err
	}

	if shortUrl == nil {
		logger.Log.Warn("short url not found",
			zap.String("alias", *req.Alias),
		)
		err := errs.NotFound("short url not found", err)
		return nil, err
	}

	res := &domain.RedirectShortUrlResponse{
		OriginalURL: shortUrl.OriginalURL,
	}

	if err := s.shortUrlRepo.UpdateClickCount(shortUrl); err != nil {
		logger.Log.Error("failed to update click count",
			zap.String("alias", *req.Alias),
			zap.Error(err),
		)
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
		logger.Log.Error("failed to create click log",
			zap.String("alias", *req.Alias),
			zap.Error(err),
		)
		err := errs.Internal("failed to create click log", err)
		return nil, err
	}

	logger.Log.Info("short url redirected successfully",
		zap.String("alias", *req.Alias),
		zap.String("original_url", shortUrl.OriginalURL),
		zap.String("ip_address", req.IPAddress),
		zap.String("short_url_id", fmt.Sprint(shortUrl.ID)),
	)

	return res, nil
}

func (s *ShortUrlUsecase) DeleteShortUrl(req *domain.DeleteShortUrlRequest) (*domain.DeleteShortUrlResponse, *errs.Error) {
	shortUrl, err := s.shortUrlRepo.FindById(req.ID)
	if err != nil && !errors.Is(err, domain.ErrAliasNotFound) {
		logger.Log.Error("failed to find short url by id",
			zap.Int("id", req.ID),
			zap.Error(err),
		)
		err := errs.Internal("failed to find short url by id", err)
		return nil, err
	}

	if shortUrl == nil {
		logger.Log.Warn("short url not found",
			zap.Int("id", req.ID),
		)
		err := errs.NotFound("short url not found", nil)
		return nil, err
	}

	if shortUrl.UserID != req.UserID {
		logger.Log.Warn("user does not own short url",
			zap.Int("id", req.ID),
			zap.Int("user_id", *req.UserID),
		)
		err := errs.Forbidden("user does not own short url", nil)
		return nil, err
	}

	if err := s.shortUrlRepo.Delete(shortUrl); err != nil {
		logger.Log.Error("failed to delete short url",
			zap.Int("id", req.ID),
			zap.Error(err),
		)
		err := errs.Internal("failed to delete short url", err)
		return nil, err
	}

	res := &domain.DeleteShortUrlResponse{
		ShortUrlID: shortUrl.ID,
		Alias:      shortUrl.Alias,
	}

	return res, nil
}
