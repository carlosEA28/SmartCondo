## ADDED Requirements

### Requirement: Síndico can list delinquent residents
The system SHALL allow the síndico to view all residents with at least one overdue payment (status = "ATRASADO").

#### Scenario: List delinquent residents
- **WHEN** the síndico sends `GET /sindico/inadimplentes`
- **THEN** the system returns status 200 with an array of delinquent residents, each containing morador info, total overdue amount, and list of overdue payments

#### Scenario: No delinquent residents
- **WHEN** the síndico sends `GET /sindico/inadimplentes` and no residents have overdue payments
- **THEN** the system returns status 200 with an empty array

#### Scenario: Unauthorized access
- **WHEN** a non-síndico user sends `GET /sindico/inadimplentes` without a valid `X-User-ID` header
- **THEN** the system returns status 401 Unauthorized

#### Scenario: Forbidden for non-síndico
- **WHEN** a user with role different from SINDICO sends `GET /sindico/inadimplentes`
- **THEN** the system returns status 403 Forbidden
