# Roteamento

Toda requisição atravessa duas malhas de decisão: primeiro o **API Gateway**, que escolhe o destino e decide se chama o authorizer; depois o **roteador da aplicação**, que escolhe o grupo de middlewares. Este documento desenha as duas.

---

## 1. Visão geral — o caminho de uma requisição

```mermaid
flowchart LR
    req(["Requisição"]) --> gw{{"API Gateway<br/>casamento de rota"}}

    gw -->|"POST /auth"| iss["Lambda issuer"]
    gw -->|"aberta"| app1["Aplicação"]
    gw -->|"protegida"| az{{"Lambda authorizer"}}

    az -->|"isAuthorized: false"| r401(["401"])
    az -->|"isAuthorized: true"| app2["Aplicação"]

    iss --> db[("RDS")]
    app1 --> mw{{"middlewares<br/>da aplicação"}}
    app2 --> mw
    mw --> h["handler"]
    h --> db

    style r401 fill:#ffe4e1,stroke:#c66
    style az fill:#fff4e0,stroke:#d90
    style gw fill:#e8f0fe,stroke:#4a7
```

O ponto que costuma passar despercebido: **passar pelo authorizer não dispensa os middlewares da aplicação**. São duas decisões independentes sobre o mesmo token, e a segunda ainda checa papel e posse do documento. Justificativa em [ADR-0007](../adr/0007-api-gateway-comunicacao.md).

---

## 2. Roteamento no API Gateway

O API Gateway v2 resolve por **especificidade**: a rota mais específica vence, e `$default` é o último recurso. A ordem abaixo é a de avaliação efetiva.

```mermaid
flowchart TD
    in(["Requisição no gateway"]) --> m1{"POST /auth ?"}

    m1 -->|sim| issuer["AWS_PROXY → Lambda issuer<br/><b>sem authorizer</b>"]
    m1 -->|não| m2{"POST /api/v1/auth/{proxy+} ?"}

    m2 -->|sim| pub["HTTP_PROXY → NLB<br/><b>sem authorizer</b><br/>login, registro, refresh"]
    m2 -->|não| m3{"começa com /api/v1/ ?"}

    m3 -->|sim| prot["ANY /api/v1/{proxy+}<br/><b>authorizer CUSTOM</b>"]
    m3 -->|não| m4{"método é GET ?"}

    m4 -->|sim| open["GET /{proxy+} → NLB<br/><b>sem authorizer</b><br/>/health, /swagger"]
    m4 -->|não| nf(["404 — rota não declarada"])

    prot --> az{{"Lambda authorizer<br/>TTL de cache = 0"}}
    az -->|"token ausente, malformado,<br/>expirado, assinatura inválida,<br/>sem sub ou sem role"| deny(["401"])
    az -->|"válido"| fwd["HTTP_PROXY → NLB<br/>+ contexto: subject, role, document"]

    style deny fill:#ffe4e1,stroke:#c66
    style nf fill:#ffe4e1,stroke:#c66
    style az fill:#fff4e0,stroke:#d90
    style issuer fill:#e6f4ea,stroke:#4a7
```

### Tabela de rotas do gateway

| Ordem | Rota | Integração | Authorizer | Por quê |
|---|---|---|---|---|
| 1 | `POST /auth` | `AWS_PROXY` → issuer | não | É onde o CPF vira token; exigir token aqui seria circular |
| 2 | `POST /api/v1/auth/{proxy+}` | `HTTP_PROXY` → NLB | não | Login da operação por e-mail e senha — outro fluxo de identidade |
| 3 | `ANY /api/v1/{proxy+}` | `HTTP_PROXY` → NLB | **sim** | Todo o resto da API |
| 4 | `GET /{proxy+}` | `HTTP_PROXY` → NLB | não | `/health` para as probes e o Swagger |

**Por que a 2 precisa existir:** sem ela, `POST /api/v1/auth/login` cairia na regra 3 e exigiria um token para obter um token. A rota mais específica vence, e é isso que abre a exceção.

**Por que a 4 é só `GET`:** limita a superfície aberta. Um `POST` para um caminho não declarado devolve 404 no gateway, sem chegar na aplicação.

---

## 3. Roteamento na aplicação

