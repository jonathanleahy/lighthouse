# region configuration
variable "region" {
  default = "eu-west-1"
  type        = string
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "irl-prod"
  type        = string
  description = "Define aws account"
}

variable "eks_cluster_name" {
  default = "irl-prod-eu-west-1-20240626"
  type        = string
  description = "Define aws cluster name"
}
