# API Reference — SmartCondo

> **Versão:** 2.0.0  
> **Base URL:** `http://localhost:8080`  
> **Formato:** JSON (exceto upload de visitante que usa `multipart/form-data`)  

---

## Sumário

- [Health Check](#health-check)
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
  - [Buscar Visitantes](#buscar-visitantes)
  - [Liberar Visitante](#liberar-visitante)
- [Comunicados (Síndico)](#comunicados-síndico)
  - [Publicar Comunicado](#publicar-comunicado)
  - [Listar Comunicados](#listar-comunicados)
  - [Obter Comunicado por ID](#obter-comunicado-por-id)
  - [Excluir Comunicado](#excluir-comunicado)
- [Autenticação](#autenticação)
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

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `full_name` | string | sim | Nome completo. Máx. 100 caracteres |
| `email` | string | sim | Email válido. Máx. 100 caracteres |
| `password` | string | sim | Senha. Mín. 8, máx. 72 caracteres |
| `phone` | string | sim | Telefone. Máx. 15 caracteres |
| `responsible` | boolean | não | Responsável pelo apartamento |
| `apartment` | object | sim | Dados do apartamento |
| `apartment.number` | number | sim | Número do apto. Deve ser > 0 |
| `apartment.block` | string | sim | Bloco. Máx. 10 caracteres |

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

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid request body"}` | JSON inválido ou campos ausentes |
| 409 | `{"error": "user already exists"}` | Email já cadastrado |
| 409 | `{"error": "apartment already registered"}` | Apartamento já cadastrado |
| 422 | `{"error": "apartment is required for residents"}` | Apartamento não informado |
| 500 | `{"error": "failed to create user"}` | Erro interno do servidor |

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
body := map[string]interface{}{
    "full_name": "Maria Silva",
    "email":     "maria@example.com",
    "password":  "password123",
    "phone":     "11999999999",
    "apartment": map[string]interface{}{
        "number": 101,
        "block":  "A",
    },
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

| Código | Body | Motivo |
|--------|------|--------|
| 500 | `{"error": "failed to list users"}` | Erro interno do servidor |

#### Exemplos

```bash
curl -X GET http://localhost:8080/users
```

---

### Obter Usuário por ID

### `GET /users/:id`

Retorna um único usuário pelo seu UUID.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do usuário |

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

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid user id"}` | UUID mal formatado |
| 404 | `{"error": "user not found"}` | Usuário não encontrado |
| 500 | `{"error": "failed to get user"}` | Erro interno do servidor |

#### Exemplos

```bash
curl -X GET http://localhost:8080/users/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

---

### Atualizar Usuário

### `PUT /users/:id`

Atualiza o nome e telefone de um usuário existente.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do usuário |

#### Request Body

```json
{
  "full_name": "Maria Santos",
  "phone": "11888888888"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `full_name` | string | sim | Novo nome completo. Máx. 100 chars |
| `phone` | string | sim | Novo telefone. Máx. 15 caracteres |

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

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid user id"}` | UUID mal formatado |
| 400 | `{"error": "invalid request body"}` | JSON inválido ou campos ausentes |
| 404 | `{"error": "user not found"}` | Usuário não encontrado |
| 422 | `{"error": "invalid user data"}` | Nome ou telefone vazio após trim |
| 500 | `{"error": "failed to update user"}` | Erro interno do servidor |

---

### Excluir Usuário

### `DELETE /users/:id`

Remove um usuário do sistema. A exclusão no AWS Cognito é tentada de forma **não crítica**.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do usuário |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid user id"}` | UUID mal formatado |
| 404 | `{"error": "user not found"}` | Usuário não encontrado |
| 409 | `{"error": "user has related records"}` | Usuário possui registros atrelados |
| 500 | `{"error": "failed to delete user"}` | Erro interno do servidor |

---

## Visitantes

### Criar Visitante

### `POST /visitors`

Cadastra um novo visitante com opção de upload de foto.

**Formato:** `multipart/form-data`

#### Validações

- CPF deve ser único no sistema
- CPF deve ter exatamente 11 dígitos
- Telefone é validado como número brasileiro
- Nome: máximo 100 caracteres

#### Request Body (multipart/form-data)

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `name` | string | sim | Nome completo. Máx. 100 caracteres |
| `cpf` | string | sim | CPF (11 dígitos, apenas números) |
| `phone` | string | sim | Telefone. Máx. 15 caracteres |
| `photo` | file | não | Foto do visitante (upload para S3) |

#### Respostas

**`201 Created`**
```json
{
  "id": "b1a2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Carlos Pereira",
  "cpf": "12345678901",
  "phone": "(11) 98888-8888",
  "photo": "https://bucket.s3.amazonaws.com/visitors/b1a2c3d4/photo.jpg",
  "liberado": false
}
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid request body"}` | Formulário inválido ou campos ausentes |
| 409 | `{"error": "visitor already exists"}` | CPF já cadastrado |
| 500 | `{"error": "failed to create visitor"}` | Erro interno do servidor |

#### Exemplos

```bash
curl -X POST http://localhost:8080/visitors \
  -F "name=Carlos Pereira" \
  -F "cpf=12345678901" \
  -F "phone=11988888888" \
  -F "photo=@foto.jpg"
```

```javascript
const form = new FormData();
form.append('name', 'Carlos Pereira');
form.append('cpf', '12345678901');
form.append('phone', '11988888888');
form.append('photo', fileInput.files[0]);

const response = await fetch('http://localhost:8080/visitors', {
  method: 'POST',
  body: form
});
const visitor = await response.json();
```

```go
var buf bytes.Buffer
w := multipart.NewWriter(&buf)
w.WriteField("name", "Carlos Pereira")
w.WriteField("cpf", "12345678901")
w.WriteField("phone", "11988888888")
part, _ := w.CreateFormFile("photo", "foto.jpg")
io.Copy(part, fileReader)
w.Close()

resp, _ := http.Post("http://localhost:8080/visitors", w.FormDataContentType(), &buf)
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
    "id": "b1a2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Carlos Pereira",
    "cpf": "12345678901",
    "phone": "(11) 98888-8888",
    "photo": "https://bucket.s3.amazonaws.com/visitors/b1a2c3d4/photo.jpg",
    "liberado": false
  }
]
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 500 | `{"error": "failed to list visitors"}` | Erro interno do servidor |

---

### Obter Visitante por ID

### `GET /visitors/:id`

Retorna um visitante específico pelo seu UUID.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do visitante |

#### Respostas

**`200 OK`**
```json
{
  "id": "b1a2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Carlos Pereira",
  "cpf": "12345678901",
  "phone": "(11) 98888-8888",
  "photo": "https://bucket.s3.amazonaws.com/visitors/b1a2c3d4/photo.jpg",
  "liberado": false
}
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid visitor id"}` | UUID mal formatado |
| 404 | `{"error": "visitor not found"}` | Visitante não encontrado |
| 500 | `{"error": "failed to get visitor"}` | Erro interno do servidor |

---

### Excluir Visitante

### `DELETE /visitors/:id`

Remove um visitante do sistema. Se possuir foto, o arquivo também é removido do S3.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do visitante |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid visitor id"}` | UUID mal formatado |
| 404 | `{"error": "visitor not found"}` | Visitante não encontrado |
| 500 | `{"error": "failed to delete visitor"}` | Erro interno do servidor |

---

## Porteiro

### Buscar Visitantes

### `GET /porteiros/visitantes`

Busca visitantes com filtros. **Pelo menos um filtro é obrigatório.**

Os filtros são combinados com AND. Campos de texto usam busca parcial (ILIKE), exceto CPF que é busca exata.

#### Parâmetros de Query

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `nome` | string | não | Busca parcial no nome (ILIKE) |
| `cpf` | string | não | Busca exata no CPF |
| `telefone` | string | não | Busca parcial no telefone (ILIKE) |
| `liberado` | boolean | não | Filtra por status de liberação |

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "b1a2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Carlos Pereira",
    "cpf": "12345678901",
    "phone": "(11) 98888-8888",
    "photo": "https://bucket.s3.amazonaws.com/visitors/b1a2c3d4/photo.jpg",
    "liberado": false
  }
]
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "at least one search filter is required"}` | Nenhum filtro informado |
| 400 | `{"error": "invalid query parameters"}` | Parâmetros inválidos |
| 500 | `{"error": "failed to search visitors"}` | Erro interno do servidor |

#### Exemplos

```bash
# Buscar por nome
curl -X GET "http://localhost:8080/porteiros/visitantes?nome=Carlos"

# Buscar por CPF (exato)
curl -X GET "http://localhost:8080/porteiros/visitantes?cpf=12345678901"

# Buscar visitantes não liberados
curl -X GET "http://localhost:8080/porteiros/visitantes?liberado=false"

# Múltiplos filtros
curl -X GET "http://localhost:8080/porteiros/visitantes?nome=Carlos&liberado=false"
```

---

### Liberar Visitante

### `PATCH /porteiros/visitantes/{id}/liberar`

Registra a liberação de entrada de um visitante. A operação é **atômica** (transação):

1. Marca o visitante como `liberado: true`
2. Cria um registro de visita com data/hora atual

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do visitante a liberar |

#### Request Body

```json
{
  "porteiro_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "morador_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `porteiro_id` | string (UUID) | sim | UUID do porteiro responsável |
| `morador_id` | string (UUID) | não | UUID do morador de destino |

#### Respostas

**`200 OK`**
```json
{
  "id": "b1a2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Carlos Pereira",
  "cpf": "12345678901",
  "phone": "(11) 98888-8888",
  "photo": "https://bucket.s3.amazonaws.com/visitors/b1a2c3d4/photo.jpg",
  "liberado": true
}
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid visitor id"}` | UUID do visitante inválido |
| 400 | `{"error": "invalid request body"}` | JSON inválido ou porteiro_id ausente |
| 400 | `{"error": "porteiro not found"}` | Porteiro não encontrado |
| 404 | `{"error": "visitor not found"}` | Visitante não encontrado |
| 500 | `{"error": "failed to release visitor"}` | Erro interno do servidor |

#### Exemplos

```bash
curl -X PATCH http://localhost:8080/porteiros/visitantes/b1a2c3d4-e5f6-7890-abcd-ef1234567890/liberar \
  -H "Content-Type: application/json" \
  -d '{
    "porteiro_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }'
```

---

## Comunicados (Síndico)

> **Atenção:** Os endpoints de criação e exclusão de comunicados **requerem autenticação**.
> Envie o header `X-User-ID` com o UUID de um usuário com role `SINDICO`.

### Publicar Comunicado

### `POST /sindico/comunicados`

Publica um novo comunicado no condomínio.

#### Headers

| Header | Obrigatório | Descrição |
|--------|-------------|-----------|
| `X-User-ID` | sim | UUID do síndico autenticado |

#### Request Body

```json
{
  "titulo": "Manutenção do elevador",
  "descricao": "O elevador estará em manutenção no dia 15/08 das 8h às 18h."
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `titulo` | string | sim | Título do comunicado. Máx. 100 caracteres |
| `descricao` | string | sim | Descrição detalhada |

#### Respostas

**`201 Created`**
```json
{
  "id": "c1d2e3f4-a5b6-7890-abcd-ef1234567890",
  "titulo": "Manutenção do elevador",
  "descricao": "O elevador estará em manutenção no dia 15/08 das 8h às 18h.",
  "dataPublicacao": "2026-07-23T14:30:00Z",
  "sindicoId": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "sindicoNome": "José Síndico"
}
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid request body"}` | JSON inválido |
| 401 | `{"error": "missing authentication header"}` | Header X-User-ID ausente |
| 401 | `{"error": "invalid user id"}` | UUID inválido |
| 403 | `{"error": "user is not authorized as sindico"}` | Usuário não é síndico |
| 500 | `{"error": "failed to publish comunicado"}` | Erro interno do servidor |

#### Exemplos

```bash
curl -X POST http://localhost:8080/sindico/comunicados \
  -H "Content-Type: application/json" \
  -H "X-User-ID: f47ac10b-58cc-4372-a567-0e02b2c3d479" \
  -d '{
    "titulo": "Manutenção do elevador",
    "descricao": "O elevador estará em manutenção no dia 15/08 das 8h às 18h."
  }'
```

```javascript
const response = await fetch('http://localhost:8080/sindico/comunicados', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-User-ID': 'f47ac10b-58cc-4372-a567-0e02b2c3d479'
  },
  body: JSON.stringify({
    titulo: 'Manutenção do elevador',
    descricao: 'O elevador estará em manutenção no dia 15/08 das 8h às 18h.'
  })
});
const comunicado = await response.json();
```

---

### Listar Comunicados

### `GET /sindico/comunicados`

Retorna todos os comunicados publicados, ordenados por data de publicação (decrescente).
Inclui o nome do síndico que publicou.

#### Requisição

Nenhum header de autenticação necessário.

#### Respostas

**`200 OK`**
```json
[
  {
    "id": "c1d2e3f4-a5b6-7890-abcd-ef1234567890",
    "titulo": "Manutenção do elevador",
    "descricao": "O elevador estará em manutenção no dia 15/08 das 8h às 18h.",
    "dataPublicacao": "2026-07-23T14:30:00Z",
    "sindicoId": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "sindicoNome": "José Síndico"
  }
]
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 500 | `{"error": "failed to list comunicados"}` | Erro interno do servidor |

---

### Obter Comunicado por ID

### `GET /sindico/comunicados/:id`

Retorna um comunicado específico.

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do comunicado |

#### Respostas

**`200 OK`**
```json
{
  "id": "c1d2e3f4-a5b6-7890-abcd-ef1234567890",
  "titulo": "Manutenção do elevador",
  "descricao": "O elevador estará em manutenção no dia 15/08 das 8h às 18h.",
  "dataPublicacao": "2026-07-23T14:30:00Z",
  "sindicoId": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "sindicoNome": "José Síndico"
}
```

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid comunicado id"}` | UUID mal formatado |
| 404 | `{"error": "comunicado not found"}` | Comunicado não encontrado |
| 500 | `{"error": "failed to get comunicado"}` | Erro interno do servidor |

---

### Excluir Comunicado

### `DELETE /sindico/comunicados/{id}`

Remove um comunicado. **Apenas o síndico autor que o publicou pode excluí-lo.**

#### Headers

| Header | Obrigatório | Descrição |
|--------|-------------|-----------|
| `X-User-ID` | sim | UUID do síndico autenticado |

#### Parâmetros de Path

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `id` | string | sim | UUID do comunicado |

#### Respostas

**`204 No Content`** — Corpo vazio.

#### Erros

| Código | Body | Motivo |
|--------|------|--------|
| 400 | `{"error": "invalid comunicado id"}` | UUID mal formatado |
| 401 | `{"error": "missing authentication header"}` | Header X-User-ID ausente |
| 403 | `{"error": "user is not authorized as sindico"}` | Usuário não é síndico |
| 403 | `{"error": "you can only delete your own comunicados"}` | Não é o autor |
| 404 | `{"error": "comunicado not found"}` | Comunicado não encontrado |
| 500 | `{"error": "failed to delete comunicado"}` | Erro interno do servidor |

---

## Autenticação

Atualmente, apenas o módulo de **Comunicados** exige autenticação.

### Método: Header `X-User-ID`

```bash
X-User-ID: f47ac10b-58cc-4372-a567-0e02b2c3d479
```

**Fluxo do middleware `RequireSindicoRole`:**

1. Lê o header `X-User-ID`
2. Busca o usuário pelo UUID no banco de dados
3. Verifica se o `role` do usuário é `SINDICO`
4. Se aprovado, injeta `user_id` no contexto da requisição

### Endpoints Protegidos

| Método | Path | Middleware |
|--------|------|-----------|
| `POST` | `/sindico/comunicados` | `RequireSindicoRole` |
| `DELETE` | `/sindico/comunicados/:id` | `RequireSindicoRole` |

### Respostas de Autenticação

| Código | Body | Motivo |
|--------|------|--------|
| 401 | `{"error": "missing authentication header"}` | Header `X-User-ID` não informado |
| 401 | `{"error": "user is not authorized as sindico"}` | Usuário não encontrado ou não é síndico |
| 403 | `{"error": "user is not authorized as sindico"}` | Usuário existe mas role não é `SINDICO` |

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

| Código | Significado |
|--------|-------------|
| 400 | Bad Request — requisição mal formatada |
| 401 | Unauthorized — autenticação ausente ou inválida |
| 403 | Forbidden — sem permissão para a operação |
| 404 | Not Found — recurso não encontrado |
| 409 | Conflict — duplicidade ou registro em uso |
| 422 | Unprocessable Entity — dados inválidos |
| 500 | Internal Server Error — erro interno |

### Todos os Erros da Aplicação

| Mensagem | Onde Ocorre | Causa |
|----------|-------------|-------|
| `invalid request body` | Usuários, Visitantes, Comunicados | JSON ou formulário inválido |
| `invalid user id` | Usuários, Autenticação | UUID mal formatado |
| `invalid visitor id` | Visitantes | UUID mal formatado |
| `invalid comunicado id` | Comunicados | UUID mal formatado |
| `invalid query parameters` | Porteiro | Parâmetros de busca inválidos |
| `user already exists` | Usuários | Email já cadastrado |
| `user not found` | Usuários | Usuário não encontrado |
| `user has related records` | Usuários | Possui registros associados |
| `apartment is required for residents` | Usuários | Apartamento não informado |
| `apartment already registered` | Usuários | Número+bloco já existe |
| `invalid user data` | Usuários | Nome/telefone vazio após trim |
| `visitor already exists` | Visitantes | CPF já cadastrado |
| `visitor not found` | Visitantes, Porteiro | Visitante não encontrado |
| `porteiro not found` | Porteiro | Porteiro informado não existe |
| `at least one search filter is required` | Porteiro | Nenhum filtro informado |
| `comunicado not found` | Comunicados | Comunicado não encontrado |
| `missing authentication header` | Comunicados (protegidos) | Header X-User-ID não enviado |
| `user is not authorized as sindico` | Comunicados (protegidos) | Usuário não tem role SINDICO |
| `you can only delete your own comunicados` | Comunicados | Síndico tentou excluir comunicado de outro |
| `failed to create/update/delete/get/list ...` | Vários | Erro interno do servidor |

---

## Modelos de Dados

### User (Usuário)

**Tabela:** `usuario`

| Campo | Tipo | Coluna | Descrição |
|-------|------|--------|-----------|
| `id` | string (UUID) | `id` | Identificador único |
| `full_name` | string | `nome` | Nome completo. Máx. 100 |
| `email` | string | `email` | Email único. Máx. 100 |
| `password` | string | `senha` | Hash bcrypt da senha. Máx. 100 |
| `phone` | string | `telefone` | Telefone formatado. Máx. 15 |
| `status` | string | `status` | `ATIVO`, `INATIVO` ou `BLOQUEADO` |
| `role` | string | `tipo` | `MORADOR`, `PORTEIRO` ou `SINDICO` |
| `apartment_id` | UUID ou null | `apartamento_id` | FK para apartamento |
| `responsible` | boolean | `responsavel` | Responsável pelo apto |

**Regras de consistência (banco de dados):**
- `responsavel = TRUE` só é permitido quando `role = 'MORADOR'`
- `apartamento_id` só pode ser não-nulo quando `role = 'MORADOR'`
- Somente `MORADOR` pode ter apartamento vinculado

### Apartment (Apartamento)

**Tabela:** `apartamento`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | string (UUID) | Identificador único |
| `number` | integer | Número do apartamento |
| `block` | string | Bloco. Máx. 10 caracteres |

### Visitor (Visitante)

**Tabela:** `visitante`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | string (UUID) | Identificador único |
| `name` | string | Nome completo. Máx. 100 |
| `cpf` | string | CPF único. 11 dígitos |
| `phone` | string | Telefone. Máx. 15 |
| `photo` | string | URL da foto no S3. Máx. 255 |
| `liberado` | boolean | Se foi liberado para entrada |

### Visit (Visita)

**Tabela:** `visita`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | string (UUID) | Identificador único |
| `dataEntrada` | timestamp | Data/hora da entrada |
| `dataSaida` | timestamp ou null | Data/hora da saída |
| `porteiro_id` | UUID (FK) | Porteiro que liberou |
| `visitante_id` | UUID (FK) | Visitante liberado |
| `morador_id` | UUID ou null (FK) | Morador de destino |

### Comunicado

**Tabela:** `comunicado`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | string (UUID) | Identificador único |
| `titulo` | string | Título. Máx. 100 |
| `descricao` | text | Descrição detalhada |
| `dataPublicacao` | timestamp | Data/hora da publicação |
| `sindico_id` | UUID (FK) | Síndico que publicou |

---

## Notas Técnicas

### Phone Validation

O telefone é validado usando a biblioteca [`libphonenumber`](https://github.com/nyaruka/phonenumbers) com região `BR`. O número é armazenado formatado no padrão nacional: `(11) 99999-9999`.

### Bcrypt

A senha do usuário é hasheada com bcrypt (custo padrão) antes de ser persistida.

### AWS Cognito

As operações no Cognito (`CreateUser`, `DeleteUser`) são executadas de forma **não fatal** — se falharem, o erro é logado mas a operação principal no banco não é revertida.

### AWS S3

Fotos de visitantes são armazenadas no S3 no path `visitors/{uuid}/photo.{ext}`.

### Transações

- Criação de **usuário**: transação atômica (usuário + apartamento)
- **Liberação** de visitante: transação atômica (update liberado + create visita)

### CORS

O servidor permite requisições de qualquer origem (`*`) com os métodos `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS` e headers `Content-Type` e `Authorization`.

### Documentação Interativa

Acesse `http://localhost:8080/` ou `http://localhost:8080/docs` para o Swagger UI.
