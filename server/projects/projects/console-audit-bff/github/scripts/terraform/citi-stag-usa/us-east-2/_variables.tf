# region configuration
variable "region" {
  default = "us-east-2"
  description = "Define aws region"
}

# AWS account configuration
variable "account" {
  default = "citi-stag-usa"
  description = "Define aws enviroment"
}

variable "eks_cluster_name" {
  default = "citi-stag-us-east-2-20221125"
  description = "Define aws account_id"
}
