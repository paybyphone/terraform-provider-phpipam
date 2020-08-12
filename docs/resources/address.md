# phpipam_address

The `phpipam_address` resource manages an IP address in PHPIPAM. You can use it to create IP address reservations for IP addresses that have been created by other Terraform resources, or if not supplied gets the first free ip address and reserves it in an atomic way. 

~> Don't use `phpipam_first_free_address` to get a free ip, because this datasource only get the next free one and since it doesn't reserve it it's always the same ip that is returned.

## Example Usage

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Reserve the address
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"

  custom_fields = {
    CustomTestAddresses = "terraform-test"
  }
}
```

## Argument Reference

The resource takes the following parameters:

* `subnet_id` - (Required) The database ID of the subnet this IP address belongs to.
* `ip_address` - (Optional) The IP address to reserve. If not defined, first free IP in subnet is used.
* `is_gateway` - (Optional) True if this IP address has been designated as a gateway.
* `description` - (Optional) The description provided to this IP address.
* `hostname` - (Optional) The hostname supplied to this IP address.
* `owner` - (Optional) The owner name provided to this IP address.
* `mac_address` - (Optional) The MAC address provided to this IP address.
* `state_tag_id` - (Optional) The tag ID in the database for the IP address' specific state. NOTE: This is currently represented as an integer but may change to the specific string representation at a later time.
* `skip_ptr_record` - (Optional) True if PTR records are not being created for this IP address.
* `ptr_record_id` - (Optional) The ID of the associated PTR record in the PHPIPAM database.
* `device_id` - (Optional) The ID of the associated device in the PHPIPAM database.
* `switch_port_label` - (Optional) A string port label that is associated with this address.
* `note` - (Optional) The note supplied to this IP address.
* `exclude_ping` - (Optional) True if this address is excluded from ping probes.
* `remove_dns_on_delete` - (Optional) Removes DNS records created by PHPIPAM when the address is deleted from Terraform. Defaults to true.
* `custom_fields` - (Optional) A key/value map of custom fields for this address.

~> Custom fields: PHPIPAM installations with custom fields must have all fields set to optional when using this plugin. For more info see here. Further to this, either ensure that your fields also do not have default values, or ensure the default is set in your TF configuration. Diff loops may happen otherwise!

## Attribute Reference

* `address_id` - The ID of the IP address in the PHPIPAM database.
* `last_seen` - The last time this IP address answered ping probes.
* `edit_date` - The last time this resource was modified.