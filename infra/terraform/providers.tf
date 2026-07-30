
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.57.0"
    }
  }
}

provider "aws" {
  # Configuration options
  region = var.aws-region
  profile = var.aws_profile
}