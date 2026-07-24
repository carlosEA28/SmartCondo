package models

import "github.com/google/uuid"

type AreaComum struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Nome        string    `gorm:"column:nome;size:100;not null"`
	Descricao   string    `gorm:"column:descricao;not null"`
	Capacidade  int       `gorm:"column:capacidade;not null"`
}

func (AreaComum) TableName() string {
	return "areacomum"
}
