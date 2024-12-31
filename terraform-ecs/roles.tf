data "aws_iam_policy" "permissions_boundary_policy" {
    name = "sso-dashboard-boundary"
}

# Execution role, permissions to log to cloudwatch
resource "aws_iam_role" "loki_execution_role" {
  name = "loki-execution-role"
  permissions_boundary = data.aws_iam_policy.permissions_boundary_policy.arn
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "loki_execution_policy" {
  name        = "loki-execution-policy"
  permissions_boundary = data.aws_iam_policy.permissions_boundary_policy.arn
  description = "Permissions for ECS task execution, including logging to CloudWatch"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "loki_execution_role_attachment" {
  role       = aws_iam_role.loki_execution_role.name
  policy_arn = aws_iam_policy.loki_execution_policy.arn
}

# Task role for running container, needs to connect to the s3 bucket for logs
resource "aws_iam_role" "loki_task_role" {
  name = "loki-task-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "loki_task_s3_policy" {
  name        = "loki-task-policy-s3"
  description = "Permissions for Loki to access S3"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:*"
        ]
        Resource = [
          "arn:aws:s3:::${aws_s3_bucket.sso_loki.bucket}",
          "arn:aws:s3:::${aws_s3_bucket.sso_loki.bucket}/*",
        ]
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "loki_task_s3_role_attachment" {
  role       = aws_iam_role.loki_task_role.name
  policy_arn = aws_iam_policy.loki_task_s3_policy.arn
}

# Allow api gateway to use lambdas
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.auth_function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.sso_loki_api.execution_arn}/*/*"
}

# Below permission can be added temporarily if needing to ssh into the loki tasks for debugging.

# resource "aws_iam_policy" "loki_task_ssh_policy" {
#   name        = "loki-task-policy-efs"
#   description = "Permissions to ssh to loki tasks"
#   policy      = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect = "Allow"
#         Action = [
#           "ssmmessages:CreateControlChannel",
#           "ssmmessages:CreateDataChannel",
#           "ssmmessages:OpenControlChannel",
#           "ssmmessages:OpenDataChannel"
#         ]
#         Resource = "*"
#       },
#     ]
#   })
# }

# resource "aws_iam_role_policy_attachment" "loki_task_efs_role_attachment" {
#   role       = aws_iam_role.loki_task_role.name
#   policy_arn = aws_iam_policy.loki_task_ssh_policy.arn
# }
