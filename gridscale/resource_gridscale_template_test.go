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

func TestAccResourceGridscaleTemplateBasic(t *testing.T) {
	var object gsclient.Template
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleTemplateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleTemplateConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleTemplateExists("gridscale_template.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_template.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleTemplateConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleTemplateExists("gridscale_template.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_template.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleTemplateExists(n string, object *gsclient.Template) resource.TestCheckFunc {
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

		foundObject, err := client.GetTemplate(context.Background(), id)

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

func testAccCheckGridscaleTemplateDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_template" {
			continue
		}

		_, err := client.GetTemplate(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleTemplateConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_storage" "foo" {
  name   = "%s"
  capacity = 1
}

resource "gridscale_snapshot" "foo" {
  name = "%s"
  storage_uuid = gridscale_storage.foo.id
}

resource "gridscale_template" "foo" {
  name   = "%s"
  snapshot_uuid = gridscale_snapshot.foo.id
}
`, name, name, name)
}

func testAccCheckResourceGridscaleTemplateConfigBasicUpdate() string {
	return `
resource "gridscale_storage" "foo" {
  name   = "newname"
  capacity = 1
}

resource "gridscale_snapshot" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.foo.id
}

resource "gridscale_template" "foo" {
  name   = "newname"
  labels = ["test"]
  snapshot_uuid = gridscale_snapshot.foo.id
}
`
}
