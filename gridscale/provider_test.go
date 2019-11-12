package gridscale

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"gridscale": testAccProvider,
	}
	serverPowerStateList = listServersPowerStatus{
		list: make(map[string]bool),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GRIDSCALE_UUID"); v == "" {
		t.Fatal("GRIDSCALE_UUID must be set for acceptance tests")
	}

	if v := os.Getenv("GRIDSCALE_TOKEN"); v == "" {
		t.Fatal("GRIDSCALE_TOKEN must be set for acceptance tests")
	}
}
