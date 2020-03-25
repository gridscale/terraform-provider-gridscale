package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleIPv6_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleIpv6DestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceIPv6Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_ipv6.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_ipv6.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceIPv6Config_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_ipv6" "foo" {
	project = "default"
	name   = "%s"
}

data "gridscale_ipv6" "foo" {
	project = gridscale_ipv6.foo.project
	resource_id   = gridscale_ipv6.foo.id
}
`, name)
}
