## Why

The síndico needs to publish announcements (comunicados) to keep residents informed about building works, maintenance, and general notices. While JWT authentication is not yet implemented, seeded users (including síndico and porteiro) already exist in the database.

To ensure that only authorized users can publish or delete announcements without dirtying request payloads, a lightweight header-based authentication middleware is required to validate the user against the database before granting access.

## What Changes

- **New Middleware (`RequireSindicoRole`)**: Intercepts requests, reads the user identification header (e.g., `X-User-ID`), checks the database to verify that the user exists and holds the **Síndico** role, and injects the authenticated user into the request context (`r.Context()`).
- **Protected Endpoints**:
  - `POST /sindico/comunicados` — Restricted to Síndico users via middleware. Creates a new announcement.
  - `DELETE /sindico/comunicados/:id` — Restricted to Síndico users via middleware. Removes an announcement.
- **Public / Resident Endpoints**:
  - `GET /sindico/comunicados` — Public/Resident access. Lists all announcements ordered by most recent.
  - `GET /sindico/comunicados/:id` — Public/Resident access. Retrieves details for a specific announcement.
- **Database**: No new migration required (`Comunicado` table and seeded users already exist).

## Capabilities

### New Capabilities

- `sindico-auth-middleware`: Header-based authorization middleware enforcing Síndico access control using database user/role verification.
- `sindico-comunicados`: Full management of announcements, enforcing role authorization for write operations while keeping read operations accessible.

### Modified Capabilities

- _(none)_

## Impact

- **New Middleware**: `internal/server/middleware/auth.go` containing header parsing, database lookup, and role enforcement.
- **New Errors**: `ErrMissingAuthHeader`, `ErrUnauthorizedSindico` (or `ErrForbidden`) in `internal/apperrors/errors.go`.
- **New Domain Models & DTOs**:
  - Model in `internal/models/comunicado.go`
  - Structs in `internal/dto/comunicado.go` (excluding `sindico_id` from payload since it is inferred from `r.Context()`)
- **New Layer Implementations**:
  - `internal/repositories/comunicado_repository.go`
  - `internal/services/comunicado_service.go`
  - `internal/server/comunicado_handler.go`
- **Route Registration**: Updated in `internal/server/server.go` to wrap management routes with `RequireSindicoRole`.
