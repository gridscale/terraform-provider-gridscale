package gridscale

import (
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"gridscale": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderInitialization(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GRIDSCALE_UUID"); v == "" {
		t.Fatal("GRIDSCALE_UUID must be set for acceptance tests")
	}

	if v := os.Getenv("GRIDSCALE_TOKEN"); v == "" {
		t.Fatal("GRIDSCALE_TOKEN must be set for acceptance tests")
	}
}

func TestConvertStrToHeaderMap(t *testing.T) {
	type testCase struct {
		InputStr       string
		ExpectedOutput map[string]string
	}
	testCases := []testCase{
		{
			InputStr:       "",
			ExpectedOutput: make(map[string]string),
		},
		{
			InputStr:       "header",
			ExpectedOutput: make(map[string]string),
		},
		{
			InputStr: "header1:value1,header2:value2",
			ExpectedOutput: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
		},
	}
	for _, tCase := range testCases {
		result := convertStrToHeaderMap(tCase.InputStr)
		if !reflect.DeepEqual(result, tCase.ExpectedOutput) {
			t.Errorf("Output: %v, Expected: %v", result, tCase.ExpectedOutput)
		}
	}
}
