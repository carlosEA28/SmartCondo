package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AreaComumRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.AreaComum, error)
	FindAll(ctx context.Context) ([]models.AreaComum, error)
}

type GormAreaComumRepository struct {
	db *gorm.DB
}

func NewGormAreaComumRepository(db *gorm.DB) *GormAreaComumRepository {
	return &GormAreaComumRepository{db: db}
}

func (r *GormAreaComumRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AreaComum, error) {
	var area models.AreaComum
	if err := r.db.WithContext(ctx).First(&area, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrAreaComumNotFound
		}
		return nil, fmt.Errorf("find area comum by id: %w", err)
	}

	return &area, nil
}

func (r *GormAreaComumRepository) FindAll(ctx context.Context) ([]models.AreaComum, error) {
	areas := make([]models.AreaComum, 0)
	if err := r.db.WithContext(ctx).Order("nome ASC").Find(&areas).Error; err != nil {
		return nil, fmt.Errorf("list areas comuns: %w", err)
	}

	return areas, nil
}
