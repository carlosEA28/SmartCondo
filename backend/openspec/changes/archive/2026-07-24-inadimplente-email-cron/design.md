## Context

The `Notificacao` table already exists in the database with fields: `id`, `tipo` (EMAIL/SMS), `destinatario_id` (FK → Usuario), `mensagem`, `status` (PENDENTE/ENVIADA/FALHA), `dataEnvio`. The `InadimplenteService` and `PagamentoRepository` already exist to find delinquent residents. Mailtrap provides SMTP credentials for sending test emails in development.

## Goals / Non-Goals

**Goals:**
- Add `MailtrapConfig` to `config.go` (host, port, username, password, from email)
- Add `MailtrapProvider` that sends emails via SMTP using `net/smtp`
- Add `Notificacao` model mapping the existing table
- Add `NotificacaoRepository` with `Create` and `FindByStatus`
- Add `EmailService` that fetches inadimplentes, builds email content, sends via MailtrapProvider, and logs each in `Notificacao`
- Add `cmd/cron/main.go` — standalone binary that runs the cronjob daily at 8 AM
- Update `Makefile` with `run-cron` target

**Non-Goals:**
- No SMS notifications (table supports it but only email for now)
- No web dashboard for sent emails
- No retry mechanism for failed sends (beyond logging as FALHA)
- No unsubscribe link in the email (can be added later)

## Decisions

### 1. Standalone cron binary vs embedded goroutine
**Decision**: Separate `cmd/cron/main.go` binary. The cron binary shares the same DB connection and repositories as the API.
**Rationale**: The cronjob is an independent process that can be scheduled separately (systemd timer, Kubernetes CronJob). It does not need to run inside the API server.

### 2. SMTP via stdlib vs third-party library
**Decision**: Use `net/smtp` (Go standard library).
**Rationale**: The only provider in dev is Mailtrap, which is plain SMTP. No third-party dependency needed. Can be swapped for a real provider later by changing the provider implementation.

### 3. Email content format
**Decision**: Plain text email with: resident name, list of overdue amounts and due dates, and a message to contact the síndico.
**Rationale**: Simple, no HTML template engine needed. Can be enhanced later.

### 4. Notification logging
**Decision**: After each email attempt, create a `Notificacao` record with `tipo = 'EMAIL'`, the send result status, and `dataEnvio`.
**Rationale**: The table exists and provides an audit trail.

### 5. Cron schedule
**Decision**: Run daily at 8 AM. The cron binary sleeps for the interval or uses a simple `time.Ticker`.
**Rationale**: Keeps it simple for deployment. Can be wrapped in a systemd timer externally.

## Data Flow

```
time.Ticker (24h)
  → cmd/cron/main.go calls EmailService.SendInadimplenteEmails()
    → PagamentoRepository.FindInadimplentes(ctx)
    → For each unique morador:
        1. Build email body (name, overdue list)
        2. MailtrapProvider.Send(to, subject, body)
        3. NotificacaoRepository.Create(record with status ENVIADA or FALHA)
    ← Summary log: "Sent N reminders, M failures"
```

## Architecture

```
cmd/cron/main.go
  → services.EmailService
    → repositories.PagamentoRepository (find inadimplentes)
    → providers.mailtrap.MailtrapProvider (send SMTP)
    → repositories.NotificacaoRepository (log results)
```

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| SMTP credentials in env vars | Already following the pattern (AWS creds in .env) |
| Cron overlaps with previous run | Not a concern — sending is idempotent, residents get duplicate emails at most |
| Email bounces / invalid addresses | Logged as FALHA in Notificacao. No retry yet. |
| Mailtrap rate limits | Low volume (condo has ~hundreds of residents). Acceptable in dev. |
