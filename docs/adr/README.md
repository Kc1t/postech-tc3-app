# Architecture Decision Records (ADRs)

Esta pasta contém o registro das decisões arquiteturais relevantes do projeto **Workshop API**. Cada ADR documenta o **contexto**, a **decisão** tomada e suas **consequências**.

## Por que ADRs?

ADRs (formato proposto por [Michael Nygard](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)) são arquivos curtos, versionados junto com o código, que registram decisões importantes no momento em que são tomadas. Servem como memória arquitetural — quem chega depois entende **por que** o sistema está como está, não apenas **como**.

## Convenções

- Numeração sequencial (`0001`, `0002`, …) sem buracos.
- Formato fixo: **Status**, **Contexto**, **Decisão**, **Alternativas consideradas**, **Consequências**.
- Status possíveis: `Proposed`, `Accepted`, `Deprecated`, `Superseded by ADR-XXXX`.
- Decisão registrada **uma vez** — se mudar, cria-se um novo ADR que supersede o anterior.

## Índice

| #    | Decisão                                           | Status   |
| ---- | ------------------------------------------------- | -------- |
| [0001](./0001-go-language.md)            | Go como linguagem do back-end                    | Accepted |
| [0002](./0002-hexagonal-architecture.md) | Arquitetura hexagonal (Ports & Adapters)         | Accepted |
| [0003](./0003-postgresql.md)             | PostgreSQL como banco de dados relacional        | Accepted |
| [0004](./0004-gorm-orm.md)               | GORM como ORM                                    | Accepted |
| [0005](./0005-jwt-refresh-token.md)      | Autenticação JWT com refresh token rotativo     | Accepted |
| [0006](./0006-domain-errors-sentinels.md)| Erros sentinela centralizados em `domainerrors` | Accepted |
| [0007](./0007-api-gateway-comunicacao.md) | Comunicação entre API Gateway e aplicação        | Accepted |
| [0008](./0008-hpa.md)                    | Escalabilidade de pods via HPA                   | Accepted |
| [0009](./0009-lambda-terraform-em-vez-de-sam.md) | Deploy da Lambda por Terraform, não por SAM | Accepted |
| [0010](./0010-cluster-unico-dois-namespaces.md) | Um cluster com dois namespaces           | Accepted |
| [0011](./0011-notificacoes-smtp-no-pod.md) | Notificações por SMTP no pod, com caminho para serverless | Accepted |
