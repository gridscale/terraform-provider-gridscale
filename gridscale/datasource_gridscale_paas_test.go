package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscalePaaSBasic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePaaSConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gridscale_paas.foo", "name", name),
					resource.TestCheckResourceAttr("data.gridscale_paas.foo", "service_template_uuid", "d7a5e8ec-fa78-4d1b-86f9-febe3e16e398"),
				),
			},
		},
	})

}

func testAccCheckDataSourcePaaSConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foo" {
  name = "%s"
  service_template_uuid = "d7a5e8ec-fa78-4d1b-86f9-febe3e16e398"
  parameter {
    param = "mysql_max_connections"
    value = "2000"
    type = "float"
  }
}

data "gridscale_paas" "foo" {
	resource_id   = gridscale_paas.foo.id
}`, name)
}
