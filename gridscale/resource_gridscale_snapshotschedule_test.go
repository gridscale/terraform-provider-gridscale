package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v2"
)

func TestAccResourceGridscaleSnapshotSchedule_Basic(t *testing.T) {
	var object gsclient.StorageSnapshotSchedule
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSnapshotScheduleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleSnapshotScheduleConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotScheduleExists("gridscale_snapshotschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshotschedule.foo", "name", name),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSnapshotScheduleConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotScheduleExists("gridscale_snapshotschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshotschedule.foo", "name", "newname"),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSnapshotScheduleConfig_forcenew_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotScheduleExists("gridscale_snapshotschedule.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshotschedule.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleSnapshotScheduleExists(n string, object *gsclient.StorageSnapshotSchedule) resource.TestCheckFunc {
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
		foundObject, err := client.GetStorageSnapshotSchedule(context.Background(), storageID, id)
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

func testAccCheckDataSourceGridscaleSnapshotScheduleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_snapshotschedule" {
			continue
		}

		_, err := client.GetStorageSnapshotSchedule(context.Background(), rs.Primary.Attributes["storage_uuid"], rs.Primary.ID)
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

func testAccCheckDataSourceGridscaleSnapshotScheduleConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
  keep_snapshots = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
}
`, name)
}

func testAccCheckDataSourceGridscaleSnapshotScheduleConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.foo.id
  labels = ["test"]
  keep_snapshots = 1
  run_interval = 60
}
`)
}

func testAccCheckDataSourceGridscaleSnapshotScheduleConfig_forcenew_update() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "new" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.new.id
  keep_snapshots = 1
  run_interval = 60
}
`)
}
