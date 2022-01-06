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
	var object gsclient.ServerWithIP
	name := fmt.Sprintf("TEST-object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleServerDHCPIPDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleServerDHCPIPConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleServerDHCPIPExists("gridscale_server_DHCP_IP.foo", &object),
					resource.TestCheckResourceAttrSet(
						"gridscale_server_DHCP_IP.foo", "server_uuid"),
					resource.TestCheckResourceAttrSet(
						"gridscale_server_DHCP_IP.foo", "network_uuid"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "ip", "192.168.121.4"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleServerDHCPIPConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleServerDHCPIPExists("gridscale_server_DHCP_IP.foo", &object),
					resource.TestCheckResourceAttrSet(
						"gridscale_server_DHCP_IP.foo", "server_uuid"),
					resource.TestCheckResourceAttrSet(
						"gridscale_server_DHCP_IP.foo", "network_uuid"),
					resource.TestCheckResourceAttr(
						"gridscale_server_DHCP_IP.foo", "ip", "192.168.121.5"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleServerDHCPIPExists(n string, object *gsclient.ServerWithIP) resource.TestCheckFunc {
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
		idList := strings.Split(id, serverDHCPIPUUIDDelimeter)
		networkUUID := idList[0]
		serverUUID := idList[1]
		ip := idList[2]

		pinnedServerList, err := client.GetPinnedServerList(context.Background(), networkUUID)
		if err != nil {
			return err
		}

		found := false
		var foundObject gsclient.ServerWithIP
		for _, pinnedServer := range pinnedServerList.List {
			if pinnedServer.ServerUUID == serverUUID && pinnedServer.IP == ip {
				found = true
				foundObject = pinnedServer
			}
		}

		if !found {
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

		id := rs.Primary.ID
		idList := strings.Split(id, serverDHCPIPUUIDDelimeter)
		networkUUID := idList[0]
		serverUUID := idList[1]
		ip := idList[2]

		pinnedServerList, err := client.GetPinnedServerList(context.Background(), networkUUID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			}
		}

		for _, pinnedServer := range pinnedServerList.List {
			if pinnedServer.ServerUUID == serverUUID && pinnedServer.IP == ip {
				return fmt.Errorf("Object %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckResourceGridscaleServerDHCPIPConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_network" "testnet"{
	name = "%s"
	dhcp_active = true
	dhcp_gateway = "192.168.121.1"
	dhcp_dns = "192.168.121.2"
	dhcp_range = "192.168.121.0/24"
	timeouts {
		create="10m"
	}
}
	
resource "gridscale_server" "testserver" {
	name   = "%s"
	cores  = 1
	memory = 2
	network {
		object_uuid = gridscale_network.testnet.id
	}
}

resource "gridscale_server_DHCP_IP" "foo" {
  server_uuid = gridscale_server.testserver.id
  network_uuid = gridscale_network.testnet.id
  ip = "192.168.121.4"
}
`, name, name)
}

func testAccCheckResourceGridscaleServerDHCPIPConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_network" "testnet"{
	name = "TEST-newname"
	dhcp_active = true
	dhcp_gateway = "192.168.121.1"
	dhcp_dns = "192.168.121.2"
	dhcp_range = "192.168.121.0/24"
	timeouts {
		create="10m"
	}
}
	
resource "gridscale_server" "testserver" {
	name   = "TEST-newname"
	cores  = 1
	memory = 2
	network {
		object_uuid = gridscale_network.testnet.id
	}
}

resource "gridscale_server_DHCP_IP" "foo" {
	server_uuid = gridscale_server.testserver.id
	network_uuid = gridscale_network.testnet.id
	ip = "192.168.121.5"
}
	`)
}
