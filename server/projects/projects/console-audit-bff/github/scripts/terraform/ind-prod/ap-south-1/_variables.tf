# region configuration
variable "region" {
  default = "ap-south-1"
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "ind-prod"
  description = "Define aws account"
}

variable "eks_cluster_name" {
  default = "ind-prod-ap-south-1-20211216"
  description = "Define aws eks_cluster_name"
}
