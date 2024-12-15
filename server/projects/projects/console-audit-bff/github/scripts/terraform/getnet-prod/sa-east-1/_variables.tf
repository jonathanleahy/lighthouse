# region configuration
variable "region" {
  default = "sa-east-1"
  type        = string
  description = "Define aws region"
}

# aws account configuration
variable "account" {
  default = "getnet-prod"
  type        = string
  description = "Define aws account"
}

variable "eks_cluster_name" {
  default = "getnet-prod-sae1"
  type        = string
  description = "Define aws cluster name"
}
