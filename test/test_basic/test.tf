provider "graphql" {
  url = "http://localhost:8080/query"
  headers = {
    "x-api-key": "5555443399"
  }
}

# data "graphql_query" "basic_query" {
#   depends_on = [graphql_mutation.basic_mutation]
#   read_query_variables = {}
#   read_query     = file("./queries/readQuery")
# }

resource "graphql_mutation" "basic_mutation" {
  mutation_variables = {
    "text" = var.todo_text
    "userId" = var.todo_user_id
  }
  # if update, create, and read variables are omitted, they will fall back to the required create_mutation_variables
  read_query_variables = {}
  # Reference files instead of inline queries to keep tf files clean. See examplquery for an example of a query file
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