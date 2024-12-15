# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-us-workload-prod"
    key    = "api/backoffice_core_bff/us-east-2/terraform.tfstate"
    region = "us-east-2"
    profile        = "remote-state"
    dynamodb_table = "terraform-us-workload-prod"
  }
}
