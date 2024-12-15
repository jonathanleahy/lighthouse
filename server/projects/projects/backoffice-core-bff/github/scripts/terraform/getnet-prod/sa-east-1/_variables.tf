# region configuration
variable "region" {
  default = "sa-east-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "getnet-prod"
  description = "aws account configuration"
}

variable "eks_cluster_name" {
  default = "getnet-prod-sae1"
  description = "eks cluster name"
}

variable "namespace" {
  default = "psm-crm"
  description = "namespace team name"
}
