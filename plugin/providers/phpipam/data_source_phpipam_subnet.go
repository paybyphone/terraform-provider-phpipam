package phpipam

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/controllers/subnets"
)

func dataSourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMSubnetRead,
		Schema: dataSourceSubnetSchema(),
	}
}

func dataSourcePHPIPAMSubnetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	s := meta.(*ProviderPHPIPAMClient).sectionsController
	var out subnets.Subnet
	// We need to determine how to get the subnet. An ID search takes priority,
	// and after that subnets.
	switch {
	case d.Get("subnet_id").(int) != 0:
		var err error
		out, err = c.GetSubnetByID(d.Get("subnet_id").(int))
		if err != nil {
			return err
		}
	case d.Get("subnet_address").(string) != "" && d.Get("subnet_mask").(int) != 0:
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
	case d.Get("section_id").(int) != 0 && (d.Get("description").(string) != "" || d.Get("description_match").(string) != "" || d.Get("custom_field_filter_key").(string) != ""):
		// If section_id and description were both defined, we do a search via
		// GetSubnetsInSection for the description and return the first match.
		v, err := s.GetSubnetsInSection(d.Get("section_id").(int))
		if err != nil {
			return err
		}
		if len(v) == 0 {
			return errors.New("No subnets were found in the supplied section")
		}
		result := -1
		for n, r := range v {
			switch {
			// Double-assert that we don't have empty strings in the conditionals
			// to ensure there there is no edge cases with matching zero values.
			case d.Get("description_match").(string) != "":
				// Don't trap error here because we should have already validated the regex via the ValidateFunc.
				if matched, _ := regexp.MatchString(d.Get("description_match").(string), r.Description); matched {
					result = n
				}
			case d.Get("description").(string) != "" && r.Description == d.Get("description").(string):
				result = n
			case d.Get("custom_field_filter_key").(string) != "":
				// Skip folders for now as there is issues pulling them down in the API.
				if r.IsFolder {
					continue
				}
				fields, err := c.GetSubnetCustomFields(r.ID)
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
			return fmt.Errorf("No subnet found in section id %d with supplied description for description/custom field filter", d.Get("section_id"))
		}
		out = v[result]
	default:
		return errors.New("No valid combination of parameters found - need one of subnet_id, subnet_address and subnet_mask, or section_id and (description|description_match|custom_field_filter_key)")
	}
	flattenSubnet(out, d)
	fields, err := c.GetSubnetCustomFields(out.ID)
	if err != nil {
		return err
	}
	trimMap(fields)
	if err := d.Set("custom_fields", fields); err != nil {
		return err
	}
	return nil
}
