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

func TestAccResourceGridscaleLoadBalancerBasic(t *testing.T) {
	var object gsclient.LoadBalancer
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleLoadBalancerConfig_basic(name, "leastconn"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleLoadBalancerExists("gridscale_loadbalancer.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_loadbalancer.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_loadbalancer.foo", "algorithm", "leastconn"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleLoadBalancerConfig_update("newname", "roundrobin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleLoadBalancerExists("gridscale_loadbalancer.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_loadbalancer.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_loadbalancer.foo", "algorithm", "roundrobin"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleLoadBalancerExists(n string, object *gsclient.LoadBalancer) resource.TestCheckFunc {
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

		foundObject, err := client.GetLoadBalancer(context.Background(), id)

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

func testAccCheckResourceGridscaleLoadBalancerConfig_basic(name string, algorithm string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "lb" {
	name   = "ipv4-%s"
}
resource "gridscale_ipv6" "lb" {
	name   = "ipv6-%s"
}
resource "gridscale_ipv4" "server" {
	name   = "server-%s"
}
resource "gridscale_loadbalancer" "foo" {
	name   = "%s"
	algorithm = "%s"
	redirect_http_to_https = false
	listen_ipv4_uuid = gridscale_ipv4.lb.id
	listen_ipv6_uuid = gridscale_ipv6.lb.id
	labels = []
	backend_server {
		weight = 100
		host   = gridscale_ipv4.server.ip
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}`, name, name, name, name, algorithm)
}

func testAccCheckResourceGridscaleLoadBalancerConfig_update(name string, algorithm string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "lb" {
	name   = "ipv4-%s"
}
resource "gridscale_ipv6" "lb" {
	name   = "ipv6-%s"
}
resource "gridscale_ipv4" "server" {
	name   = "server-%s"
}
resource "gridscale_loadbalancer" "foo" {
	name   = "%s"
	algorithm = "%s"
	redirect_http_to_https = false
	listen_ipv4_uuid = gridscale_ipv4.lb.id
	listen_ipv6_uuid = gridscale_ipv6.lb.id
	labels = []
	backend_server {
		weight = 100
		host   = gridscale_ipv4.server.ip
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}`, name, name, name, name, algorithm)
}

func testAccCheckGridscaleLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_loadbalancer" {
			continue
		}

		_, err := client.GetLoadBalancer(context.Background(), rs.Primary.ID)
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
