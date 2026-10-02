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
    "x-api-key" : "5555443399"
    # Opts requests into the test server's arrival-time recording so this
    # fixture can measure the client-side limiter without affecting the other
    # fixtures.
    "x-e2e-rate-limit" : "1"
  }

  rate_limit_per_second = var.rate_limit_per_second
  rate_limit_burst      = var.rate_limit_burst
}

resource "graphql_mutation" "rate_limited_mutation" {
  mutation_variables = {
    "text"   = "Here is something todo"
    "userId" = "T5577006791947779410"
    "list"   = "[\"this\",\"that\"]"
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

  compute_mutation_keys = {
    "id" = "todo.id"
  }
}

data "graphql_query" "rate_limited_query" {
  depends_on      = [graphql_mutation.rate_limited_mutation]
  query           = file("../../testdata/readQuery")
  query_variables = {}
}

output "query_output" {
  value = data.graphql_query.rate_limited_query
}
