# Workshop API

Sistema Integrado de Atendimento e Execucao de Servicos para oficina mecanica.

## Tecnologias

- **Go 1.23** — linguagem principal
- **Gin** — framework HTTP
- **PostgreSQL 16** — banco de dados relacional (via GORM)
- **Swagger** — documentacao da API (`swaggo/swag`)
- **JWT** — autenticacao (`golang-jwt/jwt`)
- **Docker + docker-compose** — containerizacao

### Justificativa do banco de dados

PostgreSQL foi escolhido por:
- **Integridade relacional**: clientes, veiculos e ordens de servico possuem relacoes fortes (FK) que o Postgres garante nativamente
- **UUID nativo**: `gen_random_uuid()` para IDs sem dependencia externa
- **JSONB**: permite armazenar value objects (servicos e pecas dentro da OS) como JSON sem perder a capacidade de consulta
- **Maturidade e ecossistema**: driver oficial Go (`pgx`), ORM maduro (GORM), ferramentas de administracao (pgAdmin)
- **ACID**: transacoes completas para operacoes criticas como ajuste de estoque

## Arquitetura

Arquitetura Hexagonal (Ports & Adapters) dentro de um monolito:

```
cmd/api/
  main.go              — entrypoint
  bootstrap/           — DI container (singleton)
  routes/              — setup de rotas
  middleware/           — auth JWT, CORS

internal/
  domain/              — entidades de negocio
  ports/               — interfaces (repositories + usecases)
  adapters/
    inbound/http/      — handlers HTTP (Gin)
    outbound/postgresql/ — implementacoes dos repositories (GORM)
  application/usecase/ — logica de negocio
```

## Como rodar

### Pre-requisitos

- Docker e docker-compose instalados

### Subir o ambiente completo

```bash
cp .env.example .env
docker compose up --build
```

A API estara disponivel em `http://localhost:8080`.
O Swagger UI estara em `http://localhost:8080/swagger/index.html`.
O pgAdmin estara em `http://localhost:8082` (login: `admin@workshop.com` / `admin`).

### Rodar localmente (sem Docker)

```bash
# Necessario ter Go 1.23+ e PostgreSQL rodando localmente
cp .env.example .env

# Gerar docs do swagger
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

go run ./cmd/api
```

## Variaveis de ambiente

| Variavel                 | Padrao                                                                 | Descricao                              |
|--------------------------|------------------------------------------------------------------------|----------------------------------------|
| `APP_PORT`               | `8080`                                                                 | Porta da API                           |
| `APP_ENV`                | `development`                                                          | Ambiente (development/prod)            |
| `POSTGRES_DSN`           | `postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable` | DSN do PostgreSQL                      |
| `JWT_SECRET`             | `secret`                                                               | Chave secreta JWT                      |
| `ACCESS_TOKEN_EXP_MIN`   | `15`                                                                   | Expiracao do access token (minutos)    |
| `REFRESH_TOKEN_EXP_DAYS` | `7`                                                                    | Expiracao do refresh token (dias)      |
| `BCRYPT_COST`            | `12`                                                                   | Custo do hash de senha com bcrypt      |
| `MAX_FAILED_LOGINS`      | `5`                                                                    | Tentativas falhas antes de bloquear a conta |
| `LOGIN_LOCK_MIN`         | `15`                                                                   | Duracao do bloqueio apos limite (minutos) |
| `ADMIN_EMAIL`            | `admin@workshop.com`                                                   | Email do admin seeded no boot          |
| `ADMIN_PASSWORD`         | _(obrigatorio em prod)_                                                | Senha do admin seeded no boot          |

## Endpoints

| Metodo | Rota                                    | Descricao                        |
|--------|-----------------------------------------|----------------------------------|
| GET    | `/health`                               | Health check                     |
| **Auth (publico)** | | |
| POST   | `/api/v1/auth/register`                 | Registrar usuario (role=client)  |
| POST   | `/api/v1/auth/login`                    | Login (retorna access + refresh) |
| POST   | `/api/v1/auth/refresh`                  | Rotacionar tokens                |
| POST   | `/api/v1/auth/logout`                   | Encerrar sessao (requer JWT)     |
| **Customers** | | |
| POST   | `/api/v1/customers`                     | Criar cliente                    |
| GET    | `/api/v1/customers`                     | Listar clientes                  |
| GET    | `/api/v1/customers/:id`                 | Buscar cliente por ID            |
| GET    | `/api/v1/customers/document/:document`  | Buscar cliente por CPF/CNPJ      |
| PUT    | `/api/v1/customers/:id`                 | Atualizar cliente                |
| DELETE | `/api/v1/customers/:id`                 | Deletar cliente                  |
| **Vehicles** | | |
| POST   | `/api/v1/vehicles`                      | Cadastrar veiculo                |
| GET    | `/api/v1/vehicles`                      | Listar veiculos                  |
| GET    | `/api/v1/vehicles/:id`                  | Buscar veiculo por ID            |
| GET    | `/api/v1/customers/:id/vehicles`        | Listar veiculos por cliente      |
| PUT    | `/api/v1/vehicles/:id`                  | Atualizar veiculo                |
| DELETE | `/api/v1/vehicles/:id`                  | Deletar veiculo                  |
| **Services** | | |
| POST   | `/api/v1/services`                      | Cadastrar servico                |
| GET    | `/api/v1/services`                      | Listar servicos                  |
| GET    | `/api/v1/services/:id`                  | Buscar servico por ID            |
| PUT    | `/api/v1/services/:id`                  | Atualizar servico                |
| DELETE | `/api/v1/services/:id`                  | Deletar servico                  |
| **Parts** | | |
| POST   | `/api/v1/parts`                         | Cadastrar peca/insumo            |
| GET    | `/api/v1/parts`                         | Listar pecas/insumos             |
| GET    | `/api/v1/parts/:id`                     | Buscar peca por ID               |
| PUT    | `/api/v1/parts/:id`                     | Atualizar peca                   |
| DELETE | `/api/v1/parts/:id`                     | Deletar peca                     |
| PATCH  | `/api/v1/parts/:id/stock`               | Ajustar estoque                  |
| **Service Orders** | | |
| POST   | `/api/v1/service-orders`                | Criar ordem de servico           |
| GET    | `/api/v1/service-orders`                | Listar ordens de servico         |
| GET    | `/api/v1/service-orders/:id`            | Buscar ordem por ID              |
| GET    | `/api/v1/customers/:id/service-orders`  | Listar ordens por cliente        |
| PUT    | `/api/v1/service-orders/:id/status`     | Atualizar status da OS           |
| PUT    | `/api/v1/service-orders/:id`            | Atualizar OS                     |
| DELETE | `/api/v1/service-orders/:id`            | Deletar OS                       |

As rotas `/api/v1/auth/register`, `/login` e `/refresh` sao publicas (com rate limit). Todas as demais rotas `/api/v1/*` requerem header `Authorization: Bearer <token>`. Rotas administrativas exigem `role=admin` nas claims do JWT.

## Gerando o Swagger

```bash
swag init -g cmd/api/main.go -o docs
```

## Testes

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```
