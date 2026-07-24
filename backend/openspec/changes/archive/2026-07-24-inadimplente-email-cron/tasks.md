## 1. Config: MailtrapConfig

- [x] 1.1 Add `MailtrapConfig` struct to `internal/config/config.go` with `Host`, `Port`, `Username`, `Password`, `FromEmail` fields loaded from env vars `MAILTRAP_*`

## 2. Model: Notificacao

- [x] 2.1 Create `internal/models/notificacao.go` — `Notificacao` struct mapping the existing `Notificacao` table (`id`, `tipo`, `destinatario_id`, `mensagem`, `status`, `dataEnvio`)

## 3. Provider: Mailtrap

- [x] 3.1 Create `internal/providers/mailtrap/mailtrap.go` — `MailtrapProvider` with `Send(to, subject, body string) error` using `net/smtp`

## 4. Repository: NotificacaoRepository

- [x] 4.1 Create `internal/repositories/notificacao_repository.go` — `NotificacaoRepository` interface + `GormNotificacaoRepository` with `Create` and `FindByStatus`

## 5. Service: EmailService

- [x] 5.1 Create `internal/services/email_service.go` — `EmailService` with `SendInadimplenteEmails(ctx)` that fetches inadimplentes via `PagamentoRepository`, sends emails via `MailtrapProvider`, and logs via `NotificacaoRepository`

## 6. Cron binary

- [x] 6.1 Create `cmd/cron/main.go` — standalone binary that initializes config, DB, repositories, services, and runs `EmailService.SendInadimplenteEmails` on a 24h ticker (daily at 8 AM)
- [x] 6.2 Add `run-cron` target to `Makefile`

## 7. Validation

- [x] 7.1 Run `make build` to ensure project compiles
- [x] 7.2 Run `make test` to ensure existing tests pass
