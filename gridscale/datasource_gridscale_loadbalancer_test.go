package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccdataSourceGridscaleLoadBalancer_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscaleLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceLoadBalancerConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_loadbalancer.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "name", name),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "backend_server", "leastconn"),
				),
			},
		},
	})

}

func testAccCheckDataSourceLoadBalancerConfig_basic(name string) string {
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
	algorithm = "leastconn"
	redirect_http_to_https = false
	listen_ipv4_uuid = "${gridscale_ipv4.lb.id}"
	listen_ipv6_uuid = "${gridscale_ipv6.lb.id}"
	labels = []
	backend_server {
		weight = 100
		host   = "${gridscale_ipv4.server.ip}"
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}

data "gridscale_loadbalancer" "foo" {
	resource_id   = "${gridscale_loadbalancer.foo.id}"
}`, name, name, name, name)
}
