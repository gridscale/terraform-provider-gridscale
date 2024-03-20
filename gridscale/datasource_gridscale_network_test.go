package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleNetworkBasic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleNetworkDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceNetworkConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_network.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_network.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.gridscale_network.foo", "dhcp_active", "true"),
					resource.TestCheckResourceAttr(
						"data.gridscale_network.foo", "dhcp_gateway", "192.168.121.1"),
					resource.TestCheckResourceAttr(
						"data.gridscale_network.foo", "dhcp_dns", "192.168.121.2"),
					resource.TestCheckResourceAttr(
						"data.gridscale_network.foo", "dhcp_reserved_subnet.#", "1"),
				),
			},
		},
	})

}

func testAccCheckDataSourceNetworkConfigBasic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_network" "foo" {
  name   = "%s"
  dhcp_active = true
  dhcp_gateway = "192.168.121.1"
  dhcp_dns = "192.168.121.2"
  dhcp_range = "192.168.121.0/27"
  dhcp_reserved_subnet = ["192.168.121.0/31"]
}

data "gridscale_network" "foo" {
	resource_id   = gridscale_network.foo.id
}`, name)
}
