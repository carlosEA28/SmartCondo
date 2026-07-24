## Why

The síndico needs visibility into which residents have overdue payments (inadimplentes) to manage the condominium's finances and enforce reservation rules. The `Pagamento` table already exists in the database but has no API to query it. This change adds a single síndico-only endpoint to list residents with overdue payments and their outstanding balance.

## What Changes

- New `GET /sindico/inadimplentes` endpoint — protected by `RequireSindicoRole` middleware
- Returns a list of delinquent residents with their personal info, apartment details, total overdue amount, and individual overdue payments
- A resident is considered inadimplente if they have at least one `Pagamento` with `status = 'ATRASADO'`
- No new migration needed — the `Pagamento` table already exists

## Capabilities

### New Capabilities
- `sindico-inadimplentes`: Síndico-only view of delinquent residents and their overdue payments

### Modified Capabilities
- *(none)*

## Impact

- **New model**: `Pagamento` in `internal/models/pagamento.go`
- **New DTO**: `internal/dto/inadimplente.go` with `InadimplenteResponseDTO` containing resident info, apartment, total overdue, and payment list
- **New repository**: `internal/repositories/pagamento_repository.go` with `PagamentoRepository` and a `FindInadimplentes(ctx)` method
- **New service**: `internal/services/inadimplente_service.go` with `ListInadimplentes(ctx)` method
- **New handler**: `internal/server/inadimplente_handler.go` with a `list` method
- **Route registration**: `GET /sindico/inadimplentes` in `internal/server/server.go`
- **New errors**: `ErrNenhumInadimplente` if no delinquent residents found (though empty list is fine too)
- **No breaking changes** — existing endpoints untouched
