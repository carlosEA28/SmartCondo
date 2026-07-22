package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateVisitorDTO struct {
	Name  string                `form:"name" binding:"required,max=100"`
	CPF   string                `form:"cpf" binding:"required,len=11"`
	Phone string                `form:"phone" binding:"required,max=15"`
	Photo *multipart.FileHeader `form:"photo"`
}

type VisitorResponseDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	CPF   string    `json:"cpf"`
	Phone string    `json:"phone"`
	Photo string    `json:"photo,omitempty"`
}
