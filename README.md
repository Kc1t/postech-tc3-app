# Workshop API

Sistema Integrado de Atendimento e Execucao de Servicos para oficina mecanica.

## Tecnologias

- **Go 1.22** — linguagem principal
- **Gin** — framework HTTP
- **MongoDB** — banco de dados (driver oficial `mongo-driver`)
- **Swagger** — documentacao da API (`swaggo/swag`)
- **JWT** — autenticacao (`golang-jwt/jwt`)
- **Docker + docker-compose** — containerizacao

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
    outbound/mongodb/  — implementacoes dos repositories
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
O Mongo Express estara em `http://localhost:8081`.

### Rodar localmente (sem Docker)

```bash
# Necessario ter Go 1.22+ e MongoDB rodando localmente
cp .env.example .env

# Gerar docs do swagger
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

go run ./cmd/api
```

## Variaveis de ambiente

| Variavel               | Padrao                      | Descricao                  |
|------------------------|-----------------------------|----------------------------|
| `APP_PORT`             | `8080`                      | Porta da API               |
| `APP_ENV`              | `development`               | Ambiente                   |
| `MONGO_URI`            | `mongodb://localhost:27017` | URI do MongoDB             |
| `MONGO_DB`             | `workshop`                  | Nome do banco              |
| `JWT_SECRET`           | `change-me-in-production`   | Chave secreta JWT          |
| `JWT_EXPIRATION_HOURS` | `24`                        | Expiracao do token (horas) |

## Endpoints

| Metodo | Rota                            | Descricao                        |
|--------|---------------------------------|----------------------------------|
| GET    | `/health`                       | Health check                     |
| POST   | `/api/v1/customers`             | Criar cliente                    |
| GET    | `/api/v1/customers`             | Listar clientes                  |
| GET    | `/api/v1/customers/:id`         | Buscar cliente                   |
| PUT    | `/api/v1/customers/:id`         | Atualizar cliente                |
| DELETE | `/api/v1/customers/:id`         | Deletar cliente                  |
| POST   | `/api/v1/vehicles`              | Cadastrar veiculo                |
| GET    | `/api/v1/vehicles`              | Listar veiculos                  |
| POST   | `/api/v1/service-orders`        | Criar ordem de servico           |
| GET    | `/api/v1/service-orders`        | Listar ordens de servico         |
| GET    | `/api/v1/service-orders/:id`    | Buscar ordem de servico          |
| PUT    | `/api/v1/service-orders/:id/status` | Atualizar status da OS       |
| POST   | `/api/v1/services`              | Cadastrar servico                |
| POST   | `/api/v1/parts`                 | Cadastrar peca/insumo            |

Todas as rotas `/api/v1/*` requerem header `Authorization: Bearer <token>`.

## Gerando o Swagger

```bash
swag init -g cmd/api/main.go -o docs
```
