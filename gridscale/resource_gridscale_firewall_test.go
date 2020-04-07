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

func TestAccResourceGridscaleFirewall_Basic(t *testing.T) {
	var object gsclient.Firewall
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleFirewallConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleFirewallExists("gridscale_firewall.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_firewall.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleFirewallConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleFirewallExists("gridscale_firewall.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_firewall.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleFirewallExists(n string, object *gsclient.Firewall) resource.TestCheckFunc {
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

		foundObject, err := client.GetFirewall(context.Background(), id)

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

func testAccCheckGridscaleFirewallDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_firewall" {
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

func testAccCheckResourceGridscaleFirewallConfig_basic(name string) string {
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
  labels = []
}
`, name)
}

func testAccCheckResourceGridscaleFirewallConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_firewall" "foo" {
  name   = "newname"
  rules_v4_out {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "20:80"
	comment = "test1"
  }
  rules_v6_out {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "2000:3000"
	comment = "testv6"
  }
}
`)
}
