## 1. Models

- [x] 1.1 Create `internal/models/areacomum.go` — `AreaComum` struct mapping `areas_comuns` table
- [x] 1.2 Create `internal/models/reserva.go` — `Reserva` struct mapping `reservas` table

## 2. Repository: add HasOverduePayments

- [x] 2.1 Add `HasOverduePayments(ctx, moradorID uuid.UUID) (bool, error)` to `PagamentoRepository` — queries `pagamentos WHERE morador_id = ? AND status = 'ATRASADO' LIMIT 1`

## 3. Errors

- [x] 3.1 Add `ErrAreaComumNotFound`, `ErrReservaNotFound`, `ErrReservaConflito`, `ErrMoradorInadimplente`, `ErrReservaNotOwner`, `ErrInvalidReservaData` to error package

## 4. Middleware: RequireMoradorRole

- [x] 4.1 Add `RequireMoradorRole` middleware in `internal/server/middleware/auth.go` (follow `RequireSindicoRole` pattern)

## 5. Repository: AreaComumRepository

- [x] 5.1 Create `internal/repositories/area_comum_repository.go` — `AreaComumRepository` interface + `GormAreaComumRepository` with `FindByID` and `FindAll`

## 6. Repository: ReservaRepository

- [x] 6.1 Create `internal/repositories/reserva_repository.go` — `ReservaRepository` interface + `GormReservaRepository` with `FindByID`, `FindByMorador`, `FindConflicting`, `Create`, `UpdateStatus`

## 7. DTOs

- [x] 7.1 Create `internal/dto/areacomum.go` — `AreaComumResponse`
- [x] 7.2 Create `internal/dto/reserva.go` — `CreateReservaRequest`, `ReservaResponse`

## 8. Service: AreaComumService

- [x] 8.1 Create `internal/services/area_comum_service.go` — `AreaComumService` with `ListAreas`

## 9. Service: ReservaService

- [x] 9.1 Create `internal/services/reserva_service.go` — `ReservaService` with `CreateReserva`, `ListMinhasReservas`, `CancelReserva`

## 10. Handlers

- [x] 10.1 Create `internal/server/area_comum_handler.go` — `areaComumHandler` with `listAreas`
- [x] 10.2 Create `internal/server/reserva_handler.go` — `reservaHandler` with `create`, `listMinhas`, `cancelar`

## 11. Route registration & wiring

- [x] 11.1 Add repositories and services to Server struct + constructor in `internal/server/server.go`
- [x] 11.2 Register morador routes group with `RequireMoradorRole` middleware in `internal/server/server.go`

## 12. Validation

- [x] 12.1 Run `make build` to ensure project compiles
- [x] 12.2 Run `make test` to ensure existing tests pass
