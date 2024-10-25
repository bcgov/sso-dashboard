resource "aws_ecs_cluster" "sso_ecs_cluster" {
  name = "loki-cluster"
}

resource "aws_ecs_task_definition" "loki_write" {
  family                   = "loki-task"
  execution_role_arn       = aws_iam_role.loki_execution_role.arn
  task_role_arn            = aws_iam_role.loki_task_role.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([{
    name = "loki-write"
    image = "jelanglois/grafana:flush"
    essential              = true
    memory                 = var.loki_write_memory
    cpu                    = var.loki_write_cpu
    # IMPORTANT: Make sure ingesters have time to cut any chunks in memory.
    stop_timeout = 120

    portMappings = [
      {
        name          = "app"
        containerPort = 3100
        hostPort      = 3100
      },
      {
        name          = "gossip"
        containerPort = 7946
        hostPort      = 7946
        protocol      = "tcp"
      },
      {
        name        = "grpc"
        hostPort    = 9095
        protocol    = "tcp"
        containerPort = 9095
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-create-group  = "true"
        awslogs-group         = "/ecs/${aws_ecs_cluster.sso_ecs_cluster.name}"
        awslogs-region        = "ca-central-1"
        awslogs-stream-prefix = "ecs-loki-write"
      }
    }
    environment = [
      {
        name  = "S3_BUCKETNAME"
        value = aws_s3_bucket.sso_loki.bucket
      },
      {
        name  = "S3_REGION"
        value = var.region
      },
      {
        name  = "S3_ENDPOINT"
        value = "s3.${var.region}.amazonaws.com"
      },
      {
        name  = "JOIN_MEMBERS"
        value = "${aws_lb.loki_gossip_lb.dns_name}:7946"
      },
    ]
    command = [
      "-target=write",
      "-config.file=/etc/loki/config/config.yaml",
      "-config.expand-env=true",
    ]
  }])
}

resource "aws_ecs_service" "loki_write" {
  name                   = "loki-write-service"
  cluster                = aws_ecs_cluster.sso_ecs_cluster.id
  task_definition        = aws_ecs_task_definition.loki_write.arn
  desired_count          = 2
  launch_type            = "FARGATE"
  enable_execute_command = true

  load_balancer {
    target_group_arn = aws_lb_target_group.loki_target_group_write.id
    container_name   = "loki-write"
    container_port   = 3100
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.loki_target_group_gossip.id
    container_name   = "loki-write"
    container_port   = 7946
  }

  network_configuration {
    subnets          = [data.aws_subnet.subnet_a.id, data.aws_subnet.subnet_b.id]
    security_groups  = [aws_security_group.loki_sg.id]
    assign_public_ip = false
  }
}

resource "aws_ecs_task_definition" "loki_read" {
  family                   = "loki-task"
  execution_role_arn       = aws_iam_role.loki_execution_role.arn
  task_role_arn            = aws_iam_role.loki_task_role.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([{
    name = "loki-read"
    image                  = "jelanglois/grafana:flush"
    essential              = true
    memory                 = var.loki_read_memory
    cpu                    = var.loki_read_cpu

    portMappings = [
      {
        containerPort = 3100
        hostPort      = 3100
      },
      {
        containerPort = 8500
        hostPort      = 8500
      },
      {
        name        = "grpc"
        hostPort    = 9095
        protocol    = "tcp"
        containerPort = 9095
      },
      {
        containerPort = 7946
        hostPort      = 7946
        protocol : "tcp"
      },
    ]

    environment = [
      {
        name  = "S3_BUCKETNAME"
        value = aws_s3_bucket.sso_loki.bucket
      },
      {
        name  = "S3_REGION"
        value = var.region
      },
      {
        name  = "S3_ENDPOINT"
        value = "s3.${var.region}.amazonaws.com"
      },
      {
        name  = "JOIN_MEMBERS"
        value = "${aws_lb.loki_gossip_lb.dns_name}:7946"
      },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-create-group  = "true"
        awslogs-group         = "/ecs/${aws_ecs_cluster.sso_ecs_cluster.name}"
        awslogs-region        = var.region
        awslogs-stream-prefix = "ecs-loki-read"
      }
    }
    command = [
      "-target=read",
      "-config.file=/etc/loki/config/config.yaml",
      "-config.expand-env=true",
    ]
  }])
}

resource "aws_ecs_service" "loki_read" {
  name                   = "loki-read-service"
  cluster                = aws_ecs_cluster.sso_ecs_cluster.id
  task_definition        = aws_ecs_task_definition.loki_read.arn
  desired_count          = 2
  launch_type            = "FARGATE"
  enable_execute_command = true

  load_balancer {
    target_group_arn = aws_lb_target_group.loki_target_group_read.id
    container_name   = "loki-read"
    container_port   = 3100
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.loki_target_group_gossip.id
    container_name   = "loki-read"
    container_port   = 7946
  }

  network_configuration {
    subnets          = [data.aws_subnet.subnet_a.id, data.aws_subnet.subnet_b.id]
    security_groups  = [aws_security_group.loki_sg.id]
    assign_public_ip = false
  }
}

resource "aws_appautoscaling_target" "ecs_read_service_target" {
  max_capacity       = 6
  min_capacity       = 2
  resource_id        = "service/${aws_ecs_cluster.sso_ecs_cluster.name}/${aws_ecs_service.loki_read.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "scale_out" {
  name                   = "scale_out"
  policy_type           = "StepScaling"
  resource_id           = aws_appautoscaling_target.ecs_read_service_target.id
  scalable_dimension     = "ecs:service:DesiredCount"
  service_namespace      = "ecs"

  step_scaling_policy_configuration {
    adjustment_type = "ChangeInCapacity"
    step_adjustment {
      scaling_adjustment = 1  # Increase the number of tasks by
      # The bounds are for the difference between the trigger and the actual value
      # e.g if alarm is at 60%, and you want a step for 60% - 80%, the lower bound would be zero and the upper 20.
      metric_interval_lower_bound = 0
      metric_interval_upper_bound = 20
    }
    step_adjustment {
      scaling_adjustment = 2
      metric_interval_lower_bound = 20
    }
    cooldown = 30
  }
}

resource "aws_cloudwatch_metric_alarm" "cpu_high" {
  alarm_name                = "HighCpuAlarm"
  comparison_operator       = "GreaterThanThreshold"
  evaluation_periods        = 1
  metric_name              = "CPUUtilization"
  namespace                = "AWS/ECS"
  period                   = 30
  statistic                = "Average"
  threshold                = 60
  dimensions               = {
    ClusterName = aws_ecs_cluster.sso_ecs_cluster.name
    ServiceName = aws_ecs_service.loki_read.name
  }

  alarm_actions = [
    aws_appautoscaling_policy.scale_out.arn,
  ]
}
