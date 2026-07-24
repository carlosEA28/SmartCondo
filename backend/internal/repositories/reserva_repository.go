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

type ReservaRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Reserva, error)
	FindByMorador(ctx context.Context, moradorID uuid.UUID) ([]models.Reserva, error)
	FindConflicting(ctx context.Context, areaID uuid.UUID, data string, horaInicio, horaFim string, excludeReservaID *uuid.UUID) (bool, error)
	Create(ctx context.Context, reserva *models.Reserva) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ReservaStatus) error
}

type GormReservaRepository struct {
	db *gorm.DB
}

func NewGormReservaRepository(db *gorm.DB) *GormReservaRepository {
	return &GormReservaRepository{db: db}
}

func (r *GormReservaRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Reserva, error) {
	var reserva models.Reserva
	if err := r.db.WithContext(ctx).Preload("AreaComum").First(&reserva, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrReservaNotFound
		}
		return nil, fmt.Errorf("find reserva by id: %w", err)
	}

	return &reserva, nil
}

func (r *GormReservaRepository) FindByMorador(ctx context.Context, moradorID uuid.UUID) ([]models.Reserva, error) {
	reservas := make([]models.Reserva, 0)
	if err := r.db.WithContext(ctx).
		Preload("AreaComum").
		Where("morador_id = ?", moradorID).
		Order("data DESC, horainicio ASC").
		Find(&reservas).Error; err != nil {
		return nil, fmt.Errorf("find reservas by morador: %w", err)
	}

	return reservas, nil
}

func (r *GormReservaRepository) FindConflicting(ctx context.Context, areaID uuid.UUID, data string, horaInicio, horaFim string, excludeReservaID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Reserva{}).
		Where("areacomum_id = ? AND data = ? AND status != ?", areaID, data, models.ReservaStatusCancelada).
		Where("horainicio < ? AND horafim > ?", horaFim, horaInicio)

	if excludeReservaID != nil {
		query = query.Where("id != ?", *excludeReservaID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("find conflicting reservas: %w", err)
	}

	return count > 0, nil
}

func (r *GormReservaRepository) Create(ctx context.Context, reserva *models.Reserva) error {
	if err := r.db.WithContext(ctx).Create(reserva).Error; err != nil {
		return fmt.Errorf("create reserva: %w", err)
	}

	return nil
}

func (r *GormReservaRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ReservaStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Reserva{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update reserva status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrReservaNotFound
	}

	return nil
}
