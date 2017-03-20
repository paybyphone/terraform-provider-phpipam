package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMSubnetConfig = `
data "phpipam_subnet" "subnet_by_cidr" {
	subnet_address = "10.10.2.0"
	subnet_mask = 24
}

data "phpipam_subnet" "subnet_by_id" {
	subnet_id = "${data.phpipam_subnet.subnet_by_cidr.subnet_id}"
}
`

func TestAccDataSourcePHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_id", "data.phpipam_subnet.subnet_by_id", "subnet_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_address", "data.phpipam_subnet.subnet_by_id", "subnet_address"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_mask", "data.phpipam_subnet.subnet_by_id", "subnet_mask"),
				),
			},
		},
	})
}
