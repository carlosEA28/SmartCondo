package repositories

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/models"
	"gorm.io/gorm"
)

type PagamentoRepository interface {
	FindInadimplentes(ctx context.Context) ([]models.Pagamento, error)
}

type GormPagamentoRepository struct {
	db *gorm.DB
}

func NewGormPagamentoRepository(db *gorm.DB) *GormPagamentoRepository {
	return &GormPagamentoRepository{db: db}
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
