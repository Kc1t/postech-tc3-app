# RFC-0002: Escolha do banco de dados gerenciado

- **Status:** Aceita
- **Data:** 2026-09
- **Autor:** Reservoir Devs
- **Decisão relacionada:** ADR-0003 (PostgreSQL), estendida aqui para a camada gerenciada

## Resumo

Propõe **Amazon RDS PostgreSQL 16**, instância `db.t3.micro` Single-AZ em homologação e Multi-AZ em produção, provisionado por Terraform em repositório dedicado.

## Motivação

A Fase 3 exige banco **gerenciado** e uma justificativa formal da escolha. Nas Fases 1 e 2 o PostgreSQL rodava em contêiner via `docker-compose` — adequado para desenvolvimento, mas descumpre o requisito atual e joga backup, patching e alta disponibilidade para cima do time.

Duas decisões separadas:

1. Qual engine? — já respondida no ADR-0003 e revalidada aqui.
2. Qual serviço gerenciado, e com qual configuração?

## Proposta

| Aspecto | Escolha |
|---|---|
| Engine | PostgreSQL 16 |
| Serviço | Amazon RDS |
| Classe | `db.t3.micro` |
| Armazenamento | 20 GB gp3, criptografado |
| Alta disponibilidade | Multi-AZ apenas em `prod` |
| Backup | 1 dia em `staging`, 7 dias em `prod` |
| Acesso | `publicly_accessible = false`, security group por porta 5432 |
| Credenciais | senha gerada por `random_password`, gravada no Secrets Manager |
| Observabilidade | Performance Insights e export de logs para CloudWatch |

## Justificativa da engine

O domínio é fortemente relacional e transacional:

- Abrir uma ordem de serviço grava a OS, referencia solicitante e veículo e baixa estoque de peças **em uma transação**. Falha parcial não pode existir.
- A OS é uma **máquina de estados** com transições fixas — `received → in_diagnosis → awaiting_approval → in_execution → finished → delivered`.
- CPF/CNPJ e placa são **identificadores naturais únicos**, garantidos por índice único no banco e não por código de aplicação.
- Os dashboards da fase exigem agregação por data e por status. É SQL puro.

O PostgreSQL especificamente entrega três coisas que usamos de fato:

- **`jsonb`** para o snapshot imutável de serviços e peças da OS, preservando o preço praticado sem tabela de histórico.
- **`gen_random_uuid()`** nativo desde a versão 13, sem extensão adicional.
- **Índices parciais e GIN**, disponíveis se as consultas de OS por conteúdo crescerem.

## Justificativa do serviço gerenciado

O RDS entrega sem código o que teríamos de construir e operar:

| Recurso | O que evita |
|---|---|
| Backup automático com retenção configurável | Rotina de `pg_dump`, armazenamento e teste de restauração |
| Criptografia em repouso por flag | Configuração manual de volume criptografado |
| Multi-AZ por flag | Streaming replication e failover manual |
| Performance Insights | Instrumentação de query lenta na aplicação |
| Patching gerenciado | Janela de manutenção e atualização de versão à mão |

Essas cinco linhas de Terraform substituem semanas de trabalho de operação — e nenhuma delas seria feita bem sob prazo de 11 dias.

## Alternativas avaliadas

| Alternativa | Por que não |
|---|---|
| **Aurora PostgreSQL** | Compatível e mais escalável, com a menor instância custando várias vezes mais que uma `db.t3.micro`. O volume do projeto não usa nada do que se paga a mais. |
| **Aurora Serverless v2** | Escala a zero em teoria, mas o piso de ACU cobrado é maior que o de uma `db.t3.micro` ociosa. Vantajoso para tráfego intermitente de verdade, não para tráfego de demonstração. |
| **MySQL / Aurora MySQL** | Sem ganho para este domínio, e custaria a reescrita de migrations, tipos `jsonb` e testes de repositório. |
| **DynamoDB** | Modelagem por padrão de acesso quebraria as agregações por data e status dos dashboards, e não impõe unicidade de CPF nem de placa. Exigiria emular integridade referencial na aplicação. |
| **PostgreSQL em pod no EKS** | Descumpre o requisito de banco gerenciado e devolve ao time backup, HA e patching. Também acopla o ciclo de vida do dado ao do cluster. |
| **Neon / Supabase** | Bons free tiers, mas fora da conta AWS onde está o crédito, e adicionariam um provedor não coberto por nenhuma RFC. |

## Ajustes no modelo relacional

A Fase 3 introduz `requesters.status` (`active` \| `inactive`), porque o enunciado exige que a Lambda de autenticação consulte **existência e status** do cliente, e o modelo herdado só respondia existência.

Detalhamento completo, incluindo a explicação dos relacionamentos e o diagrama ER, em [`docs/MODELAGEM_DE_DADOS.md`](../MODELAGEM_DE_DADOS.md).

## Impacto no custo

| Item | US$/hora | US$/mês 24/7 |
|---|---|---|
| `db.t3.micro` Single-AZ | 0,018 | 13,14 |
| 20 GB gp3 | — | 2,30 |
| Secrets Manager (1 segredo) | — | 0,40 |
| **Total staging** | **~0,018** | **~15,84** |

Multi-AZ em produção dobra o custo da instância. Como a demonstração roda em um único ambiente por vez, o impacto real no crédito é próximo de zero.

## Riscos

| Risco | Mitigação |
|---|---|
| Lambda não alcançar o RDS, que não é público | A Lambda roda dentro da VPC, com security group próprio; o SG do RDS libera o CIDR da VPC |
| Esgotar conexões com o HPA subindo para 10 réplicas | Monitorar `DatabaseConnections` no dashboard; `db.t3.micro` suporta ~85 conexões |
| `AutoMigrate` não cobrir mudança destrutiva | Documentado como limite conhecido; migrar para ferramenta versionada exigiria ADR própria |
| Senha em variável de ambiente da Lambda | Aceito conscientemente para evitar NAT Gateway; a alternativa é VPC endpoint para o Secrets Manager |

## Questões em aberto

- Adotar `golang-migrate` no lugar de `AutoMigrate` antes da entrega? Mais correto, mas exige converter o schema inteiro em migrations versionadas. Não cabe no prazo.
