terraform {
  required_version = ">= 0.14.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.11.0"
    }
    mysql = {
      source  = "terraform-providers/mysql"
      version = "~> 1.9.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 1.13.2"
    }
    datadog = {
      source  = "DataDog/datadog"
      version = "~> 2.16.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.0.1"
    }
  }
}

