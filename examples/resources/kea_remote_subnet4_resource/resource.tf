resource "kea_remote_subnet4_resource" "example" {
  hostname = "kea-primary.example.com"
  subnet   = "192.168.225.0/24"
  pools = [
    { pool = "192.168.225.50-192.168.225.150" }
  ]
  relay = [
    { ip_address = "192.168.225.1" }
  ]
  option_data = [
    { code = 3, name = "routers", data = "192.168.225.1" },
    { code = 15, name = "domain-name", data = "example.com" },
    { code = 6, name = "domain-name-servers", data = "4.2.2.2, 8.8.8.8", always_send = true },
  ]
  user_context = {
    "foo" = "bar"
  }
  ddns_send_updates            = true
  ddns_override_no_update      = true
  ddns_override_client_update  = true
  ddns_generated_prefix        = "host"
  ddns_qualifying_suffix       = "example.com."
  ddns_rev_dns_name            = "225.168.192.in-addr.arpa."
  ddns_use_conflict_resolution = true
}
