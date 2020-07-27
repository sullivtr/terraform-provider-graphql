package graphql

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
