package phpipam

import "github.com/hashicorp/terraform/helper/schema"

// resourcePHPIPAMVLAN returns the resource structure for the phpipam_vlan
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMVLANCreate,
		Read:   dataSourcePHPIPAMVLANRead,
		Update: resourcePHPIPAMVLANUpdate,
		Delete: resourcePHPIPAMVLANDelete,
		Schema: resourceVLANSchema(),
	}
}

func resourcePHPIPAMVLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateVLAN(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMVLANRead(d, meta)
}

func resourcePHPIPAMVLANUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	if _, err := c.UpdateVLAN(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMVLANRead(d, meta)
}

func resourcePHPIPAMVLANDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	if _, err := c.DeleteVLAN(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
