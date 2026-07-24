package services

import (
	"context"
	"fmt"
	"time"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/google/uuid"
)

type ReservaService struct {
	reservaRepo     repositories.ReservaRepository
	areaComumRepo   repositories.AreaComumRepository
	pagamentoRepo   repositories.PagamentoRepository
}

func NewReservaService(
	reservaRepo repositories.ReservaRepository,
	areaComumRepo repositories.AreaComumRepository,
	pagamentoRepo repositories.PagamentoRepository,
) *ReservaService {
	return &ReservaService{
		reservaRepo:   reservaRepo,
		areaComumRepo: areaComumRepo,
		pagamentoRepo: pagamentoRepo,
	}
}

func (s *ReservaService) CreateReserva(ctx context.Context, moradorID uuid.UUID, req dto.CreateReservaRequest) (*dto.ReservaResponse, error) {
	area, err := s.areaComumRepo.FindByID(ctx, req.AreaComumID)
	if err != nil {
		return nil, fmt.Errorf("create reserva: %w", err)
	}

	overdue, err := s.pagamentoRepo.HasOverduePayments(ctx, moradorID)
	if err != nil {
		return nil, fmt.Errorf("create reserva: %w", err)
	}
	if overdue {
		return nil, apperrors.ErrMoradorInadimplente
	}

	conflict, err := s.reservaRepo.FindConflicting(ctx, req.AreaComumID, req.Data, req.HoraInicio, req.HoraFim, nil)
	if err != nil {
		return nil, fmt.Errorf("create reserva: %w", err)
	}
	if conflict {
		return nil, apperrors.ErrReservaConflito
	}

	dataParsed, err := time.Parse("2006-01-02", req.Data)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid date format", apperrors.ErrInvalidReservaData)
	}

	reserva := &models.Reserva{
		Data:        dataParsed,
		HoraInicio:  req.HoraInicio,
		HoraFim:     req.HoraFim,
		Status:      models.ReservaStatusPendente,
		MoradorID:   moradorID,
		AreaComumID: req.AreaComumID,
	}

	if err := s.reservaRepo.Create(ctx, reserva); err != nil {
		return nil, fmt.Errorf("create reserva: %w", err)
	}

	return &dto.ReservaResponse{
		ID:         reserva.ID,
		Data:       reserva.Data,
		HoraInicio: reserva.HoraInicio,
		HoraFim:    reserva.HoraFim,
		Status:     string(reserva.Status),
		MoradorID:  reserva.MoradorID,
		AreaComum: dto.AreaComumResponse{
			ID:         area.ID,
			Nome:       area.Nome,
			Descricao:  area.Descricao,
			Capacidade: area.Capacidade,
		},
	}, nil
}

func (s *ReservaService) ListMinhasReservas(ctx context.Context, moradorID uuid.UUID) ([]dto.ReservaResponse, error) {
	reservas, err := s.reservaRepo.FindByMorador(ctx, moradorID)
	if err != nil {
		return nil, fmt.Errorf("list minhas reservas: %w", err)
	}

	result := make([]dto.ReservaResponse, 0, len(reservas))
	for _, r := range reservas {
		result = append(result, dto.ReservaResponse{
			ID:         r.ID,
			Data:       r.Data,
			HoraInicio: r.HoraInicio,
			HoraFim:    r.HoraFim,
			Status:     string(r.Status),
			MoradorID:  r.MoradorID,
			AreaComum: dto.AreaComumResponse{
				ID:         r.AreaComum.ID,
				Nome:       r.AreaComum.Nome,
				Descricao:  r.AreaComum.Descricao,
				Capacidade: r.AreaComum.Capacidade,
			},
		})
	}

	return result, nil
}

func (s *ReservaService) CancelReserva(ctx context.Context, reservaID uuid.UUID, moradorID uuid.UUID) (*dto.ReservaResponse, error) {
	reserva, err := s.reservaRepo.FindByID(ctx, reservaID)
	if err != nil {
		return nil, fmt.Errorf("cancel reserva: %w", err)
	}

	if reserva.MoradorID != moradorID {
		return nil, apperrors.ErrReservaNotOwner
	}

	if err := s.reservaRepo.UpdateStatus(ctx, reservaID, models.ReservaStatusCancelada); err != nil {
		return nil, fmt.Errorf("cancel reserva: %w", err)
	}

	reserva.Status = models.ReservaStatusCancelada

	return &dto.ReservaResponse{
		ID:         reserva.ID,
		Data:       reserva.Data,
		HoraInicio: reserva.HoraInicio,
		HoraFim:    reserva.HoraFim,
		Status:     string(reserva.Status),
		MoradorID:  reserva.MoradorID,
		AreaComum: dto.AreaComumResponse{
			ID:         reserva.AreaComum.ID,
			Nome:       reserva.AreaComum.Nome,
			Descricao:  reserva.AreaComum.Descricao,
			Capacidade: reserva.AreaComum.Capacidade,
		},
	}, nil
}
