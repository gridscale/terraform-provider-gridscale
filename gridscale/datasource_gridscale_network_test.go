package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleNetwork_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleNetworkDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceNetworkConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_network.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_network.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceNetworkConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_network" "foo" {
  name   = "%s"
}

data "gridscale_network" "foo" {
	resource_id   = gridscale_network.foo.id
}`, name)
}
