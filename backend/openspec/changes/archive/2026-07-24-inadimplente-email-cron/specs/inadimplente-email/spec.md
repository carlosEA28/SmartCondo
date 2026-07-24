## ADDED Requirements

### Requirement: Automated email notification for delinquent residents
The system SHALL send an email reminder to each resident with overdue payments via a daily cronjob.

#### Scenario: Daily cron sends emails to all delinquent residents
- **GIVEN** the cronjob runs at its scheduled time (daily at 8 AM)
- **AND** there are residents with overdue payments
- **WHEN** the cronjob executes
- **THEN** each delinquent resident receives an email with their overdue details
- **AND** a `Notificacao` record is created for each email with status "ENVIADA"

#### Scenario: No delinquent residents
- **GIVEN** there are no residents with overdue payments
- **WHEN** the cronjob executes
- **THEN** no emails are sent

#### Scenario: Email send failure is logged
- **GIVEN** the cronjob attempts to send an email to a delinquent resident
- **WHEN** the SMTP send fails
- **THEN** a `Notificacao` record is created with status "FALHA"

#### Scenario: Email contains overdue details
- **GIVEN** a delinquent resident is found
- **WHEN** the email is sent
- **THEN** the email body includes the resident's name, list of overdue amounts with due dates, and a message to contact the síndico
