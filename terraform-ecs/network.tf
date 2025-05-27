# Select pre-existing networking config into data for use in resources.
data "aws_vpc" "selected" {
  state = "available"
}

data "aws_subnet" "subnet_a" {
  filter {
    name   = "tag:Name"
    values = [var.subnet_a]
  }
}

data "aws_subnet" "subnet_b" {
  filter {
    name   = "tag:Name"
    values = [var.subnet_b]
  }
}

# Open for api gateway to receive traffic from internet, e.g our openshift alloy. Authorization token is checked on all calls
resource "aws_security_group" "loki_sg" {
  name        = "loki_sg"
  description = "Security group for loki"
  vpc_id      = data.aws_vpc.selected.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "loki_lb" {
  name                             = "loki-lb"
  internal                         = true
  load_balancer_type               = "application"
  security_groups                  = [aws_security_group.loki_sg.id]
  subnets                          = [data.aws_subnet.subnet_a.id, data.aws_subnet.subnet_b.id]
  enable_cross_zone_load_balancing = true
}

resource "aws_lb" "loki_gossip_lb" {
  name                             = "loki-gossip-lb"
  internal                         = true
  load_balancer_type               = "network"
  security_groups                  = [aws_security_group.loki_sg.id]
  subnets                          = [data.aws_subnet.subnet_a.id, data.aws_subnet.subnet_b.id]
  enable_cross_zone_load_balancing = true
}

resource "aws_lb_listener" "loki_gossip_listener" {
  load_balancer_arn = aws_lb.loki_gossip_lb.arn
  port              = "7946"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.loki_target_group_gossip.arn
  }

  depends_on = [aws_lb.loki_gossip_lb]
}

resource "aws_lb_listener" "loki_listener" {
  load_balancer_arn = aws_lb.loki_lb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.loki_target_group_read.arn
  }

  depends_on = [aws_lb.loki_lb]
}

# Each ecs container set needs its rules
resource "aws_lb_listener_rule" "write_service_rule" {
  listener_arn = aws_lb_listener.loki_listener.arn
  priority     = 98

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.loki_target_group_write.arn
  }

  condition {
    path_pattern {
      values = [
        "/api/prom/push",
        "/loki/api/v1/push",
        "/distributor/ring",
        "/flush",
        "/ingester*",
      ]
    }
  }
}

resource "aws_lb_listener_rule" "write_service_rule_2" {
  listener_arn = aws_lb_listener.loki_listener.arn
  priority     = 97

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.loki_target_group_write.arn
  }

  condition {
    path_pattern {
      values = [
        "/ring",
        "/memberlist",
      ]
    }
  }
}

# Unmatched paths all go to the read service
resource "aws_lb_listener_rule" "read_service_rule" {
  listener_arn = aws_lb_listener.loki_listener.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.loki_target_group_read.arn
  }

  condition {
    path_pattern {
      values = [
        "/*",
      ]
    }
  }
}

resource "aws_lb_target_group" "loki_target_group_read" {
  name        = "loki-target-group-read"
  port        = 3100
  protocol    = "HTTP"
  vpc_id      = data.aws_vpc.selected.id
  target_type = "ip"

  health_check {
    path                = "/ready"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

resource "aws_lb_target_group" "loki_target_group_write" {
  name        = "loki-target-group-write"
  port        = 3100
  protocol    = "HTTP"
  vpc_id      = data.aws_vpc.selected.id
  target_type = "ip"

  health_check {
    path                = "/ready"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

resource "aws_lb_target_group" "loki_target_group_gossip" {
  name        = "loki-target-group-gossip"
  port        = 7946
  protocol    = "TCP"
  vpc_id      = data.aws_vpc.selected.id
  target_type = "ip"
}
