# terraform-provider-graphql

## Synopsis

A [Terraform](http://terraform.io) plugin to manage graphql queries and mutations.

## Example:
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
data "graphql_query" "queryexample" {
  variables = {
    ID = "nfddksajf3948290dsa!f"
  }

  query = <<EOF
{
  secretByName(secretType: ENCRYPTED_TEXT, name:"secret-name") {
    name,
    id
  }
}
EOF          
}
```
#### Graphql Resource
```
resource "graphql_mutation" "mutationexample" {
  createMutationVariables = {
    "secretType" = "secret_type"
    "secret" = {
      "thing1": "thing2"
    }
  }

  # if update, create, and read variables are omitted, they will fall back to the required createMutationVariables
  updateMutationVariables = {
    "secretType" = "secret_type"
    "secret" = {
      "thing1": "thing4"
    }
  }
  deleteMutationVariables = {
    "secretId" = "123456"
  }
  readQueryVariables = {
    "secretName" = "secret_name"
  }
  # Reference files instead of inline queries to keep tf files clean. See examplquery for an example of a query file
  createMutation = "${path.module}/queries/createMutation"
  updateMutation = "${path.module}/queries/updateMutation"
  deleteMutation = "${path.module}/queries/createMudeleteMutationtation"
  readQuery      = "${path.module}/queries/readQuery"
}
```
## Data Sources

### graphql_query
#### Argument Reference
- variables (required): a map(string) of any variables that will be used in the query
- query (required): the graphql query. See [example query](./examplequery) for an example of what this looks like.
#### Outputs
- queryResponse: The resulting response body of the graphql query

## Resources

### graphql_mutation
#### Argument Reference
- createMutationVariables (Required): a map(string) of any variables that will be used in the mutation
- updateMutationVariables (Optional | falls back to createMutationVariables if omitted): a map(string) of any variables that will be used in the mutation
- deleteMutationVariables (Optional | falls back to createMutationVariables if omitted): a map(string) of any variables that will be used in the mutation
- readQueryVariables (Optional | falls back to createMutationVariables if omitted): a map(string) of any variables that will be used in the reader query

- createMutation: (Required) the graphql mutation to be used for the create operation  
- updateMutation: (Required) the graphql mutation to be used for the update operation 
- deleteMutation: (Required) the graphql mutation to be used for the delete operation 
- readQuery:      (Required) the graphql mutation to be used for the read operation

#### Outputs
- queryResponse: The resulting response body of the graphql query


# License

Apache2 - See the included LICENSE file for more details.

