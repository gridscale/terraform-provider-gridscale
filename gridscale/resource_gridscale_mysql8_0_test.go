package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMySQL8_0_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("MySQL-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMySQL8_0Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mysql8_0.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mysql8_0.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMySQL8_0Config_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mysql8_0.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mysql8_0.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMySQL8_0Config_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_mysql8_0" "test" {
	name = "%s"
	release = "8.0"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMySQL8_0Config_basic_update() string {
	return `
resource "gridscale_mysql8_0" "test" {
	name = "newname"
	release = "8.0"
	performance_class = "standard"
	max_core_count = 20
	mysql_max_connections = 2000
	labels = ["test"]
}
`
}
