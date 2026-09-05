# ADR-0009: Deploy da Lambda por Terraform, não por AWS SAM

- **Status:** Accepted
- **Data:** 2026-09
- **Decisão por:** Reservoir Devs

## Contexto

O módulo de Desenvolvimento Serverless da Fase 3 ensina o **AWS SAM** como ferramenta de empacotamento e deploy de funções Lambda. Ao mesmo tempo, o enunciado exige explicitamente **Terraform para provisionamento** da infraestrutura, e os outros dois repositórios de infra (EKS e RDS) já usam Terraform com state remoto em S3.

Havia, portanto, uma tensão: seguir a ferramenta ensinada na aula ou manter uma única toolchain de infraestrutura.

## Decisão

Provisionar as duas funções Lambda — emissora de token e authorizer — com **Terraform**, no diretório `infra/` do próprio repositório da Lambda.

O Terraform é dono do código, não só da infraestrutura: `source_code_hash = filebase64sha256(...)` faz o `apply` detectar mudança no artefato e atualizar a função. Não existe passo separado de `aws lambda update-function-code`.

O pipeline fica: testes → `make package` (build arm64 + zip) → `terraform apply`.

## Alternativas consideradas

1. **AWS SAM.** É o que a aula ensina, e entrega vantagens reais: `sam local invoke` para depurar a função na máquina, e um template bem mais curto que o HCL equivalente. Descartado por três motivos:
   - Introduziria uma **segunda ferramenta de IaC** no projeto, com um segundo formato de state (CloudFormation), num prazo de 11 dias.
   - O SAM cria a execution role da função por padrão. O **AWS Academy Learner Lab não permite criar roles IAM** — seria preciso desligar essa parte do template e apontar para a `LabRole`, perdendo boa parte da concisão que justifica o SAM.
   - O enunciado pede Terraform. Usar SAM para uma peça e Terraform para as outras duas tornaria a documentação de infraestrutura incoerente.

2. **Terraform provisionando a função + `update-function-code` no CI.** Foi o desenho inicial. Descartado porque cria duas fontes de verdade para o código da função: o `apply` reverteria para o artefato do state sempre que rodasse depois de um deploy manual.

3. **Serverless Framework.** Mesmo problema do SAM, com dependência adicional de Node.js no pipeline.

4. **Imagem de contêiner em vez de zip.** Suportado pelo Lambda e reaproveitaria o ECR já usado pela aplicação. Descartado por peso: uma imagem de ~50 MB contra um zip de poucos MB, com cold start mais lento e sem ganho — o binário Go já é estático e não precisa de runtime empacotado.

## Consequências

### Positivas

- Uma única ferramenta de IaC nos quatro repositórios, um único formato de state, um único conjunto de convenções de pipeline.
- `terraform plan` no Pull Request mostra a mudança de infraestrutura **e** a de código da função no mesmo diff.
- A execution role é referenciada por `data "aws_iam_role"`, não criada — compatível com a restrição do Learner Lab sem gambiarra.
- Rollback é `git revert` seguido de `apply`: o hash do artefato anterior volta a valer.

### Negativas

- **Perdemos o `sam local invoke`.** Testar a função localmente exige subir o Runtime Interface Emulator à mão ou testar os pacotes internos isoladamente. Mitigado pela cobertura de testes em `internal/cpf`, `internal/token` e `internal/requester`, que concentram a lógica — o `main.go` é fino de propósito.
- O HCL é mais verboso que o template SAM equivalente, especialmente na configuração de VPC e log group.
- O artefato precisa existir no disco antes do `apply`, o que acopla a ordem dos passos no pipeline. Se `make package` falhar, o `apply` falha com erro de arquivo ausente, não com mensagem clara.

## Referências

- Aula "AWS SAM e Funções Lambda" — Fase 3, Desenvolvimento Serverless
- Implementação: [`postech-tc3-lambda-auth/infra`](https://github.com/Kc1t/postech-tc3-lambda-auth/tree/main/infra)
