package phpipam

import (
	"errors"

	"github.com/Ouest-France/phpipam-sdk-go/controllers/sections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePHPIPAMSection() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMSectionRead,
		Schema: dataSourceSectionSchema(),
	}
}

func dataSourcePHPIPAMSectionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	var out sections.Section
	var err error
	// We need to determine how to get the section. An ID search takes priority,
	// and after that section name.
	switch {
	case d.Get("section_id").(int) != 0:
		out, err = c.GetSectionByID(d.Get("section_id").(int))
		if err != nil {
			return err
		}
	case d.Get("name").(string) != "":
		out, err = c.GetSectionByName(d.Get("name").(string))
		if err != nil {
			return err
		}
	default:
		return errors.New("section_id or name not defined, cannot proceed with reading data")
	}
	err = flattenSection(out, d)

	return err
}
