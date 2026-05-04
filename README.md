# Workshop API

Sistema Integrado de Atendimento e Execucao de Servicos para oficina mecanica.

## Tecnologias

- **Go 1.25** — linguagem principal
- **Gin** — framework HTTP
- **PostgreSQL 16** — banco de dados relacional (via GORM)
- **Swagger** — documentacao da API (`swaggo/swag`)
- **JWT** — autenticacao (`golang-jwt/jwt`)
- **Docker + docker-compose** — containerizacao

### Justificativa do banco de dados

PostgreSQL foi escolhido por:
- **Integridade relacional**: solicitantes, veiculos e ordens de servico possuem relacoes fortes (FK) que o Postgres garante nativamente
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
  middleware/          — auth JWT, CORS

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

| Servico       | URL                                          | Observacao                              |
|---------------|----------------------------------------------|-----------------------------------------|
| API           | `http://localhost:8080`                      |                                         |
| Swagger UI    | `http://localhost:8080/swagger/index.html`   |                                         |
| pgAdmin       | `http://localhost:8082`                      | login: `admin@workshop.com` / `admin`   |
| Documentacao  | `http://localhost:8083`                      | landing page do projeto                 |

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

| Variavel                 | Padrao                                                                 | Descricao                           |
|--------------------------|------------------------------------------------------------------------|-------------------------------------|
| `APP_PORT`               | `8080`                                                                 | Porta da API                        |
| `APP_ENV`                | `development`                                                          | Ambiente (development/prod)         |
| `POSTGRES_DSN`           | `postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable` | DSN do PostgreSQL                   |
| `JWT_SECRET`             | `change-me-in-production`                                              | Chave secreta JWT                   |
| `JWT_EXPIRATION_HOURS`   | `24`                                                                   | Expiracao do token JWT (horas)      |
| `ACCESS_TOKEN_EXP_MIN`   | `15`                                                                   | Expiracao do access token (minutos) |
| `REFRESH_TOKEN_EXP_DAYS` | `7`                                                                    | Expiracao do refresh token (dias)   |
| `BCRYPT_COST`            | `12`                                                                   | Custo do hash de senha com bcrypt   |
| `MAX_FAILED_LOGINS`      | `5`                                                                    | Tentativas antes de bloquear login  |
| `LOGIN_LOCK_MIN`         | `15`                                                                   | Tempo de bloqueio apos falhas (min) |

## Endpoints

As rotas `/api/v1/auth/register`, `/login` e `/refresh` sao publicas. Todas as demais rotas `/api/v1/*` requerem header `Authorization: Bearer <token>`. Rotas marcadas como **admin** exigem `role=admin` nas claims do JWT.

| Metodo | Rota                                          | Auth     | Descricao                              |
|--------|-----------------------------------------------|----------|----------------------------------------|
| GET    | `/health`                                     | —        | Health check                           |
| **Auth** | | | |
| POST   | `/api/v1/auth/register`                       | publico  | Registrar usuario (role=client)        |
| POST   | `/api/v1/auth/login`                          | publico  | Login (retorna access + refresh token) |
| POST   | `/api/v1/auth/refresh`                        | publico  | Rotacionar tokens                      |
| POST   | `/api/v1/auth/logout`                         | JWT      | Encerrar sessao                        |
| **Requesters** | | | |
| POST   | `/api/v1/requesters`                          | admin    | Criar solicitante                      |
| GET    | `/api/v1/requesters`                          | admin    | Listar solicitantes                    |
| GET    | `/api/v1/requesters/:id`                      | admin    | Buscar solicitante por ID              |
| GET    | `/api/v1/requesters/document/:document`       | admin    | Buscar solicitante por CPF/CNPJ        |
| PUT    | `/api/v1/requesters/:id`                      | admin    | Atualizar solicitante                  |
| DELETE | `/api/v1/requesters/:id`                      | admin    | Deletar solicitante                    |
| **Vehicles** | | | |
| POST   | `/api/v1/vehicles`                            | admin    | Cadastrar veiculo                      |
| GET    | `/api/v1/vehicles`                            | admin    | Listar veiculos                        |
| GET    | `/api/v1/vehicles/:id`                        | admin    | Buscar veiculo por ID                  |
| GET    | `/api/v1/requesters/:id/vehicles`             | admin    | Listar veiculos por solicitante        |
| PUT    | `/api/v1/vehicles/:id`                        | admin    | Atualizar veiculo                      |
| DELETE | `/api/v1/vehicles/:id`                        | admin    | Deletar veiculo                        |
| **Services** | | | |
| POST   | `/api/v1/services`                            | admin    | Cadastrar servico                      |
| GET    | `/api/v1/services`                            | admin    | Listar servicos                        |
| GET    | `/api/v1/services/:id`                        | admin    | Buscar servico por ID                  |
| PUT    | `/api/v1/services/:id`                        | admin    | Atualizar servico                      |
| DELETE | `/api/v1/services/:id`                        | admin    | Deletar servico                        |
| **Parts** | | | |
| POST   | `/api/v1/parts`                               | admin    | Cadastrar peca/insumo                  |
| GET    | `/api/v1/parts`                               | admin    | Listar pecas/insumos                   |
| GET    | `/api/v1/parts/:id`                           | admin    | Buscar peca por ID                     |
| PUT    | `/api/v1/parts/:id`                           | admin    | Atualizar peca                         |
| DELETE | `/api/v1/parts/:id`                           | admin    | Deletar peca                           |
| PATCH  | `/api/v1/parts/:id/stock`                     | admin    | Ajustar estoque (delta +/-)            |
| **Service Orders** | | | |
| POST   | `/api/v1/service-orders`                      | admin    | Criar ordem de servico                 |
| GET    | `/api/v1/service-orders`                      | admin    | Listar ordens de servico               |
| GET    | `/api/v1/service-orders/:id`                  | admin    | Buscar ordem por ID                    |
| GET    | `/api/v1/service-orders/metrics/execution-time` | admin  | Tempo medio de execucao das OSs        |
| GET    | `/api/v1/service-orders/code/:code`           | publico  | Buscar OS por codigo (portal cliente)  |
| GET    | `/api/v1/service-orders/requester`            | publico  | Listar OSs por documento do solicitante|
| PUT    | `/api/v1/service-orders/code/:code/status`    | publico  | Atualizar status por codigo            |
| PUT    | `/api/v1/service-orders/:id/status`           | admin    | Atualizar status da OS                 |
| PUT    | `/api/v1/service-orders/:id`                  | admin    | Atualizar OS                           |
| DELETE | `/api/v1/service-orders/:id`                  | admin    | Deletar OS                             |

## Gerando o Swagger

```bash
swag init -g cmd/api/main.go -o docs
```

## Testes

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```
