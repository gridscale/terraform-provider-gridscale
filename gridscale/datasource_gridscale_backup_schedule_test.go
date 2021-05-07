package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleBackupSchedule_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleBackupScheduleDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceBackupScheduleConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_backupschedule.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_backupschedule.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceBackupScheduleConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_backupschedule" "foo" {
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
  keep_backups = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
  active 	   = true
}
data "gridscale_backupschedule" "foo" {
	resource_id   = gridscale_backupschedule.foo.id
	storage_uuid   = gridscale_storage.foo.id
}
`, name)
}
