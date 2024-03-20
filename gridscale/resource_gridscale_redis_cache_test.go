package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleRedisCacheBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("redis_cache-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleRedisCacheConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_cache.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_cache.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleRedisCacheConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_cache.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_cache.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleRedisCacheConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_redis_cache" "test" {
	name = "%s"
	release = "7"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleRedisCacheConfigBasicUpdate() string {
	return `
resource "gridscale_redis_cache" "test" {
	name = "newname"
	release = "7"
	performance_class = "standard"
	labels = ["test"]
}
`
}
