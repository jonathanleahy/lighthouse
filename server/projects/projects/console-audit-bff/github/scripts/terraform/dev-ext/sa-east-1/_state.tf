# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket  = "terraform-state-pismo"
    key     = "api/console-audit-bff/sa-east-1/dev-ext.tfstate"
    region  = "sa-east-1"
    profile = "dev-ext"
  }
}
