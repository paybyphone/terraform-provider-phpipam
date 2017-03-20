package phpipam

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/controllers/addresses"
)

func dataSourcePHPIPAMAddress() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMAddressRead,
		Schema: dataSourceAddressSchema(),
	}
}

func dataSourcePHPIPAMAddressRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	var out addresses.Address
	// We need to determine how to get the address. An ID search takes priority,
	// and after that addresss.
	switch {
	case d.Get("address_id").(int) != 0:
		var err error
		out, err = c.GetAddressByID(d.Get("address_id").(int))
		if err != nil {
			return err
		}
	case d.Get("ip_address").(string) != "":
		v, err := c.GetAddressesByIP(d.Get("ip_address").(string))
		if err != nil {
			return err
		}
		// Only one result should be returned by this search. Fail on multiples.
		if len(v) != 1 {
			return errors.New("Address search returned either zero or multiple results. Please correct your search and try again")
		}
		out = v[0]
	default:
		return errors.New("address_id or ip_address not defined, cannot proceed with reading data")
	}
	flattenAddress(out, d)
	return nil
}
