package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

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
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	out, err := c.GetSubnetsInSection(d.Get("section_id").(int))
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(d.Get("section_id").(int)))

	var results []map[string]interface{}
	for _, v := range out {
		results = append(results, flattenSubnetToMap(v))
	}
	err = d.Set("subnets", results)
	if err != nil {
		return err
	}

	return nil
}
