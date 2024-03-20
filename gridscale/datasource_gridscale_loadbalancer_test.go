package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleLoadBalancerBasic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceLoadBalancerConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_loadbalancer.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "name", name),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "forwarding_rule.0.listen_port", "80"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "backend_server.#", "1"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "backend_server.0.weight", "100"),
					resource.TestCheckResourceAttr("data.gridscale_loadbalancer.foo", "backend_server.0.proxy_protocol", "v2"),
				),
			},
		},
	})

}

func testAccCheckDataSourceLoadBalancerConfigBasic(name string) string {
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
	listen_ipv4_uuid = gridscale_ipv4.lb.id
	listen_ipv6_uuid = gridscale_ipv6.lb.id
	labels = []
	backend_server {
		weight = 100
		host   = gridscale_ipv4.server.ip
		proxy_protocol = "v2"
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}

data "gridscale_loadbalancer" "foo" {
	resource_id   = gridscale_loadbalancer.foo.id
}`, name, name, name, name)
}
