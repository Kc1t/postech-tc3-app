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
| [0001](./0001-hexagonal-architecture.md) | Arquitetura hexagonal (Ports & Adapters)         | Accepted |
| [0002](./0002-postgresql.md)             | PostgreSQL como banco de dados relacional        | Accepted |
| [0003](./0003-modular-monolith.md)       | Monolito modular como ponto de partida           | Accepted |
| [0004](./0004-gorm-orm.md)               | GORM como ORM                                    | Accepted |
| [0005](./0005-jwt-refresh-token.md)      | Autenticação JWT com refresh token rotativo     | Accepted |
| [0006](./0006-ddd-bounded-contexts.md)   | DDD com bounded contexts internos                | Accepted |
| [0007](./0007-domain-errors-sentinels.md)| Erros sentinela centralizados em `domainerrors` | Accepted |
