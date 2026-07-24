## Context

The `Comunicado` table already exists in the database (migration 000001) with columns: `id`, `titulo`, `descricao`, `dataPublicacao`, `sindico_id` (FK → Usuario). No Go code exists for it — no model, repository, service, or handler. The síndico role (`RoleSindico`) is defined as a constant in `models/user.go`.

To secure operations without full JWT authentication, a header-based authentication middleware (`RequireSindicoRole`) will validate the user's role directly against the database using seeded records before granting access to management endpoints.

## Goals / Non-Goals

**Goals:**

- Implement `RequireSindicoRole` middleware to authenticate and authorize Síndico users via header (e.g., `X-User-ID`).
- Add `POST /sindico/comunicados` (protected) to publish an announcement (payload: `titulo`, `descricao`).
- Add `GET /sindico/comunicados` (public) to list all announcements ordered by most recent.
- Add `GET /sindico/comunicados/:id` (public) to view a single announcement.
- Add `DELETE /sindico/comunicados/:id` (protected) to remove an announcement.
- Extract `sindico_id` from request context (`c.Request.Context()`) injected by the middleware.

**Non-Goals:**

- No JWT token signing or verification (uses direct database lookup of user header until JWT phase).
- No update endpoint (announcements are published once; delete + recreate pattern applies).
- No soft-delete or status field.
- No pagination on the list endpoint.

## Decisions

### 1. Header-based Middleware (`RequireSindicoRole`)

**Decision**: Intercept write requests, read user identification header, verify the user exists with `RoleSindico` in the DB, and attach user/ID to the Gin context.
**Rationale**: Eliminates the need to pass `sindico_id` in request bodies or query params, creating a cleaner API contract that seamlessly transitions to JWT in the future.

### 2. DTOs reflect authenticated context

**Decision**: `CreateComunicadoDTO` contains only `titulo` and `descricao`. `sindico_id` is supplied internally via context.
**Rationale**: Prevents clients from spoofing author IDs in request payloads.

### 3. `dataPublicacao` is set server-side

**Decision**: The service sets `time.Now()` on creation; request DTOs omit timestamps.
**Rationale**: Guarantees server timestamp accuracy and satisfies DB NOT NULL constraints.

### 4. Repository: `FindByID`, `FindAll`, `Create`, `Delete`

**Decision**: Standard Gorm repository interface implementation following existing project conventions (`VisitRepository`, `UserRepository`).

## Architecture

┌──────────────────────────────────────────────────────┐
│ RequireSindicoRole Middleware (internal/server/) │
│ - Reads X-User-ID header │
│ - Validates RoleSindico via UserRepository │
│ - Sets user info in c.Set("user", user) │
├──────────────────────────────────────────────────────┤
│ ComunicadoHandler (internal/server/) │
│ - create(c *gin.Context) → Extracts ID from context │
│ - list(c *gin.Context) │
│ - getByID(c *gin.Context) │
│ - delete(c *gin.Context) → Extracts ID from context │
│ interface: comunicadoService │
├──────────────────────────────────────────────────────┤
│ ComunicadoService (internal/services/) │
│ - PublishComunicado(ctx, sindicoID, dto) │
│ - ListComunicados(ctx) → []ComunicadoResponse │
│ - GetComunicado(ctx, id) → ComunicadoResponse │
│ - DeleteComunicado(ctx, id, sindicoID) │
│ Depends on: ComunicadoRepository │
├──────────────────────────────────────────────────────┤
│ ComunicadoRepository (internal/repositories/) │
│ - FindByID(ctx, id) → \*Comunicado │
│ - FindAll(ctx) → []Comunicado │
│ - Create(ctx, comunicado) → error │
│ - Delete(ctx, id, sindicoID) → error │
├──────────────────────────────────────────────────────┤
│ Model & DTOs │
│ - Model: internal/models/comunicado.go │
│ - DTO: CreateComunicadoDTO (titulo, descricao) │
└──────────────────────────────────────────────────────┘

## Risks / Trade-offs

| Risk                                                                                     | Mitigation                                                                                                   |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| Header-based auth can be spoofed in raw HTTP requests without network/gateway protection | Acceptable during initial development with seed data; designed to be replaced 1:1 by JWT verification later. |
| DB lookup per protected request in middleware                                            | Lightweight query by Primary Key; caching or transition to JWT payload will resolve overhead in production.  |
| Preloading `Sindico` association on list endpoint                                        | Done selectively using `Preload("Sindico")` to fetch author names without extra queries.                     |
