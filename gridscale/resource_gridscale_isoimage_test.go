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

func TestAccResourceGridscaleISOImage_Basic(t *testing.T) {
	var object gsclient.ISOImage
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleISOImageDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleISOImageConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleISOImageExists("gridscale_isoimage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_isoimage.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleISOImageConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleISOImageExists("gridscale_isoimage.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_isoimage.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleISOImageExists(n string, object *gsclient.ISOImage) resource.TestCheckFunc {
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

		foundObject, err := client.GetISOImage(context.Background(), id)

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

func testAccCheckGridscaleISOImageDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_isoimage" {
			continue
		}

		_, err := client.GetISOImage(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleISOImageConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_isoimage" "foo" {
  name   = "%s"
  source_url = "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso"
}

resource "gridscale_server" "foo" {
  name   = "%s"
  cores = 1
  memory = 1
  isoimage = gridscale_isoimage.foo.id
}
`, name, name)
}

func testAccCheckResourceGridscaleISOImageConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_isoimage" "foo" {
  name   = "newname"
  source_url = "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso"
}
`)
}
