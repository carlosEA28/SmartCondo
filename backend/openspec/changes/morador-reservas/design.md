## Context

`AreaComum` and `Reserva` tables exist in the database but have no Go code. `Pagamento` model + repository already exist from the `sindico-inadimplentes` feature. The `RequireSindicoRole` middleware exists — a `RequireMoradorRole` analog is needed. Residents must pass two checks to reserve: no overdue payments and no time slot conflicts.

## Goals / Non-Goals

**Goals:**
- Add `RequireMoradorRole` middleware (reuse pattern from `RequireSindicoRole`)
- Add `GET /morador/areas-comuns` to list available common areas
- Add `POST /morador/reservas` to create a reservation (delinquency check + conflict detection)
- Add `GET /morador/reservas` to list the authenticated resident's reservations
- Add `PATCH /morador/reservas/:id/cancelar` to cancel own reservation

**Non-Goals:**
- No síndico confirmation flow (sindico confirms via existing future endpoint)
- No admin override
- No recurring reservations
- No notification on creation/cancellation

## Decisions

### 1. `RequireMoradorRole` middleware
**Decision**: New middleware reusing `X-User-ID` header, validating `user.Role == models.RoleMorador`
**Rationale**: Follows the exact same pattern as `RequireSindicoRole`. Creates a reusable foundation for future morador-only endpoints.

### 2. Delinquency check: reuse `PagamentoRepository`
**Decision**: Add `HasOverduePayments(ctx, moradorID)` method to `PagamentoRepository` — queries if any `Pagamento` with `status = 'ATRASADO'` and matching `morador_id` exists
**Rationale**: Reuses the existing repository rather than duplicating logic. The `FindInadimplentes` method returns all delinquent payments, but we just need a boolean check for a specific resident.

### 3. Time conflict: SQL overlap query in `ReservaRepository`
**Decision**: `FindConflicting(ctx, areaID, data, horaInicio, horaFim, excludeReservaID)` with overlap condition: `horaInicio < :horaFim AND horaFim > :horaInicio AND data = :data AND areacomum_id = :areaID AND status != 'CANCELADA'`
**Rationale**: Same approach as the sindico-areas-comuns design. Efficient single query.

### 4. Reservation creation flow
**Decision**: `POST /morador/reservas` creates with `status = 'PENDENTE'`. The `sindico_id` is null initially — the síndico confirms later via a separate endpoint.
**Rationale**: The Reserva table allows nullable `sindico_id`. The status workflow is PENDENTE → CONFIRMADA (by síndico) or CANCELADA (by morador or síndico).

### 5. List only own reservations
**Decision**: `GET /morador/reservas` filters by the authenticated resident's `morador_id`
**Rationale**: A resident should only see their own reservations, not everyone's.

## Data Flow

```
POST /morador/reservas { "areacomum_id", "data", "horaInicio", "horaFim" }
  → RequireMoradorRole middleware (X-User-ID → morador_id in context)
  → ReservaHandler.create
    → ReservaService.CreateReserva
      1. Check area exists → AreaComumRepository.FindByID
      2. Check delinquency → PagamentoRepository.HasOverduePayments(moradorID)
         → if true: 403 ErrMoradorInadimplente
      3. Check time conflict → ReservaRepository.FindConflicting(...)
         → if conflict: 409 ErrReservaConflito
      4. Create reserva (status = PENDENTE, morador_id from context)
    ← ReservaResponseDTO
  ← JSON 201

PATCH /morador/reservas/:id/cancelar
  → RequireMoradorRole middleware
  → ReservaHandler.cancel
    → ReservaService.CancelReserva
      1. Find reserva by ID
      2. Verify morador_id matches authenticated user
         → if not owner: 403 ErrReservaNotOwner
      3. Update status to CANCELADA
    ← ReservaResponseDTO
  ← JSON 200
```

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  RequireMoradorRole middleware                               │
│  - Reads X-User-ID, validates user exists + RoleMorador     │
│  - Injects user_id into context                              │
├──────────────────────────────────────────────────────────────┤
│  AreaComumHandler + ReservaHandler (internal/server/)        │
├──────────────────────────────────────────────────────────────┤
│  AreaComumService (internal/services/)                       │
│  - ListAreas                                                 │
├──────────────────────────────────────────────────────────────┤
│  ReservaService (internal/services/)                         │
│  - CreateReserva (delinquency + conflict checks)             │
│  - ListMinhasReservas (filtered by morador_id)               │
│  - CancelReserva (owner check)                               │
├──────────────────────────────────────────────────────────────┤
│  Repositories (internal/repositories/)                       │
│  - AreaComumRepository: FindByID, FindAll                    │
│  - ReservaRepository: FindByID, FindByMorador,              │
│    FindConflicting, Create, UpdateStatus                     │
│  - PagamentoRepository: +HasOverduePayments                  │
└──────────────────────────────────────────────────────────────┘
```

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Time conflict race condition | Use GORM transaction with `FOR UPDATE` on conflict query |
| No `sindico_id` on initial reservation | Acceptable — nullable column. Sindico sets it on confirmation |
| Resident could cancel after use (no penalty) | No penalty system yet. Can be added later |
