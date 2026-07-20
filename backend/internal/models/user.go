package models

import "github.com/google/uuid"

type Role string

const (
	RoleMorador  Role = "MORADOR"
	RolePorteiro Role = "PORTEIRO"
	RoleSindico  Role = "SINDICO"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "ATIVO"
	UserStatusInactive UserStatus = "INATIVO"
	UserStatusBlocked  UserStatus = "BLOQUEADO"
)

type User struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	FullName    string     `gorm:"column:nome;size:100;not null"`
	Email       string     `gorm:"column:email;size:100;unique;not null"`
	Password    string     `gorm:"column:senha;size:100;not null"`
	Phone       string     `gorm:"column:telefone;size:15;not null"`
	Status      UserStatus `gorm:"column:status;size:20;not null;default:ATIVO"`
	Role        Role       `gorm:"column:tipo;size:10;not null"`
	ApartmentID *uuid.UUID `gorm:"column:apartamento_id;type:uuid"`
	Apartment   *Apartment `gorm:"foreignKey:ApartmentID"`
	Responsible bool       `gorm:"column:responsavel;not null;default:false"`
}

func (User) TableName() string {
	return "usuario"
}
