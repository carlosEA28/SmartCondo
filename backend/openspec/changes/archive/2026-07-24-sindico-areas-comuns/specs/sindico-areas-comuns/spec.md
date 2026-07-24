## ADDED Requirements

### Requirement: Síndico can create a common area
The system SHALL allow a síndico to create a common area with a unique name.

#### Scenario: Successful creation
- **WHEN** the síndico sends `POST /sindico/areas-comuns` with body `{"nome": "Salão de Festas", "descricao": "Salão com capacidade para 50 pessoas", "capacidade": 50}`
- **THEN** the system returns status 201 with the created area

#### Scenario: Duplicate name
- **WHEN** the síndico sends `POST /sindico/areas-comuns` with a `nome` that already exists
- **THEN** the system returns status 409 Conflict

#### Scenario: Missing required fields
- **WHEN** the síndico sends `POST /sindico/areas-comuns` without `nome` or `descricao`
- **THEN** the system returns status 400 Bad Request

### Requirement: Síndico can list all common areas
The system SHALL return all common areas ordered by name.

#### Scenario: List all areas
- **WHEN** a user sends `GET /sindico/areas-comuns`
- **THEN** the system returns status 200 with an array of common areas

#### Scenario: Empty list
- **WHEN** there are no common areas
- **THEN** the system returns status 200 with an empty array

### Requirement: Síndico can view a common area by ID
The system SHALL return a single common area.

#### Scenario: View area by ID
- **WHEN** a user sends `GET /sindico/areas-comuns/{id}`
- **THEN** the system returns status 200 with the area object

#### Scenario: Area not found
- **WHEN** a user sends `GET /sindico/areas-comuns/{non-existent-id}`
- **THEN** the system returns status 404 Not Found

### Requirement: Síndico can update a common area
The system SHALL allow updating a common area's nome, descricao, and capacidade.

#### Scenario: Successful update
- **WHEN** the síndico sends `PUT /sindico/areas-comuns/{id}` with updated fields
- **THEN** the system returns status 200 with the updated area

#### Scenario: Update to duplicate name
- **WHEN** the síndico sends `PUT /sindico/areas-comuns/{id}` with a `nome` that another area already uses
- **THEN** the system returns status 409 Conflict

#### Scenario: Area not found on update
- **WHEN** the síndico sends `PUT /sindico/areas-comuns/{non-existent-id}`
- **THEN** the system returns status 404 Not Found

### Requirement: Síndico can delete a common area
The system SHALL allow deleting a common area.

#### Scenario: Successful delete
- **WHEN** the síndico sends `DELETE /sindico/areas-comuns/{id}`
- **THEN** the system returns status 204 No Content

#### Scenario: Area not found on delete
- **WHEN** the síndico sends `DELETE /sindico/areas-comuns/{non-existent-id}`
- **THEN** the system returns status 404 Not Found
