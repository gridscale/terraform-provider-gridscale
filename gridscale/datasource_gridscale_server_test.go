package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleServer_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleIpv4DestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceServerConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_server.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_server.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceServerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_server" "foo" {
  name   = "%s"
  cores = 1
  memory = 1
}

data "gridscale_server" "foo" {
	resource_id   = gridscale_server.foo.id
}

`, name)
}
