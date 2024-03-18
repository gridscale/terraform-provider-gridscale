package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMSSQLServer_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("MSSQLServer-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMSSQLServerConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_sqlserver.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_sqlserver.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMSSQLServerConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_sqlserver.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_sqlserver.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMSSQLServerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_sqlserver" "test" {
	name = "%s"
	release = "2019"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMSSQLServerConfig_basic_update() string {
	return `
resource "gridscale_sqlserver" "test" {
	name = "newname"
	release = "2019"
	performance_class = "standard"
	labels = ["test"]
}
`
}
