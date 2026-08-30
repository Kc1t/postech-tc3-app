# postech-tc3-app

**Aplicação principal** do sistema de gestão para oficinas mecânicas, executando em Kubernetes.

Go + PostgreSQL + arquitetura hexagonal. Tech Challenge Fase 3 — FIAP Pós Tech SOAT.

## Propósito

Continuação do [`postech-tc1`](https://github.com/Kc1t/postech-tc1) (Fases 1 e 2), agora containerizada, rodando no EKS e exposta atrás do API Gateway, com rotas sensíveis protegidas pelo JWT emitido pela lambda de autenticação.

> **Status:** o código da aplicação ainda não foi migrado do `postech-tc1`. Este repositório contém, por ora, os manifests do Kubernetes, o pipeline de CI/CD e a estrutura de documentação arquitetural.

## Tecnologias

- Go 1.23, Gin, GORM, PostgreSQL 16
- Arquitetura hexagonal (ports & adapters)
- Docker + Amazon ECR
- Kubernetes (EKS) com HPA
- CI/CD: GitHub Actions com OIDC

## Arquitetura

```
   Cliente ──▶ API Gateway ──▶ Lambda auth (CPF → JWT)
                    │
                    │ Bearer JWT
                    ▼
        ┌───────────────────────────┐
        │  EKS                      │
        │   Service (ClusterIP)     │
        │        │                  │
        │   Deployment workshop-api │
        │        │  HPA 2..10       │
        └────────┼──────────────────┘
                 ▼
           RDS PostgreSQL
```

## Repositórios relacionados

| Repo | Papel |
|---|---|
| [`postech-tc3-lambda-auth`](https://github.com/Kc1t/postech-tc3-lambda-auth) | Autenticação serverless por CPF |
| [`postech-tc3-infra-k8s`](https://github.com/Kc1t/postech-tc3-infra-k8s) | Cluster EKS + API Gateway |
| [`postech-tc3-infra-database`](https://github.com/Kc1t/postech-tc3-infra-database) | RDS PostgreSQL |

## Estrutura

```
k8s/          manifests: configmap, deployment, service, hpa
docs/adr/     decisões arquiteturais permanentes
docs/rfc/     propostas técnicas
docs/diagrams/  componentes, sequência e modelo ER
```

## Execução local

```bash
docker compose up -d
make run
```

Testes:

```bash
go test ./... -race -cover
```

## Deploy

| Evento | Ação |
|---|---|
| Pull Request | `go vet` e testes com Postgres de serviço |
| Push em `homolog` | build da imagem → ECR → `kubectl apply` no cluster de staging |
| Push em `main` | mesmo fluxo, no cluster de produção |

Secret necessário: `AWS_ROLE_ARN`. O `IMAGE_PLACEHOLDER` do `k8s/deployment.yaml` é substituído pela imagem versionada por SHA no pipeline.

O `Secret` `workshop-api` (com `DATABASE_URL` e `JWT_SECRET`) deve existir no cluster antes do primeiro deploy — veja `k8s/secret.example.yaml`.

## API

Swagger/Postman: _a definir após a migração do código._

## Pendências

- Migrar código do `postech-tc1` (com histórico).
- `Dockerfile`, `Makefile`, `docker-compose.yml` e `go.mod` vêm junto na migração.
- Instrumentação de observabilidade (Datadog ou New Relic) e logs correlacionados.
