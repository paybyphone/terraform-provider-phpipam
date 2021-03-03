package main

import (
	"github.com/Ouest-France/terraform-provider-phpipam/plugin/providers/phpipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: phpipam.Provider,
	})
}
