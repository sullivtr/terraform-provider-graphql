# <provider> GraphQL

This plugin provides a powerful way to automate GraphQL API resources using terraform.

## Example Usage

```hcl
provider "graphql" {
  url = "https://your-graphql-server-url"
  headers = {
    "header1" = "header1-value"
    "header2" = "header2-value"
    ...
  }
}
```

## Full documentation: 
This provider's extensive docmentation can be found here: [https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html](https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html)