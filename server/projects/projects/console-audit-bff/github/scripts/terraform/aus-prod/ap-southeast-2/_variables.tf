# region configuration
variable "region" {
  default = "ap-southeast-2"
  type        = string
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "aus-prod"
  type        = string
  description = "Define aws account"
}

variable "eks_cluster_name" {
  default = "aus-prod-ap-southeast-2-20240319"
  type        = string
  description = "Define aws cluster name"
}
