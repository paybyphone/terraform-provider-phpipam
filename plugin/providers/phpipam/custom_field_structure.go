package phpipam

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/controllers/addresses"
	"github.com/paybyphone/phpipam-sdk-go/controllers/subnets"
	"github.com/paybyphone/phpipam-sdk-go/controllers/vlans"
)

// customFieldFilter takes a map[string]interface{} and attempts to find a
// match based off the key, and value attributes.  The data is matched against
// value as a regex. For exact matching, ensure your match is enclosed in the ^
// (start of line) and the $ (end of line) anchors.
//
// PHPIPAM currently stringifies most, if not all, values coming out of the
// API. As such, we don't attempt to cast here - anything that is not a string
// is an error. If the need arises for this to be changed at some point in time
// this function will be updated.
func customFieldFilter(data map[string]interface{}, matchKey, matchValue string) (bool, error) {
	for k, w := range data {
		if k == matchKey {
			switch v := w.(type) {
			case string:
				return regexp.MatchString(matchValue, v)
			default:
				return false, fmt.Errorf("Key %s is not a string or stringified value, which we currently do not support", k)
			}
		}
	}
	return false, nil
}

// trimMap goes thru a map[string]interface{}, and removes keys that
// have zero or nil values.
func trimMap(in map[string]interface{}) {
	for k, v := range in {
		switch {
		case v == nil:
			fallthrough
		case reflect.ValueOf(v).Interface() == reflect.Zero(reflect.TypeOf(v)).Interface():
			delete(in, k)
		}
	}
}

// updateCustomFields performs an update of custom fields on a resource, with
// the following stipulations:
//  * If we have custom fields, we need to do a diff on what is set versus
//    what isn't set, and ensure that we clear out the keys that aren't set.
//    Since our SDK does not currently support NOT NULL custom fields in
//    PHPIPAM, we can safely set these to nil.
//  * If we don't have a value for
//    custom_fields at all, set all keys to nil and update so that all custom
//    fields get blown away.
func updateCustomFields(d *schema.ResourceData, client interface{}) error {
	customFields := make(map[string]interface{})
	if m, ok := d.GetOk("custom_fields"); ok {
		customFields = m.(map[string]interface{})
	}
	var old map[string]interface{}
	var err error
	switch c := client.(type) {
	case *addresses.Controller:
		old, err = c.GetAddressCustomFields(d.Get("address_id").(int))
	case *subnets.Controller:
		old, err = c.GetSubnetCustomFields(d.Get("subnet_id").(int))
	case *vlans.Controller:
		old, err = c.GetVLANCustomFields(d.Get("vlan_id").(int))
	}
	if err != nil {
		return fmt.Errorf("Error getting custom fields for updating: %s", err)
	}
nextKey:
	for k := range old {
		for l, v := range customFields {
			if k == l {
				customFields[l] = v
				continue nextKey
			}
		}
		customFields[k] = nil
	}

	switch c := client.(type) {
	case *addresses.Controller:
		_, err = c.UpdateAddressCustomFields(d.Get("address_id").(int), customFields)
	case *subnets.Controller:
		_, err = c.UpdateSubnetCustomFields(d.Get("subnet_id").(int), customFields)
	case *vlans.Controller:
		_, err = c.UpdateVLANCustomFields(d.Get("vlan_id").(int), d.Get("name").(string), customFields)
	}
	return err
}
