# Infraestrutura como Código (Terraform)

Scripts Terraform que provisionam a infraestrutura da **Workshop API** na AWS: um cluster **EKS** (Kubernetes) e um banco **RDS PostgreSQL 16**. Depois de aplicados, os manifestos de [`../k8s`](../k8s) são implantados no cluster (manualmente ou pela pipeline de CI/CD).

## Recursos criados

| Arquivo | Recurso | Descrição |
|---------|---------|-----------|
| `versions.tf` | `terraform` / `provider "aws"` | Versão mínima do Terraform (≥ 1.5) e provider AWS `~> 5.0`, região via `var.aws_region`. |
| `eks.tf` | `module.eks` (`terraform-aws-modules/eks ~> 19.0`) | Cluster EKS `workshop-api` (Kubernetes 1.36), endpoint público, managed node group com `var.node_instance_type` e min/max/desired de nós. |
| `rds.tf` | `data.aws_subnets.eks_vpc` | Descobre as subnets da VPC informada em `var.vpc_id`. |
| `rds.tf` | `aws_db_subnet_group.workshop` | Subnet group do RDS usando as subnets da VPC. |
| `rds.tf` | `aws_security_group.rds` | Security group liberando a porta `5432` apenas dentro da VPC (`10.0.0.0/8`, `172.31.0.0/16`). |
| `rds.tf` | `aws_db_instance.workshop` | Instância RDS PostgreSQL 16 (`db.t3.micro`, 20 GB), banco/usuário/senha via variáveis. |
| `outputs.tf` | outputs | `cluster_name`, `cluster_endpoint`, `kubeconfig_command`, `db_endpoint` e `postgres_dsn` (sensível). |

> O cluster reutiliza uma **VPC existente** (`var.vpc_id`) e suas subnets, em vez de criar uma nova — simplifica o provisionamento no ambiente da FIAP.

## Variáveis

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `aws_region` | `us-east-1` | Região AWS. |
| `cluster_name` | `workshop-api` | Nome do cluster EKS e prefixo dos recursos. |
| `node_instance_type` | `t3.medium` | Tipo das instâncias dos worker nodes. |
| `node_min` | `1` | Mínimo de worker nodes. |
| `node_max` | `4` | Máximo de worker nodes. |
| `node_desired` | `2` | Número desejado de worker nodes. |
| `db_name` | `workshop` | Nome do banco de dados. |
| `db_username` | `workshop` | Usuário do banco de dados. |
| `db_password` | — (**obrigatória**, sensível) | Senha do banco de dados. |
| `vpc_id` | `vpc-039d2bdb052aea5a8` | VPC onde o EKS e o RDS serão criados. |

## Pré-requisitos

- Terraform ≥ 1.5
- AWS CLI configurado (`aws configure`) com credenciais válidas
- `kubectl` para operar o cluster após o provisionamento
- Uma VPC existente na conta (ou ajustar `var.vpc_id`)

## Como aplicar

```bash
cd infra

# 1. Inicializa providers e módulos
terraform init

# 2. Revisa o plano (informe a senha do banco)
terraform plan  -var="db_password=<senha-forte>"

# 3. Aplica a infraestrutura
terraform apply -var="db_password=<senha-forte>"
```

Após o `apply`, configure o kubectl e recupere o DSN do banco:

```bash
# Comando pronto no output kubeconfig_command
aws eks update-kubeconfig --name workshop-api --region us-east-1

# DSN para usar no Secret do Kubernetes (output sensível)
terraform output -raw postgres_dsn
```

Em seguida, aplique os manifestos do Kubernetes (ver [`../k8s/README.md`](../k8s/README.md)).

## Como destruir

```bash
terraform destroy -var="db_password=<senha-forte>"
```

> ⚠️ O RDS está com `skip_final_snapshot = true` e `deletion_protection = false` (adequado para o ambiente acadêmico). O `destroy` apaga o banco sem snapshot final.

## Observações

- **Estado do Terraform:** este projeto usa backend local por padrão. Para trabalho em equipe, recomenda-se um backend remoto (ex.: S3 + DynamoDB lock).
- **Segredos:** nunca commitar `terraform.tfvars` com a senha real. Passe `db_password` via `-var`, variável de ambiente `TF_VAR_db_password` ou um `.tfvars` fora do versionamento.
</content>
