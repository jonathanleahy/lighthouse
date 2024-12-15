# region configuration
variable "region" {
  default = "ap-southeast-2"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "aus-prod"
  description = "aws account configuration"
}

variable "eks_cluster_name" {
  default = "aus-prod-ap-southeast-2-20240319"
  description = "eks cluster name"
}

variable "namespace" {
  default = "psm-crm"
  description = "namespace team name"
}
