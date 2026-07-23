## Why

Porteiros need to look up registered visitors and control their entry into the condominium. Currently, visitors can be created but there is no way to release them (mark as authorized to enter) or to query them with search filters tailored to the gatekeeper's workflow. The `Visita` table already exists in the schema to log entry events, but it is unused. This change enables the full gatekeeper flow: search visitors and release them, creating an audit trail.

## What Changes

- New `GET /porteiros/visitantes` endpoint with query parameters for porteiro to search visitors by name, CPF, phone, and liberado status
- New `PATCH /porteiros/visitantes/:id/liberar` endpoint for porteiro to release a visitor, setting `liberado = true` on the visitor record and inserting a row into `Visita` to log the release event
- New migration to add `porteiro_id` tracking to the release flow (for the Visita record)
- New internal models: `Visit` (Visita) — the Go model already exists implicitly via migration, needs a struct
- New repository, service, handler following established visitor pattern

## Capabilities

### New Capabilities
- `porteiro-visitor-query`: Search and list visitors with filters (name, CPF, phone, liberado status) tailored for the gatekeeper's workflow
- `porteiro-visitor-release`: Release a visitor by setting `liberado = true` and logging the event in the `Visita` table with porteiro identification

### Modified Capabilities
- *(none — no existing specs to modify)*

## Impact

- **New model**: `Visit` in `internal/models/visit.go` (maps to `Visita` table)
- **New repository**: `internal/repositories/visit_repository.go` with `VisitRepository` interface
- **New service**: `internal/services/porteiro_service.go` with query and release logic
- **New handler**: `internal/server/porteiro_handler.go` with search + release endpoints
- **New DTOs**: `internal/dto/porteiro.go` with request/response structs
- **New migration**: Logging `porteiro_id` in Visit records
- **Route registration**: Two new routes in `internal/server/server.go`
- **No breaking changes** — existing visitor endpoints remain untouched
