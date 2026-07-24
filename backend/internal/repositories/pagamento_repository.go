package repositories

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PagamentoRepository interface {
	FindInadimplentes(ctx context.Context) ([]models.Pagamento, error)
	HasOverduePayments(ctx context.Context, moradorID uuid.UUID) (bool, error)
}

type GormPagamentoRepository struct {
	db *gorm.DB
}

func NewGormPagamentoRepository(db *gorm.DB) *GormPagamentoRepository {
	return &GormPagamentoRepository{db: db}
}

func (r *GormPagamentoRepository) HasOverduePayments(ctx context.Context, moradorID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Pagamento{}).
		Where("morador_id = ? AND status = ?", moradorID, models.PaymentStatusOverdue).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("has overdue payments: %w", err)
	}
	return count > 0, nil
}

func (r *GormPagamentoRepository) FindInadimplentes(ctx context.Context) ([]models.Pagamento, error) {
	var payments []models.Pagamento
	if err := r.db.WithContext(ctx).
		Preload("Morador.Apartment").
		Where("status = ?", models.PaymentStatusOverdue).
		Order("morador_id, vencimento ASC").
		Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("find inadimplentes: %w", err)
	}

	return payments, nil
}
