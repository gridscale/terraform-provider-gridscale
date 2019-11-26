package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"testing"
)

func TestAccResourceGridscalePaaS_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscalePaaSConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_paas.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas.foopaas", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscalePaaSConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_paas.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_paas.foopaas", "name", "newname"),
				),
			},
			{
				Config: testAccCheckResourceGridscalePaaSConfig_forcenew_update(),
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
		foundObject, err := client.GetPaaSService(emptyCtx, id)
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

		_, err := client.GetPaaSService(emptyCtx, rs.Primary.ID)
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

func testAccCheckResourceGridscalePaaSConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foopaas" {
  name = "%s"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
}
`, name)
}

func testAccCheckResourceGridscalePaaSConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foopaas" {
  name = "newname"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
  resource_limit {
	resource = "cores"
	limit = 16
  }
}
`)
}

func testAccCheckResourceGridscalePaaSConfig_forcenew_update() string {
	return fmt.Sprintf(`
resource "gridscale_paas" "foopaas" {
  name = "newname"
  service_template_uuid = "136c1446-13e0-4734-bdb6-ab0a15c1d680"
  resource_limit {
	resource = "cores"
	limit = 16
  }
}
`)
}
