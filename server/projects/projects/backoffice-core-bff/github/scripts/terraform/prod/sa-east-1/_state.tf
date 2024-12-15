# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-state-pismo"
    key    = "api/backoffice_core_bff/sa-east-1/408082092235.tfstate"
    region = "sa-east-1"
  }
}
