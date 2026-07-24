## Context

The `AreaComum`, `Reserva`, and `Pagamento` tables exist in the database (migration 000001) with proper FKs and CHECK constraints but have zero Go code. The `RequireSindicoRole` middleware already exists for sindico-protected routes. This feature delivers two capabilities: CRUD for common areas (síndico-managed) and reservation management with business rules.

## Goals / Non-Goals

**Goals:**
- Full CRUD for `AreaComum` with unique name enforcement
- Create reservations with time conflict detection (same area, overlapping time)
- Block reservations for residents with overdue payments (`Pagamento.status = 'ATRASADO'`)
- Reservation status workflow: PENDENTE → CONFIRMADA (by síndico) or CANCELADA
- All sindico routes protected by `RequireSindicoRole` middleware
- Follow the same layered patterns as Comunicado and other features

**Non-Goals:**
- Resident-facing reservation endpoints (residents create requests through a future separate flow)
- Automated payment status webhook or integration
- Recurring reservations
- Notifications when a reservation is confirmed/cancelled

## Decisions

### 1. Two separate services: `AreaComumService` + `ReservaService`
**Decision**: Split into two services rather than one monolithic service
**Rationale**: Areas and reservations are distinct domain concepts with different lifecycle. A síndico manages areas independently of individual reservations.

### 2. Unique name enforcement via repository query
**Decision**: Before create/update, query for existing area with the same name (case-insensitive). No unique DB constraint — the table doesn't have one, and adding a migration for this alone isn't warranted.
**Rationale**: Avoids a migration. The service-level check is sufficient since writes go through the service.

### 3. Time conflict detection via SQL overlap query
**Decision**: Repository method `FindConflicting(ctx, areaID, data, horaInicio, horaFim, excludeID)` that checks for overlapping time ranges on the same date and area. Only checks reservations with `status != 'CANCELADA'`.
**Rationale**: GORM can express the overlap condition: `horaInicio < :horaFim AND horaFim > :horaInicio`.

### 4. Delinquency check via Pagamento query
**Decision**: Repository method `HasOverduePayments(ctx, moradorID)` that checks if any `Pagamento` with `status = 'ATRASADO'` exists for the resident.
**Rationale**: Simple existence check. The `Pagamento.status` CHECK constraint ensures only valid statuses.

### 5. Reservation creation includes sindico confirmation flag
**Decision**: `POST /sindico/reservas` creates with `status = 'PENDENTE'`. The síndico then explicitly confirms via `PATCH /sindico/reservas/:id/confirmar`. This gives the síndico control over the reservation workflow.
**Rationale**: Matches the requirement doc's business rules — síndico must have oversight.

## Architecture

```
┌────────────────────────────────────────────────────────────┐
│  AreaComumHandler + ReservaHandler (internal/server/)      │
│  - CRUD areas, CRUD reservas, confirmar, cancelar         │
│  interfaces: areaComumService, reservaService              │
├────────────────────────────────────────────────────────────┤
│  AreaComumService (internal/services/)                     │
│  - Create, List, Get, Update, Delete                       │
│  - Unique name validation                                  │
│  Depends on: AreaComumRepository                           │
├────────────────────────────────────────────────────────────┤
│  ReservaService (internal/services/)                       │
│  - Create (with conflict + payment checks)                 │
│  - List, Get, Confirmar, Cancelar                          │
│  Depends on: ReservaRepository, PagamentoRepository        │
├────────────────────────────────────────────────────────────┤
│  Repositories (internal/repositories/)                     │
│  - AreaComumRepository: FindByID, FindAll, FindByName,     │
│    Create, Update, Delete                                  │
│  - ReservaRepository: FindByID, FindAll, FindConflicting,   │
│    Create, UpdateStatus                                    │
│  - PagamentoRepository: HasOverduePayments(moradorID)      │
├────────────────────────────────────────────────────────────┤
│  Models (internal/models/)                                 │
│  - AreaComum (table: areacomum)                            │
│  - Reserva (table: reserva, FKs to AreaComum + Usuario)    │
│  - Pagamento (table: pagamento)                            │
└────────────────────────────────────────────────────────────┘
```

## Data Flow: Reservation Creation

```
POST /sindico/reservas { morador_id, areacomum_id, data, horaInicio, horaFim }
  → RequireSindicoRole middleware (validates X-User-ID)
  → ReservaHandler.create
    → ReservaService.CreateReserva
      1. Validate morador exists + has RoleMorador
      2. Check PagamentoRepository.HasOverduePayments(moradorID)
         → if true: return ErrMoradorInadimplente
      3. Check ReservaRepository.FindConflicting(areaID, data, horaInicio, horaFim)
         → if conflict: return ErrReservaConflito
      4. Create reserva with status = 'PENDENTE'
    ← ReservaResponseDTO
  ← JSON 201
```

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Race condition: two simultaneous reservations for same time slot | Use GORM transaction with `FOR UPDATE` on conflicting query inside the transaction |
| No unique constraint on AreaComum.nome in DB | Service-level check before create/update. Acceptable since all writes go through the API |
| `Pagamento` table may have many records per resident | `HasOverduePayments` uses `EXISTS` — efficient query, no full table scan |
| "Inadimplente" definition is simplified (any overdue = blocked) | Current requirement. Can be refined later (e.g., only if overdue > 30 days) |
