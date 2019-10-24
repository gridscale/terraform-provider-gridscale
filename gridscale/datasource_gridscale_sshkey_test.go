package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccdataSourceGridscaleSSHKey_basic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleSshkeyDestroyCheck,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDataSourceSSHKeyConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gridscale_sshkey.foo", "id"),
					resource.TestCheckResourceAttr("data.gridscale_sshkey.foo", "name", name),
				),
			},
		},
	})

}

func testAccCheckDataSourceSSHKeyConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "gridscale_sshkey" "foo" {
  name   = "%s"
  sshkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKea3u6cuJ/2ZoMA4fpnXRK8ZIZWQz8ddXJv+iul9gTAc4fbm30IjZNnBBxiFOETc5ev1mcxvi6XvW99gLmxJAGwUrHylxYODXl1fLhc2G5czwQS9Qk57ED+IYb7AGOWPxGYeDaDka6gxJal/aaUx0C42fQErpUiJj2mJlF8yUOqyygtQOZhT2XUBU5UBZd50r8die8oRgdKJrbcn48q1Eu60vpx4S4JgH+krrHoXuCRydQ31KfOXmD8Y3/oGlZQ40luhfnj6g1jpm6PIQEBehGyZl6Dyh0MeeJsePWAGmXMEA33FcDkUiQPLoaalr4QQZdAUS74/irf+mgRcSRPvL root@475d4232363a"
}

data "gridscale_sshkey" "foo" {
	resource_id   = "${gridscale_sshkey.foo.id}"
}`, name)
}
