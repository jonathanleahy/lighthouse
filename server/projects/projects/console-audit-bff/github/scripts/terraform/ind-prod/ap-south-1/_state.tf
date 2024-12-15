# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket         = "terraform-ind-workload-prod"
    key            = "api/console_audit_bff/ap-south-1/terraform.tfstate"
    region         = "ap-south-1"
    profile        = "remote-state"
    dynamodb_table = "terraform-ind-workload-prod"
  }
}

