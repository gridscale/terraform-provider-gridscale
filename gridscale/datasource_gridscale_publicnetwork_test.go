package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscalePublicNetworkBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourcePublicNetworkConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_public_network.foo", "id"),
					resource.TestCheckResourceAttrSet("data.gridscale_public_network.foo", "name"),
				),
			},
		},
	})

}

func testAccCheckDataSourcePublicNetworkConfigBasic() string {
	return `
data "gridscale_public_network" "foo" {
}`
}
