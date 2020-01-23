package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleFirewall_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceFirewallConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_firewall.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_firewall.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceFirewallConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_firewall" "foo" {
  name   = "%s"
  rules_v4_in {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "20:80"
	comment = "test"
  }
  rules_v6_in {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "2000:3000"
	comment = "testv6"
  }
}

data "gridscale_firewall" "foo" {
	resource_id   = gridscale_firewall.foo.id
}

`, name)
}
