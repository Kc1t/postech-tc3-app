module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "workshop-api"
  cluster_version = "1.36"

  cluster_endpoint_public_access = true

  vpc_id     = var.vpc_id
  subnet_ids = data.aws_subnets.eks_vpc.ids

  eks_managed_node_groups = {
    workers = {
      instance_types = [var.node_instance_type]
      min_size       = var.node_min
      max_size       = var.node_max
      desired_size   = var.node_desired
    }
  }

  tags = {
    Project = var.cluster_name
  }
}
