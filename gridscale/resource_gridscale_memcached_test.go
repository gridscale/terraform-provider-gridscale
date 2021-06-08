package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMemcached_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("memcached-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMemcachedConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_memcached.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_memcached.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMemcachedConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_memcached.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_memcached.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMemcachedConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_memcached" "test" {
	name = "%s"
	release = "1.5"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMemcachedConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_memcached" "test" {
	name = "newname"
	release = "1.5"
	performance_class = "standard"
	max_core_count = 20
	labels = ["test"]
}
`)
}
