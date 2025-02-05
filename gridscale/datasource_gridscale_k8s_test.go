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
					resource.TestCheckResourceAttrSet(
						"data.gridscale_k8s.test", "k8s_private_network_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.gridscale_k8s.test", "kubeconfig"),
					resource.TestCheckResourceAttrSet(
						"data.gridscale_k8s.test", "labels"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleK8sConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "test" {
    name   = "%s"
	release = "1.30" # instead, gsk_version can be set.

	node_pool {
		name = "pool-0"
		node_count = 2
		cores = 2
		memory = 4
		storage = 30
		storage_type = "storage_insane"
	}
}

data "gridscale_k8s" "test" {
    resource_id = gridscale_k8s.test.id
}`, name)
}
