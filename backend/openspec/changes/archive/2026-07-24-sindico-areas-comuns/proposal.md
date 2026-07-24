## Why

The síndico needs to manage common areas (salão de festas, churrasqueira, piscina, etc.) that residents can reserve. The `AreaComum` and `Reserva` tables already exist in the database but have no API or Go code. Residents also need protection — areas must have unique names to avoid confusion, reservations cannot overlap in time for the same area, and residents with overdue payments (inadimplentes) must be blocked from reserving.

## What Changes

- New `POST /sindico/areas-comuns` — síndico creates a common area (unique name enforced)
- New `GET /sindico/areas-comuns` — list all common areas
- New `GET /sindico/areas-comuns/:id` — view a single area
- New `PUT /sindico/areas-comuns/:id` — update area details
- New `DELETE /sindico/areas-comuns/:id` — remove an area
- New `POST /sindico/reservas` — create a reservation with time conflict detection and payment check
- New `GET /sindico/reservas` — list all reservations
- New `GET /sindico/reservas/:id` — view a single reservation
- New `PATCH /sindico/reservas/:id/confirmar` — síndico confirms a reservation
- New `PATCH /sindico/reservas/:id/cancelar` — síndico cancels a reservation
- All sindico routes protected by the existing `RequireSindicoRole` middleware
- No new migrations needed — tables already exist

## Capabilities

### New Capabilities
- `sindico-areas-comuns`: Síndico CRUD for common areas with unique name constraint
- `sindico-reservas`: Reservation management with time conflict detection, payment status check, and status workflow

### Modified Capabilities
- *(none)*

## Impact

- **3 new models**: `AreaComum`, `Reserva`, `Pagamento` in `internal/models/`
- **3 new DTOs**: `internal/dto/areacomum.go` and `internal/dto/reserva.go`
- **2 new repositories**: `AreaComumRepository` and `ReservaRepository` in `internal/repositories/`
- **2 new services**: `AreaComumService` and `ReservaService` in `internal/services/`
- **2 new handlers**: `areaComumHandler` and `reservaHandler` in `internal/server/`
- **New errors**: `ErrAreaComumNotFound`, `ErrAreaComumAlreadyExists`, `ErrReservaNotFound`, `ErrReservaConflito`, `ErrMoradorInadimplente`, etc.
- **Route registration**: routes grouped under `/sindico/areas-comuns` and `/sindico/reservas`
- **No breaking changes** — existing endpoints untouched
