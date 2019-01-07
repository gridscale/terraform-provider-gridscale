package gridscale

import (
	"fmt"
	"testing"

	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGridscaleIpv4_Basic(t *testing.T) {
	var object gsclient.Ip
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleIpDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleIpv4Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleIpv4Exists("gridscale_ipv4.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv4.foo", "name", name),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleIpv4Exists(n string, object *gsclient.Ip) resource.TestCheckFunc {
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

		foundObject, err := client.GetIp(id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUuid != id {
			return fmt.Errorf("Object not found")
		}

		*object = *foundObject

		return nil
	}
}

func testAccCheckDataSourceGridscaleIpDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_server" {
			continue
		}

		_, err := client.GetIp(rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(*gsclient.RequestError); ok {
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

func testAccCheckDataSourceGridscaleIpv4Config_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "foo" {
  name   = "%s"
}
`, name)
}
