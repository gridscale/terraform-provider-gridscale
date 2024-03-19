package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

func TestAccResourceGridscaleK8sBasic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("k8s-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleK8sConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleK8sConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "%s"
	release = "1.26"
	node_pool {
		name = "my_node_pool"
		node_count = 2
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
		rocket_storage = 90
	}
}
`, name)
}

func testAccCheckResourceGridscaleK8sConfigBasicUpdate() string {
	return `
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	release = "1.26"
	node_pool {
		name = "my_node_pool"
		node_count = 2
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
		rocket_storage = 90
	}
}
`
}
