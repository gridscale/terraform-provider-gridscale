package gridscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGridscaleK8sBasic(t *testing.T) {
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceGridscaleServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleK8sConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.gridscale_k8s.test", "name", name),
					resource.TestCheckResourceAttrSet(
						"data.gridscale_k8s.test", "id"),
					resource.TestCheckResourceAttr(
						"data.gridscale_k8s.test", "k8s_private_network_uuid", "f5d1b4e1-4f3b-4f6b-8e1e-3e6b4e1f3b4f"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleK8sConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "test" {
    name   = "%s"
	k8s_private_network_uuid = "f5d1b4e1-4f3b-4f6b-8e1e-3e6b4e1f3b4f"
}

data "gridscale_k8s" "test" {
    resource_id = gridscale_k8s.test.id
}`, name)
}
