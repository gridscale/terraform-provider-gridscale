package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleBackupScheduleBasic(t *testing.T) {
	var object gsclient.StorageBackupSchedule
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleBackupScheduleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleBackupScheduleConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleBackupScheduleExists("gridscale_backupschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_backupschedule.foo", "name", name),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleBackupScheduleConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleBackupScheduleExists("gridscale_backupschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_backupschedule.foo", "name", "newname"),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleBackupScheduleConfigForceNewUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleBackupScheduleExists("gridscale_backupschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_backupschedule.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleBackupScheduleExists(n string, object *gsclient.StorageBackupSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}
		client := testAccProvider.Meta().(*gsclient.Client)
		id := rs.Primary.ID
		storageID := rs.Primary.Attributes["storage_uuid"]
		foundObject, err := client.GetStorageBackupSchedule(context.Background(), storageID, id)
		if err != nil {
			return err
		}
		if foundObject.Properties.ObjectUUID != id {
			return fmt.Errorf("Object not found")
		}
		*object = foundObject
		return nil
	}
}

func testAccCheckDataSourceGridscaleBackupScheduleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_backupschedule" {
			continue
		}

		_, err := client.GetStorageBackupSchedule(context.Background(), rs.Primary.Attributes["storage_uuid"], rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("Unable to fetching object %s", rs.Primary.ID)
			}
		} else {
			return fmt.Errorf("Object %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDataSourceGridscaleBackupScheduleConfigBasic(name string) string {
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
  active = false
}
`, name)
}

func testAccCheckDataSourceGridscaleBackupScheduleConfigBasicUpdate() string {
	return `
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_backupschedule" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.foo.id
  keep_backups = 1
  run_interval = 60
  next_runtime = "2030-12-30 15:04:05"
  active = true
}
`
}

func testAccCheckDataSourceGridscaleBackupScheduleConfigForceNewUpdate() string {
	return `
resource "gridscale_storage" "new" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_backupschedule" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.new.id
  keep_backups = 1
  run_interval = 60
  next_runtime = "2030-12-30 15:04:05"
  active = true
}
`
}
