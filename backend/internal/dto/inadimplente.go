package dto

import (
	"time"

	"github.com/google/uuid"
)

type PagamentoResumoDTO struct {
	ID        uuid.UUID `json:"id"`
	Valor     float64   `json:"valor"`
	Vencimento time.Time `json:"vencimento"`
	Status    string    `json:"status"`
}

type InadimplenteResponseDTO struct {
	Morador      UserResponseDTO      `json:"morador"`
	TotalOverdue float64              `json:"total_overdue"`
	Payments     []PagamentoResumoDTO `json:"payments"`
}
