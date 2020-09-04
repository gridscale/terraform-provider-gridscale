package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleMarketplaceApplication_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleMarketplaceApplicationDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceMarketplaceApplicationConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_marketplace_application.foo", "id"),
					resource.TestCheckResourceAttrSet("data.gridscale_marketplace_application.foo", "category"),
					resource.TestCheckResourceAttrSet("data.gridscale_marketplace_application.foo", "setup_cores"),
					resource.TestCheckResourceAttrSet("data.gridscale_marketplace_application.foo", "setup_memory"),
					resource.TestCheckResourceAttrSet("data.gridscale_marketplace_application.foo", "setup_storage_capacity"),
					resource.TestCheckResourceAttr("data.gridscale_marketplace_application.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceMarketplaceApplicationConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_marketplace_application" "foo" {
	name = "%s"
	object_storage_path = "s3://testsnapshot/test.gz"
	category = "Archiving"
	setup_cores = 1
	setup_memory = 1
	setup_storage_capacity = 1
}

data "gridscale_marketplace_application" "foo" {
	resource_id   = gridscale_marketplace_application.foo.id
}

`, name)
}
