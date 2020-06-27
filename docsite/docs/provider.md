---
id: provider
title: Provider Setup
---

>This provider provides a powerful way to automate GraphQL api resources using terraform.

## Provider Installation

The latest release can be downloaded from the graphql provider [release page](https://github.com/sullivtr/terraform-provider-graphql/releases/latest).
- The easiest way to install third-party terraform plugins like this on is to place the dowloaded binary in `~/.terraform.d/plugins/`. 
  
## Provider Configuration

This provider has only two configuration inputs: `url` & `headers`

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

## Inputs

### url
  - **Type**: `string`
  - **Description**: `The graphql api url to the provider will use to make requests`

### headers
  - **Type**: `map(string)`
  - **Desciption**: `Any http headers that the graphql api requires. (eg; Authentication; x-api-key; etc)`