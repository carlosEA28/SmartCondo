package repositories

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/models"
	"gorm.io/gorm"
)

type VisitRepository interface {
	Create(ctx context.Context, visit *models.Visit) error
}

type GormVisitRepository struct {
	db *gorm.DB
}

func NewGormVisitRepository(db *gorm.DB) *GormVisitRepository {
	return &GormVisitRepository{db: db}
}

func (r *GormVisitRepository) Create(ctx context.Context, visit *models.Visit) error {
	if err := r.db.WithContext(ctx).Create(visit).Error; err != nil {
		return fmt.Errorf("create visit: %w", err)
	}

	return nil
}
