# ADR-0004: GORM como ORM

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

A Workshop API tem um modelo de dados relacional não trivial (ver [ADR-0002](./0002-postgresql.md)) e precisa de:

- Mapeamento entre entidades de domínio e tabelas do banco.
- Suporte a transações para operações atômicas (ex.: criação de OS com ajuste de estoque).
- Migrations para o esquema durante desenvolvimento.
- Consultas com filtros, joins e paginação.
- Bom desempenho e baixo overhead.

O ecossistema Go oferece várias alternativas: GORM, `database/sql` puro, `sqlx`, `sqlc`, `ent`, `bun`.

## Decisão

Adotar **GORM v2** (`gorm.io/gorm`) com driver `gorm.io/driver/postgres`.

## Alternativas consideradas

1. **`database/sql` puro.** Máximo controle e performance, mas exige escrever toda conversão manualmente — verboso e propenso a bugs em projeto com muitos CRUDs.
2. **`sqlx`.** Camada fina sobre `database/sql`. Boa opção, mas ainda exige escrever SQL manualmente para todos os CRUDs.
3. **`sqlc`.** Gera código a partir de SQL. Muito performático, mas o fluxo de "edita SQL → roda gerador → ajusta tipos" não compensa para CRUDs simples.
4. **`ent` (Facebook).** Schema-first com forte tipagem, porém com curva de aprendizado e DSL próprio.
5. **`bun`.** ORM moderno e performante, mas com comunidade menor.

## Justificativa

### Produtividade para CRUDs

Boa parte das operações da Workshop API são CRUDs (clientes, veículos, peças, serviços). GORM oferece `Create`, `First`, `Find`, `Save`, `Delete`, `Updates` prontos, reduzindo boilerplate drasticamente.

### Suporte de primeira classe a Postgres

GORM tem driver oficial para Postgres com suporte a UUID, JSONB, timestamps com timezone e migrations automáticas via `AutoMigrate`.

### Hooks para conversão domain ↔ persistence

Os adapters em `internal/adapters/outbound/postgresql/model/` implementam `FromDomain()` e `ToDomain()` para isolar o domínio do schema do banco. GORM se comporta bem com esse padrão.

### Transações expressivas

```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&serviceOrder).Error; err != nil { return err }
    if err := tx.Save(&part).Error; err != nil { return err }
    return nil
})
```

Mais legível que abrir transação manualmente com `database/sql`.

### Comunidade e documentação

GORM é o ORM Go mais usado por uma margem grande, com documentação extensa e respostas para a maioria dos problemas no Stack Overflow.

## Consequências

### Positivas

- Código de repositório enxuto: ~10–15 linhas por método CRUD.
- Migrations automáticas em desenvolvimento (`AutoMigrate`) aceleram iteração no MVP.
- Conversão domain ↔ model isolada nos arquivos de `model/`.
- Tradução de `gorm.ErrRecordNotFound` → `domainerrors.ErrNotFound` centralizada nos repositórios.

### Negativas

- **Performance**: GORM é mais lento que SQL puro para queries muito otimizadas — aceitável para o volume esperado, monitorado em `EXPLAIN ANALYZE` se necessário.
- **N+1 queries** podem ocorrer silenciosamente — mitigado por `Preload` explícito quando há joins necessários.
- **Magic** — algumas operações (callbacks, hooks, scopes) podem surpreender. Mitigado pela disciplina de manter os repositórios simples.
- **AutoMigrate em produção é arriscado** — para a próxima fase, será substituído por migrations versionadas (`goose` ou `migrate`).

## Próximos passos / dívidas técnicas

- Substituir `AutoMigrate` por ferramenta de migration versionada antes de produção.
- Adicionar `gorm.Logger` configurável por nível para diagnóstico em ambiente de homologação.
- Avaliar uso de `Prometheus` middleware para métricas de query.

## Referências

- GORM v2 — https://gorm.io/docs/