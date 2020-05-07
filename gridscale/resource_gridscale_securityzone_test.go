package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccDataSourceGridscaleSecurityZone_Basic(t *testing.T) {
	var object gsclient.PaaSSecurityZone
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleSecurityZoneDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleSecurityZoneConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSecurityZoneExists("gridscale_paas_securityzone.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas_securityzone.foo", "name", name),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleSecurityZoneConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleSecurityZoneExists("gridscale_paas_securityzone.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas_securityzone.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleSecurityZoneExists(n string, object *gsclient.PaaSSecurityZone) resource.TestCheckFunc {
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
		foundObject, err := client.GetPaaSSecurityZone(context.Background(), id)
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

func testAccCheckDataSourceGridscaleSecurityZoneDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_securityzone" {
			continue
		}

		_, err := client.GetPaaSSecurityZone(context.Background(), rs.Primary.ID)
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

func testAccCheckDataSourceGridscaleSecurityZoneConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas_securityzone" "foo" {
  name = "%s"
}
`, name)
}

func testAccCheckDataSourceGridscaleSecurityZoneConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_paas_securityzone" "foo" {
  name = "newname"
}
`)
}
