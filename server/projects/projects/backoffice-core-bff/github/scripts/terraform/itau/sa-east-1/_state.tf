# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-state-pismo"
    key    = "api/backoffice_core_bff/sa-east-1/756778449919.tfstate"
    region = "sa-east-1"
  }
}
