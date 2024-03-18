package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGridscaleBucket_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleBucketConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_bucket.foo", "access_key"),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_bucket.foo", "secret_key"),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_bucket.foo", "bucket_name"),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_bucket.foo", "s3_host"),
					resource.TestCheckResourceAttrSet(
						"gridscale_object_storage_bucket.foo", "loc_constrain"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleBucketConfig_basic() string {
	return `
resource "gridscale_object_storage_accesskey" "test" {
   timeouts {
      create="10m"
  }
}

resource "gridscale_object_storage_bucket" "foo" {
   access_key = gridscale_object_storage_accesskey.test.access_key
   secret_key = gridscale_object_storage_accesskey.test.secret_key
   bucket_name = "myterraformbucket"
}
`
}
