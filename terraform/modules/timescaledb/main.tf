data "aws_region" "timescaledb" {}

/* ECS CLUSTER SETUP */

resource aws_security_group "timescaledb" {
  name        = "${var.name}_timescaledb"
  description = "Allow traffic into the TimescaleDB port"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

module "ecs_cluster" {
  source                  = "github.com/schramm-famm/bespin//modules/ecs_cluster"
  name                    = "${var.name}-timescaledb"
  security_group_ids      = [aws_security_group.timescaledb.id]
  subnets                 = var.subnets
  ec2_instance_profile_id = var.ecs_instance_profile_id
  enable_efs              = false
}

/* ECS SERVICE SETUP */

resource "aws_cloudwatch_log_group" "timescaledb" {
  name              = "${var.name}_timescaledb"
  retention_in_days = 1
}

resource "aws_ecs_task_definition" "timescaledb" {
  family       = "${var.name}_timescaledb"
  network_mode = "bridge"

  container_definitions = <<EOF
[
  {
    "name": "${var.name}_timescaledb",
    "image": "timescale/timescaledb:latest",
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-group": "${aws_cloudwatch_log_group.timescaledb.name}",
        "awslogs-region": "${data.aws_region.timescaledb.name}",
        "awslogs-stream-prefix": "${var.name}"
      }
    },
    "cpu": 512,
    "memory": 512,
    "essential": true,
    "environment": [
      {
        "name": "POSTGRES_PASSWORD",
        "value": "${var.db_password}"
      }
    ],
    "portMappings": [
      {
        "containerPort": 5432,
        "hostPort": 5432,
        "protocol": "tcp"
      }
    ]
  }
]
EOF
}

resource "aws_elb" "timescaledb" {
  name            = "${var.name}-timescaledb"
  subnets         = var.subnets
  security_groups = [aws_security_group.timescaledb.id]
  internal        = var.internal

  listener {
    instance_port     = 5432
    instance_protocol = "tcp"
    lb_port           = 5432
    lb_protocol       = "tcp"
  }
}

resource "aws_ecs_service" "timescaledb" {
  name            = "${var.name}_timescaledb"
  cluster         = module.ecs_cluster.cluster_id
  task_definition = aws_ecs_task_definition.timescaledb.arn

  load_balancer {
    elb_name       = aws_elb.timescaledb.name
    container_name = "${var.name}_timescaledb"
    container_port = 5432
  }

  desired_count = 1
}
