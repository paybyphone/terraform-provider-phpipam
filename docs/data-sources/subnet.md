# phpipam_subnet

The `phpipam_subnet` data source gets information on a subnet such as its ID (required for creating addresses), description, and more.

## Example Usage

Simple:
```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Reserve the address.
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.2.10"
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"
}
```

With `description_match`:
```hcl
// Look up the subnet (matching on either case of "customer")
data "phpipam_subnet" "subnet" {
  section_id        = 1
  description_match = "[Cc]ustomer 2"
}

// Reserve the address
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"
}
```

With `custom_field_filter`:
```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  section_id = 1

  custom_field_filter = {
    CustomTestSubnets = ".*terraform.*"
  }
}

// Reserve the address
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"
}
```

## Argument Reference

The data source takes the following parameters:


* `section_id` - The ID of the section of the subnet. Required if you are looking up a subnet using the description or description_match arguments.
* `subnet_id` - The ID of the subnet to look up.
* `subnet_address` - The network address of the subnet to look up.
* `subnet_mask` - The subnet mask, in bits, of the subnet to look up.
* `description` - The subnet's description. section_id is required if you want to use this option.
* `description_match` - A regular expression to match against when searching for a subnet. section_id is required if you want to use this option.
* `custom_field_filter` - A map of custom fields to search for. The filter values are regular expressions that follow the RE2 syntax for which you can find documentation here. All fields need to match for the match to succeed.

~> Searches with the `description`, `description_match` and `custom_field_filter` fields return the first match found without any warnings. Conversely, the resource fails if it somehow finds multiple results on a CIDR (subnet and mask) search - this is to assert that you are getting the subnet you requested. If you want to return multiple results, combine this data source with the `phpipam_subnets` data source.

~> An empty or unspecified `custom_field_filter` value is the equivalent to a regular expression that matches everything, and hence will return the first subnet it sees in the section.

-> Arguments are processed in the following order of precedence: `subnet_id`, `subnet_address` and `subnet_mask`, `section_id`, and either one of `description`, `description_match`, or `custom_field_filter`.


## Attribute Reference

The following attributes are exported:

* `subnet_id` - The ID of the subnet in the PHPIPAM database.
* `subnet_address` - The network address of the subnet.
* `subnet_mask` - The subnet mask, in bits.
* `description` - The description set for the subnet.
* `section_id` - The ID of the section for this address in the PHPIPAM database.
* `linked_subnet_id` - The ID of the linked subnet in the PHPIPAM database.
* `vlan_id` - The ID of the VLAN for this subnet in the PHPIPAM database.
* `vrf_id` - The ID of the VRF for this subnet in the PHPIPAM database.
* `master_subnet_id` - The ID of the parent subnet for this subnet in the PHPIPAM database.
* `nameserver_id` - The ID of the nameserver used to assign PTR records for this subnet.
* `show_name` - true if the subnet name is are shown in the section, instead of the network address.
* `permissions` - A JSON representation of the permissions associated with this subnet.
* `create_ptr_records` - true if PTR records are created for addresses in this subnet.
* `display_hostnames` - true if hostnames are displayed instead of IP addresses in the address listing for this subnet.
* `allow_ip_requests` - true if the subnet allows IP requests in PHPIPAM.
* `scan_agent_id` - The ID of the ping scan agent that is used for this subnet.
* `include_in_ping` - true if this subnet is included in ping probes.
* `host_discovery_enabled` - true if this subnet is included in new host scans.
* `is_folder` - true if this subnet is a folder and not an actual subnet.
* `is_full` - true if the subnet has been marked as full.
* `state_tag_id` - The ID of the state tag for this subnet. This may become an actual string representation of this at a later time (example: Used).
* `utilization_threshold` - The subnet's utilization threshold.
* `location_id` - The ID of the location for this subnet.
* `edit_date` - The date this resource was last updated.
* `custom_fields` - A key/value map of custom fields for this subnet.