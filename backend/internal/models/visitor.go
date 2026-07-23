package models

import "github.com/google/uuid"

type Visitor struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"column:nome;size:100;not null"`
	CPF      string    `gorm:"column:cpf;size:11;unique;not null"`
	Phone    string    `gorm:"column:telefone;size:15;not null"`
	Photo    string    `gorm:"column:foto;size:255"`
	Liberado bool      `gorm:"column:liberado;not null;default:false"`
}

func (Visitor) TableName() string {
	return "visitante"
}
