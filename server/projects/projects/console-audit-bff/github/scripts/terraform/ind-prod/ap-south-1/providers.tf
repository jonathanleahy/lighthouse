provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Squad   = "psm-console"
      Service = "console-audit-bff"
    }
  }
}

