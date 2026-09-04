# ADR-0001: Go como linguagem do back-end

- **Status:** Accepted
- **Data:** 2026-04
- **Decisão por:** Reservoir Devs

## Contexto

O Tech Challenge pede um back-end monolítico para o MVP do sistema da oficina, com foco em qualidade, segurança e DDD. A escolha de linguagem é livre, mas impacta diretamente:

- Tempo de boot e footprint do container (relevante para `docker-compose` em ambiente local de avaliação).
- Maturidade do ORM e do driver de Postgres.
- Curva de adoção pelo time (todos os integrantes têm experiência em ambas as linguagens consideradas).
- Distribuição: binário único vs. JVM + artefato.

A discussão se restringiu a duas opções viáveis dentro do que o time domina: **Go** e **Java**.

## Decisão

Adotar **Go 1.25** como linguagem do back-end, com **Gin** como framework HTTP e **GORM** como ORM.

## Alternativas consideradas

1. **Java + Spring Boot.** Stack extremamente madura, ecossistema gigantesco, ORM consolidado (Hibernate/JPA). Familiar à maior parte do mercado. Contra: boot mais lento, footprint de memória maior no container, build/empacotamento mais pesado e verbosidade que não agrega num MVP enxuto.

## Decisão detalhada — por que Go

### Footprint e startup

Binário único compilado estaticamente, container final (`scratch`/`distroless`) com poucos MB e startup em centenas de milissegundos. Em ambiente de avaliação (`docker compose up`) isso significa o avaliador ter o ambiente de pé em segundos, não em ~30s+ esperando a JVM aquecer.

### Concorrência nativa

Goroutines + channels são parte da linguagem, sem dependência de pool/executor configurado por fora. Para um monolito com requisições HTTP concorrentes acessando o mesmo Postgres, isso reduz boilerplate e classes inteiras de bug de threading.

### Tipagem forte e simplicidade

Sistema de tipos suficiente para modelar Value Objects (CPF/CNPJ, Placa) e estados da OS sem framework por trás. A linguagem é deliberadamente pequena — menos features = menos jeitos de fazer a mesma coisa, código mais homogêneo entre os 4 integrantes.

### Hexagonal natural

Interfaces implícitas em Go (`structural typing`) casam muito bem com Ports & Adapters: o domínio define a interface, qualquer struct que satisfaça os métodos vira adapter sem `implements` explícito. Reduz acoplamento sem framework de DI.

### Ecossistema suficiente

- `gin-gonic/gin` para HTTP.
- `gorm.io/gorm` + driver `pgx` para Postgres.
- `golang-jwt/jwt` para tokens.
- `swaggo/swag` para gerar Swagger a partir de annotations.
- `gosec`, `govulncheck`, `trivy` para análise de segurança.

### Familiaridade do time

Todos os integrantes já trabalharam com Go em outros contextos. A escolha não introduz curva de aprendizado dentro do prazo do TC1.

## Consequências

### Positivas

- Container final pequeno, deploy e CI rápidos.
- Concorrência primeira-classe sem configuração extra.
- Hexagonal expresso de forma idiomática via interfaces.
- Ferramental de segurança maduro (`gosec`, `govulncheck`).

### Negativas

- Ecossistema de bibliotecas menor que o de Java/Spring para casos muito específicos (ex.: integrações enterprise). Não impacta o escopo do MVP.
- Tratamento explícito de erros (`if err != nil`) gera verbosidade — mitigado por convenção de guard clauses e erros de domínio centralizados (ver [ADR-0006](./0006-domain-errors-sentinels.md)).
- Falta de generics em código mais antigo do ecossistema; a partir do Go 1.18 esse ponto perdeu relevância.

## Referências

- Go release notes — https://go.dev/doc/devel/release
- Effective Go — https://go.dev/doc/effective_go
