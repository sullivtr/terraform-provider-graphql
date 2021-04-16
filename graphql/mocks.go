package graphql

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

const (
	readDataResponse   = `{"data": {"todo": {"id": "1", "text": "something todo", "userId": "900"}}}`
	createDataResponse = `{"data": {"createTodo": {"id": "2"}}}`
	queryUrl           = "http://mock-gql-url.io"

	dataSourceConfig = `
	data "graphql_query" "basic_query" {
		query_variables = {}
		query     =  file("../testdata/readQuery")
	}
`
	resourceConfigCreate = `
	resource "graphql_mutation" "basic_mutation" {
		mutation_variables = {
			"text" = "something todo"
			"userId" = "900"
		}
		delete_mutation_variables = {
			"testvar1" = "testval1"
		}
		read_query_variables = {}
		create_mutation = file("../testdata/createMutation")
		update_mutation = file("../testdata/updateMutation")
		delete_mutation = file("../testdata/deleteMutation")
		read_query      = file("../testdata/readQuery")

		compute_mutation_keys = {
			"id" = "todo.id"
		}
	}
`

	resourceConfigComputeMutationKeysOnCreate = `
	resource "graphql_mutation" "basic_mutation" {
		mutation_variables = {
			"text" = "something todo"
			"userId" = "900"
		}
		delete_mutation_variables = {
			"testvar1" = "testval1"
		}
		read_query_variables = {
			"testvar1" = "testval1"
		}
		create_mutation = file("../testdata/createMutation")
		update_mutation = file("../testdata/updateMutation")
		delete_mutation = file("../testdata/deleteMutation")
		read_query      = file("../testdata/readQuery")

		compute_from_create = true

		compute_mutation_keys = {
			"id" = "createTodo.id"
		}
	}
`
	resourceConfigUpdate = `
	resource "graphql_mutation" "basic_mutation" {
		mutation_variables = {
			"text" = "something else"
			"userId" = "500"
		}
		delete_mutation_variables = {
			"testvar1" = "testval1"
		}
		read_query_variables = {}
		create_mutation = file("../testdata/createMutation")
		update_mutation = file("../testdata/updateMutation")
		delete_mutation = file("../testdata/deleteMutation")
		read_query      = file("../testdata/readQuery")

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
		return httpmock.NewStringResponse(200, readDataResponse), nil
	}

	return httpmock.NewStringResponse(200, ""), nil
}

func mockGqlServerResponseCreate(req *http.Request) (*http.Response, error) {
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	reqBody := string(reqBytes)

	if strings.Contains(reqBody, "createTodo") {
		return httpmock.NewStringResponse(200, createDataResponse), nil
	}

	return httpmock.NewStringResponse(200, ""), nil
}
