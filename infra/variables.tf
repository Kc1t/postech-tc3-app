variable "aws_region" {
  description = "Região AWS"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "Nome do cluster EKS"
  type        = string
  default     = "workshop-api"
}

variable "node_instance_type" {
  description = "Tipo de instância dos worker nodes"
  type        = string
  default     = "t3.medium"
}

variable "node_min" {
  description = "Mínimo de worker nodes"
  type        = number
  default     = 1
}

variable "node_max" {
  description = "Máximo de worker nodes"
  type        = number
  default     = 4
}

variable "node_desired" {
  description = "Número desejado de worker nodes"
  type        = number
  default     = 2
}

variable "db_name" {
  description = "Nome do banco de dados"
  type        = string
  default     = "workshop"
}

variable "db_username" {
  description = "Usuário do banco de dados"
  type        = string
  default     = "workshop"
}

variable "db_password" {
  description = "Senha do banco de dados (sensível)"
  type        = string
  sensitive   = true
}

variable "vpc_id" {
  description = "ID da VPC onde o EKS está rodando"
  type        = string
  default     = "vpc-039d2bdb052aea5a8"
}
