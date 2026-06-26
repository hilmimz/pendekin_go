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
	err := s.db.Where("alias = ?", alias).First(&shortLink).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAliasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (s *ShortUrlRepository) FindById(id int) (*domain.ShortUrl, error) {
	var shortUrl domain.ShortUrl
	err := s.db.Where("id = ?", id).First(&shortUrl).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAliasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &shortUrl, nil
}

func (s *ShortUrlRepository) Create(shortLink *domain.ShortUrl) error {
	if err := s.db.Create(shortLink).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortUrlRepository) Delete(shortLink *domain.ShortUrl) error {
	if err := s.db.Delete(shortLink).Error; err != nil {
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
