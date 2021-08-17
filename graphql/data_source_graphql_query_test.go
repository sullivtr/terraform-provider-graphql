package graphql

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
)

func init() {
	os.Setenv("TF_GRAPHQL_URL", queryUrl)
	os.Setenv("TF_ACC", "1")
}

func TestAccGraphqlQuery_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", queryUrl, mockGqlServerResponse)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graphql_query.basic_query", "query_response", readDataResponse),
				),
			},
		},
	})
}

func TestAccGraphqlQuery_basic_expectError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", queryUrl, mockGqlServerResponseError)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSourceConfig,
				ExpectError: regexp.MustCompile("bad things happened"),
			},
		},
	})
}
