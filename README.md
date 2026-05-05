<div align="center">

# Workshop API

Sistema integrado de atendimento e execução de serviços para oficinas mecânicas.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-HTTP%20Framework-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

</div>

## Sumário

- [Sobre o projeto](#sobre-o-projeto)
- [Principais recursos](#principais-recursos)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Banco de dados](#banco-de-dados)
- [Como rodar](#como-rodar)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Endpoints](#endpoints)
- [Testes](#testes)
- [Contribuidores](#contribuidores)

## Sobre o projeto

O **Workshop API** é uma API REST para gestão do fluxo operacional de uma oficina mecânica. A aplicação centraliza o cadastro de clientes, veículos, serviços, peças e ordens de serviço, além de oferecer autenticação JWT, controle de perfis administrativos e consulta pública de ordens por código.

O projeto foi desenvolvido como parte do Tech Challenge da FIAP Pós Tech, com foco em organização de domínio, separação de responsabilidades, documentação de API e execução simplificada via Docker.

## Principais recursos

- Autenticação com access token e refresh token.
- Controle de acesso por perfil administrativo.
- Cadastro e manutenção de solicitantes, veículos, serviços e peças.
- Criação, atualização e acompanhamento de ordens de serviço.
- Ajuste de estoque de peças e insumos.
- Consulta pública de ordens de serviço por código.
- Métricas de tempo médio de execução das ordens.
- Documentação interativa via Swagger UI.
- Ambiente completo com API, PostgreSQL, pgAdmin e site de documentação.

## Tecnologias

| Tecnologia | Uso no projeto |
|------------|----------------|
| Go 1.25 | Linguagem principal da API |
| Gin | Framework HTTP e roteamento |
| GORM | ORM para persistência em PostgreSQL |
| PostgreSQL 16 | Banco de dados relacional |
| JWT | Autenticação e autorização |
| Swagger / swaggo | Geração da documentação OpenAPI |
| Docker Compose | Orquestração local dos serviços |
| pgAdmin | Administração visual do banco de dados |

## Arquitetura

A aplicação segue uma abordagem de **Arquitetura Hexagonal (Ports & Adapters)** em um monolito modular. O objetivo é manter a regra de negócio isolada de detalhes externos, como framework HTTP, banco de dados e infraestrutura.

```text
cmd/api/
  main.go                 # entrypoint da aplicação
  bootstrap/              # injeção de dependências
  routes/                 # configuração das rotas
  middleware/             # autenticação JWT e CORS

internal/
  domain/                 # entidades e regras de negócio
  ports/                  # contratos de repositories e use cases
  adapters/
    inbound/http/         # handlers HTTP com Gin
    outbound/postgresql/  # repositories com GORM
  application/usecase/    # casos de uso da aplicação

config/                   # configurações da aplicação
docs/                     # Swagger e documentação do projeto
scripts/                  # scripts auxiliares
```

## Banco de dados

O PostgreSQL foi escolhido por oferecer recursos importantes para o domínio da aplicação:

- **Integridade relacional:** solicitantes, veículos e ordens de serviço possuem relações fortes que são protegidas por chaves estrangeiras.
- **UUID nativo:** geração de identificadores com `gen_random_uuid()` sem dependência externa.
- **JSONB:** armazenamento flexível de value objects, como serviços e peças dentro da ordem de serviço, sem perder capacidade de consulta.
- **Transações ACID:** consistência em operações críticas, como atualização de estoque e mudanças de status.
- **Ecossistema maduro:** integração com Go via `pgx`, GORM e ferramentas como pgAdmin.

## Como rodar

### Pré-requisitos

- Docker e Docker Compose.
- Go 1.23+ para execução local sem container.
- `swag` para regenerar a documentação Swagger localmente.

### Ambiente completo com Docker

```bash
cp .env.example .env
docker compose up --build
```

Também é possível usar o Makefile:

```bash
make up
```

| Serviço | URL | Observação |
|---------|-----|------------|
| API | `http://localhost:8080` | Serviço principal |
| Swagger UI | `http://localhost:8080/swagger/index.html` | Documentação interativa da API |
| pgAdmin | `http://localhost:8082` | Login: `admin@workshop.com` / senha: `admin` |
| Documentação | `http://localhost:8083` | Site estático do projeto |

### Execução local sem Docker

```bash
cp .env.example .env

go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

go run ./cmd/api
```

Com Makefile:

```bash
make swagger
make run
```

### Comandos úteis

| Comando | Descrição |
|---------|-----------|
| `make up` | Sobe todos os containers com build |
| `make up-d` | Sobe todos os containers em background |
| `make down` | Para e remove os containers |
| `make logs` | Exibe os logs dos containers |
| `make swagger` | Gera a documentação Swagger |
| `make tidy` | Organiza as dependências Go |
| `make test` | Executa os testes com cobertura |

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `APP_PORT` | `8080` | Porta da API |
| `APP_ENV` | `development` | Ambiente da aplicação |
| `POSTGRES_DSN` | `postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable` | DSN do PostgreSQL |
| `JWT_SECRET` | `change-me-in-production` | Chave secreta para assinatura JWT |
| `JWT_EXPIRATION_HOURS` | `24` | Expiração do token JWT em horas |
| `ACCESS_TOKEN_EXP_MIN` | `15` | Expiração do access token em minutos |
| `REFRESH_TOKEN_EXP_DAYS` | `7` | Expiração do refresh token em dias |
| `BCRYPT_COST` | `12` | Custo do hash de senha com bcrypt |
| `MAX_FAILED_LOGINS` | `5` | Tentativas permitidas antes do bloqueio de login |
| `LOGIN_LOCK_MIN` | `15` | Tempo de bloqueio após falhas de login, em minutos |

## Endpoints

As rotas de cadastro, login e refresh são públicas. As demais rotas sob `/api/v1/*` exigem o header `Authorization: Bearer <token>`. Rotas marcadas como `admin` também exigem `role=admin` nas claims do JWT.

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/health` | público | Health check da aplicação |
| **Auth** | | | |
| POST | `/api/v1/auth/register` | público | Registra usuário com perfil `client` |
| POST | `/api/v1/auth/login` | público | Autentica usuário e retorna access token e refresh token |
| POST | `/api/v1/auth/refresh` | público | Rotaciona tokens |
| POST | `/api/v1/auth/logout` | JWT | Encerra sessão |
| **Requesters** | | | |
| POST | `/api/v1/requesters` | admin | Cria solicitante |
| GET | `/api/v1/requesters` | admin | Lista solicitantes |
| GET | `/api/v1/requesters/:id` | admin | Busca solicitante por ID |
| GET | `/api/v1/requesters/document/:document` | admin | Busca solicitante por CPF/CNPJ |
| PUT | `/api/v1/requesters/:id` | admin | Atualiza solicitante |
| DELETE | `/api/v1/requesters/:id` | admin | Remove solicitante |
| **Vehicles** | | | |
| POST | `/api/v1/vehicles` | admin | Cadastra veículo |
| GET | `/api/v1/vehicles` | admin | Lista veículos |
| GET | `/api/v1/vehicles/:id` | admin | Busca veículo por ID |
| GET | `/api/v1/requesters/:id/vehicles` | admin | Lista veículos por solicitante |
| PUT | `/api/v1/vehicles/:id` | admin | Atualiza veículo |
| DELETE | `/api/v1/vehicles/:id` | admin | Remove veículo |
| **Services** | | | |
| POST | `/api/v1/services` | admin | Cadastra serviço |
| GET | `/api/v1/services` | admin | Lista serviços |
| GET | `/api/v1/services/:id` | admin | Busca serviço por ID |
| PUT | `/api/v1/services/:id` | admin | Atualiza serviço |
| DELETE | `/api/v1/services/:id` | admin | Remove serviço |
| **Parts** | | | |
| POST | `/api/v1/parts` | admin | Cadastra peça ou insumo |
| GET | `/api/v1/parts` | admin | Lista peças e insumos |
| GET | `/api/v1/parts/:id` | admin | Busca peça por ID |
| PUT | `/api/v1/parts/:id` | admin | Atualiza peça |
| DELETE | `/api/v1/parts/:id` | admin | Remove peça |
| PATCH | `/api/v1/parts/:id/stock` | admin | Ajusta estoque por delta positivo ou negativo |
| **Service Orders** | | | |
| POST | `/api/v1/service-orders` | admin | Cria ordem de serviço |
| GET | `/api/v1/service-orders` | admin | Lista ordens de serviço |
| GET | `/api/v1/service-orders/:id` | admin | Busca ordem por ID |
| GET | `/api/v1/service-orders/metrics/execution-time` | admin | Retorna tempo médio de execução das ordens |
| GET | `/api/v1/service-orders/code/:code` | público | Busca ordem por código |
| GET | `/api/v1/service-orders/requester` | público | Lista ordens por documento do solicitante |
| PUT | `/api/v1/service-orders/code/:code/status` | público | Atualiza status por código |
| PUT | `/api/v1/service-orders/:id/status` | admin | Atualiza status da ordem |
| PUT | `/api/v1/service-orders/:id` | admin | Atualiza ordem de serviço |
| DELETE | `/api/v1/service-orders/:id` | admin | Remove ordem de serviço |

## Testes

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Ou via Makefile:

```bash
make test
```

## Contribuidores

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/yuriLpadlipskas">
        <img src="https://github.com/yuriLpadlipskas.png?size=100" width="100px;" alt="Avatar de yuriLpadlipskas"/><br />
        <sub><b>yuriLpadlipskas</b></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/kc1t">
        <img src="https://github.com/kc1t.png?size=100" width="100px;" alt="Avatar de kc1t"/><br />
        <sub><b>kc1t</b></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/diegoliveiraa">
        <img src="https://github.com/diegoliveiraa.png?size=100" width="100px;" alt="Avatar de diegoliveiraa"/><br />
        <sub><b>diegoliveiraa</b></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/PedroHCarlini">
        <img src="https://github.com/PedroHCarlini.png?size=100" width="100px;" alt="Avatar de PedroHCarlini"/><br />
        <sub><b>PedroHCarlini</b></sub>
      </a>
    </td>
  </tr>
</table>
