# ADR-0006: DDD com bounded contexts internos

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O Tech Challenge exige **modelagem com Domain-Driven Design (DDD)**, incluindo:

- Event Storming dos fluxos principais.
- Diagramas (Context Map, agregados).
- Linguagem Ubíqua aplicada.

O domínio (oficina mecânica) tem áreas funcionalmente distintas: cadastro de clientes, gestão de veículos, catálogo de serviços, controle de estoque, ordens de serviço, autenticação. Modelar tudo num "saco grande" leva ao acoplamento e dificulta evolução.

## Decisão

Modelar o sistema em **quatro bounded contexts internos**, mantendo os limites visíveis no código:

### 1. Identidade & Acesso (Identity)

- **Agregados**: `User`, `RefreshToken`.
- **Responsabilidades**: cadastro, autenticação, autorização (roles), revogação de sessões.
- **Pacotes**: `application/usecase/auth/`, `adapters/inbound/http/auth/`.

### 2. Atendimento (Customer Service)

- **Agregados**: `Customer` (solicitante), `Vehicle`.
- **Responsabilidades**: identificação por CPF/CNPJ, vínculo de múltiplos veículos por cliente.
- **Pacotes**: `application/usecase/customer/`, `application/usecase/vehicle/`.

### 3. Catálogo & Estoque (Catalog & Inventory)

- **Agregados**: `Service` (catálogo), `Part` (peça/insumo com estoque).
- **Responsabilidades**: CRUD de itens oferecidos, controle de estoque (ajustes positivos e negativos), invariante de não permitir estoque negativo.
- **Pacotes**: `application/usecase/service/`, `application/usecase/part/`.

### 4. Oficina / Execução (Workshop)

- **Agregados**: `ServiceOrder` (raiz), com `ServiceItem` e `PartItem` como entidades internas.
- **Responsabilidades**: ciclo de vida da OS (máquina de estados), cálculo automático de orçamento, métricas operacionais (tempo médio de execução).
- **Pacotes**: `application/usecase/service_order/`.

### Linguagem Ubíqua

Glossário entre vocabulário do negócio (PT) e código (EN) mantido no Miro e no [site de documentação](../site/index.html#ubiquitous-language). Termos centrais:

| Negócio                  | Código                |
| ------------------------ | --------------------- |
| Solicitante              | `Requester`/`Customer` |
| Veículo                  | `Vehicle`             |
| Ordem de Serviço (OS)    | `ServiceOrder`        |
| Serviço (mão de obra)    | `Service`             |
| Peça / Insumo            | `Part`                |
| Orçamento                | `ServiceOrder.Total`  |
| Código da OS             | `ServiceOrder.Code`   |

### Context Map

Os contextos se integram pelo padrão **Customer/Supplier**:

- **Workshop** consome IDs de **Atendimento** (cliente, veículo) e **Catálogo & Estoque** (serviços, peças).
- **Identidade** é **Conformist** para todos os outros (apenas autentica/autoriza, não conhece o domínio).

Não há *Anti-Corruption Layers* explícitas internamente porque o domínio é compartilhado e os contextos são módulos internos do mesmo monolito (ver [ADR-0003](./0003-modular-monolith.md)).

## Alternativas consideradas

1. **Modelagem plana (sem bounded contexts).** Conflita com requisitos do PDF e dificulta extração futura de microserviços.
2. **Microserviços por contexto desde o início.** Custo operacional desproporcional para um MVP.
3. **CQRS + Event Sourcing.** Excelente para domínios complexos, mas overkill para o escopo atual.

## Consequências

### Positivas

- Estrutura de pastas reflete os contextos — fácil para qualquer dev novo entender o domínio.
- Cada contexto pode evoluir independentemente.
- Preparação para extração futura de microserviços, se necessário.
- Atende ao requisito de DDD do Tech Challenge.

### Negativas

- Pode haver duplicação ocasional (ex.: campos de auditoria) entre contextos. Mitigado pelo uso de tipos compartilhados em `domain/entities/` quando faz sentido.
- Risco de "context bleeding" se a disciplina não for seguida (ex.: importar uma entidade de outro contexto diretamente). Mitigado por code review.

## Referências

- Eric Evans, *Domain-Driven Design* (2003).
- Vaughn Vernon, *Implementing Domain-Driven Design* (2013).
- Board do Miro do projeto (link no documento de entrega e no site de documentação).
- ADRs relacionados: [0001](./0001-hexagonal-architecture.md), [0003](./0003-modular-monolith.md).
