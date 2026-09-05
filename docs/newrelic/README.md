# Observabilidade — New Relic

Tudo que a Fase 3 exige de monitoramento: dashboards, alertas e a origem de cada métrica.

## Por que New Relic e não Datadog

O enunciado deixa a escolha livre entre os dois, e o curso cobre ambos com a mesma profundidade. O critério foi o free tier:

| | New Relic | Datadog |
|---|---|---|
| Ingestão gratuita | 100 GB/mês, sem cartão | trial de 14 dias |
| APM no plano gratuito | incluído | pago à parte |
| Hosts | ilimitado no free | 5 no trial |

Com orçamento fechado e prazo de 11 dias, só o New Relic fecha sem custo.

## Como os dados chegam

| Sinal | Caminho |
|---|---|
| **Logs** | A aplicação escreve JSON em `stdout`; o `newrelic-logging` do `nri-bundle` coleta o stdout dos pods e envia. Os campos do JSON viram atributos consultáveis. |
| **APM e traces** | Agente Go (`newrelic/go-agent`) com `nrgin`, ativado apenas quando `NEW_RELIC_LICENSE_KEY` está presente. |
| **Infra e Kubernetes** | `newrelic-infrastructure`, `kube-state-metrics` e `nri-kube-events`, instalados pelo pipeline do `postech-tc3-infra-k8s`. |
| **Lambdas** | CloudWatch Logs. |

O encaminhamento de log pelo próprio agente Go fica **desligado** (`ConfigAppLogForwardingEnabled(false)`): o coletor do cluster já lê o `stdout`, e ligar os dois duplicaria cada linha e consumiria o dobro da cota.

## Importar o dashboard

1. New Relic → **Dashboards** → **Import dashboard**.
2. Cole o conteúdo de [`dashboard.json`](./dashboard.json).
3. **Substitua os `"accountIds": [0]`** pelo id da sua conta antes de importar — são 18 ocorrências.

```bash
sed -i "s/\"accountIds\": \[0\]/\"accountIds\": [SEU_ACCOUNT_ID]/g" docs/newrelic/dashboard.json
```

O dashboard tem três páginas:

| Página | Cobre |
|---|---|
| **Negócio** | Volume diário de OS, tempo médio de execução, transições por status, erros de integração |
| **APIs** | Latência p50/p95/p99, taxa de sucesso, rotas mais lentas, resultado da autenticação por CPF, rastreio por `correlation_id` |
| **Kubernetes** | CPU e memória por pod, réplicas do HPA, pods indisponíveis, reinícios, uptime |

## Origem de cada métrica exigida

| Requisito do enunciado | Fonte | Query |
|---|---|---|
| Latência das APIs | campo `latency_ms` do `RequestLogger` | `SELECT percentile(latency_ms, 50, 95, 99) FROM Log WHERE msg = 'request'` |
| Consumo de CPU e memória | `K8sContainerSample` | `SELECT average(cpuUsedCores) FROM K8sContainerSample WHERE containerName = 'workshop-api'` |
| Healthchecks e uptime | `K8sDeploymentSample` | `SELECT percentage(count(*), WHERE podsReady >= 1) FROM K8sDeploymentSample` |
| Alertas de falha no processamento de OS | logs `request` com status 5xx em `/api/v1/service-orders` | ver alerta 1 abaixo |
| Logs estruturados com correlação | `correlation_id` presente em toda linha | `SELECT * FROM Log WHERE correlation_id = '...'` |
| Volume diário de OS | evento `service_order_created` | `SELECT count(*) FROM Log WHERE msg = 'service_order_created' TIMESERIES 1 day` |
| Tempo médio por status | evento `service_order_status_changed` com `execution_seconds` | `SELECT average(execution_seconds) FROM Log WHERE execution_seconds IS NOT NULL` |
| Erros e falhas nas integrações | evento `integration_failure` | `SELECT count(*) FROM Log WHERE msg = 'integration_failure' FACET integration, stage` |

O `execution_seconds` é calculado no momento da transição, a partir de `started_at` e `finished_at` da própria OS — as mesmas colunas que existem no banco. A métrica é auditável contra o `SELECT`.

## Alertas

Criar em **Alerts → Alert conditions → NRQL**. Todos usam janela de 5 minutos.

### 1. Falha no processamento de ordens de serviço

Exigido nominalmente pelo enunciado.

```sql
SELECT count(*) FROM Log
WHERE msg = 'request' AND path LIKE '/api/v1/service-orders%' AND status >= 500
```

- Crítico: acima de 0 por 5 minutos.

### 2. Falha de integração

```sql
SELECT count(*) FROM Log WHERE msg = 'integration_failure'
```

- Aviso: acima de 3 em 5 minutos.
- Crítico: acima de 10 em 5 minutos.

### 3. Uptime — pods indisponíveis

```sql
SELECT latest(podsReady) FROM K8sDeploymentSample WHERE deploymentName = 'workshop-api'
```

- Crítico: abaixo de 1 por 3 minutos.
- Aviso: abaixo de 2 por 10 minutos — abaixo do `minReplicas` do HPA.

### 4. Latência degradada

```sql
SELECT percentile(latency_ms, 95) FROM Log WHERE msg = 'request' AND path NOT LIKE '/health%'
```

- Aviso: acima de 1000 ms por 5 minutos.
- Crítico: acima de 3000 ms por 5 minutos.

### 5. Saturação de recursos

```sql
SELECT average(cpuUsedCores) / average(cpuLimitCores) * 100
FROM K8sContainerSample WHERE containerName = 'workshop-api'
```

- Aviso: acima de 85% por 10 minutos. Indica que o HPA chegou ao teto ou não está escalando.

### 6. Autenticação sob ataque

Não é exigido, mas o CPF é um dado público — vale monitorar.

```sql
SELECT count(*) FROM Log WHERE msg IN ('cliente nao encontrado', 'cpf invalido')
```

- Aviso: acima de 50 em 5 minutos, o que sugere enumeração de CPFs.

## Rastrear uma requisição ponta a ponta

Toda resposta traz o header `X-Correlation-ID`. Com ele:

```sql
SELECT timestamp, level, msg, path, status, latency_ms
FROM Log WHERE correlation_id = '<id da resposta>'
ORDER BY timestamp
```

Se o cliente enviar o header na requisição, o valor dele é reaproveitado em vez de gerado — o que permite correlacionar com sistemas externos.

## Configuração

| Variável | Onde | Efeito |
|---|---|---|
| `NEW_RELIC_LICENSE_KEY` | Secret do Kubernetes, vindo do GitHub Secret | Sem ela, o APM não sobe e a aplicação registra que a observabilidade está desligada |
| `NEW_RELIC_APP_NAME` | ConfigMap, com o namespace injetado pelo pipeline | Separa `workshop-api-postech` de `workshop-api-postech-homolog` no APM |
| `LOG_LEVEL` | ConfigMap | `info` em produção; `debug` aumenta bastante a ingestão |

O mesmo `NEW_RELIC_LICENSE_KEY` alimenta o agente do cluster, instalado pelo pipeline do `postech-tc3-infra-k8s`.
