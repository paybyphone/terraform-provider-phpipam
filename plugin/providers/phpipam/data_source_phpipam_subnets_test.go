package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMSubnetsConfig = `
data "phpipam_subnets" "subnets" {
	section_id = 1
}

output "customer_1_subnet_addr" {
	value = "${element(data.phpipam_subnets.subnets.subnets.*.subnet_address, 0)}"
}
`

func TestAccDataSourcePHPIPAMSubnets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("customer_1_subnet_addr", "10.10.1.0"),
				),
			},
		},
	})
}
