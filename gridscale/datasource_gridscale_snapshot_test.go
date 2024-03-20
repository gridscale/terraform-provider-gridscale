package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleSnapshotBasic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceSnapshotConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_snapshot.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_snapshot.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceSnapshotConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
}

data "gridscale_snapshot" "foo" {
	resource_id   = gridscale_snapshot.foo.id
  	storage_uuid = gridscale_storage.foo.id
}`, name)
}
