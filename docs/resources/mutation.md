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
* `mutation_variables` - (Required) A map of any variables that will be used in your create & update mutation. Each variable's value is interpreted as JSON when possible. In order for the provider to keep track of state with respect to inputs provided in `mutation_variables`, you should ensure that your `read_query` includes the property that holds the value for your mutation variable. The keys do not need to be the same, as the provider with automatically map the keys based on value. See `query_response_input_key_map` description for details on how this works. 
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
* `force_replace` - A bool to determine if the resource should always be replaced (deleted and recreated) during update lifecycle hooks. Default: false
* `enable_remote_state_verification` - A pre v2.4.0 backward-compatibility flag. Set to `false` to disable resource remote state verification during reads.


## Attribute Reference

* `query_response` - A computed json encoded http response object received from the query.
    - To use properties from this response, leverage Terraform's built in [jsondecode](https://www.terraform.io/docs/configuration/functions/jsondecode.html) function.

* `query_response_input_key_map` - A computed map between the values represented by mutation_variable inputs and query response object keys. The purpose of this comnputed map is for the provider to keep track of drift between a resource's server-state and its state representation in terraform. Its value is calculated by creating a map of key-value pairs such that the key is the key for a particular input of the `mutation_variables`, and the value is the key that represents the corresponding value in the `query_response`. For example, consider we have a query response with the following structure: 
  ```javascript
  "data": {
    "id": "123456"
    "user": {
      "id": "654321"
    } 
  }
  ```

  Consider there is a mutation variable input of `userID = "654321"`. The computed `query_response_input_key_map` would be: 
  ```javascript
  {
    "userID": "data.user.id"
  }
  ```
  Internally, the provider uses this map to check the state of given input properties against their value returned from server state via `query_response` during a read lifecycle. 

* `computed_update_operation_variables` - A computed map that combines any computed variables with the `mutation_variables` input based on what is provided in the `compute_mutation_keys` input. This is also useful for outputing properties of the response object and using it on other resources.
  
* `computed_delete_operation_variables` - A computed map that combines any computed variables with the `delete_mutation_variables` input based on what is provided in the `compute_mutation_keys` input.

* `computed_read_operation_variables` - A computed map that combines any computed variables with the `read_query_variables` input based on what is provided in the `compute_mutation_keys` input. 


->**Note** For a full guide on using this provider, see the full documentation site located at [https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html](https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html)
