package gridscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGridscaleBucketBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleBucketConfigBasic(),
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

func testAccCheckResourceGridscaleBucketConfigBasic() string {
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

func TestAccResourceGridscaleBucketLifecycleRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleBucketConfigWithLifecycleRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.expiration_days", "30"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleBucketConfigWithUpdatedLifecycleRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.enabled", "false"),
					resource.TestCheckResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule.0.expiration_days", "60"),
				),
			},
			{
				Config: testAccCheckResourceGridscaleBucketConfigWithoutLifecycleRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"gridscale_object_storage_bucket.foo", "lifecycle_rule"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleBucketConfigWithLifecycleRules() string {
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

   lifecycle_rule {
     id = "rule1"
     enabled = true
     expiration_days = 30
   }
}
`
}

func testAccCheckResourceGridscaleBucketConfigWithUpdatedLifecycleRules() string {
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

   lifecycle_rule {
     id = "rule1"
     enabled = false
     expiration_days = 60
   }
}
`
}

func testAccCheckResourceGridscaleBucketConfigWithoutLifecycleRules() string {
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
