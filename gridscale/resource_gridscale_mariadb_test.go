package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleMariaDB_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("postgres-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMariaDBConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mariadb.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mariadb.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMariaDBConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_mariadb.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_mariadb.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMariaDBConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_mariadb" "test" {
	name = "%s"
	release = "10.5"
	performance_class = "standard"
}
`, name)
}

func testAccCheckResourceGridscaleMariaDBConfig_basic_update() string {
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
