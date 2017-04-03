package phpipam

import (
	"errors"
	"fmt"

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
	s := meta.(*ProviderPHPIPAMClient).subnetsController
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
	case d.Get("subnet_id").(int) != 0 && (d.Get("description").(string) != "" || d.Get("hostname").(string) != "" || d.Get("custom_field_filter_key").(string) != ""):
		// If subnet_id and one of description or hostname were defined, we do
		// search via GetAddressesInSubnet and return the first found for one of
		// the fields.
		v, err := s.GetAddressesInSubnet(d.Get("subnet_id").(int))
		if err != nil {
			return err
		}
		if len(v) == 0 {
			return errors.New("No addresses were found in the supplied subnet")
		}
		result := -1
		for n, r := range v {
			switch {
			// Double-assert that we don't have empty strings in the conditionals
			// to ensure there there is no edge cases with matching zero values.
			case d.Get("description").(string) != "" && r.Description == d.Get("description").(string):
				result = n
			case d.Get("hostname").(string) != "" && r.Hostname == d.Get("hostname").(string):
				result = n
			case d.Get("custom_field_filter_key").(string) != "":
				fields, err := c.GetAddressCustomFields(r.ID)
				if err != nil {
					return err
				}
				matchKey := d.Get("custom_field_filter_key").(string)
				matchValue := d.Get("custom_field_filter_value").(string)
				matched, err := customFieldFilter(fields, matchKey, matchValue)
				if err != nil {
					return err
				}
				if matched {
					result = n
				}
			}
		}
		if result == -1 {
			return fmt.Errorf("No address found in subnet id %d with supplied description, hostname or custom field value", d.Get("subnet_id"))
		}
		out = v[result]
	default:
		return errors.New("No valid combination of parameters found - need one of address_id, ip_address, or subnet_id and (description|hostname|custom_field_filter_key)")
	}
	flattenAddress(out, d)
	fields, err := c.GetAddressCustomFields(out.ID)
	if err != nil {
		return err
	}
	trimMap(fields)
	if err := d.Set("custom_fields", fields); err != nil {
		return err
	}
	return nil
}
