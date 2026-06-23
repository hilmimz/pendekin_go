package repository

import (
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

func (s *ShortLinkRepository) FindByAlias(alias *string) (bool, error) {
	var count int64
	if err := s.db.Model(&domain.ShortLink{}).Where("alias = ?", alias).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *ShortLinkRepository) Create(shortLink *domain.ShortLink) error {
	if err := s.db.Create(shortLink).Error; err != nil {
		return err
	}
	return nil
}
