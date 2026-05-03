# Tech Challenge — Fase 1 · Reservoir Devs
## Documento de Entrega

**Grupo:** Reservoir Devs · **Pós-Tech FIAP — SOAT (15SOAT)** · **Data:** 2026-05-03

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

## 2. Hub central da entrega

> ### 📚 **https://tc1-doc.vercel.app/**
> Site único de documentação. A partir dele se navega para **todos** os artefatos da entrega: arquitetura hexagonal, ADRs, modelo de dados, bounded contexts, linguagem ubíqua, Swagger, relatório de segurança, instruções de execução e mais.
>
> **Use este link como ponto de partida da avaliação.** Os links abaixo estão duplicados aqui apenas para acesso direto aos itens cobrados pelo enunciado.

---

## 3. Links diretos (itens exigidos pelo enunciado)

| Entregável (PDF Fase 1) | Link |
|---|---|
| 🎬 Vídeo de apresentação (até 15 min) | https://www.youtube.com/watch?v=UtZBA1bVbcs |
| 📂 Repositório (privado, acesso `soat-architecture`) | https://github.com/Kc1t/postech-tc1 |
| 🧠 Documentação DDD — Event Storming, diagramas e Linguagem Ubíqua (Miro) | https://miro.com/app/board/uXjVHcNJp-s=/?share_link_id=374317179210 |
| 📖 README com instruções de execução | https://github.com/Kc1t/postech-tc1/blob/main/README.md |
| 🔧 Swagger (após `docker compose up`) | http://localhost:8080/swagger/index.html |
| 🏛️ Documentação completa (hub) | https://tc1-doc.vercel.app/ |

---

## 4. Cobertura dos requisitos do enunciado

| Requisito | Onde verificar |
|---|---|
| Back-end monolítico em Go (Gin + GORM + PostgreSQL) | Repositório + seção *Stack* do site |
| Arquitetura em camadas / Hexagonal (Ports & Adapters) | Site → *Arquitetura Hexagonal* / *Estrutura do projeto* / *ADRs* |
| Justificativa do banco (PostgreSQL) | Site → *Banco de dados* + ADR correspondente |
| Fluxos da OS (criação, status, acompanhamento) | Site → *Ordem de Serviço* + Swagger |
| CRUDs (clientes, veículos, serviços, peças com estoque) | Swagger + Postman collection no repositório |
| Tempo médio de execução dos serviços | Site → *Endpoints de métricas* |
| JWT em rotas administrativas + validação CPF/CNPJ/placa | Site → *Segurança* |
| Testes unitários e integração (≥ 80% domínios críticos) | Site → *Testes* + `make test` |
| Dockerfile + docker-compose | Raiz do repositório |
| Event Storming + Linguagem Ubíqua | Miro (link acima) |

---

## 5. Análise de vulnerabilidades

Três scanners executados via Docker em **2026-04-27**: **`gosec`** (SAST), **`govulncheck`** (call-graph) e **`trivy`** (deps + Dockerfile + segredos).

| Severidade | Qtd. | Origem |
|---|---|---|
| 🔴 Crítica | 1 | `pgx/v5` (CVE-2026-33816) |
| 🟠 Alta | 1 | Dockerfile sem `USER` não-root |
| 🟡 Média | 2 | `golang.org/x/crypto` |
| ⚪ Baixa | 3 | Código, Dockerfile, dependência |
| Segredos expostos | **0** | — |

> **Risco residual: BAIXO.** A análise de call-graph (`govulncheck`) confirma que **0 vulnerabilidades são exploráveis** pelos fluxos da aplicação. As CVEs em dependências não são alcançáveis pelo código.

📑 **Relatório técnico completo, com detalhamento por ferramenta e mitigações:**
[`docs/security-reports/RELATORIO.md`](https://github.com/Kc1t/postech-tc1/blob/main/docs/security-reports/RELATORIO.md)
