provider "aws" {
  default_tags {
    tags = {
      Squad   = "psm-crm"
      Service = "crm-core-bff"
    }
  }
}

