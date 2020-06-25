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

output "mutation_output" {
  value = graphql_mutation.basic_mutation
}