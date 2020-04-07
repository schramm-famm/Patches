variable "name" {
  type        = string
  description = "Name used to identify resources"
}

variable "container_tag" {
  type        = string
  description = "Tag of the patches container in the registry to be used"
  default     = "latest"
}

variable "port" {
  type        = number
  description = "The port that patches' container port will map to on the host"
  default     = 80
}

variable "container_count" {
  type        = number
  description = "The number of containers to deploy in the patches service"
  default     = 1
}

variable "cluster_id" {
  type        = string
  description = "ID of the ECS cluster that the patches service will run in"
}

variable "security_groups" {
  type        = list(string)
  description = "VPC security groups for the patches service load balancer"
}

variable "subnets" {
  type        = list(string)
  description = "VPC subnets for the patches service load balancer"
}

variable "internal" {
  type        = bool
  description = "Toggle whether the load balancer will be internal"
}

variable "db_host" {
  type        = string
  description = "Host of the TimescaleDB server"
}

variable "db_port" {
  type        = string
  description = "Port of the TimescaleDB server"
}

variable "db_username" {
  type        = string
  description = "Username for accessing the TimescaleDB server"
}

variable "db_password" {
  type        = string
  description = "Password for accessing the TimescaleDB server"
}

variable "kafka_server" {
  type        = string
  description = "Server where Kafka is running"
}

variable "kafka_topic" {
  type        = string
  description = "Kafka topic to read from"
}

variable "heimdall_endpoint" {
  type        = string
  description = "Endpoint for accessing the heimdall service"
}

variable "ether_endpoint" {
  type        = string
  description = "Endpoint for accessing the ether service"
}
