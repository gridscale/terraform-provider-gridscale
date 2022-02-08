package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

func TestAccResourceGridscaleLocation_Basic(t *testing.T) {
	var object gsclient.Location
	name := fmt.Sprintf("Test-TF-Location-%s", acctest.RandString(10))
	parentLocationUUID := "45ed677b-3702-4b36-be2a-a2eab9827950"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleLocationDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleLocationConfig_basic(name, parentLocationUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleLocationExists("gridscale_location.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_location.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_location.foo", "parent_location_uuid", parentLocationUUID),
					resource.TestCheckResourceAttr(
						"gridscale_location.foo", "cpunode_count", "10"),
					resource.TestCheckResourceAttr(
						"gridscale_location.foo", "product_no", "1500001"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleLocationExists(n string, object *gsclient.Location) resource.TestCheckFunc {
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

		foundObject, err := client.GetLocation(context.Background(), id)

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

func testAccCheckGridscaleLocationDestroyCheck(s *terraform.State) error {
	return nil
}

func testAccCheckResourceGridscaleLocationConfig_basic(name, parentLocationUUID string) string {
	return fmt.Sprintf(`
resource "gridscale_location" "foo" {
  name   = "%s"
  parent_location_uuid = "%s"
  product_no = 1500001
  cpunode_count = 10
}
`, name, parentLocationUUID)
}
