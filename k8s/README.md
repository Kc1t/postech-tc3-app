# Manifestos Kubernetes

Manifestos para deploy da **Workshop API** em um cluster Kubernetes (AWS EKS em produção; Minikube/kind/Docker Desktop em ambiente local). Todos os recursos vivem no namespace `postech`.

## Manifestos

| Arquivo | Recurso | Descrição |
|---------|---------|-----------|
| `namespace.yaml` | `Namespace` | Cria o namespace `postech`, isolando os recursos da aplicação. |
| `configmap.yaml` | `ConfigMap` `workshop-api-config` | Variáveis **não sensíveis** (porta, ambiente, host/porta/remetente SMTP, expiração de tokens, política de lockout). |
| `secret.yaml` | `Secret` `workshop-api-secret` | Variáveis **sensíveis** (`JWT_SECRET`, `POSTGRES_DSN`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `ADMIN_PASSWORD`). Valores de exemplo — em produção são criados pela pipeline a partir dos GitHub Secrets. |
| `deployment.yaml` | `Deployment` `workshop-api` | 2 réplicas, imagem via `IMAGE_PLACEHOLDER` (substituída no deploy), `envFrom` ConfigMap + Secret, requests/limits de CPU e memória, readiness/liveness em `/health`. |
| `service.yaml` | `Service` `workshop-api` | `LoadBalancer` (NLB na AWS) expondo a porta `80` → `8080` do container. |
| `hpa.yaml` | `HorizontalPodAutoscaler` | Escala de **2 a 10 réplicas** com base em CPU (70%) e memória (80%). |

## Como aplicar (manual)

A ordem importa: namespace → config/secret → deployment → service → hpa.

```bash
kubectl apply -f namespace.yaml

# Ajuste os valores do Secret antes de aplicar (base64)
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml

# Substitua IMAGE_PLACEHOLDER pela imagem publicada (ex.: no ECR)
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f hpa.yaml
```

Acompanhar o rollout e descobrir a URL pública:

```bash
kubectl rollout status deployment/workshop-api -n postech
kubectl get svc workshop-api -n postech        # coluna EXTERNAL-IP = endpoint do LoadBalancer
```

## Pré-requisitos do HPA

O HPA precisa de métricas de CPU/memória. Garanta que o **metrics-server** esteja instalado:

```bash
# EKS geralmente não traz o metrics-server por padrão
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Minikube
minikube addons enable metrics-server
```

## Validar a escalabilidade automática

```bash
kubectl get hpa -n postech -w      # observa TARGETS e REPLICAS
kubectl get pods -n postech -w     # observa novos pods sob carga
```

Para simular carga, gere requisições contra o endpoint do Service (ex.: `hey`, `k6`, `ab` ou múltiplas ordens de serviço) e observe as réplicas subindo.

## Secret na pipeline

Em produção, o Secret **não** usa os valores de exemplo de `secret.yaml`. A pipeline de CI/CD ([`../.github/workflows/ci.yml`](../.github/workflows/ci.yml)) recria o Secret a partir dos GitHub Secrets:

```bash
kubectl create secret generic workshop-api-secret \
  --namespace=postech \
  --from-literal=JWT_SECRET="..." \
  --from-literal=POSTGRES_DSN="..." \
  --from-literal=SMTP_USERNAME="..." \
  --from-literal=SMTP_PASSWORD="..." \
  --from-literal=ADMIN_PASSWORD="..." \
  --dry-run=client -o yaml | kubectl apply -f -
```
</content>
