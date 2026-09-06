<div align="center">

# Workshop API

Sistema integrado de atendimento e execução de serviços para oficinas mecânicas.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-HTTP%20Framework-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/Kubernetes-EKS-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-IaC-7B42BC?style=for-the-badge&logo=terraform&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-CI%2FCD-2088FF?style=for-the-badge&logo=githubactions&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

</div>

## Fase 3 — segurança, escalabilidade e observabilidade

A Fase 3 quebrou o projeto em **quatro repositórios**, cada um com CI/CD próprio. Este aqui é a aplicação; a infraestrutura mora nos outros três.

| Repositório | Responsabilidade |
|---|---|
| **`postech-tc3-app`** (este) | Aplicação Go rodando no EKS |
| [`postech-tc3-lambda-auth`](https://github.com/Kc1t/postech-tc3-lambda-auth) | Autenticação serverless por CPF e authorizer do gateway |
| [`postech-tc3-infra-k8s`](https://github.com/Kc1t/postech-tc3-infra-k8s) | Cluster EKS, metrics-server, API Gateway e agente New Relic |
| [`postech-tc3-infra-database`](https://github.com/Kc1t/postech-tc3-infra-database) | RDS PostgreSQL gerenciado |

```mermaid
flowchart TB
    cli(["Cliente<br/>dono do veículo"]) -->|CPF| gw["API Gateway"]
    op(["Operação<br/>atendente · mecânico · admin"]) -->|e-mail e senha| gw

    gw -.->|autoriza| az["Lambda authorizer"]
    gw -->|"POST /auth"| iss["Lambda issuer"]
    gw -->|"HTTP_PROXY"| nlb["NLB"]

    nlb --> app["<b>workshop-api</b><br/>este repositório<br/>2 a 10 pods no EKS"]

    iss --> db[("RDS PostgreSQL")]
    app --> db
    app -->|"JSON + correlation id"| nr["New Relic<br/>APM · logs · dashboards"]

    style app fill:#e6f4ea,stroke:#4a7
    style az fill:#fff4e0,stroke:#d90
```

O que mudou na aplicação nesta fase:

- **Autenticação por CPF.** As rotas de consulta e aprovação de OS, públicas na Fase 2, passaram a exigir o JWT emitido pela Lambda a partir do CPF. Ver [`docs/rfc/0003-autenticacao-por-cpf.md`](docs/rfc/0003-autenticacao-por-cpf.md).
- **`requesters.status`.** Nova coluna que permite à Lambda consultar existência **e** status do cliente. Ver [`docs/MODELAGEM_DE_DADOS.md`](docs/MODELAGEM_DE_DADOS.md).
- **Logs estruturados em JSON com correlação.** Todo request recebe um `X-Correlation-ID`, devolvido na resposta e presente em cada log.
- **Deploy por branch em namespaces separados.** `homolog` entrega em `postech-homolog`, `main` em `postech`, no mesmo cluster. Ver [`docs/adr/0010-cluster-unico-dois-namespaces.md`](docs/adr/0010-cluster-unico-dois-namespaces.md).
- **Posse do documento.** Não basta ter um token válido: o documento consultado precisa ser o do próprio token, senão a resposta é 403. Sem isso, qualquer cliente autenticado leria as OS de outro conhecendo apenas o CPF.

### Camadas de proteção

```mermaid
flowchart LR
    r(["Requisição"]) --> l1{{"1 · Gateway<br/>authorizer Lambda"}}
    l1 -->|nega| e1(["401"])
    l1 --> l2{{"2 · Middleware Auth<br/>valida o mesmo JWT"}}
    l2 -->|nega| e2(["401"])
    l2 --> l3{{"3 · Papel ou posse<br/>RequireRole · authorizeDocument"}}
    l3 -->|nega| e3(["403"])
    l3 --> h(["handler"])

    style e1 fill:#ffe4e1,stroke:#c66
    style e2 fill:#ffe4e1,stroke:#c66
    style e3 fill:#ffe0b2,stroke:#e80
```

As três camadas são independentes. Alcançar o NLB direto, contornando o gateway, ainda esbarra na segunda e na terceira. O grafo completo das 37 rotas está em [`docs/diagrams/roteamento.md`](docs/diagrams/roteamento.md).

### Documentação da Fase 3

| Documento | Conteúdo |
|---|---|
| [`docs/diagrams/componentes.md`](docs/diagrams/componentes.md) | Diagrama de componentes com visão de nuvem |
| [`docs/diagrams/roteamento.md`](docs/diagrams/roteamento.md) | Grafos de roteamento: gateway, aplicação, inventário das 37 rotas e mapa de códigos de erro |
| [`docs/diagrams/pipelines.md`](docs/diagrams/pipelines.md) | Grafos dos quatro pipelines e matriz de secrets por repositório |
| [`docs/diagrams/sequencia-autenticacao.md`](docs/diagrams/sequencia-autenticacao.md) | Sequência da autenticação por CPF |
| [`docs/diagrams/sequencia-ordem-servico.md`](docs/diagrams/sequencia-ordem-servico.md) | Sequência da abertura de OS |
| [`docs/MODELAGEM_DE_DADOS.md`](docs/MODELAGEM_DE_DADOS.md) | Justificativa do banco, relacionamentos e diagrama ER |
| [`docs/rfc/README.md`](docs/rfc/README.md) | RFCs: nuvem, banco gerenciado, autenticação |
| [`docs/adr/README.md`](docs/adr/README.md) | ADRs 0001–0010 |
| [`docs/newrelic/README.md`](docs/newrelic/README.md) | Observabilidade: dashboards, alertas e origem de cada métrica |
| [`docs/DOCUMENTO_ENTREGA_FASE3.md`](docs/DOCUMENTO_ENTREGA_FASE3.md) | Documento de entrega, base do PDF do portal |
| [`docs/ROTEIRO_VIDEO.md`](docs/ROTEIRO_VIDEO.md) | Roteiro do vídeo de demonstração |

---

## Guia rápido para avaliação

Esta seção abaixo é o **hub técnico da entrega da Fase 2**, mantido como referência histórica.

Para avaliar o projeto com mais contexto, recomenda-se seguir esta ordem:

1. **Entendimento da entrega:** [`docs/DOCUMENTO_ENTREGA_FASE2.md`](docs/DOCUMENTO_ENTREGA_FASE2.md)
2. **Visão da solução:** [Sobre o projeto](#sobre-o-projeto), [Objetivos da Fase 2](#objetivos-da-fase-2) e [Arquitetura](#arquitetura)
3. **Fluxo de negócio da Fase 2:** [Ordem de Serviço](#ordem-de-serviço-fluxo-da-fase-2) e [`postman_collection.json`](postman_collection.json)
4. **Infraestrutura e deploy:** [`k8s/README.md`](k8s/README.md), [CI/CD](#cicd) e os repositórios de infraestrutura da Fase 3
5. **Decisões técnicas:** [`docs/adr/README.md`](docs/adr/README.md)
6. **Segurança e qualidade:** [`docs/security-reports/RELATORIO.md`](docs/security-reports/RELATORIO.md) e [Testes](#testes)
7. **Demonstração em vídeo:** https://www.youtube.com/watch?v=PiraAX3RVzg

## Checklist da entrega

| Item exigido | Referência |
|--------------|------------|
| Repositório do projeto | https://github.com/Kc1t/postech-tc1 |
| README com solução, arquitetura, deploy e execução | Este arquivo |
| Desenho da arquitetura com componentes, infraestrutura e fluxo de deploy | [Arquitetura](#arquitetura) e [`docs/documentation-diagram.drawio`](docs/documentation-diagram.drawio) |
| Collection das APIs | [`postman_collection.json`](postman_collection.json) |
| Swagger estático | [`docs/swagger.yaml`](docs/swagger.yaml) |
| Hub de documentação | https://tc-doc.vercel.app/ |
| Vídeo demonstrativo da Fase 2 | https://www.youtube.com/watch?v=PiraAX3RVzg |
| Documento formal da entrega | [`docs/DOCUMENTO_ENTREGA_FASE2.md`](docs/DOCUMENTO_ENTREGA_FASE2.md) |
| Conteúdo para o PDF do portal | Repositório, desenho da arquitetura e link do vídeo estão consolidados neste checklist e no documento formal da entrega. |

## Sumário

- [Guia rápido para avaliação](#guia-rápido-para-avaliação)
- [Checklist da entrega](#checklist-da-entrega)
- [Sobre o projeto](#sobre-o-projeto)
- [Objetivos da Fase 2](#objetivos-da-fase-2)
- [Rastreabilidade da Fase 2](#rastreabilidade-da-fase-2)
- [Principais recursos](#principais-recursos)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
  - [Componentes da aplicação](#componentes-da-aplicação)
  - [Infraestrutura provisionada](#infraestrutura-provisionada)
  - [Fluxo de deploy (CI/CD)](#fluxo-de-deploy-cicd)
- [Ordem de Serviço (fluxo da Fase 2)](#ordem-de-serviço-fluxo-da-fase-2)
- [Banco de dados](#banco-de-dados)
- [Índice de documentação](#índice-de-documentação)
- [Como rodar](#como-rodar)
  - [Execução local com Docker](#execução-local-com-docker)
  - [Deploy em Kubernetes](#deploy-em-kubernetes)
  - [Provisionamento da infraestrutura com Terraform](#provisionamento-da-infraestrutura-com-terraform)
- [CI/CD](#cicd)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Endpoints](#endpoints)
- [Testes](#testes)
- [Entregáveis](#entregáveis)
- [Contribuidores](#contribuidores)

## Sobre o projeto

O **Workshop API** é uma API REST para gestão do fluxo operacional de uma oficina mecânica. A aplicação centraliza o cadastro de clientes, veículos, serviços, peças e ordens de serviço, além de oferecer autenticação JWT, controle de perfis administrativos e consulta pública de ordens por código.

O projeto foi desenvolvido como parte do Tech Challenge da **FIAP Pós Tech — SOAT**. Na **Fase 1**, a aplicação foi construída como um monolito modular em Arquitetura Hexagonal com foco em domínio e documentação. Na **Fase 2**, a aplicação evoluiu para garantir **qualidade, resiliência e escalabilidade**, incorporando containerização, orquestração com Kubernetes, Infraestrutura como Código (Terraform) e uma pipeline de CI/CD completa.

## Objetivos da Fase 2

A Fase 2 evolui a aplicação da Fase 1 com práticas modernas de infraestrutura e automação, atendendo às seguintes metas:

- **Reduzir riscos operacionais** com infraestrutura escalável (Kubernetes + HPA na AWS EKS).
- **Automatizar provisionamento e deploy** do ambiente (Terraform + GitHub Actions).
- **Melhorar a qualidade e a organização do código**, mantendo a evolução sustentável (Clean Code + Arquitetura Hexagonal + testes com cobertura mínima de 80%).
- **Suportar picos de demanda** com escalabilidade dinâmica (Horizontal Pod Autoscaler por CPU e memória).

## Rastreabilidade da Fase 2

| Requisito avaliado | Implementação / evidência |
|--------------------|---------------------------|
| Evolução da aplicação da Fase 1 | Monolito Go com Arquitetura Hexagonal, novos fluxos de OS, autenticação JWT, notificação por e-mail e endpoints públicos para cliente. |
| Clean Code e organização | Separação entre `internal/domain`, `internal/application/usecase`, `internal/ports` e `internal/adapters`; um use case por operação; ADRs registrando decisões. |
| Testes automatizados e cobertura | Arquivos `_test.go` por domínio/use case/adapter; pipeline com gate de cobertura mínima de 80%. |
| Abertura de Ordem de Serviço | `POST /api/v1/service-orders`, use case `internal/application/usecase/service_order/create.go` e handler HTTP correspondente. |
| Consulta de status da OS | `GET /api/v1/service-orders/code/:code?document=<cpf/cnpj>` e `GET /api/v1/service-orders/:id`. |
| Aprovação/recusa de orçamento | `PUT /api/v1/service-orders/code/:code/status`, endpoint público com validação do documento do solicitante. |
| Listagem ordenada de OS | `GET /api/v1/service-orders`, ordenação por prioridade de status e data de criação, excluindo finalizadas/entregues da visão operacional. |
| Notificação por e-mail | Adapter SMTP em `internal/adapters/outbound/smtp/email_notifier.go`. |
| Docker e execução local | [`Dockerfile`](Dockerfile), [`docker-compose.yml`](docker-compose.yml), [Execução local com Docker](#execução-local-com-docker). |
| Kubernetes | Manifestos em [`k8s/`](k8s/) com Namespace, Deployment, Service, ConfigMap, Secret e HPA. |
| Escalabilidade automática | [`k8s/hpa.yaml`](k8s/hpa.yaml), 2 a 10 réplicas por CPU/memória. |
| Infraestrutura como Código | [`infra/`](infra/) provisionando EKS e RDS PostgreSQL via Terraform. |
| CI/CD | [`.github/workflows/ci.yml`](.github/workflows/ci.yml) com lint, dependências, testes, build/push no ECR, deploy no EKS e aplicação dos manifestos YAML. |
| Banco de dados em infraestrutura | RDS PostgreSQL 16 provisionado por Terraform em [`postech-tc3-infra-database`](https://github.com/Kc1t/postech-tc3-infra-database), com DSN consumido pelo Secret do Kubernetes. |
| APIs documentadas | Swagger UI, [`docs/swagger.yaml`](docs/swagger.yaml) e [`postman_collection.json`](postman_collection.json). |
| Entrega formal da Fase 2 | [`docs/DOCUMENTO_ENTREGA_FASE2.md`](docs/DOCUMENTO_ENTREGA_FASE2.md), [Checklist da entrega](#checklist-da-entrega), vídeo e links oficiais em [Entregáveis](#entregáveis). |
| Vídeo demonstrativo | Demonstra deploy, execução da pipeline, consumo das APIs e escalabilidade automática: https://www.youtube.com/watch?v=PiraAX3RVzg |

## Principais recursos

- Autenticação com access token e refresh token rotativo.
- Controle de acesso por perfil administrativo.
- Cadastro e manutenção de solicitantes, veículos, serviços e peças.
- **Abertura de Ordem de Serviço** com cliente, veículo, serviços e peças, retornando o identificador único da OS.
- **Consulta do status da OS** (Recebida, Diagnóstico, Aguardando Aprovação, Execução, Finalizada, Entregue).
- **Aprovação/recusa de orçamento** por endpoint público (recebe notificações externas do cliente).
- **Listagem de OS ordenada** por prioridade de status e mais antigas primeiro, excluindo (logicamente) as finalizadas e entregues.
- **Notificação por e-mail** a cada mudança de status relevante (SMTP).
- Ajuste de estoque de peças e insumos.
- Métricas de tempo médio de execução das ordens.
- Documentação interativa via Swagger UI.

## Tecnologias

| Camada | Tecnologia | Uso no projeto |
|--------|------------|----------------|
| Aplicação | Go 1.25 | Linguagem principal da API |
| Aplicação | Gin | Framework HTTP e roteamento |
| Aplicação | GORM | ORM para persistência em PostgreSQL |
| Aplicação | JWT | Autenticação e autorização |
| Aplicação | SMTP (Brevo) | Notificação de mudança de status por e-mail |
| Aplicação | Swagger / swaggo | Geração da documentação OpenAPI |
| Dados | PostgreSQL 16 | Banco de dados relacional (local e AWS RDS) |
| Container | Docker / Docker Compose | Build e orquestração local dos serviços |
| Orquestração | Kubernetes (AWS EKS) | Deploy, Services, ConfigMaps, Secrets e HPA |
| IaC | Terraform | Provisionamento do cluster EKS e do RDS |
| Registry | AWS ECR | Registro das imagens Docker |
| CI/CD | GitHub Actions | Lint, testes, build, push e deploy automatizados |

## Arquitetura

A aplicação segue **Arquitetura Hexagonal (Ports & Adapters)** em um monolito modular. A regra de negócio fica isolada de detalhes externos (framework HTTP, banco de dados, JWT, SMTP e infraestrutura), o que garante testabilidade e permite trocar adapters sem tocar no domínio.

### Componentes da aplicação

```text
                 Inbound Adapters                 Application            Domain
              (handlers HTTP - Gin)               (use cases)      (entidades + VOs)
        ┌──────────────────────────────┐     ┌────────────────┐   ┌──────────────────┐
HTTP →  │ Auth / Requester / Vehicle /  │ →   │  Execute()     │ → │ ServiceOrder,    │
        │ Service / Part / ServiceOrder │     │  (1 por caso)  │   │ User, Vehicle,   │
        └──────────────────────────────┘     └───────┬────────┘   │ Part, Service... │
                                                      │            │ regras + máquina │
                                                      │            │ de estados da OS │
                                              ports (interfaces)   └──────────────────┘
                                                      │
        ┌─────────────────────────────────────────────┴───────────────────────────┐
        │                         Outbound Adapters                                 │
        │  PostgreSQL (GORM)   ·   JWT (token service)   ·   SMTP (email notifier)  │
        └───────────────────────────────────────────────────────────────────────────┘

cmd/api/            # entrypoint, bootstrap (DI), rotas e middleware
internal/domain/    # entidades, value objects e erros de domínio
internal/ports/     # contratos (repositories, use cases, token, hasher, notifier)
internal/application/usecase/   # casos de uso (1 arquivo por operação)
internal/adapters/inbound/http/ # handlers HTTP (Gin)
internal/adapters/outbound/     # postgresql (GORM), jwt, smtp
pkg/                # código genérico e reutilizável (env, hasher, ioc)
```

**Regra central:** um use case nunca importa GORM, JWT ou bcrypt — apenas interfaces de `ports/`. As implementações concretas vivem em `internal/adapters/outbound/<tech>/` ou em `pkg/` quando são genéricas.

### Infraestrutura provisionada

O ambiente de produção roda na **AWS**, provisionado por Terraform e orquestrado por Kubernetes (EKS). O tráfego externo entra por um Network Load Balancer, e o HPA escala os pods conforme a carga.

```text
                          ┌────────────────────────── AWS ──────────────────────────┐
                          │                                                          │
   Internet  ───────────► │  NLB (Service LoadBalancer)                              │
                          │        │                                                 │
                          │        ▼                                                 │
                          │  ┌──────────────── EKS Cluster (namespace: postech) ───┐ │
                          │  │                                                      │ │
                          │  │   Deployment: workshop-api (2..10 réplicas)          │ │
                          │  │     ├─ ConfigMap  (config não sensível)              │ │
                          │  │     ├─ Secret     (JWT, DSN, SMTP, admin)            │ │
                          │  │     └─ HPA        (CPU 70% / Mem 80% → 2..10 pods)   │ │
                          │  │            │ readiness/liveness em /health           │ │
                          │  └────────────┼─────────────────────────────────────────┘ │
                          │               │                                          │
                          │               ▼                                          │
                          │        RDS PostgreSQL 16 (db.t3.micro)                   │
                          │                                                          │
                          │        ECR (imagens Docker: workshop-api)               │
                          └──────────────────────────────────────────────────────────┘
```

Recursos criados pelo Terraform (agora em [`postech-tc3-infra-k8s`](https://github.com/Kc1t/postech-tc3-infra-k8s) e [`postech-tc3-infra-database`](https://github.com/Kc1t/postech-tc3-infra-database)):

- **Cluster EKS** (`terraform-aws-modules/eks`) com managed node group.
- **RDS PostgreSQL 16** com subnet group e security group dedicado.
- **Outputs** com endpoint do cluster, comando de kubeconfig e DSN do banco.

Recursos Kubernetes (detalhes em [`k8s/README.md`](k8s/README.md)):

- **Namespace** `postech`, **Deployment**, **Service** (LoadBalancer/NLB).
- **ConfigMap** (variáveis não sensíveis) e **Secret** (JWT, DSN, SMTP, admin).
- **HorizontalPodAutoscaler** (2 a 10 réplicas, CPU 70% / memória 80%).

### Fluxo de deploy (CI/CD)

```text
  git push (main)
       │
       ▼
  GitHub Actions ── lint (golangci-lint)
       ├──────────── dependencies (go mod tidy + govulncheck)
       ├──────────── tests & coverage (Postgres service, gate ≥ 80%)
       │
       └── deploy (apenas em push na main, após os jobs acima passarem)
              1. Configure AWS credentials
              2. Login no ECR
              3. docker build → push (tag = git SHA + latest)
              4. aws eks update-kubeconfig
              5. kubectl apply: namespace → secret → configmap → deployment → service → hpa
              6. kubectl rollout status (aguarda o rollout, timeout 300s)
```

## Ordem de Serviço (fluxo da Fase 2)

A OS possui uma **máquina de estados** validada no domínio (`internal/domain/entities/service_order.go`). Os nomes internos (enum) mapeiam os status do enunciado:

| Status (enunciado) | Enum interno | Transições permitidas |
|--------------------|--------------|-----------------------|
| Recebida | `received` | → Diagnóstico |
| Diagnóstico | `in_diagnosis` | → Aguardando Aprovação |
| Aguardando Aprovação | `awaiting_approval` | → Execução (aprovar) · → Recebida (recusar) |
| Execução | `in_execution` | → Finalizada |
| Finalizada | `finished` | → Entregue |
| Entregue | `delivered` | — |

- **Abertura da OS:** `POST /api/v1/service-orders` recebe cliente, veículo, serviços e peças e retorna o identificador único.
- **Consulta de status:** `GET /api/v1/service-orders/code/:code` (público) e `GET /api/v1/service-orders/:id` (admin).
- **Aprovação/recusa de orçamento:** `PUT /api/v1/service-orders/code/:code/status` (público) — o cliente aprova (→ Execução) ou recusa (→ Recebida) o orçamento. As demais transições são restritas ao fluxo operacional da oficina.
- **Listagem:** `GET /api/v1/service-orders` ordena por **Execução > Aguardando Aprovação > Diagnóstico > Recebida**, mais antigas primeiro (`created_at ASC`), e **exclui logicamente** as OS *finalizadas* e *entregues* (elas não somem do banco, apenas da listagem).
- **Atualização de status via e-mail:** a cada mudança relevante (Aguardando Aprovação, Execução, Finalizada, Entregue), o cliente é notificado por e-mail via SMTP (adapter `internal/adapters/outbound/smtp`).

### Exemplos de request

> Endpoints públicos do cliente identificam o dono da OS pelo **documento** (CPF/CNPJ). Nos GETs ele vai na query string; na aprovação, no corpo. Documento que não bate com a OS retorna `404` (não vaza a existência da ordem).

```bash
# 1. Abrir OS (admin) — retorna a OS com o código único
curl -X POST "$BASE/api/v1/service-orders" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"requester_document":"50891498877","vehicle_plate":"ABC1234","notes":"Barulho no motor"}'

# 2. Consultar status (público) — exige o document do dono
curl "$BASE/api/v1/service-orders/code/1?document=50891498877"

# 3. Aprovar orçamento (público) — status "in_execution" aprova, "received" recusa
curl -X PUT "$BASE/api/v1/service-orders/code/1/status" \
  -H "Content-Type: application/json" \
  -d '{"requester_document":"50891498877","status":"in_execution"}'

# 4. Avançar status pelo fluxo da oficina (admin) — services/parts são opcionais
curl -X PUT "$BASE/api/v1/service-orders/<id>/status" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"status":"finished"}'
```

## Banco de dados

O PostgreSQL foi escolhido por oferecer recursos importantes para o domínio da aplicação:

- **Integridade relacional:** solicitantes, veículos e ordens de serviço possuem relações fortes protegidas por chaves estrangeiras.
- **UUID nativo:** identificadores gerados com `gen_random_uuid()` sem dependência externa.
- **JSONB:** armazenamento flexível de value objects (serviços e peças da OS) sem perder capacidade de consulta.
- **Transações ACID:** consistência em operações críticas como atualização de estoque e mudança de status.
- **Ecossistema maduro:** integração com Go via `pgx`, GORM e ferramentas como pgAdmin. Em produção, o mesmo engine roda no **AWS RDS**.

## Índice de documentação

| Documento / artefato | Finalidade |
|----------------------|------------|
| [`docs/DOCUMENTO_ENTREGA_FASE2.md`](docs/DOCUMENTO_ENTREGA_FASE2.md) | Documento objetivo da entrega da Fase 2, com grupo, links, requisitos e roteiro do vídeo. |
| [`docs/DOCUMENTO_ENTREGA.md`](docs/DOCUMENTO_ENTREGA.md) | Documento da Fase 1, útil para entender a origem do domínio e a evolução da solução. |
| [`postech-tc3-infra-k8s`](https://github.com/Kc1t/postech-tc3-infra-k8s) e [`postech-tc3-infra-database`](https://github.com/Kc1t/postech-tc3-infra-database) | Terraform da infraestrutura, movido para repositórios próprios na Fase 3. |
| [`k8s/README.md`](k8s/README.md) | Detalha os manifestos Kubernetes, ordem de aplicação, HPA, metrics-server e Secret da pipeline. |
| [`docs/adr/README.md`](docs/adr/README.md) | Índice das decisões arquiteturais registradas no projeto. |
| [`docs/adr/0001-go-language.md`](docs/adr/0001-go-language.md) | Justificativa da escolha de Go. |
| [`docs/adr/0002-hexagonal-architecture.md`](docs/adr/0002-hexagonal-architecture.md) | Justificativa da Arquitetura Hexagonal. |
| [`docs/adr/0003-postgresql.md`](docs/adr/0003-postgresql.md) | Justificativa da escolha do PostgreSQL. |
| [`docs/adr/0004-gorm-orm.md`](docs/adr/0004-gorm-orm.md) | Justificativa do uso do GORM. |
| [`docs/adr/0005-jwt-refresh-token.md`](docs/adr/0005-jwt-refresh-token.md) | Decisão sobre autenticação JWT com refresh token rotativo. |
| [`docs/adr/0006-domain-errors-sentinels.md`](docs/adr/0006-domain-errors-sentinels.md) | Decisão sobre erros de domínio centralizados. |
| [`docs/security-reports/RELATORIO.md`](docs/security-reports/RELATORIO.md) | Relatório de análise de vulnerabilidades com `gosec`, `govulncheck` e `trivy`. |
| [`docs/swagger.yaml`](docs/swagger.yaml) | Especificação OpenAPI estática. |
| [`postman_collection.json`](postman_collection.json) | Collection para validação prática dos endpoints. |
| [`docs/database-model.dbml`](docs/database-model.dbml) | Modelo de dados em DBML. |
| [`docs/documentation-diagram.drawio`](docs/documentation-diagram.drawio) | Diagrama editável da documentação/arquitetura. |
| [`docs/site/`](docs/site/) | Site estático de documentação publicado no hub da entrega. |

## Como rodar

### Pré-requisitos

- **Local:** Docker e Docker Compose (ou Go 1.25 para rodar sem container).
- **Kubernetes:** `kubectl` e acesso a um cluster (EKS, Minikube, kind ou Docker Desktop).
- **Terraform:** Terraform ≥ 1.5, AWS CLI configurado com credenciais válidas.

### Execução local com Docker

```bash
cp .env.example .env
docker compose up --build
```

Também é possível usar o Makefile: `make up`.

| Serviço | URL | Observação |
|---------|-----|------------|
| API | `http://localhost:8080` | Serviço principal |
| Swagger UI | `http://localhost:8080/swagger/index.html` | Documentação interativa da API |
| Health check | `http://localhost:8080/health` | Usado por readiness/liveness no K8s |
| pgAdmin | `http://localhost:8082` | Login: `admin@workshop.com` / senha: `admin` |
| Documentação | `http://localhost:8083` | Site estático do projeto |

Execução local sem Docker:

```bash
cp .env.example .env
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
go run ./cmd/api
```

### Deploy em Kubernetes

Os manifestos estão em [`k8s/`](k8s/). Aplicação manual (ordem importa):

Os manifestos não declaram namespace: ele é escolhido na aplicação, o que permite usar os mesmos arquivos nos dois ambientes.

```bash
NAMESPACE=postech-homolog   # ou postech, para produção

# 1. Namespace
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# 2. Configuração e segredos
#    Ajuste os valores do Secret antes (JWT_SECRET, POSTGRES_DSN, SMTP_*, ADMIN_PASSWORD)
kubectl apply -n "$NAMESPACE" -f k8s/configmap.yaml
kubectl apply -n "$NAMESPACE" -f k8s/secret.yaml

# 3. Deployment, Service e HPA
#    Substitua IMAGE_PLACEHOLDER pela imagem publicada no ECR
kubectl apply -n "$NAMESPACE" -f k8s/deployment.yaml
kubectl apply -n "$NAMESPACE" -f k8s/service.yaml
kubectl apply -n "$NAMESPACE" -f k8s/hpa.yaml

# 4. Acompanhar o rollout e a URL pública
kubectl rollout status deployment/workshop-api -n "$NAMESPACE"
kubectl get svc workshop-api -n "$NAMESPACE"
```

> O `JWT_SECRET` precisa ser **o mesmo** configurado nas Lambdas de autenticação, senão os tokens emitidos por CPF são rejeitados com 401.

Verificar a escalabilidade automática:

```bash
kubectl get hpa -n postech -w
kubectl get pods -n postech -w   # observar as réplicas subindo sob carga
```

> Em produção o deploy é automatizado pela pipeline (ver [CI/CD](#cicd)), que também substitui a imagem no manifesto e cria o Secret a partir dos GitHub Secrets. Consulte [`k8s/README.md`](k8s/README.md) para detalhes de cada manifesto.

### Provisionamento da infraestrutura com Terraform

Os scripts estão em [`infra/`](infra/). Resumo:

```bash
cd infra
terraform init
terraform plan  -var="db_password=<senha-forte>"
terraform apply -var="db_password=<senha-forte>"

# Após o apply, configurar o kubectl para o cluster criado:
aws eks update-kubeconfig --name workshop-api --region us-east-1
```

O `terraform apply` cria o cluster EKS e o RDS PostgreSQL. Os outputs incluem o endpoint do cluster, o comando de kubeconfig e o `POSTGRES_DSN` (sensível) a ser usado no Secret do Kubernetes. **O Terraform foi movido para os repositórios de infraestrutura na Fase 3; consulte o README de cada um.**

## CI/CD

A pipeline está em [`.github/workflows/ci.yml`](.github/workflows/ci.yml) (GitHub Actions) e executa:

| Job | Gatilho | O que faz |
|-----|---------|-----------|
| **lint** | PR e push | `golangci-lint` |
| **dependencies** | PR e push | Verifica `go mod tidy` e roda `govulncheck` |
| **tests** | PR e push | Sobe Postgres, roda os testes com cobertura e falha se `< 80%` |
| **deploy** | push na `main` | Build/push da imagem no ECR e deploy no EKS (aplica os manifestos + rollout) |

Secrets necessários no GitHub (Settings → Secrets and variables → Actions): `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`, `AWS_REGION`, `EKS_CLUSTER_NAME`, `JWT_SECRET`, `POSTGRES_DSN`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `ADMIN_PASSWORD`.

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `APP_PORT` | `8080` | Porta da API |
| `APP_ENV` | `development` | Ambiente da aplicação |
| `POSTGRES_DSN` | `postgres://postgres:postgres@localhost:5432/workshop?sslmode=disable` | DSN do PostgreSQL |
| `JWT_SECRET` | `change-me-in-production` | Chave secreta para assinatura JWT |
| `ACCESS_TOKEN_EXP_MIN` | `15` | Expiração do access token em minutos |
| `REFRESH_TOKEN_EXP_DAYS` | `7` | Expiração do refresh token em dias |
| `BCRYPT_COST` | `12` | Custo do hash de senha com bcrypt |
| `MAX_FAILED_LOGINS` | `5` | Tentativas permitidas antes do bloqueio de login |
| `LOGIN_LOCK_MIN` | `15` | Tempo de bloqueio após falhas de login, em minutos |
| `ADMIN_EMAIL` | `admin@workshop.com` | E-mail do admin criado no seed |
| `ADMIN_PASSWORD` | `change-me-in-production` | Senha do admin criado no seed |
| `SMTP_HOST` | `smtp-relay.brevo.com` | Host do servidor SMTP (vazio desativa o envio) |
| `SMTP_PORT` | `587` | Porta SMTP |
| `SMTP_USERNAME` | — | Usuário/login do SMTP |
| `SMTP_PASSWORD` | — | Chave/senha do SMTP |
| `SMTP_FROM` | `noreply@workshop.com` | Remetente das notificações |

## Endpoints

As rotas de cadastro, login e refresh são públicas. As demais rotas sob `/api/v1/*` exigem o header `Authorization: Bearer <token>`. Rotas marcadas como `admin` também exigem `role=admin` nas claims do JWT.

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/health` | público | Health check da aplicação |
| **Auth** | | | |
| POST | `/api/v1/auth/register` | público | Registra usuário com perfil `client` |
| POST | `/api/v1/auth/login` | público | Autentica usuário e retorna access e refresh token |
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
| POST | `/api/v1/service-orders` | admin | **Abre OS** (cliente, veículo, serviços e peças) → retorna ID único |
| GET | `/api/v1/service-orders` | admin | **Lista OS** ordenada por status e mais antigas primeiro (exclui finalizadas/entregues) |
| GET | `/api/v1/service-orders/:id` | admin | Busca OS por ID |
| GET | `/api/v1/service-orders/metrics/execution-time` | admin | Tempo médio de execução das ordens |
| GET | `/api/v1/service-orders/code/:code?document=<cpf/cnpj>` | público | **Consulta status** da OS por código (exige o `document` do dono; documento errado → 404) |
| GET | `/api/v1/service-orders/requester?document=<cpf/cnpj>` | público | Lista ordens do solicitante (exige o query param `document`) |
| PUT | `/api/v1/service-orders/code/:code/status` | público | **Aprova/recusa orçamento** — body com `requester_document` + `status` |
| PUT | `/api/v1/service-orders/:id/status` | admin | Atualiza status da OS (fluxo operacional) |
| PUT | `/api/v1/service-orders/:id` | admin | Atualiza OS |
| DELETE | `/api/v1/service-orders/:id` | admin | Remove OS |

Collection completa das APIs: [`postman_collection.json`](postman_collection.json) · Swagger estático: [`docs/swagger.yaml`](docs/swagger.yaml).

## Testes

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Ou via Makefile: `make test`. A cobertura é validada na pipeline com **gate mínimo de 80%** nos domínios críticos.

## Entregáveis

| Item | Link |
|------|------|
| 📂 Repositório (privado, acesso `soat-architecture`) | https://github.com/Kc1t/postech-tc1 |
| 🔧 Collection das APIs (Postman) | [`postman_collection.json`](postman_collection.json) |
| 📄 Swagger estático | [`docs/swagger.yaml`](docs/swagger.yaml) |
| 🏛️ Hub de documentação | https://tc-doc.vercel.app/ |
| 🎬 Vídeo demonstrativo (Fase 2 · até 15 min) | https://www.youtube.com/watch?v=PiraAX3RVzg |

O vídeo deve demonstrar: deploy da aplicação, execução do CI/CD, consumo das APIs e escalabilidade automática (simulação de carga).

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
      <a href="https://github.com/PedroHCarlini">
        <img src="https://github.com/PedroHCarlini.png?size=100" width="100px;" alt="Avatar de PedroHCarlini"/><br />
        <sub><b>PedroHCarlini</b></sub>
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
  </tr>
</table>
