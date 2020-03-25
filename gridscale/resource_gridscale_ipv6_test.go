package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v2"
)

func TestAccResourceGridscaleIpv6_Basic(t *testing.T) {
	var object gsclient.IP
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleIpv6DestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleIpv6Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleIpv6Exists("gridscale_ipv6.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv6.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleIpv6Config_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleIpv6Exists("gridscale_ipv6.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv6.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_ipv6.foo", "failover", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_ipv6.foo", "reverse_dns", "test.test"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleIpv6Exists(n string, object *gsclient.IP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}

		projectName := rs.Primary.Attributes["project"]
		client, err := getProjectClientFromMeta(projectName, testAccProvider.Meta())
		if err != nil {
			return err
		}

		id := rs.Primary.ID

		foundObject, err := client.GetIP(emptyCtx, id)

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

func testAccCheckGridscaleIpv6DestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_ipv6" {
			continue
		}
		projectName := rs.Primary.Attributes["project"]
		client, err := getProjectClientFromMeta(projectName, testAccProvider.Meta())
		if err != nil {
			return err
		}
		_, err = client.GetIP(emptyCtx, rs.Primary.ID)
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

func testAccCheckResourceGridscaleIpv6Config_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv6" "foo" {
  project = "default"
  name   = "%s"
}
`, name)
}

func testAccCheckResourceGridscaleIpv6Config_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_ipv6" "foo" {
  project = "default"
  name   = "newname"
  failover = true
  reverse_dns = "test.test"
}
`)
}
