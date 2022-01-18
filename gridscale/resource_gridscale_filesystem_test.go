package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleFilesystem_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("TEST-Filesystem-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleFilesystemConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_filesystem.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_filesystem.test", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_filesystem.test", "allowed_ip_ranges.0", "192.14.2.2"),
					resource.TestCheckResourceAttr(
						"gridscale_filesystem.test", "allowed_ip_ranges.1", "192.168.0.0/16"),
					resource.TestCheckResourceAttrSet(
						"gridscale_filesystem.test", "root_squash"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleFilesystemConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_filesystem.test", &object),
					resource.TestCheckResourceAttr(
						"gridscale_filesystem.test", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_filesystem.test", "allowed_ip_ranges.0", "192.14.15.15"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleFilesystemConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_filesystem" "test" {
	name = "%s"
	release = "1"
	performance_class = "standard"
	root_squash = true
	allowed_ip_ranges = ["192.14.2.2", "192.168.0.0/16"]
}
`, name)
}

func testAccCheckResourceGridscaleFilesystemConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_filesystem" "test" {
	name = "newname"
	root_squash = false
	release = "1"
	performance_class = "standard"
	allowed_ip_ranges = ["192.14.15.15"]
	labels = ["test"]
}
`)
}
