package gridscale

import (
	"fmt"
	"testing"

	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGridscaleStorage_Basic(t *testing.T) {
	var object gsclient.Storage
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleStorageConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleStorageExists("gridscale_storage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_storage.foo", "capacity", "1"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleStorageExists(n string, object *gsclient.Storage) resource.TestCheckFunc {
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

		foundObject, err := client.GetStorage(id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUuid != id {
			return fmt.Errorf("Object not found")
		}

		*object = *foundObject

		return nil
	}
}

func testAccCheckDataSourceGridscaleStorageConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "%s"
  capacity = 1
}
`, name)
}
