# region configuration
variable "region" {
  type = string
  default = "ap-southeast-3"
  description = "aws region"
}

# aws account configuration
variable "account" {
  type = string
  default = "aus-prod"
  description = "account"
}

