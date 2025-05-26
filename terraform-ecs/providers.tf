provider "aws" {
  region = var.region

  default_tags {
    tags = {
      repository = "sso-dashboard"
    }
  }
}

terraform {
  backend "s3" {}
}
