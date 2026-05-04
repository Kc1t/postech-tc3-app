# Roteiro — Kauã

**Bloco:** Documentação, ADRs e C4
**Tempo alvo:** ~3 a 4 min
**Objetivo:** mostrar **onde** as decisões e estrutura estão documentadas. Sem aprofundar em fluxo, DDD ou API — isso é com Pedro e Diego.

---

## 1. Abertura (~15s)

> "Agora eu mostro **onde** a gente documentou tudo: arquitetura, decisões e diagramas. Quem quiser ir fundo em qualquer ponto, tudo tá linkado a partir do site."

Abrir: **https://tc1-doc.vercel.app/**

---

## 2. O hub de documentação (~30s)

- Mostrar o menu lateral do site rolando rápido:
  - Visão geral · Arquitetura · DDD · API · Segurança · Operação
- Falar:
  > "Esse site é o hub. Tem Swagger, Miro, README, ADRs e relatório de segurança — tudo num lugar só, hospedado no Vercel pra avaliador acessar sem precisar subir nada."

---

## 3. Diagramas C4 (~1 min)

Abrir o arquivo `documentation-diagram.drawio` (ou a página *Arquitetura* do site).

Passar pelas 3 páginas:

| Nível | O que dizer |
|---|---|
| **C1 — System Context** | "Dois atores: **Mecânico** (admin) e **Solicitante** (cliente). Ambos falam com a Workshop API, que persiste em PostgreSQL." |
| **C2 — Container** | "Monolito Go/Gin na 8080, PostgreSQL na 5432, PgAdmin na 8082 e o site de docs na 8083. Tudo orquestrado pelo `docker-compose`." |
| **C3 — Component** | "Aqui dá pra ver a Hexagonal: **Inbound (handlers HTTP)** → **Application (use cases)** → **Outbound (repositórios GORM)** → **Domain (entidades + value objects)**. Cada recurso (Auth, Requester, Vehicle, Service, Part, ServiceOrder) atravessa essas 4 camadas." |

> "Não vou abrir cada componente — quem quiser ver os métodos de cada use case e repositório, está tudo no diagrama."

---

## 4. Arquitetura Hexagonal (~30s)

Voltar pro site, seção *Arquitetura Hexagonal* / *Estrutura do projeto*.

- Apontar a estrutura de pastas: `domain/entities`, `application/usecase`, `ports`, `adapters/inbound`, `adapters/outbound`.
- Frase-chave:
  > "**Use case nunca importa GORM, JWT ou bcrypt** — só interfaces de `ports/`. Isso é o que mantém o domínio isolado."

---

## 5. ADRs — Decisões registradas (~1 min)

Ir em *Decisões (ADRs)* e citar **3 principais** (uma frase cada):

- **PostgreSQL** — relacional forte, suporte nativo a UUID, justificativa pedida pelo enunciado.
- **Hexagonal sobre camadas tradicionais** — testabilidade e troca de adapters sem tocar no domínio.
- **JWT + refresh token rotativo** — segurança em rotas administrativas, refresh com rotação no banco.

> "Cada ADR tem contexto, alternativas e consequências. Qualquer decisão futura passa por uma nova ADR."

---

## 6. Fechamento (~15s)

> "Documentação é viva e linkada a partir do site. Agora o **Diego** mostra a API rodando."

---

## Cheat sheet — abas pré-abertas

1. https://tc1-doc.vercel.app/ (página inicial)
2. `documentation-diagram.drawio` aberto na página C1
3. Site na seção *Arquitetura Hexagonal*
4. Site na seção *Decisões (ADRs)*

## Se estourar o tempo, corte nesta ordem

1. C1 (pula direto pra C2/C3)
2. Cheat das pastas (já tá no diagrama C3)
3. ADRs vira "tem 3 ADRs principais documentadas no site"
