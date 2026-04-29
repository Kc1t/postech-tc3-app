# Tech Challenge — Fase 1
## Documento de Entrega

**Pós-Tech FIAP — SOAT (15SOAT)** · **Data:** _<preencher>_

---

## 1. Grupo

**Nome:** Reservoir Devs

| Nome | Discord |
|---|---|
| Kauã Miguel da Cunha | `kc1t` |
| Pedro Henrique Carlini de Oliveira Ribeiro | `phcor` |
| Yuri Lopes Padlipskas | `Dyn4mo` |
| Diego Oliveira Andrade Rafael | `.diegooliveira` |

---

## 2. Acesso à entrega

> ### 📚 Documentação técnica
> **`http://localhost:8083`** — disponível após `docker compose up`.
> Hub central com arquitetura, DDD, API, segurança e operação. A partir dele se navega para Swagger, Miro, README e demais artefatos.

> ### 📂 Repositório
> **https://github.com/Kc1t/postech-tc1** — privado, com acesso concedido ao usuário `soat-architecture`.

> ### 🎬 Vídeo de apresentação (até 15 min)
> https://www.youtube.com/watch?v=UtZBA1bVbcs

---

## 3. Análise de vulnerabilidades

Três scanners executados via Docker em 2026-04-27: **`gosec`** (SAST), **`govulncheck`** (call-graph) e **`trivy`** (deps + Dockerfile + segredos).

| Severidade | Qtd. | Origem |
|---|---|---|
| 🔴 Crítica | 1 | `pgx/v5` (CVE-2026-33816) |
| 🟠 Alta | 1 | Dockerfile sem `USER` não-root |
| 🟡 Média | 2 | `golang.org/x/crypto` |
| ⚪ Baixa | 3 | Código, Dockerfile, dependência |
| Segredos expostos | **0** | — |

> **Risco residual: BAIXO.** A análise de call-graph (`govulncheck`) confirma que **0 vulnerabilidades são exploráveis** pelos fluxos da aplicação. As CVEs em dependências não são alcançáveis pelo código.

📑 **Relatório técnico completo, com detalhamento por ferramenta e mitigações:** [`docs/security-reports/RELATORIO.md`](https://github.com/Kc1t/postech-tc1/blob/main/docs/security-reports/RELATORIO.md)
