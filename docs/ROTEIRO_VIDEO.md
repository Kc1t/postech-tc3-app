# Roteiro do vídeo — Tech Challenge Fase 3

Limite de **15 minutos**. O enunciado lista seis itens obrigatórios; o roteiro abaixo cobre todos com folga de tempo.

| Item exigido | Bloco |
|---|---|
| Autenticação com CPF | 3 |
| Execução da pipeline CI/CD | 2 |
| Deploy automatizado | 2 |
| Consumo das APIs protegidas | 4 |
| Dashboard de monitoramento com análise ao vivo | 6 |
| Logs e traces em execução | 5 |

---

## Antes de gravar

- [ ] Infraestrutura de pé: RDS, EKS, gateway e as duas Lambdas
- [ ] Dashboard do New Relic importado, com dados de pelo menos 1 hora
- [ ] Um cliente **ativo** e um **inativo** cadastrados, com os CPFs anotados
- [ ] Uma OS já aberta, para consultar sem precisar criar na hora
- [ ] Terminal com fonte grande; abas pré-abertas: GitHub Actions, New Relic, terminal `kubectl`
- [ ] Um branch pronto com um commit trivial, para abrir o PR ao vivo
- [ ] Roteiro de `curl` já em um arquivo, para copiar e colar sem digitar errado

---

## Bloco 1 — Contexto e arquitetura (2 min)

Abrir com o diagrama de componentes (`docs/diagrams/componentes.md` renderizado no GitHub).

Falar:
- O que mudou da Fase 2 para a 3: um repositório virou quatro, cada um com pipeline próprio.
- Mostrar os quatro repositórios lado a lado e o `soat-architecture` na lista de colaboradores.
- Apontar no diagrama o caminho da requisição: cliente → gateway → authorizer → NLB → pod → RDS.

Não gastar mais de 2 minutos aqui. O avaliador quer ver funcionando.

---

## Bloco 2 — CI/CD e deploy automatizado (3 min)

1. Abrir o PR preparado. Mostrar a pipeline rodando: lint, dependências, testes e **o gate de cobertura de 80%**.
2. Mostrar que a `main` está protegida e não aceita push direto.
3. Fazer o merge. Acompanhar o job de deploy: build da imagem, push no ECR, `kubectl apply`, `rollout status`.
4. No terminal, confirmar:

```bash
kubectl get pods -n postech
kubectl get deployment workshop-api -n postech -o wide
```

Falar enquanto roda: o mesmo pipeline entrega em `postech-homolog` quando o push é na `homolog` — um cluster, dois namespaces, decisão de custo registrada em ADR.

---

## Bloco 3 — Autenticação com CPF (3 min)

Este é o coração da fase. Mostrar os quatro casos.

**Sucesso:**

```bash
curl -X POST "$GATEWAY/auth" \
  -H 'Content-Type: application/json' \
  -d '{"cpf": "529.982.247-25"}'
```

Mostrar o token na resposta e colar em <https://jwt.io> para exibir as claims `sub`, `role`, `document` e `exp`.

**CPF inválido → 400:**

```bash
curl -i -X POST "$GATEWAY/auth" -d '{"cpf": "111.111.111-11"}'
```

**Cliente inexistente → 404:**

```bash
curl -i -X POST "$GATEWAY/auth" -d '{"cpf": "390.533.447-05"}'
```

**Cliente inativo → 403:** usar o CPF do cliente desativado.

Falar: a Lambda valida os dígitos verificadores antes de tocar no banco, e depois consulta existência **e** status — a coluna `status` foi adicionada nesta fase justamente para isso.

---

## Bloco 4 — APIs protegidas (2 min)

**Sem token → 401 no gateway, sem chegar na aplicação:**

```bash
curl -i "$GATEWAY/api/v1/service-orders/code/1042"
```

**Com token → 200:**

```bash
curl -i "$GATEWAY/api/v1/service-orders/code/1042" \
  -H "Authorization: Bearer $TOKEN"
```

**Aprovar o orçamento como cliente:**

```bash
curl -i -X PUT "$GATEWAY/api/v1/service-orders/code/1042/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"status": "in_execution", "requester_document": "52998224725"}'
```

