package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleSnapshotSchedule_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSnapshotScheduleDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceSnapshotScheduleConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_snapshotschedule.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_snapshotschedule.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceSnapshotScheduleConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "%s"
  storage_uuid = "${gridscale_storage.foo.id}"
  keep_snapshots = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
}
data "gridscale_snapshotschedule" "foo" {
	resource_id   = "${gridscale_snapshotschedule.foo.id}"
	storage_uuid   = "${gridscale_storage.foo.id}"
}
`, name)
}
