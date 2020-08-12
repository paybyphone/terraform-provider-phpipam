# phpipam_first_free_address

The `phpipam_first_free_address` data source allows you to get the next available IP address in a specific subnet in PHPIPAM. Note that not having any addresses available will cause the Terraform run to fail. Conversely, marking a subnet as unavailable or used will not prevent this data source from returning an IP address, so be aware of this while using this resource.

!> Don't use `phpipam_first_free_address` to get a free ip for `phpipam_address` resource, because this datasource only get the next free one and since it doesn't reserve it it's always the same ip that is returned.

## Example Usage

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Get the first available address
data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}
```

## Argument Reference

The data source takes the following parameters:

* `subnet_id` - The ID of the subnet to look up the address in.

## Attribute Reference

The following attributes are exported:

* `ip_address` - The next available IP address.