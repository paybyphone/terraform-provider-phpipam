package phpipam

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/phpipam"
)

// resourcePHPIPAMAddress returns the resource structure for the phpipam_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMAddressCreate,
		Read:   dataSourcePHPIPAMAddressRead,
		Update: resourcePHPIPAMAddressUpdate,
		Delete: resourcePHPIPAMAddressDelete,
		Schema: resourceAddressSchema(),
	}
}

func resourcePHPIPAMAddressCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateAddress(in); err != nil {
		return err
	}

	// If we have custom fields, set them now. We need to get the IP address's ID
	// beforehand.
	if customFields, ok := d.GetOk("custom_fields"); ok {
		addrs, err := c.GetAddressesByIP(in.IPAddress)
		if err != nil {
			return fmt.Errorf("Could not read IP address after creating: %s", err)
		}

		if len(addrs) != 1 {
			return errors.New("IP address either missing or multiple results returned by reading IP after creation")
		}

		if _, err := c.UpdateAddressCustomFields(addrs[0].ID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	// IPAddress and SubnetID need to be removed for update requests.
	in.IPAddress = ""
	in.SubnetID = 0
	if _, err := c.UpdateAddress(in); err != nil {
		return err
	}

	// Set custom fields here. We need to make sure of a couple of things:
	//  * If we have custom fields, we need to do a diff on what is set versus
	//    what isn't set, and ensure that we clear out the keys that aren't set.
	//    Since our SDK does not currently support NOT NULL custom fields in
	//    PHPIPAM, we can safely set these to nil.
	//  * If we don't have a value for
	//    custom_fields at all, set all keys to nil and update so that all custom
	//    fields get blown away.
	fields, err := c.GetAddressCustomFields(in.ID)
	if err != nil {
		return fmt.Errorf("Error getting custom fields for updating: %s", err)
	}
	for k, v := range fields {
	}

	if _, err := c.UpdateAddressCustomFields(in.ID, customFields.(map[string]interface{})); err != nil {
		return err
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMAddressDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	if _, err := c.DeleteAddress(in.ID, phpipam.BoolIntString(d.Get("remove_dns_on_delete").(bool))); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
