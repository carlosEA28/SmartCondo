## ADDED Requirements

### Requirement: Síndico can create a reservation
The system SHALL allow a síndico to create a reservation with time conflict detection and payment status validation.

#### Scenario: Successful reservation
- **WHEN** the síndico sends `POST /sindico/reservas` with body `{"morador_id": "<uuid>", "areacomum_id": "<uuid>", "data": "2026-08-01", "horaInicio": "14:00", "horaFim": "16:00"}`
- **THEN** the system returns status 201 with the reservation (status = "PENDENTE")

#### Scenario: Time conflict with existing reservation
- **WHEN** the síndico creates a reservation for the same area, same date with overlapping times
- **THEN** the system returns status 409 Conflict with conflict error

#### Scenario: Resident has overdue payments
- **WHEN** the síndico creates a reservation for a resident with `Pagamento.status = 'ATRASADO'`
- **THEN** the system returns status 403 Forbidden with delinquency error

#### Scenario: Missing required fields
- **WHEN** the síndico sends `POST /sindico/reservas` without required fields
- **THEN** the system returns status 400 Bad Request

### Requirement: Síndico can list reservations
The system SHALL list all reservations ordered by date descending.

#### Scenario: List all reservations
- **WHEN** a user sends `GET /sindico/reservas`
- **THEN** the system returns status 200 with an array of reservations

### Requirement: Síndico can view a reservation by ID
The system SHALL return a single reservation with related data.

#### Scenario: View reservation
- **WHEN** a user sends `GET /sindico/reservas/{id}`
- **THEN** the system returns status 200 with the reservation including area name and resident name

#### Scenario: Reservation not found
- **WHEN** a user sends `GET /sindico/reservas/{non-existent-id}`
- **THEN** the system returns status 404 Not Found

### Requirement: Síndico can confirm a reservation
The system SHALL allow the síndico to change a reservation status to "CONFIRMADA".

#### Scenario: Confirm reservation
- **WHEN** the síndico sends `PATCH /sindico/reservas/{id}/confirmar`
- **THEN** the system returns status 200 with the reservation (status = "CONFIRMADA")

#### Scenario: Confirm non-existent reservation
- **WHEN** the síndico sends `PATCH /sindico/reservas/{non-existent-id}/confirmar`
- **THEN** the system returns status 404 Not Found

### Requirement: Síndico can cancel a reservation
The system SHALL allow the síndico to change a reservation status to "CANCELADA".

#### Scenario: Cancel reservation
- **WHEN** the síndico sends `PATCH /sindico/reservas/{id}/cancelar`
- **THEN** the system returns status 200 with the reservation (status = "CANCELADA")

#### Scenario: Cancel non-existent reservation
- **WHEN** the síndico sends `PATCH /sindico/reservas/{non-existent-id}/cancelar`
- **THEN** the system returns status 404 Not Found
