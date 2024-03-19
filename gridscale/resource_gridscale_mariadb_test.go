package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMariaDBBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("postgres-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMariaDBConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mariadb.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mariadb.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMariaDBConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mariadb.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mariadb.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMariaDBConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_mariadb" "test" {
	name = "%s"
	release = "10.5"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMariaDBConfigBasicUpdate() string {
	return `
resource "gridscale_mariadb" "test" {
	name = "newname"
	release = "10.5"
	performance_class = "standard"
	max_core_count = 20
	labels = ["test"]
	mariadb_query_cache_limit = "2M"
	mariadb_default_time_zone = "Europe/Berlin"
	mariadb_sql_mode = "NO_AUTO_CREATE_USER,ERROR_FOR_DIVISION_BY_ZERO"
	mariadb_server_id = 2
	mariadb_binlog_format = "STATEMENT"
}
`
}
