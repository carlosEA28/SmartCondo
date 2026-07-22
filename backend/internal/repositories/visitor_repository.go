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

type VisitorRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Visitor, error)
	FindByCPF(ctx context.Context, cpf string) (*models.Visitor, error)
	List(ctx context.Context) ([]models.Visitor, error)
	Create(ctx context.Context, visitor *models.Visitor) error
}

type GormVisitorRepository struct {
	db *gorm.DB
}

func NewGormVisitorRepository(db *gorm.DB) *GormVisitorRepository {
	return &GormVisitorRepository{db: db}
}

func (r *GormVisitorRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Visitor, error) {
	var visitor models.Visitor
	if err := r.db.WithContext(ctx).First(&visitor, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrVisitorNotFound
		}
		return nil, fmt.Errorf("find visitor by id: %w", err)
	}

	return &visitor, nil
}

func (r *GormVisitorRepository) FindByCPF(ctx context.Context, cpf string) (*models.Visitor, error) {
	var visitor models.Visitor
	if err := r.db.WithContext(ctx).Where("cpf = ?", cpf).First(&visitor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrVisitorNotFound
		}
		return nil, fmt.Errorf("find visitor by cpf: %w", err)
	}

	return &visitor, nil
}

func (r *GormVisitorRepository) List(ctx context.Context) ([]models.Visitor, error) {
	visitors := make([]models.Visitor, 0)
	if err := r.db.WithContext(ctx).Order("nome ASC").Find(&visitors).Error; err != nil {
		return nil, fmt.Errorf("list visitors: %w", err)
	}

	return visitors, nil
}

func (r *GormVisitorRepository) Create(ctx context.Context, visitor *models.Visitor) error {
	if err := r.db.WithContext(ctx).Create(visitor).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrVisitorAlreadyExists
		}
		return fmt.Errorf("create visitor: %w", err)
	}

	return nil
}
