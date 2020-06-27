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
 
### delete_mutation_variables
  - **Type**: map(string)
  - **Desciption**: A computed map based on the result of what is provided in the `mutation_keys` input. 
  
  >NOTE: delete mutation variables are fully calculated at the moment. We expect to support merged delete mutation variables in a future release (similar to computed_update_operation_variables).

### computed_update_operation_variables (This is where the magic happens)
  - **Type**: map(string)
  - **Desciption**: A computed map that combines any computed variables with the `mutation_variables` input based on what is provided in the `mutation_keys` input. 
    - This is also useful for outputing properties of the response object and using it on other resources (if want to avoid that whole json decode thing mentioned above).

## Handling Update & Destroy

>This provider makes it simple to update and destroy api resources using computed properties (such as IDs). Since most delete and update mutations require a computed identifier for the object, this provider will keep track of the object's computed identifiers in state (or any other properties you ask it to keep track of).

### Defining computed variables:

As mentioned above, you define variables that _you_ want terraform to keep track of using the `mutation_keys` input. 

  **Example**: We have a read query that returns an object with this structure: 
  ```json
  { 
    data: { 
      todos: { 
        id
        text 
        } 
      } 
  }
  ```
  We can define our `mutation_keys` as:
  ```hcl
  mutation_keys = {
    "id" = "todo.id"
  }
  ```

  In this example, `todo.id` describes the property we want to collect from the response object. 
  >NOTE: Since it is idiomatic for GraphQL server responses to return objects with a "data" parent property, the "data" property is implcit. However, you can define the mutation key as "data.todo.id" if that makes you sleep better at night.

  To add to this, we can collect N... variables using `mutation_keys`. 

  For example, we can collect both the "id" and the "text" property off of a todo in the above example by defining `mutation_keys` as:
  ```hcl
    mutation_keys = {
      "id" = "todo.id"
      "my_todo_text" = "todo.text"
    }
  ```
  
### Using computed variables

The only thing that we have to do to make use of the properties collected from `mutation_keys` is to use those variables in your update and/or delete mutations.

  **Example**: We define a delete mutation that looks like this: 
  ```
  mutation deleteTodo($id: String!) {
    deleteTodo(input: $id) {
      id
    }
  }
  ```

  Since we told `mutation_keys` to collect the `id` property, and we defined it as `id` in the `mutation_keys` map, the delete mutation will automatically utilize the value returned from `todo.id` (which is collected during the read_query execution after a create or update execution). You could similary pass in a variable called `my_todo_text` to the mutation.

  This resource outputs `computed_update_operation_variables` and `delete_mutation_variables`, so you can always verify that they are reading values that you expect.

  The principles outlined above apply the same way to the `update_mutation`. If you need to utilize computed values in your update mutation, define them in your `mutation_keys` input. 


## Full lifecyle graphql_mutation examples

### Create Mutation Example
```hcl
mutation createUser($firstName: String!, $lastName: String!, $email: String!) {
  createUser(userInput: {
    givenName: $firstName,
    surname: $lastName,
    email: $email
  }) {
    id
  }
}
```

### Update Mutation Example
```hcl
mutation updateUser($userID: String!, $firstName: String!, $lastName: String!, $email: String!) {
  updateUser(userInput: {
    id: $userID,
    givenName: $firstName,
    surname: $lastName,
    email: $email
  }) {
    id,
    givenName,
    surname,
    email
  }
}
```

### Delete Mutation Example
```hcl
mutation deleteUser($userID: String!) {
  deleteUser(userID: $userID) {
    id
  }
}
```

### Read Query Example
```hcl
mutation getUserByEmail($email: String!) {
  deleteUser(userInput: {
    email: $email
  }) {
    id,
    givenName,
    surname,
    email
  }
}
```


