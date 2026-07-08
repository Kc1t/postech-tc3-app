output "cluster_name" {
  description = "Nome do cluster EKS"
  value       = module.eks.cluster_name
}

output "cluster_endpoint" {
  description = "Endpoint do cluster EKS"
  value       = module.eks.cluster_endpoint
}

output "kubeconfig_command" {
  description = "Comando para configurar o kubectl"
  value       = "aws eks update-kubeconfig --name ${module.eks.cluster_name} --region ${var.aws_region}"
}

output "db_endpoint" {
  description = "Endpoint do RDS PostgreSQL"
  value       = aws_db_instance.workshop.endpoint
}

output "postgres_dsn" {
  description = "DSN completo para o POSTGRES_DSN"
  value       = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.workshop.endpoint}/${var.db_name}?sslmode=require"
  sensitive   = true
}
