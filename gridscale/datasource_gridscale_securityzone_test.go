package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleSecurityZone_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSecurityZoneDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceSecurityZoneConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_paas_securityzone.foo", "location_uuid"),
					resource.TestCheckResourceAttr("data.gridscale_paas_securityzone.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceSecurityZoneConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas_securityzone" "foo" {
  project = "default"
  name = "%s"
}

data "gridscale_paas_securityzone" "foo" {
	project = gridscale_paas_securityzone.foo.project
	resource_id   = gridscale_paas_securityzone.foo.id
}`, name)
}
