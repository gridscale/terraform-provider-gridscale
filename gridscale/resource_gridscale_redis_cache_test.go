package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleRedisCache_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("redis_cache-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleRedisCacheConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_cache.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_cache.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleRedisCacheConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_cache.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_cache.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleRedisCacheConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_redis_cache" "test" {
	name = "%s"
	release = "7"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleRedisCacheConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_redis_cache" "test" {
	name = "newname"
	release = "7"
	performance_class = "standard"
	labels = ["test"]
}
`)
}
