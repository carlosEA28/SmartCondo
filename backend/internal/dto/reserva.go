package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateReservaRequest struct {
	AreaComumID uuid.UUID `json:"areacomum_id" binding:"required"`
	Data        string    `json:"data" binding:"required"`        // "2006-01-02"
	HoraInicio  string    `json:"horaInicio" binding:"required"`  // "15:04"
	HoraFim     string    `json:"horaFim" binding:"required"`     // "15:04"
}

type ReservaResponse struct {
	ID          uuid.UUID          `json:"id"`
	Data        time.Time          `json:"data"`
	HoraInicio  string             `json:"horaInicio"`
	HoraFim     string             `json:"horaFim"`
	Status      string             `json:"status"`
	MoradorID   uuid.UUID          `json:"morador_id"`
	AreaComum   AreaComumResponse  `json:"area_comum"`
}
