package repository

import (
	"errors"
	"pendekin_go/internal/domain"

	"gorm.io/gorm"
)

type ShortUrlRepository struct {
	db *gorm.DB
}

func NewShortUrlRepository(db *gorm.DB) *ShortUrlRepository {
	return &ShortUrlRepository{
		db: db,
	}
}

func (s *ShortUrlRepository) FindByAlias(alias *string) (*domain.ShortUrl, error) {
	var shortLink domain.ShortUrl
	err := s.db.Where("alias = ?", alias).First(&shortLink).Error // langsung ambil error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAliasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (s *ShortUrlRepository) Create(shortLink *domain.ShortUrl) error {
	if err := s.db.Create(shortLink).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortUrlRepository) UpdateClickCount(shortLink *domain.ShortUrl) error {
	if err := s.db.Model(shortLink).Update("click_count", shortLink.ClickCount+1).Error; err != nil {
		return err
	}
	return nil
}
