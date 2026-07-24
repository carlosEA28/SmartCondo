## Why

Residents (moradores) need to be able to reserve common areas (salão de festas, churrasqueira, piscina, etc.) for their use, and cancel their own reservations when plans change. The `AreaComum`, `Reserva`, and `Pagamento` tables already exist in the database — residents must not have overdue payments to reserve, and reservations cannot overlap in time for the same area.

## What Changes

- New `RequireMoradorRole` middleware (reads `X-User-ID`, validates user has `RoleMorador`)
- `AreaComum` model + repository for listing available areas
- `Reserva` model + repository for creating, listing, and canceling reservations
- `GET /morador/areas-comuns` — list all available common areas (public, no auth or morador-only)
- `POST /morador/reservas` — create a reservation with delinquency check + time conflict detection
- `GET /morador/reservas` — list the resident's own reservations
- `PATCH /morador/reservas/:id/cancelar` — cancel a reservation (owner check)
- No new migrations needed — `AreaComum`, `Reserva`, and `Pagamento` tables already exist

## Capabilities

### New Capabilities
- `morador-reservas`: Resident-facing reservation flow with delinquency check, time conflict detection, and cancellation

### Modified Capabilities
- *(none)*

## Impact

- **New middleware**: `RequireMoradorRole` in `internal/server/middleware/auth.go`
- **2 new models**: `AreaComum` + `Reserva` in `internal/models/`
- **2 new DTOs**: `internal/dto/reserva.go` + `internal/dto/areacomum.go`
- **2 new repositories**: `AreaComumRepository` + `ReservaRepository` in `internal/repositories/`
- **2 new services**: `AreaComumService` (list areas) + `ReservaService` (create with checks, list, cancel)
- **2 new handlers**: `areaComumHandler` + `reservaHandler` in `internal/server/`
- **New errors**: `ErrAreaComumNotFound`, `ErrReservaNotFound`, `ErrReservaConflito`, `ErrMoradorInadimplente`, `ErrReservaNotOwner`, `ErrInvalidReservaData`
- **Route registration**: 4 new routes under `/morador/`
- **No breaking changes** — existing endpoints untouched
