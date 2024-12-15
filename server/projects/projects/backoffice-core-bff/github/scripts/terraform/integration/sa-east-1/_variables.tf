# region configuration
variable "region" {
  default = "sa-east-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "integration"
}

variable "eks_cluster_name" {
  default = "integration-sa-east-1-20210528"
  description = "eks cluster name"
}
