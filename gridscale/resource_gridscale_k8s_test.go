package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"testing"
)

func TestAccResourceGridscaleK8s_Basic(t *testing.T) {
	var object gsclient.PaaSService
	name := fmt.Sprintf("k8s-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscalePaaSDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleK8sConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", name),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleK8sConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "%s"
	release = "1.19"
	node_pool {
		name = "my_node_pool"
		node_count = 2
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
	}
}
`, name)
}

func testAccCheckResourceGridscaleK8sConfig_basic_update() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	release = "1.19"
	node_pool {
		name = "my_node_pool"
		node_count = 2
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
	}
}
`)
}
