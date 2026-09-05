# Diagrama de Componentes

Visão de nuvem da solução: APIs, banco e monitoramento.

```mermaid
flowchart TB
    cliente["Cliente da oficina<br/>(navegador / app)"]
    operacao["Operação<br/>(atendente, mecânico, admin)"]

    subgraph aws["AWS · us-east-1"]
        gw["API Gateway v2<br/>HTTP API<br/><i>postech-tc3-infra-k8s</i>"]

        subgraph serverless["Lambda · postech-tc3-lambda-auth"]
            issuer["issuer<br/>CPF → JWT"]
            authz["authorizer<br/>valida HS256"]
        end

        subgraph eks["EKS · postech-tc3-infra-k8s"]
            nlb["NLB<br/>Service LoadBalancer"]
            subgraph ns["namespace postech / postech-homolog"]
                api["workshop-api<br/>Deployment 2–10 réplicas<br/><i>postech-tc3-app</i>"]
                hpa["HPA<br/>CPU 70% · Mem 80%"]
                ms["metrics-server<br/>(addon EKS)"]
            end
            nragent["New Relic<br/>nri-bundle"]
        end

        rds[("RDS PostgreSQL 16<br/><i>postech-tc3-infra-database</i>")]
        secrets["Secrets Manager<br/>credenciais do banco"]
        ecr["ECR<br/>imagem da aplicação"]
        cw["CloudWatch Logs<br/>gateway e lambdas"]
    end

    nr["New Relic<br/>APM · logs · dashboards · alertas"]
    gha["GitHub Actions<br/>CI/CD dos 4 repositórios"]

    cliente -->|"POST /auth<br/>{cpf}"| gw
    cliente -->|"Bearer JWT"| gw
    operacao -->|"POST /api/v1/auth/login"| gw

    gw -->|"POST /auth"| issuer
    gw -.->|"autoriza ANY /api/v1/{proxy+}"| authz
    gw -->|"HTTP_PROXY"| nlb
    nlb --> api

    issuer -->|"SELECT document, status"| rds
    api --> rds
    secrets -.->|"DSN"| api
    secrets -.->|"DSN"| issuer

    ms --> hpa
    hpa -->|"escala"| api

    api -->|"stdout JSON<br/>+ correlation id"| nragent
    nragent --> nr
    gw --> cw
    issuer --> cw
    authz --> cw

    gha -->|"build + push"| ecr
    ecr --> api
    gha -->|"terraform apply"| eks
    gha -->|"terraform apply"| rds
    gha -->|"terraform apply"| serverless
```

## Camadas de autenticação

| Rota no gateway | Autorização no gateway | Validação na aplicação |
|---|---|---|
| `POST /auth` | aberta | — |
| `POST /api/v1/auth/{proxy+}` | aberta | — |
| `GET /{proxy+}` | aberta | — |
| `ANY /api/v1/{proxy+}` | authorizer Lambda | middleware `Auth` |

As duas camadas validam o mesmo JWT HS256 com o mesmo segredo. Justificativa em [ADR-0007](../adr/0007-api-gateway-comunicacao.md).

## Fronteiras de repositório

| Repositório | Provisiona |
|---|---|
| `postech-tc3-infra-database` | RDS, subnet group, security group, secret |
| `postech-tc3-infra-k8s` | EKS, node group, metrics-server, API Gateway, rotas, authorizer, agente New Relic |
| `postech-tc3-lambda-auth` | Duas Lambdas, security group, log groups |
| `postech-tc3-app` | Imagem no ECR e manifests aplicados no cluster |
