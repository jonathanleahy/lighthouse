# region configuration
variable "region" {
  default = "ap-south-1"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "ind-prod"
}

variable "eks_cluster_name" {
  default = "ind-prod-ap-south-1-20211216"
  description = "eks cluster name"
}
