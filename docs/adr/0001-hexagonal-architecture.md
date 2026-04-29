# ADR-0001: Arquitetura hexagonal (Ports & Adapters)

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O Tech Challenge da Fase 1 da FIAP SOAT permite, explicitamente, a adoção de **arquitetura em camadas tradicional** (controller → service → repository) por se tratar de um MVP em monolito. Apesar disso, identificamos os seguintes desafios:

- O domínio (oficina mecânica) tem regras de negócio claras: máquina de estados da Ordem de Serviço, cálculo de orçamento, controle de estoque com transações, validação de CPF/CNPJ e placa.
- A escolha de framework HTTP (Gin), ORM (GORM) e biblioteca de JWT são detalhes de infraestrutura que podem mudar.
- Queremos cobertura de testes mínima de 80% nos domínios críticos — o que exige isolamento do domínio para permitir mocks de portas externas.
- O grupo prevê evolução do sistema em fases futuras (talvez microserviços, talvez troca de banco). Quanto menos a regra de negócio depender de tecnologia, mais barata é a evolução.

## Decisão

Adotar **arquitetura hexagonal (Ports & Adapters)** dentro de um monolito. O domínio fica no centro e expõe interfaces (ports). Frameworks e infraestrutura são adapters que dependem do domínio, nunca o contrário.

### Camadas

```
internal/
  domain/                  → núcleo (entidades + erros). Zero dependências externas.
  ports/                   → interfaces puras (repositories, usecases, security).
  application/usecase/     → lógica de negócio (depende só de ports + entities).
  adapters/
    inbound/http/          → adapters Gin (driving adapters).
    outbound/postgresql/   → adapters GORM (driven adapters).
    outbound/jwt/          → adapter de TokenService.
```

### Regra de dependência

Dependências apontam **sempre para dentro**. O único componente que conhece todas as camadas é o container de DI manual em `cmd/api/bootstrap/container.go`.

## Alternativas consideradas

1. **Camadas tradicionais (controller → service → repository).** Mais simples, mas o `service` acabaria importando GORM/Gin diretamente, dificultando testes e troca de tecnologia.
2. **Clean Architecture (Uncle Bob).** Muito similar ao hexagonal, porém com terminologia mais carregada (entities/use cases/interface adapters/frameworks). Optamos por hexagonal por ser mais conciso e explícito sobre "o que é porta e o que é adapter".
3. **Vertical Slice Architecture.** Boa para sistemas com features muito independentes, mas exigiria mais boilerplate e ainda é incomum em Go.

## Consequências

### Positivas

- Domínio testável sem subir banco nem servidor HTTP (usando mocks gerados via `mockery`).
- Fácil substituição de frameworks (trocar Gin por gRPC ou GORM por `sqlc` impacta só os adapters).
- Separação clara: handler é "burro", use case é regra de negócio, adapter é tradução.
- Revisores sêniores (incluindo a banca) reconhecem o padrão como sinal de maturidade arquitetural.

### Negativas

- Boilerplate maior comparado a camadas tradicionais (interfaces explícitas em `ports/`, adapters de tradução `FromDomain`/`ToDomain`).
- Curva de aprendizado para quem nunca trabalhou com hexagonal — é necessário disciplina para não vazar dependências.
- Risco de "anemia" no domínio se os use cases acabarem fazendo todo o trabalho — mitigado pela regra de comportamentos no agregado (ex.: `ServiceOrder.UpdateStatus`).

## Referências

- Alistair Cockburn, *Hexagonal Architecture* (2005).
- Vaughn Vernon, *Implementing Domain-Driven Design* (2013).
