# region configuration
variable "region" {
  default = "us-east-1"
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "integration"
  description = "Define aws enviroment"
}

variable "eks_cluster_name" {
  default = "integration-us-east-1-20210528"
  description = "Define aws eks_cluster_name"
}
