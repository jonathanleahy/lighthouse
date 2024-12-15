provider "aws" {
  region  = var.region
  profile = var.environment
  default_tags {
    tags = {
      Env     = var.environment
      Squad   = "psm-enablement"
      Service = "backoffice-core-bff"
      Pci: "no"
      Repo: "github.com/pismo/backoffice-core-bff"
    }
  }
}

