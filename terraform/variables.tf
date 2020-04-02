variable "name" {
  type        = string
  description = "Name used to identify resources"
}

variable "access_key" {
  type        = string
  description = "AWS access key ID"
}

variable "secret_key" {
  type        = string
  description = "AWS secret access key"
}

variable "region" {
  type        = string
  description = "AWS region to deploy where resources will be deployed"
  default     = "us-east-2"
}

variable "timescaledb_password" {
  type        = string
  description = "Password for the TimescaleDB user"
}

variable "rds_username" {
  type        = string
  description = "Username for the master RDS user"
}

variable "rds_password" {
  type        = string
  description = "Password for the master RDS user"
}

variable "patches_container_tag" {
  type        = string
  description = "Tag of the Docker container to be used in the patches container definition"
  default     = "latest"
}

variable "ether_container_tag" {
  type        = string
  description = "Tag of the Docker container to be used in the ether container definition"
  default     = "latest"
}

variable "karen_container_tag" {
  type        = string
  description = "Tag of the Docker container to be used in the karen container definition"
  default     = "latest"
}

variable "heimdall_container_tag" {
  type        = string
  description = "Tag of the Docker container to be used in the heimdall container definition"
  default     = "latest"
}

variable "private_key_cert" {
  type        = string
  description = "Local path to the private RSA key for the TLS certificate for heimdall"
}

variable "cert" {
  type        = string
  description = "Local path to the TLS certificate for heimdall"
}
