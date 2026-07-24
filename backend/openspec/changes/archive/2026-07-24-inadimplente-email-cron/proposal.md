## Why

Delinquent residents (inadimplentes) need to be reminded of their overdue payments via email. Currently, only the síndico can see who is delinquent via the `GET /sindico/inadimplentes` endpoint — residents receive no automatic notification. A daily cronjob will send reminder emails using Mailtrap (SMTP), improving payment collection and resident communication.

## What Changes

- New `MailtrapProvider` in `internal/providers/mailtrap/` — SMTP email sender using Mailtrap credentials
- New `Notificacao` model (`internal/models/notificacao.go`) — maps to the existing `Notificacao` table
- New `NotificacaoRepository` — creates notification records and queries by status
- New `EmailService` — builds and sends delinquency reminder emails, logs to `Notificacao` table
- New cronjob in `cmd/cron/main.go` — standalone binary or goroutine that runs daily, fetches inadimplentes, and sends emails
- New env vars: `MAILTRAP_HOST`, `MAILTRAP_PORT`, `MAILTRAP_USERNAME`, `MAILTRAP_PASSWORD`, `MAILTRAP_FROM_EMAIL`
- No new DB migrations — the `Notificacao` table already exists
- No breaking changes

## Capabilities

### New Capabilities
- `inadimplente-email`: Automatic email notification for delinquent residents via cronjob

### Modified Capabilities
- *(none)*

## Impact

- **New provider**: `internal/providers/mailtrap/mailtrap.go` — SMTP email client
- **New model**: `internal/models/notificacao.go` — maps `Notificacao` table
- **New repository**: `internal/repositories/notificacao_repository.go` — `Create` + `FindByStatus`
- **New service**: `internal/services/email_service.go` — `SendInadimplenteEmails(ctx)`
- **New cmd**: `cmd/cron/main.go` — standalone cron binary with daily schedule
- **Config changes**: `MailtrapConfig` added to `internal/config/config.go`
- **Dependencies**: `net/smtp` (stdlib) — no new third-party deps
