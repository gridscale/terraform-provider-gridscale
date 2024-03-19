package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscalePostgresBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("postgres-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscalePostgresConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_postgresql.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_postgresql.test", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscalePostgresConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_postgresql.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_postgresql.test", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscalePostgresConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_postgresql" "test" {
	name = "%s"
	release = "13"
	performance_class = "standard"
	pgaudit_log_bucket = "foo"
	pgaudit_log_server_url = "https://gos3.io"
	pgaudit_log_access_key = "TESTINGPOSTGRESQLRESOURCEACCESSKEY"
	pgaudit_log_secret_key = "testing"
	pgaudit_log_rotation_frequency = 30
}
`, name)
}

func testAccCheckResourceGridscalePostgresConfigBasicUpdate() string {
	return `
resource "gridscale_postgresql" "test" {
	name = "newname"
	release = "13"
	performance_class = "standard"
	max_core_count = 20
	labels = ["test"]
	pgaudit_log_bucket = "foo"
	pgaudit_log_server_url = "https://gos3.io"
	pgaudit_log_access_key = "TESTINGPOSTGRESQLRESOURCEACCESSKEY"
	pgaudit_log_secret_key = "testing"
	pgaudit_log_rotation_frequency = 25
}
`
}
