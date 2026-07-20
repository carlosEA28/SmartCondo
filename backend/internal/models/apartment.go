package models

import "github.com/google/uuid"

type Apartment struct {
	ID     uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Number int       `gorm:"column:numero;not null"`
	Block  string    `gorm:"column:bloco;size:10;not null"`
}

func (Apartment) TableName() string {
	return "apartamento"
}
