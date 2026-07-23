## ADDED Requirements

### Requirement: Porteiro can release a visitor
The system SHALL allow a porteiro to release a visitor by setting `liberado = true` and creating a `Visita` log entry.

#### Scenario: Successful release
- **WHEN** the porteiro sends `PATCH /porteiros/visitantes/{id}/liberar` with body `{"porteiro_id": "<uuid>"}`
- **THEN** the system sets `liberado = true` on the visitor
- **AND** creates a `Visita` record with `dataEntrada = NOW()`, `porteiro_id`, `visitante_id`, and `morador_id = NULL`

#### Scenario: Release returns updated visitor
- **WHEN** the porteiro releases a visitor successfully
- **THEN** the system returns status 200 with the visitor object where `liberado` is `true`

### Requirement: Release with optional resident
The system SHALL allow specifying which resident the visit is for.

#### Scenario: Release with morador_id
- **WHEN** the porteiro sends `PATCH /porteiros/visitantes/{id}/liberar` with body `{"porteiro_id": "<uuid>", "morador_id": "<uuid>"}`
- **THEN** the system creates the `Visita` record with the provided `morador_id`

### Requirement: Release is idempotent for liberado
The system SHALL allow releasing an already-released visitor without error — subsequent releases still create new `Visita` records.

#### Scenario: Release already released visitor
- **WHEN** the porteiro releases a visitor that already has `liberado = true`
- **THEN** the system returns status 200
- **AND** creates a new `Visita` record for this release event

### Requirement: Validations for release request
The system SHALL validate that the provided UUIDs exist in the database.

#### Scenario: Visitor not found
- **WHEN** the porteiro sends `PATCH /porteiros/visitantes/{non-existent-id}/liberar`
- **THEN** the system returns 404 Not Found

#### Scenario: Porteiro not found
- **WHEN** the porteiro sends a release request with a non-existent `porteiro_id`
- **THEN** the system returns 400 Bad Request

### Requirement: Release is transactional
The system SHALL ensure that both the `liberado` update and the `Visita` creation succeed or fail together.

#### Scenario: Transaction rollback on failure
- **WHEN** the visitor update succeeds but the visit creation fails
- **THEN** the `liberado` update SHALL be rolled back
