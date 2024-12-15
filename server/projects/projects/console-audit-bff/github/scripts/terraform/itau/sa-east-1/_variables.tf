# region configuration
variable "region" {
  default = "sa-east-1"
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "itau-prod"
  description = "Define aws account"
}

variable "eks_cluster_name" {
  default = "itau-sa-east-1-20220223"
  description = "Define aws eks_cluster_name"
}
