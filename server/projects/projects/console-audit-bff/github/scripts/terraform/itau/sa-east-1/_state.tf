# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-state-pismo"
    key    = "api/console_audit_bff/sa-east-1/756778449919.tfstate"
    region = "sa-east-1"
  }
}
