package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleStorage_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleStorageDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceStorageConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_storage.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_storage.foo", "name", name),
					resource.TestCheckResourceAttr("data.gridscale_storage.foo", "capacity", "1"),
				),
			},
		},
	})

}

func testAccCheckDataSourceStorageConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_storage" "foo" {
  name   = "%s"
  capacity = 1
}

data "gridscale_storage" "foo" {
	resource_id   = "${gridscale_storage.foo.id}"
}`, name)
}
