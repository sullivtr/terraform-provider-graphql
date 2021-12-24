# <provider> GraphQL

This plugin provides a powerful way to automate GraphQL API resources using terraform.

## Example Usage

This provider has several configuration inputs.

In most cases, only `url` & `headers` are used, e.g.:

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

In some advanced cases where the GraphQL server exposes an endpoint to perform OAuth 2.0 authentication, instead of using `headers`, you can use `oauth2_login_query`, `oauth2_login_query_variables` and `oauth2_login_query_value_attribute`, e.g.:

```hcl
provider "graphql" {
  url = "https://your-graphql-server-url"

  oauth2_login_query = "mutation loginAPI($apiKey: String!) {loginAPI(apiKey: $apiKey) {accessToken}}"
  oauth2_login_query_variables = {
    "apiKey" = "5555-44-33-99"
  }
  oauth2_login_query_value_attribute = "data.loginAPI.accessToken"
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html) (i.e., `alias` and `version`), the following arguments are supported in the GraphQL `provider` block:

* `url` - (Required) The GraphQL API url that the provider will use to make requests.

* `headers` - (Optional) Any http headers that the GraphQL API server requires (e.g. `Authentication`, `x-api-key`, etc.).

* `oauth2_login_query` - (Optional) The GraphQL query or mutation used to retrieve an OAuth 2.0 access token from the GraphQL server. It will be executed once during provider initialization. Be aware that renewal is not implemented, which may cause issue with short-lived access tokens. Note: you must also define `oauth2_login_query_variables` and `oauth2_login_query_value_attribute` when using `oauth2_login_query`.

* `oauth2_login_query_variables` - (Optional) A map of any variables that will be used in your OAuth 2.0 login query. Each variable's value is interpreted as JSON when possible. Note: you must also define `oauth2_login_query` and `oauth2_login_query_value_attribute` when using `oauth2_login_query_variables`.

* `oauth2_login_query_value_attribute` - (Optional) The dot-separated path to the attribute containing the access token value that will be extracted from the OAuth 2.0 login query or mutation response (it must start with `data.`, e.g. `data.loginAPI.accessToken`). Note: you must also define `oauth2_login_query` and `oauth2_login_query_variables` when using `oauth2_login_query_value_attribute`.

## Full documentation

This provider's extensive documentation can be found here: [https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html](https://sullivtr.github.io/terraform-provider-graphql/docs/provider.html)