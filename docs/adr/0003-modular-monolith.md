# ADR-0003: Monolito modular como ponto de partida

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O Tech Challenge da Fase 1 exige explicitamente um **back-end monolítico**, citando que "como será um MVP, é possível criar um Monolito utilizando a arquitetura em camadas".

Apesar do enunciado mencionar "camadas", o grupo decidiu adotar arquitetura hexagonal (ver [ADR-0001](./0001-hexagonal-architecture.md)) e organizar o domínio em **bounded contexts** (ver [ADR-0006](./0006-ddd-bounded-contexts.md)).

A pergunta é: **como estruturar o monolito** — totalmente plano (uma pasta `internal` com tudo solto) ou modular (separação por contexto/recurso)?

## Decisão

Adotar **monolito modular**: um único deployable, mas com código organizado por recurso de domínio dentro das camadas hexagonais.

```
internal/application/usecase/
  customer/        ← bounded context de Atendimento
  vehicle/
  part/            ← bounded context de Catálogo & Estoque
  service/
  service_order/   ← bounded context de Oficina
  auth/            ← bounded context de Identidade & Acesso

internal/adapters/inbound/http/
  customer/
  vehicle/
  part/
  service/
  service_order/
  auth/
```

Cada recurso tem seu próprio pacote Go com handler, use cases e (se necessário) DTOs específicos.

## Alternativas consideradas

1. **Monolito plano com tudo em um pacote.** Rápido para protótipos, mas vira "big ball of mud" rapidamente.
2. **Microserviços do início.** Conflita com o requisito do PDF e adiciona complexidade operacional desproporcional ao escopo de MVP.
3. **Serverless (lambdas isoladas).** Mesmo problema dos microserviços, com custo de cold start.
4. **Modular monolith com módulos Go separados (multi-module repository).** Considerado, mas adiciona complexidade de gerenciamento de versões internas para pouco ganho num projeto desse porte.

## Consequências

### Positivas

- Atende ao requisito explícito do PDF.
- Deploy simples (um único container), compatível com `docker-compose`.
- Operação simplificada: um log, um banco, um processo.
- Modularidade interna prepara o terreno para extração futura de serviços se necessário (cada bounded context já é um candidato natural).
- Refactors entre contextos são apenas mudanças de código (sem coordenação de versões de API).

### Negativas

- Risco de acoplamento entre contextos se a disciplina de respeitar os bounded contexts não for seguida — mitigado pela revisão de PRs e pela estrutura de pastas que torna importações cruzadas visíveis.
- Escalabilidade horizontal granular não é possível (não dá para escalar só "service-order" de forma independente). Aceitável para o escopo atual.
- Banco de dados único compartilhado entre contextos — mitigado por tabelas separadas, sem joins entre tabelas de contextos diferentes (o que seria um *anti-pattern* em monolito modular).

## Caminho de evolução

Se for necessário extrair um contexto para microserviço no futuro:

1. O contexto já tem use cases e ports independentes.
2. Substitui-se o adapter HTTP local por um cliente HTTP/gRPC para o novo serviço.
3. As tabelas do banco do contexto extraído migram para um banco dedicado.

A modularidade interna **antecipa** essa possibilidade sem pagar o custo agora.

## Referências

- Sam Newman, *Monolith to Microservices* (2019).
- Repositório do projeto: estrutura de pastas em `internal/`.
- ADRs relacionados: [0001](./0001-hexagonal-architecture.md), [0006](./0006-ddd-bounded-contexts.md).
