provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}

