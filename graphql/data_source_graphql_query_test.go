package graphql

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/jarcoal/httpmock"
)

func init() {
	os.Setenv("TF_GRAPHQL_URL", query_url)
	os.Setenv("TF_ACC", "1")
}

func TestAccGraphqlQuery_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", query_url, mockGqlServerResponse)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: data_source_config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graphql_query.basic_query", "query_response", read_data_response),
				),
			},
		},
	})
}
