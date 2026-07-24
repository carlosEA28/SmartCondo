package services

import (
	"context"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/repositories"
)

type InadimplenteService struct {
	pagamentoRepo repositories.PagamentoRepository
}

func NewInadimplenteService(pagamentoRepo repositories.PagamentoRepository) *InadimplenteService {
	return &InadimplenteService{pagamentoRepo: pagamentoRepo}
}

func (s *InadimplenteService) ListInadimplentes(ctx context.Context) ([]dto.InadimplenteResponseDTO, error) {
	payments, err := s.pagamentoRepo.FindInadimplentes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list inadimplentes: %w", err)
	}

	grouped := make(map[string]*dto.InadimplenteResponseDTO)
	order := make([]string, 0)

	for i := range payments {
		p := &payments[i]
		moradorID := p.MoradorID.String()

		entry, exists := grouped[moradorID]
		if !exists {
			apartment := p.Morador.Apartment
			var aptDTO *dto.ApartmentResponseDTO
			if apartment != nil {
				aptDTO = &dto.ApartmentResponseDTO{
					ID:     apartment.ID,
					Number: apartment.Number,
					Block:  apartment.Block,
				}
			}

			entry = &dto.InadimplenteResponseDTO{
				Morador: dto.UserResponseDTO{
					ID:          p.Morador.ID,
					FullName:    p.Morador.FullName,
					Email:       p.Morador.Email,
					Phone:       p.Morador.Phone,
					Status:      string(p.Morador.Status),
					Role:        string(p.Morador.Role),
					Responsible: p.Morador.Responsible,
					Apartment:   aptDTO,
				},
				Payments: make([]dto.PagamentoResumoDTO, 0),
			}
			grouped[moradorID] = entry
			order = append(order, moradorID)
		}

		entry.TotalOverdue += p.Valor
		entry.Payments = append(entry.Payments, dto.PagamentoResumoDTO{
			ID:        p.ID,
			Valor:     p.Valor,
			Vencimento: p.Vencimento,
			Status:    string(p.Status),
		})
	}

	result := make([]dto.InadimplenteResponseDTO, 0, len(order))
	for _, id := range order {
		result = append(result, *grouped[id])
	}

	return result, nil
}
