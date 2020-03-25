package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleSnapshot_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceSnapshotConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_snapshot.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_snapshot.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceSnapshotConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  project = "default"
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  project = gridscale_storage.foo.project
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
}

data "gridscale_snapshot" "foo" {
	project = gridscale_storage.foo.project
	resource_id   = gridscale_snapshot.foo.id
  	storage_uuid = gridscale_storage.foo.id
}`, name)
}
