package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateComunicadoDTO struct {
	Titulo    string `json:"titulo" binding:"required,max=100"`
	Descricao string `json:"descricao" binding:"required"`
}

type ComunicadoResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	Titulo         string    `json:"titulo"`
	Descricao      string    `json:"descricao"`
	DataPublicacao time.Time `json:"dataPublicacao"`
	SindicoID      uuid.UUID `json:"sindicoId"`
	SindicoNome    string    `json:"sindicoNome"`
}
