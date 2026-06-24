package repository

import (
	"errors"
	"pendekin_go/internal/domain"

	"gorm.io/gorm"
)

type ShortLinkRepository struct {
	db *gorm.DB
}

func NewShortLinkRepository(db *gorm.DB) *ShortLinkRepository {
	return &ShortLinkRepository{
		db: db,
	}
}

func (s *ShortLinkRepository) FindByAlias(alias *string) (*domain.ShortLink, error) {
	var shortLink domain.ShortLink
	err := s.db.Where("alias = ?", alias).First(&shortLink).Error // langsung ambil error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAliasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (s *ShortLinkRepository) Create(shortLink *domain.ShortLink) error {
	if err := s.db.Create(shortLink).Error; err != nil {
		return err
	}
	return nil
}

func (s *ShortLinkRepository) UpdateClickCount(shortLink *domain.ShortLink) error {
	if err := s.db.Model(shortLink).Update("click_count", shortLink.ClickCount+1).Error; err != nil {
		return err
	}
	return nil
}
