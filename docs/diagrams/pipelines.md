# Pipelines de CI/CD

Quatro repositórios, quatro pipelines independentes. Todos autenticam na AWS com credencial de sessão do Learner Lab e disparam nos mesmos eventos: Pull Request valida, push entrega.

---

## 1. Modelo comum

```mermaid
flowchart LR
    pr(["Pull Request<br/>→ main ou homolog"]) --> val["Validação"]
    push(["Push<br/>em homolog ou main"]) --> val2["Validação"]

    val --> plan["Plano / dry-run<br/><b>não altera nada</b>"]
    val2 --> deploy["Entrega"]

    deploy --> amb{"qual branch?"}
    amb -->|homolog| stg["Environment: staging"]
    amb -->|main| prd["Environment: prod"]

    style plan fill:#e8f0fe,stroke:#4a7
    style prd fill:#ffe0b2,stroke:#e80
```

A separação é rígida: **nenhum job que roda em Pull Request tem permissão de escrita na nuvem**. O `plan` do Terraform e o `validate` são leitura; o `apply` só existe no gatilho de push.

---

## 2. `postech-tc3-app` — aplicação

```mermaid
flowchart TD
    trig(["PR ou push"]) --> lint["<b>lint</b><br/>golangci-lint"]
    trig --> deps["<b>dependencies</b><br/>go mod tidy verificado<br/>govulncheck"]
    trig --> test["<b>tests</b><br/>Postgres 16 como service<br/>go test -race -cover"]

    test --> gate{"cobertura ≥ 80%?"}
    gate -->|não| falha(["pipeline reprovada"])
    gate -->|sim| ok["aprovado"]

    lint --> join{{"os três passaram"}}
    deps --> join
    ok --> join

    join --> ev{"evento é push?"}
    ev -->|não| fim(["fim — PR validado"])
    ev -->|sim| dep["<b>deploy</b>"]

    dep --> creds["credenciais de sessão AWS"]
    creds --> ecr["docker build<br/>push no ECR<br/>tag = SHA + nome da branch"]
    ecr --> kube["aws eks update-kubeconfig"]
    kube --> ns["cria o namespace<br/>postech ou postech-homolog"]
    ns --> sec["cria o Secret a partir<br/>dos GitHub Secrets"]
    sec --> sed["substitui IMAGE_PLACEHOLDER<br/>e NAMESPACE_PLACEHOLDER"]
    sed --> apply["kubectl apply<br/>configmap, deployment, service, hpa"]
    apply --> roll["kubectl rollout status<br/>timeout 300s"]

    style falha fill:#ffe4e1,stroke:#c66
    style gate fill:#fff4e0,stroke:#d90
```

A cota entra **antes** do Deployment: o `LimitRange` só vale para pods criados depois dele, e aplicá-lo na ordem inversa deixaria o primeiro rollout fora do teto.

O gate de cobertura filtra mocks, `cmd/api`, `docs`, `scripts`, `pkg/env` e `pkg/ioc` antes de calcular — código gerado e ponto de entrada não inflam o número. A cobertura efetiva hoje é **82,4%**.

---

## 3. `postech-tc3-lambda-auth` — funções serverless

```mermaid
flowchart TD
    trig(["PR ou push"]) --> test["<b>test</b><br/>go mod tidy verificado<br/>gofmt · go vet<br/>go test -race -cover"]

    test --> ev{"evento"}

    ev -->|Pull Request| plan["<b>plan</b><br/>terraform fmt -check<br/>init -backend=false<br/>terraform validate"]
    ev -->|push| dep["<b>deploy</b>"]

    dep --> pack["make package<br/>build arm64 dos dois binários<br/>issuer.zip e authorizer.zip"]
    pack --> creds["credenciais de sessão AWS"]
    creds --> init["terraform init<br/>key = lambda-auth/{ambiente}.tfstate"]
    init --> apply["terraform apply<br/>TF_VAR_database_url<br/>TF_VAR_jwt_secret"]
    apply --> out["publica os invoke_arn<br/>no resumo do job"]

    style plan fill:#e8f0fe,stroke:#4a7
```

O Terraform é dono do código: o `source_code_hash` do zip dispara a atualização da função no `apply`. Não existe passo de `update-function-code` — decisão no [ADR-0009](../adr/0009-lambda-terraform-em-vez-de-sam.md).

O resumo do job imprime os quatro valores que o `infra-k8s` precisa, prontos para colar.

---

