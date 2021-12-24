---
id: provider
title: Provider Setup
---

## Synopsis 
This plugin provides a powerful way to automate GraphQL API resources using terraform.

## Provider Installation

The latest release can be downloaded from the graphql provider [release page](https://github.com/sullivtr/terraform-provider-graphql/releases/latest).
- The easiest way to install third-party terraform plugins like this on is to place the dowloaded binary in `~/.terraform.d/plugins/`. 
  
## Provider Configuration

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
  oauth2_login_query_value_attribute = "loginAPI.accessToken"
}
```

## Inputs

### url
  - **Required**: `true`
  - **Type**: `string`
  - **Description**: The GraphQL API url that the provider will use to make requests.

### headers
  - **Required**: `false`
  - **Type**: `map(string)`
  - **Desciption**: Any http headers that the GraphQL API server requires (e.g. `Authentication`, `x-api-key`, etc.).

### oauth2_login_query
  - **Required**: `false`
  - **Type**: `string`
  - **Description**: The GraphQL query or mutation used to retrieve an OAuth 2.0 access token from the GraphQL server. It will be executed once during provider initialization. Be aware that renewal is not implemented, which may cause issue with short-lived access tokens. Note: you must also define `oauth2_login_query_variables` and `oauth2_login_query_value_attribute` when using `oauth2_login_query`.

### oauth2_login_query_variables
  - **Required**: `false`
  - **Type**: `map(string)`
  - **Desciption**: A map of any variables that will be used in your OAuth 2.0 login query. Each variable's value is interpreted as JSON when possible. Note: you must also define `oauth2_login_query` and `oauth2_login_query_value_attribute` when using `oauth2_login_query_variables`.

### oauth2_login_query_value_attribute
  - **Required**: `false`
  - **Type**: `string`
  - **Description**: The dot-separated path to the attribute containing the access token value that will be extracted from the OAuth 2.0 login query or mutation response `data` (e.g. `loginAPI.accessToken`). Note: you must also define `oauth2_login_query` and `oauth2_login_query_variables` when using `oauth2_login_query_value_attribute`.