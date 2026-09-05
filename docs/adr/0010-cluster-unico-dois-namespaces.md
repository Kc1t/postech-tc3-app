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

O RDS segue a mesma lógica de economia: uma instância com dois databases lógicos.

## Alternativas consideradas

1. **Dois clusters completos.** Isolamento real de plano de controle e de rede, e é o desenho correto para produção de verdade. Descartado pelo custo: dobra a conta e não acrescenta nada demonstrável na avaliação.

2. **Um cluster, um namespace só, diferenciando por tag de imagem.** Mais barato ainda, mas eliminaria a noção de ambiente — não haveria onde validar antes de produção, e o requisito de "deploy automático das branches de homologação e produção" ficaria descumprido de fato, não só na forma.

3. **Namespaces mais isolamento por NetworkPolicy e ResourceQuota.** É a evolução natural desta decisão e teria sido incluída se houvesse tempo. Registrada como dívida.

4. **`kind` ou `k3s` para homologação, EKS só para produção.** Eliminaria o custo de staging, mas homologar em um Kubernetes diferente do de produção invalida boa parte do valor de homologar.

## Consequências

### Positivas

- Custo cai pela metade: ~US$ 0,21/hora em vez de ~US$ 0,42.
- O requisito de deploy automático por branch é atendido, com ambientes separados e aprovação independente via GitHub Environments.
- Manifests sem namespace fixo ficam reutilizáveis, e o mesmo YAML serve aos dois ambientes sem duplicação nem template engine.
- Um cluster só significa um `metrics-server`, um agente do New Relic e um HPA para observar — menos superfície para configurar errado sob prazo.

### Negativas

- **Não há isolamento de plano de controle.** Um erro de configuração com escopo de cluster — RBAC, CRD, webhook — afeta os dois ambientes.
- **Não há isolamento de recursos.** Sem `ResourceQuota`, um teste de carga em homologação pode consumir capacidade de nó e provocar `Pending` em produção. É a consequência mais séria desta decisão.
- Namespaces não isolam rede por padrão: os pods de `postech-homolog` alcançam os de `postech` sem `NetworkPolicy`.
- Um RDS compartilhado significa que uma migration destrutiva em homologação atinge a mesma instância que serve produção. Databases lógicos separados limitam o dano, mas não isolam recursos de I/O.

Nenhuma dessas consequências é aceitável em produção real. São aceitáveis aqui porque o sistema é acadêmico, tem prazo fixo e orçamento fechado — e a decisão está registrada para que quem herdar o projeto saiba que é dívida deliberada, não descuido.

## Referências

- Amazon EKS — [Pricing](https://aws.amazon.com/eks/pricing/)
- Implementação: [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml), job `deploy`
