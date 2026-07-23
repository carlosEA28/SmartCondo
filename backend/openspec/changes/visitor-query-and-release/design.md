## Context

The codebase has a complete layered architecture (models → repositories → services → handlers → routes) with GORM, PostgreSQL, and Gin. The `Visita` table already exists in the schema with `porteiro_id`, `visitante_id`, `morador_id`, `dataEntrada`, and `dataSaida` columns but has no Go model or code. The `Visitante` table has a `liberado` boolean column (added in migration 000002). Porteiros (users with `tipo = 'PORTEIRO'`) currently have no dedicated endpoints.

Two capabilities are needed:
1. **Porteiro visitor query** — search visitors by name, CPF, phone, and liberado status
2. **Porteiro visitor release** — set `liberado = true` on a visitor and create a `Visita` record

## Goals / Non-Goals

**Goals:**
- Add `GET /porteiros/visitantes` with query parameters (`nome`, `cpf`, `telefone`, `liberado`) for porteiro-facing search
- Add `PATCH /porteiros/visitantes/:id/liberar` to release a visitor and log the visit
- Add `Visit` model mapping the existing `Visita` table
- Add `VisitRepository` for creating visit records
- Follow the exact layered patterns from the existing visitor feature

**Non-Goals:**
- No authentication/authorization middleware — porteiro identity will be passed as a request parameter (porteiro_id) since auth is not implemented yet
- No `dataSaida` tracking (exit logging) — future work
- No changes to existing visitor endpoints
- No S3 interactions for this feature

## Decisions

### 1. New `PorteiroService` vs extending `VisitorService`
**Decision**: New `PorteiroService` in `internal/services/porteiro_service.go`
**Rationale**: Though it touches visitors, this is a distinct actor (porteiro) with different responsibilities. Keeping it separate avoids bloating `VisitorService` with porteiro-specific concerns (visit logging). The pattern matches how the codebase separates concerns.

### 2. Query endpoint: reuse `VisitorRepository` directly
**Decision**: Add `Search` method to `VisitorRepository` interface and implementation
**Rationale**: The query needs filtering on the `visitante` table. Adding a `Search` method with dynamic filters keeps the data access in the repository layer where it belongs. The porteiro service orchestrates the call.

### 3. Release endpoint: two repositories, one service method
**Decision**: `PorteiroService.ReleaseVisitor` calls `VisitorRepository.Update` (set liberado=true) and then `VisitRepository.Create` in a single transactional unit
**Rationale**: Both operations must succeed or fail together. GORM's `Transaction` ensures atomicity. If visitor update fails, no visit record is created.

### 4. Visit model: porteiro_id, visitante_id required; morador_id optional
**Decision**: Use `*uuid.UUID` for `morador_id` (nullable) since the porteiro may not know which resident the visitor is going to see at release time. The `dataEntrada` is set to `time.Now()` automatically.
**Rationale**: The schema has `morador_id` as NOT NULL, so a migration is needed to make it nullable, OR we use a default/sentinel resident. Making it nullable better reflects reality — during initial release the porteiro may not have this info. A new migration will alter the column.

### 5. New migration required
**Decision**: Migration `000003_alter-visita-morador-nullable` to make `morador_id` nullable in `Visita`
**Rationale**: Current schema enforces `morador_id NOT NULL`, but in practice the porteiro releasing a visitor may not know which apartment/resident the visitor is going to. Making it optional allows the release flow to work without requiring this information up front.

## Data Flow

```
[Porteiro] → GET /porteiros/visitantes?nome=...
  → PorteiroHandler.Search
    → PorteiroService.SearchVisitors
      → VisitorRepository.Search(nome, cpf, telefone, liberado)
        → DB: SELECT * FROM visitante WHERE ...
    ← []VisitorResponseDTO
  ← JSON 200

[Porteiro] → PATCH /porteiros/visitantes/:id/liberar
  → PorteiroHandler.Release (body: { porteiro_id, morador_id? })
    → PorteiroService.ReleaseVisitor
      → DB Transaction:
        1. UPDATE visitante SET liberado = true WHERE id = :id
        2. INSERT INTO Visita (dataEntrada, porteiro_id, visitante_id, morador_id)
           VALUES (NOW(), :porteiro_id, :visitante_id, :morador_id)
    ← VisitorResponseDTO (with liberado = true)
  ← JSON 200
```

## Architecture (Layered)

```
┌────────────────────────────────────────────────────┐
│  PorteiroHandler (internal/server/)                │
│  - search(c *gin.Context)                          │
│  - release(c *gin.Context)                         │
│  interface: porteiroService                        │
├────────────────────────────────────────────────────┤
│  PorteiroService (internal/services/)              │
│  - SearchVisitors(ctx, filter) → []VisitorResponse │
│  - ReleaseVisitor(ctx, id, porteiro_id, morador_id)│
│  Depends on: VisitorRepository, VisitRepository    │
├────────────────────────────────────────────────────┤
│  VisitorRepository + VisitRepository               │
│  (internal/repositories/)                          │
│  - Visitor.Search(ctx, filter) → []Visitor         │
│  - Visitor.Update(ctx, id, fields) → error         │
│  - Visit.Create(ctx, visit) → error                │
├────────────────────────────────────────────────────┤
│  Models (internal/models/)                         │
│  - Visitor (existing)                              │
│  - Visit (new)                                     │
├────────────────────────────────────────────────────┤
│  DTOs (internal/dto/)                              │
│  - VisitorFilterDTO (query params)                 │
│  - ReleaseRequestDTO (request body)                │
│  - VisitorResponseDTO (reused from existing)       │
└────────────────────────────────────────────────────┘
```

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| `morador_id` is NOT NULL in current schema | Migration 000003 makes it nullable; rollback is a standard down migration |
| No auth means any client can release any visitor | Acceptable for now — auth is a separate change. The `porteiro_id` is explicitly required in the request body |
| Race condition: two porteiros releasing same visitor simultaneously | GORM transaction + row-level locking (`FOR UPDATE`) on the visitor row within the transaction |
| Search might return too many results if no filters provided | Require at least one filter parameter; return 400 if all are empty |
