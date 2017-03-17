package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccPHPIPAMSubnetName = "phpipam_subnet.subnet"
const testAccPHPIPAMSubnetCIDR = "10.10.3.0/24"
const testAccPHPIPAMSubnetConfig = `
resource "phpipam_subnet" "subnet" {
	subnet_address = "10.10.3.0"
	subnet_mask = 24
	description = "Terraform test subnet"
	section_id = 1
}
`

func TestAccPHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPHPIPAMSubnetDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPHPIPAMSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPHPIPAMSubnetCreated,
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "description", "Terraform test subnet"),
				),
			},
		},
	})
}

func testAccCheckPHPIPAMSubnetCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccPHPIPAMSubnetName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccPHPIPAMSubnetName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).subnetsController
	if _, err := c.GetSubnetByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckPHPIPAMSubnetDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).subnetsController
	_, err := c.GetSubnetsByCIDR(testAccPHPIPAMSubnetCIDR)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): No subnets found":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
