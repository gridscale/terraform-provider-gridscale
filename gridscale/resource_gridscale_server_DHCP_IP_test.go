package gridscale

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleServerDHCPIP_Basic(t *testing.T) {
	var object gsclient.ServerDHCPIP
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleServerDHCPIPDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleServerDHCPIPConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleServerDHCPIPExists("gridscale_server_DHCP_IP.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_active", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_gateway", "192.168.121.1"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_dns", "192.168.121.2"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_reserved_subnet.#", "1"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleServerDHCPIPConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleServerDHCPIPExists("gridscale_server_DHCP_IP.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "l2security", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_active", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_gateway", "192.168.122.1"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_dns", "192.168.122.2"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "dhcp_reserved_subnet.#", "1"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleServerDHCPIPExists(n string, object *gsclient.ServerDHCPIP) resource.TestCheckFunc {
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
		idList := strings.Split(id, "-")
		networkUUID := idList[0]
		serverUUID := idList[1]
		ip := idList[2]

		pinnedServerList, err := client.GetPinnedServerList(context.Background(), networkUUID)

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

func testAccCheckGridscaleServerDHCPIPDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_server_DHCP_IP" {
			continue
		}

		_, err := client.GetServerDHCPIP(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleServerDHCPIPConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_server_DHCP_IP" "foo" {
  name   = "%s"
  dhcp_active = true
  dhcp_gateway = "192.168.121.1"
  dhcp_dns = "192.168.121.2"
  dhcp_range = "192.168.121.0/27"
  dhcp_reserved_subnet = ["192.168.121.0/31"]
}
`, name)
}

func testAccCheckResourceGridscaleServerDHCPIPConfig_basic_update() string {
	return fmt.Sprint(`
resource "gridscale_server_DHCP_IP" "foo" {
  name   = "newname"
  l2security = true
  dhcp_active = true
  dhcp_gateway = "192.168.122.1"
  dhcp_dns = "192.168.122.2"
  dhcp_range = "192.168.122.0/27"
  dhcp_reserved_subnet = ["192.168.122.0/31"]
}
`)
}
