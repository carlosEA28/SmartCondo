# SmartCondo

## Project Structure

- **backend/** — PostgreSQL database (Docker) + SQL schema
- **frontend/** — Astro 7 web application
- **docs/** — Project documentation and diagrams

## Database

Start PostgreSQL:
```bash
cd backend && docker compose up -d
```

- Host: `localhost:5432`
- User: `postgres` / Password: `postgres`
- Database: `smartcondo`

Schema is auto-loaded from `backend/database/schema.sql`.

## Frontend

```bash
cd frontend
npm install
npm run dev
```

Node.js >=22.12.0 required. See `frontend/AGENTS.md` for Astro-specific guidance.

## Database Tables

| Table | Purpose |
|-------|---------|
| `Apartamento` | Units (number + block) |
| `Usuario` | Users (types: ADMINISTRADOR, PORTEIRO, MORADOR, SINDICO) |
| `Pagamento` | Payments (status: PENDENTE, PAGO, ATRASADO) |
| `AreaComum` | Common areas for reservations |
| `Reserva` | Reservations (status: PENDENTE, CONFIRMADA, CANCELADA) |
| `Visitante` | Visitor data |
| `Visita` | Entry/exit log |
| `Comunicado` | Announcements |
| `Notificacao` | Email/SMS notifications |

## Key Constraints

- Users with `tipo = 'MORADOR'` must have `apartamento_id` set
- `responsavel = TRUE` only allowed for MORADOR type
- Visitante CPF must be unique (11 digits)
