# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket         = "terraform-citi-stag-workload"
    key            = "api/console-audit-bff/us-east-2/terraform.tfstate"
    region         = "us-east-1"
    profile        = "remote-state"
    dynamodb_table = "terraform-citi-stag-workload"
  }
}

