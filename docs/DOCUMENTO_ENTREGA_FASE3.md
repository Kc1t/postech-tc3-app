# Tech Challenge — Fase 3 · Reservoir Devs

## Documento de Entrega

**Grupo 112** · Pós-Tech FIAP — SOAT (13SOAT) · Prazo: 15/09/2026

> Este documento é a base do PDF a ser submetido no Portal do Aluno. Os campos marcados com `PREENCHER` dependem do deploy e da gravação.

---

## 1. Grupo

| Nome | E-mail | GitHub |
|---|---|---|
| Kauã Miguel da Cunha | kauamigueldev@gmail.com | [@Kc1t](https://github.com/Kc1t) |
| Pedro Henrique Carlini de Oliveira Ribeiro | phcarlini@gmail.com | [@PedroHCarlini](https://github.com/PedroHCarlini) |
| Yuri Lopes Padlipskas | yupadlipskas@gmail.com | [@yuriLpadlipskas](https://github.com/yuriLpadlipskas) |
| Diego Oliveira Andrade Rafael | diegoliveiraa1@gmail.com | [@diegoliveiraa](https://github.com/diegoliveiraa) |

---

## 2. Os quatro repositórios

| # | Repositório | Responsabilidade | `soat-architecture` |
|---|---|---|---|
| 1 | [postech-tc3-lambda-auth](https://github.com/Kc1t/postech-tc3-lambda-auth) | Function serverless: emissor de token por CPF e authorizer do gateway | ✅ adicionado |
| 2 | [postech-tc3-infra-k8s](https://github.com/Kc1t/postech-tc3-infra-k8s) | Terraform do cluster EKS, metrics-server e API Gateway | ✅ adicionado |
| 3 | [postech-tc3-infra-database](https://github.com/Kc1t/postech-tc3-infra-database) | Terraform do RDS PostgreSQL gerenciado | ✅ adicionado |
| 4 | [postech-tc3-app](https://github.com/Kc1t/postech-tc3-app) | Aplicação principal executando em Kubernetes | ✅ adicionado |

Histórico das Fases 1 e 2 preservado: [postech-tc1](https://github.com/Kc1t/postech-tc1).

---

## 3. Links da entrega

| Entregável | Link |
|---|---|
| 🎬 Vídeo demonstrativo (até 15 min) | `PREENCHER` |
| 🌐 API Gateway (endpoint público) | `PREENCHER` |
| 📊 Dashboard New Relic | `PREENCHER` |
| 📖 Documentação arquitetural | [`docs/`](.) deste repositório |
| 🔧 Collection Postman | [`postman_collection.json`](../postman_collection.json) |
| 📄 Swagger | [`docs/swagger.yaml`](./swagger.yaml) |

---

## 4. Cobertura dos requisitos

### 4.1 Autenticação e API Gateway

| Requisito | Onde verificar |
|---|---|
| API Gateway implementado | Amazon API Gateway v2 — [`api_gateway.tf`](https://github.com/Kc1t/postech-tc3-infra-k8s/blob/main/api_gateway.tf) |
| Rotas sensíveis protegidas por autenticação via CPF | Rota `ANY /api/v1/{proxy+}` com authorizer Lambda; validação repetida pelo middleware `Auth` da aplicação |
| Function serverless valida o CPF | [`internal/cpf`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/cpf) — dígitos verificadores, cobertura 94% |
| Consulta existência e status do cliente | [`internal/requester`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/requester) — cobertura 100%; 404 para inexistente, 403 para inativo |
| Gera e devolve JWT | [`internal/token`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/token) — HS256 com `sub`, `role`, `email`, `name`, `document` |

### 4.2 Estrutura de repositórios e CI/CD

| Requisito | Status |
|---|---|
| Quatro repositórios separados | ✅ |
| CI/CD em cada um (GitHub Actions) | ✅ |
| Deploy automático para a nuvem | ✅ push em `homolog` → staging; push em `main` → produção |
| Branch `main` protegida, sem commit direto | `PREENCHER` — aplicar antes da entrega |
| Pull Request obrigatório para merge | `PREENCHER` — aplicar junto com a proteção |
| Deploy automático de homologação e produção | ✅ namespaces `postech-homolog` e `postech` |
| Dockerfiles | ✅ [`Dockerfile`](../Dockerfile) multi-stage |
| README com instruções claras | ✅ nos quatro |

> **Pendência conhecida:** a proteção de branch exige repositório público ou plano pago. Decisão do grupo: manter privado durante o desenvolvimento e tornar público antes da submissão, aplicando a proteção nesse momento.

### 4.3 Infraestrutura

| Requisito | Serviço |
|---|---|
| API Gateway | Amazon API Gateway v2 (HTTP API) |
| Function serverless | AWS Lambda `provided.al2023`, arm64, duas funções |
| Banco de dados gerenciado | Amazon RDS PostgreSQL 16 |
| Cluster Kubernetes com escalabilidade | Amazon EKS 1.35, node group 2–5 nós, HPA 2–10 pods por CPU (70%) e memória (80%), com o addon `metrics-server` que o HPA exige |
| Isolamento entre ambientes | `ResourceQuota` + `LimitRange` + `NetworkPolicy` por namespace ([`k8s/quota.yaml`](../k8s/quota.yaml), [`k8s/networkpolicy.yaml`](../k8s/networkpolicy.yaml)), com o addon `vpc-cni` que a política de rede exige no EKS |
| Terraform | Nos três repositórios de infraestrutura, state em S3 |

### 4.4 Monitoramento e observabilidade

| Requisito | Onde |
|---|---|
| Ferramenta (Datadog ou New Relic) | **New Relic** — justificativa em [`docs/newrelic/README.md`](./newrelic/README.md) |
| Latência das APIs | Dashboard, página APIs — p50/p95/p99 |
| CPU e memória do Kubernetes | Dashboard, página Kubernetes |
| Healthchecks e uptime | Dashboard, página Kubernetes |
| Alertas para falha no processamento de OS | Alerta 1 em [`docs/newrelic/README.md`](./newrelic/README.md) |
| Logs estruturados em JSON com correlação | `X-Correlation-ID` em toda requisição e em todo log |
| Dashboard: volume diário de OS | Evento `service_order_created` |
| Dashboard: tempo médio por status | Evento `service_order_status_changed` com `execution_seconds` |
| Dashboard: erros nas integrações | Evento `integration_failure` |

### 4.5 Documentação da arquitetura

| Requisito | Documento |
|---|---|
| Diagrama de Componentes | [`docs/diagrams/componentes.md`](./diagrams/componentes.md) |
| Diagrama de Sequência — autenticação | [`docs/diagrams/sequencia-autenticacao.md`](./diagrams/sequencia-autenticacao.md) |
| Diagrama de Sequência — abertura de OS | [`docs/diagrams/sequencia-ordem-servico.md`](./diagrams/sequencia-ordem-servico.md) |
| RFCs | [`docs/rfc/`](./rfc/) — nuvem, banco gerenciado, autenticação |
| ADRs | [`docs/adr/`](./adr/) — 0001 a 0010 |
| Justificativa do banco e ajustes no modelo, com ER | [`docs/MODELAGEM_DE_DADOS.md`](./MODELAGEM_DE_DADOS.md) e [`docs/database-model.dbml`](./database-model.dbml) |

---

## 5. Qualidade

| Métrica | Valor |
|---|---|
| Cobertura de testes da aplicação | **82,5%**, com gate de 80% travando a pipeline |
| Cobertura da Lambda | `requester` 100%, `cpf` 94%, `token` 82% |
| Lint | `golangci-lint` na pipeline |
| Vulnerabilidades | `govulncheck` na pipeline |
| Terraform | `fmt -check`, `validate` e `tfsec` na pipeline |

---

## 6. Decisões que valem destacar

1. **Duas Lambdas, não uma.** Um authorizer `REQUEST` do API Gateway recebe um evento diferente do emissor e precisa responder `isAuthorized` — não é o mesmo papel. [ADR-0007](./adr/0007-api-gateway-comunicacao.md)

2. **Cognito avaliado e descartado.** É o caminho ensinado na aula, mas não modela CPF como credencial, e o enunciado pede consulta de existência e status na base da aplicação. [RFC-0003](./rfc/0003-autenticacao-por-cpf.md)

3. **Terraform em vez de SAM.** O enunciado exige Terraform para provisionamento, e o Learner Lab não permite criar as roles IAM que o SAM cria por padrão. [ADR-0009](./adr/0009-lambda-terraform-em-vez-de-sam.md)

4. **Um cluster com dois namespaces.** Corta o custo pela metade e atende ao requisito de deploy automático por branch. As consequências negativas estão registradas. [ADR-0010](./adr/0010-cluster-unico-dois-namespaces.md)

5. **Nova coluna `requesters.status`.** O modelo herdado só respondia existência; o enunciado pede existência **e** status. [`MODELAGEM_DE_DADOS.md`](./MODELAGEM_DE_DADOS.md)

---

## 7. Checklist antes de submeter

- [ ] Tornar os quatro repositórios públicos
- [ ] Aplicar proteção de `main` com PR obrigatório nos quatro
- [ ] Confirmar `soat-architecture` como colaborador nos quatro (print)
- [ ] Preencher os links de vídeo, gateway e dashboard neste documento
- [ ] Gerar o PDF a partir deste documento e submeter no Portal do Aluno
