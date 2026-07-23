package models

import (
	"time"

	"github.com/google/uuid"
)

type Comunicado struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Titulo         string    `gorm:"column:titulo;size:100;not null"`
	Descricao      string    `gorm:"column:descricao;not null"`
	DataPublicacao time.Time `gorm:"column:datapublicacao;not null"`
	SindicoID      uuid.UUID `gorm:"column:sindico_id;type:uuid;not null"`
	Sindico        User      `gorm:"foreignKey:SindicoID"`
}

func (Comunicado) TableName() string {
	return "comunicado"
}
