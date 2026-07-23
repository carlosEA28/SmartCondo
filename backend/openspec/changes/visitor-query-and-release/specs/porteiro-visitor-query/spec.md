## ADDED Requirements

### Requirement: Porteiro can search visitors by name
The system SHALL allow a porteiro to search visitors by their name using a partial match (LIKE/ILIKE).

#### Scenario: Search by full name
- **WHEN** the porteiro sends `GET /porteiros/visitantes?nome=João`
- **THEN** the system returns a list of visitors whose name contains "João"

#### Scenario: Search by partial name
- **WHEN** the porteiro sends `GET /porteiros/visitantes?nome=Silv`
- **THEN** the system returns a list of visitors whose name contains "Silv" (e.g., "Silva", "Silveira")

#### Scenario: No results found
- **WHEN** the porteiro sends `GET /porteiros/visitantes?nome=NobodyXYZ`
- **THEN** the system returns an empty list with status 200

### Requirement: Porteiro can search visitors by CPF
The system SHALL allow a porteiro to search visitors by their exact CPF.

#### Scenario: Search by CPF
- **WHEN** the porteiro sends `GET /porteiros/visitantes?cpf=12345678901`
- **THEN** the system returns the visitor with that exact CPF

#### Scenario: Search by non-existent CPF
- **WHEN** the porteiro sends `GET /porteiros/visitantes?cpf=00000000000`
- **THEN** the system returns an empty list with status 200

### Requirement: Porteiro can search visitors by phone
The system SHALL allow a porteiro to search visitors by their phone number using partial match.

#### Scenario: Search by phone
- **WHEN** the porteiro sends `GET /porteiros/visitantes?telefone=1199999`
- **THEN** the system returns visitors whose phone contains "1199999"

### Requirement: Porteiro can filter visitors by liberado status
The system SHALL allow a porteiro to filter visitors by their release status.

#### Scenario: Filter by liberado true
- **WHEN** the porteiro sends `GET /porteiros/visitantes?liberado=true`
- **THEN** the system returns only visitors with `liberado = true`

#### Scenario: Filter by liberado false
- **WHEN** the porteiro sends `GET /porteiros/visitantes?liberado=false`
- **THEN** the system returns only visitors with `liberado = false`

### Requirement: Combined search filters
The system SHALL allow combining multiple search parameters.

#### Scenario: Search by name and liberado status
- **WHEN** the porteiro sends `GET /porteiros/visitantes?nome=João&liberado=false`
- **THEN** the system returns visitors whose name contains "João" AND have `liberado = false`

### Requirement: At least one filter is required
The system SHALL require at least one search parameter, returning 400 if all are empty.

#### Scenario: Missing filters
- **WHEN** the porteiro sends `GET /porteiros/visitantes`
- **THEN** the system returns 400 Bad Request with error message

### Requirement: Visitor search response format
The system SHALL return a list of visitor objects in the same format as `VisitorResponseDTO`.

#### Scenario: Successful search
- **WHEN** the porteiro searches with valid parameters
- **THEN** the system returns status 200 with `[{id, name, cpf, phone, photo, liberado}, ...]`
