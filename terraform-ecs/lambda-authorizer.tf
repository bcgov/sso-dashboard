resource "aws_lambda_function" "auth_function" {
  function_name = "lokiApiAuth"
  runtime       = "nodejs20.x"  # or your preferred runtime
  role          = aws_iam_role.lambda_exec.arn
  handler       = "authorize.handler"
  filename      = "loki-authorizer.zip"

  source_code_hash = filebase64sha256("./loki-authorizer.zip")

  environment {
    variables = {
        AUTH_SECRET = var.auth_secret
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Effect = "Allow"
        Sid    = ""
      }
    ]
  })
}
