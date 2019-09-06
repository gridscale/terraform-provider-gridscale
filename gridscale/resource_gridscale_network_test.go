package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gridscale/gsclient-go"
)

func TestAccDataSourceGridscaleNetwork_Basic(t *testing.T) {
	var object gsclient.Network
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleNetworkDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleNetworkConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", name),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleNetworkConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "l2security", "true"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleNetworkExists(n string, object *gsclient.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}
		client := testAccProvider.Meta().(*gsclient.Client)
		id := rs.Primary.ID
		foundObject, err := client.GetNetwork(id)
		if err != nil {
			return err
		}
		if foundObject.Properties.ObjectUUID != id {
			return fmt.Errorf("Object not found")
		}
		*object = foundObject
		return nil
	}
}

func testAccCheckDataSourceGridscaleNetworkDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_network" {
			continue
		}
		_, err := client.GetNetwork(rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("Unable to fetching object %s", rs.Primary.ID)
			}
		} else {
			return fmt.Errorf("Object %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDataSourceGridscaleNetworkConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_network" "foo" {
  name   = "%s"
}
`, name)
}

func testAccCheckDataSourceGridscaleNetworkConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_network" "foo" {
  name   = "newname"
  l2security = true
}
`)
}
