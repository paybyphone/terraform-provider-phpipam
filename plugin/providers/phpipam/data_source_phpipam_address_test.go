package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMAddressConfig = `
data "phpipam_address" "address_by_address" {
	ip_address = "10.10.1.245"
}

data "phpipam_address" "address_by_id" {
	address_id = "${data.phpipam_address.address_by_address.address_id}"
}
`

func TestAccDataSourcePHPIPAMAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "address_id", "data.phpipam_address.address_by_id", "address_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "ip_address", "data.phpipam_address.address_by_id", "ip_address"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_address", "description", "Gateway"),
				),
			},
		},
	})
}
