# ADR-0010: Um cluster com dois namespaces em vez de dois clusters

- **Status:** Accepted
- **Data:** 2026-09
- **Decisão por:** Reservoir Devs

## Contexto

O enunciado exige **deploy automático das branches de homologação e produção**. A leitura literal seria provisionar dois ambientes completos e independentes.

O projeto roda com o crédito do AWS Academy Learner Lab, limitado a **US$ 100 por aluno**. O control plane do EKS custa **US$ 0,10/hora por cluster**, cobrado mesmo com o cluster ocioso. Um ambiente completo — control plane, dois nós `t3.medium`, RDS `db.t3.micro` e armazenamento — sai por aproximadamente **US$ 0,21/hora**.

Dois ambientes completos dobram isso para ~US$ 0,42/hora, o que consumiria o crédito em menos de dez dias de operação contínua.

## Decisão

**Um cluster EKS**, com dois namespaces:

| Branch | Namespace | GitHub Environment |
|---|---|---|
| `homolog` | `postech-homolog` | `staging` |
| `main` | `postech` | `prod` |

O pipeline resolve o namespace a partir de `github.ref_name` e aplica os manifests com `-n "$NAMESPACE"`. Os manifests **não declaram namespace** — foi removido de todos eles, e o `namespace.yaml` fixo foi descartado em favor de criação dinâmica pelo pipeline.

O RDS segue a mesma lógica de economia: **uma instância compartilhada**. Hoje os dois namespaces usam o mesmo `POSTGRES_DSN` e o mesmo `JWT_SECRET`, porque os secrets estão no nível do repositório e não nos GitHub Environments. Separar um database lógico e um segredo por ambiente é o próximo passo, e não exige mudança de infraestrutura — só secrets por environment.

Como os dois ambientes dividem os mesmos nós, cada namespace recebe **teto de recursos e política de rede próprios**, aplicados pelo pipeline antes do Deployment:

- `ResourceQuota` + `LimitRange` (`k8s/quota.yaml`) — teto dimensionado a partir do `maxReplicas: 10` do HPA mais um pod de surge: 2 CPU / 2Gi de requests e 6 CPU / 6Gi de limits. O `LimitRange` garante que nenhum container entre no namespace sem requests declarados, o que furaria o teto por omissão.
- `NetworkPolicy` (`k8s/networkpolicy.yaml`) — entrada liberada para pods do mesmo namespace e para a porta `8080` (NLB e probes do kubelet, que chegam com IP de nó, não de pod). Como a `8080` é a única porta da aplicação, na prática a política documenta a intenção mais do que bloqueia tráfego; o isolamento real da API depende de fechar o NLB (ver ADR-0007).

A cota fica **acima** do pior caso de propósito. Uma cota abaixo do `maxReplicas` do HPA é pior que nenhuma: o autoscaler segue tentando escalar, o ReplicaSet acumula `FailedCreate`, e o sintoma se disfarça de "o HPA não funciona".

No EKS a `NetworkPolicy` exige o addon `vpc-cni` com `enableNetworkPolicy = "true"` — sem ele o objeto é aceito e silenciosamente ignorado. O addon está declarado em `postech-tc3-infra-k8s/addons.tf`, pela mesma razão que o `metrics-server` está: no EKS, o recurso que parece existir só funciona se o addon correspondente estiver ligado.

O detalhamento do dimensionamento está em [`k8s/README.md`](../../k8s/README.md).

### Os outros três repositórios

O mesmo raciocínio de custo define o que a branch `homolog` faz nos repositórios que não são a aplicação:

| Repositório | Push em `homolog` | Push em `main` |
|---|---|---|
| `postech-tc3-lambda-auth` | `terraform apply` com `envs/staging.tfvars`: functions `postech-tc3-staging-auth-*`, state `lambda-auth/staging.tfstate` | `terraform apply` em produção |
| `postech-tc3-infra-k8s` | `terraform plan` contra o state de produção | `terraform apply` |
| `postech-tc3-infra-database` | `terraform plan` contra o state de produção | `terraform apply` |

