# region configuration
variable "region" {
  default = "ap-south-2"
  description = "aws region"
}

# aws account configuration
variable "account" {
  default = "ind-prod"
}

variable "eks_cluster_name" {
  default = "ind-prod-ap-south-2-20230607"
  description = "eks cluster name"
}
