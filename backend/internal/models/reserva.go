package models

import (
	"time"

	"github.com/google/uuid"
)

type ReservaStatus string

const (
	ReservaStatusPendente   ReservaStatus = "PENDENTE"
	ReservaStatusConfirmada ReservaStatus = "CONFIRMADA"
	ReservaStatusCancelada  ReservaStatus = "CANCELADA"
)

type Reserva struct {
	ID          uuid.UUID     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Data        time.Time     `gorm:"column:data;type:date;not null"`
	HoraInicio  string        `gorm:"column:horainicio;type:time;not null"`
	HoraFim     string        `gorm:"column:horafim;type:time;not null"`
	Status      ReservaStatus `gorm:"column:status;size:20;not null;default:PENDENTE"`
	MoradorID   uuid.UUID     `gorm:"column:morador_id;type:uuid;not null"`
	Morador     User          `gorm:"foreignKey:MoradorID"`
	AreaComumID uuid.UUID     `gorm:"column:areacomum_id;type:uuid;not null"`
	AreaComum   AreaComum     `gorm:"foreignKey:AreaComumID"`
	SindicoID   *uuid.UUID    `gorm:"column:sindico_id;type:uuid"`
	Sindico     *User         `gorm:"foreignKey:SindicoID"`
}

func (Reserva) TableName() string {
	return "reserva"
}
