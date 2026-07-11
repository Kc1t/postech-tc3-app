# Tech Challenge — Fase 2 · Reservoir Devs
## Documento de Entrega

**Grupo:** Reservoir Devs · **Pós-Tech FIAP — SOAT (15SOAT)**

---

## 1. Grupo

| Nome | Discord | GitHub |
|---|---|---|
| Kauã Miguel da Cunha | `kc1t` | [@kc1t](https://github.com/kc1t) |
| Pedro Henrique Carlini de Oliveira Ribeiro | `phcor` | [@PedroHCarlini](https://github.com/PedroHCarlini) |
| Yuri Lopes Padlipskas | `Dyn4mo` | [@yuriLpadlipskas](https://github.com/yuriLpadlipskas) |
| Diego Oliveira Andrade Rafael | `.diegooliveira` | [@diegoliveiraa](https://github.com/diegoliveiraa) |

---

## 2. Links da entrega (itens exigidos pelo enunciado da Fase 2)

| Entregável | Link |
|---|---|
| 📂 Repositório (privado, acesso `soat-architecture`) | https://github.com/Kc1t/postech-tc1 |
| 📖 README (solução, arquitetura, deploy e instruções) | https://github.com/Kc1t/postech-tc1/blob/main/README.md |
| 🔧 Collection das APIs (Postman) | [`postman_collection.json`](https://github.com/Kc1t/postech-tc1/blob/main/postman_collection.json) |
| 📄 Swagger estático | [`docs/swagger.yaml`](https://github.com/Kc1t/postech-tc1/blob/main/docs/swagger.yaml) |
| 🏛️ Hub de documentação | https://tc1-doc.vercel.app/ |
| 🎬 Vídeo demonstrativo (até 15 min) | _adicionar link do YouTube/Vimeo_ |

---

## 3. Cobertura dos requisitos da Fase 2

### Evolução da aplicação

| Requisito | Onde verificar |
|---|---|
| Clean Code (nomes claros, simplicidade, coesão) | Código-fonte + `CLAUDE.md` (regras de estilo/arquitetura) |
| Arquitetura Hexagonal (separação de camadas e dependências) | `internal/` (domain, ports, application, adapters) + README → *Arquitetura* |
| Testes automatizados (unitários e integração) | Arquivos `_test.go` + gate de cobertura ≥ 80% na pipeline |
| **Abertura de OS** (cliente, veículo, serviços, peças → ID único) | `POST /api/v1/service-orders` |
| **Consulta de status da OS** | `GET /api/v1/service-orders/code/:code` e `/:id` |
| **Aprovação de orçamento** (notificação externa de aprovação/recusa) | `PUT /api/v1/service-orders/code/:code/status` (público) |
| **Listagem de OS** ordenada (Execução > Aguardando > Diagnóstico > Recebida; mais antigas primeiro; exclui finalizadas/entregues) | `GET /api/v1/service-orders` + `serviceorder_repository.go` (`FindAll`) |
| **Atualização de status via e-mail** | Adapter SMTP `internal/adapters/outbound/smtp/email_notifier.go` |

### Infraestrutura

| Requisito | Onde verificar |
|---|---|
| Dockerfile atualizado | [`Dockerfile`](https://github.com/Kc1t/postech-tc1/blob/main/Dockerfile) (multi-stage) |
| docker-compose para dev local | [`docker-compose.yml`](https://github.com/Kc1t/postech-tc1/blob/main/docker-compose.yml) |
| Manifestos K8s: Deployments, Services, ConfigMaps, Secrets | [`k8s/`](https://github.com/Kc1t/postech-tc1/tree/main/k8s) + `k8s/README.md` |
| Horizontal Pod Autoscaler (CPU/memória) | [`k8s/hpa.yaml`](https://github.com/Kc1t/postech-tc1/blob/main/k8s/hpa.yaml) (2–10 réplicas) |
| Terraform: cluster K8s + banco de dados | [`infra/`](https://github.com/Kc1t/postech-tc1/tree/main/infra) (EKS + RDS) + `infra/README.md` |
| Documentação dos recursos IaC e como aplicar | [`infra/README.md`](https://github.com/Kc1t/postech-tc1/blob/main/infra/README.md) |
| Pipeline CI/CD (build, testes, imagem, deploy K8s + banco) | [`.github/workflows/ci.yml`](https://github.com/Kc1t/postech-tc1/blob/main/.github/workflows/ci.yml) |

### README

| Item exigido | Status |
|---|---|
| Descrição da solução e objetivos da fase | ✅ |
| Desenho da arquitetura (componentes, infra provisionada, fluxo de deploy) | ✅ (diagramas em ASCII + `docs/documentation-diagram.drawio`) |
| Instruções de execução local | ✅ |
| Instruções de deploy em Kubernetes | ✅ |
| Instruções de provisionamento com Terraform | ✅ |
| Link para a collection das APIs | ✅ (`postman_collection.json` / Swagger) |
| Link para o vídeo demonstrativo | ⏳ adicionar |

---

## 4. Arquitetura (resumo)

- **Aplicação:** monolito Go (Gin + GORM) em Arquitetura Hexagonal — domínio isolado de HTTP, banco, JWT e SMTP.
- **Infraestrutura AWS (Terraform):** cluster **EKS** (Kubernetes) + **RDS PostgreSQL 16**. Imagens no **ECR**.
- **Kubernetes:** Deployment (2 réplicas), Service (NLB), ConfigMap, Secret e **HPA** (2–10 pods por CPU/memória).
- **Fluxo de deploy (GitHub Actions):** lint → dependências/govulncheck → testes (gate 80%) → build/push no ECR → `kubectl apply` dos manifestos no EKS → `rollout status`.

O diagrama completo está no README e em `docs/documentation-diagram.drawio`.

---

## 5. O que demonstrar no vídeo (até 15 min)

1. **Deploy da aplicação** — `terraform apply` (ou cluster já provisionado) + aplicação dos manifestos.
2. **Execução do CI/CD** — push na `main` disparando a pipeline (lint, testes, build, deploy).
3. **Consumo das APIs** — abertura de OS, consulta de status, aprovação de orçamento, listagem ordenada e e-mail de notificação.
4. **Escalabilidade automática** — simular carga e mostrar o HPA subindo réplicas (`kubectl get hpa -w` / `kubectl get pods -w`).
</content>
