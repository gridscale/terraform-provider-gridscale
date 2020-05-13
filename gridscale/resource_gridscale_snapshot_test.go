package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleSnapshot_Basic(t *testing.T) {
	var object gsclient.StorageSnapshot
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSnapshotDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleSnapshotConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotExists("gridscale_snapshot.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshot.foo", "name", name),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSnapshotConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotExists("gridscale_snapshot.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshot.foo", "name", "newname"),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSnapshotConfig_forcenew_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSnapshotExists("gridscale_snapshot.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_snapshot.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleSnapshotExists(n string, object *gsclient.StorageSnapshot) resource.TestCheckFunc {
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
		foundObject, err := client.GetStorageSnapshot(context.Background(), storageID, id)
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

func testAccCheckDataSourceGridscaleSnapshotDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_snapshot" {
			continue
		}

		_, err := client.GetStorageSnapshot(context.Background(), rs.Primary.Attributes["storage_uuid"], rs.Primary.ID)
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

func testAccCheckDataSourceGridscaleSnapshotConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
  rollback {
	id = "first"
  }
}
`, name)
}

func testAccCheckDataSourceGridscaleSnapshotConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.foo.id
  labels = ["test"]
  rollback {
	id = "first"
  }
  rollback {
	id = "second"
  }
}
`)
}

func testAccCheckDataSourceGridscaleSnapshotConfig_forcenew_update() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "new" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.new.id
}
`)
}
