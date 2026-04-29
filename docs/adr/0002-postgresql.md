# ADR-0002: PostgreSQL como banco de dados relacional

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O domínio da oficina mecânica é fortemente relacional:

- **Hierarquia clara**: `Requester (1) → Vehicle (N) → ServiceOrder (N) → ServiceItem / PartItem (N)`.
- **Transações críticas**: criação de OS atualiza estoque, calcula total e cria itens — tudo precisa ser atômico.
- **Buscas frequentes por chave de negócio**: CPF/CNPJ, placa, código da OS.
- **Histórico imutável**: ordens finalizadas não podem ser silenciosamente alteradas.

O Tech Challenge permite escolha livre de banco, **desde que justificada**.

## Decisão

Adotar **PostgreSQL 16** como banco de dados único, executado em container dedicado via `docker-compose`.

## Alternativas consideradas

1. **MySQL 8.** Maduro e comum, porém o suporte a tipos avançados (JSONB, arrays) e a *check constraints* expressivas do Postgres é superior.
2. **MongoDB.** Atrativo para desenvolvimento rápido, mas o domínio é fortemente relacional — usar Mongo exigiria simular FKs na aplicação, perdendo integridade nativa.
3. **SQLite.** Ótimo para protótipos, mas não suporta o uso concorrente esperado em produção e não tem `gen_random_uuid()` nativo.
4. **DynamoDB / outros NoSQL.** Lock-in de cloud e modelagem mais complexa para joins entre OS, peças e serviços.

## Decisão detalhada — por que PostgreSQL

### ACID e integridade transacional

Uma Ordem de Serviço agrega múltiplas entidades (serviços, peças) e atualiza o estoque numa única operação. Se qualquer parte falhar, a transação inteira é revertida. Postgres garante isso nativamente sem necessidade de compensações no código de aplicação.

### Integridade referencial

Chaves estrangeiras com `ON DELETE RESTRICT` impedem:

- Remoção de solicitantes que ainda têm veículos cadastrados.
- Remoção de peças referenciadas em OS abertas.
- OS órfãs sem veículo válido.

Isso elimina classes inteiras de bugs sem código adicional.

### UUID nativo

`gen_random_uuid()` (extensão `pgcrypto`, embutida) elimina dependência externa para gerar IDs. UUIDs são preferíveis a `int auto_increment` para evitar enumeração de recursos via API.

### JSONB para Value Objects

Itens de serviço e peças vinculados a uma OS são armazenados como `JSONB`, permitindo:

- Snapshot do preço no momento da OS (não muda se o catálogo for atualizado depois).
- Queries com índices GIN se necessário (`SELECT ... WHERE items @> '[{"part_id": "..."}]'`).

### Buscas indexadas em campos de negócio

Índices B-tree em `requesters.document` e `vehicles.plate` garantem busca em O(log n) independente do volume.

### Maturidade e ecossistema

- Driver oficial Go (`pgx`).
- ORM maduro (GORM tem suporte de primeira classe a Postgres).
- pgAdmin (incluso no `docker-compose`) para inspeção visual sem ferramentas externas.
- Documentação extensa, comunidade gigante.

### Escalabilidade

Para o volume esperado de uma oficina (centenas/milhares de OSs por mês), single-node basta. Se crescer:

- **Read replicas** nativas via streaming replication.
- **Particionamento de tabela** para histórico antigo de OSs.
- **Connection pooling** via `pgbouncer`.

Sem necessidade de migração de tecnologia.

## Consequências

### Positivas

- Garantias ACID sem código extra na aplicação.
- Modelo de dados expresso fielmente no schema (FKs, constraints, types).
- Ferramental rico (pgAdmin, `psql`, `pg_dump`, etc.).
- Driver e ORM Go maduros.

### Negativas

- Necessidade de gerenciar migrations (mitigado por `gorm.AutoMigrate` para o MVP; produção exigiria ferramenta dedicada como `goose` ou `migrate`).
- Maior consumo de memória que SQLite num cenário de "tudo em um container".
- Lock-in moderado em queries específicas do Postgres (JSONB, `ON CONFLICT`) — aceitável dado o ganho de produtividade.

## Referências

- PostgreSQL 16 release notes — https://www.postgresql.org/docs/16/release-16.html
- Repositório do projeto: `internal/adapters/outbound/postgresql/`.
- `docker-compose.yml` na raiz.
