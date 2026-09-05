# Manifestos Kubernetes

Manifestos para deploy da **Workshop API** em um cluster Kubernetes (AWS EKS em produção; Minikube/kind/Docker Desktop em ambiente local). Os manifestos **não declaram namespace**: ele é escolhido no momento da aplicação. O pipeline usa `postech` para a branch `main` e `postech-homolog` para a `homolog`, no mesmo cluster — decisão registrada no [ADR-0010](../docs/adr/0010-cluster-unico-dois-namespaces.md).

## Manifestos

| Arquivo | Recurso | Descrição |
|---------|---------|-----------|
| `configmap.yaml` | `ConfigMap` `workshop-api-config` | Variáveis **não sensíveis** (porta, ambiente, host/porta/remetente SMTP, expiração de tokens, política de lockout). |
| `secret.yaml` | `Secret` `workshop-api-secret` | Variáveis **sensíveis** (`JWT_SECRET`, `POSTGRES_DSN`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `ADMIN_PASSWORD`). Valores de exemplo — em produção são criados pela pipeline a partir dos GitHub Secrets. |
| `deployment.yaml` | `Deployment` `workshop-api` | 2 réplicas, imagem via `IMAGE_PLACEHOLDER` (substituída no deploy), `envFrom` ConfigMap + Secret, requests/limits de CPU e memória, readiness/liveness em `/health`. |
| `service.yaml` | `Service` `workshop-api` | `LoadBalancer` (NLB na AWS) expondo a porta `80` → `8080` do container. |
| `hpa.yaml` | `HorizontalPodAutoscaler` | Escala de **2 a 10 réplicas** com base em CPU (70%) e memória (80%). |

## Como aplicar (manual)

A ordem importa: namespace → config/secret → deployment → service → hpa. Defina `NAMESPACE` antes de começar.

```bash
NAMESPACE=postech-homolog   # ou postech
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Ajuste os valores do Secret antes de aplicar (base64)
kubectl apply -n "$NAMESPACE" -f configmap.yaml
kubectl apply -n "$NAMESPACE" -f secret.yaml

# Substitua IMAGE_PLACEHOLDER pela imagem publicada (ex.: no ECR)
kubectl apply -n "$NAMESPACE" -f deployment.yaml
kubectl apply -n "$NAMESPACE" -f service.yaml
kubectl apply -n "$NAMESPACE" -f hpa.yaml
```

Acompanhar o rollout e descobrir a URL pública:

```bash
kubectl rollout status deployment/workshop-api -n "$NAMESPACE"
kubectl get svc workshop-api -n "$NAMESPACE"        # coluna EXTERNAL-IP = endpoint do LoadBalancer
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
kubectl get hpa -n "$NAMESPACE" -w      # observa TARGETS e REPLICAS
kubectl get pods -n "$NAMESPACE" -w     # observa novos pods sob carga
```

Para simular carga, gere requisições contra o endpoint do Service (ex.: `hey`, `k6`, `ab` ou múltiplas ordens de serviço) e observe as réplicas subindo.

## Secret na pipeline

Em produção, o Secret **não** usa os valores de exemplo de `secret.yaml`. A pipeline de CI/CD ([`../.github/workflows/ci.yml`](../.github/workflows/ci.yml)) recria o Secret a partir dos GitHub Secrets:

```bash
kubectl create secret generic workshop-api-secret \
  --namespace="$NAMESPACE" \
  --from-literal=JWT_SECRET="..." \
  --from-literal=POSTGRES_DSN="..." \
  --from-literal=SMTP_USERNAME="..." \
  --from-literal=SMTP_PASSWORD="..." \
  --from-literal=ADMIN_PASSWORD="..." \
  --dry-run=client -o yaml | kubectl apply -f -
```
</content>
