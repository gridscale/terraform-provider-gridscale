package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gridscale/gsclient-go"
)

func TestAccDataSourceGridscaleSshkey_Basic(t *testing.T) {
	var object gsclient.Sshkey
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSshkeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleSshkeyConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSshkeyExists("gridscale_sshkey.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_sshkey.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_sshkey.foo", "sshkey", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKea3u6cuJ/2ZoMA4fpnXRK8ZIZWQz8ddXJv+iul9gTAc4fbm30IjZNnBBxiFOETc5ev1mcxvi6XvW99gLmxJAGwUrHylxYODXl1fLhc2G5czwQS9Qk57ED+IYb7AGOWPxGYeDaDka6gxJal/aaUx0C42fQErpUiJj2mJlF8yUOqyygtQOZhT2XUBU5UBZd50r8die8oRgdKJrbcn48q1Eu60vpx4S4JgH+krrHoXuCRydQ31KfOXmD8Y3/oGlZQ40luhfnj6g1jpm6PIQEBehGyZl6Dyh0MeeJsePWAGmXMEA33FcDkUiQPLoaalr4QQZdAUS74/irf+mgRcSRPvL root@475d4232363a"),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSshkeyConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSshkeyExists("gridscale_sshkey.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_sshkey.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_sshkey.foo", "sshkey", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvOYJ3xXtPXaPOacFQ97+nGq5QDkl17/JeTaY36RLPKgYBt2Z5YSPSROdzh/5GgZ0p6E3W84gKNaedUo3v+zvgmdGZeDFk+cxlC0HtXwQN87GQRtYTMsucbI6OJT7p4qntl70MIBzvIrmheGZqXnpeRxA7PjVcjkA3nxps3XJsuMDd0Ft0Ue3j0lmOno779mfgg34VeTgE2GZlH31gFqxWz3fXUgaZoLdO7HbLKu4ybfFWdCzqBt4B8RG9xMq0220gJR6ZwAaiMc1CGIknK7C6EKeCx9LOWDjCaHg6pA2iPAb/PoxDuiqbUIzfRmkgMf0lYmrf0kqx529ALm92ulSx root@33c294c5235e"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleSshkeyExists(n string, object *gsclient.Sshkey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}
		client := testAccProvider.Meta().(*gsclient.Client)
		id := rs.Primary.ID
		foundObject, err := client.GetSshkey(id)
		if err != nil {
			return err
		}
		if foundObject.Properties.ObjectUUID != id {
			return fmt.Errorf("Object not found")
		}
		*object = foundObject
		return nil
	}
}

func testAccCheckDataSourceGridscaleSshkeyDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_sshkey" {
			continue
		}
		_, err := client.GetSshkey(rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("Unable to fetching object %s", rs.Primary.ID)
			}
		} else {
			return fmt.Errorf("Object %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDataSourceGridscaleSshkeyConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_sshkey" "foo" {
  name   = "%s"
  sshkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKea3u6cuJ/2ZoMA4fpnXRK8ZIZWQz8ddXJv+iul9gTAc4fbm30IjZNnBBxiFOETc5ev1mcxvi6XvW99gLmxJAGwUrHylxYODXl1fLhc2G5czwQS9Qk57ED+IYb7AGOWPxGYeDaDka6gxJal/aaUx0C42fQErpUiJj2mJlF8yUOqyygtQOZhT2XUBU5UBZd50r8die8oRgdKJrbcn48q1Eu60vpx4S4JgH+krrHoXuCRydQ31KfOXmD8Y3/oGlZQ40luhfnj6g1jpm6PIQEBehGyZl6Dyh0MeeJsePWAGmXMEA33FcDkUiQPLoaalr4QQZdAUS74/irf+mgRcSRPvL root@475d4232363a"
}
`, name)
}

func testAccCheckDataSourceGridscaleSshkeyConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_sshkey" "foo" {
  name   = "newname"
  sshkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvOYJ3xXtPXaPOacFQ97+nGq5QDkl17/JeTaY36RLPKgYBt2Z5YSPSROdzh/5GgZ0p6E3W84gKNaedUo3v+zvgmdGZeDFk+cxlC0HtXwQN87GQRtYTMsucbI6OJT7p4qntl70MIBzvIrmheGZqXnpeRxA7PjVcjkA3nxps3XJsuMDd0Ft0Ue3j0lmOno779mfgg34VeTgE2GZlH31gFqxWz3fXUgaZoLdO7HbLKu4ybfFWdCzqBt4B8RG9xMq0220gJR6ZwAaiMc1CGIknK7C6EKeCx9LOWDjCaHg6pA2iPAb/PoxDuiqbUIzfRmkgMf0lYmrf0kqx529ALm92ulSx root@33c294c5235e"
}
`)
}
