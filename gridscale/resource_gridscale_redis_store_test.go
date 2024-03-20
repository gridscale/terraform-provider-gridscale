package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleRedisStoreBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("redis_store-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleRedisStoreConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_store.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_store.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleRedisStoreConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_redis_store.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_redis_store.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleRedisStoreConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_redis_store" "test" {
	name = "%s"
	release = "7"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleRedisStoreConfigBasicUpdate() string {
	return `
resource "gridscale_redis_store" "test" {
	name = "newname"
	release = "7"
	performance_class = "standard"
	labels = ["test"]
}
`
}
