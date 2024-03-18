package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	return `
data "gridscale_public_network" "foo" {
}`
}
