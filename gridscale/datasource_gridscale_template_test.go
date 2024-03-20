package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTemplateBasic(t *testing.T) {
	name := "Ubuntu 22.04 LTS (Jammy Jellyfish) "
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceGridscaleTemplateConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_template.foo", "id"),
				),
			},
		},
	})

}

func testAccCheckDataSourceGridscaleTemplateConfigBasic(name string) string {
	return fmt.Sprintf(`
data "gridscale_template" "foo" {
	name   = "%s"
}
`, name)
}
