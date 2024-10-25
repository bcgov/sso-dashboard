provider "aws" {
  region = var.region
}

terraform {
  backend "s3" {
    bucket = var.state_bucket
    key    = "tf-state"
    region = var.region
  }
}
