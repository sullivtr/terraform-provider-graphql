package graphql

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

const (
	read_data_response = `{"data": {"todo": {"id": "1", "text": "something todo", "userId": "900"}}}`
	query_url          = "http://mock-gql-url.io"

	data_source_config = `
	data "graphql_query" "basic_query" {
		query_variables = {}
		query     =  file("../e2e/test_basic/queries/readQuery")
	}
`
	resource_config_create = `
	resource "graphql_mutation" "basic_mutation" {
		mutation_variables = {
			"text" = "something todo"
			"userId" = "900"
		}
		delete_mutation_variables = {
			"testvar1" = "testval1"
		}
		read_query_variables = {}
		create_mutation = file("../e2e/test_basic/queries/createMutation")
		update_mutation = file("../e2e/test_basic/queries/updateMutation")
		delete_mutation = file("../e2e/test_basic/queries/deleteMutation")
		read_query      = file("../e2e/test_basic/queries/readQuery")

		compute_mutation_keys = {
			"id" = "todo.id"
		}
	}
`
	resource_config_update = `
	resource "graphql_mutation" "basic_mutation" {
		mutation_variables = {
			"text" = "something else"
			"userId" = "500"
		}
		delete_mutation_variables = {
			"testvar1" = "testval1"
		}
		read_query_variables = {}
		create_mutation = file("../e2e/test_basic/queries/createMutation")
		update_mutation = file("../e2e/test_basic/queries/updateMutation")
		delete_mutation = file("../e2e/test_basic/queries/deleteMutation")
		read_query      = file("../e2e/test_basic/queries/readQuery")

		compute_mutation_keys = {
			"id" = "todo.id"
		}
	}
`
)

func mockGqlServerResponse(req *http.Request) (*http.Response, error) {
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)

	if strings.Contains(reqBody, "findTodos") {
		return httpmock.NewStringResponse(200, read_data_response), nil
	}

	return httpmock.NewStringResponse(200, ""), nil
}
