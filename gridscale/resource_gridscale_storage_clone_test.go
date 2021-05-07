package gridscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleStorageClone_Basic(t *testing.T) {
	var object gsclient.Storage

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleStorageDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleStorageCloneConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleStorageExists("gridscale_storage_clone.foo", &object),
					resource.TestCheckResourceAttrSet("gridscale_storage_clone.foo", "name"),
					resource.TestCheckResourceAttr(
						"gridscale_storage_clone.foo", "capacity", "1"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleStorageCloneConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleStorageExists("gridscale_storage_clone.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_storage_clone.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_storage_clone.foo", "capacity", "2"),
				),
			},
		},
	})
}

func testAccCheckGridscaleStorageCloneDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_storage_clone" {
			continue
		}

		//We wait a while for the storage to delete, since it is not instant
		time.Sleep(time.Second * 5)

		_, err := client.GetStorage(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleStorageCloneConfig_basic() string {
	return fmt.Sprint(`
resource "gridscale_storage" "foo" {
  name   = "test"
  capacity = 1
}

resource "gridscale_storage_clone" "foo" {
  source_storage_id   = gridscale_storage.foo.id
  name = "desired_name"
  storage_type = "storage_high"
}
`)
}

func testAccCheckResourceGridscaleStorageCloneConfig_basic_update() string {
	return fmt.Sprint(`
resource "gridscale_storage" "foo" {
	name   = "test"
	capacity = 1
}

resource "gridscale_storage_clone" "foo" {
  source_storage_id   = gridscale_storage.foo.id
  name = "newname"
  capacity = 2
  storage_type = "storage_insane"
}
`)
}
