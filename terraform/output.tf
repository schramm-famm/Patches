output "patches_host" {
  value = module.patches.elb_dns_name
}

output "heimdall_host" {
  value = module.heimdall.external_lb_dns_name
}
