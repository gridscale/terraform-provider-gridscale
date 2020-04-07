package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v2"
)

func TestAccResourceGridscaleIpv4_Basic(t *testing.T) {
	var object gsclient.IP
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleIpv4DestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleIpv4Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleIpv4Exists("gridscale_ipv4.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv4.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleIpv4Config_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleIpv4Exists("gridscale_ipv4.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv4.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_ipv4.foo", "failover", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_ipv4.foo", "reverse_dns", "test.test"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleIpv4Exists(n string, object *gsclient.IP) resource.TestCheckFunc {
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

		foundObject, err := client.GetIP(context.Background(), id)

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

func testAccCheckGridscaleIpv4DestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_ipv4" {
			continue
		}

		_, err := client.GetIP(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleIpv4Config_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "foo" {
  name   = "%s"
}
`, name)
}

func testAccCheckResourceGridscaleIpv4Config_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "foo" {
  name   = "newname"
  failover = true
  reverse_dns = "test.test"
}
`)
}
