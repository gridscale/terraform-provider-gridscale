package gridscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleStorage_Basic(t *testing.T) {
	var object gsclient.Storage
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleStorageDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleStorageConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleStorageExists("gridscale_storage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "capacity", "1"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleStorageConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleStorageExists("gridscale_storage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "name", "newname"),
				),
			},
		},
	})
}

func TestAccResourceGridscaleStorage_Advanced(t *testing.T) {
	var object gsclient.Storage
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleStorageDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleStorageConfig_advanced(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleStorageExists("gridscale_storage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "storage_type", "storage"),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "last_used_template", "4db64bfc-9fb2-4976-80b5-94ff43b1233a"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleStorageExists(n string, object *gsclient.Storage) resource.TestCheckFunc {
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

		foundObject, err := client.GetStorage(context.Background(), id)

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

func testAccCheckGridscaleStorageDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_storage" {
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

func testAccCheckResourceGridscaleStorageConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "%s"
  capacity = 1
}
`, name)
}

func testAccCheckResourceGridscaleStorageConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "newname"
  capacity = 1
}
`)
}

func testAccCheckResourceGridscaleStorageConfig_advanced(name string) string {
	return fmt.Sprintf(`
resource "gridscale_sshkey" "sshkey" {
  name = "%s"
  sshkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQClJCCOAFyBNIWUpzU4/mFqns5G4+nXzf5iFblNZqtAJmPzKnl0m0Gxj9GV27EkaWpqivVSUblmw3KRWMgCAiJUrMoQt4VAUKUzwdlNZ+6cIDSncEg671SLmCGZmWmVdOR5KaHWlkIRnowfB7UIDyubu/B7r+9L5IPdVgqw3KQW4jZRSsaOOG+I6z0J46c0j+/uJBxuqsr0QD0RQYc2n2Q8O9oNvp3U/L0B5ZYkecAZCCTuGpfNnJdpjj4ww+Qgq/qt4WEIWgVIPEU3B5PlqKZDTO+0JjCsAaQIkN6HOSVHP7h9b+grBnTxSc55CPqBGEBP8zlcne29olJttseJgnBT"
}
resource "gridscale_storage" "foo" {
  name   = "%s"
  capacity = 10
  storage_type= "storage"
  labels = []
  template {
    template_uuid = "4db64bfc-9fb2-4976-80b5-94ff43b1233a"
    hostname = "ubuntu"
    sshkeys = [ gridscale_sshkey.sshkey.id ]
  }
}
`, name, name)
}
