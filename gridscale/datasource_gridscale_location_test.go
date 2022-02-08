package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceGridscaleLocation_basic(t *testing.T) {
	locUUID := "45ed677b-3702-4b36-be2a-a2eab9827950"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceLocationConfig_basic(locUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.gridscale_location.foo", "id", locUUID),
					resource.TestCheckResourceAttrSet("data.gridscale_location.foo", "name"),
					resource.TestCheckResourceAttrSet("data.gridscale_location.foo", "cpunode_count"),
				),
			},
		},
	})

}

func testAccCheckDataSourceLocationConfig_basic(locUUID string) string {
	return fmt.Sprintf(`
data "gridscale_location" "foo" {
	resource_id = "%s"
}`, locUUID)
}
