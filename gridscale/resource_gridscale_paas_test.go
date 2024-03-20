package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"testing"
)

func TestAccResourceGridscalePaaSBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscalePaaSConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_paas.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas.foopaas", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscalePaaSConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_paas.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas.foopaas", "name", "newname"),
				),
			},
			{
				Config: testAccCheckResourceGridscalePaaSConfigTMPUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_paas.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas.foopaas", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscalePaaSExists(n string, object *gsclient.PaaSService) resource.TestCheckFunc {
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
		foundObject, err := client.GetPaaSService(context.Background(), id)
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

func testAccCheckResourceGridscalePaaSDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_paas" {
			continue
		}

		_, err := client.GetPaaSService(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscalePaaSConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foopaas" {
  name = "%s"
  service_template_uuid = "d7a5e8ec-fa78-4d1b-86f9-febe3e16e398"
}
`, name)
}

func testAccCheckResourceGridscalePaaSConfigBasicUpdate() string {
	return `
resource "gridscale_paas" "foopaas" {
  name = "newname"
  service_template_uuid = "d7a5e8ec-fa78-4d1b-86f9-febe3e16e398"
  resource_limit {
	resource = "cores"
	limit = 16
  }
  parameter {
    param = "mysql_max_connections"
    value = "2000"
    type = "float"
  }
  parameter {
    param = "mysql_default_time_zone"
    value = "UTC"
    type = "string"
  }
}
`
}

// TO DO: update `service_template_uuid` when the backend enables the option to
// update `service_template_uuid`.
func testAccCheckResourceGridscalePaaSConfigTMPUpdate() string {
	return `
resource "gridscale_paas" "foopaas" {
  name = "newname"
  service_template_uuid = "d7a5e8ec-fa78-4d1b-86f9-febe3e16e398"
  resource_limit {
	resource = "cores"
	limit = 16
  }
  parameter {
    param = "mysql_max_connections"
    value = "2000"
    type = "float"
  }
  parameter {
    param = "mysql_default_time_zone"
    value = "UTC"
    type = "string"
  }
}
`
}
