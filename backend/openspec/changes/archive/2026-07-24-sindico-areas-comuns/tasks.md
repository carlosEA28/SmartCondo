## 1. Models

- [ ] 1.1 Create `internal/models/areacomum.go` — `AreaComum` struct mapping `areacomum` table
- [ ] 1.2 Create `internal/models/reserva.go` — `Reserva` struct mapping `reserva` table with FKs to AreaComum and User
- [ ] 1.3 Create `internal/models/pagamento.go` — `Pagamento` struct mapping `pagamento` table with status and morador FK

## 2. DTOs

- [ ] 2.1 Create `internal/dto/areacomum.go` — `CreateAreaComumDTO`, `UpdateAreaComumDTO`, `AreaComumResponseDTO`
- [ ] 2.2 Create `internal/dto/reserva.go` — `CreateReservaDTO`, `ReservaResponseDTO`

## 3. Error Definitions

- [ ] 3.1 Add errors to `internal/apperrors/errors.go`: `ErrAreaComumNotFound`, `ErrAreaComumAlreadyExists`, `ErrReservaNotFound`, `ErrReservaConflito`, `ErrMoradorInadimplente`

## 4. Repositories

- [ ] 4.1 Create `internal/repositories/areacomum_repository.go` — `AreaComumRepository` interface + `GormAreaComumRepository` with FindByID, FindAll, FindByName, Create, Update, Delete
- [ ] 4.2 Create `internal/repositories/reserva_repository.go` — `ReservaRepository` interface + `GormReservaRepository` with FindByID, FindAll, FindConflicting, Create, UpdateStatus
- [ ] 4.3 Create `internal/repositories/pagamento_repository.go` — `PagamentoRepository` interface with `HasOverduePayments` + `GormPagamentoRepository`

## 5. Service Layer

- [ ] 5.1 Create `internal/services/areacomum_service.go` — `AreaComumService` with Create (unique name), List, Get, Update, Delete
- [ ] 5.2 Create `internal/services/reserva_service.go` — `ReservaService` with Create (conflict + payment checks), List, Get, Confirmar, Cancelar

## 6. Handler Layer

- [ ] 6.1 Create `internal/server/areacomum_handler.go` — `areaComumHandler` with CRUD methods
- [ ] 6.2 Create `internal/server/reserva_handler.go` — `reservaHandler` with create, list, getByID, confirmar, cancelar

## 7. Route Registration & Wiring

- [ ] 7.1 Add `areaComumRepository`, `reservaRepository`, `pagamentoRepository` to Server struct + constructor
- [ ] 7.2 Register routes in `server.go`: `/sindico/areas-comuns` and `/sindico/reservas` with middleware protection

## 8. Validation

- [ ] 8.1 Run `make build` to ensure project compiles
- [ ] 8.2 Run `make test` to ensure existing tests pass
