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

resource "gridscale_ipv4" "foo1" {
  name   = "newname"
}
resource "gridscale_network" "foo" {
  name   = "newname"
}
resource "gridscale_storage" "foo1" {
  name   = "newname"
  capacity = 1
}
resource "gridscale_server" "foo" {
  name   = "%s"
  cores = 1
  memory = 1
  power = true
  ipv4 = gridscale_ipv4.foo1.id
  network {
		object_uuid = gridscale_network.foo.id
		rules_v4_in {
				order = 0
				protocol = "tcp"
				action = "drop"
				dst_port = "20:80"
				comment = "test"
		}
		rules_v6_in	{
				order = 1
				protocol = "tcp"
				action = "drop"
				dst_port = "10:20"
				comment = "test1"
		}
  	}
  storage {
  	object_uuid = gridscale_storage.foo1.id
  }
}


data "gridscale_server" "foo" {
	resource_id   = gridscale_server.foo.id
}

`, name)
}
