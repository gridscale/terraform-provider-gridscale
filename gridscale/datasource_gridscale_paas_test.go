package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscalePaaS_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePaaSConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gridscale_paas.foo", "name", name),
					resource.TestCheckResourceAttr("data.gridscale_paas.foo", "service_template_uuid", "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"),
				),
			},
		},
	})

}

func testAccCheckDataSourcePaaSConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foo" {
  name = "%s"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
}

data "gridscale_paas" "foo" {
	resource_id   = "${gridscale_paas.foo.id}"
}`, name)
}
