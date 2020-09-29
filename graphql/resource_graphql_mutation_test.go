package graphql

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jarcoal/httpmock"
)

func init() {
	os.Setenv("TF_GRAPHQL_URL", queryUrl)
	os.Setenv("TF_ACC", "1")
}

func TestAccGraphqlMutation_full(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", queryUrl, mockGqlServerResponse)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGraphqlMutationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.id", "1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.text", "something todo"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.userId", "900"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.id", "1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.testvar1", "testval1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "query_response", readDataResponse),
				),
			},
			{
				Config: resourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.id", "1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.text", "something else"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.userId", "500"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.id", "1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.testvar1", "testval1"),
				),
			}},
	})
}

func TestAccGraphqlMutation_computefromcreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", queryUrl, mockGqlServerResponseCreate)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGraphqlMutationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigComputeMutationKeysOnCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_update_operation_variables.id", "2"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.id", "2"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_read_operation_variables.id", "2"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_read_operation_variables.testvar1", "testval1"),
					resource.TestCheckResourceAttr("graphql_mutation.basic_mutation", "computed_delete_operation_variables.testvar1", "testval1"),
				),
			},
		},
	})
}

func testAccGraphqlMutationResourceDestroy(s *terraform.State) error {
	resource.TestCheckNoResourceAttr("graphql_mutation.basic_mutation", "id")
	return nil
}
