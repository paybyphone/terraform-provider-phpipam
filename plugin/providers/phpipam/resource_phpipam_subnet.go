package phpipam

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/controllers/subnets"
)

func resourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMSubnetCreate,
		Read:   resourcePHPIPAMSubnetRead,
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

	return resourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMSubnetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	var out subnets.Subnet
	// We need to determine how to get the subnet. An ID search takes priority,
	// and after that subnets.
	if id := d.Get("subnet_id").(int); id != 0 {
		var err error
		out, err = c.GetSubnetByID(id)
		if err != nil {
			return err
		}
	} else {
		v, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", d.Get("subnet_address"), d.Get("subnet_mask")))
		if err != nil {
			return err
		}
		// This should not happen, as a CIDR search from what I have seen so far
		// generally only returns 1 subnet ever. Nontheless, the API spec and the
		// function return a slice of subnets, so we need to assert that the search
		// has only retuned a single match.
		if len(v) != 1 {
			return errors.New("CIDR search returned either zero or multiple results. Please correct your search and try again")
		}
		out = v[0]
	}
	flattenSubnet(out, d)
	return nil
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

	return resourcePHPIPAMSubnetRead(d, meta)
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
