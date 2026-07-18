## Project: SmartCondo

## Dev environment tips
- Run `make run` to run the backend locally.
- Run `make test` to run all the tests in the backend 
- Run `make migrate-create` to generate a new sql migration file
- Run `make migrate-up` to run the migrations files
- Run `make migrate-down` to run down the migrations
- Run `make deps` to install the project dependencies or a dependecie that is being used but not installed
- Run `make build` to build the application

## Do

- Follow idiomatic Go.
- Keep functions small and focused.
- Prefer simple solutions.
- Write clean and readable code.
- Reuse existing code whenever possible.
- Keep changes minimal and scoped to the requested task.
- Ask for clarification instead of making assumptions.

## Don't
- do not hard code anything
- do not execute any migration up or down without approvence
- do not commit changes before my review
- do not create any other folder or file unless i approve

## Project Structure
- `cmd/api/main.go` — Main file of the application, where everything its loaded an the run the app

- `db/migrations` — Where are found all the migration files of the application

- **`internal/config/`** — Application configuration. Loads environment variables from `.env` into typed structs (`ServerConfig`, `DatabaseConfig`) using `godotenv`. Provides defaults for port, Gin mode, and database connection fields.

- **`internal/database/`** — Database connection layer. Establishes a PostgreSQL connection via GORM using the config values. Returns a `*gorm.DB` instance used throughout the application.

- **`internal/dto/`** — Data Transfer Objects. Intended to hold request/response payload structs (e.g., `CreateUserDTO`, `LoginRequest`) that decouple internal domain models from the API layer. Currently empty.

- **`internal/interfaces/`** — Interface definitions. Intended to define contracts between layers (repository interfaces, service interfaces, provider interfaces) to enable dependency injection and testability. Currently empty.

- **`internal/logger/`** — Structured logging. Configures a global `zerolog` logger with console output in debug mode and JSON output in release mode. Uses RFC3339 timestamp format.

- **`internal/models/`** — Domain models / entities. Intended to contain GORM model structs representing database tables (e.g., `User`, `Condo`, `Resident`, `Unit`, `Payment`). Currently empty.

- **`internal/providers/`** — External service providers. Intended to hold integrations with third-party APIs (payment gateways, email/SMS services, OAuth providers). Currently empty.

- **`internal/repositories/`** — Data access layer. Intended to contain database query implementations (`FindByID`, `Create`, `Update`, `Delete`) that wrap GORM calls for each model. Currently empty.

- **`internal/server/`** — HTTP server and router setup. Creates and configures the Gin engine with middleware (logger, recovery, CORS), defines routes, and includes a `GET /health` endpoint. Handles graceful shutdown wiring.

- **`internal/services/`** — Business logic layer. Intended to contain core application logic (user registration, authentication, condo management, payment processing). Services depend on repository interfaces, not concrete implementations. Currently empty.

- **`internal/tests/`** — Test utilities and helpers. Intended to hold shared test utilities, mock implementations, fixtures, and database setup/teardown helpers. Currently empty.

- **`internal/utils/`** — Shared utility functions. Intended to contain general-purpose helpers (password hashing, JWT handling, pagination, input sanitization, UUID generation). Currently empty.

## Architecture Rules

- Keep a clear separation between HTTP, business logic and persistence.
- Handlers must not contain business logic.
- Services must not depend on Gin.
- Repositories must only contain database operations.
- DTOs must never be used as database models.
- Domain models should not contain HTTP-specific tags unless required.
- Depend on interfaces instead of concrete implementations whenever possible.

## Database

- Never perform database queries directly from handlers.
- All database access must go through repositories.
- Always create migrations for schema changes.
- Never modify old migrations.
- Create a new migration for every schema change.

## API

- Return JSON only.
- Use proper HTTP status codes.
- Validate request payloads before calling services.
- Keep handlers thin.

## Security

- Never hardcode secrets.
- Read configuration from environment variables.
- Validate every user input.
- Never trust client-provided IDs.

## Workflow

When implementing a feature:

1. Update or create database migrations if needed.
2. Create or update domain models.
3. Implement repository methods.
4. Implement business logic in services.
5. Add HTTP handlers.
6. Register routes.
7. Write tests.
8. Run formatting.

## Avoid

- Do not introduce unnecessary packages.
- Do not duplicate business logic.
- Do not place SQL inside handlers.
- Do not bypass repositories.
- Do not create utility functions unless they are reusable.

## Logging

- Use the global zerolog logger.
- Never use fmt.Println for application logs.
- Do not log passwords, tokens or secrets.

## Error Handling

- Never ignore returned errors.
- Wrap errors with context using fmt.Errorf(... %w ...).
- Return meaningful errors.
- Do not panic except during application startup.
