---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kea_remote_subnet4_resource Resource - terraform-provider-kea"
subcategory: ""
description: |-
  Remote Subnet4 resource
---

# kea_remote_subnet4_resource (Resource)

Remote Subnet4 resource

## Example Usage

```terraform
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
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `hostname` (String) Hostname of the kea server to connect to. e.g. `kea.example.com`
- `pools` (Attributes List) List of pools to configure in the subnet. e.g. `['192.168.230.10-192.168.230.200'] (see [below for nested schema](#nestedatt--pools))
- `subnet` (String) Subnet4 prefix to configure in Kea. e.g. `192.168.230.0/24`

### Optional

- `boot_file_name` (String) Optional conveys the boot configuration file, can be up to 128 bytes long, and is sent using the `file` field.
- `next_server` (String) Optional TFTP boot server IP address, packets sent in the `siaddr` field.
- `option_data` (Attributes List) List of option-data to configure on the pool. e.g. `[{code = 6, name = "domain-name-servers", data = "8.8.8.8, 4.2.2.2"}]` (see [below for nested schema](#nestedatt--option_data))
- `relay` (Attributes List) List of relay IPs to configure in Kea. e.g. `['192.168.230.1']` (see [below for nested schema](#nestedatt--relay))
- `server_hostname` (String) Optional, conveys a server hostname, can be up to 64 bytes long, and is in the `sname` field.
- `user_context` (Map of String) Arbitrary string data to tie to the subnet. e.g. `{site = "AUS", name = "Austin, Tx"}`

### Read-Only

- `id` (Number) The ID of this resource.

<a id="nestedatt--pools"></a>
### Nested Schema for `pools`

Required:

- `pool` (String)


<a id="nestedatt--option_data"></a>
### Nested Schema for `option_data`

Required:

- `always_send` (Boolean)
- `code` (Number)
- `data` (String)
- `name` (String)


<a id="nestedatt--relay"></a>
### Nested Schema for `relay`

Required:

- `ip_address` (String)

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import kea_remote_subnet4_resource.example 100
```
