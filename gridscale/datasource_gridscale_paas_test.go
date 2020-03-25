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
					resource.TestCheckResourceAttr("data.gridscale_paas.foo", "service_template_uuid", "8bcb216c-65ec-4c93-925d-1b8feaa5c2c5"),
				),
			},
		},
	})

}

func testAccCheckDataSourcePaaSConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foo" {
  project = "default"
  name = "%s"
  service_template_uuid = "8bcb216c-65ec-4c93-925d-1b8feaa5c2c5"
  parameter {
    param = "mysql_max_connections"
    value = "2000"
    type = "float"
  }
}

data "gridscale_paas" "foo" {
	project = gridscale_paas.foo.project
	resource_id   = gridscale_paas.foo.id
}`, name)
}
