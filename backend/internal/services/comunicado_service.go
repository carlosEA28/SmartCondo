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
)

type ComunicadoService struct {
	comunicadoRepo repositories.ComunicadoRepository
	userRepo       repositories.UserRepository
}

func NewComunicadoService(comunicadoRepo repositories.ComunicadoRepository, userRepo repositories.UserRepository) *ComunicadoService {
	return &ComunicadoService{
		comunicadoRepo: comunicadoRepo,
		userRepo:       userRepo,
	}
}

func (s *ComunicadoService) PublishComunicado(ctx context.Context, sindicoID uuid.UUID, input *dto.CreateComunicadoDTO) (*dto.ComunicadoResponseDTO, error) {
	if _, err := s.userRepo.FindByID(ctx, sindicoID); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrUnauthorizedSindico
		}
		return nil, fmt.Errorf("find sindico: %w", err)
	}

	comunicado := &models.Comunicado{
		Titulo:         input.Titulo,
		Descricao:      input.Descricao,
		DataPublicacao: time.Now(),
		SindicoID:      sindicoID,
	}

	if err := s.comunicadoRepo.Create(ctx, comunicado); err != nil {
		return nil, fmt.Errorf("publish comunicado: %w", err)
	}

	comunicado, err := s.comunicadoRepo.FindByID(ctx, comunicado.ID)
	if err != nil {
		return nil, fmt.Errorf("find created comunicado: %w", err)
	}

	return comunicadoToResponse(comunicado), nil
}

func (s *ComunicadoService) ListComunicados(ctx context.Context) ([]dto.ComunicadoResponseDTO, error) {
	comunicados, err := s.comunicadoRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list comunicados: %w", err)
	}

	response := make([]dto.ComunicadoResponseDTO, 0, len(comunicados))
	for i := range comunicados {
		response = append(response, *comunicadoToResponse(&comunicados[i]))
	}

	return response, nil
}

func (s *ComunicadoService) GetComunicado(ctx context.Context, id uuid.UUID) (*dto.ComunicadoResponseDTO, error) {
	comunicado, err := s.comunicadoRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrComunicadoNotFound) {
			return nil, apperrors.ErrComunicadoNotFound
		}
		return nil, fmt.Errorf("get comunicado: %w", err)
	}

	return comunicadoToResponse(comunicado), nil
}

func (s *ComunicadoService) DeleteComunicado(ctx context.Context, id uuid.UUID, sindicoID uuid.UUID) error {
	if err := s.comunicadoRepo.Delete(ctx, id, sindicoID); err != nil {
		if errors.Is(err, apperrors.ErrComunicadoNotFound) {
			return apperrors.ErrComunicadoNotFound
		}
		if errors.Is(err, apperrors.ErrComunicadoNotOwner) {
			return apperrors.ErrComunicadoNotOwner
		}
		return fmt.Errorf("delete comunicado: %w", err)
	}

	return nil
}

func comunicadoToResponse(c *models.Comunicado) *dto.ComunicadoResponseDTO {
	sindicoNome := ""
	if c.Sindico.FullName != "" {
		sindicoNome = c.Sindico.FullName
	}

	return &dto.ComunicadoResponseDTO{
		ID:             c.ID,
		Titulo:         c.Titulo,
		Descricao:      c.Descricao,
		DataPublicacao: c.DataPublicacao,
		SindicoID:      c.SindicoID,
		SindicoNome:    sindicoNome,
	}
}
