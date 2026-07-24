# Tasks: Publicar Comunicado (Síndico)

## 1. Middleware de Autenticação / Autorização

- [x] 1.1 Criar middleware `RequireSindicoRole` em `internal/server/middleware/auth.go` (ou `sindico_auth.go`) para ler o header (ex: `X-User-ID`), validar no `UserRepository` se o usuário existe e possui a role `RoleSindico`, e injetar o usuário/ID no contexto do Gin (`c.Set("user", user)`).

## 2. Models

- [x] 2.1 Criar `internal/models/comunicado.go` — struct `Comunicado` mapeando a tabela `comunicado` com os campos: `ID`, `Titulo`, `Descricao`, `DataPublicacao`, `SindicoID` e a associação `Sindico` (`User`).

## 3. DTOs

- [x] 3.1 Criar `internal/dto/comunicado.go`:
  - `CreateComunicadoDTO` contendo apenas `Titulo` e `Descricao` (o `SindicoID` é injetado via contexto).
  - `ComunicadoResponseDTO` com todos os campos do comunicado + `SindicoNome`.

## 4. Error Definitions

- [x] 4.1 Adicionar sentinel errors em `internal/apperrors/errors.go`: `ErrComunicadoNotFound`, `ErrInvalidComunicadoData`, `ErrMissingAuthHeader` e `ErrUnauthorizedSindico`.

## 5. Repository

- [x] 5.1 Criar `internal/repositories/comunicado_repository.go` — interface `ComunicadoRepository` com `FindByID`, `FindAll`, `Create`, `Delete` + implementação `GormComunicadoRepository`.

## 6. Service Layer

- [x] 6.1 Criar `internal/services/comunicado_service.go` — `ComunicadoService` contendo os métodos:
  - `PublishComunicado(ctx, sindicoID, dto)`
  - `ListComunicados(ctx)`
  - `GetComunicado(ctx, id)`
  - `DeleteComunicado(ctx, id, sindicoID)`

## 7. Handler Layer

- [x] 7.1 Criar `internal/server/comunicado_handler.go` — `comunicadoHandler` com os métodos `create`, `list`, `getByID` e `delete`, extraindo o `sindico_id` diretamente do contexto nas ações protegidas (`create` e `delete`).

## 8. Route Registration & Wiring

- [x] 8.1 Adicionar `comunicadoRepository` à struct `Server` e ao construtor `New()` em `internal/server/server.go`.
- [x] 8.2 Registrar as rotas em `internal/server/server.go`:
  - Rotas protegidas pelo middleware `RequireSindicoRole`: `POST /sindico/comunicados` e `DELETE /sindico/comunicados/:id`.
  - Rotas públicas: `GET /sindico/comunicados` e `GET /sindico/comunicados/:id`.

## 9. Validation & Testing

- [x] 9.1 Rodar `make build` para garantir que o projeto compila sem erros.
- [x] 9.2 Rodar `make test` para verificar se a suíte de testes passou com sucesso.
