# Diagrama de Sequência — Abertura de Ordem de Serviço

Da abertura da OS pela operação até a aprovação do orçamento pelo cliente autenticado por CPF.

```mermaid
sequenceDiagram
    autonumber
    actor OP as Atendente
    actor C as Cliente
    participant GW as API Gateway
    participant AZ as Lambda authorizer
    participant API as workshop-api
    participant UC as CreateServiceOrder
    participant DB as RDS PostgreSQL
    participant SMTP as Notificador SMTP
    participant NR as New Relic

    rect rgb(238, 245, 255)
    note over OP,DB: Abertura da OS — operação autenticada por e-mail e senha
    OP->>GW: POST /api/v1/service-orders<br/>Bearer <jwt admin>
    GW->>AZ: autoriza
    AZ-->>GW: isAuthorized = true
    GW->>API: HTTP_PROXY via NLB
    API->>API: CorrelationID + RequestLogger
    API->>API: middleware Auth + RequireRole(admin)

    API->>UC: Execute(input)
    UC->>DB: FindByDocument(cpf do solicitante)
    alt solicitante inexistente
        DB-->>UC: sem linhas
        UC-->>OP: 404
    else encontrado
        DB-->>UC: requester
        UC->>DB: FindByPlate(placa)
        alt veículo de outro solicitante
            UC-->>OP: 422 ErrVehicleNotFromRequester
        else vínculo válido
            UC->>UC: NewServiceOrder(requesterID, vehicleID)<br/>status = received
            UC->>DB: INSERT service_orders<br/>code sequencial
            DB-->>UC: OS criada
            UC-->>API: service order
            API-->>OP: 201 {id, code, status}
        end
    end
    API->>NR: log JSON com correlation_id, status, latency_ms
    end

    rect rgb(255, 248, 235)
    note over API,SMTP: Diagnóstico e orçamento
    OP->>API: PUT status → in_diagnosis
    OP->>API: PUT status → awaiting_approval
    API->>SMTP: notifica o cliente
    SMTP-->>C: e-mail com o código da OS
    end

    rect rgb(240, 250, 240)
    note over C,DB: Aprovação pelo cliente — protegida por CPF
    C->>GW: POST /auth {cpf}
    GW-->>C: 200 {access_token}

    C->>GW: GET /api/v1/service-orders/code/1042<br/>Bearer <jwt client>
    GW->>AZ: autoriza
    AZ-->>GW: isAuthorized = true
    GW->>API: HTTP_PROXY
    API-->>C: 200 orçamento

    C->>GW: PUT /api/v1/service-orders/code/1042/status<br/>{"status": "in_execution"}
    GW->>AZ: autoriza
    AZ-->>GW: isAuthorized = true
    GW->>API: HTTP_PROXY
    API->>API: AuthorizeRequesterTransition
    note right of API: cliente só pode ir para<br/>in_execution (aprova) ou<br/>received (recusa)
    API->>DB: UPDATE status, started_at = now()
    API-->>C: 200
    API->>NR: log + métrica de tempo por status
    end
```

## Máquina de estados da OS

```mermaid
stateDiagram-v2
    [*] --> received
    received --> in_diagnosis
    in_diagnosis --> awaiting_approval
    awaiting_approval --> in_execution: cliente aprova
    awaiting_approval --> received: cliente recusa
    in_execution --> finished
    finished --> delivered
    delivered --> [*]
```

Transição fora deste grafo retorna `ErrInvalidStatus`. O cliente autenticado por CPF só pode iniciar as duas transições marcadas; as demais pertencem ao fluxo operacional da oficina.

## Origem das métricas dos dashboards

| Dashboard exigido | Fonte |
|---|---|
| Volume diário de ordens de serviço | `created_at` de `service_orders` |
| Tempo médio por status | `started_at` (entrada em `in_execution`) e `finished_at` (entrada em `finished`) |
| Erros e falhas nas integrações | Logs de nível `error` com `correlation_id`, incluindo falhas do notificador SMTP |
| Latência das APIs | `latency_ms` do `RequestLogger` e o access log do gateway |

As colunas `started_at` e `finished_at` são gravadas por `UpdateStatus` no momento exato da transição — é o que torna a métrica de tempo médio auditável contra o banco.
