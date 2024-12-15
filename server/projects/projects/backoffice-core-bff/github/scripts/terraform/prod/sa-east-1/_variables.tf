# region configuration
variable "region" {
  default = "sa-east-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "production"
}

variable "eks_cluster_name" {
  default = "prod-sa-east-1-20210712"
  description = "eks cluster name"
}
