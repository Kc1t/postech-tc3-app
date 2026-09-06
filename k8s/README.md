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
| `quota.yaml` | `ResourceQuota` + `LimitRange` `workshop-api` | Teto de consumo do namespace e valores padrão de requests/limits por container. |
| `networkpolicy.yaml` | `NetworkPolicy` `workshop-api` | Só a porta `8080` é alcançável de fora do namespace; o resto do tráfego de entrada é negado. |

## Como aplicar (manual)

A ordem importa: namespace → cota/rede → config/secret → deployment → service → hpa. A cota entra antes do Deployment porque o `LimitRange` só se aplica a pods criados **depois** dele. Defina `NAMESPACE` antes de começar.

```bash
NAMESPACE=postech-homolog   # ou postech
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Ajuste os valores do Secret antes de aplicar (base64)
kubectl apply -n "$NAMESPACE" -f quota.yaml
kubectl apply -n "$NAMESPACE" -f networkpolicy.yaml

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

## Cota de recursos e isolamento

Os dois ambientes dividem o mesmo cluster ([ADR-0010](../docs/adr/0010-cluster-unico-dois-namespaces.md)), então cada namespace precisa de um teto próprio — sem ele, um teste de carga em homologação consome a capacidade dos nós e deixa produção em `Pending`.

O teto é derivado do próprio HPA, não escolhido no chute. Com `maxReplicas: 10` e mais um pod de surge no rolling update, o pior caso do namespace é **11 pods**:

| | Por pod | × 11 pods | Cota |
|---|---|---|---|
| `requests.cpu` | 100m | 1100m | **2** |
| `requests.memory` | 128Mi | 1408Mi | **2Gi** |
| `limits.cpu` | 500m | 5500m | **6** |
| `limits.memory` | 512Mi | 5632Mi | **6Gi** |

A cota fica acima do pior caso de propósito. Uma cota **abaixo** do `maxReplicas` do HPA é pior que nenhuma: o autoscaler continua tentando escalar, o ReplicaSet acumula eventos `FailedCreate` e o sintoma aparece como "o HPA não funciona", não como "a cota barrou".

O `LimitRange` cobre o outro lado: qualquer container aplicado no namespace sem `resources` declarado herda 100m/128Mi de request e 500m/512Mi de limit, e nenhum container pode pedir mais de 1 CPU / 1Gi. Sem ele, um pod solto sem requests entraria na cota como zero e furaria o teto.

O mesmo arquivo vale para os dois namespaces: a carga de trabalho é idêntica, o que muda é o volume de tráfego — e esse já é absorvido pelo HPA.

## Política de rede

A `NetworkPolicy` declara duas regras de entrada:

1. de qualquer pod **do mesmo namespace**, em qualquer porta;
2. de qualquer origem, **apenas na porta 8080/TCP**.

A porta 8080 precisa ficar aberta a origens externas ao namespace porque quem bate nela não é um pod: o NLB em modo `instance` chega com o IP do **nó**, e as probes de readiness/liveness partem do **kubelet**. Restringi-las por `podSelector` derrubaria o serviço e o rollout.

O efeito real da política é negar todo o resto — nenhum pod de `postech-homolog` alcança um pod de `postech` fora da porta da API.

> **Atenção no EKS:** `NetworkPolicy` é **silenciosamente ignorada** se o addon `vpc-cni` não estiver com `enableNetworkPolicy = "true"`. O objeto é aceito pelo servidor, aparece em `kubectl get netpol` e não filtra nada. O addon é configurado em [`postech-tc3-infra-k8s/addons.tf`](https://github.com/Kc1t/postech-tc3-infra-k8s/blob/main/addons.tf). É a mesma classe de armadilha do `metrics-server` ausente com o HPA.

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
