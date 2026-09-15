# Tech Challenge — Fase 3 · Reservoir Devs

## Documento de Entrega

**Grupo 112** · Pós-Tech FIAP — SOAT (13SOAT) · Prazo: 15/09/2026

> Este documento é a base do PDF submetido no Portal do Aluno.

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
| 🎬 Vídeo demonstrativo (até 15 min) | https://youtu.be/pxW0rp4isL0 |
| 🌐 API Gateway (endpoint público) | https://tkh5cum8g8.execute-api.us-east-1.amazonaws.com — Swagger em `/swagger/index.html` |
| 📊 Dashboard New Relic (snapshots, válidos por 30 dias) | [Negócio](https://web-snapshots.newrelic.com/snapshot/1789177336_1791769336_894d3355-b437-4cc0-ab1b-9dadf67d483a.pdf?token=70a5d90cd8f031d9b0968fb33995088e50a7e9ac94500e38d5c6af7ec21e9456) · [APIs](https://web-snapshots.newrelic.com/snapshot/1789177340_1791769340_7463b348-e349-4420-b056-1578aa47a2a7.pdf?token=f9c373e10a688bd424a354488057a241fb9cfc29a7977729b4e65c0464666ace) · [Kubernetes](https://web-snapshots.newrelic.com/snapshot/1789177344_1791769344_31fc75d4-6e34-415e-9ed6-6db4fffa6bba.pdf?token=c6c8db3cca86059c97904bfcbd7f25603d325b7fe7fbf4b92f79885af8a689bb) — o painel ao vivo aparece no vídeo |
| 📖 Documentação arquitetural | [`docs/`](https://github.com/Kc1t/postech-tc3-app/tree/main/docs) do repositório da aplicação |
| 🔧 Collection Postman | [`postman_collection.json`](https://github.com/Kc1t/postech-tc3-app/blob/main/postman_collection.json) |
| 📄 Swagger | ao vivo em https://tkh5cum8g8.execute-api.us-east-1.amazonaws.com/swagger/index.html · estático em [`docs/swagger.yaml`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/swagger.yaml) |

---

## 4. Cobertura dos requisitos

### 4.1 Autenticação e API Gateway

| Requisito | Onde verificar |
|---|---|
| API Gateway implementado | Amazon API Gateway v2 — [`api_gateway.tf`](https://github.com/Kc1t/postech-tc3-infra-k8s/blob/main/api_gateway.tf) |
| Rotas sensíveis protegidas por autenticação via CPF | Rota `ANY /api/v1/{proxy+}` com authorizer Lambda; o middleware `Auth` da aplicação valida o token de novo e confere que o CPF do token é o dono da OS (`403` caso contrário) |
| Function serverless valida o CPF | [`internal/cpf`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/cpf) — dígitos verificadores, cobertura 94% |
| Consulta existência e status do cliente | [`internal/requester`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/requester) — cobertura 100%; 404 para inexistente, 403 para inativo |
| Gera e devolve JWT | [`internal/token`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/internal/token) — HS256 com `sub`, `role`, `email`, `name`, `document` |

### 4.2 Estrutura de repositórios e CI/CD

| Requisito | Status |
|---|---|
| Quatro repositórios separados | ✅ |
| CI/CD em cada um (GitHub Actions) | ✅ |
| Deploy automático para a nuvem | ✅ push em `homolog` → staging; push em `main` → produção |
| Branch `main` protegida, sem commit direto | ✅ nos quatro repositórios, que são públicos |
| Pull Request obrigatório para merge | ✅ merge na `main` só por PR, com checks obrigatórios — app: `Lint`, `Dependencies`, `Tests & Coverage`; lambda: `Test`, `Terraform Plan`; infra: `Validate` |
| Deploy automático de homologação e produção | ✅ app: namespaces `postech-homolog` e `postech` · lambda: functions `postech-tc3-staging-auth-*` e `postech-tc3-prod-auth-*` · infra: `plan` automático na `homolog` e `apply` na `main`, porque cluster e banco são compartilhados ([ADR-0010](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0010-cluster-unico-dois-namespaces.md)) |
| Dockerfiles | ✅ [`Dockerfile`](https://github.com/Kc1t/postech-tc3-app/blob/main/Dockerfile) multi-stage no app; os outros três não geram imagem (as Lambdas sobem como zip, a infra é só Terraform) |
| README com instruções claras | ✅ nos quatro |

### 4.3 Infraestrutura

| Requisito | Serviço |
|---|---|
| API Gateway | Amazon API Gateway v2 (HTTP API) |
| Function serverless | AWS Lambda `provided.al2023`, arm64, duas funções |
| Banco de dados gerenciado | Amazon RDS PostgreSQL 16 |
| Cluster Kubernetes com escalabilidade | Amazon EKS 1.35, node group de 2 nós `t3.medium` (máximo 5), HPA 2–10 pods por CPU (70%) e memória (80%), com o addon `metrics-server` que o HPA exige |
| Isolamento entre ambientes | `ResourceQuota` + `LimitRange` + `NetworkPolicy` por namespace ([`k8s/quota.yaml`](https://github.com/Kc1t/postech-tc3-app/blob/main/k8s/quota.yaml), [`k8s/networkpolicy.yaml`](https://github.com/Kc1t/postech-tc3-app/blob/main/k8s/networkpolicy.yaml)), com o addon `vpc-cni` que a política de rede exige no EKS |
| Terraform | Nos três repositórios de infraestrutura, state em S3 versionado; os recursos de produção estão no state e o pipeline aplica no push da `main` |

### 4.4 Monitoramento e observabilidade

| Requisito | Onde |
|---|---|
| Ferramenta (Datadog ou New Relic) | **New Relic** — justificativa em [`docs/newrelic/README.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/newrelic/README.md) |
| Latência das APIs | Dashboard, página APIs — p50/p95/p99 |
| CPU e memória do Kubernetes | Dashboard, página Kubernetes |
| Healthchecks e uptime | Dashboard, página Kubernetes |
| Alertas para falha no processamento de OS | Alerta 1 em [`docs/newrelic/README.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/newrelic/README.md) |
| Logs estruturados em JSON com correlação | `X-Correlation-ID` em toda requisição e em todo log |
| Dashboard: volume diário de OS | Evento `service_order_created` |
| Dashboard: tempo médio por status | Evento `service_order_status_changed` com `from_status` e `seconds_in_status` (Diagnóstico, Execução, Finalização) |
| Dashboard: erros nas integrações | Evento `integration_failure` |

### 4.5 Documentação da arquitetura

| Requisito | Documento |
|---|---|
| Diagrama de Componentes | [`docs/diagrams/componentes.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/diagrams/componentes.md) |
| Diagrama de Sequência — autenticação | [`docs/diagrams/sequencia-autenticacao.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/diagrams/sequencia-autenticacao.md) |
| Diagrama de Sequência — abertura de OS | [`docs/diagrams/sequencia-ordem-servico.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/diagrams/sequencia-ordem-servico.md) |
| RFCs | [`docs/rfc/`](https://github.com/Kc1t/postech-tc3-app/tree/main/docs/rfc) — nuvem, banco gerenciado, autenticação |
| ADRs | [`docs/adr/`](https://github.com/Kc1t/postech-tc3-app/tree/main/docs/adr) — 0001 a 0011 |
| Justificativa do banco e ajustes no modelo, com ER | [`docs/MODELAGEM_DE_DADOS.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/MODELAGEM_DE_DADOS.md) e [`docs/database-model.dbml`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/database-model.dbml) |

---

## 5. Qualidade

| Métrica | Valor |
|---|---|
| Cobertura de testes da aplicação | **82,4%**, com gate de 80% travando a pipeline |
| Cobertura da Lambda | `requester` 100%, `cpf` 94%, `token` 82% |
| Lint | `golangci-lint` na pipeline |
| Vulnerabilidades | `govulncheck` na pipeline dos **dois** repositórios Go, sem achados alcançáveis |
| Terraform | `fmt -check`, `validate` e `tfsec` na pipeline, **zero achados** nos dois repositórios de infra |
| Alertas | definidos como código em [`newrelic/alerts.json`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/newrelic/alerts.json), aplicados por script |

---

## 6. Decisões que valem destacar

1. **Duas Lambdas, não uma.** Um authorizer `REQUEST` do API Gateway recebe um evento diferente do emissor e precisa responder `isAuthorized` — não é o mesmo papel. [ADR-0007](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0007-api-gateway-comunicacao.md)

2. **Cognito avaliado e descartado.** É o caminho ensinado na aula, mas não modela CPF como credencial, e o enunciado pede consulta de existência e status na base da aplicação. [RFC-0003](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/rfc/0003-autenticacao-por-cpf.md)

3. **Terraform em vez de SAM.** O enunciado exige Terraform para provisionamento, e o Learner Lab não permite criar as roles IAM que o SAM cria por padrão. [ADR-0009](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0009-lambda-terraform-em-vez-de-sam.md)

4. **Um cluster com dois namespaces.** Corta o custo pela metade e atende ao requisito de deploy automático por branch. A Lambda, que cobra por invocação, tem um par de staging próprio; cluster e banco, que cobram por hora, são compartilhados. As consequências negativas estão registradas. [ADR-0010](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0010-cluster-unico-dois-namespaces.md)

5. **Nova coluna `requesters.status`.** O modelo herdado só respondia existência; o enunciado pede existência **e** status. [`MODELAGEM_DE_DADOS.md`](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/MODELAGEM_DE_DADOS.md)

6. **Notificações ainda saem do pod.** O e-mail de mudança de status usa SMTP de dentro da aplicação, com falha observável no dashboard e em alerta; a migração para SNS → Lambda → SES está desenhada, não implementada. [ADR-0011](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0011-notificacoes-smtp-no-pod.md)

---

## 7. Limitações conhecidas

Todas vêm do orçamento e das restrições do Learner Lab ou do prazo, e estão registradas para quem herdar o projeto.

| Limitação | Por quê | Caminho |
|---|---|---|
| Os nós do EKS não escalam sozinhos | Não há Cluster Autoscaler; o node group fica em 2 nós (máximo 5) e o HPA escala só os pods | Cluster Autoscaler ou Karpenter |
| Homologação e produção dividem banco e `JWT_SECRET` | Secrets no nível do repositório e um RDS só ([ADR-0010](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0010-cluster-unico-dois-namespaces.md)) | Secrets por GitHub Environment e um database lógico por ambiente |
| As Lambdas de staging não têm rota no gateway | O gateway é único e aponta para as de produção | Um stage ou gateway de homologação |
| Notificações por SMTP dentro do pod | [ADR-0011](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0011-notificacoes-smtp-no-pod.md) | SNS → Lambda → SES |
| As Lambdas não reportam ao New Relic | Os logs delas ficam no CloudWatch; o widget e o alerta de "Autenticação por CPF" ficam sem dados | Layer/extension do New Relic nas Lambdas |
| `NetworkPolicy` com a porta `8080` aberta | É por onde chegam o NLB e as probes do kubelet | VPC Link com NLB interno ([ADR-0007](https://github.com/Kc1t/postech-tc3-app/blob/main/docs/adr/0007-api-gateway-comunicacao.md)) |
| State do Terraform sem lock | Backend S3 sem DynamoDB nem `use_lockfile` | `use_lockfile = true` no backend |
| A imagem da API roda como root | O `Dockerfile` não declara `USER` | Usuário não-root na imagem de runtime |

---

## 8. Checklist antes de submeter

- [x] Tornar os quatro repositórios públicos
- [x] Aplicar proteção de `main` com PR obrigatório nos quatro
- [x] Confirmar `soat-architecture` como colaborador nos quatro (verificado pela API do GitHub em 15/09/2026)
- [x] Preencher os links de vídeo, gateway e dashboard neste documento
- [ ] Gerar o PDF a partir deste documento e submeter no Portal do Aluno
