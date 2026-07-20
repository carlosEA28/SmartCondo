package dto

import "github.com/google/uuid"

type CreateApartmentDTO struct {
	Number int    `json:"number" binding:"required,gt=0"`
	Block  string `json:"block" binding:"required,max=10"`
}

type CreateUserDTO struct {
	FullName    string              `json:"full_name" binding:"required,max=100"`
	Email       string              `json:"email" binding:"required,email,max=100"`
	Password    string              `json:"password" binding:"required,min=8,max=72"`
	Phone       string              `json:"phone" binding:"required,max=15"`
	Responsible bool                `json:"responsible"`
	Role        string              `json:"role" binding:"required,oneof=MORADOR PORTEIRO SINDICO"`
	Apartment   *CreateApartmentDTO `json:"apartment"`
}

type ApartmentResponseDTO struct {
	ID     uuid.UUID `json:"id"`
	Number int       `json:"number"`
	Block  string    `json:"block"`
}

type UserResponseDTO struct {
	ID          uuid.UUID             `json:"id"`
	FullName    string                `json:"full_name"`
	Email       string                `json:"email"`
	Phone       string                `json:"phone"`
	Status      string                `json:"status"`
	Role        string                `json:"role"`
	Responsible bool                  `json:"responsible"`
	Apartment   *ApartmentResponseDTO `json:"apartment,omitempty"`
}
