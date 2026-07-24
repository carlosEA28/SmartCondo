package repositories

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/models"
	"gorm.io/gorm"
)

type NotificacaoRepository interface {
	Create(ctx context.Context, notificacao *models.Notificacao) error
	FindByStatus(ctx context.Context, status models.NotificacaoStatus) ([]models.Notificacao, error)
}

type GormNotificacaoRepository struct {
	db *gorm.DB
}

func NewGormNotificacaoRepository(db *gorm.DB) *GormNotificacaoRepository {
	return &GormNotificacaoRepository{db: db}
}

func (r *GormNotificacaoRepository) Create(ctx context.Context, notificacao *models.Notificacao) error {
	if err := r.db.WithContext(ctx).Create(notificacao).Error; err != nil {
		return fmt.Errorf("create notificacao: %w", err)
	}

	return nil
}

func (r *GormNotificacaoRepository) FindByStatus(ctx context.Context, status models.NotificacaoStatus) ([]models.Notificacao, error) {
	var notificacoes []models.Notificacao
	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("dataenvio DESC").
		Find(&notificacoes).Error; err != nil {
		return nil, fmt.Errorf("find notificacoes by status: %w", err)
	}

	return notificacoes, nil
}


