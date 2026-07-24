## ADDED Requirements

### Requirement: Síndico can publish an announcement
The system SHALL allow a síndico to publish a new announcement by providing `sindico_id`, `titulo`, and `descricao`.

#### Scenario: Successful publish
- **WHEN** the síndico sends `POST /sindico/comunicados` with body `{"sindico_id": "<uuid>", "titulo": "Manutenção", "descricao": "Descrição detalhada"}`
- **THEN** the system returns status 201 with the created announcement including `id`, `titulo`, `descricao`, `dataPublicacao`, `sindico_id`, and `sindico_nome`

#### Scenario: Missing required fields
- **WHEN** the síndico sends `POST /sindico/comunicados` without `titulo` or `descricao`
- **THEN** the system returns status 400 Bad Request

#### Scenario: Síndico not found
- **WHEN** the síndico sends `POST /sindico/comunicados` with a non-existent `sindico_id`
- **THEN** the system returns status 400 Bad Request

### Requirement: Anyone can list all announcements
The system SHALL return all announcements ordered by most recent publication date.

#### Scenario: List all announcements
- **WHEN** a user sends `GET /sindico/comunicados`
- **THEN** the system returns status 200 with an array of announcements ordered by `dataPublicacao` descending

#### Scenario: Empty list
- **WHEN** a user sends `GET /sindico/comunicados` and there are no announcements
- **THEN** the system returns status 200 with an empty array

### Requirement: Anyone can view a single announcement
The system SHALL return a single announcement by its ID.

#### Scenario: View announcement by ID
- **WHEN** a user sends `GET /sindico/comunicados/{id}`
- **THEN** the system returns status 200 with the announcement object

#### Scenario: Announcement not found
- **WHEN** a user sends `GET /sindico/comunicados/{non-existent-id}`
- **THEN** the system returns status 404 Not Found

### Requirement: Síndico can delete their own announcement
The system SHALL allow a síndico to delete an announcement, verifying the caller is the owner.

#### Scenario: Successful delete
- **WHEN** the síndico sends `DELETE /sindico/comunicados/{id}?sindico_id={uuid}`
- **THEN** the system returns status 204 No Content

#### Scenario: Announcement not found
- **WHEN** the síndico sends `DELETE /sindico/comunicados/{non-existent-id}?sindico_id={uuid}`
- **THEN** the system returns status 404 Not Found

#### Scenario: Not the owner
- **WHEN** a different síndico sends `DELETE /sindico/comunicados/{id}?sindico_id={different-uuid}`
- **THEN** the system returns status 403 Forbidden
