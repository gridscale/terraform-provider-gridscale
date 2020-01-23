package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v2"
)

func TestAccResourceGridscaleNetwork_Basic(t *testing.T) {
	var object gsclient.Network
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleNetworkDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleNetworkConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleNetworkConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleNetworkExists("gridscale_network.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_network.foo", "l2security", "true"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleNetworkExists(n string, object *gsclient.Network) resource.TestCheckFunc {
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

		foundObject, err := client.GetNetwork(emptyCtx, id)

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

func testAccCheckGridscaleNetworkDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_network" {
			continue
		}

		_, err := client.GetNetwork(emptyCtx, rs.Primary.ID)
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

func testAccCheckResourceGridscaleNetworkConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_network" "foo" {
  name   = "%s"
}
`, name)
}

func testAccCheckResourceGridscaleNetworkConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_network" "foo" {
  name   = "newname"
  l2security = true
}
`)
}
