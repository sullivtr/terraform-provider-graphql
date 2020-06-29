provider "graphql" {
  url = "http://localhost:8080/query"
  headers = {
    "x-api-key": "5555443399"
  }
}

resource "graphql_mutation" "basic_mutation" {
  mutation_variables = {
    "text" = var.todo_text
    "userId" = var.todo_user_id
  }
  read_query_variables = {}
  create_mutation = file("./queries/createMutation")
  update_mutation = file("./queries/updateMutation")
  delete_mutation = file("./queries/deleteMutation")
  read_query      = file("./queries/readQuery")

  mutation_keys = {
    "id" = "todo.id"
  }
}

data "graphql_query" "basic_query" {
  depends_on = [graphql_mutation.basic_mutation]
  query = file("./queries/readQuery")
  query_variables = {}
}


output "mutation_output" {
  value = graphql_mutation.basic_mutation
}

output "query_output" {
  value = data.graphql_query.basic_query
}