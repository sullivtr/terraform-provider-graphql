package graphql

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"graphql": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProvider_rateLimitValidation(t *testing.T) {
	cases := []struct {
		name      string
		config    map[string]interface{}
		expectErr bool
	}{
		{name: "defaults", config: map[string]interface{}{}},
		{name: "zero rate disables limiting", config: map[string]interface{}{"rate_limit_per_second": 0}},
		{name: "fractional rate", config: map[string]interface{}{"rate_limit_per_second": 2.5}},
		{name: "minimum burst", config: map[string]interface{}{"rate_limit_burst": 1}},
		{name: "negative rate", config: map[string]interface{}{"rate_limit_per_second": -1}, expectErr: true},
		{name: "zero burst", config: map[string]interface{}{"rate_limit_burst": 0}, expectErr: true},
		{name: "negative burst", config: map[string]interface{}{"rate_limit_burst": -1}, expectErr: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{"url": "http://localhost:8080/query"}
			for k, v := range tc.config {
				raw[k] = v
			}

			diags := Provider().Validate(terraform.NewResourceConfigRaw(raw))
			if diags.HasError() != tc.expectErr {
				t.Fatalf("expectErr = %t, got diagnostics: %v", tc.expectErr, diags)
			}
		})
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TF_GRAPHQL_URL"); v == "" {
		t.Fatal("TF_GRAPHQL_URL must be set for acceptance tests")
	}
}
