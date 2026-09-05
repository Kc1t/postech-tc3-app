# Diagrama de Sequência — Autenticação por CPF

Fluxo completo: o cliente troca o CPF por um JWT e usa esse token em uma rota protegida.

```mermaid
sequenceDiagram
    autonumber
    actor C as Cliente
    participant GW as API Gateway
    participant IS as Lambda issuer
    participant AZ as Lambda authorizer
    participant DB as RDS PostgreSQL
    participant API as workshop-api (EKS)

    rect rgb(238, 245, 255)
    note over C,DB: Emissão do token — rota aberta
    C->>GW: POST /auth {"cpf": "529.982.247-25"}
    GW->>IS: invoca (AWS_PROXY)
    IS->>IS: normaliza e valida dígitos verificadores

    alt CPF inválido
        IS-->>C: 400 cpf invalido
    else CPF válido
        IS->>DB: SELECT id, name, email, status<br/>FROM requesters WHERE document = $1
        alt não encontrado
            DB-->>IS: sem linhas
            IS-->>C: 404 cliente nao encontrado
        else status ≠ active
            DB-->>IS: status = inactive
            IS-->>C: 403 cliente inativo
        else ativo
            DB-->>IS: id, name, email, status
            IS->>IS: assina HS256<br/>sub, role=client, email, name, document, exp
            IS-->>C: 200 {access_token, token_type, expires_at}
        end
    end
    end

    rect rgb(240, 250, 240)
    note over C,API: Consumo de rota protegida — dupla validação
    C->>GW: GET /api/v1/service-orders/code/1042<br/>Authorization: Bearer <jwt>
    GW->>AZ: invoca authorizer (TTL de cache = 0)
    AZ->>AZ: valida assinatura, exp, sub e role<br/>rejeita algoritmo none

    alt token inválido ou ausente
        AZ-->>GW: isAuthorized = false
        GW-->>C: 401
    else token válido
        AZ-->>GW: isAuthorized = true<br/>context: subject, role, document
        GW->>API: HTTP_PROXY via NLB<br/>repassa o header Authorization
        API->>API: middleware Auth valida o mesmo JWT
        API->>DB: consulta a ordem de serviço
        DB-->>API: dados da OS
        API-->>C: 200 + X-Correlation-ID
    end
    end
```

## Pontos de decisão

| Ponto | Regra |
|---|---|
| Validação do CPF | Dígitos verificadores calculados localmente, antes de qualquer ida ao banco |
| Existência | `SELECT` por `document`, que tem índice único |
| Status | `status = 'active'`; qualquer outro valor resulta em 403 |
| Cache do authorizer | Zero — desativar um cliente tem efeito na requisição seguinte |
| Segunda validação | O middleware `Auth` da aplicação repete a verificação, cobrindo acesso direto ao NLB |

## Distinção entre 404 e 403

A resposta separa "não existe" de "existe mas está inativo". É informação útil para o cliente legítimo e, ao mesmo tempo, permite enumerar CPFs cadastrados. Mitigado por throttling no stage do gateway e registrado como risco aceito na [RFC-0003](../rfc/0003-autenticacao-por-cpf.md).
