package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscalePublicNetwork_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourcePublicNetworkConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_public_network.foo", "id"),
					resource.TestCheckResourceAttrSet("data.gridscale_public_network.foo", "name"),
				),
			},
		},
	})

}

func testAccCheckDataSourcePublicNetworkConfig_basic() string {
	return fmt.Sprint(`
data "gridscale_public_network" "foo" {
	project = "default"
}`)
}
