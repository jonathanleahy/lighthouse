# region configuration
variable "region" {
  default = "sa-east-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "dev-ext"
}

variable "eks_cluster_name" {
  default = "eks-sa-east-1-20210521"
  description = "eks cluster name"
}
