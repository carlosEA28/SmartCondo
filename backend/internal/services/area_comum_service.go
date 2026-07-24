package services

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/repositories"
)

type AreaComumService struct {
	areaComumRepo repositories.AreaComumRepository
}

func NewAreaComumService(areaComumRepo repositories.AreaComumRepository) *AreaComumService {
	return &AreaComumService{areaComumRepo: areaComumRepo}
}

func (s *AreaComumService) ListAreas(ctx context.Context) ([]dto.AreaComumResponse, error) {
	areas, err := s.areaComumRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list areas: %w", err)
	}

	result := make([]dto.AreaComumResponse, 0, len(areas))
	for _, a := range areas {
		result = append(result, dto.AreaComumResponse{
			ID:         a.ID,
			Nome:       a.Nome,
			Descricao:  a.Descricao,
			Capacidade: a.Capacidade,
		})
	}

	return result, nil
}
