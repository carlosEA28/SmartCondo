package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PorteiroService struct {
	db          *gorm.DB
	visitorRepo repositories.VisitorRepository
	visitRepo   repositories.VisitRepository
	userRepo    repositories.UserRepository
}

func NewPorteiroService(
	db *gorm.DB,
	visitorRepo repositories.VisitorRepository,
	visitRepo repositories.VisitRepository,
	userRepo repositories.UserRepository,
) *PorteiroService {
	return &PorteiroService{
		db:          db,
		visitorRepo: visitorRepo,
		visitRepo:   visitRepo,
		userRepo:    userRepo,
	}
}

func (s *PorteiroService) SearchVisitors(ctx context.Context, filter *dto.VisitorFilterDTO) ([]dto.VisitorResponseDTO, error) {
	if filter.Nome == "" && filter.CPF == "" && filter.Telefone == "" && filter.Liberado == nil {
		return nil, apperrors.ErrFilterRequired
	}

	visitors, err := s.visitorRepo.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("search visitors: %w", err)
	}

	response := make([]dto.VisitorResponseDTO, 0, len(visitors))
	for i := range visitors {
		response = append(response, *visitorToResponse(&visitors[i]))
	}

	return response, nil
}

func (s *PorteiroService) ReleaseVisitor(ctx context.Context, visitorID uuid.UUID, porteiroID uuid.UUID, moradorID *uuid.UUID) (*dto.VisitorResponseDTO, error) {
	visitor, err := s.visitorRepo.FindByID(ctx, visitorID)
	if err != nil {
		if errors.Is(err, apperrors.ErrVisitorNotFound) {
			return nil, apperrors.ErrVisitorNotFound
		}
		return nil, fmt.Errorf("find visitor: %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, porteiroID); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrPorteiroNotFound
		}
		return nil, fmt.Errorf("find porteiro: %w", err)
	}

	visit := &models.Visit{
		DataEntrada: time.Now(),
		PorteiroID:  porteiroID,
		VisitanteID: visitorID,
		MoradorID:   moradorID,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Visitor{}).
			Where("id = ?", visitorID).
			Update("liberado", true).Error; err != nil {
			return fmt.Errorf("update liberado: %w", err)
		}

		if err := tx.Create(visit).Error; err != nil {
			return fmt.Errorf("create visit: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	visitor.Liberado = true
	return visitorToResponse(visitor), nil
}
