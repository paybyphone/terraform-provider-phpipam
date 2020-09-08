# phpIPAM Provider

The phpIPAM provider is used to interact with the resources supported by [phpIPAM](https://phpipam.net/) API.

## Example Usage

```hcl
# Provider configuration
terraform {
  required_providers {
    phpipam = {
      source  = "Ouest-France/phpipam"
    }
  }
}

provider "phpipam" {
  app_id   = "test"
  endpoint = "https://phpipam.example.com/api"
  password = "PHPIPAM_PASSWORD"
  username = "Admin"
}

...
```

## Argument Reference

* `app_id` - (Optional) The API application ID, configured in the PHPIPAM API panel. This application ID should have read/write access if you are planning to use the resources, but read-only access should be sufficient if you are only using the data sources. Can also be supplied by the PHPIPAM_APP_ID environment variable.

* `endpoint` - (Optional) The full URL to the PHPIPAM API endpoint, such as `https://phpipam.example.com/api`. Can also be supplied by the PHPIPAM_ENDPOINT_ADDR environment variable.

* `user` - (Optional) The user name to access the PHPIPAM API with. Can also be supplied via the PHPIPAM_USER_NAME variable.

* `password` - (Optional) The password to access the PHPIPAM API with. Can also be supplied via PHPIPAM_PASSWORD to prevent plain text password storage in config.

* `insecure` - (Optional) If true, disable SSL cert verification. Defaults to false.