A Lambda cobra por invocação, não por hora, então um par de staging de verdade não pesa no orçamento. O cluster e o banco cobram por hora: aplicar a homolog nesses dois repositórios criaria um segundo EKS e um segundo RDS, justamente o que esta decisão evita. Neles, a homolog valida a mudança com um `plan` contra a infraestrutura que existe, e só a `main` aplica.

As functions de staging não têm rota no gateway, que é único e aponta para as de produção. Elas são exercitadas por invocação direta (`aws lambda invoke`).

## Alternativas consideradas

1. **Dois clusters completos.** Isolamento real de plano de controle e de rede, e é o desenho correto para produção de verdade. Descartado pelo custo: dobra a conta e não acrescenta nada demonstrável na avaliação.

2. **Um cluster, um namespace só, diferenciando por tag de imagem.** Mais barato ainda, mas eliminaria a noção de ambiente — não haveria onde validar antes de produção, e o requisito de "deploy automático das branches de homologação e produção" ficaria descumprido de fato, não só na forma.

3. **Namespaces sem nenhum isolamento adicional.** Era o desenho inicial. Descartado depois de dimensionar o risco: o `t3.medium` tem 2 vCPU, o node group parte de 2 nós, e o HPA de homologação sozinho pode pedir 11 pods. Sem teto, a demonstração de escalabilidade em homologação derruba produção. O isolamento por `ResourceQuota` e `NetworkPolicy` deixou de ser evolução futura e entrou no escopo.

4. **`kind` ou `k3s` para homologação, EKS só para produção.** Eliminaria o custo de staging, mas homologar em um Kubernetes diferente do de produção invalida boa parte do valor de homologar.

## Consequências

### Positivas

- Custo cai pela metade: ~US$ 0,21/hora em vez de ~US$ 0,42.
- O requisito de deploy automático por branch é atendido: cada push em `homolog` ou `main` dispara o pipeline do ambiente correspondente nos quatro repositórios.
- Manifests sem namespace fixo ficam reutilizáveis, e o mesmo YAML serve aos dois ambientes sem duplicação nem template engine.
- Um cluster só significa um `metrics-server`, um agente do New Relic e um HPA para observar — menos superfície para configurar errado sob prazo.

### Negativas

- **Não há isolamento de plano de controle.** Um erro de configuração com escopo de cluster — RBAC, CRD, webhook — afeta os dois ambientes.
- **O teto de recursos é por namespace, não por nó.** A `ResourceQuota` impede que um namespace peça mais do que lhe cabe, mas os dois somados ainda podem exceder a capacidade instalada: 2 × 2 CPU de requests contra 4 vCPU de dois `t3.medium`. O node group aceita até 5 nós, mas não há Cluster Autoscaler: ele fica nos 2 nós desejados, e réplicas acima da capacidade ficam em `Pending` até alguém aumentar o `desired_size`. O risco caiu de "homologação derruba produção" para "produção não escala além da capacidade instalada".
- A `NetworkPolicy` mantém a porta `8080` aberta a origens de fora do namespace, porque é por onde chegam o NLB e as probes. O isolamento de rede é do restante do tráfego, não da API em si.
- Um RDS compartilhado significa que uma migration destrutiva em homologação atinge a mesma instância — e hoje o mesmo database — que serve produção. Databases lógicos separados limitariam o dano, mas não isolariam recursos de I/O.
- A infraestrutura de cluster e banco não tem um ambiente onde uma mudança seja aplicada antes de produção; o `plan` na homolog mostra o efeito, mas não o exercita.

Nenhuma dessas consequências é aceitável em produção real. São aceitáveis aqui porque o sistema é acadêmico, tem prazo fixo e orçamento fechado — e a decisão está registrada para que quem herdar o projeto saiba que é dívida deliberada, não descuido.

## Referências

- Amazon EKS — [Pricing](https://aws.amazon.com/eks/pricing/)
- Kubernetes — [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
- AWS — [Amazon VPC CNI network policies](https://docs.aws.amazon.com/eks/latest/userguide/cni-network-policy.html)
- Implementação: [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml) job `deploy`, [`k8s/quota.yaml`](../../k8s/quota.yaml), [`k8s/networkpolicy.yaml`](../../k8s/networkpolicy.yaml)
