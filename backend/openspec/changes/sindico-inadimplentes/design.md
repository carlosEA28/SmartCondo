## Context

The `Pagamento` table exists in the database (migration 000001) with columns: `id`, `valor`, `vencimento`, `dataPagamento`, `status` (PENDENTE/PAGO/ATRASADO), `morador_id` (FK to Usuario). There is zero Go code for it. The `RequireSindicoRole` middleware already exists and will protect the endpoint.

## Goals / Non-Goals

**Goals:**
- Add `GET /sindico/inadimplentes` returning residents with at least one `Pagamento.status = 'ATRASADO'`
- Return resident info (name, email, phone, apartment), total overdue amount, and individual overdue payments
- Protect with `RequireSindicoRole` middleware

**Non-Goals:**
- No CRUD for payments (this is read-only)
- No payment status change (this is a separate future feature)
- No filtering or pagination (reasonable for condominium scale)

## Decisions

### 1. Repository query: JOIN Pagamento + Usuario + Apartamento
**Decision**: `PagamentoRepository.FindInadimplentes` does a single query joining `pagamento` with `usuario` and `apartamento`, grouping by morador_id, filtering where `status = 'ATRASADO'`, and aggregating the total.
**Rationale**: Single query is more efficient than N+1 queries. GORM's `Preload` and `Joins` can handle this.

### 2. Response structure: grouped by resident
**Decision**: Return a list of `InadimplenteResponseDTO` where each item has the resident's info + list of overdue payments + total amount.
**Rationale**: Clearer than flattening payments. The síndico sees who owes and what specific payments are overdue.

### 3. Empty result: return empty array
**Decision**: If no residents are delinquent, return `[]` with status 200, not 404.
**Rationale**: Follows REST conventions. An empty list means "no delinquent residents" which is a valid, successful response.

## Response Structure

```json
[
  {
    "morador": {
      "id": "uuid",
      "full_name": "João Silva",
      "email": "joao@email.com",
      "phone": "11999999999",
      "apartment": { "id": "uuid", "number": 101, "block": "A" }
    },
    "total_overdue": 1250.50,
    "payments": [
      {
        "id": "uuid",
        "valor": 625.25,
        "vencimento": "2026-06-10",
        "status": "ATRASADO"
      }
    ]
  }
]
```

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| No pagination could be slow with many delinquent residents | Acceptable for condominium scale (max hundreds of residents). Can add pagination later |
| Simple definition of "inadimplente" = any overdue payment | Matches current requirements. Can refine later (e.g., overdue > 30 days, minimum amount) |
