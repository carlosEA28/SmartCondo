## 1. Database Migration

- [x] 1.1 Create migration `000003_alter-visita-morador-nullable.up.sql` altering `Visita.morador_id` to nullable
- [x] 1.2 Create migration `000003_alter-visita-morador-nullable.down.sql` reverting `morador_id` back to NOT NULL

## 2. Models

- [x] 2.1 Create `internal/models/visit.go` — `Visit` struct mapping `Visita` table with fields: ID, DataEntrada, DataSaida, PorteiroID, VisitanteID, MoradorID

## 3. DTOs

- [x] 3.1 Create `internal/dto/porteiro.go` — `VisitorFilterDTO` with `Nome`, `CPF`, `Telefone`, `Liberado` query params + `ReleaseRequestDTO` with `PorteiroID`, `MoradorID`

## 4. Repository Changes

- [x] 4.1 Add `Search(ctx, filter)` and `UpdateLiberado(ctx, id)` methods to `VisitorRepository` interface + `GormVisitorRepository` implementation in `internal/repositories/visitor_repository.go`
- [x] 4.2 Create `internal/repositories/visit_repository.go` — `VisitRepository` interface with `Create(ctx, visit)` + `GormVisitRepository` implementation

## 5. Error Definitions

- [x] 5.1 Add `ErrVisitNotFound`, `ErrPorteiroNotFound`, `ErrInvalidPorteiroData` sentinel errors to `internal/apperrors/errors.go`

## 6. Service Layer

- [x] 6.1 Create `internal/services/porteiro_service.go` — `PorteiroService` with `SearchVisitors(ctx, filter)` and `ReleaseVisitor(ctx, id, porteiroID, moradorID)` methods, using `VisitorRepository` and `VisitRepository` interfaces

## 7. Handler Layer

- [x] 7.1 Create `internal/server/porteiro_handler.go` — `porteiroHandler` with `search` and `release` methods, following the same patterns as `visitorHandler`

## 8. Route Registration

- [x] 8.1 Add `PorteiroService` and `porteiroHandler` wiring + routes `GET /porteiros/visitantes` and `PATCH /porteiros/visitantes/:id/liberar` in `internal/server/server.go`

## 9. Verification

- [x] 9.1 Run `make build` to ensure project compiles
- [x] 9.2 Run `make test` to ensure existing tests still pass
