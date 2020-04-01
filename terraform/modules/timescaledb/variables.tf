variable "name" {
  type        = string
  description = "Name used to identify resources"
}

variable "vpc_id" {
  type        = string
  description = "ID of the VPC to run the TimescaleDB cluster in"
}

variable "subnets" {
  type        = list(string)
  description = "VPC subnets for the TimescaleDB cluster and load balancer"
}

variable "ecs_instance_profile_id" {
  type        = string
  description = "ID of the instance profile to attach to the EC2 container instances in the cluster"
}

variable "internal" {
  type        = bool
  description = "Toggle whether the load balancer will be internal"
}

variable "db_password" {
  type        = string
  description = "Port of the TimescaleDB server"
}
