# terraform-provider-graphql
### NOTE: 
This provider plugin is currently in under active development, and is not yet production ready. Contributions and feedback are welcome. 

## Synopsis

A [Terraform](https://terraform.io) plugin to manage [GraphQL](https://graphql.org/) queries and mutations. 
  

## Example: 
Open the [./test/test_basic]("./test/test_basic") directory for a basic example usage of this provider.

#### Provider setup:
```
provider "graphql" {
  url = "https://my-graphql-service-url.io"
  headers = {
    "x-api-key": "4324nsfkdsanj32k!!4FakeApiKey8873"
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
}
```
#### Data source
```
data "graphql_query" "basic_query" {
  read_query_variables = {}
  read_query     = "${file("./queries/readQuery")}"
}
```
#### Graphql Resource
```
resource "graphql_mutation" "basic_mutation" {
  # variables
  create_mutation_variables = {
    "text" = "Here is something todo"
    "userId" = "98"
  }
  update_mutation_variables = {}
  delete_mutation_variables = {}
  read_query_variables = {}

  # Reference files instead of inline queries to keep tf files clean. See examplquery for an example of a query file
  create_mutation = "${file("./queries/createMutation")}"
  update_mutation = "${file("./queries/updateMutation")}"
  delete_mutation = "${file("./queries/deleteMutation")}"
  read_query      = "${file("./queries/readQuery")}"
}
```
## Data Sources

### graphql_query
#### Argument Reference
- read_query_variables (required): a map(string) of any variables that will be used in the query
- read_query (required): the graphql query. See [example query](./examplequery) for an example of what this looks like.
#### Outputs
- query_response: The resulting response body of the graphql query

## Resources

### graphql_mutation
#### Argument Reference
- create_mutation_variables (Required): a map(string) of any variables that will be used in the mutation
- update_mutation_variables (Optional): a map(string) of any variables that will be used in the mutation
- delete_mutation_variables (Optional): a map(string) of any variables that will be used in the mutation
- read_query_variables (Optional): a map(string) of any variables that will be used in the reader query

- create_mutation: (Required) the graphql mutation to be used for the create operation  
- update_mutation: (Required) the graphql mutation to be used for the update operation 
- delete_mutation: (Required) the graphql mutation to be used for the delete operation 
- read_query:      (Required) the graphql mutation to be used for the read operation

#### Outputs
- query_response: The resulting response body of the graphql query

## Testing
- First, in the root of the project run `make build && make copyplugins`
- Go to [./test/gql-server]("./test/gql-server") and run `go run server.go`.
- Open the [./test/test_basic]("./test/test_basic") directory. 
  - run `tf init && tf apply`
  
 note: development of automated testing is in progress.
 
# License

Apache2 - See the included LICENSE file for more details.

