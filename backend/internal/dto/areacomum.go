package dto

import "github.com/google/uuid"

type AreaComumResponse struct {
	ID          uuid.UUID `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao"`
	Capacidade  int       `json:"capacidade"`
}