## 4. `postech-tc3-infra-k8s` e `postech-tc3-infra-database`

Os dois seguem a mesma forma.

```mermaid
flowchart TD
    trig(["PR ou push"]) --> val["<b>validate</b><br/>terraform fmt -check -recursive<br/>init -backend=false<br/>terraform validate<br/>tfsec"]

    val --> ev{"evento"}
    ev -->|Pull Request| plan["<b>plan</b> em staging<br/>mostra o diff no log"]
    ev -->|push| apply["<b>apply</b>"]

    apply --> amb{"branch"}
    amb -->|homolog| stg["-var-file=envs/staging.tfvars"]
    amb -->|main| prd["-var-file=envs/prod.tfvars"]

    stg --> tf["terraform apply -auto-approve"]
    prd --> tf

    tf --> extra{"é o repo do cluster?"}
    extra -->|não| fim(["fim"])
    extra -->|sim| kube["aws eks update-kubeconfig"]
    kube --> helm["helm upgrade --install<br/>nri-bundle do New Relic"]
    helm --> ep["publica o endpoint do gateway<br/>no resumo do job"]

    style plan fill:#e8f0fe,stroke:#4a7
    style prd fill:#ffe0b2,stroke:#e80
```

O passo do New Relic é tolerante: sem `NEW_RELIC_LICENSE_KEY` configurada, ele registra que está pulando e o `apply` segue normalmente. Isso permite subir a infraestrutura antes de existir conta de observabilidade.

---

## 5. Secrets por repositório

```mermaid
flowchart LR
    lab(["AWS Academy<br/>Learner Lab"]) -->|"AWS Details → CLI: Show"| file["`.aws-lab-credentials`<br/>não versionado"]
    file --> script["scripts/sync-aws-secrets.sh"]
    script --> r1["app"]
    script --> r2["lambda-auth"]
    script --> r3["infra-k8s"]
    script --> r4["infra-database"]

    style lab fill:#fff4e0,stroke:#d90
```

| Secret | app | lambda-auth | infra-k8s | infra-database |
|---|:-:|:-:|:-:|:-:|
| `AWS_ACCESS_KEY_ID` | ✅ | ✅ | ✅ | ✅ |
| `AWS_SECRET_ACCESS_KEY` | ✅ | ✅ | ✅ | ✅ |
| `AWS_SESSION_TOKEN` | ✅ | ✅ | ✅ | ✅ |
| `AWS_REGION` | ✅ | ✅ | ✅ | ✅ |
| `TF_STATE_BUCKET` | — | ✅ | ✅ | ✅ |
| `EKS_CLUSTER_NAME` | ✅ | — | — | — |
| `POSTGRES_DSN` | ✅ | ✅ | — | — |
| `JWT_SECRET` | ✅ | ✅ | — | — |
| `NEW_RELIC_LICENSE_KEY` | ✅ | — | ✅ | — |
| `SMTP_USERNAME` / `SMTP_PASSWORD` | ✅ | — | — | — |
| `ADMIN_PASSWORD` | ✅ | — | — | — |

O `JWT_SECRET` aparece em dois repositórios e **precisa ser idêntico**. Se divergir, o token emitido pela Lambda é recusado com 401 pelo authorizer e pela aplicação. É o erro de configuração mais provável do projeto.

**As três credenciais da AWS expiram junto com a sessão do lab, a cada ~4 horas.** Rodar o script antes de qualquer deploy é parte do fluxo, não uma exceção.

---

## 6. Proteção de branch

```mermaid
flowchart LR
    dev["commit local"] --> br["branch de trabalho"]
    br --> pr(["Pull Request"])
    pr --> ci{{"CI obrigatório"}}
    ci -->|reprovou| blk(["merge bloqueado"])
    ci -->|passou| rev{{"revisão"}}
    rev --> mrg["merge em homolog"]
    mrg --> stg["deploy em staging"]
    stg --> pr2(["PR homolog → main"])
    pr2 --> prd["deploy em produção"]

    dev -.->|"bloqueado"| main["main"]

    style blk fill:#ffe4e1,stroke:#c66
    style main fill:#ffe0b2,stroke:#e80
```

> **Pendente.** A seta pontilhada — push direto na `main` bloqueado — ainda não está aplicada. Proteção de branch exige repositório público ou plano pago, e a decisão do grupo foi manter privado durante o desenvolvimento e tornar público antes da submissão. Enquanto isso, a disciplina de PR é acordo entre as pessoas, não trava técnica.
