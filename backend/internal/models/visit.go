package models

import (
	"time"

	"github.com/google/uuid"
)

type Visit struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	DataEntrada time.Time  `gorm:"column:dataentrada;not null"`
	DataSaida   *time.Time `gorm:"column:datasaida"`
	PorteiroID  uuid.UUID  `gorm:"column:porteiro_id;type:uuid;not null"`
	VisitanteID uuid.UUID  `gorm:"column:visitante_id;type:uuid;not null"`
	MoradorID   *uuid.UUID `gorm:"column:morador_id;type:uuid"`
	Porteiro    User       `gorm:"foreignKey:PorteiroID"`
	Visitante   Visitor    `gorm:"foreignKey:VisitanteID"`
	Morador     *User      `gorm:"foreignKey:MoradorID"`
}

func (Visit) TableName() string {
	return "visita"
}
