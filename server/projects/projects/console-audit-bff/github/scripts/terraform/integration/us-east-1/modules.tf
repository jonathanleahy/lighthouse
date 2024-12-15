module "aws_static_parameters" {
  source = "github.com/pismo/tf-module-aws-static-parameters?ref=1.34.0"

  region           = data.aws_region.current.name
  eks_cluster_name = var.eks_cluster_name
}
