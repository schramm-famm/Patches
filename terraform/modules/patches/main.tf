data "aws_region" "patches" {}

resource "aws_cloudwatch_log_group" "patches" {
  name              = "${var.name}_patches"
  retention_in_days = 1
}

resource "aws_ecs_task_definition" "patches" {
  family       = "${var.name}_patches"
  network_mode = "bridge"

  container_definitions = <<EOF
[
  {
    "name": "${var.name}_patches",
    "image": "343660461351.dkr.ecr.us-east-2.amazonaws.com/patches:${var.container_tag}",
    "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
            "awslogs-group": "${aws_cloudwatch_log_group.patches.name}",
            "awslogs-region": "${data.aws_region.patches.name}",
            "awslogs-stream-prefix": "${var.name}"
        }
    },
    "cpu": 10,
    "memory": 128,
    "essential": true,
    "environment": [
        {
            "name": "PATCHES_DB_HOST",
            "value": "${var.db_host}"
        },
        {
            "name": "PATCHES_DB_PORT",
            "value": "${var.db_port}"
        },
        {
            "name": "PATCHES_DB_USERNAME",
            "value": "${var.db_username}"
        },
        {
            "name": "PATCHES_DB_PASSWORD",
            "value": "${var.db_password}"
        },
        {
            "name": "PATCHES_KAFKA_SERVER",
            "value": "${var.kafka_server}"
        },
        {
            "name": "PATCHES_KAFKA_TOPIC",
            "value": "${var.kafka_topic}"
        },
        {
            "name": "PATCHES_HEIMDALL_SERVER",
            "value": "${var.heimdall_endpoint}"
        },
        {
            "name": "PATCHES_ETHER_SERVER",
            "value": "${var.ether_endpoint}"
        }
    ],
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": ${var.port},
        "protocol": "tcp"
      }
    ]
  }
]
EOF
}

resource "aws_elb" "patches" {
  name            = "${var.name}-patches"
  subnets         = var.subnets
  security_groups = var.security_groups
  internal        = var.internal

  listener {
    instance_port     = var.port
    instance_protocol = "tcp"
    lb_port           = 80
    lb_protocol       = "tcp"
  }
}

resource "aws_ecs_service" "patches" {
  name            = "${var.name}_patches"
  cluster         = var.cluster_id
  task_definition = aws_ecs_task_definition.patches.arn

  load_balancer {
    elb_name       = aws_elb.patches.name
    container_name = "${var.name}_patches"
    container_port = 80
  }

  desired_count = 1
}
