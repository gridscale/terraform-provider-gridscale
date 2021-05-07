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

func TestAccResourceGridscaleMarketplaceApplicationImport_Basic(t *testing.T) {
	var object gsclient.MarketplaceApplication
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleMarketplaceApplicationImportDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleMarketplaceApplicationImportConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleMarketplaceApplicationImportExists("gridscale_marketplace_application_import.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_marketplace_application_import.foo", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleMarketplaceApplicationImportConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleMarketplaceApplicationImportExists("gridscale_marketplace_application_import.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_marketplace_application_import.foo", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleMarketplaceApplicationImportExists(n string, object *gsclient.MarketplaceApplication) resource.TestCheckFunc {
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

		foundObject, err := client.GetMarketplaceApplication(context.Background(), id)

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

func testAccCheckGridscaleMarketplaceApplicationImportDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_marketplace_application_import" {
			continue
		}

		_, err := client.GetMarketplaceApplication(context.Background(), rs.Primary.ID)
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

func testAccCheckResourceGridscaleMarketplaceApplicationImportConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_marketplace_application" "foo" {
	name = "%s"
	object_storage_path = "s3://testsnapshot/test.gz"
	category = "CMS"
	setup_cores = 1
	setup_memory = 1
	setup_storage_capacity = 1
}

resource "gridscale_marketplace_application_import" "foo" {
	import_unique_hash = gridscale_marketplace_application.foo.unique_hash
}
`, name)
}

func testAccCheckResourceGridscaleMarketplaceApplicationImportConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_marketplace_application" "foonew" {
  	name   = "newname"
	object_storage_path = "s3://testsnapshot/test.gz"
	category = "CMS"
	setup_cores = 2
	setup_memory = 4
	setup_storage_capacity = 5
}
resource "gridscale_marketplace_application_import" "foo" {
	import_unique_hash = gridscale_marketplace_application.foonew.unique_hash
}
`)
}
