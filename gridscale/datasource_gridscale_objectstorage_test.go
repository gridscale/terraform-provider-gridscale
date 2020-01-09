package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleObjectStorage_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleObjectStorageDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceObjectStorageConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "id"),
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "access_key"),
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "secret_key"),
				),
			},
		},
	})

}

func testAccCheckDataSourceObjectStorageConfig_basic() string {
	return fmt.Sprint(`
resource "gridscale_object_storage_accesskey" "foo" {
}

data "gridscale_object_storage_accesskey" "foo" {
	resource_id   = "${gridscale_object_storage_accesskey.foo.id}"
}
`)
}
