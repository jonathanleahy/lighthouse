# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-irl-workload-prod"
    key    = "api/backoffice_core_bff/eu-west-1/terraform.tfstate"
    region = "eu-west-1"
    profile        = "remote-state"
    dynamodb_table = "terraform-irl-workload-prod"
  }
}
