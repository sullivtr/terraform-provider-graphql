# <resource name> Query

This data source is appropriate when you need to read/create any resources without managing its full lifecycle as you would with the `graphql_mutation` resource. 

## Example Usage

```hcl
data "graphql_query" "basic_query" {
  query_variables = {}
  query     = file("./path/to/query/file")
}
```

## Argument Reference

* `query_variables` - (Required) A map of any variables that will be used in your query

* `query` - (Required) The graphql query. (See basic example below for what that looks like.)

## Attribute Reference

* `query_response` - A computed json encoded http response object received from the query.
    - To use properties from this response, leverage Terraform's built in [jsondecode](https://www.terraform.io/docs/configuration/functions/jsondecode.html) function.


->**Note** For a full guide on using this provider, see the full documentation site located at [https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html](https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html)
