package phpipam

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["app_id"],
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["endpoint"],
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["password"],
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["username"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"phpipam_address": resourcePHPIPAMAddress(),
			"phpipam_section": resourcePHPIPAMSection(),
			"phpipam_subnet":  resourcePHPIPAMSubnet(),
			"phpipam_vlan":    resourcePHPIPAMVLAN(),
			"phpipam_first_free_subnet":  resourcePHPIPAMFirstFreeSubnet(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"phpipam_address":            dataSourcePHPIPAMAddress(),
			"phpipam_addresses":          dataSourcePHPIPAMAddresses(),
			"phpipam_first_free_address": dataSourcePHPIPAMFirstFreeAddress(),
			"phpipam_first_free_subnet":  dataSourcePHPIPAMFirstFreeSubnet(),
			"phpipam_section":            dataSourcePHPIPAMSection(),
			"phpipam_subnet":             dataSourcePHPIPAMSubnet(),
			"phpipam_subnets":            dataSourcePHPIPAMSubnets(),
			"phpipam_vlan":               dataSourcePHPIPAMVLAN(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"app_id":   "The application ID required for API requests",
		"endpoint": "The full URL (plus path) to the API endpoint",
		"password": "The password of the PHPIPAM account",
		"username": "The username of the PHPIPAM account",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AppID:    d.Get("app_id").(string),
		Endpoint: d.Get("endpoint").(string),
		Password: d.Get("password").(string),
		Username: d.Get("username").(string),
	}
	return config.Client()
}
