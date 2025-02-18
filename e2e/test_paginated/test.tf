terraform {
  required_providers {
    graphql = {
      source  = "gqlprovidertf.com/examplecorp/graphql"
      version = "2.0.0"
    }
  }
}

provider "graphql" {
  url = "http://localhost:8080/query"
  headers = {
    "x-api-key": "5555443399"
  }
}

resource "graphql_mutation" "basic_mutation" {
  compute_from_create = var.compute_from_create
  mutation_variables = {
    "text" = var.todo_text
    "userId" = var.todo_user_id
    "list" = "[\"this\",\"that\"]"
  }
  
  delete_mutation_variables = {
    "testvar1" = "testval2"
  }
  read_query_variables = {
    "testvar1" = "testval2"
  }
  create_mutation = file("../../testdata/createMutation")
  update_mutation = file("../../testdata/updateMutation")
  delete_mutation = file("../../testdata/deleteMutation")
  read_query      = file("../../testdata/readQuery")

  compute_mutation_keys = var.compute_mutation_keys
}

resource "graphql_mutation" "basic_mutation_2" {
  compute_from_create = var.compute_from_create
  mutation_variables = {
    "text" = var.todo_text
    "userId" = var.todo_user_id_2
    "list" = "[\"this\",\"those\"]"
  }
  
  delete_mutation_variables = {
    "testvar1" = "testval2"
  }
  read_query_variables = {
    "testvar1" = "testval2"
  }
  create_mutation = file("../../testdata/createMutation")
  update_mutation = file("../../testdata/updateMutation")
  delete_mutation = file("../../testdata/deleteMutation")
  read_query      = file("../../testdata/readQuery")

  compute_mutation_keys = var.compute_mutation_keys
}

data "graphql_query" "basic_query" {
  depends_on = [graphql_mutation.basic_mutation]
  query = file("../../testdata/readQuery")
  query_variables = {}
  paginated = true
}

output "query_output" {
  value = data.graphql_query.basic_query
}

output "computed_read_variables" {
  value = graphql_mutation.basic_mutation.computed_read_operation_variables
}

output "computed_delete_variables" {
  value = graphql_mutation.basic_mutation.computed_delete_operation_variables
}