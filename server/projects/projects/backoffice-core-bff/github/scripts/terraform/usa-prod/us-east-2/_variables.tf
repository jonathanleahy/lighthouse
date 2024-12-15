# region configuration
variable "region" {
  default = "us-east-2"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "usa-prod"
}

variable "eks_cluster_name" {
  default = "prod-us-east-2-20230713"
  description = "eks cluster name"
}

variable "namespace" {
  default = "psm-enablement"
}
