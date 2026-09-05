# RFC-0001: Escolha da nuvem

- **Status:** Aceita
- **Data:** 2026-09
- **Autor:** Reservoir Devs
- **Decisão relacionada:** consolidada em ADR quando estabilizar

## Resumo

Propõe manter a **AWS** como provedor único para a Fase 3, apesar de o control plane do EKS ser o item mais caro do orçamento e de existirem alternativas com control plane gratuito.

## Motivação

O enunciado dá **livre escolha de nuvem**, exigindo apenas API Gateway, function serverless, banco gerenciado, cluster Kubernetes e Terraform. A escolha precisa ser justificada formalmente, e o momento de revisá-la é agora — antes de existir qualquer state de Terraform aplicado.

Duas forças puxam em direções opostas:

- **Custo:** o control plane do EKS custa US$ 73/mês. AKS oferece control plane gratuito no tier Free; GKE dá crédito mensal de US$ 74,40 que cobre um cluster zonal. Sair da AWS zeraria essa linha.
- **Coerência:** o módulo de Desenvolvimento Serverless da fase é integralmente AWS — Lambda, SAM, API Gateway e Cognito. As Fases 1 e 2 já entregaram Terraform de EKS e RDS.

## Proposta

Manter **AWS**, região `us-east-1`, com:

| Requisito | Serviço |
|---|---|
| API Gateway | Amazon API Gateway v2 (HTTP API) |
| Function serverless | AWS Lambda (`provided.al2023`, arm64) |
| Banco gerenciado | Amazon RDS PostgreSQL 16 |
| Cluster Kubernetes | Amazon EKS |
| IaC | Terraform com state em S3 |

## Alternativas avaliadas

### Azure

| A favor | Contra |
|---|---|
| Control plane do AKS é gratuito no tier Free | Toda a Fase 2 seria reescrita |
| O módulo de API Gateway da fase ensina Azure API Management | O módulo de Serverless é AWS: Lambda viraria Azure Functions, sem correspondência com a aula |
| US$ 200 de crédito por 30 dias | O crédito do Learner Lab da FIAP é AWS e ficaria ocioso |

### Google Cloud

| A favor | Contra |
|---|---|
| Crédito mensal cobre um cluster zonal | Nenhum conteúdo da fase usa GCP |
| US$ 300 por 90 dias | Mesma reescrita completa da Azure |

### Arquitetura híbrida — cluster em AKS, Lambda na AWS

Economizaria os US$ 73 mantendo o alinhamento do módulo serverless. Rejeitada: a latência e a complexidade de rede entre o gateway da AWS e um cluster no Azure não se pagam, e a justificativa de escolha de nuvem — item explicitamente cobrado — ficaria incoerente com duas nuvens em produção.

## Justificativa da decisão

1. **A economia real é menor do que a etiqueta sugere.** Os US$ 73/mês pressupõem o cluster ligado 24/7. Com `terraform destroy` ao fim de cada sessão de trabalho, a operação dos 11 dias custa cerca de **US$ 14**, dos quais o control plane responde por ~US$ 7. Migrar de nuvem para economizar US$ 7 não se justifica.

2. **O crédito disponível é da AWS.** O Learner Lab dá US$ 100 por aluno em conta AWS. Migrar para Azure ou GCP significaria não usar o crédito que já existe e depender de trial com cartão de crédito.

3. **A fase é majoritariamente AWS no que importa para a entrega.** A autenticação serverless — item mais pesado da avaliação — é ensinada com Lambda, SAM e Cognito. Sair da AWS descolaria a entrega do conteúdo.

4. **Continuidade de Fases 1 e 2.** Existe Terraform de EKS e RDS funcionando, e pipeline de deploy no EKS já validado. Esse é capital que se perde inteiro numa migração.

5. **O prazo é o recurso escasso, não o dinheiro.** São 11 dias. Reescrever dois repositórios de infraestrutura consumiria o tempo que precisa ir para autenticação, observabilidade e documentação.

## Impacto no custo

Estimativa `us-east-1`, um ambiente, on-demand:

| Item | US$/hora | US$/mês 24/7 |
|---|---|---|
| EKS control plane | 0,100 | 73,00 |
| 2× `t3.medium` | 0,083 | 60,74 |
| RDS `db.t3.micro` | 0,018 | 13,14 |
| Armazenamento (EBS + RDS) | 0,008 | 5,50 |
| Gateway, Lambda, ECR, logs | ~0 | ~1,40 |
| **Total** | **~0,21** | **~154** |

Medidas de contenção adotadas:

- Um cluster com dois namespaces em vez de dois clusters (ADR-0010) — corta 50%.
- Sem NAT Gateway: a Lambda recebe a string de conexão por variável de ambiente em vez de consultar o Secrets Manager em runtime — evita US$ 32,85/mês fixos.
- Versão do EKS mantida em suporte padrão. A 1.31 estava em *extended support*, cobrado a US$ 0,60/hora — **seis vezes** o valor normal. Corrigido para 1.35, com suporte padrão até 27/03/2027.

## Riscos

| Risco | Mitigação |
|---|---|
| Crédito do Learner Lab acabar antes da entrega | Destruir a infraestrutura ao fim de cada sessão; monitorar o consumo no painel do lab |
| Sessão de 4 h do lab derrubar recursos no meio de um `apply` de EKS (~15 min) | Iniciar `apply` no começo da sessão, nunca no fim |
| Credencial temporária expirar e quebrar os pipelines | `scripts/sync-aws-secrets.sh` propaga a credencial nova para os quatro repositórios |
| Cluster esquecido ligado | Verificação no encerramento de cada sessão; é o maior risco de consumo do crédito |

## Questões em aberto

- Vale mover o cluster para instâncias Spot? Reduziria ~70% do custo de nós, ao preço de interrupções durante a gravação do vídeo de demonstração. Não adotado.
