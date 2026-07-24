# API Reference — SmartCondo

> **Versão:** 1.0.0  
> **Base URL:** `http://localhost:8080`  
> **Formato:** JSON  
> **Autenticação:** Nenhuma (rotas públicas)

---

## Sumário

- [Health Check](#health-check)
- [Usuários](#usuários)
  - [Criar Usuário](#criar-usuário)
  - [Listar Usuários](#listar-usuários)
  - [Obter Usuário por ID](#obter-usuário-por-id)
  - [Atualizar Usuário](#atualizar-usuário)
  - [Excluir Usuário](#excluir-usuário)
- [Erros](#erros)
- [Modelos de Dados](#modelos-de-dados)

---

## Health Check

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
| 404    | Not Found — recurso não encontrado |
| 409    | Conflict — conflito (duplicidade, registro em uso) |
| 422    | Unprocessable Entity — dados inválidos |
| 500    | Internal Server Error — erro interno |

### Erros da Aplicação

| Erro (mensagem)                          | Causa                                          |
|------------------------------------------|------------------------------------------------|
| `invalid request body`                   | JSON mal formatado ou campos obrigatórios ausentes |
| `invalid user id`                        | UUID do usuário inválido                       |
| `user already exists`                    | Tentativa de cadastro com email já existente   |
| `user not found`                         | Usuário não encontrado pelo ID informado       |
| `user has related records`               | Tentativa de excluir usuário com registros associados |
| `apartment is required for residents`    | Morador criado sem informar apartamento        |
| `apartment already registered`           | Número e bloco do apartamento já cadastrados   |
| `invalid user data`                      | Nome ou telefone enviado vazio                 |
| `failed to create user`                  | Erro interno ao criar usuário                  |
| `failed to get user`                     | Erro interno ao buscar usuário                 |
| `failed to list users`                   | Erro interno ao listar usuários                |
| `failed to update user`                  | Erro interno ao atualizar usuário              |
| `failed to delete user`                  | Erro interno ao excluir usuário                |

---

## Modelos de Dados

### User (Usuário)

| Campo        | Tipo              | Descrição                                    |
|--------------|-------------------|----------------------------------------------|
| `id`         | string (UUID)     | Identificador único                          |
| `full_name`  | string            | Nome completo (coluna: `nome`)               |
| `email`      | string            | Email único (coluna: `email`)                |
| `password`   | string            | Hash bcrypt da senha (coluna: `senha`)       |
| `phone`      | string            | Telefone formatado (coluna: `telefone`)      |
| `status`     | string            | `ATIVO`, `INATIVO` ou `BLOQUEADO`            |
| `role`       | string            | `MORADOR`, `PORTEIRO` ou `SINDICO`           |
| `apartment_id` | string (UUID) ou null | FK para o apartamento (coluna: `apartamento_id`) |
| `apartment`  | object or null    | Dados do apartamento                         |
| `responsible`| boolean           | Se é o responsável pelo apto (coluna: `responsavel`) |

**Regras de consistência (banco de dados):**
- `responsavel = TRUE` só é permitido quando `role = 'MORADOR'`
- `apartamento_id` só pode ser não-nulo quando `role = 'MORADOR'`

### Apartment (Apartamento)

| Campo    | Tipo          | Descrição                     |
|----------|---------------|-------------------------------|
| `id`     | string (UUID) | Identificador único           |
| `number` | integer       | Número do apartamento         |
| `block`  | string        | Bloco (máx. 10 caracteres)    |

---

## Notas Técnicas

### Phone Validation

O telefone é validado usando a biblioteca [`libphonenumber`](https://github.com/nyaruka/phonenumbers) com região `BR`. O número é armazenado formatado no padrão nacional: `(11) 99999-9999`.

### Bcrypt

A senha é hasheada com bcrypt (custo padrão) antes de ser persistida.

### AWS Cognito

As operações no Cognito (`CreateUser`, `DeleteUser`) são executadas de forma **não fatal** — se falharem, o erro é logado mas a operação principal no banco de dados não é revertida.

### Transação

A criação do usuário é atômica: usuário e apartamento são criados dentro de uma transação do banco de dados.

### CORS

O servidor permite requisições de qualquer origem (`*`) com os métodos `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS` e headers `Content-Type` e `Authorization`.
