package repository

import (
	"pendekin_go/internal/domain"

	"gorm.io/gorm"
)

type ClickLogRepository struct {
	db *gorm.DB
}

func NewClickLogRepository(db *gorm.DB) *ClickLogRepository {
	return &ClickLogRepository{
		db: db,
	}
}

func (c *ClickLogRepository) Create(clickLog *domain.ClickLog) error {
	if err := c.db.Create(clickLog).Error; err != nil {
		return err
	}
	return nil
}