Falar: na Fase 2 essas rotas eram públicas. Agora exigem o JWT emitido a partir do CPF, e a validação acontece duas vezes — no gateway e na aplicação.

---

## Bloco 5 — Logs e traces (2 min)

Pegar o `X-Correlation-ID` do header da resposta anterior. No New Relic:

```sql
SELECT timestamp, level, msg, path, status, latency_ms
FROM Log WHERE correlation_id = '<id>' ORDER BY timestamp
```

Mostrar a requisição inteira reconstruída, incluindo o evento de domínio `service_order_status_changed`.

Depois, o APM: **APM & Services → workshop-api-postech → Distributed tracing**. Abrir um trace e mostrar o tempo gasto em banco.

---

## Bloco 6 — Dashboard e alertas ao vivo (3 min)

Abrir o dashboard e passar pelas três páginas.

**Negócio:** volume diário de OS, tempo médio de execução, transições por status.

**APIs:** latência p95, taxa de sucesso, rotas mais lentas.

**Kubernetes:** CPU, memória e réplicas.

**Análise ao vivo — gerar carga e mostrar o HPA reagindo:**

```bash
# terminal 1
kubectl get hpa -n postech -w

# terminal 2
kubectl get pods -n postech -w

# terminal 3 — o alvo é o NLB, não o gateway
NLB=$(kubectl get svc workshop-api -n postech -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
kubectl run carga -n postech --rm -it --image=williamyeh/hey --restart=Never -- \
  -z 180s -c 50 "http://$NLB/health"
```

**A carga vai no NLB de propósito.** O API Gateway está com `throttling_rate_limit` de 100 req/s e burst de 200; com `-c 50` a carga passa disso em segundos e o gateway devolve 429 antes de a requisição chegar no pod. O resultado na gravação seria o pior possível: erro na tela, CPU parada e o HPA imóvel, parecendo que a escalabilidade não funciona quando o que funcionou foi o rate limit. O NLB não tem esse teto e é o caminho que exercita o pod de verdade.

Narrar: CPU sobe, o HPA aumenta as réplicas, novos pods aparecem em `Pending` e depois `Running`, a latência estabiliza. Voltar ao dashboard e mostrar o mesmo movimento no gráfico de réplicas.

Se quiser mostrar o rate limit também, vale como cena curta e separada — aí sim contra o gateway, com o 429 como comportamento esperado, não como acidente.

**Alertas:** abrir a condição "Falha no processamento de ordens de serviço" e explicar o gatilho.

---

## Encerramento (30 s)

Fechar com o checklist do enunciado no `DOCUMENTO_ENTREGA_FASE3.md`, apontando onde cada item foi demonstrado.

---

## Distribuição do tempo

| Bloco | Minutos | Acumulado |
|---|---|---|
| 1 — Arquitetura | 2 | 2 |
| 2 — CI/CD e deploy | 3 | 5 |
| 3 — Autenticação por CPF | 3 | 8 |
| 4 — APIs protegidas | 2 | 10 |
| 5 — Logs e traces | 2 | 12 |
| 6 — Dashboard ao vivo | 3 | 15 |

Sem folga. Se estourar, cortar do bloco 1 — é o único que não é exigido nominalmente.

## Armadilhas de gravação

- **Credencial do Learner Lab expira em 4 h.** Gravar no começo da sessão.
- **O `terraform apply` do EKS leva ~15 min.** Subir a infraestrutura antes, nunca durante.
- **O HPA demora para reagir.** O `metrics-server` publica métricas a cada ~15 s e o HPA reavalia na mesma cadência. Gerar carga por pelo menos 2 minutos antes de esperar movimento.
- **A `ResourceQuota` do namespace tem teto de 14 pods** e o pod `carga` conta nele. Com o HPA em 10 réplicas ainda sobra folga, mas conferir `kubectl get pods -n postech` antes de começar, para não entrar na gravação com pods de teste pendurados.
- **Cold start da Lambda em VPC.** A primeira chamada ao `/auth` pode levar alguns segundos; fazer uma chamada de aquecimento antes de gravar.
- **Não mostrar a tela de credenciais da AWS nem o `JWT_SECRET`** em nenhum momento.
