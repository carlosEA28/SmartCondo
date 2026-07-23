package dto

import "github.com/google/uuid"

type VisitorFilterDTO struct {
	Nome     string `form:"nome"`
	CPF      string `form:"cpf"`
	Telefone string `form:"telefone"`
	Liberado *bool  `form:"liberado"`
}

type ReleaseRequestDTO struct {
	PorteiroID uuid.UUID  `json:"porteiro_id" binding:"required"`
	MoradorID  *uuid.UUID `json:"morador_id"`
}
