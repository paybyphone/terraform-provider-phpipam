package phpipam

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourcePHPIPAMSubnet returns the resource structure for the phpipam_subnet
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMFirstFreeSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMFirstFreeSubnetCreate,
		Read:   dataSourcePHPIPAMSubnetRead,
		Update: resourcePHPIPAMSubnetUpdate,
		Delete: resourcePHPIPAMSubnetDelete,
		Schema: resourceFirstFreeSubnetSchema(),
	}
}

func resourcePHPIPAMFirstFreeSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController

	id := d.Get("master_subnet_id").(int)
	mask := d.Get("subnet_mask").(int)

	message, err := c.CreateFirstFreeSubnet(id,mask);
	if err != nil {
		return err
	} 
	cidr_mask := strings.Split(message, "/");
	d.Set("subnet_address", cidr_mask[0])
	if customFields, ok := d.GetOk("custom_fields"); ok {
		subnets, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%s", cidr_mask[0], cidr_mask[1]))

		if err != nil {
			return fmt.Errorf("Could not read subnet after creating: %s", err)
		}

		if len(subnets) != 1 {
			return errors.New("Subnet either missing or multiple results returned by reading subnet after creation")
		}

		d.SetId(strconv.Itoa(subnets[0].ID))
		d.Set("subnet_id", subnets[0].ID)
		d.Set("subnet_address", subnets[0].SubnetAddress)
		d.Set("subnet_mask", subnets[0].Mask)
		if _, err := c.UpdateSubnetCustomFields(subnets[0].ID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}
	return dataSourcePHPIPAMSubnetRead(d, meta)
}

// flattenSubnet(out[0], d)

	// in.SetId(strconv.Itoa(subnet.ID))
	// If we have custom fields, set them now. We need to get the subnet's ID
	// beforehand.


// func resourcePHPIPAMSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
// 	c := meta.(*ProviderPHPIPAMClient).subnetsController
// 	in := expandSubnet(d)
// 	// Remove the CIDR fields from the request, as these fields being present
// 	// implies that the subnet will be either split or renamed, which is not
// 	// supported by UpdateSubnet. These are implemented in the API but not in the
// 	// SDK, so support may be added at a later time.
// 	in.SubnetAddress = ""
// 	in.Mask = 0
// 	if _, err := c.UpdateSubnet(in); err != nil {
// 		return err
// 	}

// 	if err := updateCustomFields(d, c); err != nil {
// 		return err
// 	}

// 	return dataSourcePHPIPAMSubnetRead(d, meta)
// }

// func resourcePHPIPAMFirstChildSubnetDelete(d *schema.ResourceData, meta interface{}) error {
// 	c := meta.(*ProviderPHPIPAMClient).subnetsController
// 	in := expandSubnet(d)

// 	if _, err := c.DeleteSubnet(in.ID); err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
