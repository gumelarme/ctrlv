terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }

    random = {
      source  = "hashicorp/random"
      version = "3.1.2"
    }
  }

  required_version = ">= 0.14.9"
}

locals {
  region             = "ap-northeast-1"
  unique_id          = "${random_pet.ctrlv.id}-${random_string.ctrlv.id}"
  dynamodb_table     = "ctrlv-db-${local.unique_id}"
  dynamodb_table_gsi = "Alias-Id-Index-${local.unique_id}"
}

resource "random_pet" "ctrlv" {
  separator = "-"
}

resource "random_string" "ctrlv" {
  length  = 6
  lower   = false
  special = false
}

provider "aws" {
  profile = "ctrlv"
  region  = local.region
}
