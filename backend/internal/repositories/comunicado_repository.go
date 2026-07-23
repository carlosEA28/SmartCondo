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

type ComunicadoRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Comunicado, error)
	FindAll(ctx context.Context) ([]models.Comunicado, error)
	Create(ctx context.Context, comunicado *models.Comunicado) error
	Delete(ctx context.Context, id uuid.UUID, sindicoID uuid.UUID) error
}

type GormComunicadoRepository struct {
	db *gorm.DB
}

func NewGormComunicadoRepository(db *gorm.DB) *GormComunicadoRepository {
	return &GormComunicadoRepository{db: db}
}

func (r *GormComunicadoRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Comunicado, error) {
	var comunicado models.Comunicado
	if err := r.db.WithContext(ctx).Preload("Sindico").First(&comunicado, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrComunicadoNotFound
		}
		return nil, fmt.Errorf("find comunicado by id: %w", err)
	}

	return &comunicado, nil
}

func (r *GormComunicadoRepository) FindAll(ctx context.Context) ([]models.Comunicado, error) {
	comunicados := make([]models.Comunicado, 0)
	if err := r.db.WithContext(ctx).Preload("Sindico").Order("datapublicacao DESC").Find(&comunicados).Error; err != nil {
		return nil, fmt.Errorf("list comunicados: %w", err)
	}

	return comunicados, nil
}

func (r *GormComunicadoRepository) Create(ctx context.Context, comunicado *models.Comunicado) error {
	if err := r.db.WithContext(ctx).Create(comunicado).Error; err != nil {
		return fmt.Errorf("create comunicado: %w", err)
	}

	return nil
}

func (r *GormComunicadoRepository) Delete(ctx context.Context, id uuid.UUID, sindicoID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND sindico_id = ?", id, sindicoID).
		Delete(&models.Comunicado{})
	if result.Error != nil {
		return fmt.Errorf("delete comunicado: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		exists, err := r.existsByID(ctx, id)
		if err != nil {
			return fmt.Errorf("check comunicado existence: %w", err)
		}
		if !exists {
			return apperrors.ErrComunicadoNotFound
		}
		return apperrors.ErrComunicadoNotOwner
	}

	return nil
}

func (r *GormComunicadoRepository) existsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Comunicado{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count comunicado: %w", err)
	}
	return count > 0, nil
}
