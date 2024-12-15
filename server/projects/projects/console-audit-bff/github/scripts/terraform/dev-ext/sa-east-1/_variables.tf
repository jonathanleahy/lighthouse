# region configuration
variable "region" {
  default = "sa-east-1"
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "dev-ext"
  description = "Define aws enviroment"
}

variable "eks_cluster_name" {
  default = "eks-sa-east-1-20210521"
  description = "Define aws eks_cluster_name"
}
