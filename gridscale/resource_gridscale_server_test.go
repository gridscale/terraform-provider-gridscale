package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gridscale/gsclient-go"
)

func TestAccDataSourceGridscaleServer_Basic(t *testing.T) {
	var object gsclient.Server
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceGridscaleServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleServerConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleServerExists("gridscale_server.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "memory", "1"),
				),
			},
			{
				Config: testAccCheckDataSourceGridscaleServerConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleServerExists("gridscale_server.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "cores", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "power", "true"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleServerExists(n string, object *gsclient.Server) resource.TestCheckFunc {
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
		foundObject, err := client.GetServer(id)
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

func testAccCheckDataSourceGridscaleServerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_server" {
			continue
		}
		_, err := client.GetServer(rs.Primary.ID)
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

func testAccCheckDataSourceGridscaleServerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv4" "ipserver" {
	name = "test-ip"
}

resource "gridscale_storage" "boot" {
  name   = "bootstorage"
  capacity = 1
}

resource "gridscale_storage" "additional" {
  name   = "additionalstorage"
  capacity = 1
}


resource "gridscale_server" "foo" {
  name   = "%s"
  cores = 1
  memory = 1
  ipv4 = "${gridscale_ipv4.ipserver.id}"
  storage {
    object_uuid = "${gridscale_storage.boot.id}"
	bootdevice = true
  }
  storage {
    object_uuid = "${gridscale_storage.additional.id}"
  }
}
`, name)
}

func testAccCheckDataSourceGridscaleServerConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_server" "foo" {
  name   = "newname"
  cores = 2
  memory = 2
  power = true
}
`)
}
