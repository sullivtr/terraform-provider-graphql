package graphql

import (
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/jarcoal/httpmock"
)

const (
	data_response = `"{"data": {"todo": {"id": "1", "text": "something to do"}}}"`
	query_url     = "http://mock-gql-url.io"
)

func init() {
	os.Setenv("TF_GRAPHQL_URL", query_url)
	os.Setenv("TF_ACC", "1")
}

func TestAccGraphqlQuery_basic(t *testing.T) {
	config := `
		data "graphql_query" "basic_query" {
			query_variables = {} # this query does not take any variables as input
			query     = <<EOF
		query findTodos{
			todo {
				id
				text
			}
		}
		EOF
		}
	`

	// MOCK THE HTTP REQUEST
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", query_url, mockGqlServerResponse)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graphql_query.basic_query", "query_response", data_response),
				),
			},
		},
	})
}

func mockGqlServerResponse(req *http.Request) (*http.Response, error) {
	return httpmock.NewStringResponse(200, data_response), nil
}
