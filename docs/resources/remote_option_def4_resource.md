---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kea_remote_option_def4_resource Resource - terraform-provider-kea"
subcategory: ""
description: |-
  Remote OptionDef4 resource
---

# kea_remote_option_def4_resource (Resource)

Remote OptionDef4 resource

## Example Usage

```terraform
resource "kea_remote_option_def4_resource" "example" {
  hostname = "kea-primary.example.com"
  code     = 223
  space    = "dhcp4"
  type     = "string"
  name     = "custom-option"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `code` (Number) DHCP option code. e.g. `222`
- `hostname` (String) Hostname of the kea server to connect to. e.g. `kea.example.com`
- `name` (String) DHCP option name. e.g. `location-identifier`
- `space` (String) The DHCP space for the option-def. e.g. `dhcp4`.
- `type` (String) DHCP option type. e.g. `string`, `uint32`

### Optional

- `array` (Boolean) The false value of the array parameter determines that the option does NOT comprise an array of uint32 values but is, instead, a single value..
- `encapsulate` (String) The name of the option space in which the sub-options are defined.
- `record_types` (String) The record_types value should be non-empty if type is set to "record"; otherwise it must be left blank.

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import kea_remote_option_def4_resource.example 100
```
