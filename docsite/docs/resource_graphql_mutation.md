---
id: resource_graphql_mutation
title: graphql_mutation
---

## Usage Overview
```hcl
resource "graphql_mutation" "basic_mutation" {
  # variables
  mutation_variables = {
    "name" = "Jimmy Dean"
    "email" = "thewurst@jimmydean.com"
    "phone" = "1234567890"
  }
  read_query_variables = {
    "email" = "thewurst@jimmydean.com"
  }

  mutation_keys = {
    "id" = "user.id"
  }

  # Reference files instead of inline queries to keep tf files clean. See examplquery for an example of a query file
  create_mutation = file("./queries/createMutation")
  update_mutation = file("./queries/updateMutation")
  delete_mutation = file("./queries/deleteMutation")
  read_query      = file("./queries/readQuery")
}
```

## Inputs

### mutation_variables
  - **Required**: true
  - **Type**: map(string)
  - **Description**: A map of any variables that will be used in your create & update mutation. 
  >NOTE: Any variables that are not actually used in mutations will be ignored. 

### read_query_variables
  - **Required**: false
  - **Type**: map(string)
  - **Description**: A map of any variables that will be used in the read query for teh resources lifecycle. 

### create_mutation
  - **Required**: true
  - **Type**: string (multi-line)
  - **Description**: A GraphQL mutation that will be used to create the api resource. (See basic example below for an example of what this looks like)
   
### update_mutation
   - **Required**: true
  - **Type**: string (multi-line)
  - **Description**: A GraphQL mutation that will be used to update the api resource. (See basic example below for an example of what this looks like)
  
### delete_mutation
  - **Required**: true
  - **Type**: string (multi-line)
  - **Description**: A GraphQL mutation that will be used to delete the api resource. (See basic example below for an example of what this looks like)

### read_query
  - **Required**: true
  - **Type**: string (multi-line)
  - **Description**: A GraphQL query that will be used to query the api resource after it has been created. (See basic example below for an example of what this looks like)

### mutation_keys
  - **Required**: true
  - **Type**: map(string)
  - **Description**: A map representing the hierarchy of your response object leading to the object properties that will be used during a terraform destroy & update operation.
  
  **See the "Handling Update & Destroy" section below** for an overview of the `mutation_keys` input usage. 

## Outputs

### query_response
  - **Type**: string (json encoded http response)
  - **Desciption**: A computed json encoded http response object received from the query.
    - To use properties from this response, leverage Terraform's built in [jsondecode](https://www.terraform.io/docs/configuration/functions/jsondecode.html) function.
 
### computed_update_operation_variables (This is where the magic happens)
  - **Type**: map(string)
  - **Desciption**: A computed map that combines any computed variables with your `mutation_variables` input based on what is provided in the `mutation_keys` input. 
    - This is also useful for outputing properties of your response object and using it on other resources (if you want to avoid that whole json decode thing mentioned above).

## Handling Update & Destroy

>This provider makes it simple to update and destroy api resources using computed properties (such as IDs). Since most delete and update mutations require a computed identifier for the object, this provider will keep track of the object's computed identifiers in state (or any other properties you ask it to keep track of).



