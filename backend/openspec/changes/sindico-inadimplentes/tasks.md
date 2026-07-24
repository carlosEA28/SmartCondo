## 1. Model

- [x] 1.1 Create `internal/models/pagamento.go` — `Pagamento` struct mapping `pagamento` table with fields: ID, Valor, Vencimento, DataPagamento, Status, MoradorID, Morador (User association)

## 2. DTO

- [x] 2.1 Create `internal/dto/inadimplente.go` — `InadimplenteResponseDTO` with `Morador` (UserResponseDTO), `TotalOverdue`, `Payments` (list of overdue payment summaries)

## 3. Repository

- [x] 3.1 Create `internal/repositories/pagamento_repository.go` — `PagamentoRepository` interface with `FindInadimplentes(ctx)` + `GormPagamentoRepository` implementation that queries payments with status ATRASADO grouped by morador

## 4. Service

- [x] 4.1 Create `internal/services/inadimplente_service.go` — `InadimplenteService` with `ListInadimplentes(ctx)` method using `PagamentoRepository`

## 5. Handler

- [x] 5.1 Create `internal/server/inadimplente_handler.go` — `inadimplenteHandler` with `list` method, protected by `RequireSindicoRole` middleware

## 6. Route Registration & Wiring

- [x] 6.1 Add `pagamentoRepository` to Server struct + constructor in `internal/server/server.go`
- [x] 6.2 Register `GET /sindico/inadimplentes` route with middleware in `internal/server/server.go`

## 7. Validation

- [x] 7.1 Run `make build` to ensure project compiles
- [x] 7.2 Run `make test` to ensure existing tests pass
