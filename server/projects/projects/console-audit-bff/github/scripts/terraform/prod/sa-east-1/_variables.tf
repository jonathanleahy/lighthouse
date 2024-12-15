# region configuration
variable "region" {
  default = "sa-east-1"
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "production"
  description = "Define aws enviroment"
}

variable "eks_cluster_name" {
  default = "prod-sa-east-1-20210712"
  description = "Define aws eks_cluster_name"
}
