resource "aws_apigatewayv2_api" "sso_loki_api" {
  name          = "loki-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_vpc_link" "loki_vpc_link" {
  name               = "loki_vpc_link"
  subnet_ids         = [data.aws_subnet.subnet_a.id, data.aws_subnet.subnet_b.id]
  security_group_ids = [aws_security_group.loki_sg.id]
}

resource "aws_apigatewayv2_integration" "sso_loki_api_integration" {
  api_id             = aws_apigatewayv2_api.sso_loki_api.id
  integration_type   = "HTTP_PROXY"
  connection_id      = aws_apigatewayv2_vpc_link.loki_vpc_link.id
  connection_type    = "VPC_LINK"
  integration_method = "ANY"
  integration_uri    = aws_lb_listener.loki_listener.arn
}

resource "aws_apigatewayv2_route" "sso_grafana_route_any" {
  api_id    = aws_apigatewayv2_api.sso_loki_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.sso_loki_api_integration.id}"

  authorization_type = "CUSTOM"
  authorizer_id      = aws_apigatewayv2_authorizer.loki_authorizer.id
}

resource "aws_apigatewayv2_authorizer" "loki_authorizer" {
  api_id                            = aws_apigatewayv2_api.sso_loki_api.id
  authorizer_type                   = "REQUEST"
  enable_simple_responses           = true
  authorizer_uri                    = aws_lambda_function.auth_function.invoke_arn
  identity_sources                  = ["$request.header.Authorization"]
  name                              = "loki-authorizer"
  authorizer_payload_format_version = "2.0"
}

resource "aws_apigatewayv2_stage" "sso_grafana_api_default_stage" {
  api_id      = aws_apigatewayv2_api.sso_loki_api.id
  name        = "$default"
  auto_deploy = true
}
