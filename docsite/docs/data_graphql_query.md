---
id: data_graphql_query
title: graphql_query
---


## Usage Overview
```hcl
data "graphql_query" "basic_query" {
  read_query_variables = {}
  read_query     = file("./path/to/query/file")
}
```

## Inputs

### read_query_variables
  - **Required**: true
  - **Type**: map(string)
  - **Description**: A map of any variables that will be used in your query

### read_query
  - **Required**: true
  - **Type**: string (multi-line)
  - **Desciption**: The graphql query. See basic example below for what that looks like.

## Outputs

### query_response
  - **Type**: string (json encoded http response)
  - **Desciption**: A computed json encoded http response object received from the query.
    - To use properties from this response, leverage Terraform's built in [jsondecode](https://www.terraform.io/docs/configuration/functions/jsondecode.html) function.

## Basic Example

Just like graphql on its own, this data source takes in the query varibles, and the query itself:
```hcl
data "graphql_query" "basic_query" {
  read_query_variables = {} # this query does not take any variables as input
  read_query     = <<EOF
query findTodos{
  todo {
    id
    text
    done
    user {
      name
    }
  }
}
EOF
}
```

The query itself can be referenced in-line, as shown above, or it can be referenced from a file using terraform's built in [file](https://www.terraform.io/docs/configuration/functions/file.html) function.

## Advanced Example

Just like graphql on its own, this data source takes in the query varibles, and the query itself:
```hcl
data "graphql_query" "advanced_query" {

  read_query_variables = {
      "name"  = "Jimmy Dean"
      "email" = "jimmydean@jdthesausageman.com"
  }

  read_query  = file("${path.module}/queries/readQuery")
}
```

With the above data source, you would create a file in the module path at `./queries/readQuery`:
```javascript
query getUser($name: String!, $email: String!) {
  user(input: {name: $name, email: $email}) {
    id
    name
    email
    phone
  }
}
```

As you can see above, if a query requires a user object to be pass as a parameter, you can build the object inline on the query and fill in the properties using variables. 

> NOTE: This provider does not currently support usage of complex objects as variables. read_query_variables must be a map of string.

The `query_response` output would be a json encoded object with this structure: 

```json
{
    "data": {
      "user": {
        "id": "XXXXXX",
        "name": "Jimmy Dean",
        "email": "jimmydean@jdthesausageman.com",
        "phone": "1234567890",
      }   
    }
}
```