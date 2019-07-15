package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gridscale/gsclient-go"
)

func TestAccDataSourceGridscaleLoadBalancerBasic(t *testing.T) {
	var object gsclient.LoadBalancer
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleLoadBalancerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleLoadBalancerExists("gridscale_loadbalancer.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_loadbalancer.foo", "name", name),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleLoadBalancerExists(n string, object *gsclient.LoadBalancer) resource.TestCheckFunc {
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

		foundObject, err := client.GetLoadBalancer(id)

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

func testAccCheckDataSourceGridscaleLoadBalancerConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "foo" {
	name   = "%s"
}
resource "gridscale_ipv4" "server" {
	name   = "server"
}
resource "gridscale_ipv6" "foo" {
	name   = "%s"
}
resource "gridscale_server" "foo" {
	name   = "%s"
	cores = 2
	memory = 2
	power = true
	ipv4 = "${gridscale_server.foo.id}"
}
resource "gridscale_loadbalancer" "foo" {
	name   = "%s"
	algorithm = "leastconn"
	redirect_http_to_https = false
	listen_ipv4_uuid = "${gridscale_ipv4.foo.id}"
	listen_ipv6_uuid = "${gridscale_ipv6.foo.id}"
	backend_servers {
		backend_server {
			weight = 100
			host   = "${gridscale_ipv4.server.ip}"
		},
	}
	forwarding_rules {
		forwarding_rule {
			listen_port=     8080
			mode       =     "http"
			target_port=     8000
		},
	}
	
}`, name, name, name, name)
}

func testAccCheckDataSourceGridscaleLoadBalancerConfigUpdate() string {
	return fmt.Sprintf(`
resource "gridscale_loadbalancer" "foo" {
  name   = "newname"
}
`)
}

func testAccCheckDataSourceGridscaleLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_loadbalancer" {
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
