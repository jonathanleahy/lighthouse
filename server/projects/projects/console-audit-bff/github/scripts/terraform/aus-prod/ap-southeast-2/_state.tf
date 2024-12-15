# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket = "terraform-aus-workload-prod"
    key    = "api/console_audit_bff/ap-southeast-2/terraform.tfstate"
    region = "ap-southeast-2"
    profile = "aus-prod"
  }
}
