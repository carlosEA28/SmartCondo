package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MailtrapProvider interface {
	Send(to, subject, body string) error
}

type EmailService struct {
	pagamentoRepo   repositories.PagamentoRepository
	notificacaoRepo repositories.NotificacaoRepository
	mailtrap        MailtrapProvider
	fromEmail       string
}

func NewEmailService(
	pagamentoRepo repositories.PagamentoRepository,
	notificacaoRepo repositories.NotificacaoRepository,
	mailtrap MailtrapProvider,
	fromEmail string,
) *EmailService {
	return &EmailService{
		pagamentoRepo:   pagamentoRepo,
		notificacaoRepo: notificacaoRepo,
		mailtrap:        mailtrap,
		fromEmail:       fromEmail,
	}
}

func (s *EmailService) SendInadimplenteEmails(ctx context.Context) error {
	payments, err := s.pagamentoRepo.FindInadimplentes(ctx)
	if err != nil {
		return fmt.Errorf("fetch inadimplentes: %w", err)
	}

	grouped := groupByMorador(payments)

	var sent, failed int
	for _, entry := range grouped {
		to := entry.morador.Email
		if to == "" {
			log.Warn().Str("morador_id", entry.moradorID.String()).Msg("morador has no email, skipping")
			continue
		}

		subject := "Lembrete de Pagamento em Atraso - SmartCondo"
		body := buildEmailBody(entry.morador.FullName, entry.payments)

		if err := s.mailtrap.Send(to, subject, body); err != nil {
			log.Error().Err(err).Str("morador_id", entry.moradorID.String()).Msg("failed to send email")
			s.logNotification(ctx, entry.moradorID, "FALHA ao enviar email: "+err.Error(), models.NotificacaoStatusFalha)
			failed++
			continue
		}

		s.logNotification(ctx, entry.moradorID, "Email de cobrança enviado com sucesso", models.NotificacaoStatusEnviada)
		sent++
	}

	log.Info().Int("sent", sent).Int("failed", failed).Msg("inadimplente email cron finished")
	return nil
}

type moradorPayments struct {
	moradorID uuid.UUID
	morador   models.User
	payments  []models.Pagamento
}

func groupByMorador(payments []models.Pagamento) []*moradorPayments {
	grouped := make(map[uuid.UUID]*moradorPayments)
	order := make([]uuid.UUID, 0)
	for i := range payments {
		p := &payments[i]
		if _, ok := grouped[p.MoradorID]; !ok {
			grouped[p.MoradorID] = &moradorPayments{
				moradorID: p.MoradorID,
				morador:   p.Morador,
			}
			order = append(order, p.MoradorID)
		}
		grouped[p.MoradorID].payments = append(grouped[p.MoradorID].payments, *p)
	}
	result := make([]*moradorPayments, 0, len(order))
	for _, id := range order {
		result = append(result, grouped[id])
	}
	return result
}

func buildEmailBody(name string, payments []models.Pagamento) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Olá %s,\n\n", name))
	b.WriteString("Identificamos que você possui os seguintes pagamentos em atraso:\n\n")

	for _, p := range payments {
		b.WriteString(fmt.Sprintf("  - Valor: R$ %.2f | Vencimento: %s\n",
			p.Valor, p.Vencimento.Format("02/01/2006")))
	}

	b.WriteString("\nPor favor, regularize sua situação o quanto antes.\n")
	b.WriteString("Em caso de dúvidas, entre em contato com o síndico.\n\n")
	b.WriteString("Atenciosamente,\nEquipe SmartCondo")
	return b.String()
}

func (s *EmailService) logNotification(ctx context.Context, destinatarioID uuid.UUID, mensagem string, status models.NotificacaoStatus) {
	now := time.Now()
	notificacao := &models.Notificacao{
		Tipo:           models.NotificacaoTipoEmail,
		DestinatarioID: destinatarioID,
		Mensagem:       mensagem,
		Status:         status,
		DataEnvio:      &now,
	}

	if err := s.notificacaoRepo.Create(ctx, notificacao); err != nil {
		log.Error().Err(err).Msg("failed to log notification")
	}
}
