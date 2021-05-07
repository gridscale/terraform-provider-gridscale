package gridscale

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceISOImage_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceGridscaleISOImageConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_isoimage.foo", "id"),
				),
			},
		},
	})

}

func testAccCheckDataSourceGridscaleISOImageConfig_basic(name string) string {
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

data "gridscale_isoimage" "foo" {
	resource_id   = gridscale_isoimage.foo.id
}
`, name, name)
}
