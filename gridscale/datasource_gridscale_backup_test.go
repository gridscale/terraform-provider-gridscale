package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleBackupBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceBackupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_backup_list.foo", "id"),
				),
			},
		},
	})

}

func testAccCheckDataSourceBackupConfigBasic() string {
	return `
resource "gridscale_storage" "foo" {
	name   = "storage"
	capacity = 1
}
data "gridscale_backup_list" "foo" {
  	storage_uuid = gridscale_storage.foo.id
}`
}