O roteador do Gin monta quatro grupos, cada um com uma pilha de middlewares diferente.

```mermaid
flowchart TD
    entra(["Requisição no pod"]) --> glob["<b>Middlewares globais</b><br/>Recovery → New Relic → CorrelationID → RequestLogger"]

    glob --> raiz{"caminho"}

    raiz -->|"/health"| health["200 — sem log,<br/>excluído para não poluir<br/>com probe do kubelet"]
    raiz -->|"/swagger/*"| swagger["UI do Swagger"]
    raiz -->|"/api/v1/..."| g{"grupo"}

    g -->|"POST /auth/register<br/>POST /auth/login<br/>POST /auth/refresh"| aberto["<b>Aberto</b><br/>sem middleware adicional"]

    g -->|"POST /auth/logout"| prot["<b>protected</b><br/>+ Auth"]

    g -->|"GET /service-orders/code/:code<br/>GET /service-orders/requester<br/>PUT /service-orders/code/:code/status"| cli["<b>protected</b><br/>+ Auth<br/>+ posse do documento"]

    g -->|"todo o restante"| adm["<b>admin</b><br/>+ Auth<br/>+ RequireRole(admin)"]

    prot --> auth{{"Auth<br/>valida HS256, exige sub e role"}}
    cli --> auth
    adm --> auth

    auth -->|inválido| e401(["401"])
    auth -->|válido| papel{"grupo exige admin?"}

    papel -->|sim| role{{"RequireRole"}}
    papel -->|não| dono{"grupo do cliente?"}

    role -->|"role ≠ admin"| e403a(["403"])
    role -->|"role = admin"| handler["handler"]

    dono -->|sim| own{{"authorizeDocument"}}
    dono -->|não| handler

    own -->|"documento ≠ o do token"| e403b(["403<br/>ErrDocumentMismatch"])
    own -->|"igual, ou token sem documento"| handler

    style e401 fill:#ffe4e1,stroke:#c66
    style e403a fill:#ffe4e1,stroke:#c66
    style e403b fill:#ffe4e1,stroke:#c66
    style auth fill:#fff4e0,stroke:#d90
    style role fill:#fff4e0,stroke:#d90
    style own fill:#fff4e0,stroke:#d90
```

### A checagem de posse

`authorizeDocument` compara o documento pedido com a claim `document` do token:

| Token | Documento pedido | Resultado |
|---|---|---|
| `document = 529…725` | `529…725` | passa |
| `document = 529…725` | `390…705` | **403** |
| sem claim `document` (token de operação) | qualquer | passa — a operação já tem acesso pelas rotas administrativas |

Sem ela, qualquer cliente autenticado leria as ordens de serviço de qualquer outro conhecendo apenas o CPF. Registrada na [RFC-0003](../rfc/0003-autenticacao-por-cpf.md).

---

## 4. Inventário completo — 37 rotas

Todas sob o prefixo `/api/v1`.

### Abertas (3)

| Método | Rota | Quem usa |
|---|---|---|
| `POST` | `/auth/register` | operação |
| `POST` | `/auth/login` | operação |
| `POST` | `/auth/refresh` | operação |

Fora do prefixo: `GET /health` e `GET /swagger/*any`.

### Autenticadas, qualquer papel (1)

| Método | Rota |
|---|---|
| `POST` | `/auth/logout` |

### Cliente — autenticado por CPF, com posse do documento (3)

| Método | Rota | Papel no fluxo |
|---|---|---|
| `GET` | `/service-orders/code/{code}` | consulta o orçamento |
| `GET` | `/service-orders/requester` | lista as próprias OS |
| `PUT` | `/service-orders/code/{code}/status` | aprova (`in_execution`) ou recusa (`received`) |

### Administrativas — exigem `role: admin` (30)

