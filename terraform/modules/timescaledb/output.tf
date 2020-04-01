output "db_host" {
  value = aws_elb.timescaledb.dns_name
}

output "db_port" {
  value = 5432
}

output "db_username" {
  value = "postgres"
}
