locals {
  ecs_services = ["loki-write-service", "loki-read-service"]
}


resource "aws_sns_topic" "rocket_chat" {
  name = "rocketchat"
}

resource "aws_sns_topic_subscription" "rocket_chat_subscription" {
  topic_arn = aws_sns_topic.rocket_chat.arn
  protocol  = "https"
  endpoint  = var.rc_webhook
}

resource "aws_cloudwatch_metric_alarm" "loki_tasks_low" {
  for_each            = toset(local.ecs_services)
  alarm_name          = "${var.rc_prefix}: ${each.key} tasks low"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  threshold           = 0
  alarm_description   = "Alarm if DesiredTaskCount is greater than the RunningTaskCount for the ECS service ${each.key}."

  metric_query {
    id = "running_task_count"
    metric {
      metric_name = "RunningTaskCount"
      namespace   = "ECS/ContainerInsights"
      period      = 300 # 5 minute
      stat        = "Average"
      dimensions = {
        ClusterName = "loki-cluster"
        ServiceName = each.key
      }
    }
  }

  metric_query {
    id = "desired_task_count"
    metric {
      metric_name = "DesiredTaskCount"
      namespace   = "ECS/ContainerInsights"
      period      = 300 # 5 minute
      stat        = "Average"
      dimensions = {
        ClusterName = "loki-cluster"
        ServiceName = each.key
      }
    }
  }

  metric_query {
    id          = "task_deficit"
    expression  = "desired_task_count - running_task_count"
    label       = "Difference between Desired and Running Task Counts"
    return_data = true
  }

  alarm_actions = [aws_sns_topic.rocket_chat.arn]
}
