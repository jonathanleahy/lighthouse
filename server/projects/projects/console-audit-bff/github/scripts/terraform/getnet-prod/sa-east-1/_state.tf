# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket         = "terraform-getnet-workload-prod"
    key            = "api/console_audit_bff/sa-east-1/terraform.tfstate"
    region         = "sa-east-1"
    profile        = "remote-state"
  }
}

