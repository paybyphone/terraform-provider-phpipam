package phpipam

import "github.com/hashicorp/terraform/helper/schema"

func dataSourcePHPIPAMSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMSubnetsRead,
		Schema: map[string]*schema.Schema{
			"section_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"subnets": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceSubnetsSchema(),
				},
			},
		},
	}
}

func dataSourcePHPIPAMSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
