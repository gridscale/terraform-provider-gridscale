package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleBackup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceBackupConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_backup_list.foo", "id"),
				),
			},
		},
	})

}

func testAccCheckDataSourceBackupConfig_basic() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
	name   = "storage"
	capacity = 1
}
data "gridscale_backup_list" "foo" {
  	storage_uuid = gridscale_storage.foo.id
}`)
}
