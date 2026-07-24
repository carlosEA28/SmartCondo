# API Reference — SmartCondo

> **Versão:** 3.0.0  
> **Base URL:** `http://localhost:8080`  
> **Formato:** JSON  
> **Autenticação:** Header-based (X-User-ID) para rotas de síndico

---

## Sumário

- [Health Check & Docs](#health-check--docs)
- [Usuários](#usuários)
  - [Criar Usuário](#criar-usuário)
  - [Listar Usuários](#listar-usuários)
  - [Obter Usuário por ID](#obter-usuário-por-id)
  - [Atualizar Usuário](#atualizar-usuário)
  - [Excluir Usuário](#excluir-usuário)
- [Visitantes](#visitantes)
  - [Criar Visitante](#criar-visitante)
  - [Listar Visitantes](#listar-visitantes)
  - [Obter Visitante por ID](#obter-visitante-por-id)
  - [Excluir Visitante](#excluir-visitante)
- [Porteiro](#porteiro)
  - [Pesquisar Visitantes](#pesquisar-visitantes)
  - [Liberar Visitante](#liberar-visitante)
- [Síndico — Comunicados](#síndico--comunicados)
  - [Publicar Comunicado](#publicar-comunicado)
  - [Listar Comunicados](#listar-comunicados)
  - [Obter Comunicado por ID](#obter-comunicado-por-id)
  - [Excluir Comunicado](#excluir-comunicado)
- [Síndico — Inadimplentes](#síndico--inadimplentes)
  - [Listar Inadimplentes](#listar-inadimplentes)
- [Erros](#erros)
- [Modelos de Dados](#modelos-de-dados)
- [Notas Técnicas](#notas-técnicas)

---

## Health Check & Docs

### `GET /health`

Verifica se o servidor está ativo. Usado por balanceadores de carga e monitoramento.

#### Requisição

Nenhum body, nenhum parâmetro.

#### Respostas

**`200 OK`**
```json
{
  "status": "ok"
}
```

#### Exemplos

```bash
curl -X GET http://localhost:8080/health
```

```javascript
fetch('http://localhost:8080/health')
  .then(res => res.json())
  .then(console.log);
```

```go
resp, err := http.Get("http://localhost:8080/health")
```

### `GET /` e `GET /docs`

Exibe a documentação interativa da API via Swagger UI.

### `GET /docs/openapi.yaml`

Retorna o arquivo OpenAPI Specification (YAML) estático.

---

## Usuários

### Criar Usuário

### `POST /users`

Cria um novo morador com vínculo a um apartamento. O usuário é criado com role `MORADOR` e status `ATIVO`.

O cadastro no AWS Cognito é feito de forma **não crítica** — falhas no Cognito não impedem a criação local.

#### Validações

- Email deve ser único no sistema
- Apartamento (número + bloco) deve ser único
- Senha deve ter entre 8 e 72 caracteres
- Telefone é validado como número brasileiro (formato `(XX) XXXXX-XXXX`)
- O campo `apartment` é obrigatório

#### Request Body

```json
{
  "full_name": "Maria Silva",
  "email": "maria@example.com",
  "password": "password123",
  "phone": "11999999999",
  "responsible": false,
  "apartment": {
    "number": 101,
    "block": "A"
  }
}
```

| Campo          | Tipo    | Obrigatório | Descrição                          |
|----------------|---------|-------------|------------------------------------|
| `full_name`    | string  | sim         | Nome completo. Máx. 100 caracteres |
| `email`        | string  | sim         | Email válido. Máx. 100 caracteres  |
| `password`     | string  | sim         | Senha. Mín. 8, máx. 72 caracteres  |
| `phone`        | string  | sim         | Telefone. Máx. 15 caracteres       |
| `responsible`  | boolean | não         | Responsável pelo apartamento       |
| `apartment`    | object  | sim         | Dados do apartamento               |
| `apartment.number` | number | sim      | Número do apto. Deve ser > 0       |
| `apartment.block`  | string | sim      | Bloco. Máx. 10 caracteres          |

#### Respostas

**`201 Created`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "full_name": "Maria Silva",
  "email": "maria@example.com",
  "phone": "(11) 99999-9999",
  "status": "ATIVO",
  "role": "MORADOR",
  "responsible": false,
  "apartment": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": 101,
    "block": "A"
  }
}
```

#### Erros

| Código | Body                                     | Motivo                            |
|--------|------------------------------------------|-----------------------------------|
| 400    | `{"error": "invalid request body"}`      | JSON inválido ou campos ausentes  |
| 409    | `{"error": "user already exists"}`       | Email já cadastrado               |
| 409    | `{"error": "apartment already registered"}` | Apartamento já cadastrado      |
| 422    | `{"error": "apartment is required for residents"}` | Apartamento não informado |
| 500    | `{"error": "failed to create user"}`     | Erro interno do servidor          |

#### Exemplos

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Maria Silva",
    "email": "maria@example.com",
    "password": "password123",
    "phone": "11999999999",
    "apartment": {
      "number": 101,
      "block": "A"
    }
  }'
```

```javascript
const response = await fetch('http://localhost:8080/users', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    full_name: 'Maria Silva',
    email: 'maria@example.com',
    password: 'password123',
    phone: '11999999999',
    apartment: { number: 101, block: 'A' }
  })
});
const data = await response.json();
```

```go
type CreateUserRequest struct {
    FullName    string      `json:"full_name"`
    Email       string      `json:"email"`
    Password    string      `json:"password"`
    Phone       string      `json:"phone"`
    Responsible bool        `json:"responsible"`
    Apartment   ApartmentReq `json:"apartment"`
}

type ApartmentReq struct {
    Number int    `json:"number"`
    Block  string `json:"block"`
}

body := CreateUserRequest{
    FullName: "Maria Silva",
    Email:    "maria@example.com",
    Password: "password123",
    Phone:    "11999999999",
    Apartment: ApartmentReq{Number: 101, Block: "A"},
}

b, _ := json.Marshal(body)
resp, _ := http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(b))
```

---

### Listar Usuários

### `GET /users`

Retorna todos os usuários cadastrados, ordenados por nome (ASC).

#### Requisição

Nenhum body, nenhum parâmetro.

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "full_name": "Maria Silva",
    "email": "maria@example.com",
    "phone": "(11) 99999-9999",
    "status": "ATIVO",
    "role": "MORADOR",
    "responsible": false,
    "apartment": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "number": 101,
      "block": "A"
    }
  }
]
```

#### Erros

| Código | Body                                 | Motivo                    |
|--------|--------------------------------------|---------------------------|
| 500    | `{"error": "failed to list users"}`  | Erro interno do servidor  |

#### Exemplos

```bash
curl -X GET http://localhost:8080/users
```

```javascript
const response = await fetch('http://localhost:8080/users');
const users = await response.json();
```

```go
resp, err := http.Get("http://localhost:8080/users")
var users []UserResponseDTO
json.NewDecoder(resp.Body).Decode(&users)
```

---

### Obter Usuário por ID

### `GET /users/:id`

Retorna um único usuário pelo seu UUID.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição           |
|-----------|--------|-------------|---------------------|
| `id`      | string | sim         | UUID do usuário     |

#### Respostas

**`200 OK`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "full_name": "Maria Silva",
  "email": "maria@example.com",
  "phone": "(11) 99999-9999",
  "status": "ATIVO",
  "role": "MORADOR",
  "responsible": false,
  "apartment": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": 101,
    "block": "A"
  }
}
```

#### Erros

| Código | Body                                    | Motivo                           |
|--------|-----------------------------------------|----------------------------------|
| 400    | `{"error": "invalid user id"}`          | UUID mal formatado               |
| 404    | `{"error": "user not found"}`           | Usuário não encontrado           |
| 500    | `{"error": "failed to get user"}`       | Erro interno do servidor         |

#### Exemplos

```bash
curl -X GET http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

```javascript
const response = await fetch('http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479');
const user = await response.json();
```

```go
id := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s", id))
var user UserResponseDTO
json.NewDecoder(resp.Body).Decode(&user)
```

---

### Atualizar Usuário

### `PUT /users/:id`

Atualiza o nome e telefone de um usuário existente.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição           |
|-----------|--------|-------------|---------------------|
| `id`      | string | sim         | UUID do usuário     |

#### Request Body

```json
{
  "full_name": "Maria Santos",
  "phone": "11888888888"
}
```

| Campo       | Tipo    | Obrigatório | Descrição                          |
|-------------|---------|-------------|------------------------------------|
| `full_name` | string  | sim         | Novo nome completo. Máx. 100 chars |
| `phone`     | string  | sim         | Novo telefone. Máx. 15 caracteres  |

#### Respostas

**`200 OK`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "full_name": "Maria Santos",
  "email": "maria@example.com",
  "phone": "(11) 88888-8888",
  "status": "ATIVO",
  "role": "MORADOR",
  "responsible": false,
  "apartment": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "number": 101,
    "block": "A"
  }
}
```

#### Erros

| Código | Body                                      | Motivo                            |
|--------|-------------------------------------------|-----------------------------------|
| 400    | `{"error": "invalid user id"}`            | UUID mal formatado                |
| 400    | `{"error": "invalid request body"}`       | JSON inválido ou campos ausentes  |
| 404    | `{"error": "user not found"}`             | Usuário não encontrado            |
| 422    | `{"error": "invalid user data"}`          | Nome ou telefone vazio após trim  |
| 500    | `{"error": "failed to update user"}`      | Erro interno do servidor          |

#### Exemplos

```bash
curl -X PUT http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Maria Santos",
    "phone": "11888888888"
  }'
```

```javascript
const response = await fetch('http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    full_name: 'Maria Santos',
    phone: '11888888888'
  })
});
const user = await response.json();
```

```go
body := UpdateUserRequest{
    FullName: "Maria Santos",
    Phone:    "11888888888",
}
b, _ := json.Marshal(body)
url := fmt.Sprintf("http://localhost:8080/users/%s", id)
req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
req.Header.Set("Content-Type", "application/json")
resp, _ := http.DefaultClient.Do(req)
```

---

### Excluir Usuário

### `DELETE /users/:id`

Remove um usuário do sistema. A exclusão no AWS Cognito é tentada de forma **não crítica**.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição           |
|-----------|--------|-------------|---------------------|
| `id`      | string | sim         | UUID do usuário     |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body                                        | Motivo                            |
|--------|---------------------------------------------|-----------------------------------|
| 400    | `{"error": "invalid user id"}`              | UUID mal formatado                |
| 404    | `{"error": "user not found"}`               | Usuário não encontrado            |
| 409    | `{"error": "user has related records"}`     | Usuário possui registros atrelados |
| 500    | `{"error": "failed to delete user"}`        | Erro interno do servidor          |

#### Exemplos

```bash
curl -X DELETE http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

```javascript
const response = await fetch('http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479', {
  method: 'DELETE'
});
// response.status === 204
```

```go
url := fmt.Sprintf("http://localhost:8080/users/%s", id)
req, _ := http.NewRequest(http.MethodDelete, url, nil)
resp, _ := http.DefaultClient.Do(req)
// resp.StatusCode == http.StatusNoContent
```

---

## Visitantes

### Criar Visitante

### `POST /visitors`

Cadastra um novo visitante no sistema com foto opcional. A foto é enviada como arquivo e armazenada no S3.

#### Request Body (multipart/form-data)

| Campo   | Tipo     | Obrigatório | Descrição                     |
|---------|----------|-------------|-------------------------------|
| `name`  | string   | sim         | Nome completo. Máx. 100 chars |
| `cpf`   | string   | sim         | CPF com 11 dígitos            |
| `phone` | string   | sim         | Telefone. Máx. 15 caracteres  |
| `photo` | file     | não         | Foto do visitante (upload S3) |

#### Respostas

**`201 Created`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "João Visitante",
  "cpf": "12345678901",
  "phone": "(11) 99999-9999",
  "photo": "https://bucket.s3.us-east-1.amazonaws.com/visitors/uuid/photo.jpg",
  "liberado": false
}
```

#### Erros

| Código | Body                                       | Motivo                          |
|--------|--------------------------------------------|----------------------------------|
| 400    | `{"error": "invalid request body"}`        | Formulário inválido              |
| 409    | `{"error": "visitor already exists"}`      | CPF já cadastrado                |
| 500    | `{"error": "failed to create visitor"}`    | Erro interno do servidor         |

#### Exemplos

```bash
curl -X POST http://localhost:8080/visitors \
  -F "name=João Visitante" \
  -F "cpf=12345678901" \
  -F "phone=11999999999" \
  -F "photo=@foto.jpg"
```

---

### Listar Visitantes

### `GET /visitors`

Retorna todos os visitantes cadastrados, ordenados por nome (ASC).

#### Requisição

Nenhum body, nenhum parâmetro.

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "name": "João Visitante",
    "cpf": "12345678901",
    "phone": "(11) 99999-9999",
    "photo": "https://bucket.s3.us-east-1.amazonaws.com/visitors/uuid/photo.jpg",
    "liberado": false
  }
]
```

#### Erros

| Código | Body                                       | Motivo                    |
|--------|--------------------------------------------|---------------------------|
| 500    | `{"error": "failed to list visitors"}`     | Erro interno do servidor  |

---

### Obter Visitante por ID

### `GET /visitors/:id`

Retorna um único visitante pelo seu UUID.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição           |
|-----------|--------|-------------|---------------------|
| `id`      | string | sim         | UUID do visitante   |

#### Respostas

**`200 OK`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "João Visitante",
  "cpf": "12345678901",
  "phone": "(11) 99999-9999",
  "photo": "https://bucket.s3.us-east-1.amazonaws.com/visitors/uuid/photo.jpg",
  "liberado": false
}
```

#### Erros

| Código | Body                                     | Motivo                    |
|--------|------------------------------------------|---------------------------|
| 400    | `{"error": "invalid visitor id"}`        | UUID mal formatado        |
| 404    | `{"error": "visitor not found"}`         | Visitante não encontrado  |
| 500    | `{"error": "failed to get visitor"}`     | Erro interno do servidor  |

---

### Excluir Visitante

### `DELETE /visitors/:id`

Remove um visitante do sistema. A foto associada no S3 também é removida.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição           |
|-----------|--------|-------------|---------------------|
| `id`      | string | sim         | UUID do visitante   |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body                                        | Motivo                     |
|--------|---------------------------------------------|----------------------------|
| 400    | `{"error": "invalid visitor id"}`           | UUID mal formatado         |
| 404    | `{"error": "visitor not found"}`            | Visitante não encontrado   |
| 500    | `{"error": "failed to delete visitor"}`     | Erro interno do servidor   |

---

## Porteiro

### Pesquisar Visitantes

### `GET /porteiros/visitantes`

Pesquisa visitantes por filtros. Pelo menos um filtro deve ser informado.

#### Parâmetros de Query

| Parâmetro  | Tipo    | Obrigatório | Descrição                    |
|------------|---------|-------------|------------------------------|
| `nome`     | string  | não         | Nome (busca parcial ILIKE)   |
| `cpf`      | string  | não         | CPF (busca exata)            |
| `telefone` | string  | não         | Telefone (busca parcial)     |
| `liberado` | boolean | não         | Status de liberação          |

> ⚠️ Pelo menos um dos filtros deve ser informado.

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "name": "João Visitante",
    "cpf": "12345678901",
    "phone": "(11) 99999-9999",
    "photo": "https://bucket.s3.us-east-1.amazonaws.com/visitors/uuid/photo.jpg",
    "liberado": false
  }
]
```

#### Erros

| Código | Body                                                | Motivo                           |
|--------|-----------------------------------------------------|----------------------------------|
| 400    | `{"error": "invalid query parameters"}`              | Parâmetros inválidos             |
| 400    | `{"error": "at least one search filter is required"}` | Nenhum filtro informado         |
| 500    | `{"error": "failed to search visitors"}`             | Erro interno do servidor         |

#### Exemplos

```bash
curl -X GET "http://localhost:8080/porteiros/visitantes?nome=João&liberado=false"
```

---

### Liberar Visitante

### `PATCH /porteiros/visitantes/:id/liberar`

Registra a entrada de um visitante, marcando-o como liberado e criando um registro de visita.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição             |
|-----------|--------|-------------|-----------------------|
| `id`      | string | sim         | UUID do visitante     |

#### Request Body

```json
{
  "porteiro_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "morador_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

| Campo         | Tipo    | Obrigatório | Descrição                        |
|---------------|---------|-------------|----------------------------------|
| `porteiro_id` | string  | sim         | UUID do porteiro que liberou     |
| `morador_id`  | string  | não         | UUID do morador autorizador      |

#### Respostas

**`200 OK`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "João Visitante",
  "cpf": "12345678901",
  "phone": "(11) 99999-9999",
  "photo": "https://bucket.s3.us-east-1.amazonaws.com/visitors/uuid/photo.jpg",
  "liberado": true
}
```

#### Erros

| Código | Body                                             | Motivo                        |
|--------|--------------------------------------------------|-------------------------------|
| 400    | `{"error": "invalid visitor id"}`                | UUID do visitante inválido    |
| 400    | `{"error": "invalid request body"}`              | JSON inválido                 |
| 400    | `{"error": "porteiro not found"}`                | Porteiro não encontrado       |
| 404    | `{"error": "visitor not found"}`                 | Visitante não encontrado      |
| 500    | `{"error": "failed to release visitor"}`         | Erro interno do servidor      |

#### Exemplos

```bash
curl -X PATCH http://localhost:8080/porteiros/visitantes/f47ac10b-58cc-4372-a567-0e02b2c3d479/liberar \
  -H "Content-Type: application/json" \
  -d '{
    "porteiro_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
  }'
```

---

## Síndico — Comunicados

### Autenticação

As rotas de criação e exclusão de comunicados exigem o middleware `RequireSindicoRole`, que valida o usuário através do header:

| Header       | Obrigatório | Descrição                    |
|-------------|-------------|------------------------------|
| `X-User-ID` | sim         | UUID do usuário síndico      |

### Publicar Comunicado

### `POST /sindico/comunicados`

Publica um novo comunicado. O `sindico_id` é extraído automaticamente do header `X-User-ID`.

#### Headers

| Header       | Obrigatório | Descrição                    |
|-------------|-------------|------------------------------|
| `X-User-ID` | sim         | UUID do síndico autenticado  |

#### Request Body

```json
{
  "titulo": "Manutenção no elevador",
  "descricao": "O elevador será desligado para manutenção no dia 15/08."
}
```

| Campo       | Tipo    | Obrigatório | Descrição                        |
|-------------|---------|-------------|----------------------------------|
| `titulo`    | string  | sim         | Título do comunicado. Máx. 100   |
| `descricao` | string  | sim         | Descrição/conteúdo do comunicado |

#### Respostas

**`201 Created`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "titulo": "Manutenção no elevador",
  "descricao": "O elevador será desligado para manutenção no dia 15/08.",
  "dataPublicacao": "2026-07-24T10:00:00Z",
  "sindicoId": "550e8400-e29b-41d4-a716-446655440000",
  "sindicoNome": "Carlos Síndico"
}
```

#### Erros

| Código | Body                                                | Motivo                           |
|--------|-----------------------------------------------------|----------------------------------|
| 400    | `{"error": "invalid request body"}`                 | JSON inválido                    |
| 401    | `{"error": "missing authentication header"}`        | Header X-User-ID ausente         |
| 401    | `{"error": "user is not authorized as sindico"}`    | Usuário não encontrado           |
| 403    | `{"error": "user is not authorized as sindico"}`    | Usuário não é síndico            |
| 500    | `{"error": "failed to create comunicado"}`          | Erro interno do servidor         |

---

### Listar Comunicados

### `GET /sindico/comunicados`

Retorna todos os comunicados ordenados do mais recente para o mais antigo. Rota pública.

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "titulo": "Manutenção no elevador",
    "descricao": "O elevador será desligado para manutenção no dia 15/08.",
    "dataPublicacao": "2026-07-24T10:00:00Z",
    "sindicoId": "550e8400-e29b-41d4-a716-446655440000",
    "sindicoNome": "Carlos Síndico"
  }
]
```

#### Erros

| Código | Body                                                | Motivo                    |
|--------|-----------------------------------------------------|---------------------------|
| 500    | `{"error": "failed to list comunicados"}`           | Erro interno do servidor  |

---

### Obter Comunicado por ID

### `GET /sindico/comunicados/:id`

Retorna um único comunicado pelo seu UUID. Rota pública.

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição               |
|-----------|--------|-------------|-------------------------|
| `id`      | string | sim         | UUID do comunicado      |

#### Respostas

**`200 OK`**
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "titulo": "Manutenção no elevador",
  "descricao": "O elevador será desligado para manutenção no dia 15/08.",
  "dataPublicacao": "2026-07-24T10:00:00Z",
  "sindicoId": "550e8400-e29b-41d4-a716-446655440000",
  "sindicoNome": "Carlos Síndico"
}
```

#### Erros

| Código | Body                                           | Motivo                     |
|--------|------------------------------------------------|----------------------------|
| 400    | `{"error": "invalid comunicado id"}`           | UUID mal formatado         |
| 404    | `{"error": "comunicado not found"}`            | Comunicado não encontrado  |
| 500    | `{"error": "failed to get comunicado"}`        | Erro interno do servidor   |

---

### Excluir Comunicado

### `DELETE /sindico/comunicados/:id`

Remove um comunicado. Apenas o síndico que criou o comunicado pode excluí-lo.

#### Headers

| Header       | Obrigatório | Descrição                    |
|-------------|-------------|------------------------------|
| `X-User-ID` | sim         | UUID do síndico autenticado  |

#### Parâmetros de Path

| Parâmetro | Tipo   | Obrigatório | Descrição               |
|-----------|--------|-------------|-------------------------|
| `id`      | string | sim         | UUID do comunicado      |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body                                                | Motivo                           |
|--------|-----------------------------------------------------|----------------------------------|
| 400    | `{"error": "invalid comunicado id"}`                | UUID mal formatado               |
| 401    | `{"error": "missing authentication header"}`        | Header X-User-ID ausente         |
| 403    | `{"error": "user is not authorized as sindico"}`    | Usuário não é síndico            |
| 403    | `{"error": "you can only delete your own comunicados"}` | Não é o autor do comunicado  |
| 404    | `{"error": "comunicado not found"}`                 | Comunicado não encontrado        |
| 500    | `{"error": "failed to delete comunicado"}`          | Erro interno do servidor         |

---

## Síndico — Inadimplentes

### Listar Inadimplentes

### `GET /sindico/inadimplentes`

Retorna a lista de moradores com pagamentos atrasados, agrupados por morador com o total devido.

#### Headers

| Header       | Obrigatório | Descrição                    |
|-------------|-------------|------------------------------|
| `X-User-ID` | sim         | UUID do síndico autenticado  |

#### Respostas

**`200 OK`**
```json
[
  {
    "morador": {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "full_name": "Maria Silva",
      "email": "maria@example.com",
      "phone": "(11) 99999-9999",
      "status": "ATIVO",
      "role": "MORADOR",
      "responsible": false,
      "apartment": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "number": 101,
        "block": "A"
      }
    },
    "total_overdue": 1500.00,
    "payments": [
      {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
        "valor": 750.00,
        "vencimento": "2026-06-10T00:00:00Z",
        "status": "ATRASADO"
      },
      {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c9",
        "valor": 750.00,
        "vencimento": "2026-07-10T00:00:00Z",
        "status": "ATRASADO"
      }
    ]
  }
]
```

#### Erros

| Código | Body                                                | Motivo                           |
|--------|-----------------------------------------------------|----------------------------------|
| 401    | `{"error": "missing authentication header"}`        | Header X-User-ID ausente         |
| 403    | `{"error": "user is not authorized as sindico"}`    | Usuário não é síndico            |
| 500    | `{"error": "failed to list inadimplentes"}`         | Erro interno do servidor         |

#### Exemplos

```bash
curl -X GET http://localhost:8080/sindico/inadimplentes \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Erros

### Formato Padrão

Todas as respostas de erro seguem o formato:

```json
{
  "error": "mensagem descritiva do erro"
}
```

### Códigos de Erro

| Código | Significado                  |
|--------|------------------------------|
| 400    | Bad Request — requisição mal formatada |
| 401    | Unauthorized — autenticação ausente ou inválida |
| 403    | Forbidden — permissão negada |
| 404    | Not Found — recurso não encontrado |
| 409    | Conflict — conflito (duplicidade, registro em uso) |
| 422    | Unprocessable Entity — dados inválidos |
| 500    | Internal Server Error — erro interno |

### Erros da Aplicação

| Erro (mensagem)                                     | Causa                                                   |
|-----------------------------------------------------|---------------------------------------------------------|
| `invalid request body`                              | JSON mal formatado ou campos obrigatórios ausentes      |
| `invalid user id`                                   | UUID do usuário inválido                                |
| `invalid visitor id`                                | UUID do visitante inválido                              |
| `invalid comunicado id`                             | UUID do comunicado inválido                             |
| `user already exists`                               | Tentativa de cadastro com email já existente            |
| `visitor already exists`                            | Tentativa de cadastro com CPF já existente              |
| `user not found`                                    | Usuário não encontrado pelo ID informado                |
| `visitor not found`                                 | Visitante não encontrado pelo ID informado              |
| `comunicado not found`                              | Comunicado não encontrado pelo ID informado             |
| `user has related records`                          | Tentativa de excluir usuário com registros associados   |
| `apartment is required for residents`               | Morador criado sem informar apartamento                 |
| `apartment already registered`                      | Número e bloco do apartamento já cadastrados            |
| `invalid user data`                                 | Nome ou telefone enviado vazio                          |
| `at least one search filter is required`            | Nenhum filtro informado na pesquisa de visitantes       |
| `porteiro not found`                                | Porteiro não encontrado na liberação                    |
| `missing authentication header`                     | Header X-User-ID ausente em rota protegida              |
| `user is not authorized as sindico`                 | Usuário não possui role de síndico                      |
| `you can only delete your own comunicados`          | Tentativa de excluir comunicado de outro síndico        |
| `failed to create user`                             | Erro interno ao criar usuário                           |
| `failed to get user`                                | Erro interno ao buscar usuário                          |
| `failed to list users`                              | Erro interno ao listar usuários                         |
| `failed to update user`                             | Erro interno ao atualizar usuário                       |
| `failed to delete user`                             | Erro interno ao excluir usuário                         |
| `failed to create visitor`                          | Erro interno ao criar visitante                         |
| `failed to list visitors`                           | Erro interno ao listar visitantes                       |
| `failed to get visitor`                             | Erro interno ao buscar visitante                        |
| `failed to delete visitor`                          | Erro interno ao excluir visitante                       |
| `failed to search visitors`                         | Erro interno ao pesquisar visitantes                    |
| `failed to release visitor`                         | Erro interno ao liberar visitante                       |
| `failed to create comunicado`                       | Erro interno ao criar comunicado                        |
| `failed to list comunicados`                        | Erro interno ao listar comunicados                      |
| `failed to get comunicado`                          | Erro interno ao buscar comunicado                       |
| `failed to delete comunicado`                       | Erro interno ao excluir comunicado                      |
| `failed to list inadimplentes`                      | Erro interno ao listar inadimplentes                    |

---

## Modelos de Dados

### User (Usuário) — Tabela `usuario`

| Campo          | Tipo              | Descrição                                    |
|----------------|-------------------|----------------------------------------------|
| `id`           | string (UUID)     | Identificador único                          |
| `full_name`    | string            | Nome completo (coluna: `nome`)               |
| `email`        | string            | Email único (coluna: `email`)                |
| `password`     | string            | Hash bcrypt da senha (coluna: `senha`)       |
| `phone`        | string            | Telefone formatado (coluna: `telefone`)      |
| `status`       | string            | `ATIVO`, `INATIVO` ou `BLOQUEADO`            |
| `role`         | string            | `MORADOR`, `PORTEIRO` ou `SINDICO`           |
| `apartment_id` | string (UUID)/null| FK para o apartamento (coluna: `apartamento_id`) |
| `apartment`    | object or null    | Dados do apartamento                         |
| `responsible`  | boolean           | Se é o responsável pelo apto (coluna: `responsavel`) |

**Regras de consistência (banco de dados):**
- `responsavel = TRUE` só é permitido quando `role = 'MORADOR'`
- `apartamento_id` só pode ser não-nulo quando `role = 'MORADOR'`

### Apartment (Apartamento) — Tabela `apartamento`

| Campo    | Tipo          | Descrição                     |
|----------|---------------|-------------------------------|
| `id`     | string (UUID) | Identificador único           |
| `number` | integer       | Número do apartamento         |
| `block`  | string        | Bloco (máx. 10 caracteres)    |

### Visitor (Visitante) — Tabela `visitante`

| Campo      | Tipo          | Descrição                                |
|------------|---------------|------------------------------------------|
| `id`       | string (UUID) | Identificador único                      |
| `name`     | string        | Nome completo (coluna: `nome`)           |
| `cpf`      | string        | CPF com 11 dígitos (coluna: `cpf`)       |
| `phone`    | string        | Telefone (coluna: `telefone`)            |
| `photo`    | string        | URL da foto no S3 (coluna: `foto`)       |
| `liberado` | boolean       | Status de liberação (coluna: `liberado`) |

### Visit (Visita) — Tabela `visita`

| Campo          | Tipo              | Descrição                              |
|----------------|-------------------|----------------------------------------|
| `id`           | string (UUID)     | Identificador único                    |
| `dataEntrada`  | string (datetime) | Data/hora de entrada (coluna: `dataentrada`) |
| `dataSaida`    | string (datetime) or null | Data/hora de saída (coluna: `datasaida`) |
| `porteiroId`   | string (UUID)     | FK para o porteiro (coluna: `porteiro_id`) |
| `visitanteId`  | string (UUID)     | FK para o visitante (coluna: `visitante_id`) |
| `moradorId`    | string (UUID) or null | FK para o morador (coluna: `morador_id`) |

### Comunicado (Comunicado) — Tabela `comunicado`

| Campo            | Tipo              | Descrição                              |
|------------------|-------------------|----------------------------------------|
| `id`             | string (UUID)     | Identificador único                    |
| `titulo`         | string            | Título (coluna: `titulo`)              |
| `descricao`      | string            | Conteúdo (coluna: `descricao`)         |
| `dataPublicacao` | string (datetime) | Data de publicação (coluna: `datapublicacao`) |
| `sindicoId`      | string (UUID)     | FK para o síndico (coluna: `sindico_id`) |
| `sindicoNome`    | string            | Nome do síndico (relacionamento)       |

### Pagamento (Pagamento) — Tabela `pagamento`

| Campo            | Tipo              | Descrição                              |
|------------------|-------------------|----------------------------------------|
| `id`             | string (UUID)     | Identificador único                    |
| `valor`          | number            | Valor do pagamento (decimal(10,2))     |
| `vencimento`     | string (date)     | Data de vencimento                     |
| `dataPagamento`  | string (date) or null | Data de pagamento                  |
| `status`         | string            | `PENDENTE`, `PAGO` ou `ATRASADO`       |
| `moradorId`      | string (UUID)     | FK para o morador (coluna: `morador_id`) |

---

## Notas Técnicas

### Autenticação (Header-Based)

Enquanto o JWT não é implementado, as rotas protegidas de síndico utilizam autenticação via header:

- Header `X-User-ID`: UUID do usuário
- O middleware `RequireSindicoRole` valida se o usuário existe e possui role `SINDICO`
- O `sindico_id` é injetado no contexto da requisição

### Phone Validation

O telefone é validado usando a biblioteca [`libphonenumber`](https://github.com/nyaruka/phonenumbers) com região `BR`. O número é armazenado formatado no padrão nacional: `(11) 99999-9999`.

### Bcrypt

A senha é hasheada com bcrypt (custo padrão) antes de ser persistida.

### AWS Cognito

As operações no Cognito (`CreateUser`, `DeleteUser`) são executadas de forma **não fatal** — se falharem, o erro é logado mas a operação principal no banco de dados não é revertida.

### AWS S3

Fotos de visitantes são armazenadas no S3. O `UploadFile` retorna a URL pública completa. O `DeleteFile` extrai a key da URL e remove o objeto.

### Transação

A criação do usuário é atômica: usuário e apartamento são criados dentro de uma transação do banco de dados.
A liberação do visitante também é transacional: atualiza o campo `liberado` e cria o registro de `visita` na mesma transação.

### CORS

O servidor permite requisições de qualquer origem (`*`) com os métodos `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS` e headers `Content-Type` e `Authorization`.
