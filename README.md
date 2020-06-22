# terraform-provider-graphql

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
  read_query     = file("./queries/readQuery")
}
```
#### Graphql Resource
```
resource "graphql_mutation" "basic_mutation" {
  # variables
  mutation_variables = {
    "text" = "Here is something todo"
    "userId" = "98"
  }
  read_query_variables = {}

  mutation_keys = ["todo.id"]

  # Reference files instead of inline queries to keep tf files clean. See examplquery for an example of a query file
  create_mutation = file("./queries/createMutation")
  update_mutation = file("./queries/updateMutation")
  delete_mutation = file("./queries/deleteMutation")
  read_query      = file("./queries/readQuery")
}
```
## Data Sources

### graphql_query
#### Argument Reference
- `read_query_variables` (required): a map(string) of any variables that will be used in the query
- `read_query` (required): the graphql query. See [example query](./examplequery) for an example of what this looks like.
#### Outputs
- `query_response`: The resulting response body of the graphql query

## Resources

### graphql_mutation
#### Argument Reference
- `mutation_variables` (Required): a map(string) of any variables that will be used in the mutation
- `read_query_variables` (Optional): a map(string) of any variables that will be used in the reader query

- `create_mutation`: (Required): the graphql mutation to be used for the create operation  
- `update_mutation`: (Required): the graphql mutation to be used for the update operation 
- `delete_mutation`: (Required): the graphql mutation to be used for the delete operation 
- `read_query`:      (Required): the graphql mutation to be used for the read operation

- `mutation_keys` (Required): list of string representing the hierarchy of your response object leading to the key(s) that will be used during a terraform destroy operation.
  **See "Handling tf update & destroy operations" below in the outputs section.**

#### Outputs
- `query_response`: The resulting response body of the graphql query

#### Handling tf update & destroy operations:

**Delete Operations**
- `delete_mutation_variables`: The delete mutation variables are computed based on the `mutation_keys` variable.
  Example: Your read query returns an object that has this structure: 
  ```
  { 
    data: { 
      todos: { 
        id
        text 
        } 
      } 
  }
  ```
  such that `id` is the property you use to delete your object during a destroy event. 

  Your delete mutation would look something like this: 
  ```
  mutation deleteTodo($id: String!) {
    deleteTodo(input: $id) {
      id
    }
  }
  ```
  You would set the `mutation_keys` variable on the resource as `["todo.id"]`. NOTE: Since the standard for GraphQL is to return objects with the `data` parent object, the root `data` key is implied. However, you can use `["data.todo.id"]` if that makes you sleep better at night. 

  If your delete events require more than one key/variable, you can pass unlimited maps to the `mutation_keys` list. For example, for two keys you would use this: `["todo.id", "todo.text"]`

  The result is a map(string) that is used as the variables object in your delete mutation. Example:
  ```
    "delete_mutation_variables" = {
      "id" = "T8674665223082153551"
      "text" = "Here is something todo"
    }
  ```
  Your delete mutation variables are automatically computed this way. 

**Update operations**:
- Any computed keys (such as an object's ID) will be computed for your update mutations. The keys are combined with your `mutation_variables` during an update event. 
  For example:

  Your `mutation_variables` are set as:
  ```
  { 
    "text" = "Here is my todo text"
    "userId" = "12"
  }
  ```

  The Todo.Id property is calculated from the `mutation_keys` variables, and merged with your `mutation_variables` to looks like this:
  ```
  { 
    "id"     = "computed_id"
    "text"   = "Here is my todo text"
    "userId" = "12"
  }
  ```

  See the basic [test project]("./test/test_basic") for examples.

## Testing
- First, in the root of the project run `make build && make copyplugins`
- Go to [./test/gql-server]("./test/gql-server") and run `go run server.go`.
- Open the [./test/test_basic]("./test/test_basic") directory. 
  - run `tf init && tf apply`
  
 note: development of automated testing is in progress.
 
# License

Apache2 - See the included LICENSE file for more details.


