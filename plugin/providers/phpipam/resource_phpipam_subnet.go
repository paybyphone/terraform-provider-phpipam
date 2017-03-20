package phpipam

import "github.com/hashicorp/terraform/helper/schema"

// resourcePHPIPAMSubnet returns the resource structure for the phpipam_subnet
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMSubnetCreate,
		Read:   dataSourcePHPIPAMSubnetRead,
		Update: resourcePHPIPAMSubnetUpdate,
		Delete: resourcePHPIPAMSubnetDelete,
		Schema: resourceSubnetSchema(),
	}
}

func resourcePHPIPAMSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnetResource(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateSubnet(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnetResource(d)
	// Remove the CIDR fields from the request, as these fields being present
	// implies that the subnet will be either split or renamed, which is not
	// supported by UpdateSubnet. These are implemented in the API but not in the
	// SDK, so support may be added at a later time.
	in.SubnetAddress = ""
	in.Mask = 0
	if _, err := c.UpdateSubnet(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnetResource(d)

	if _, err := c.DeleteSubnet(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
