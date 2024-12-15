# region configuration
variable "region" {
  default = "us-east-2"
  description = "Defines region"
}

# aws account configuration
variable "account" {
  default = "usa-prod"
  description = "Defines account"
}

variable "eks_cluster_name" {
  default = "prod-us-east-2-20230713"
   description = "Defines eks cluster name"
}
