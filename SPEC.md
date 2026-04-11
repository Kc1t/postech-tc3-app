# Spec — O que falta no codigo

Mapeamento baseado no PDF "15SOAT - Fase 1 - Tech Challenge".

Legenda: OK | PARCIAL | STUB (estrutura existe mas nao funciona) | FALTA

---

## Handlers (30 stubs)

Todos os handlers retornam `StatusNotImplemented`.
Repo + UseCase + Model existem e funcionam. Falta conectar no handler.

| Handler | Arquivo | Status |
|---------|---------|--------|
| POST /customers | customer/create.go | STUB |
| GET /customers | customer/list.go | STUB |
| GET /customers/:id | customer/find_by_id.go | STUB |
| GET /customers/document/:doc | customer/find_by_document.go | STUB |
| PUT /customers/:id | customer/update.go | STUB |
| DELETE /customers/:id | customer/delete.go | STUB |
| POST /vehicles | vehicle/create.go | STUB |
| GET /vehicles | vehicle/list.go | STUB |
| GET /vehicles/:id | vehicle/find_by_id.go | STUB |
| GET /customers/:id/vehicles | vehicle/list_by_customer.go | STUB |
| PUT /vehicles/:id | vehicle/update.go | STUB |
| DELETE /vehicles/:id | vehicle/delete.go | STUB |
| POST /services | service/create.go | STUB |
| GET /services | service/list.go | STUB |
| GET /services/:id | service/find_by_id.go | STUB |
| PUT /services/:id | service/update.go | STUB |
| DELETE /services/:id | service/delete.go | STUB |
| POST /parts | part/create.go | STUB |
| GET /parts | part/list.go | STUB |
| GET /parts/:id | part/find_by_id.go | STUB |
| PUT /parts/:id | part/update.go | STUB |
| DELETE /parts/:id | part/delete.go | STUB |
| PATCH /parts/:id/stock | part/adjust_stock.go | STUB |
| POST /service-orders | service_order/create.go | STUB |
| GET /service-orders | service_order/list.go | STUB |
| GET /service-orders/:id | service_order/find_by_id.go | STUB |
| GET /customers/:id/service-orders | service_order/list_by_customer.go | STUB |
| PUT /service-orders/:id/status | service_order/update_status.go | STUB |
| PUT /service-orders/:id | service_order/update.go | STUB |
| DELETE /service-orders/:id | service_order/delete.go | STUB |

---

## Autenticacao (JWT)

| Item | Status | Detalhe |
|------|--------|---------|
| Middleware de validacao JWT | OK | `cmd/api/middleware/auth.go` funciona |
| Rotas protegidas por JWT | OK | Todas as `/api/v1/*` passam pelo middleware |
| Endpoint de login | FALTA | Nao existe forma de gerar token |
| Endpoint de registro | FALTA | Nao existe cadastro de usuario do sistema |
| Entidade User/Admin | FALTA | Nao existe dominio de usuario |
| Repositorio de usuarios | FALTA | Nao existe |
| Hash de senha | FALTA | Nao existe |
| Claims no token | PARCIAL | Middleware seta `c.Set("claims", claims)` mas ninguem consome |
| Expiracao configuravel | PARCIAL | `Config.JWTExpirationHours` existe mas ninguem usa (nao gera token) |

---

## Validacoes de dominio (TODOs nos use cases)

| Validacao | Arquivo | Linha |
|-----------|---------|-------|
| Validar CPF/CNPJ + duplicidade | usecase/customer/create.go | 19 |
| Validar placa + cliente existe | usecase/vehicle/create.go | 20 |
| Validar cliente e veiculo + calcular total | usecase/service_order/create.go | 35 |
| Maquina de estados (transicoes validas) | usecase/service_order/update_status.go | 19 |
| Apenas OS "received" pode ser deletada | usecase/service_order/delete.go | 18 |
| Nao permitir estoque negativo | usecase/part/adjust_stock.go | 18 |
| Verificar dependencias ao deletar cliente | usecase/customer/delete.go | 18 |
| Verificar OS abertas ao deletar veiculo | usecase/vehicle/delete.go | 18 |
| Verificar uso em OS ao deletar servico | usecase/service/delete.go | 18 |

---

## Funcionalidades faltantes

| Funcionalidade | Spec diz | Codigo tem |
|----------------|----------|------------|
| Orcamento automatico | "Orcamento gerado automaticamente com base nos servicos e pecas" | `recalcTotal()` existe na entidade mas nao e chamado no fluxo de criacao |
| Baixa de estoque na OS | Implicito pelo controle de estoque | Nao desconta estoque ao criar OS com pecas |
| Tempo medio de execucao | "Monitoramento do tempo medio de execucao dos servicos" | Nenhum endpoint ou logica |
| Consulta do cliente | "Permitir consulta por parte do cliente via API" | Endpoint existe (stub), mas nao tem diferenciacao de role (admin vs cliente) |

---

## Testes

| Item | Status |
|------|--------|
| Arquivos _test.go | ZERO |
| Testes unitarios | ZERO |
| Testes de integracao | ZERO |
| Cobertura 80% dominios criticos | ZERO |

---

## Infra e docs

| Item | Status | Detalhe |
|------|--------|---------|
| Dockerfile | OK | Multi-stage build funcionando |
| docker-compose.yml | OK | API + Postgres + pgAdmin |
| README.md | DESATUALIZADO | Menciona MongoDB, variaveis erradas |
| Swagger | PARCIAL | Funciona mas todos endpoints retornam 501 |
| Scan de vulnerabilidades | FALTA | Nao foi rodado (spec pede relatorio) |
