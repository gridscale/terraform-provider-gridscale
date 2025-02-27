package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"
)

const (
	originalGSKVersion = "1.28.15-gs1"
	updatedGSKVersion  = "1.29.13-gs0"
	originalGSKRelease = "1.28"
	updatedGSKRelease  = "1.29"
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
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "30"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", originalGSKVersion),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigBasicUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", originalGSKVersion),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "30"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigVersionUpgrade(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", updatedGSKVersion),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "30"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigNodePoolSpecsUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", updatedGSKVersion),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "4"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "50"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigNodeCountIncrease(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", updatedGSKVersion),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "4"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "50"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigNodeCountDecrease(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "name", "newname"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "gsk_version", updatedGSKVersion),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.cores", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.memory", "4"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage", "50"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.storage_type", "storage_insane"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas", "node_pool.0.rocket_storage", "10"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigBasicRelease(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas2", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "release", originalGSKRelease),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.storage", "30"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.storage_type", "storage_insane"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleK8sConfigReleaseUpgrade(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscalePaaSExists("gridscale_k8s.foopaas2", &object),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "release", updatedGSKRelease),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.name", "my_node_pool"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.memory", "2"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.storage", "30"),
					resource.TestCheckResourceAttr(
						"gridscale_k8s.foopaas2", "node_pool.0.storage_type", "storage_insane"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleK8sConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "%s"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 1
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, name, originalGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigBasicUpdate() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 1
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, originalGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigVersionUpgrade() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 1
		cores = 1
		memory = 2
		storage = 30
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, updatedGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigNodePoolSpecsUpdate() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 1
		cores = 2
		memory = 4
		storage = 50
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, updatedGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigNodeCountIncrease() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 2
		cores = 2
		memory = 4
		storage = 50
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, updatedGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigNodeCountDecrease() string {
	return fmt.Sprintf(`
resource "gridscale_k8s" "foopaas" {
	name   = "newname"
	gsk_version = "%s"
	node_pool {
		name = "my_node_pool"
		node_count = 1
		cores = 2
		memory = 4
		storage = 50
		storage_type = "storage_insane"
		rocket_storage = 10
	}
}
`, updatedGSKVersion)
}

func testAccCheckResourceGridscaleK8sConfigBasicRelease() string {
	return fmt.Sprintf(`
	resource "gridscale_k8s" "foopaas2" {
		name   = "gsk-release-test"
		release = "%s"
		node_pool {
			name = "my_node_pool"
			node_count = 1
			cores = 1
			memory = 2
			storage = 30
			storage_type = "storage_insane"
		}
	}
	`, originalGSKRelease)
}

func testAccCheckResourceGridscaleK8sConfigReleaseUpgrade() string {
	return fmt.Sprintf(`
	resource "gridscale_k8s" "foopaas2" {
		name   = "gsk-release-test"
		release = "%s"
		node_pool {
			name = "my_node_pool"
			node_count = 1
			cores = 1
			memory = 2
			storage = 30
			storage_type = "storage_insane"
		}
	}
	`, updatedGSKRelease)
}