| Recurso | Rotas |
|---|---|
| **Solicitantes** | `POST /requesters` · `GET /requesters` · `GET /requesters/{id}` · `GET /requesters/document/{document}` · `PUT /requesters/{id}` · `DELETE /requesters/{id}` · `GET /requesters/{id}/vehicles` |
| **Veículos** | `POST /vehicles` · `GET /vehicles` · `GET /vehicles/{id}` · `PUT /vehicles/{id}` · `DELETE /vehicles/{id}` |
| **Serviços** | `POST /services` · `GET /services` · `GET /services/{id}` · `PUT /services/{id}` · `DELETE /services/{id}` |
| **Peças** | `POST /parts` · `GET /parts` · `GET /parts/{id}` · `PUT /parts/{id}` · `DELETE /parts/{id}` · `PATCH /parts/{id}/stock` |
| **Ordens de serviço** | `POST /service-orders` · `GET /service-orders` · `GET /service-orders/{id}` · `PUT /service-orders/{id}` · `PUT /service-orders/{id}/status` · `DELETE /service-orders/{id}` · `GET /service-orders/metrics/execution-time` |

---

## 5. Mapa de códigos de resposta

Onde cada código pode nascer, e quem o produz.

```mermaid
flowchart LR
    subgraph gwl["No gateway"]
        g404["404<br/>rota não declarada"]
        g401["401<br/>authorizer negou"]
        g429["429<br/>throttling<br/>100 req/s, burst 200"]
    end

    subgraph lam["Na Lambda issuer"]
        l400["400<br/>payload inválido<br/>ou CPF inválido"]
        l404["404<br/>cliente não encontrado"]
        l403["403<br/>cliente inativo"]
        l500["500<br/>falha ao consultar<br/>ou assinar"]
    end

    subgraph appl["Na aplicação"]
        a401["401<br/>token ausente, inválido<br/>ou sem sub/role"]
        a403["403<br/>papel insuficiente,<br/>documento alheio,<br/>transição proibida"]
        a404["404<br/>recurso inexistente"]
        a409["409<br/>já existe"]
        a422["422<br/>regra de domínio violada"]
        a500["500<br/>erro interno"]
    end

    style g401 fill:#ffe4e1,stroke:#c66
    style a401 fill:#ffe4e1,stroke:#c66
    style a403 fill:#ffe0b2,stroke:#e80
    style l403 fill:#ffe0b2,stroke:#e80
```

O mesmo código pode vir de camadas diferentes, e isso importa ao depurar. Um `401` do gateway **não aparece nos logs da aplicação** — a requisição nunca chega no pod. Se o log da aplicação está vazio para uma chamada que o cliente jura ter feito, o suspeito é o authorizer, e a evidência está no access log do gateway em CloudWatch.

---

## 6. Fronteira entre repositórios

Quem declara o quê:

```mermaid
flowchart TB
    subgraph k8s["postech-tc3-infra-k8s"]
        api["aws_apigatewayv2_api"]
        rotas["4 rotas + 2 integrações"]
        authz["aws_apigatewayv2_authorizer"]
        perm["aws_lambda_permission"]
    end

    subgraph lam["postech-tc3-lambda-auth"]
        fi["função issuer"]
        fa["função authorizer"]
    end

    subgraph app["postech-tc3-app"]
        gin["grupos de rota do Gin"]
        svc["Service LoadBalancer → NLB"]
    end

    fi -.->|"invoke_arn"| rotas
    fa -.->|"invoke_arn"| authz
    svc -.->|"app_backend_url"| rotas
    perm -.->|"autoriza invocação"| lam

    style k8s fill:#e8f0fe,stroke:#4a7
    style lam fill:#fff4e0,stroke:#d90
    style app fill:#e6f4ea,stroke:#4a7
```

As três ligações pontilhadas são **valores passados à mão** entre repositórios, não referências de Terraform — os states são separados. A ordem de aplicação é:

1. `infra-database` → RDS de pé.
2. `infra-k8s` → cluster e gateway sem rotas (as variáveis de Lambda e backend ficam vazias).
3. `lambda-auth` → publica as funções e imprime os `invoke_arn` no resumo do job.
4. `app` → deploy cria o NLB; `kubectl get svc workshop-api` dá o hostname.
5. `infra-k8s` de novo → agora com `lambda_issuer_invoke_arn`, `lambda_authorizer_invoke_arn` e `app_backend_url` preenchidos.

O passo 5 existe porque as variáveis são opcionais de propósito: sem elas o gateway sobe vazio, o que quebra a dependência circular entre gateway, Lambda e NLB.
