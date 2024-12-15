# region configuration
variable "region" {
  type    = string
  default = "ap-south-22"
  description = "aws region"
}

# aws account configuration
variable "account" {
  type    = string
  default = "ind-prod"
  description = "account"
}

