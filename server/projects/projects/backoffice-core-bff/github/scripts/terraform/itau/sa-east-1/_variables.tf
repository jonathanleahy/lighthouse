# region configuration
variable "region" {
  default = "sa-east-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "itau-prod"
}

variable "eks_cluster_name" {
  default = "itau-sa-east-1-20220223"
  description = "eks cluster name"
}
