package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTemplate_basic(t *testing.T) {
	name := "Ubuntu 18.04 LTS"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceGridscaleTemplateConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_template.foo", "id"),
				),
			},
		},
	})

}

func testAccCheckDataSourceGridscaleTemplateConfig_basic(name string) string {
	return fmt.Sprintf(`
data "gridscale_template" "foo" {
	name   = "%s"
}
`, name)
}
