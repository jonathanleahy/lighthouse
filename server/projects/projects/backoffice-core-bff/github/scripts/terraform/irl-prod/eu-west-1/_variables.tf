# region configuration
variable "region" {
  default = "eu-west-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "irl-prod"
  description = "aws account configuration"
}

variable "eks_cluster_name" {
  default = "irl-prod-eu-west-1-20240626"
  description = "eks cluster name"
}

variable "namespace" {
  default = "psm-crm"
  description = "namespace team name"
}
