package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDENTE"
	PaymentStatusPaid      PaymentStatus = "PAGO"
	PaymentStatusOverdue   PaymentStatus = "ATRASADO"
)

type Pagamento struct {
	ID            uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Valor         float64        `gorm:"column:valor;type:decimal(10,2);not null"`
	Vencimento    time.Time      `gorm:"column:vencimento;type:date;not null"`
	DataPagamento *time.Time     `gorm:"column:datapagamento;type:date"`
	Status        PaymentStatus  `gorm:"column:status;size:20;not null;default:PENDENTE"`
	MoradorID     uuid.UUID      `gorm:"column:morador_id;type:uuid;not null"`
	Morador       User           `gorm:"foreignKey:MoradorID"`
}

func (Pagamento) TableName() string {
	return "pagamento"
}
