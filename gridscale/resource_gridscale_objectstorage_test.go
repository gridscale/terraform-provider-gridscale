package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleObjectStorageBasic(t *testing.T) {
	var object gsclient.ObjectStorageAccessKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleObjectStorageDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleObjectStorageConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleObjectStorageExists("gridscale_object_storage_accesskey.foo", &object),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_accesskey.foo", "access_key"),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_accesskey.foo", "secret_key"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleObjectStorageExists(n string, object *gsclient.ObjectStorageAccessKey) resource.TestCheckFunc {
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

		foundObject, err := client.GetObjectStorageAccessKey(context.Background(), id)

		if err != nil {
			return err
		}

		if foundObject.Properties.AccessKey != id {
			return fmt.Errorf("Object not found")
		}

		*object = foundObject

		return nil
	}
}

func testAccCheckGridscaleObjectStorageDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_object_storage_accesskey" {
			continue
		}

		_, err := client.GetObjectStorageAccessKey(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleObjectStorageConfigBasic() string {
	return `
resource "gridscale_object_storage_accesskey" "foo" {
}
`
}
