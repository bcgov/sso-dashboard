resource "aws_s3_bucket" "sso_loki" {
  bucket = var.bucket_name
}
