# Terraform Provider Plugin for PHPIPAM

This repository holds a external plugin for a [Terraform][1] provider to manage
resources within [PHPIPAM][2], an open source IP address management system.

[1]: https://www.terraform.io/
[2]: https://phpipam.net/

## About PHPIPAM

[PHPIPAM][2] is an open source IP address management system written in PHP. It
has an evolving [API][3] that allows for the management and lookup of data that
has been entered into the system. Through our Go integration
[phpipam-sdk-go][4], we have been able to take this API and integrate it into
Terraform, allowing for the management and lookup of sections, VLANs, subnets,
and IP addresses, entirely within Terraform.

[3]: https://phpipam.net/api/api_documentation/
[4]: https://github.com/paybyphone/phpipam-sdk-go

## Installing

See the [Plugin Basics][5] page of the Terraform docs to see how to plunk this
into your config. Check the [releases page][6] of this repo to get releases for
Linux, OS X, and Windows.

## Usage

After installation, to use the plugin, simply use any of its resources or data
sources (such as `phpipam_subnet` or `phpipam_address` in a Terraform
configuration.

Credentials can be supplied via configuration variables to the `phpipam`
provider instance, or via environment variables. These are documented in the
next section.

You can see the following example below for a simple usage example that reserves
the first available IP address in a subnet. This address could then be passed
along to the configuration for a VM, say, for example, a
[`vsphere_virtual_machine`][7] resource.

[7]: https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html

```
provider "phpipam" {
  app_id   = "test"
  endpoint = "https://phpipam.example.com/api"
  password = "PHPIPAM_PASSWORD"
  username = "Admin"
}

data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  mask           = 24
}

data "phpipam_first_free_address" "next_address" {
  subnet_id = "${data.phpipam_subnet.subnet.subnet_id}"
}

resource "phpipam_address" {
  subnet_id   = "${data.phpipam_subnet.subnet.subnet_id}"
  ip_address  = "${data.phpipam_first_free_address.next_address.ip_address}"
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"

  lifecycle {
    ignore_changes = [
      "subnet_id",
      "ip_address",
    ]
  }
}
```


### Data Sources

### Resources


## LICENSE

```
Copyright 2017 PayByPhone Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
