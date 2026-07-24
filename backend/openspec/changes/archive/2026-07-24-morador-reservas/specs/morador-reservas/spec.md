## ADDED Requirements

### Requirement: Resident can list common areas
The system SHALL allow any authenticated resident to list all available common areas.

#### Scenario: List all areas
- **WHEN** a resident sends `GET /morador/areas-comuns`
- **THEN** the system returns status 200 with an array of common areas

### Requirement: Resident can reserve a common area
The system SHALL allow a resident with no overdue payments to reserve a common area for a specific date and time.

#### Scenario: Successful reservation
- **WHEN** a non-delinquent resident sends `POST /morador/reservas` with body `{"areacomum_id": "<uuid>", "data": "2026-08-15", "horaInicio": "14:00", "horaFim": "16:00"}`
- **THEN** the system returns status 201 with the reservation (status = "PENDENTE")

#### Scenario: Resident is delinquent
- **WHEN** a resident with overdue payments sends `POST /morador/reservas`
- **THEN** the system returns status 403 Forbidden with delinquency error

#### Scenario: Time conflict
- **WHEN** a resident sends `POST /morador/reservas` for an area/date/time that overlaps an existing non-canceled reservation
- **THEN** the system returns status 409 Conflict with conflict error

#### Scenario: Area not found
- **WHEN** a resident sends `POST /morador/reservas` with a non-existent `areacomum_id`
- **THEN** the system returns status 404 Not Found

#### Scenario: Missing required fields
- **WHEN** a resident sends `POST /morador/reservas` without required fields
- **THEN** the system returns status 400 Bad Request

### Requirement: Resident can list their own reservations
The system SHALL list only the authenticated resident's reservations.

#### Scenario: List my reservations
- **WHEN** a resident sends `GET /morador/reservas`
- **THEN** the system returns status 200 with an array of that resident's reservations

#### Scenario: No reservations
- **WHEN** a resident with no reservations sends `GET /morador/reservas`
- **THEN** the system returns status 200 with an empty array

### Requirement: Resident can cancel their own reservation
The system SHALL allow a resident to cancel their own reservation.

#### Scenario: Successful cancellation
- **WHEN** a resident sends `PATCH /morador/reservas/{id}/cancelar`
- **THEN** the system returns status 200 with the reservation (status = "CANCELADA")

#### Scenario: Cancel another resident's reservation
- **WHEN** a resident sends `PATCH /morador/reservas/{id}/cancelar` for a reservation belonging to a different resident
- **THEN** the system returns status 403 Forbidden

#### Scenario: Reservation not found
- **WHEN** a resident sends `PATCH /morador/reservas/{non-existent-id}/cancelar`
- **THEN** the system returns status 404 Not Found
