# <resource name> Mutation

This resource provides everything you need to create, read, update, and delete an api resource using GraphQL. 

## Example Usage

```hcl
resource "graphql_mutation" "basic_mutation" {
  mutation_variables = {
    "name" = "Jimmy Dean"
    "email" = "thewurst@jimmydean.com"
    "phone" = "\"1234567890\"" // Interpret as string
    "age" = "25" // This is interpreted as an integer
  }
  read_query_variables = {
    "email" = "thewurst@jimmydean.com"
  }

  compute_mutation_keys = {
    "id" = "user.id"
  }

  create_mutation = file("./queries/createMutation")
  update_mutation = file("./queries/updateMutation")
  delete_mutation = file("./queries/deleteMutation")
  read_query      = file("./queries/readQuery")
}
```

## Argument Reference
* `mutation_variables` - (Required) A map of any variables that will be used in your create & update mutation. Each variable's value is interpreted as JSON when possible.
  >NOTE: If a mutation variable is a number that must be interpreted as a string, it should be wrapped in quotations. For example `"marVar" = "\"123\""`.

* `read_query_variables` - (Optional) A map of any variables that will be used in the read query for the resource's lifecycle. Each variable's value is interpreted as JSON when possible.
   >NOTE: If a query variable is a number that must be interpreted as a string, it should be wrapped in quotations. For example `"marVar" = "\"123\""`.

* `delete_mutation_variables` - (Optional) A map of any variables that will be used in the delete mutation for the resource's lifecycle (This is automatically combined with any computed variables). Each variable's value is interpreted as JSON when possible.
  >NOTE: delete_mutation_variables are merged with any variables that are computed based on the compute_mutation_keys input. The result is the computed_delete_operation_variables output (similar to computed_update_operation_variables). If a delete mutation variable is a number that must be interpreted as a string, it should be wrapped in quotations. For example `"marVar" = "\"123\""`.

* `create_mutation` - (Required) A GraphQL mutation that will be used to create the api resource.
   
* `update_mutation` - (Required) A GraphQL mutation that will be used to update the api resource.
  
* `delete_mutation` - (Required) A GraphQL mutation that will be used to delete the api resource.

* `read_query` - (Required) A GraphQL query that will be used to query the api resource after it has been created.

* `compute_mutation_keys` - (Required) A map representing the hierarchy of your response object leading to the object properties that will be used during a terraform destroy & update operation.
* `compute_from_create` - A bool to determine if computed keys should be computed based off of the response from the create request, or the read request. Default: false


## Attribute Reference

* `query_response` - A computed json encoded http response object received from the query.
    - To use properties from this response, leverage Terraform's built in [jsondecode](https://www.terraform.io/docs/configuration/functions/jsondecode.html) function.

* `computed_update_operation_variables` - A computed map that combines any computed variables with the `mutation_variables` input based on what is provided in the `compute_mutation_keys` input. This is also useful for outputing properties of the response object and using it on other resources.
  
* `computed_delete_operation_variables` - A computed map that combines any computed variables with the `delete_mutation_variables` input based on what is provided in the `compute_mutation_keys` input.

* `computed_read_operation_variables` - A computed map that combines any computed variables with the `read_query_variables` input based on what is provided in the `compute_mutation_keys` input. 


->**Note** For a full guide on using this provider, see the full documentation site located at [https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html](https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html)
