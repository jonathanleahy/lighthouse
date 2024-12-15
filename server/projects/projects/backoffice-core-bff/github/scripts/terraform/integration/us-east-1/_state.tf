# s3 bucket to persist state
terraform {
  backend "s3" {
    bucket         = "terraform-pismo-workload-integration"
    key            = "api/backoffice-core-bff/us-east-1/terraform.tfstate"
    region         = "us-east-1"
    profile        = "remote-state"
    dynamodb_table = "terraform-pismo-workload-integration"
  }
}
