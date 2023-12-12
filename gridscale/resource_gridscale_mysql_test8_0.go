package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMySQL_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("MySQL-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMySQLConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mysql.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mysql.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMySQLConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mysql.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mysql.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMySQLConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_mysql" "test" {
	name = "%s"
	release = "5.7"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMySQLConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_mysql" "test" {
	name = "newname"
	release = "5.7"
	performance_class = "standard"
	max_core_count = 20
	mysql_max_connections = 2000
	labels = ["test"]
}
`)
}
