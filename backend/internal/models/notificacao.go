package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificacaoTipo string

const (
	NotificacaoTipoEmail NotificacaoTipo = "EMAIL"
	NotificacaoTipoSMS   NotificacaoTipo = "SMS"
)

type NotificacaoStatus string

const (
	NotificacaoStatusPendente NotificacaoStatus = "PENDENTE"
	NotificacaoStatusEnviada  NotificacaoStatus = "ENVIADA"
	NotificacaoStatusFalha    NotificacaoStatus = "FALHA"
)

type Notificacao struct {
	ID             uuid.UUID          `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Tipo           NotificacaoTipo    `gorm:"column:tipo;size:10;not null"`
	DestinatarioID uuid.UUID          `gorm:"column:destinatario_id;type:uuid;not null"`
	Destinatario   User               `gorm:"foreignKey:DestinatarioID"`
	Mensagem       string             `gorm:"column:mensagem;not null"`
	Status         NotificacaoStatus  `gorm:"column:status;size:20;not null;default:PENDENTE"`
	DataEnvio      *time.Time         `gorm:"column:dataenvio"`
}

func (Notificacao) TableName() string {
	return "notificacao"
}
