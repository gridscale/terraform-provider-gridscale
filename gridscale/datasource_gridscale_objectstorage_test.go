package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleObjectStorageBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleObjectStorageDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceObjectStorageConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "id"),
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "access_key"),
					resource.TestCheckResourceAttrSet("data.gridscale_object_storage_accesskey.foo", "secret_key"),
				),
			},
		},
	})

}

func testAccCheckDataSourceObjectStorageConfigBasic() string {
	return `
resource "gridscale_object_storage_accesskey" "foo" {
}

data "gridscale_object_storage_accesskey" "foo" {
	resource_id   = "${gridscale_object_storage_accesskey.foo.id}"
}
`
}
